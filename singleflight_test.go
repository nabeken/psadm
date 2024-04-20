package psadm

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestSingleflightCient(t *testing.T) {
	assert := assert.New(t)
	mockctrl := gomock.NewController(t)

	mockSSM := NewMockssmClient(mockctrl)
	client := &Client{SSM: mockSSM}

	ch := make(chan struct{})
	mockSSM.EXPECT().GetParameter(gomock.Any(), gomock.Any()).
		// make sure the client only call the underlying client once
		Times(1).
		DoAndReturn(func(_ context.Context, _ *ssm.GetParameterInput, _ ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
			log.Print("ssm client is waiting for goroutines launched...")
			<-ch
			log.Print("ssm client is going to return a result")
			return &ssm.GetParameterOutput{
				Parameter: &types.Parameter{
					Value: aws.String("value"),
				},
			}, nil
		})

	c := cache.New(time.Minute, 10*time.Minute)
	sfc := client.SingleflightClientWithCache(c)

	// launch 10 goroutines
	const numG = 10
	var launched int

	var wg sync.WaitGroup
	wg.Add(numG)

	cond := sync.NewCond(&sync.Mutex{})
	for i := 0; i < numG; i++ {
		log.Print("Launching goroutine...")
		go func() {
			defer wg.Done()
			cond.L.Lock()
			launched++
			cond.L.Unlock()

			// let the main goroutine check launched again
			cond.Signal()

			actual, err := sfc.GetParameter(context.TODO(), "key")
			assert.NoError(err)
			assert.Equal("value", actual)
		}()
	}

	log.Print("Waiting for goroutines launched...")
	cond.L.Lock()
	for launched != numG {
		cond.Wait()
	}
	cond.L.Unlock()
	log.Print("Goroutines launched")

	close(ch)

	wg.Wait()

	cvalue, found := c.Get("GetParameter/key")
	assert.True(found)
	assert.Equal("value", cvalue)
}
