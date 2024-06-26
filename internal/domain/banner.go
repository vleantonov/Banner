package domain

type Banner struct {
	ID      int                    `json:"id,omitempty" db:"id"`
	Content map[string]interface{} `json:"content" db:"content"`
	Tags    []int                  `json:"tag_ids" db:"tag_ids"`
	Feature int                    `json:"feature_id" db:"feature_id"`
	Active  bool                   `json:"is_active" db:"is_active"`

	CreatedAT string `json:"created_at" db:"created_at"`
	UpdatedAT string `json:"updated_at" db:"updated_at"`
}

type UpdBanner struct {
	ID      int                     `json:"id" db:"id"`
	Content *map[string]interface{} `json:"content" db:"content"`
	Tags    *[]int                  `json:"tag_ids" db:"tag_ids"`
	Feature *int                    `json:"feature_id" db:"feature_id"`
	Active  *bool                   `json:"is_active" db:"is_active"`
}

type TagFeatureBanner struct {
	TagID     int `json:"tag_id" db:"tag_id"`
	FeatureID int `json:"feature_id" db:"feature_id"`
	BannerID  int `json:"banner_id" db:"banner_id"`
}

type FilterBanner struct {
	FeatureID *int `db:"feature_id"`
	TagID     *int `db:"banner_id"`

	Limit  *int
	Offset *int
}
