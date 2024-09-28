package main

import (
	"google.golang.org/grpc/reflection"
	"log"
	"net"

	"google.golang.org/grpc"
	"notification-service/grpcservice"
	"notification-service/mailer"
	"notification-service/rest"
)

func main() {
	go func() {
		rest.StartServer(":8080")
	}()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	mailer.RegisterMailerServer(grpcServer, &grpcservice.MailService{})

	reflection.Register(grpcServer)

	log.Println("Starting gRPC server on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
