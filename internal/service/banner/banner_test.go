package banner_test

import (
	"banner/internal/service/banner"
	"banner/internal/service/banner/mock"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBannerService_GetActiveByTagFeatureID_Valid(t *testing.T) {
	rep := &mock.Repo{}
	c := &mock.Cache{}

	c.On("GetByTagFeatureID", 1, 1).Return(&mock.TestBanner.Content, true)
	rep.On("GetActiveContentByTagFeatureID", 1, 1).Return(&mock.TestBanner.Content)

	s := banner.New(
		rep,
		c,
		&mock.Deleter{},
	)

	ctx := context.Background()

	b, err := s.GetActiveContentByTagFeatureID(ctx, 1, 1, false)

	c.AssertCalled(t, "GetByTagFeatureID", 1, 1)
	rep.AssertNotCalled(t, "GetActiveContentByTagFeatureID")
	require.NoError(t, err)
	require.NotNil(t, b)
	assert.Equal(t, mock.TestBanner.Content, *b)

	b, err = s.GetActiveContentByTagFeatureID(ctx, 1, 1, true)

	rep.AssertCalled(t, "GetActiveContentByTagFeatureID", 1, 1)
	require.NotNil(t, b)
	assert.Equal(t, mock.TestBanner.Content, *b)

}

func TestBannerService_GetActiveByTagFeatureID_Invalid(t *testing.T) {
	rep := &mock.Repo{}
	c := &mock.Cache{}

	c.On("GetByTagFeatureID", 10, 10).Return(nil, false)
	rep.On("GetActiveContentByTagFeatureID", 10, 10).Return(nil, false)

	s := banner.New(
		rep,
		c,
		&mock.Deleter{},
	)

	ctx := context.Background()

	b, _ := s.GetActiveContentByTagFeatureID(ctx, 10, 10, false)

	c.AssertCalled(t, "GetByTagFeatureID", 10, 10)
	rep.AssertCalled(t, "GetActiveContentByTagFeatureID", 10, 10)

	assert.Nil(t, b)
}
