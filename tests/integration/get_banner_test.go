package integration

import (
	"banner/internal/domain"
	"banner/tests/suite"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	httpurl "net/url"
	"strconv"
	"testing"
)

func getPointеrWithValue[T any](v T) *T {
	newVal := v
	return &newVal
}

func TestGetBanner_SuccessCheckLen(t *testing.T) {
	_, st := suite.New(t)
	st.PG.MustExec(initQueryForUserGet)

	url := fmt.Sprintf("http://%s:%d/banner", st.Cfg.ServerCfg.Host, st.Cfg.ServerCfg.Port)

	tests := []struct {
		name            string
		filter          domain.FilterBanner
		expectedLenData int
		status          int
	}{
		{
			name:            "Get Full Records",
			filter:          domain.FilterBanner{},
			expectedLenData: 3,
			status:          http.StatusOK,
		},
		{
			name:            "Get by tag",
			filter:          domain.FilterBanner{TagID: getPointеrWithValue(2)},
			expectedLenData: 2,
			status:          http.StatusOK,
		},
		{
			name:            "Get by feature",
			filter:          domain.FilterBanner{FeatureID: getPointеrWithValue(1)},
			expectedLenData: 2,
			status:          http.StatusOK,
		},
		{
			name:            "Get with limit 2",
			filter:          domain.FilterBanner{Limit: getPointеrWithValue(2)},
			expectedLenData: 2,
			status:          http.StatusOK,
		},
		{
			name: "Get with tag and feature",
			filter: domain.FilterBanner{
				TagID:     getPointеrWithValue(1),
				FeatureID: getPointеrWithValue(1),
			},
			expectedLenData: 1,
			status:          http.StatusOK,
		},
		{
			name: "Get nonexistent with tag and feature",
			filter: domain.FilterBanner{
				TagID:     getPointеrWithValue(10),
				FeatureID: getPointеrWithValue(10),
			},
			expectedLenData: 0,
			status:          http.StatusNotFound,
		},
		{
			name: "Get nonexistent with tag",
			filter: domain.FilterBanner{
				TagID: getPointеrWithValue(10),
			},
			expectedLenData: 0,
			status:          http.StatusNotFound,
		},
		{
			name: "Get nonexistent with feature",
			filter: domain.FilterBanner{
				FeatureID: getPointеrWithValue(10),
			},
			expectedLenData: 0,
			status:          http.StatusNotFound,
		},
	}

	for _, test := range tests {
		q := httpurl.Values{}
		if test.filter.TagID != nil {
			q.Add("tag_id", strconv.Itoa(*test.filter.TagID))
		}
		if test.filter.FeatureID != nil {
			q.Add("feature_id", strconv.Itoa(*test.filter.FeatureID))
		}
		if test.filter.Limit != nil {
			q.Add("limit", strconv.Itoa(*test.filter.Limit))
		}
		if test.filter.Offset != nil {
			q.Add("offset", strconv.Itoa(*test.filter.Offset))
		}

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			st.Fatal(err)
		}

		req.URL.RawQuery = q.Encode()
		req.Header.Set("token", st.Cfg.ServerCfg.AdminToken)

		resp, err := st.HttpClient.Do(req)
		if err != nil {
			st.Fatal(err)
		}
		defer resp.Body.Close()

		require.Equal(t, test.status, resp.StatusCode)

		if test.expectedLenData != 0 {
			bodyJson, err := io.ReadAll(resp.Body)
			if err != nil {
				st.Fatal(err)
			}

			var resBody []domain.Banner
			err = json.Unmarshal(bodyJson, &resBody)

			require.NoError(t, err)

			assert.Equal(t, test.expectedLenData, len(resBody))
		}
	}
}
