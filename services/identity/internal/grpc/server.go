package grpc

import (
	"context"
	"go-audio-stream/pkg/database"
	"go-audio-stream/pkg/models"
	pb "go-audio-stream/pkg/proto/auth"
	"log"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
	db          database.Service
	firebaseApp *firebase.App
}

func NewServer(db database.Service, firebaseApp *firebase.App) *Server {
	return &Server{
		db:          db,
		firebaseApp: firebaseApp,
	}
}

func (s *Server) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	client, err := s.firebaseApp.Auth(ctx)
	if err != nil {
		log.Printf("error getting Auth client: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to initialize auth client")
	}

	verifiedToken, err := client.VerifyIDToken(ctx, req.Token)
	if err != nil {
		log.Printf("error verifying ID token: %v\n", err)
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	var user models.User
	_, err = s.db.Find(&user, "email = ?", verifiedToken.Claims["email"])
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	if user.IsSuspended {
		return nil, status.Errorf(codes.PermissionDenied, "user is suspended")
	}

	return &pb.VerifyTokenResponse{
		Id:          user.ID,
		Email:       user.Email,
		Name:        user.FirstName + " " + user.LastName,
		IsSuspended: user.IsSuspended,
	}, nil
}
