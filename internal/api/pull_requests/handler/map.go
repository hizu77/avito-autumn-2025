package pullrequest

import (
	"time"

	"github.com/hizu77/avito-autumn-2025/internal/api/httperr"
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

func mapDomainReassignedPullRequestToResponsePullRequest(
	req model.ReassignedPullRequest,
) response.PullRequest {
	return response.PullRequest{
		ID:        req.ID,
		Name:      req.Name,
		AuthorID:  req.AuthorID,
		Status:    req.Status,
		Reviewers: req.ReviewersIDs,
	}
}

func mapDomainReassignedPullRequestToResponseReassignPullRequest(
	req model.ReassignedPullRequest,
) response.ReassignPullRequest {
	mappedPullRequest := mapDomainReassignedPullRequestToResponsePullRequest(req)

	return response.ReassignPullRequest{
		PullRequest: mappedPullRequest,
		ReplacedBy:  req.ReassignedBy,
	}
}

func mapDomainPullRequestErrorToCode(err error) httperr.ErrorCode {
	switch {
	case errors.Is(err, model.ErrPullRequestIsMerged):
		return httperr.CodePrMerged
	case errors.Is(err, model.ErrPullRequestAlreadyExists):
		return httperr.CodePrExists
	case errors.Is(err, model.ErrPullRequestDoesNotExist):
		return httperr.CodeNotFound
	case errors.Is(err, model.ErrTeamDoesNotExist):
		return httperr.CodeNotFound
	case errors.Is(err, model.ErrUserDoesNotExist):
		return httperr.CodeNotFound
	case errors.Is(err, model.ErrReviewerNotAssign):
		return httperr.CodeNotAssigned
	case errors.Is(err, model.ErrNoCandidate):
		return httperr.CodeNoCandidate
	default:
		return httperr.CodeInternal
	}
}
