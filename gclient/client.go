package gclient

import (
	"context"
	"log"
	"time"

	"qmaru-api/configs"

	"google.golang.org/grpc"
)

var extCfg = configs.ExtCfg()
var BindAddress = extCfg["api-ig"].(string)

func RPCData(url string) []string {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, BindAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Panicf("Connect Failed: %v", err)
		return []string{}
	}
	defer conn.Close()
	c := NewInstaMediaClient(conn)

	if url != "" {
		shareURL := url
		res, err := c.GetMedia(ctx, &ShareURL{
			Url:  shareURL,
			Code: "Instago",
		})
		if err != nil {
			log.Panicf("Client Error: %v", err)
			return []string{}
		}
		return res.GetUrls()
	}
	return []string{}
}
