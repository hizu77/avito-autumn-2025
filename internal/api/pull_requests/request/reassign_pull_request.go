package request

type ReassignPullRequest struct {
	ID            string `json:"pull_request_id"`
	OldReviewerID string `json:"old_reviewer_id"`
}
