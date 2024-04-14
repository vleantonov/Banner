package integration

import (
	"banner/tests/suite"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	httpurl "net/url"
	"strconv"
	"testing"
	"time"
)

const (
	initQueryForDelete = `
	INSERT INTO banners (id) SELECT generate_series(1, 4);
	INSERT INTO tag_feature_banners (tag_id, feature_id, banner_id) 
	VALUES 
	(1, 1, 1),
	(2, 1, 1),
	(3, 1, 2),
	(3, 2, 3),
	(4, 4, 4);
	`
	checkRawsTagQuery = `
		SELECT id FROM banners JOIN public.tag_feature_banners tfb on banners.id = tfb.banner_id
		WHERE tag_id=$1;
	`
	checkRawsFeatureQuery = `
		SELECT id FROM banners JOIN public.tag_feature_banners tfb on banners.id = tfb.banner_id
		WHERE feature_id=$1;
	`
	checkLenTable = `
		SELECT COUNT(*) FROM banners;
	`
)

func TestDeleteBanner(t *testing.T) {
	_, st := suite.New(t)
	url := fmt.Sprintf("http://%s:%d/banner", st.Cfg.ServerCfg.Host, st.Cfg.ServerCfg.Port)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		st.Fatal(err)
	}
	req.Header.Set("token", st.AdminToken)

	tests := []struct {
		name        string
		tagID       int
		featureID   int
		expectedCnt int
	}{
		{
			name:        "Delete with equal feature_id",
			tagID:       0,
			featureID:   1,
			expectedCnt: 2,
		},
		{
			name:        "Delete with equal tag_id",
			tagID:       0,
			featureID:   1,
			expectedCnt: 2,
		},
		{
			name:        "Delete all with feature_id and tag_id",
			tagID:       3,
			featureID:   1,
			expectedCnt: 1,
		},
		{
			name:        "Delete one with feature_id and tag_id",
			tagID:       4,
			featureID:   4,
			expectedCnt: 3,
		},
		{
			name:        "Delete nonexistent",
			tagID:       10,
			featureID:   10,
			expectedCnt: 4,
		},
	}

	for _, test := range tests {
		st.PG.MustExec("DELETE FROM banners")
		st.PG.MustExec(initQueryForDelete)

		t.Log(test.name)

		q := httpurl.Values{}
		if test.tagID > 0 {
			q.Add("tag_id", strconv.Itoa(test.tagID))
		}
		if test.featureID > 0 {
			q.Add("feature_id", strconv.Itoa(test.featureID))
		}
		req.URL.RawQuery = q.Encode()

		resp, err := st.HttpClient.Do(req)
		if err != nil {
			st.Fatal(err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusAccepted, resp.StatusCode)

		time.Sleep(time.Duration(500) * time.Millisecond)

		if test.tagID > 0 {
			rows := st.PG.MustExec(checkRawsTagQuery, test.tagID)
			n, err := rows.RowsAffected()
			require.NoError(t, err)
			assert.Equal(t, 0, int(n))
		}
		if test.featureID > 0 {
			rows := st.PG.MustExec(checkRawsFeatureQuery, test.featureID)
			n, err := rows.RowsAffected()
			require.NoError(t, err)
			assert.Equal(t, 0, int(n))
		}

		var cnt int
		err = st.PG.Get(&cnt, checkLenTable)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, test.expectedCnt, cnt)
		resp.Body.Close()
	}
}
