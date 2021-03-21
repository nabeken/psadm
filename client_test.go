package psadm

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	assert := assert.New(t)
	mockctrl := gomock.NewController(t)

	t.Run("GetParameter", func(t *testing.T) {
		mockSSM := NewMockssmClient(mockctrl)
		mockSSM.EXPECT().
			GetParameter(&ssm.GetParameterInput{
				Name:           aws.String("key/1/2/3"),
				WithDecryption: aws.Bool(true),
			}).
			Return(&ssm.GetParameterOutput{
				Parameter: &ssm.Parameter{
					Value: aws.String("value"),
				},
			}, nil)

		client := &Client{SSM: mockSSM}
		v, err := client.GetParameter("key/1/2/3")
		assert.Equal("value", v)
		assert.NoError(err)
	})
}
