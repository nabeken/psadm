//go:generate mockgen -source=client.go -package psadm -destination mock_client.go
package psadm

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
)

// ssmClient allows us to inject a fake API client for testing.
type ssmClient interface {
	DescribeParameters(context.Context, *ssm.DescribeParametersInput, ...func(*ssm.Options)) (*ssm.DescribeParametersOutput, error)
	GetParameter(context.Context, *ssm.GetParameterInput, ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
	GetParameterHistory(context.Context, *ssm.GetParameterHistoryInput, ...func(*ssm.Options)) (*ssm.GetParameterHistoryOutput, error)
	PutParameter(context.Context, *ssm.PutParameterInput, ...func(*ssm.Options)) (*ssm.PutParameterOutput, error)
}

// Parameter is the parameter exported by psadm.
// This should be sufficient for import and export.
type Parameter struct {
	Description string `yaml:"description"`
	KMSKeyID    string `yaml:"kmskeyid"`
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Value       string `yaml:"value"`
}

// client is an internal interface that can be chained with the standard client.
type client interface {
	GetParameterWithDescription(context.Context, string) (*Parameter, error)
	GetParameter(context.Context, string) (string, error)
	GetParameterByTime(context.Context, string, time.Time) (*Parameter, error)
	PutParameter(context.Context, *Parameter, bool) error
	GetParametersByPath(context.Context, string) ([]*Parameter, error)
}

// Client wraps the SSM client for psadm.
type Client struct {
	SSM ssmClient
}

// NewClient returns a psadm client.
func NewClient(cfg aws.Config) *Client {
	return &Client{
		SSM: ssm.NewFromConfig(cfg),
	}
}

// SingleflightClient returns a client with single flight caching.
func (c *Client) SingleflightClientWithCache(cache *cache.Cache) *SingleflightClient {
	return &SingleflightClient{
		client: c.CachedClient(cache),
	}
}

// CachedClient returns a client with caching.
func (c *Client) CachedClient(cache *cache.Cache) *CachedClient {
	return &CachedClient{
		cache:  cache,
		client: c,
	}
}

func (c *Client) GetParameterWithDescription(ctx context.Context, key string) (*Parameter, error) {
	keyName := string(types.ParametersFilterKeyName)
	desc, err := c.describeParameters(ctx, []types.ParameterStringFilter{
		{
			Key:    &keyName,
			Option: aws.String("Equals"),
			Values: []string{key},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(desc) == 0 {
		return nil, errors.Errorf("'%s' is not found.", key)
	}

	val, err := c.getParameter(ctx, key)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get parameter '%s'", key)
	}

	p := desc[0]
	return &Parameter{
		Description: aws.ToString(p.Description),
		KMSKeyID:    aws.ToString(p.KeyId),
		Name:        aws.ToString(p.Name),
		Type:        string(p.Type),
		Value:       aws.ToString(val.Parameter.Value),
	}, nil
}

// GetParameter returns the decrypted parameter.
func (c *Client) GetParameter(ctx context.Context, key string) (string, error) {
	resp, err := c.getParameter(ctx, key)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get parameter '%s'", key)
	}
	return aws.ToString(resp.Parameter.Value), nil
}

// GetParameterByTime returns the latest parameter.
func (c *Client) GetParameterByTime(ctx context.Context, key string, at time.Time) (*Parameter, error) {
	keyName := string(types.ParametersFilterKeyName)
	desc, err := c.describeParameters(ctx, []types.ParameterStringFilter{
		{
			Key:    &keyName,
			Option: aws.String("Equals"),
			Values: []string{key},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(desc) == 0 {
		return nil, errors.Errorf("'%s' is not found.", key)
	}

	latest := aws.ToTime(desc[0].LastModifiedDate)

	if latest.Before(at) {
		return c.GetParameterWithDescription(ctx, key)
	}

	// dig into history
	history, err := c.getParameterHistory(ctx, key)
	if err != nil {
		return nil, err
	}
	if len(history) == 0 {
		return nil, errors.Errorf("'%s' is not found.", key)
	}

	// history is sorted by LastModifiedDate in ascending order
	var p types.ParameterHistory
	for _, h := range history {
		if aws.ToTime(h.LastModifiedDate).After(at) {
			continue
		}
		p = h
	}

	return &Parameter{
		Description: aws.ToString(p.Description),
		KMSKeyID:    aws.ToString(p.KeyId),
		Name:        aws.ToString(p.Name),
		Type:        string(p.Type),
		Value:       aws.ToString(p.Value),
	}, nil
}

// PutParameter puts param into Parameter Store.
func (c *Client) PutParameter(ctx context.Context, param *Parameter, overwrite bool) error {
	input := &ssm.PutParameterInput{
		Name:      aws.String(param.Name),
		Type:      types.ParameterType(param.Type),
		Value:     aws.String(param.Value),
		Overwrite: aws.Bool(overwrite),
	}
	if param.Description != "" {
		input.Description = aws.String(param.Description)
	}
	if param.KMSKeyID != "" {
		input.KeyId = aws.String(param.KMSKeyID)
	}
	_, err := c.SSM.PutParameter(ctx, input)
	return errors.Wrap(err, "failed to put parameters")
}

func (c *Client) getParameter(ctx context.Context, key string) (*ssm.GetParameterOutput, error) {
	return c.SSM.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	})
}

// GetParametersByPath gets all parameters having given path prefix.
func (c *Client) GetParametersByPath(ctx context.Context, pathPrefix string) ([]*Parameter, error) {
	var filters []types.ParameterStringFilter

	keyName := string(types.ParametersFilterKeyName)

	if pathPrefix != "" {
		filters = []types.ParameterStringFilter{
			{
				Key:    &keyName,
				Option: aws.String("BeginsWith"),
				Values: []string{pathPrefix},
			},
		}
	}

	desc, err := c.describeParameters(ctx, filters)
	if err != nil {
		return nil, err
	}

	var params []*Parameter
	for _, p := range desc {
		val, err := c.getParameter(ctx, *p.Name)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get parameters")
		}

		params = append(params, &Parameter{
			Description: aws.ToString(p.Description),
			KMSKeyID:    aws.ToString(p.KeyId),
			Name:        aws.ToString(p.Name),
			Type:        string(p.Type),
			Value:       aws.ToString(val.Parameter.Value),
		})
	}

	return params, nil
}

func (c *Client) getParameterHistory(ctx context.Context, key string) ([]types.ParameterHistory, error) {
	input := &ssm.GetParameterHistoryInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	}

	var history []types.ParameterHistory
	for {
		resp, err := c.SSM.GetParameterHistory(ctx, input)
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

func (c *Client) describeParameters(ctx context.Context, filters []types.ParameterStringFilter) ([]types.ParameterMetadata, error) {
	input := &ssm.DescribeParametersInput{
		ParameterFilters: filters,
	}

	var params []types.ParameterMetadata
	for {
		desc, err := c.SSM.DescribeParameters(ctx, input)
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
