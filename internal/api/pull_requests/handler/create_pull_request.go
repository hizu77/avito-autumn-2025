package pullrequest

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/hizu77/avito-autumn-2025/internal/api/common"
	"github.com/hizu77/avito-autumn-2025/internal/api/pull_requests/request"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (h *Handler) CreatePullRequest(w http.ResponseWriter, r *http.Request) {
	const op = "pullrequest.CreatePullRequest"

	var createPullRequestRequest request.CreatePullRequest
	if err := render.DecodeJSON(r.Body, &createPullRequestRequest); err != nil {
		h.logger.Error("decoding request body",
			zap.String("op", op),
			zap.Error(err),
		)

		common.WriteError(w, r, common.CodeBadRequest, "invalid json body")
		return
	}

	if err := validateCreatePullRequestRequest(createPullRequestRequest); err != nil {
		h.logger.Error("validate request",
			zap.String("op", op),
			zap.Error(err),
		)

		common.WriteError(w, r, common.CodeBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	mappedPullRequest := mapRequestCreatePullRequestToDomainPullRequest(createPullRequestRequest)

	pullRequest, err := h.service.CreatePullRequest(ctx, mappedPullRequest)
	if err != nil {
		h.logger.Error("creating pull request",
			zap.String("op", op),
			zap.Error(err),
		)

		code := mapDomainPullRequestErrorToCode(err)
		common.WriteError(w, r, code)
		return
	}

	mappedResponse := mapDomainPullRequestToResponseCreatePullRequest(pullRequest)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, mappedResponse)
}

func validateCreatePullRequestRequest(req request.CreatePullRequest) error {
	if req.ID == "" {
		return errors.New("id is required")
	}

	if req.Name == "" {
		return errors.New("name is required")
	}

	if req.AuthorID == "" {
		return errors.New("author_id is required")
	}

	return nil
}
