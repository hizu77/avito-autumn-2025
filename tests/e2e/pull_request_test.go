package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// Basic create flow: PR appears with OPEN status and correct fields.
func TestPR_Create_Succeeds(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()
	tn := uniqueID("e2e-pr-open")
	author := "u1-" + tn

	// prepare team with author
	status, body := post(t, base+teamAddPath, map[string]any{
		"team_name": tn,
		"members": []any{
			map[string]any{"user_id": author, "username": "author", "is_active": true},
		},
	}, nil)
	require.Equal(t, http.StatusCreated, status, string(body))

	// create PR
	status, body = post(t, base+prCreatePath, map[string]any{
		"pull_request_id":   "pr-" + tn,
		"pull_request_name": "open-endpoint",
		"author_id":         author,
	}, nil)
	require.Equal(t, http.StatusCreated, status, string(body))

	var resp map[string]any
	require.NoError(t, json.Unmarshal(body, &resp))
	pr := asMap(t, resp["pr"])
	require.Equal(t, "pr-"+tn, getString(t, pr, "pull_request_id"))
	require.Equal(t, "open-endpoint", getString(t, pr, "pull_request_name"))
	require.Equal(t, author, getString(t, pr, "author_id"))
	require.Equal(t, "OPEN", getString(t, pr, "status"))
}

// Auto-assigns exactly two active reviewers when available; excludes author and inactive users.
func TestPR_Create_AutoAssignTwoReviewers(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()

	tn := uniqueID("e2e-pr-two")
	author := "u1-" + tn
	r1 := "u2-" + tn
	r2 := "u3-" + tn
	inactive := "u4-" + tn

	status, body := post(t, base+teamAddPath, map[string]any{
		"team_name": tn,
		"members": []any{
			map[string]any{"user_id": author, "username": "author", "is_active": true},
			map[string]any{"user_id": r1, "username": "r1", "is_active": true},
			map[string]any{"user_id": r2, "username": "r2", "is_active": true},
			map[string]any{"user_id": inactive, "username": "r3", "is_active": false},
		},
	}, nil)
	require.Equal(t, http.StatusCreated, status, string(body))

	status, body = post(t, base+prCreatePath, map[string]any{
		"pull_request_id":   "pr-" + tn,
		"pull_request_name": "add feature",
		"author_id":         author,
	}, nil)
	require.Equal(t, http.StatusCreated, status, string(body))

	var resp map[string]any
	require.NoError(t, json.Unmarshal(body, &resp))
	pr := asMap(t, resp["pr"])

	require.Equal(t, "OPEN", getString(t, pr, "status"))
	revs := getArray(t, pr, "assigned_reviewers")
	require.Len(t, revs, 2)
	require.True(t, containsString(revs, r1))
	require.True(t, containsString(revs, r2))
	require.False(t, containsString(revs, author))
	require.False(t, containsString(revs, inactive))
}

// Assigns one reviewer when only one candidate exists.
func TestPR_Create_AssignsOneWhenOnlyOneCandidate(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()

	tn := uniqueID("e2e-pr-one")
	author := "u1-" + tn
	only := "u2-" + tn

	status, body := post(t, base+teamAddPath, map[string]any{
		"team_name": tn,
		"members": []any{
			map[string]any{"user_id": author, "username": "author", "is_active": true},
			map[string]any{"user_id": only, "username": "only", "is_active": true},
		},
	}, nil)
	require.Equal(t, http.StatusCreated, status, string(body))

	status, body = post(t, base+prCreatePath, map[string]any{
		"pull_request_id":   "pr-" + tn,
		"pull_request_name": "add A",
		"author_id":         author,
	}, nil)
	require.Equal(t, http.StatusCreated, status, string(body))

	var resp map[string]any
	require.NoError(t, json.Unmarshal(body, &resp))
	pr := asMap(t, resp["pr"])
	revs := getArray(t, pr, "assigned_reviewers")
	require.Len(t, revs, 1)

	rev0, ok := revs[0].(string)
	require.True(t, ok)
	require.Equal(t, only, rev0)
}

