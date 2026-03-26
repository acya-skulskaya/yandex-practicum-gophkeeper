package interceptors

import (
	"context"
	"errors"
	"strings"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/config"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/logger"
	auth2 "github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/service/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	Cfg config.Auth
}

func NewAuthInterceptor(cfg config.Auth) *AuthInterceptor {
	return &AuthInterceptor{
		Cfg: cfg,
	}
}

func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		// skipping auth if request is to any of auth methods
		// info.FullMethod /acyaskulskaya.gophkeeper.auth.AuthService/Register
		// info.FullMethod /acyaskulskaya.gophkeeper.auth.AuthService/Login
		if strings.Contains(info.FullMethod, "auth.AuthService") {
			resp, handlerErr := handler(ctx, req)
			return resp, handlerErr
		}

		md, ok := metadata.FromIncomingContext(ctx)
		var token string
		var err error
		if ok && len(md[config.MetadataAuthKeyName]) > 0 {
			token = md[config.MetadataAuthKeyName][0]
		} else {
			return nil, status.Error(codes.Unauthenticated, "unauthenticated: could not read authorization token")
		}

		userID, err := auth2.GetUserIDFromToken(i.Cfg.SecretKey, token)
		if err != nil {
			if errors.Is(err, auth2.ErrTokenIsNotValid) {
				return nil, status.Errorf(codes.Unauthenticated, "unauthenticated: could not validate token: %v", err)
			} else {
				logger.Log.Debug("error getting user id from auth token", zap.Error(err))
				return nil, status.Errorf(codes.Internal, "could not set user grom token: %v", err)
			}
		}

		if userID == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "unauthenticated: user is empty")
		}

		logger.Log.Info("AuthUnaryInterceptor: got user id", zap.Any("userID", userID))
		ctx = context.WithValue(ctx, auth2.AuthContextKey(auth2.AuthContextKeyUserID), userID)

		resp, handlerErr := handler(ctx, req)
		return resp, handlerErr
	}
}
