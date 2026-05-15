package main

import (
	"fmt"
	"gator/internal/config"
)

func main() {
	cfg := config.Read()
	cfg.SetUser("andrew")
	newCfg := config.Read()
	fmt.Printf("db_url: '%v'\n", newCfg.DBUrl)
	fmt.Printf("current_user_name: '%v'\n", newCfg.CurrentUsername)
}