// Assigns zero reviewers when no candidates exist.
func TestPR_Create_AssignsZeroWhenNoCandidates(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()

	tn := uniqueID("e2e-pr-zero")
	author := "u1-" + tn

	status, body := post(t, base+teamAddPath, map[string]any{
		"team_name": tn,
		"members":   []any{map[string]any{"user_id": author, "username": "author", "is_active": true}},
	}, nil)
	require.Equal(t, http.StatusCreated, status, string(body))

	status, body = post(t, base+prCreatePath, map[string]any{
		"pull_request_id":   "pr-" + tn,
		"pull_request_name": "empty",
		"author_id":         author,
	}, nil)
	require.Equal(t, http.StatusCreated, status, string(body))

	var resp map[string]any
	require.NoError(t, json.Unmarshal(body, &resp))
	pr := asMap(t, resp["pr"])
	require.Equal(t, "OPEN", getString(t, pr, "status"))
	require.Len(t, getArray(t, pr, "assigned_reviewers"), 0)
}

// Creating a PR with the same ID twice returns 409.
func TestPR_Create_Duplicate_Returns409(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()

	tn := uniqueID("e2e-pr-dup")
	author := "u1-" + tn
	r1 := "u2-" + tn

	status, body := post(t, base+teamAddPath, map[string]any{
		"team_name": tn,
		"members": []any{
			map[string]any{"user_id": author, "username": "author", "is_active": true},
			map[string]any{"user_id": r1, "username": "r1", "is_active": true},
		},
	}, nil)
	require.Equal(t, http.StatusCreated, status, string(body))

	req := map[string]any{
		"pull_request_id":   "pr-" + tn,
		"pull_request_name": "feature",
		"author_id":         author,
	}

	status, body = post(t, base+prCreatePath, req, nil)
	require.Equal(t, http.StatusCreated, status, string(body))

	status, body = post(t, base+prCreatePath, req, nil)
	require.Equal(t, http.StatusConflict, status, string(body))

	var er map[string]any
	require.NoError(t, json.Unmarshal(body, &er))
	errObj := asMap(t, er["error"])
	require.Equal(t, "PR_EXISTS", getString(t, errObj, "code"))
}

// Reassign replaces exactly one reviewer; author MUST NOT be chosen.
func TestPR_Reassign_Success(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()

	tn := uniqueID("e2e-pr-reassign-ok")
	author := "u1-" + tn
	a := "u2-" + tn
	b := "u3-" + tn
	c := "u4-" + tn
	inactive := "u5-" + tn

	status, body := post(t, base+teamAddPath, map[string]any{
		"team_name": tn,
		"members": []any{
			map[string]any{"user_id": author, "username": "author", "is_active": true},
			map[string]any{"user_id": a, "username": "A", "is_active": true},
			map[string]any{"user_id": b, "username": "B", "is_active": true},
			map[string]any{"user_id": c, "username": "C", "is_active": true},
			map[string]any{"user_id": inactive, "username": "X", "is_active": false},
		},
	}, nil)
	require.Equal(t, http.StatusCreated, status, string(body))

	status, body = post(t, base+prCreatePath, map[string]any{
		"pull_request_id":   "pr-" + tn,
		"pull_request_name": "reassign",
		"author_id":         author,
	}, nil)
	require.Equal(t, http.StatusCreated, status, string(body))

	var created map[string]any
	require.NoError(t, json.Unmarshal(body, &created))
	pr := asMap(t, created["pr"])
	revs := getArray(t, pr, "assigned_reviewers")
	require.Len(t, revs, 2)

	toReplace := getString(t, map[string]any{"x": revs[0]}, "x")
	otherAssigned := getString(t, map[string]any{"x": revs[1]}, "x")

	all := map[string]bool{a: true, b: true, c: true}
	delete(all, toReplace)
	delete(all, otherAssigned)
	var candidate string
	for id := range all {
		candidate = id
	}
	require.NotEmpty(t, candidate)

	status, body = post(t, base+prReassignPath, map[string]any{
		"pull_request_id": "pr-" + tn,
		"old_reviewer_id": toReplace,
	}, nil)
	require.Equal(t, http.StatusOK, status, string(body))

	var rr map[string]any
	require.NoError(t, json.Unmarshal(body, &rr))
	newPR := asMap(t, rr["pr"])
	newRevs := getArray(t, newPR, "assigned_reviewers")

	require.Len(t, newRevs, 2)
	require.False(t, containsString(newRevs, toReplace))
	require.True(t, containsString(newRevs, otherAssigned))
	require.True(t, containsString(newRevs, candidate))
	require.False(t, containsString(newRevs, author))
	require.False(t, containsString(newRevs, inactive))

	replacedBy := getString(t, rr, "replaced_by")
	require.Equal(t, candidate, replacedBy)
	require.True(t, containsString(newRevs, replacedBy))
}

