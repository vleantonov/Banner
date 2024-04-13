package integration

import (
	"banner/tests/suite"
	"encoding/json"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	httpurl "net/url"
	"strconv"
	"testing"
)

const (
	content1 = `{"hello_banner": {"key1": "value1", "key2": [1, 2, 3],"key3": 3}}`
	content2 = `{"hello_banner": {}}`
	content3 = `{"banner": {"key1": "value1"}}`
)

var initQueryForUserGet = fmt.Sprintf(
	`
	INSERT INTO banners (id, content, is_active)
		VALUES
		(1, '%s', true),
		(2, '%s', true),
		(3, '%s', false);
		
	INSERT INTO tag_feature_banners (tag_id, feature_id, banner_id)
		VALUES
		(1, 1, 1),
		(2, 1, 1),
		(2, 2, 2),
		(3, 1, 3);
	`, content1, content2, content3,
)

func TestGetUserBanner_Success(t *testing.T) {
	_, st := suite.New(t)
	st.PG.MustExec(initQueryForUserGet)

	url := fmt.Sprintf("http://%s:%d/user_banner", st.Cfg.ServerCfg.Host, st.Cfg.ServerCfg.Port)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		st.Fatal(err)
	}
	req.Header.Set("token", st.Cfg.ServerCfg.UserToken)

	testMap := []struct {
		name             string
		tagID            int
		FeatureID        int
		expectedStatus   int
		expectedBodyJson string
	}{
		{
			name:             "Get full banner 1",
			tagID:            1,
			FeatureID:        1,
			expectedStatus:   http.StatusOK,
			expectedBodyJson: content1,
		},
		{
			name:             "Get full banner 1 with new tag",
			tagID:            2,
			FeatureID:        1,
			expectedStatus:   http.StatusOK,
			expectedBodyJson: content1,
		},
		{
			name:             "Get full banner 2",
			tagID:            2,
			FeatureID:        2,
			expectedStatus:   http.StatusOK,
			expectedBodyJson: content2,
		},
		{
			name:             "Get inactive banner",
			tagID:            3,
			FeatureID:        1,
			expectedStatus:   http.StatusNotFound,
			expectedBodyJson: "",
		},
		{
			name:             "Get unexciting banner",
			tagID:            10,
			FeatureID:        10,
			expectedStatus:   http.StatusNotFound,
			expectedBodyJson: "",
		},
	}

	for _, test := range testMap {

		q := httpurl.Values{}

		q.Add("tag_id", strconv.Itoa(test.tagID))
		q.Add("feature_id", strconv.Itoa(test.FeatureID))
		req.URL.RawQuery = q.Encode()

		resp, err := st.HttpClient.Do(req)
		if err != nil {
			st.Fatal(err)
		}
		defer resp.Body.Close()

		bodyJson, err := io.ReadAll(resp.Body)
		if err != nil {
			st.Fatal(err)
		}

		assert.Equal(t, test.expectedStatus, resp.StatusCode)

		if test.expectedBodyJson != "" {
			var expBody, body map[string]interface{}
			err = json.Unmarshal([]byte(test.expectedBodyJson), &expBody)
			if err != nil {
				st.Fatal(err)
			}
			err = json.Unmarshal(bodyJson, &body)
			if err != nil {
				st.Fatal(err)
			}

			assert.Equal(t, expBody, body)
		}
	}
}

func TestUserBanner_InvalidParams(t *testing.T) {
	_, st := suite.New(t)

	st.PG.MustExec(initQueryForUserGet)

	url := fmt.Sprintf("http://%s:%d/user_banner", st.Cfg.ServerCfg.Host, st.Cfg.ServerCfg.Port)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		st.Fatal(err)
	}
	req.Header.Set("token", st.Cfg.ServerCfg.UserToken)

	resp, err := st.HttpClient.Do(req)
	if err != nil {
		st.Fatal(err)
	}
	resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	q := httpurl.Values{}
	q.Add("tag_id", strconv.Itoa(1))
	req.URL.RawQuery = q.Encode()

	resp, err = st.HttpClient.Do(req)
	if err != nil {
		st.Fatal(err)
	}
	resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	q = httpurl.Values{}
	q.Add("feature_id", strconv.Itoa(1))
	req.URL.RawQuery = q.Encode()

	resp, err = st.HttpClient.Do(req)
	if err != nil {
		st.Fatal(err)
	}
	resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
