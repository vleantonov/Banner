package integration

import (
	"banner/tests/suite"
	"bytes"
	"encoding/json"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

const (
	checkBannerQuery = `SELECT id, content, is_active FROM banners`
	checkTagsQuery   = `SELECT banner_id, feature_id, array_agg(tag_id) FROM tag_feature_banners GROUP BY banner_id, feature_id`
)

var reqBodyPost = map[string]interface{}{
	"feature_id": 1,
	"tag_ids":    []int{1, 2, 3},
	"content": map[string]interface{}{
		"Content": "Content",
	},
	"is_active": true,
}

func TestPostBanner_Success(t *testing.T) {
	_, st := suite.New(t)
	url := fmt.Sprintf("http://%s:%d/banner", st.Cfg.ServerCfg.Host, st.Cfg.ServerCfg.Port)

	j, err := json.Marshal(reqBodyPost)
	if err != nil {
		st.Fatal(err)
	}

	rows, err := st.PG.Query(checkBannerQuery)
	require.NoError(t, err)
	require.False(t, rows.Next())
	rows, err = st.PG.Query(checkTagsQuery)
	require.NoError(t, err)
	require.False(t, rows.Next())

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

	bodyJson, err := io.ReadAll(resp.Body)
	if err != nil {
		st.Fatal(err)
	}

	var resBody map[string]int
	err = json.Unmarshal(bodyJson, &resBody)
	if err != nil {
		st.Fatal(err)
	}

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, 1, len(resBody))
	assert.Contains(t, resBody, "banner_id")
}
