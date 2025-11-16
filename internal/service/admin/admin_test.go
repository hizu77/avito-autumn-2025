package admin_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hizu77/avito-autumn-2025/internal/mock/admin"
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/hizu77/avito-autumn-2025/internal/service/admin"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

const (
	testAdminID   = "4d4a8cd8-501b-4bd4-8589-6be8dcca7c09"
	testPassword  = "secret123"
	wrongPassword = "wrongpass"

	testJWTSecret = "super-secret-for-tests"
)

func newService(t *testing.T) (*admin.Service, *mock.AdminStorage) {
	t.Helper()
	ctrl := gomock.NewController(t)
	storage := mock.NewAdminStorage(ctrl)
	service := admin.New(storage, []byte(testJWTSecret))
	return service, storage
}

func TestLoginAdmin(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx      context.Context
		id       string
		password string
	}

	tests := []struct {
		name    string
		args    args
		mock    func(storage *mock.AdminStorage)
		want    string
		wantErr error
	}{
		{
			name: "admin not found",
			args: args{
				ctx:      context.Background(),
				id:       testAdminID,
				password: testPassword,
			},
			mock: func(storage *mock.AdminStorage) {
				storage.EXPECT().GetAdmin(gomock.Any(), testAdminID).
					Return(model.Admin{}, model.ErrAdminDoesNotExist)
			},
			want:    "",
			wantErr: model.ErrAdminDoesNotExist,
		},
		{
			name: "invalid password",
			args: args{
				ctx:      context.Background(),
				id:       testAdminID,
				password: wrongPassword,
			},
			mock: func(storage *mock.AdminStorage) {
				hash, _ := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
				storage.EXPECT().GetAdmin(gomock.Any(), testAdminID).
					Return(model.Admin{ID: testAdminID, PasswordHash: string(hash)}, nil)
			},
			want:    "",
			wantErr: model.ErrInvalidAdminPassword,
		},
		{
			name: "success",
			args: args{
				ctx:      context.Background(),
				id:       testAdminID,
				password: testPassword,
			},
			mock: func(storage *mock.AdminStorage) {
				hash, _ := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
				storage.EXPECT().GetAdmin(gomock.Any(), testAdminID).
					Return(model.Admin{ID: testAdminID, PasswordHash: string(hash)}, nil)
			},
			want:    "",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, storage := newService(t)
			tt.mock(storage)

			got, err := service.LoginAdmin(tt.args.ctx, tt.args.id, tt.args.password)

			if tt.wantErr != nil {
				require.Equal(t, tt.want, got)
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			if tt.wantErr == nil && tt.want == "" {
				require.NotEmpty(t, got)
			} else {
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestRegisterAdmin(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx      context.Context
		id       string
		password string
	}

	tests := []struct {
		name    string
		args    args
		mock    func(storage *mock.AdminStorage)
		want    model.Admin
		wantErr error
	}{
		{
			name: "admin already exists",
			args: args{
				ctx:      context.Background(),
				id:       testAdminID,
				password: testPassword,
			},
			mock: func(storage *mock.AdminStorage) {
				storage.EXPECT().InsertAdmin(gomock.Any(), gomock.Any()).
					Return(model.Admin{}, model.ErrAdminAlreadyExists)
			},
			want:    model.Admin{},
			wantErr: model.ErrAdminAlreadyExists,
		},
		{
			name: "success",
			args: args{
				ctx:      context.Background(),
				id:       testAdminID,
				password: testPassword,
			},
			mock: func(storage *mock.AdminStorage) {
				storage.EXPECT().InsertAdmin(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, a model.Admin) (model.Admin, error) {
						return a, nil
					})
			},
			want: model.Admin{
				ID: testAdminID,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, storage := newService(t)
			tt.mock(storage)

			got, err := service.RegisterAdmin(tt.args.ctx, tt.args.id, tt.args.password)

			require.ErrorIs(t, err, tt.wantErr)
			if tt.wantErr == nil {
				require.Equal(t, tt.want.ID, got.ID)
				require.NotEmpty(t, got.PasswordHash)
			} else {
				require.Equal(t, tt.want, got)
			}
		})
	}
}
