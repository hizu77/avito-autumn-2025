package pullrequest

import (
	"net/http"
	"time"

	common_response "github.com/hizu77/avito-autumn-2025/internal/api/common/response"
	"github.com/hizu77/avito-autumn-2025/internal/api/pull_requests/request"
	"github.com/hizu77/avito-autumn-2025/internal/api/pull_requests/response"
	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/pkg/errors"
)

func mapRequestCreatePullRequestToDomainPullRequest(req request.CreatePullRequest) model.PullRequest {
	return model.PullRequest{
		ID:       req.ID,
		Name:     req.Name,
		AuthorID: req.AuthorID,
	}
}

func mapDomainPullRequestToResponsePullRequest(req model.PullRequest) response.PullRequest {
	return response.PullRequest{
		ID:        req.ID,
		Name:      req.Name,
		AuthorID:  req.AuthorID,
		Status:    req.Status,
		Reviewers: req.ReviewersIDs,
	}
}

func mapDomainPullRequestToResponseCreatePullRequest(req model.PullRequest) response.CreatePullRequest {
	mappedPullRequest := mapDomainPullRequestToResponsePullRequest(req)

	return response.CreatePullRequest{
		PullRequest: mappedPullRequest,
	}
}

func mapDomainPullRequestToResponseMergedPullRequest(
	req model.PullRequest,
) response.MergedPullRequest {
	var mergedAt time.Time
	if req.MergedAt != nil {
		mergedAt = *req.MergedAt
	}
	return response.MergedPullRequest{
		ID:        req.ID,
		Name:      req.Name,
		AuthorID:  req.AuthorID,
		Status:    req.Status,
		Reviewers: req.ReviewersIDs,
		MergedAt:  mergedAt,
	}
}

func mapDomainPullRequestToResponseMergePullRequest(
	req model.PullRequest,
) response.MergePullRequest {
	mappedPullRequest := mapDomainPullRequestToResponseMergedPullRequest(req)

	return response.MergePullRequest{
		MergedPullRequest: mappedPullRequest,
	}
}

func mapDomainPullRequestToResponseReassignPullRequest(
	req model.PullRequest,
	replacedBy string,
) response.ReassignPullRequest {
	mappedPullRequest := mapDomainPullRequestToResponsePullRequest(req)

	return response.ReassignPullRequest{
		PullRequest: mappedPullRequest,
		ReplacedBy:  replacedBy,
	}
}

func mapDomainPullRequestErrorToResponseErrorWithStatusCode(err error) (common_response.Error, int) {
	switch {
	case errors.Is(err, model.ErrPullRequestIsMerged):
		return common_response.NewPRMergedError(), http.StatusConflict
	case errors.Is(err, model.ErrPullRequestAlreadyExists):
		return common_response.NewPRExistsError(), http.StatusConflict
	case errors.Is(err, model.ErrPullRequestDoesNotExist):
		return common_response.NewNotFoundError(), http.StatusNotFound
	case errors.Is(err, model.ErrTeamDoesNotExist):
		return common_response.NewNotFoundError(), http.StatusNotFound
	case errors.Is(err, model.ErrUserDoesNotExist):
		return common_response.NewNotFoundError(), http.StatusNotFound
	case errors.Is(err, model.ErrReviewerNotAssign):
		return common_response.NewNotAssignedError(), http.StatusConflict
	case errors.Is(err, model.ErrNoCandidate):
		return common_response.NewNoCandidateError(), http.StatusConflict
	default:
		return common_response.NewInternalServerError(), http.StatusInternalServerError
	}
}
