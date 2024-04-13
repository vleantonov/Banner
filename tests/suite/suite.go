package suite

import (
	"banner/internal/config"
	"context"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

const cleanDBPath = "../suite/clean.sql"

type Suite struct {
	*testing.T
	Cfg        *config.Config
	HttpClient *http.Client
	PG         *sqlx.DB
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()

	cfg, err := config.New()
	if err != nil {
		t.Fatalf("can't load test config: %v", err)
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)

	db, err := sqlx.Connect("pgx", cfg.StorageCfg.PGUrl)
	if err != nil {
		t.Fatalf("can't connect to db: %v", err)
	}

	cleanQueries, err := getQueriesFromFile(cleanDBPath)
	if err != nil {
		t.Fatalf("open clean file in %s: %v", cleanDBPath, err)
	}

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()

		for _, query := range cleanQueries {
			db.Exec(query)
		}
	})

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		HttpClient: &http.Client{},
		PG:         db,
	}
}

func getQueriesFromFile(filepath string) ([]string, error) {

	file, err := os.ReadFile(filepath)

	if err != nil {
		return nil, err
	}

	return strings.Split(string(file), ";"), nil

}
