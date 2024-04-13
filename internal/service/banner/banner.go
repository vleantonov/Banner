package banner

import (
	"banner/internal/domain"
	"context"
)

type Deleter interface {
	Delete(ctx context.Context, tagID, featureID *int) error
}

type Repository interface {
	GetActiveContentByTagFeatureID(ctx context.Context, tagId, featureId int) (*map[string]interface{}, error)
	GetByFilter(ctx context.Context, f domain.FilterBanner) (*[]domain.Banner, error)
	Insert(ctx context.Context, b domain.Banner) (int, error)
	Update(ctx context.Context, b domain.UpdBanner) error
	DeleteByID(ctx context.Context, id int) error
}

type Cache interface {
	GetByTagFeatureID(ctx context.Context, tagID, featureID int) (*map[string]interface{}, bool)
	StoreByTagFeatureID(ctx context.Context, tagID, featureID int, content *map[string]interface{})
}

type Service struct {
	r Repository
	c Cache
	d Deleter
}

func New(r Repository, c Cache, d Deleter) *Service {
	return &Service{
		r: r,
		c: c,
		d: d,
	}
}

func (b *Service) GetActiveContentByTagFeatureID(ctx context.Context, tagID, FeatureID int, useLast bool) (*map[string]interface{}, error) {
	if content, ok := b.c.GetByTagFeatureID(ctx, tagID, FeatureID); !useLast && ok {
		return content, nil
	}

	uncachedBannerContent, err := b.r.GetActiveContentByTagFeatureID(ctx, tagID, FeatureID)
	if err != nil {
		return nil, err
	}

	b.c.StoreByTagFeatureID(ctx, tagID, FeatureID, uncachedBannerContent)

	return uncachedBannerContent, nil
}

func (b *Service) GetByFilter(ctx context.Context, f domain.FilterBanner) (*[]domain.Banner, error) {
	return b.r.GetByFilter(ctx, f)
}

func (b *Service) Update(ctx context.Context, banner domain.UpdBanner) error {
	if banner.Tags != nil {
		*banner.Tags = getUnique(*banner.Tags)
	}
	return b.r.Update(ctx, banner)
}

func (b *Service) Create(ctx context.Context, banner domain.Banner) (int, error) {
	banner.Tags = getUnique(banner.Tags)
	return b.r.Insert(ctx, banner)
}

func (b *Service) DeleteById(ctx context.Context, id int) error {
	return b.r.DeleteByID(ctx, id)
}

func (b *Service) Delete(ctx context.Context, tagID, featureID *int) error {
	return b.d.Delete(ctx, tagID, featureID)
}

func getUnique[T comparable](sl []T) []T {
	m, unique := make(map[T]struct{}), make([]T, 0, len(sl))
	for _, v := range sl {
		if _, ok := m[v]; !ok {
			m[v], unique = struct{}{}, append(unique, v)
		}
	}
	return unique
}
