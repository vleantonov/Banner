package mock

import (
	"banner/internal/domain"
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"slices"
)

type Repo struct {
	mock.Mock
}

func (r *Repo) GetActiveContentByTagFeatureID(ctx context.Context, tagId, featureId int) (*map[string]interface{}, error) {
	r.Called(tagId, featureId)

	if slices.Contains(TestBanner.Tags, tagId) && TestBanner.Feature == featureId {
		return &TestBanner.Content, nil
	}
	return nil, errors.New("mock error")
}

func (r *Repo) GetByFilter(ctx context.Context, f domain.FilterBanner) (*[]domain.Banner, error) {
	return &[]domain.Banner{
		*TestBanner,
	}, nil
}

func (r *Repo) Insert(ctx context.Context, b domain.Banner) (int, error) {
	return 0, nil
}
func (r *Repo) Update(ctx context.Context, b domain.UpdBanner) error {
	return nil
}
func (r *Repo) DeleteByID(ctx context.Context, id int) error {
	return nil
}
