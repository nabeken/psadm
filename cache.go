package psadm

import (
	"context"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

// CachedClient is a cache-aware psadmin client.
type CachedClient struct {
	cache  *cache.Cache
	client client
}

var _ client = &CachedClient{}

func (c *CachedClient) GetParameterWithDescription(ctx context.Context, key string) (*Parameter, error) {
	ck := buildCacheKey("GetParameterWithDescription", key)

	if v, found := c.cache.Get(ck); found {
		return v.(*Parameter), nil
	}

	param, err := c.client.GetParameterWithDescription(ctx, key)
	if err != nil {
		return nil, err
	}

	c.cache.Set(ck, param, cache.DefaultExpiration)
	return param, nil
}

func (c *CachedClient) GetParameter(ctx context.Context, key string) (string, error) {
	ck := buildCacheKey("GetParameter", key)

	if v, found := c.cache.Get(ck); found {
		return v.(string), nil
	}

	param, err := c.client.GetParameter(ctx, key)
	if err != nil {
		return "", err
	}

	c.cache.Set(ck, param, cache.DefaultExpiration)
	return param, nil
}

func (c *CachedClient) GetParameterByTime(ctx context.Context, key string, at time.Time) (*Parameter, error) {
	ck := fmt.Sprintf("%s/%s/%d", "GetParameterByTime", key, at.Unix())

	if v, found := c.cache.Get(ck); found {
		return v.(*Parameter), nil
	}

	param, err := c.client.GetParameterByTime(ctx, key, at)
	if err != nil {
		return nil, err
	}

	c.cache.Set(ck, param, cache.DefaultExpiration)
	return param, nil
}

func (c *CachedClient) GetParametersByPath(ctx context.Context, pathPrefix string) ([]*Parameter, error) {
	ck := buildCacheKey("GetParametersByPath", pathPrefix)

	if v, found := c.cache.Get(ck); found {
		return v.([]*Parameter), nil
	}

	params, err := c.client.GetParametersByPath(ctx, pathPrefix)
	if err != nil {
		return nil, err
	}

	c.cache.Set(ck, params, cache.DefaultExpiration)
	return params, nil
}

// PutParameter forwards a call to the underlying client. It doesn't do any caching.
func (c *CachedClient) PutParameter(ctx context.Context, p *Parameter, overwrite bool) error {
	return c.client.PutParameter(ctx, p, overwrite)
}

func buildCacheKey(prefix, key string) string {
	return fmt.Sprintf("%s/%s", prefix, key)
}