// Reassign returns 409 when no replacement candidate exists.
func TestPR_Reassign_NoCandidate_409(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()

	tn := uniqueID("e2e-pr-reassign-nocand")
	author := "u1-" + tn
	r1 := "u2-" + tn
	r2 := "u3-" + tn

	status, body := post(t, base+teamAddPath, map[string]any{
		"team_name": tn,
		"members": []any{
			map[string]any{"user_id": author, "username": "author", "is_active": true},
			map[string]any{"user_id": r1, "username": "r1", "is_active": true},
			map[string]any{"user_id": r2, "username": "r2", "is_active": true},
		},
	}, nil)
	require.Equal(t, http.StatusCreated, status, string(body))

	status, body = post(t, base+prCreatePath, map[string]any{
		"pull_request_id":   "pr-" + tn,
		"pull_request_name": "no-cand",
		"author_id":         author,
	}, nil)
	require.Equal(t, http.StatusCreated, status, string(body))

	status, body = post(t, base+prReassignPath, map[string]any{
		"pull_request_id": "pr-" + tn,
		"old_reviewer_id": r1,
	}, nil)
	require.Equal(t, http.StatusConflict, status, string(body))

	var er map[string]any
	require.NoError(t, json.Unmarshal(body, &er))
	errObj := asMap(t, er["error"])
	require.Equal(t, "NO_CANDIDATE", getString(t, errObj, "code"))
}

// Reassign returns 404 when the specified old_reviewer_id isn't an assigned reviewer.
func TestPR_Reassign_NotAssigned_404(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()

	tn := uniqueID("e2e-pr-reassign-notass")
	author := "u1-" + tn
	r1 := "u2-" + tn
	r2 := "u3-" + tn
	outsider := "ux-" + tn

	status, body := post(t, base+teamAddPath, map[string]any{
		"team_name": tn,
		"members": []any{
			map[string]any{"user_id": author, "username": "author", "is_active": true},
			map[string]any{"user_id": r1, "username": "r1", "is_active": true},
			map[string]any{"user_id": r2, "username": "r2", "is_active": true},
		},
	}, nil)
	require.Equal(t, http.StatusCreated, status, string(body))

	status, body = post(t, base+prCreatePath, map[string]any{
		"pull_request_id":   "pr-" + tn,
		"pull_request_name": "not-assigned",
		"author_id":         author,
	}, nil)
	require.Equal(t, http.StatusCreated, status, string(body))

	status, body = post(t, base+prReassignPath, map[string]any{
		"pull_request_id": "pr-" + tn,
		"old_reviewer_id": outsider,
	}, nil)
	require.Equal(t, http.StatusNotFound, status, string(body))

	var er map[string]any
	require.NoError(t, json.Unmarshal(body, &er))
	errObj := asMap(t, er["error"])
	require.Equal(t, "NOT_ASSIGNED", getString(t, errObj, "code"))
}

