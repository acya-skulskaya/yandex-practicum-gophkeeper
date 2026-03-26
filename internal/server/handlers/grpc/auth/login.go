package auth

import (
	"context"
	"errors"

	pb "github.com/acya-skulskaya/yandex-practicum-gophkeeper/api/gophkeeper"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/handlers/request"
	userRepo "github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/repository/user"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/service/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func (s *AuthServer) Login(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error) {
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

	user, err := s.ServiceAuth.Repo.Get(ctx, userCreds.Login)
	if err != nil {
		if errors.Is(err, userRepo.ErrNotFound) {
			//nolint:wrapcheck
			return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
		} else {
			//nolint:wrapcheck
			return nil, status.Errorf(codes.Internal, "could not get user: %v", err)
		}
	}

	if err = auth.VerifyPassword(user.Password, password); err != nil {
		//nolint:wrapcheck
		return nil, status.Errorf(codes.Unauthenticated, "could not verify password: %v", err)
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
