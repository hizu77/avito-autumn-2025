package users

import (
	"github.com/hizu77/avito-autumn-2025/internal/api/httperr"
	"github.com/hizu77/avito-autumn-2025/internal/api/user/response"
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/hizu77/avito-autumn-2025/pkg/utils/collection"
	"github.com/pkg/errors"
)

func mapDomainUserToResponseUser(user model.User) response.User {
	return response.User{
		ID:       user.ID,
		Name:     user.Name,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}
}

func mapDomainUserToResponseSetActive(user model.User) response.SetActive {
	mappedUser := mapDomainUserToResponseUser(user)

	return response.SetActive{
		User: mappedUser,
	}
}

func mapDomainPullRequestToResponseUserReviewRequest(request model.PullRequest) response.ReviewRequest {
	return response.ReviewRequest{
		ID:       request.ID,
		Name:     request.Name,
		AuthorID: request.AuthorID,
		Status:   request.Status,
	}
}

func mapDomainPullRequestsToResponseGetUserReviewRequests(
	userID string,
	requests []model.PullRequest,
) response.GetUserReviewRequests {
	mappedRequests := collection.Map(requests, mapDomainPullRequestToResponseUserReviewRequest)

	return response.GetUserReviewRequests{
		UserID:   userID,
		Requests: mappedRequests,
	}
}

func mapDomainUserErrorToCode(err error) httperr.ErrorCode {
	switch {
	case errors.Is(err, model.ErrUserDoesNotExist):
		return httperr.CodeNotFound
	default:
		return httperr.CodeInternal
	}
}
