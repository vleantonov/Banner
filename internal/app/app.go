package app

import (
	"banner/internal/config"
	banRoutes "banner/internal/handler/http/v1"
	api "banner/internal/handler/http/v1/gen"
	"banner/internal/pkg/logger"
	repo "banner/internal/repository/postresql"
	"banner/internal/service"
	"fmt"
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"log"
)

const (
	driverName = "pgx"
)

type App struct {
	e *gin.Engine
	c *config.Config
	l *zap.Logger
}

func New() *App {

	cfg, err := config.New(config.FetchConfigPath())
	if err != nil {
		log.Fatalf("can't create config: %v", err)
	}

	l, err := logger.New(cfg.Env)
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	httpEngine := gin.New()
	httpEngine.Use(
		ginzap.Ginzap(l, "", false),
		ginzap.RecoveryWithZap(l, true),
	)

	db, err := createDBConnection(cfg.StorageCfg.PGUrl, l)
	if err != nil {
		l.Fatal("can't create pgpool for app", zap.Error(err))
	}

	pgRepo := repo.New(db)
	bannerService := service.New(l, pgRepo)

	if err != nil {
		l.Fatal("can't create connection with postgres database", zap.Error(err))
	}

	api.RegisterHandlers(
		httpEngine,
		banRoutes.New(bannerService),
	)

	l.Info("app has been successfully built")
	return &App{
		e: httpEngine,
		c: cfg,
		l: l,
	}
}

func (a *App) MustRun() {

	defer a.l.Sync()
	a.l.Info("server started")

	err := a.e.Run(
		fmt.Sprintf("%s:%d", a.c.ServerCfg.Host, a.c.ServerCfg.Port),
	)
	if err != nil {
		a.l.Fatal("can't run http engine", zap.Error(err))
	}
}

func createDBConnection(url string, log *zap.Logger) (*sqlx.DB, error) {

	pgCfg, err := pgx.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("can't parse pg config: %w", err)
	}

	pgLog := logger.NewPgxLogger(log)
	pgCfg.Tracer = pgLog

	nativeDB := stdlib.OpenDB(*pgCfg)
	if err != nil {
		return nil, err
	}

	nativeDB.SetMaxOpenConns(10)
	nativeDB.SetMaxIdleConns(5)

	return sqlx.NewDb(nativeDB, driverName), nil

}
