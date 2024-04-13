package domain

import "errors"

var (
	ErrBannerNotFound          = errors.New("can't find any banner")
	ErrInternalServerError     = errors.New("internal Server Error")
	ErrTagFeatureAlreadyExists = errors.New("banner with that tag_id and feature_id already exists")
	ErrTagFeatureRequired      = errors.New("tag_ids and feature_id must be provided")
	ErrTagOrFeatureRequired    = errors.New("tag_id or feature_id must be provided")
)
