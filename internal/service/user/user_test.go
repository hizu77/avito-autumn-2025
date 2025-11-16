package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mock "github.com/hizu77/avito-autumn-2025/internal/mock/user"
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/hizu77/avito-autumn-2025/internal/service/user"
	"github.com/stretchr/testify/require"
)

const (
	testUserID   = "user-1"
	testUserName = "Alice"
	testTeamName = "backend"
	testPRID1    = "pr-1"
	testPRID2    = "pr-2"
	testAuthorID = "author-1"
)

var mockTime = time.Now()

func newService(t *testing.T) (*user.Service, *mock.UserStorage, *mock.PullRequestStorage) {
	t.Helper()
	ctrl := gomock.NewController(t)
	userStorage := mock.NewUserStorage(ctrl)
	pullRequestStorage := mock.NewPullRequestStorage(ctrl)
	service := user.New(userStorage, pullRequestStorage)
	return service, userStorage, pullRequestStorage
}

func TestSetActive(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx    context.Context
		id     string
		active bool
	}

	tests := []struct {
		name    string
		args    args
		mock    func(storage *mock.UserStorage)
		want    model.User
		wantErr error
	}{
		{
			name: "user not found",
			args: args{
				ctx:    context.Background(),
				id:     testUserID,
				active: true,
			},
			mock: func(storage *mock.UserStorage) {
				storage.EXPECT().UpdateActivity(gomock.Any(), testUserID, true).
					Return(model.User{}, model.ErrUserDoesNotExist)
			},
			want:    model.User{},
			wantErr: model.ErrUserDoesNotExist,
		},
		{
			name: "success - set active",
			args: args{
				ctx:    context.Background(),
				id:     testUserID,
				active: true,
			},
			mock: func(storage *mock.UserStorage) {
				storage.EXPECT().UpdateActivity(gomock.Any(), testUserID, true).
					Return(model.User{
						ID:       testUserID,
						Name:     testUserName,
						TeamName: testTeamName,
						IsActive: true,
					}, nil)
			},
			want: model.User{
				ID:       testUserID,
				Name:     testUserName,
				TeamName: testTeamName,
				IsActive: true,
			},
			wantErr: nil,
		},
		{
			name: "success - set inactive",
			args: args{
				ctx:    context.Background(),
				id:     testUserID,
				active: false,
			},
			mock: func(storage *mock.UserStorage) {
				storage.EXPECT().UpdateActivity(gomock.Any(), testUserID, false).
					Return(model.User{
						ID:       testUserID,
						Name:     testUserName,
						TeamName: testTeamName,
						IsActive: false,
					}, nil)
			},
			want: model.User{
				ID:       testUserID,
				Name:     testUserName,
				TeamName: testTeamName,
				IsActive: false,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, userStorage, _ := newService(t)
			tt.mock(userStorage)

			got, err := service.SetActive(tt.args.ctx, tt.args.id, tt.args.active)

			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestGetUserReviewRequests(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		id  string
	}

	tests := []struct {
		name    string
		args    args
		mock    func(storage *mock.PullRequestStorage)
		want    []model.PullRequest
		wantErr error
	}{
		{
			name: "success - empty list",
			args: args{
				ctx: context.Background(),
				id:  testUserID,
			},
			mock: func(storage *mock.PullRequestStorage) {
				storage.EXPECT().GetPullRequestsByReviewer(gomock.Any(), testUserID).
					Return([]model.PullRequest{}, nil)
			},
			want:    []model.PullRequest{},
			wantErr: nil,
		},
		{
			name: "success - with pull requests",
			args: args{
				ctx: context.Background(),
				id:  testUserID,
			},
			mock: func(storage *mock.PullRequestStorage) {
				storage.EXPECT().GetPullRequestsByReviewer(gomock.Any(), testUserID).
					Return([]model.PullRequest{
						{
							ID:           testPRID1,
							Name:         "PR 1",
							AuthorID:     testAuthorID,
							Status:       model.StatusOpen,
							ReviewersIDs: []string{testUserID},
							CreatedAt:    &mockTime,
							MergedAt:     nil,
						},
						{
							ID:           testPRID2,
							Name:         "PR 2",
							AuthorID:     testAuthorID,
							Status:       model.StatusOpen,
							ReviewersIDs: []string{testUserID},
							CreatedAt:    &mockTime,
							MergedAt:     nil,
						},
					}, nil)
			},
			want: []model.PullRequest{
				{
					ID:           testPRID1,
					Name:         "PR 1",
					AuthorID:     testAuthorID,
					Status:       model.StatusOpen,
					ReviewersIDs: []string{testUserID},
					CreatedAt:    &mockTime,
					MergedAt:     nil,
				},
				{
					ID:           testPRID2,
					Name:         "PR 2",
					AuthorID:     testAuthorID,
					Status:       model.StatusOpen,
					ReviewersIDs: []string{testUserID},
					CreatedAt:    &mockTime,
					MergedAt:     nil,
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, _, pullRequestStorage := newService(t)
			tt.mock(pullRequestStorage)

			got, err := service.GetUserReviewRequests(tt.args.ctx, tt.args.id)

			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)
		})
	}
}
