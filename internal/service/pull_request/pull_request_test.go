package pullrequest_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mock "github.com/hizu77/avito-autumn-2025/internal/mock/pull_request"
	trmanager "github.com/hizu77/avito-autumn-2025/internal/mock/tr_manager"
	"github.com/hizu77/avito-autumn-2025/internal/model"
	pullrequest "github.com/hizu77/avito-autumn-2025/internal/service/pull_request"
	"github.com/stretchr/testify/require"
)

const (
	testPRID        = "pr-1"
	testPRName      = "Test PR"
	testAuthorID    = "author-1"
	testTeamName    = "backend"
	testUserID1     = "user-1"
	testUserID2     = "user-2"
	testUserID3     = "user-3"
	testUserID4     = "user-4"
	testReviewerID1 = "reviewer-1"
	testReviewerID2 = "reviewer-2"
)

var mockTime = time.Now()

func newService(t *testing.T) (*pullrequest.Service, *mock.TeamStorage, *mock.PullRequestStorage) {
	t.Helper()
	ctrl := gomock.NewController(t)
	teamStorage := mock.NewTeamStorage(ctrl)
	pullRequestStorage := mock.NewPullRequestStorage(ctrl)
	trManager := trmanager.NewMockTrManager()
	service := pullrequest.New(teamStorage, pullRequestStorage, trManager)
	return service, teamStorage, pullRequestStorage
}

