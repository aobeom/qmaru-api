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
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := NewInstaMediaClient(conn)

	if url != "" {
		shareURL := url
		r, err := c.GetMedia(ctx, &ShareURL{Url: shareURL})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		return r.GetUrls()
	}
	return []string{}
}
