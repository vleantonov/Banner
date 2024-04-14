package integration

import (
	"banner/tests/suite"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

const (
	initQueryForDeleteByID = `
	INSERT INTO banners (id, is_active)
		VALUES
		(1, true),
		(2, true),
		(3, false);
		
	INSERT INTO tag_feature_banners (tag_id, feature_id, banner_id)
		VALUES
		(1, 1, 1),
		(1, 2, 1),
		(2, 1, 2),
		(3, 1, 3);
	`

	checkRawsQuery = `
	SELECT * FROM banners JOIN public.tag_feature_banners tfb on banners.id = tfb.banner_id
	WHERE id=1
	`
)

func TestDeleteByIDBanner_Success(t *testing.T) {
	_, st := suite.New(t)
	url := fmt.Sprintf("http://%s:%d/banner/1", st.Cfg.ServerCfg.Host, st.Cfg.ServerCfg.Port)

	st.PG.MustExec(initQueryForDeleteByID)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		st.Fatal(err)
	}
	req.Header.Set("token", st.AdminToken)

	resp, err := st.HttpClient.Do(req)
	if err != nil {
		st.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	rows, err := st.PG.Queryx(checkRawsQuery)
	if err != nil {
		st.Fatal(err)
	}

	assert.False(t, rows.Next())
}
