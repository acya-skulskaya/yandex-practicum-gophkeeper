package secret

import (
	"context"
	"testing"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/models"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/crypto"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/repository/secret/mocks"
)

func TestService_Delete(t *testing.T) {
	ctx := context.Background()
	crpt, _ := crypto.NewCryptoMock("", "")

	type args struct {
		ctx       context.Context
		userID    uint
		id        uint
		versionID uint
	}
	tests := []struct {
		name      string
		setupMock func(m *mocks.MockSecretRepositoryInterface)
		args      args
		want      models.Secret
		wantErr   bool
	}{
		{
			name: "does not exist",
			setupMock: func(m *mocks.MockSecretRepositoryInterface) {
				m.EXPECT().
					Exists(ctx, uint(1), uint(1)).
					Return(false, nil).
					Maybe()
			},
			args: args{
				ctx:       context.Background(),
				userID:    uint(1),
				id:        uint(1),
				versionID: uint(1),
			},
			wantErr: true,
		},
		{
			name: "version id 0",
			setupMock: func(m *mocks.MockSecretRepositoryInterface) {
				m.EXPECT().
					Exists(ctx, uint(1), uint(1)).
					Return(true, nil).
					Maybe()

				m.EXPECT().
					Delete(ctx, uint(1)).
					Return(nil).
					Maybe()

				m.EXPECT().
					GetWithVersionsList(ctx, uint(1), uint(1)).
					Return(models.Secret{
						UserID:   uint(1),
						ID:       uint(1),
						Versions: []models.SecretVersion{},
					}, nil).
					Maybe()
			},
			args: args{
				ctx:       context.Background(),
				userID:    uint(1),
				id:        uint(1),
				versionID: uint(0),
			},
			wantErr: false,
		},
		{
			name: "version id != 0",
			setupMock: func(m *mocks.MockSecretRepositoryInterface) {
				m.EXPECT().
					Exists(ctx, uint(1), uint(1)).
					Return(true, nil).
					Maybe()

				m.EXPECT().
					ExistsVersion(ctx, uint(1)).
					Return(true, nil).
					Maybe()

				m.EXPECT().
					GetFilePath(uint(1), uint(1)).
					Return("", nil).
					Maybe()

				m.EXPECT().
					VersionsCount(ctx, uint(1)).
					Return(3, nil).
					Maybe()

				m.EXPECT().
					DeleteVersion(ctx, uint(1)).
					Return(nil).
					Maybe()
			},
			args: args{
				ctx:       context.Background(),
				userID:    uint(1),
				id:        uint(1),
				versionID: uint(1),
			},
			wantErr: false,
		},
		{
			name: "only one version",
			setupMock: func(m *mocks.MockSecretRepositoryInterface) {
				m.EXPECT().
					Exists(ctx, uint(1), uint(1)).
					Return(true, nil).
					Maybe()

				m.EXPECT().
					ExistsVersion(ctx, uint(1)).
					Return(true, nil).
					Maybe()

				m.EXPECT().
					GetFilePath(uint(1), uint(1)).
					Return("", nil).
					Maybe()

				m.EXPECT().
					VersionsCount(ctx, uint(1)).
					Return(1, nil).
					Maybe()

				m.EXPECT().
					DeleteVersion(ctx, uint(1)).
					Return(nil).
					Maybe()
			},
			args: args{
				ctx:       context.Background(),
				userID:    uint(1),
				id:        uint(1),
				versionID: uint(1),
			},
			wantErr: true,
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

			err := s.Delete(tt.args.ctx, tt.args.userID, tt.args.id, tt.args.versionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
