package secret

import (
	"context"
	"testing"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/crypto"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/repository/secret/mocks"
	"github.com/stretchr/testify/mock"
)

func TestService_Store(t *testing.T) {
	ctx := context.Background()
	crpt, _ := crypto.NewCryptoMock("", "")

	type args struct {
		ctx    context.Context
		userID uint
	}
	tests := []struct {
		name      string
		setupMock func(m *mocks.MockSecretRepositoryInterface)
		args      args
		want      []models.Secret
		wantErr   bool
	}{
		{
			name: "creates a new secret",
			setupMock: func(m *mocks.MockSecretRepositoryInterface) {
				m.EXPECT().
					Create(ctx, uint(1), mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(models.Secret{
						UserID: uint(1),
						ID:     uint(1),
					}, nil).
					Once()
			},
			args: args{
				ctx:    context.Background(),
				userID: uint(1),
			},
			want: []models.Secret{
				{
					UserID: uint(1),
					ID:     uint(1),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockSecretRepositoryInterface(t)
			tt.setupMock(mockRepo)

			s := &Service{
				Repo:      mockRepo,
				Crypto:    crpt,
				SecretKey: "",
			}
			_, err := s.Store(tt.args.ctx, tt.args.userID, mock.Anything, mock.Anything, []byte{}, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Store() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
