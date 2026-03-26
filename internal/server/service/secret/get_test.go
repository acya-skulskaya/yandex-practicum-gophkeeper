package secret

import (
	"context"
	"testing"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/crypto"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/repository/secret/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_Get(t *testing.T) {
	ctx := context.Background()
	crpt, _ := crypto.NewCryptoMock("", "")

	type args struct {
		ctx              context.Context
		userID           uint
		id               uint
		versionID        uint
		getLatestVersion bool
	}
	tests := []struct {
		name      string
		setupMock func(m *mocks.MockSecretRepositoryInterface)
		args      args
		want      models.Secret
		wantErr   bool
	}{
		{
			name: "gets a specified secret version",
			setupMock: func(m *mocks.MockSecretRepositoryInterface) {
				m.EXPECT().
					GetWithVersion(ctx, uint(1), uint(1), uint(1)).
					Return(models.Secret{
						UserID: uint(1),
						ID:     uint(1),
						Versions: []models.SecretVersion{
							{
								CreatedAt: nil,
								Data:      "",
								SecretID:  0,
								ID:        0,
							},
						},
					}, nil).
					Maybe()
				m.EXPECT().
					GetWithVersionsList(ctx, uint(1), uint(1)).
					Return(models.Secret{
						UserID: uint(2),
						ID:     uint(2),
						Versions: []models.SecretVersion{
							{
								CreatedAt: nil,
								Data:      "",
								SecretID:  0,
								ID:        0,
							},
						},
					}, nil).
					Maybe()
			},
			args: args{
				ctx:              context.Background(),
				userID:           uint(1),
				id:               uint(1),
				versionID:        uint(1),
				getLatestVersion: false,
			},
			want: models.Secret{
				UserID: 1,
				ID:     1,
			},
			wantErr: false,
		},
		{
			name: "gets latest secret version when version is 0 and getLatestVersion=true",
			setupMock: func(m *mocks.MockSecretRepositoryInterface) {
				m.EXPECT().
					GetWithVersion(ctx, uint(1), uint(1), mock.Anything).
					Return(models.Secret{
						UserID: uint(1),
						ID:     uint(1),
						Versions: []models.SecretVersion{
							{
								CreatedAt: nil,
								Data:      "",
								SecretID:  0,
								ID:        0,
							},
						},
					}, nil).
					Maybe()
				m.EXPECT().
					GetWithVersionsList(ctx, uint(1), uint(1)).
					Return(models.Secret{
						UserID: uint(2),
						ID:     uint(2),
						Versions: []models.SecretVersion{
							{
								CreatedAt: nil,
								Data:      "",
								SecretID:  0,
								ID:        0,
							},
						},
					}, nil).
					Maybe()
			},
			args: args{
				ctx:              context.Background(),
				userID:           uint(1),
				id:               uint(1),
				versionID:        uint(0),
				getLatestVersion: true,
			},
			want: models.Secret{
				UserID: 1,
				ID:     1,
			},
			wantErr: false,
		},
		{
			name: "gets secret with version list when version is 0 and getLatestVersion=false",
			setupMock: func(m *mocks.MockSecretRepositoryInterface) {
				m.EXPECT().
					GetWithVersion(ctx, uint(1), uint(1), mock.Anything).
					Return(models.Secret{
						UserID: uint(1),
						ID:     uint(1),
						Versions: []models.SecretVersion{
							{
								CreatedAt: nil,
								Data:      "",
								SecretID:  0,
								ID:        0,
							},
						},
					}, nil).
					Maybe()
				m.EXPECT().
					GetWithVersionsList(ctx, uint(1), uint(1)).
					Return(models.Secret{
						UserID: uint(2),
						ID:     uint(2),
						Versions: []models.SecretVersion{
							{
								CreatedAt: nil,
								Data:      "",
								SecretID:  0,
								ID:        0,
							},
						},
					}, nil).
					Maybe()
			},
			args: args{
				ctx:              context.Background(),
				userID:           uint(1),
				id:               uint(1),
				versionID:        uint(0),
				getLatestVersion: false,
			},
			want: models.Secret{
				UserID: 2,
				ID:     2,
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
			got, err := s.Get(tt.args.ctx, tt.args.userID, tt.args.id, tt.args.versionID, tt.args.getLatestVersion)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, got.UserID, tt.want.UserID)
			assert.Equal(t, got.ID, tt.want.ID)
		})
	}
}
