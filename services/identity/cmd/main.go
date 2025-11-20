package main

import (
	"context"
	"log"
	"net"
	"os/signal"
	"syscall"

	pb "go-audio-stream/pkg/proto/auth"
	grpcHandler "go-audio-stream/services/identity/internal/grpc"
	"go-audio-stream/services/identity/internal/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func gracefulShutdown(grpcServer *grpc.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// Gracefully stop gRPC server
	grpcServer.GracefulStop()

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {
	srv := server.NewServer()

	// gRPC Server
	grpcServer := grpc.NewServer()
	authServer := grpcHandler.NewServer(srv.DB, srv.FirebaseApp)
	pb.RegisterAuthServiceServer(grpcServer, authServer)

	// Health Check Server
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(grpcServer, done)

	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		log.Println("gRPC server listening on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}
