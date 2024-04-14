package e2e

import (
	"banner/tests/suite"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

const initAuth = `
	INSERT INTO banners (id, content, is_active) VALUES (1, '{"content": "content"}', true);
	INSERT INTO tag_feature_banners (tag_id, feature_id, banner_id) VALUES (1, 1, 1);
`

func getUserBannerWithToken(st *suite.Suite, token string, expStatusCode int) {
	url := fmt.Sprintf("http://%s:%d/user_banner?tag_id=1&feature_id=1", st.Cfg.ServerCfg.Host, st.Cfg.ServerCfg.Port)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		st.Fatal(err)
	}
	if token != "" {
		req.Header.Add("token", token)
	}
	resp, err := st.HttpClient.Do(req)
	if err != nil {
		st.Fatal(err)
	}
	resp.Body.Close()

	assert.Equal(st, expStatusCode, resp.StatusCode)
}

func TestAuth(t *testing.T) {
	_, st := suite.New(t)
	st.PG.MustExec(initAuth)

	getUserBannerWithToken(st, "", http.StatusUnauthorized)

	url := fmt.Sprintf("http://%s:%d/auth/register", st.Cfg.ServerCfg.Host, st.Cfg.ServerCfg.Port)
	loginData := map[string]string{
		"login":    "randomuser",
		"password": "randomuser",
	}
	jBody, err := json.Marshal(loginData)
	if err != nil {
		st.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jBody))
	if err != nil {
		st.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := st.HttpClient.Do(req)
	if err != nil {
		st.Fatal(err)
	}
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	url = fmt.Sprintf("http://%s:%d/auth/login", st.Cfg.ServerCfg.Host, st.Cfg.ServerCfg.Port)
	req, err = http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jBody))
	if err != nil {
		st.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err = st.HttpClient.Do(req)
	if err != nil {
		st.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		st.Fatal(err)
	}
	var r struct {
		Token string `json:"token"`
	}
	err = json.Unmarshal(respBody, &r)
	if err != nil {
		st.Fatal(err)
	}
	resp.Body.Close()
	getUserBannerWithToken(st, r.Token, http.StatusOK)
}