func TestCreatePullRequest(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx     context.Context
		request model.PullRequest
	}

	tests := []struct {
		name    string
		args    args
		mock    func(teamStorage *mock.TeamStorage, prStorage *mock.PullRequestStorage)
		want    model.PullRequest
		wantErr error
	}{
		{
			name: "team not found",
			args: args{
				ctx: context.Background(),
				request: model.PullRequest{
					ID:       testPRID,
					Name:     testPRName,
					AuthorID: testAuthorID,
				},
			},
			mock: func(teamStorage *mock.TeamStorage, _ *mock.PullRequestStorage) {
				teamStorage.EXPECT().GetTeamByUserID(gomock.Any(), testAuthorID).
					Return(model.Team{}, model.ErrTeamDoesNotExist)
			},
			want:    model.PullRequest{},
			wantErr: model.ErrTeamDoesNotExist,
		},
		{
			name: "no active teammates - 0 reviewers",
			args: args{
				ctx: context.Background(),
				request: model.PullRequest{
					ID:       testPRID,
					Name:     testPRName,
					AuthorID: testAuthorID,
				},
			},
			mock: func(teamStorage *mock.TeamStorage, prStorage *mock.PullRequestStorage) {
				teamStorage.EXPECT().GetTeamByUserID(gomock.Any(), testAuthorID).
					Return(model.Team{
						Name: testTeamName,
						Members: []model.User{
							{
								ID:       testAuthorID,
								Name:     "Author",
								TeamName: testTeamName,
								IsActive: true,
							},
							{
								ID:       testUserID1,
								Name:     "Inactive User",
								TeamName: testTeamName,
								IsActive: false,
							},
						},
					}, nil)
				prStorage.EXPECT().InsertPullRequest(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, pr model.PullRequest) (model.PullRequest, error) {
						return pr, nil
					})
			},
			want: model.PullRequest{
				ID:           testPRID,
				Name:         testPRName,
				AuthorID:     testAuthorID,
				Status:       model.StatusOpen,
				ReviewersIDs: []string{},
			},
			wantErr: nil,
		},
		{
			name: "one active teammate - 1 reviewer",
			args: args{
				ctx: context.Background(),
				request: model.PullRequest{
					ID:       testPRID,
					Name:     testPRName,
					AuthorID: testAuthorID,
				},
			},
			mock: func(teamStorage *mock.TeamStorage, prStorage *mock.PullRequestStorage) {
				teamStorage.EXPECT().GetTeamByUserID(gomock.Any(), testAuthorID).
					Return(model.Team{
						Name: testTeamName,
						Members: []model.User{
							{
								ID:       testAuthorID,
								Name:     "Author",
								TeamName: testTeamName,
								IsActive: true,
							},
							{
								ID:       testUserID1,
								Name:     "Active User",
								TeamName: testTeamName,
								IsActive: true,
							},
						},
					}, nil)
				prStorage.EXPECT().InsertPullRequest(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, pr model.PullRequest) (model.PullRequest, error) {
						return pr, nil
					})
			},
			want: model.PullRequest{
				ID:           testPRID,
				Name:         testPRName,
				AuthorID:     testAuthorID,
				Status:       model.StatusOpen,
				ReviewersIDs: []string{testUserID1},
			},
			wantErr: nil,
		},
		{
			name: "two active teammates - 2 reviewers",
			args: args{
				ctx: context.Background(),
				request: model.PullRequest{
					ID:       testPRID,
					Name:     testPRName,
					AuthorID: testAuthorID,
				},
			},
			mock: func(teamStorage *mock.TeamStorage, prStorage *mock.PullRequestStorage) {
				teamStorage.EXPECT().GetTeamByUserID(gomock.Any(), testAuthorID).
					Return(model.Team{
						Name: testTeamName,
						Members: []model.User{
							{
								ID:       testAuthorID,
								Name:     "Author",
								TeamName: testTeamName,
								IsActive: true,
							},
							{
								ID:       testUserID1,
								Name:     "Active User 1",
								TeamName: testTeamName,
								IsActive: true,
							},
							{
								ID:       testUserID2,
								Name:     "Active User 2",
								TeamName: testTeamName,
								IsActive: true,
							},
						},
					}, nil)
				prStorage.EXPECT().InsertPullRequest(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, pr model.PullRequest) (model.PullRequest, error) {
						return pr, nil
					})
			},
			want: model.PullRequest{
				ID:           testPRID,
				Name:         testPRName,
				AuthorID:     testAuthorID,
				Status:       model.StatusOpen,
				ReviewersIDs: []string{testUserID1, testUserID2},
			},
			wantErr: nil,
		},
		{
			name: "more than two active teammates - 2 reviewers selected",
			args: args{
				ctx: context.Background(),
				request: model.PullRequest{
					ID:       testPRID,
					Name:     testPRName,
					AuthorID: testAuthorID,
				},
			},
			mock: func(teamStorage *mock.TeamStorage, prStorage *mock.PullRequestStorage) {
				teamStorage.EXPECT().GetTeamByUserID(gomock.Any(), testAuthorID).
					Return(model.Team{
						Name: testTeamName,
						Members: []model.User{
							{
								ID:       testAuthorID,
								Name:     "Author",
								TeamName: testTeamName,
								IsActive: true,
							},
							{
								ID:       testUserID1,
								Name:     "Active User 1",
								TeamName: testTeamName,
								IsActive: true,
							},
							{
								ID:       testUserID2,
								Name:     "Active User 2",
								TeamName: testTeamName,
								IsActive: true,
							},
							{
								ID:       testUserID3,
								Name:     "Active User 3",
								TeamName: testTeamName,
								IsActive: true,
							},
						},
					}, nil)
				prStorage.EXPECT().InsertPullRequest(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, pr model.PullRequest) (model.PullRequest, error) {
						return pr, nil
					})
			},
			want: model.PullRequest{
				ID:           testPRID,
				Name:         testPRName,
				AuthorID:     testAuthorID,
				Status:       model.StatusOpen,
				ReviewersIDs: []string{}, // its random so empty
			},
			wantErr: nil,
		},
		{
			name: "inactive users excluded",
			args: args{
				ctx: context.Background(),
				request: model.PullRequest{
					ID:       testPRID,
					Name:     testPRName,
					AuthorID: testAuthorID,
				},
			},
			mock: func(teamStorage *mock.TeamStorage, prStorage *mock.PullRequestStorage) {
				teamStorage.EXPECT().GetTeamByUserID(gomock.Any(), testAuthorID).
					Return(model.Team{
						Name: testTeamName,
						Members: []model.User{
							{
								ID:       testAuthorID,
								Name:     "Author",
								TeamName: testTeamName,
								IsActive: true,
							},
							{
								ID:       testUserID1,
								Name:     "Active User",
								TeamName: testTeamName,
								IsActive: true,
							},
							{
								ID:       testUserID2,
								Name:     "Inactive User",
								TeamName: testTeamName,
								IsActive: false,
							},
						},
					}, nil)
				prStorage.EXPECT().InsertPullRequest(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, pr model.PullRequest) (model.PullRequest, error) {
						return pr, nil
					})
			},
			want: model.PullRequest{
				ID:           testPRID,
				Name:         testPRName,
				AuthorID:     testAuthorID,
				Status:       model.StatusOpen,
				ReviewersIDs: []string{testUserID1},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, teamStorage, prStorage := newService(t)
			tt.mock(teamStorage, prStorage)

			got, err := service.CreatePullRequest(tt.args.ctx, tt.args.request)

			require.ErrorIs(t, err, tt.wantErr)

			if tt.wantErr == nil {
				require.Equal(t, tt.want.ID, got.ID)
				require.Equal(t, tt.want.Name, got.Name)
				require.Equal(t, tt.want.AuthorID, got.AuthorID)
				require.Equal(t, tt.want.Status, got.Status)
				if len(tt.want.ReviewersIDs) > 0 {
					require.ElementsMatch(t, tt.want.ReviewersIDs, got.ReviewersIDs)
				} else {
					require.LessOrEqual(t, len(got.ReviewersIDs), 2)
					require.NotContains(t, got.ReviewersIDs, tt.args.request.AuthorID)
				}
				require.NotNil(t, got.CreatedAt)
				require.Nil(t, got.MergedAt)
			} else {
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestMergePullRequest(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		id  string
	}

	tests := []struct {
		name    string
		args    args
		mock    func(prStorage *mock.PullRequestStorage)
		want    model.PullRequest
		wantErr error
	}{
		{
			name: "pull request not found",
			args: args{
				ctx: context.Background(),
				id:  testPRID,
			},
			mock: func(prStorage *mock.PullRequestStorage) {
				prStorage.EXPECT().GetPullRequestByID(gomock.Any(), testPRID).
					Return(model.PullRequest{}, model.ErrPullRequestDoesNotExist)
			},
			want:    model.PullRequest{},
			wantErr: model.ErrPullRequestDoesNotExist,
		},
		{
			name: "success - merge open PR",
			args: args{
				ctx: context.Background(),
				id:  testPRID,
			},
			mock: func(prStorage *mock.PullRequestStorage) {
				prStorage.EXPECT().GetPullRequestByID(gomock.Any(), testPRID).
					Return(model.PullRequest{
						ID:           testPRID,
						Name:         testPRName,
						AuthorID:     testAuthorID,
						Status:       model.StatusOpen,
						ReviewersIDs: []string{testReviewerID1},
						CreatedAt:    &mockTime,
						MergedAt:     nil,
					}, nil)
				prStorage.EXPECT().UpdatePullRequestInfo(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, pr model.PullRequest) (model.PullRequest, error) {
						return pr, nil
					})
			},
			want: model.PullRequest{
				ID:           testPRID,
				Name:         testPRName,
				AuthorID:     testAuthorID,
				Status:       model.StatusMerged,
				ReviewersIDs: []string{testReviewerID1},
				CreatedAt:    &mockTime,
			},
			wantErr: nil,
		},
		{
			name: "idempotent - already merged",
			args: args{
				ctx: context.Background(),
				id:  testPRID,
			},
			mock: func(prStorage *mock.PullRequestStorage) {
				mergedTime := mockTime.Add(-time.Hour)
				prStorage.EXPECT().GetPullRequestByID(gomock.Any(), testPRID).
					Return(model.PullRequest{
						ID:           testPRID,
						Name:         testPRName,
						AuthorID:     testAuthorID,
						Status:       model.StatusMerged,
						ReviewersIDs: []string{testReviewerID1},
						CreatedAt:    &mockTime,
						MergedAt:     &mergedTime,
					}, nil)
			},
			want: model.PullRequest{
				ID:           testPRID,
				Name:         testPRName,
				AuthorID:     testAuthorID,
				Status:       model.StatusMerged,
				ReviewersIDs: []string{testReviewerID1},
				CreatedAt:    &mockTime,
				MergedAt:     func() *time.Time { t := mockTime.Add(-time.Hour); return &t }(),
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, _, prStorage := newService(t)
			tt.mock(prStorage)

			got, err := service.MergePullRequest(tt.args.ctx, tt.args.id)

			require.ErrorIs(t, err, tt.wantErr)

			require.Equal(t, tt.want.ID, got.ID)
			require.Equal(t, tt.want.Name, got.Name)
			require.Equal(t, tt.want.AuthorID, got.AuthorID)
			require.Equal(t, tt.want.Status, got.Status)
			require.Equal(t, tt.want.ReviewersIDs, got.ReviewersIDs)
			require.Equal(t, tt.want.CreatedAt, got.CreatedAt)
		})
	}
}

