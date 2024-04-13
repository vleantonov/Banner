package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"slices"
)

type Cache struct {
	mock.Mock
}

func (c *Cache) GetByTagFeatureID(ctx context.Context, tagID, featureID int) (*map[string]interface{}, bool) {
	c.Called(tagID, featureID)
	if slices.Contains(TestBanner.Tags, tagID) && TestBanner.Feature == featureID {
		return &TestBanner.Content, true
	}
	return nil, false
}

func (c *Cache) StoreByTagFeatureID(ctx context.Context, tagID, featureID int, content *map[string]interface{}) {
}
