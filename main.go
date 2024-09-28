package main

import (
	"fmt"
	"log"

	"github.com/CoDanTheBarbarian/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v\n", err)
	}
	fmt.Printf("Read config: %+v", cfg)
	err = cfg.SetUser("Daniel")
	if err != nil {
		log.Fatalf("Error setting user: %v\n", err)
	}
	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v\n", err)
	}
	fmt.Printf("Read config again: %+v", cfg)
}
