package domain

const DeleteQueue = "DeleteQueue"

type DeleteBodyQueue struct {
	TagID     *int `json:"tag_id,omitempty"`
	FeatureID *int `json:"feature_id,omitempty"`
}
