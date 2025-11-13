package users

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/hizu77/avito-autumn-2025/internal/api/common/response"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	idQueryParam = "user_id"
)

func (h *Handler) GetUserReviewRequests(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(idQueryParam)

	if err := validateUserID(id); err != nil {
		h.logger.Error("validate user id", zap.Error(err))

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.NewBadRequestError(err.Error()))
		return
	}

	ctx := r.Context()
	requests, err := h.service.GetUserReviewRequests(ctx, id)
	if err != nil {
		h.logger.Error("getting user review requests", zap.Error(err))

		mappedErr, code := mapDomainUserErrorToResponseErrorWithStatusCode(err)
		render.Status(r, code)
		render.JSON(w, r, mappedErr)
		return
	}

	mappedRequests := mapDomainPullRequestsToResponseGetUserReviewRequests(
		id,
		requests,
	)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, mappedRequests)
}

func validateUserID(id string) error {
	if id == "" {
		return errors.New("invalid user id")
	}

	return nil
}
