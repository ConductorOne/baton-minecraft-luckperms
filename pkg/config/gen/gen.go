package main

import (
	cfg "github.com/conductorone/baton-minecraft-luckperms/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	config.Generate("minecraft-luckperms", cfg.Config)
}
