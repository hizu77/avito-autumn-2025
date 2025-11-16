package team_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hizu77/avito-autumn-2025/internal/mock/team"
	trmanager "github.com/hizu77/avito-autumn-2025/internal/mock/tr_manager"
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/hizu77/avito-autumn-2025/internal/service/team"
	"github.com/stretchr/testify/require"
)

const (
	testTeamName  = "backend"
	testUserID1   = "user-1"
	testUserID2   = "user-2"
	testUserName1 = "Alice"
	testUserName2 = "Bob"
)

func newService(t *testing.T) (*team.Service, *mock.TeamStorage, *mock.UserStorage) {
	t.Helper()
	ctrl := gomock.NewController(t)
	teamStorage := mock.NewTeamStorage(ctrl)
	userStorage := mock.NewUserStorage(ctrl)
	trManager := trmanager.NewMockTrManager()
	service := team.New(userStorage, teamStorage, trManager)
	return service, teamStorage, userStorage
}

func TestGetTeamByName(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		name string
	}

	tests := []struct {
		name    string
		args    args
		mock    func(storage *mock.TeamStorage)
		want    model.Team
		wantErr error
	}{
		{
			name: "team not found",
			args: args{
				ctx:  context.Background(),
				name: testTeamName,
			},
			mock: func(storage *mock.TeamStorage) {
				storage.EXPECT().GetTeamByName(gomock.Any(), testTeamName).
					Return(model.Team{}, model.ErrTeamDoesNotExist)
			},
			want:    model.Team{},
			wantErr: model.ErrTeamDoesNotExist,
		},
		{
			name: "success",
			args: args{
				ctx:  context.Background(),
				name: testTeamName,
			},
			mock: func(storage *mock.TeamStorage) {
				storage.EXPECT().GetTeamByName(gomock.Any(), testTeamName).
					Return(model.Team{
						Name: testTeamName,
						Members: []model.User{
							{
								ID:       testUserID1,
								Name:     testUserName1,
								TeamName: testTeamName,
								IsActive: true,
							},
							{
								ID:       testUserID2,
								Name:     testUserName2,
								TeamName: testTeamName,
								IsActive: false,
							},
						},
					}, nil)
			},
			want: model.Team{
				Name: testTeamName,
				Members: []model.User{
					{
						ID:       testUserID1,
						Name:     testUserName1,
						TeamName: testTeamName,
						IsActive: true,
					},
					{
						ID:       testUserID2,
						Name:     testUserName2,
						TeamName: testTeamName,
						IsActive: false,
					},
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, teamStorage, _ := newService(t)
			tt.mock(teamStorage)

			got, err := service.GetTeamByName(tt.args.ctx, tt.args.name)

			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestSaveTeam(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		team model.Team
	}

	tests := []struct {
		name    string
		args    args
		mock    func(teamStorage *mock.TeamStorage, userStorage *mock.UserStorage)
		want    model.Team
		wantErr error
	}{
		{
			name: "team already exists",
			args: args{
				ctx: context.Background(),
				team: model.Team{
					Name: testTeamName,
					Members: []model.User{
						{
							ID:       testUserID1,
							Name:     testUserName1,
							TeamName: "",
							IsActive: true,
						},
					},
				},
			},
			mock: func(teamStorage *mock.TeamStorage, userStorage *mock.UserStorage) {
				teamStorage.EXPECT().SaveTeam(gomock.Any(), gomock.Any()).
					Return(model.Team{}, model.ErrTeamAlreadyExists)
			},
			want:    model.Team{},
			wantErr: model.ErrTeamAlreadyExists,
		},
		{
			name: "user storage error",
			args: args{
				ctx: context.Background(),
				team: model.Team{
					Name: testTeamName,
					Members: []model.User{
						{
							ID:       testUserID1,
							Name:     testUserName1,
							TeamName: "",
							IsActive: true,
						},
					},
				},
			},
			mock: func(teamStorage *mock.TeamStorage, userStorage *mock.UserStorage) {
				teamStorage.EXPECT().SaveTeam(gomock.Any(), gomock.Any()).
					Return(model.Team{
						Name: testTeamName,
					}, nil)
				userStorage.EXPECT().SaveUsers(gomock.Any(), gomock.Any()).
					Return(nil, model.ErrUserDoesNotExist)
			},
			want:    model.Team{},
			wantErr: model.ErrUserDoesNotExist,
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				team: model.Team{
					Name: testTeamName,
					Members: []model.User{
						{
							ID:       testUserID1,
							Name:     testUserName1,
							TeamName: "",
							IsActive: true,
						},
						{
							ID:       testUserID2,
							Name:     testUserName2,
							TeamName: "",
							IsActive: false,
						},
					},
				},
			},
			mock: func(teamStorage *mock.TeamStorage, userStorage *mock.UserStorage) {
				teamStorage.EXPECT().SaveTeam(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, team model.Team) (model.Team, error) {
						return model.Team{
							Name: team.Name,
						}, nil
					})
				userStorage.EXPECT().SaveUsers(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, users []model.User) ([]model.User, error) {
						return users, nil
					})
			},
			want: model.Team{
				Name: testTeamName,
				Members: []model.User{
					{
						ID:       testUserID1,
						Name:     testUserName1,
						TeamName: testTeamName,
						IsActive: true,
					},
					{
						ID:       testUserID2,
						Name:     testUserName2,
						TeamName: testTeamName,
						IsActive: false,
					},
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, teamStorage, userStorage := newService(t)
			tt.mock(teamStorage, userStorage)

			got, err := service.SaveTeam(tt.args.ctx, tt.args.team)

			require.Equal(t, tt.want, got)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
