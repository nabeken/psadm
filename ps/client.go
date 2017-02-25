package ps

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type svcClient interface {
	DescribeParameters(*ssm.DescribeParametersInput) (*ssm.DescribeParametersOutput, error)
	GetParameters(*ssm.GetParametersInput) (*ssm.GetParametersOutput, error)
}

// Client wraps SSM client for psadm.
type Client struct {
	SSM svcClient
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
