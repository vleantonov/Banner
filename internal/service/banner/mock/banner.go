package mock

import "banner/internal/domain"

var TestBanner = &domain.Banner{
	ID: 1,
	Content: map[string]interface{}{
		"BannerKey": "BannerValue",
	},
	Tags:    []int{1, 2},
	Feature: 1,
}