// Merge is idempotent; after MERGED, reassignment is blocked.
func TestPR_Merge_Idempotent_And_BlockReassign(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()

	tn := uniqueID("e2e-pr-merge")
	author := "u1-" + tn
	r1 := "u2-" + tn
	prID := "pr-" + tn

	// team
	status, body := post(t, base+teamAddPath, map[string]any{
		"team_name": tn,
		"members": []any{
			map[string]any{"user_id": author, "username": "author", "is_active": true},
			map[string]any{"user_id": r1, "username": "r1", "is_active": true},
		},
	}, nil)
	require.Equal(t, http.StatusCreated, status, string(body))

	// create PR
	status, body = post(t, base+prCreatePath, map[string]any{
		"pull_request_id":   prID,
		"pull_request_name": "merge-me",
		"author_id":         author,
	}, nil)
	require.Equal(t, http.StatusCreated, status, string(body))

	// first merge
	status, body = post(t, base+prMergePath, map[string]any{
		"pull_request_id": prID,
	}, nil)
	require.Equal(t, http.StatusOK, status, string(body))
	var m1 map[string]any
	require.NoError(t, json.Unmarshal(body, &m1))
	pr := asMap(t, m1["pr"])
	require.Equal(t, "MERGED", getString(t, pr, "status"))

	// second merge (idempotent)
	status, body = post(t, base+prMergePath, map[string]any{
		"pull_request_id": prID,
	}, nil)
	require.Equal(t, http.StatusOK, status, string(body))
	var m2 map[string]any
	require.NoError(t, json.Unmarshal(body, &m2))
	pr2 := asMap(t, m2["pr"])
	require.Equal(t, "MERGED", getString(t, pr2, "status"))

	// reassign after merge -> 409 PR_MERGED
	status, body = post(t, base+prReassignPath, map[string]any{
		"pull_request_id": prID,
		"old_reviewer_id": r1,
	}, nil)
	require.Equal(t, http.StatusConflict, status, string(body))

	var er map[string]any
	require.NoError(t, json.Unmarshal(body, &er))
	errObj := asMap(t, er["error"])
	require.Equal(t, "PR_MERGED", getString(t, errObj, "code"))
}

// Author must never be among assigned_reviewers.
func TestPR_Create_ExcludesAuthor_From_Reviewers(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()

	tn := uniqueID("e2e-pr-exclude-author")
	author := "u1-" + tn

	status, body := post(t, base+teamAddPath, map[string]any{
		"team_name": tn,
		"members":   []any{map[string]any{"user_id": author, "username": "author", "is_active": true}},
	}, nil)
	require.Equal(t, http.StatusCreated, status, string(body))

	status, body = post(t, base+prCreatePath, map[string]any{
		"pull_request_id":   "pr-" + tn,
		"pull_request_name": "exclude-author",
		"author_id":         author,
	}, nil)
	require.Equal(t, http.StatusCreated, status, string(body))

	var resp map[string]any
	require.NoError(t, json.Unmarshal(body, &resp))
	pr := asMap(t, resp["pr"])
	require.False(t, containsString(getArray(t, pr, "assigned_reviewers"), author))
}

// Creating a PR for an unknown author returns 404.
func TestPR_Create_UnknownAuthor_Returns404(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()

	tn := uniqueID("e2e-pr-404")

	status, body := post(t, base+prCreatePath, map[string]any{
		"pull_request_id":   "pr-" + tn,
		"pull_request_name": "unknown-author",
		"author_id":         "no-such-user-" + tn,
	}, nil)
	require.Equal(t, http.StatusNotFound, status, string(body))

	var er map[string]any
	require.NoError(t, json.Unmarshal(body, &er))
	errObj := asMap(t, er["error"])
	require.Equal(t, "NOT_FOUND", getString(t, errObj, "code"))
}
