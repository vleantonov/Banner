package domain

import "errors"

var (
	ErrBannerNotFound          = errors.New("can't find any banner")
	ErrInternalServerError     = errors.New("internal Server Error")
	ErrTagFeatureAlreadyExists = errors.New("banner with that tag_id and feature_id already exists")
	ErrTagFeatureRequired      = errors.New("tag_ids and feature_id must be indicated")
	ErrContentRequired         = errors.New("content is required")
)
