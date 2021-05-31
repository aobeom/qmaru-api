package main

import (
	"qmaru-api/apis"
	"qmaru-api/configs"
)

func main() {
	apis.Run(configs.Deployment())
}
