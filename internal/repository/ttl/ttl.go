package ttl

import (
	"context"
	"fmt"
	"github.com/ReneKroon/ttlcache"
)

type BannerContent struct {
	c *ttlcache.Cache
}

func (b *BannerContent) GetByTagFeatureID(ctx context.Context, tagID, featureID int) (*map[string]interface{}, bool) {
	key := getStringKey(tagID, featureID)

	if contentValue, ok := b.c.Get(key); ok {
		return contentValue.(*map[string]interface{}), true
	}

	return nil, false
}

func (b *BannerContent) StoreByTagFeatureID(ctx context.Context, tagID, featureID int, content *map[string]interface{}) {
	key := getStringKey(tagID, featureID)
	b.c.Set(key, content)

	return
}

func New(c *ttlcache.Cache) *BannerContent {
	return &BannerContent{
		c: c,
	}
}

func getStringKey(tagID, featureID int) string {
	return fmt.Sprintf("%d:%d", tagID, featureID)
}
