package ps

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type ssmClient interface {
	DescribeParameters(*ssm.DescribeParametersInput) (*ssm.DescribeParametersOutput, error)
	GetParameters(*ssm.GetParametersInput) (*ssm.GetParametersOutput, error)
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
	return err
}

// GetAllParameters gets all parameters having prefix.
// If decrypt is true, returned parameters will be decrypted.
func (c *Client) GetAllParameters(prefix string) ([]*Parameter, error) {
	input := &ssm.DescribeParametersInput{}
	var params []*Parameter

	for {
		desc, err := c.SSM.DescribeParameters(input)
		if err != nil {
			return nil, err
		}
		for _, p := range desc.Parameters {
			if prefix == "" || strings.HasPrefix(aws.StringValue(p.Name), prefix) {
				val, err := c.SSM.GetParameters(&ssm.GetParametersInput{
					Names:          []*string{p.Name},
					WithDecryption: aws.Bool(true),
				})
				if err != nil {
					return nil, err
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
		if desc.NextToken == nil {
			break
		}
		input.NextToken = desc.NextToken
	}

	return params, nil
}
