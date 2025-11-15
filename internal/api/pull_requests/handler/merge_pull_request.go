package pullrequest

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/hizu77/avito-autumn-2025/internal/api/common"
	"github.com/hizu77/avito-autumn-2025/internal/api/user/request"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (h *Handler) MergePullRequest(w http.ResponseWriter, r *http.Request) {
	const op = "pullrequest.MergePullRequest"

	var mergePullRequestRequest request.MergePullRequest
	if err := render.DecodeJSON(r.Body, &mergePullRequestRequest); err != nil {
		h.logger.Error("decoding request body",
			zap.String("op", op),
			zap.Error(err),
		)

		common.WriteError(w, r, common.CodeBadRequest, "invalid json body")
		return
	}

	if err := validateMergePullRequestRequest(mergePullRequestRequest); err != nil {
		h.logger.Error("validating request",
			zap.String("op", op),
			zap.Error(err),
		)

		common.WriteError(w, r, common.CodeBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	pullRequest, err := h.service.MergePullRequest(ctx, mergePullRequestRequest.ID)
	if err != nil {
		h.logger.Error("merging pull request",
			zap.String("op", op),
			zap.Error(err),
		)

		code := mapDomainPullRequestErrorToCode(err)
		common.WriteError(w, r, code)
		return
	}

	mappedResponse := mapDomainPullRequestToResponseMergePullRequest(pullRequest)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, mappedResponse)
}

func validateMergePullRequestRequest(req request.MergePullRequest) error {
	if req.ID == "" {
		return errors.New("id required")
	}

	return nil
}
