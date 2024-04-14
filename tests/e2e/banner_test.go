package e2e

import (
	"banner/internal/domain"
	"banner/tests/suite"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	httpurl "net/url"
	"sort"
	"strconv"
	"testing"
)

var BannerMap = map[string]interface{}{
	"feature_id": 1,
	"tag_ids":    []int{1},
	"content": map[string]interface{}{
		"Content": "Content",
	},
	"is_active": true,
}

func createBanner(st *suite.Suite) {

	url := fmt.Sprintf("http://%s:%d/banner", st.Cfg.ServerCfg.Host, st.Cfg.ServerCfg.Port)

	j, err := json.Marshal(BannerMap)
	if err != nil {
		st.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(j))
	if err != nil {
		st.Fatal(err)
	}

	req.Header.Set("token", st.AdminToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := st.HttpClient.Do(req)
	if err != nil {
		st.Fatal(err)
	}
	defer resp.Body.Close()
}

func getUserBanner(st *suite.Suite, t *testing.T, useLast bool, expStatus int, expContent map[string]interface{}) {
	url := fmt.Sprintf("http://%s:%d/user_banner", st.Cfg.ServerCfg.Host, st.Cfg.ServerCfg.Port)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		st.Fatal(err)
	}
	req.Header.Set("token", st.UserToken)
	q := httpurl.Values{}
	q.Add("tag_id", strconv.Itoa(1))
	q.Add("feature_id", strconv.Itoa(1))
	q.Add("use_last_revision", strconv.FormatBool(useLast))
	req.URL.RawQuery = q.Encode()

	resp, err := st.HttpClient.Do(req)
	if err != nil {
		st.Fatal(err)
	}
	defer resp.Body.Close()

	require.Equal(t, expStatus, resp.StatusCode)
	if expContent != nil {
		bodyJson, err := io.ReadAll(resp.Body)
		if err != nil {
			st.Fatal(err)
		}
		var body map[string]interface{}
		err = json.Unmarshal(bodyJson, &body)
		if err != nil {
			st.Fatal(err)
		}

		assert.Equal(t, expContent, body)
	}
}

func getAdminBanner(st *suite.Suite, t *testing.T, expBody *domain.Banner) int {
	url := fmt.Sprintf("http://%s:%d/banner", st.Cfg.ServerCfg.Host, st.Cfg.ServerCfg.Port)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		st.Fatal(err)
	}
	req.Header.Set("token", st.AdminToken)

	resp, err := st.HttpClient.Do(req)
	if err != nil {
		st.Fatal(err)
	}
	defer resp.Body.Close()
	bodyJson, err := io.ReadAll(resp.Body)
	if err != nil {
		st.Fatal(err)
	}

	var resBody []domain.Banner
	err = json.Unmarshal(bodyJson, &resBody)

	if expBody != nil {
		require.Equal(t, 1, len(resBody))

		sort.Ints(resBody[0].Tags)
		sort.Ints(expBody.Tags)

		assert.Equal(t, resBody[0].Tags, expBody.Tags)
		assert.Equal(t, len(resBody[0].Tags), len(expBody.Tags))
		assert.Equal(t, expBody.Feature, resBody[0].Feature)
		assert.Equal(t, expBody.Content, resBody[0].Content)
		assert.Equal(t, expBody.Active, resBody[0].Active)

		return resBody[0].ID
	}
	return 0
}

func setIsActiveFalse(st *suite.Suite, t *testing.T, id int) {
	url := fmt.Sprintf("http://%s:%d/banner/%d", st.Cfg.ServerCfg.Host, st.Cfg.ServerCfg.Port, id)

	j, err := json.Marshal(map[string]interface{}{
		"is_active": false,
	})
	if err != nil {
		st.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(j))
	if err != nil {
		st.Fatal(err)
	}
	req.Header.Set("token", st.AdminToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := st.HttpClient.Do(req)
	if err != nil {
		st.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestBannerActiveToggle(t *testing.T) {
	_, st := suite.New(t)

	createBanner(st)
	getUserBanner(st, t, true, http.StatusOK, BannerMap["content"].(map[string]interface{}))
	id := getAdminBanner(st, t, &domain.Banner{
		Content: BannerMap["content"].(map[string]interface{}),
		Tags:    BannerMap["tag_ids"].([]int),
		Feature: BannerMap["feature_id"].(int),
		Active:  BannerMap["is_active"].(bool),
	})

	setIsActiveFalse(st, t, id)
	getUserBanner(st, t, false, http.StatusOK, BannerMap["content"].(map[string]interface{}))
	getUserBanner(st, t, true, http.StatusNotFound, nil)
	getAdminBanner(st, t, &domain.Banner{
		Content: BannerMap["content"].(map[string]interface{}),
		Tags:    BannerMap["tag_ids"].([]int),
		Feature: BannerMap["feature_id"].(int),
		Active:  false,
	})
}
