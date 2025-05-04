package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"publisher-service/cmd/webservice"
	"publisher-service/internal/config"
)

const logTagMain = "[main]"

func main() {
	config.Init()
	conf := config.Get()
	webservice.Start(conf)
	confBytes, _ := json.Marshal(conf)
	slog.Info(fmt.Sprintf("%s starting service with the config of: %s", logTagMain, confBytes))
}
