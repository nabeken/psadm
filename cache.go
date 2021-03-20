package psadm

import (
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

func (c *CachedClient) GetParameterWithDescription(key string) (*Parameter, error) {
	ck := buildCacheKey("GetParameterWithDescription", key)

	if v, found := c.cache.Get(ck); found {
		return v.(*Parameter), nil
	}

	param, err := c.client.GetParameterWithDescription(key)
	if err != nil {
		return nil, err
	}

	c.cache.Set(ck, param, cache.DefaultExpiration)
	return param, nil
}

func (c *CachedClient) GetParameter(key string) (string, error) {
	ck := buildCacheKey("GetParameter", key)

	if v, found := c.cache.Get(ck); found {
		return v.(string), nil
	}

	param, err := c.client.GetParameter(key)
	if err != nil {
		return "", err
	}

	c.cache.Set(ck, param, cache.DefaultExpiration)
	return param, nil
}

func (c *CachedClient) GetParameterByTime(key string, at time.Time) (*Parameter, error) {
	ck := fmt.Sprintf("%s/%s/%d", "GetParameterByTime", key, at.Unix())

	if v, found := c.cache.Get(ck); found {
		return v.(*Parameter), nil
	}

	param, err := c.client.GetParameterByTime(key, at)
	if err != nil {
		return nil, err
	}

	c.cache.Set(ck, param, cache.DefaultExpiration)
	return param, nil
}

func (c *CachedClient) GetParametersByPath(pathPrefix string) ([]*Parameter, error) {
	ck := buildCacheKey("GetParametersByPath", pathPrefix)

	if v, found := c.cache.Get(ck); found {
		return v.([]*Parameter), nil
	}

	params, err := c.client.GetParametersByPath(pathPrefix)
	if err != nil {
		return nil, err
	}

	c.cache.Set(ck, params, cache.DefaultExpiration)
	return params, nil
}

// PutParameter forwards a call to the underlying client. It doesn't do any caching.
func (c *CachedClient) PutParameter(p *Parameter, overrite bool) error {
	return c.client.PutParameter(p, overrite)
}

func buildCacheKey(prefix, key string) string {
	return fmt.Sprintf("%s/%s", prefix, key)
}
