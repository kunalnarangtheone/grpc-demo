package main

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	logger "github.com/sirupsen/logrus"

	"my_grpc_service/config"
	pb "my_grpc_service/proto"
)

// printFeature gets the feature for the given point.
func printFeature(client pb.RouteGuideClient, point *pb.Point) {
	logger.Debugf("Getting feature for point (%d, %d)", point.Latitude, point.Longitude)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if feature, err := client.GetFeature(ctx, point); err != nil {
		logger.Errorf("GetFeature failed: %v", err)
	} else {
		logger.Infof("GetFeature success: %v", feature)
	}
}

func main() {
	logger.SetLevel(logger.DebugLevel)
	creds, err := credentials.NewClientTLSFromFile(config.CertFile, "")
	if err != nil {
		logger.Fatalf("client: failed to get TLS credentials: %v", err)
	}
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(creds))

	serverAddr := fmt.Sprintf("%v:%d", config.HostName, config.Port)
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		logger.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewRouteGuideClient(conn)
	printFeature(client, &pb.Point{Latitude: 409146138, Longitude: -746188906})
}
