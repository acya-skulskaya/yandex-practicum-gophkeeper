package secret

import (
	"context"
	"testing"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/crypto"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/repository/secret/mocks"
	"github.com/stretchr/testify/assert"
)

func TestService_List(t *testing.T) {
	ctx := context.Background()
	crpt, _ := crypto.NewCryptoMock("", "")

	type args struct {
		ctx              context.Context
		userID           uint
		getLatestVersion bool
	}
	tests := []struct {
		name      string
		setupMock func(m *mocks.MockSecretRepositoryInterface)
		args      args
		want      []models.Secret
		wantErr   bool
	}{
		{
			name: "gets a specified secret version",
			setupMock: func(m *mocks.MockSecretRepositoryInterface) {
				m.EXPECT().
					List(ctx, uint(1)).
					Return([]models.Secret{
						{
							UserID: uint(1),
							ID:     uint(1),
						},
					}, nil).
					Once()
			},
			args: args{
				ctx:              context.Background(),
				userID:           uint(1),
				getLatestVersion: false,
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
			got, err := s.List(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, got[0].UserID, tt.want[0].UserID)
		})
	}
}
