package auth

import (
	"context"

	pb "github.com/acya-skulskaya/yandex-practicum-gophkeeper/api/gophkeeper"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/handlers/request"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func (s *AuthServer) Register(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error) {
	username := in.GetUsername()
	password := in.GetPassword()

	userCreds := request.UserCredits{
		Login:    username,
		Password: password,
	}
	problems, err := request.IsValid(userCreds)
	if err != nil {
		//nolint:wrapcheck
		return nil, status.Error(codes.InvalidArgument, "request is invalid: "+problems.String())
	}

	user, err := s.ServiceAuth.Register(ctx, userCreds)
	if err != nil {
		//nolint:wrapcheck
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	tokenString, err := s.ServiceAuth.Token(user)
	if err != nil {
		//nolint:wrapcheck
		return nil, status.Errorf(codes.Internal, "could not get token: %v", err)
	}

	response := pb.AuthResponse_builder{
		Token: proto.String(tokenString),
	}

	return response.Build(), nil
}
