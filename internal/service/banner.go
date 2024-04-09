package service

import (
	"banner/internal/domain"
	"context"
	"go.uber.org/zap"
)

type Repository interface {
	GetByTagFeatureID(ctx context.Context, tagId, featureId int) (*domain.Banner, error)
	GetByFilter(ctx context.Context, f domain.FilterBanner) (*[]domain.Banner, error)
	Insert(ctx context.Context, b domain.Banner) (int, error)
	Update(ctx context.Context, b domain.UpdBanner) error
	DeleteByID(ctx context.Context, id int) error
}

type BannerService struct {
	l *zap.Logger
	r Repository
}

func New(l *zap.Logger, r Repository) *BannerService {
	return &BannerService{
		l: l,
		r: r,
	}
}

func (b *BannerService) GetByTagFeatureID(ctx context.Context, tagID, FeatureID int) (*domain.Banner, error) {
	return b.r.GetByTagFeatureID(ctx, tagID, FeatureID)
}

func (b *BannerService) GetByFilter(ctx context.Context, f domain.FilterBanner) (*[]domain.Banner, error) {
	return b.r.GetByFilter(ctx, f)
}

func (b *BannerService) Update(ctx context.Context, banner domain.UpdBanner) error {
	if banner.Tags != nil {
		*banner.Tags = getUnique(*banner.Tags)
	}
	return b.r.Update(ctx, banner)
}

func (b *BannerService) Create(ctx context.Context, banner domain.Banner) (int, error) {
	banner.Tags = getUnique(banner.Tags)
	return b.r.Insert(ctx, banner)
}

func (b *BannerService) Delete(ctx context.Context, id int) error {
	return b.r.DeleteByID(ctx, id)
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
