package ps

import (
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pkg/errors"
)

type ssmClient interface {
	DescribeParameters(*ssm.DescribeParametersInput) (*ssm.DescribeParametersOutput, error)
	GetParameters(*ssm.GetParametersInput) (*ssm.GetParametersOutput, error)
	GetParameterHistory(*ssm.GetParameterHistoryInput) (*ssm.GetParameterHistoryOutput, error)
	PutParameter(*ssm.PutParameterInput) (*ssm.PutParameterOutput, error)
}

// Client wraps SSM client for psadm.
type Client struct {
	SSM ssmClient
}

// NewClient returns an AWS wrapper client fr psadm.
func NewClient(sess *session.Session) *Client {
	return &Client{
		SSM: ssm.New(sess),
	}
}

// GetParameter returns the latest parameter.
func (c *Client) GetParameter(key string) (*Parameter, error) {
	desc, err := c.describeParameters([]*ssm.ParametersFilter{
		{
			Key:    aws.String(ssm.ParametersFilterKeyName),
			Values: []*string{aws.String(key)},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(desc) == 0 {
		return nil, errors.Errorf("'%s' is not found.", key)
	}

	val, err := c.SSM.GetParameters(&ssm.GetParametersInput{
		Names:          []*string{aws.String(key)},
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get parameters")
	}

	p := desc[0]
	return &Parameter{
		Description: aws.StringValue(p.Description),
		KMSKeyID:    aws.StringValue(p.KeyId),
		Name:        aws.StringValue(p.Name),
		Type:        aws.StringValue(p.Type),
		Value:       aws.StringValue(val.Parameters[0].Value),
	}, nil
}

// GetParameterByTime returns the latest parameter.
func (c *Client) GetParameterByTime(key string, at time.Time) (*Parameter, error) {
	desc, err := c.describeParameters([]*ssm.ParametersFilter{
		{
			Key:    aws.String(ssm.ParametersFilterKeyName),
			Values: []*string{aws.String(key)},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(desc) == 0 {
		return nil, errors.Errorf("'%s' is not found.", key)
	}

	latest := aws.TimeValue(desc[0].LastModifiedDate)

	if latest.Before(at) {
		return c.GetParameter(key)
	}

	// dig into history
	history, err := c.getParameterHistory(key)
	if err != nil {
		return nil, err
	}
	if len(history) == 0 {
		return nil, errors.Errorf("'%s' is not found.", key)
	}

	// history is sorted by LastModifiedDate in ascending order
	var p *ssm.ParameterHistory
	for _, h := range history {
		if aws.TimeValue(h.LastModifiedDate).After(at) {
			continue
		}
		p = h
	}

	if p == nil {
		return nil, errors.Errorf("'%s' is not found at give time.", key)
	}

	return &Parameter{
		Description: aws.StringValue(p.Description),
		KMSKeyID:    aws.StringValue(p.KeyId),
		Name:        aws.StringValue(p.Name),
		Type:        aws.StringValue(p.Type),
		Value:       aws.StringValue(p.Value),
	}, nil
}

// PutParameter puts param into Parameter Store.
func (c *Client) PutParameter(param *Parameter, overwrite bool) error {
	input := &ssm.PutParameterInput{
		Name:      aws.String(param.Name),
		Type:      aws.String(param.Type),
		Value:     aws.String(param.Value),
		Overwrite: aws.Bool(overwrite),
	}
	if param.Description != "" {
		input.Description = aws.String(param.Description)
	}
	if param.KMSKeyID != "" {
		input.KeyId = aws.String(param.KMSKeyID)
	}
	_, err := c.SSM.PutParameter(input)
	return errors.Wrap(err, "failed to put parameters")
}

// GetAllParameters gets all parameters having prefix.
func (c *Client) GetAllParameters(prefix string) ([]*Parameter, error) {
	desc, err := c.describeParameters(nil)
	if err != nil {
		return nil, err
	}

	var params []*Parameter
	for _, p := range desc {
		if prefix == "" || strings.HasPrefix(aws.StringValue(p.Name), prefix) {
			val, err := c.SSM.GetParameters(&ssm.GetParametersInput{
				Names:          []*string{p.Name},
				WithDecryption: aws.Bool(true),
			})
			if err != nil {
				return nil, errors.Wrap(err, "failed to get parameters")
			}

			params = append(params, &Parameter{
				Description: aws.StringValue(p.Description),
				KMSKeyID:    aws.StringValue(p.KeyId),
				Name:        aws.StringValue(p.Name),
				Type:        aws.StringValue(p.Type),
				Value:       aws.StringValue(val.Parameters[0].Value),
			})
		}
	}

	return params, nil
}

func (c *Client) getParameterHistory(key string) ([]*ssm.ParameterHistory, error) {
	input := &ssm.GetParameterHistoryInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	}

	var history []*ssm.ParameterHistory
	for {
		resp, err := c.SSM.GetParameterHistory(input)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get parameter history")
		}
		history = append(history, resp.Parameters...)

		if resp.NextToken == nil {
			break
		}
		input.NextToken = resp.NextToken
	}

	return history, nil
}

func (c *Client) describeParameters(filters []*ssm.ParametersFilter) ([]*ssm.ParameterMetadata, error) {
	input := &ssm.DescribeParametersInput{
		Filters: filters,
	}

	var params []*ssm.ParameterMetadata
	for {
		desc, err := c.SSM.DescribeParameters(input)
		if err != nil {
			return nil, errors.Wrap(err, "failed to describe parameters")
		}
		params = append(params, desc.Parameters...)
		if desc.NextToken == nil {
			break
		}
		input.NextToken = desc.NextToken
	}

	return params, nil
}
