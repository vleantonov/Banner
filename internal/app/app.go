package app

import (
	api "banner/internal/api/gen"
	"banner/internal/config"
	"banner/internal/logger"
	"banner/internal/repository/postresql"
	banRoutes "banner/internal/routes/banner"
	"fmt"
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"time"
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
		ginzap.Ginzap(l, time.RFC3339, false),
		ginzap.RecoveryWithZap(l, true),
	)

	// TODO: take out creating *sqlx.DB
	postgresRepo, err := postresql.New(
		cfg.StorageCfg.PostgresCfg.Host,
		cfg.StorageCfg.PostgresCfg.Port,
		cfg.StorageCfg.PostgresCfg.Username,
		cfg.StorageCfg.PostgresCfg.Password,
		cfg.StorageCfg.PostgresCfg.DBName,
		cfg.StorageCfg.PostgresCfg.SSLMode,
		l,
	)

	if err != nil {
		l.Fatal("can't create connection with postgres database", zap.Error(err))
	}

	api.RegisterHandlers(
		httpEngine,
		banRoutes.New(postgresRepo),
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
