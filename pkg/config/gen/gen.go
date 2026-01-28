package main

import (
	cfg "github.com/conductorone/baton-privx/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	config.Generate("privx", cfg.Config)
}
