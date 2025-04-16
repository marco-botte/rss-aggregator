package main

import (
	"fmt"
	"rss-aggregator/internal/config"
)

func main() {
	conf := config.Read()
	conf.SetUser("marco")
	conf = config.Read()
	fmt.Printf("Current config: %v", *conf)
}
