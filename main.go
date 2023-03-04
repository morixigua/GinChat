package main

import (
	"ginchat/models"
	"ginchat/router"
	"ginchat/utils"
	"github.com/spf13/viper"
	"time"
)

func main() {
	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()
	InitTimer()
	r := router.Router()
	r.Run(viper.GetString("port.server"))
	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func InitTimer() {
	utils.Timer(
		time.Duration(viper.GetInt("timeout.DelayHeartbeat")),
		time.Duration(viper.GetInt("timeout.HeartbeatHz")),
		models.CleanConnection,
		"",
	)
}
