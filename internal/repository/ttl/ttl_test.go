package ttl_test

import (
	"banner/internal/domain"
	"banner/internal/repository/ttl"
	"context"
	"github.com/ReneKroon/ttlcache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	testTagID1                  = 1
	testTagID2                  = 2
	testFeatureID               = 1
	testDuration  time.Duration = time.Duration(3) * time.Second
)

func TestBannerCache(t *testing.T) {

	testBanner := &domain.Banner{
		ID: 1,
		Content: map[string]interface{}{
			"BannerKey": "BannerValue",
		},
		Tags:    []int{testTagID1, testTagID2},
		Feature: testFeatureID,
	}

	c := ttlcache.NewCache()
	c.SetTTL(testDuration)
	bc := ttl.New(
		c,
	)
	ctx := context.Background()

	_, ok := bc.GetByTagFeatureID(ctx, testTagID1, testFeatureID)
	require.False(t, ok)

	bc.StoreByTagFeatureID(ctx, testTagID1, testFeatureID, &testBanner.Content)

	b, ok := bc.GetByTagFeatureID(ctx, testTagID1, testFeatureID)
	require.True(t, ok)
	assert.Equal(t, testBanner.Content, *b)

	<-time.After(testDuration)

	_, ok = bc.GetByTagFeatureID(ctx, testTagID1, testFeatureID)
	require.False(t, ok)

}
