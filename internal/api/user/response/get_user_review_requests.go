package response

type GetUserReviewRequests struct {
	UserID   string          `json:"user_id"`
	Requests []ReviewRequest `json:"pull_requests"`
}
