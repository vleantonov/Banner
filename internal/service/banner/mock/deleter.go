package mock

import "context"

type Deleter struct{}

func (d Deleter) Delete(ctx context.Context, tagID, featureID *int) error {
	return nil
}