func TestReassignPullRequest(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx        context.Context
		id         string
		reviewerID string
	}

	tests := []struct {
		name    string
		args    args
		mock    func(teamStorage *mock.TeamStorage, prStorage *mock.PullRequestStorage)
		want    model.ReassignedPullRequest
		wantErr error
	}{
		{
			name: "pull request not found",
			args: args{
				ctx:        context.Background(),
				id:         testPRID,
				reviewerID: testReviewerID1,
			},
			mock: func(_ *mock.TeamStorage, prStorage *mock.PullRequestStorage) {
				prStorage.EXPECT().GetPullRequestByID(gomock.Any(), testPRID).
					Return(model.PullRequest{}, model.ErrPullRequestDoesNotExist)
			},
			want:    model.ReassignedPullRequest{},
			wantErr: model.ErrPullRequestDoesNotExist,
		},
		{
			name: "pull request is merged - cannot reassign",
			args: args{
				ctx:        context.Background(),
				id:         testPRID,
				reviewerID: testReviewerID1,
			},
			mock: func(_ *mock.TeamStorage, prStorage *mock.PullRequestStorage) {
				mergedTime := mockTime.Add(-time.Hour)
				prStorage.EXPECT().GetPullRequestByID(gomock.Any(), testPRID).
					Return(model.PullRequest{
						ID:           testPRID,
						Name:         testPRName,
						AuthorID:     testAuthorID,
						Status:       model.StatusMerged,
						ReviewersIDs: []string{testReviewerID1},
						CreatedAt:    &mockTime,
						MergedAt:     &mergedTime,
					}, nil)
			},
			want:    model.ReassignedPullRequest{},
			wantErr: model.ErrPullRequestIsMerged,
		},
		{
			name: "reviewer not assigned",
			args: args{
				ctx:        context.Background(),
				id:         testPRID,
				reviewerID: testUserID1,
			},
			mock: func(_ *mock.TeamStorage, prStorage *mock.PullRequestStorage) {
				prStorage.EXPECT().GetPullRequestByID(gomock.Any(), testPRID).
					Return(model.PullRequest{
						ID:           testPRID,
						Name:         testPRName,
						AuthorID:     testAuthorID,
						Status:       model.StatusOpen,
						ReviewersIDs: []string{testReviewerID1, testReviewerID2},
						CreatedAt:    &mockTime,
					}, nil)
			},
			want:    model.ReassignedPullRequest{},
			wantErr: model.ErrReviewerNotAssign,
		},
		{
			name: "no candidate to reassign",
			args: args{
				ctx:        context.Background(),
				id:         testPRID,
				reviewerID: testReviewerID1,
			},
			mock: func(teamStorage *mock.TeamStorage, prStorage *mock.PullRequestStorage) {
				prStorage.EXPECT().GetPullRequestByID(gomock.Any(), testPRID).
					Return(model.PullRequest{
						ID:           testPRID,
						Name:         testPRName,
						AuthorID:     testAuthorID,
						Status:       model.StatusOpen,
						ReviewersIDs: []string{testReviewerID1},
						CreatedAt:    &mockTime,
					}, nil)
				teamStorage.EXPECT().GetTeamByUserID(gomock.Any(), testReviewerID1).
					Return(model.Team{
						Name: testTeamName,
						Members: []model.User{
							{
								ID:       testReviewerID1,
								Name:     "Reviewer",
								TeamName: testTeamName,
								IsActive: true,
							},
						},
					}, nil)
			},
			want:    model.ReassignedPullRequest{},
			wantErr: model.ErrNoCandidate,
		},
		{
			name: "success - reassign reviewer",
			args: args{
				ctx:        context.Background(),
				id:         testPRID,
				reviewerID: testReviewerID1,
			},
			mock: func(teamStorage *mock.TeamStorage, prStorage *mock.PullRequestStorage) {
				prStorage.EXPECT().GetPullRequestByID(gomock.Any(), testPRID).
					Return(model.PullRequest{
						ID:           testPRID,
						Name:         testPRName,
						AuthorID:     testAuthorID,
						Status:       model.StatusOpen,
						ReviewersIDs: []string{testReviewerID1, testReviewerID2},
						CreatedAt:    &mockTime,
					}, nil)

				teamStorage.EXPECT().GetTeamByUserID(gomock.Any(), testReviewerID1).
					Return(model.Team{
						Name: testTeamName,
						Members: []model.User{
							{
								ID:       testReviewerID1,
								Name:     "Reviewer 1",
								TeamName: testTeamName,
								IsActive: true,
							},
							{
								ID:       testUserID1,
								Name:     "Active User",
								TeamName: testTeamName,
								IsActive: true,
							},
							{
								ID:       testUserID2,
								Name:     "Inactive User",
								TeamName: testTeamName,
								IsActive: false,
							},
						},
					}, nil)

				prStorage.EXPECT().UpdatePullRequestReviewers(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, pr model.PullRequest) (model.PullRequest, error) {
						return pr, nil
					})
			},
			want: model.ReassignedPullRequest{
				ID:           testPRID,
				Name:         testPRName,
				AuthorID:     testAuthorID,
				Status:       model.StatusOpen,
				ReviewersIDs: []string{testUserID1, testReviewerID2},
				CreatedAt:    &mockTime,
				ReassignedBy: testUserID1,
			},
			wantErr: nil,
		},
		{
			name: "all team members already reviewers - no candidate",
			args: args{
				ctx:        context.Background(),
				id:         testPRID,
				reviewerID: testReviewerID1,
			},
			mock: func(teamStorage *mock.TeamStorage, prStorage *mock.PullRequestStorage) {
				prStorage.EXPECT().GetPullRequestByID(gomock.Any(), testPRID).
					Return(model.PullRequest{
						ID:           testPRID,
						Name:         testPRName,
						AuthorID:     testAuthorID,
						Status:       model.StatusOpen,
						ReviewersIDs: []string{testReviewerID1, testUserID1},
						CreatedAt:    &mockTime,
					}, nil)
				teamStorage.EXPECT().GetTeamByUserID(gomock.Any(), testReviewerID1).
					Return(model.Team{
						Name: testTeamName,
						Members: []model.User{
							{
								ID:       testReviewerID1,
								Name:     "Reviewer 1",
								TeamName: testTeamName,
								IsActive: true,
							},
							{
								ID:       testUserID1,
								Name:     "User 1",
								TeamName: testTeamName,
								IsActive: true,
							},
						},
					}, nil)
			},
			want:    model.ReassignedPullRequest{},
			wantErr: model.ErrNoCandidate,
		},
		{
			name: "inactive users excluded from candidates",
			args: args{
				ctx:        context.Background(),
				id:         testPRID,
				reviewerID: testReviewerID1,
			},
			mock: func(teamStorage *mock.TeamStorage, prStorage *mock.PullRequestStorage) {
				prStorage.EXPECT().GetPullRequestByID(gomock.Any(), testPRID).
					Return(model.PullRequest{
						ID:           testPRID,
						Name:         testPRName,
						AuthorID:     testAuthorID,
						Status:       model.StatusOpen,
						ReviewersIDs: []string{testReviewerID1},
						CreatedAt:    &mockTime,
					}, nil)
				teamStorage.EXPECT().GetTeamByUserID(gomock.Any(), testReviewerID1).
					Return(model.Team{
						Name: testTeamName,
						Members: []model.User{
							{
								ID:       testReviewerID1,
								Name:     "Reviewer 1",
								TeamName: testTeamName,
								IsActive: true,
							},
							{
								ID:       testUserID1,
								Name:     "Inactive User",
								TeamName: testTeamName,
								IsActive: false,
							},
						},
					}, nil)
			},
			want:    model.ReassignedPullRequest{},
			wantErr: model.ErrNoCandidate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, teamStorage, prStorage := newService(t)
			tt.mock(teamStorage, prStorage)

			got, err := service.ReassignPullRequest(tt.args.ctx, tt.args.id, tt.args.reviewerID)

			require.ErrorIs(t, err, tt.wantErr)

			if tt.wantErr == nil {
				require.Equal(t, tt.want.ID, got.ID)
				require.Equal(t, tt.want.Name, got.Name)
				require.Equal(t, tt.want.AuthorID, got.AuthorID)
				require.Equal(t, tt.want.Status, got.Status)
				require.ElementsMatch(t, tt.want.ReviewersIDs, got.ReviewersIDs)
				require.Equal(t, tt.want.ReassignedBy, got.ReassignedBy)
			} else {
				require.Equal(t, model.ReassignedPullRequest{}, got)
			}
		})
	}
}
