package main

import (
	"fmt"

	"github.com/Goga211/go-auth/internal/config"
)

func main() {
	cfg := config.Load()
	fmt.Println(cfg)
}
