package main

import (
	"banner/internal/config"
	"banner/internal/domain"
	"banner/internal/pkg/jwt"
	"fmt"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("can't open config: %v", err)
	}

	t, err := jwt.NewToken(&domain.User{
		Login:    "user",
		PassHash: []byte("1234"),
		IsAdmin:  true,
	}, cfg.ServerCfg.TTLToken, cfg.ServerCfg.AppSecret)

	if err != nil {
		log.Fatalf("can't create token: %v", err)
	}

	fmt.Println(t)
}
