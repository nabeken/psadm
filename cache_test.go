package psadm

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestCachedClient(t *testing.T) {
	assert := assert.New(t)
	mockctrl := gomock.NewController(t)

	t.Run("GetParameter", func(t *testing.T) {
		mockSSM := NewMockssmClient(mockctrl)
		mockSSM.EXPECT().
			GetParameter(gomock.Any(), &ssm.GetParameterInput{
				Name:           aws.String("key/1/2/3"),
				WithDecryption: aws.Bool(true),
			}).
			Return(&ssm.GetParameterOutput{
				Parameter: &types.Parameter{
					Value: aws.String("value"),
				},
			}, nil)

		c := cache.New(time.Minute, 10*time.Minute)
		client := (&Client{SSM: mockSSM}).CachedClient(c)

		v, err := client.GetParameter(context.TODO(), "key/1/2/3")
		assert.Equal("value", v)
		assert.NoError(err)

		cvalue, found := c.Get("GetParameter/key/1/2/3")
		assert.True(found)
		assert.Equal("value", cvalue)
	})
}
