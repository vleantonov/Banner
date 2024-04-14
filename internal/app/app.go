package app

import (
	"banner/internal/config"
	banRoutes "banner/internal/handler/http/v1"
	api "banner/internal/handler/http/v1/gen"
	"banner/internal/pkg/cache"
	"banner/internal/pkg/database"
	"banner/internal/pkg/logger"
	"banner/internal/pkg/rabbitmq"
	repo "banner/internal/repository/postresql"
	"banner/internal/service/auth"
	"banner/internal/service/banner"
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

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("can't create config: %v", err)
	}

	l, err := logger.New()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	httpEngine := gin.New()
	httpEngine.Use(
		ginzap.Ginzap(l, "", false),
		ginzap.RecoveryWithZap(l, true),
	)

	db, err := database.CreateDBConnection(cfg.StorageCfg.PGUrl, l)
	if err != nil {
		l.Fatal("can't create pgpool for app", zap.Error(err))
	}

	pgRepo := repo.New(db)

	rmqProducer, err := rabbitmq.SetupRMQ(cfg.QueueCfg.RMQUrl)
	if err != nil {
		l.Fatal("can't setup RMQ", zap.Error(err))
	}

	bannerCache := cache.SetupCache(cfg.StorageCfg.TTLCache, l)

	bannerService := banner.New(
		pgRepo,
		bannerCache,
		rmqProducer,
	)

	authService := auth.New(l, pgRepo, time.Hour, cfg.ServerCfg.AppSecret)

	if err != nil {
		l.Fatal("can't create connection with postgres database", zap.Error(err))
	}

	api.RegisterHandlers(
		httpEngine,
		banRoutes.New(
			bannerService,
			authService,
		),
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
