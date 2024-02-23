package main

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	logger "github.com/sirupsen/logrus"

	"my_grpc_service/config"
	pb "my_grpc_service/proto"
)

type routeGuideServer struct {
	pb.UnimplementedRouteGuideServer
}

func (s *routeGuideServer) GetFeature(ctx context.Context, point *pb.Point) (*pb.Feature, error) {
	logger.Debugf("GetFeature called with point (%d, %d)", point.Latitude, point.Longitude)
	return &pb.Feature{Name: "example-feature", Location: point}, nil
}

func main() {
	logger.SetLevel(logger.DebugLevel)
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%d", config.HostName, config.Port))
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	tlsCredentials, err := credentials.NewServerTLSFromFile("creds/server-cert.pem", "creds/server-key.pem")
	if err != nil {
		logger.Fatal("server: failed to get tls credentials")
	}

	grpcServer := grpc.NewServer(grpc.Creds(tlsCredentials))
	s := routeGuideServer{}
	pb.RegisterRouteGuideServer(grpcServer, &s)

	logger.Info("Starting server")
	if err = grpcServer.Serve(lis); err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}

	defer grpcServer.Stop()
}
