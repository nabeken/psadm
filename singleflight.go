package psadm

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/singleflight"
)

// SingleflightClient is a duplicate function call suppression client.
type SingleflightClient struct {
	g      singleflight.Group
	client client
}

var _ client = &SingleflightClient{}

func (c *SingleflightClient) GetParameterWithDescription(ctx context.Context, key string) (*Parameter, error) {
	ck := buildCacheKey("GetParameterWithDescription", key)
	v, err, _ := c.g.Do(ck, func() (interface{}, error) {
		return c.client.GetParameterWithDescription(ctx, key)
	})
	return v.(*Parameter), err
}

func (c *SingleflightClient) GetParameter(ctx context.Context, key string) (string, error) {
	ck := buildCacheKey("GetParameter", key)

	v, err, _ := c.g.Do(ck, func() (interface{}, error) {
		return c.client.GetParameter(ctx, key)
	})

	return v.(string), err
}

func (c *SingleflightClient) GetParameterByTime(ctx context.Context, key string, at time.Time) (*Parameter, error) {
	ck := fmt.Sprintf("%s/%s/%d", "GetParameterByTime", key, at.Unix())

	v, err, _ := c.g.Do(ck, func() (interface{}, error) {
		return c.client.GetParameterByTime(ctx, key, at)
	})

	return v.(*Parameter), err
}

func (c *SingleflightClient) GetParametersByPath(ctx context.Context, pathPrefix string) ([]*Parameter, error) {
	ck := buildCacheKey("GetParametersByPath", pathPrefix)

	v, err, _ := c.g.Do(ck, func() (interface{}, error) {
		return c.client.GetParametersByPath(ctx, pathPrefix)
	})

	return v.([]*Parameter), err
}

// PutParameter forwards a call to the underlying client. It doesn't do any deduplication.
func (c *SingleflightClient) PutParameter(ctx context.Context, p *Parameter, overwrite bool) error {
	return c.client.PutParameter(ctx, p, overwrite)
}
