package pullrequest

import (
	"net/http"

	"github.com/go-chi/render"
	common_response "github.com/hizu77/avito-autumn-2025/internal/api/common/response"
	"github.com/hizu77/avito-autumn-2025/internal/api/pull_requests/request"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (h *Handler) ReassignPullRequest(w http.ResponseWriter, r *http.Request) {
	var reassignPullRequestRequest request.ReassignPullRequest
	if err := render.DecodeJSON(r.Body, &reassignPullRequestRequest); err != nil {
		h.logger.Error("decoding request body", zap.Error(err))

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, common_response.NewBadRequestError("invalid json body"))
		return
	}

	if err := validateReassignPullRequestRequest(reassignPullRequestRequest); err != nil {
		h.logger.Error("validating request", zap.Error(err))

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, common_response.NewBadRequestError(err.Error()))
		return
	}

	ctx := r.Context()
	pullRequest, err := h.service.ReassignPullRequest(
		ctx,
		reassignPullRequestRequest.ID,
		reassignPullRequestRequest.OldReviewerID,
	)
	if err != nil {
		h.logger.Error("reassigning pull request", zap.Error(err))

		mappedErr, code := mapDomainPullRequestErrorToResponseErrorWithStatusCode(err)
		render.Status(r, code)
		render.JSON(w, r, mappedErr)
		return
	}

	mappedResponse := mapDomainPullRequestToResponseReassignPullRequest(
		pullRequest,
		reassignPullRequestRequest.OldReviewerID,
	)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, mappedResponse)
}

func validateReassignPullRequestRequest(req request.ReassignPullRequest) error {
	if req.ID == "" {
		return errors.New("is required")
	}

	if req.OldReviewerID == "" {
		return errors.New("old_reviewer_id is required")
	}

	return nil
}
