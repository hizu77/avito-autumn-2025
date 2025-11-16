package tests

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUsers_SetActive_Unauthorized(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()
	tn := uniqueID("e2e-users-unauth")
	u1 := "u1-" + tn

	// Seed a team and a single active user
	status, body := post(t, base+teamAddPath, map[string]any{
		"team_name": tn,
		"members": []any{
			map[string]any{"user_id": u1, "username": "User", "is_active": true},
		},
	}, nil)
	require.Equal(t, http.StatusCreated, status, string(body))

	// No Authorization header -> must be 401 with code=UNAUTHORIZED
	status, body = post(t, base+usersSetActive, map[string]any{
		"user_id":   u1,
		"is_active": false,
	}, nil)
	require.Equal(t, http.StatusUnauthorized, status, string(body))

	var er map[string]any
	require.NoError(t, json.Unmarshal(body, &er))
	errObj := asMap(t, er["error"])
	require.Equal(t, "UNAUTHORIZED", getString(t, errObj, "code"))
}

func TestUsers_SetActive_NotFound(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()
	token := loginAsDefaultAdmin(t)

	// Authorized call for a non-existing user -> must be 404 NOT_FOUND
	status, body := post(t, base+usersSetActive, map[string]any{
		"user_id":   "no-such-user-" + uniqueID("x"),
		"is_active": false,
	}, map[string]string{"Authorization": "Bearer " + token})
	require.Equal(t, http.StatusNotFound, status, string(body))

	var er map[string]any
	require.NoError(t, json.Unmarshal(body, &er))
	errObj := asMap(t, er["error"])
	require.Equal(t, "NOT_FOUND", getString(t, errObj, "code"))
}

func TestUsers_SetActive_Toggle_Success(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()
	token := loginAsDefaultAdmin(t)

	tn := uniqueID("e2e-users-toggle")
	u1 := "u1-" + tn

	// Seed an active user
	status, body := post(t, base+teamAddPath, map[string]any{
		"team_name": tn,
		"members": []any{
			map[string]any{"user_id": u1, "username": "U1", "is_active": true},
		},
	}, nil)
	require.Equal(t, http.StatusCreated, status, string(body))

	// Deactivate
	status, body = post(t, base+usersSetActive, map[string]any{
		"user_id":   u1,
		"is_active": false,
	}, map[string]string{"Authorization": "Bearer " + token})
	require.Equal(t, http.StatusOK, status, string(body))

	var r1 map[string]any
	require.NoError(t, json.Unmarshal(body, &r1))
	user1 := asMap(t, r1["user"])
	require.Equal(t, u1, getString(t, user1, "user_id"))
	require.Equal(t, tn, getString(t, user1, "team_name"))
	require.False(t, getBool(t, user1, "is_active"))

	// Activate back
	status, body = post(t, base+usersSetActive, map[string]any{
		"user_id":   u1,
		"is_active": true,
	}, map[string]string{"Authorization": "Bearer " + token})
	require.Equal(t, http.StatusOK, status, string(body))

	var r2 map[string]any
	require.NoError(t, json.Unmarshal(body, &r2))
	user2 := asMap(t, r2["user"])
	require.Equal(t, u1, getString(t, user2, "user_id"))
	require.True(t, getBool(t, user2, "is_active"))
}

// /users/getReview must return 200 with an empty list if the user has no PR assignments.
func TestUsers_GetReview_EmptyList_WhenNoPRs(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()
	tn := uniqueID("e2e-users-empty")
	u1 := "u1-" + tn

	// Seed a team and a single active user
	status, body := post(t, base+teamAddPath, map[string]any{
		"team_name": tn,
		"members": []any{
			map[string]any{"user_id": u1, "username": "U1", "is_active": true},
		},
	}, nil)
	require.Equal(t, http.StatusCreated, status, string(body))

	// No PRs for this user -> empty pull_requests array
	q := url.Values{}
	q.Set("user_id", u1)
	status, body = get(t, base+usersGetReview+"?"+q.Encode())
	require.Equal(t, http.StatusOK, status, string(body))

	var resp map[string]any
	require.NoError(t, json.Unmarshal(body, &resp))
	require.Equal(t, u1, getString(t, resp, "user_id"))
	prs := getArray(t, resp, "pull_requests")
	require.Len(t, prs, 0)
}

// /users/getReview must return 200 with an empty list for a non-existing user.
func TestUsers_GetReview_UnknownUser_Empty200(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()
	uid := "no-such-" + uniqueID("u")

	q := url.Values{}
	q.Set("user_id", uid)

	status, body := get(t, base+usersGetReview+"?"+q.Encode())
	require.Equal(t, http.StatusOK, status, string(body))

	var resp map[string]any
	require.NoError(t, json.Unmarshal(body, &resp))

	require.Equal(t, uid, getString(t, resp, "user_id"))
	require.Len(t, getArray(t, resp, "pull_requests"), 0)
}

// /users/getReview must list PRs where the user is assigned as a reviewer.
func TestUsers_GetReview_ListsAssignedPRs(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()
	token := loginAsDefaultAdmin(t)

	tn := uniqueID("e2e-users-list")
	author := "a1-" + tn
	r1 := "r1-" + tn
	r2 := "r2-" + tn

	// Seed a team: author + 2 active reviewers
	status, body := post(t, base+teamAddPath, map[string]any{
		"team_name": tn,
		"members": []any{
			map[string]any{"user_id": author, "username": "Author", "is_active": true},
			map[string]any{"user_id": r1, "username": "R1", "is_active": true},
			map[string]any{"user_id": r2, "username": "R2", "is_active": true},
		},
	}, nil)
	require.Equal(t, http.StatusCreated, status, string(body))

	// Create a PR (the service will auto-assign up to two reviewers)
	prID := "pr-" + tn
	status, body = post(t, base+prCreatePath, map[string]any{
		"pull_request_id":   prID,
		"pull_request_name": "feat",
		"author_id":         author,
	}, map[string]string{"Authorization": "Bearer " + token})
	require.Equal(t, http.StatusCreated, status, string(body))

	// getReview for r1 must include that PR
	q := url.Values{}
	q.Set("user_id", r1)
	status, body = get(t, base+usersGetReview+"?"+q.Encode())
	require.Equal(t, http.StatusOK, status, string(body))

	var resp map[string]any
	require.NoError(t, json.Unmarshal(body, &resp))
	require.Equal(t, r1, getString(t, resp, "user_id"))
	prs := getArray(t, resp, "pull_requests")

	found := false
	for _, it := range prs {
		pr := asMap(t, it)
		if getString(t, pr, "pull_request_id") == prID {
			require.Equal(t, "feat", getString(t, pr, "pull_request_name"))
			require.Equal(t, author, getString(t, pr, "author_id"))
			require.Equal(t, "OPEN", getString(t, pr, "status"))
			found = true
			break
		}
	}
	require.True(t, found, "expected PR %s in user's review list", prID)
}

func TestUsers_Deactivated_NotAssigned_To_NewPRs(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()
	token := loginAsDefaultAdmin(t)

	tn := uniqueID("e2e-users-deact")
	author := "a1-" + tn
	willOff := "r1-" + tn
	other := "r2-" + tn

	// Seed a team: author + two active reviewers
	status, body := post(t, base+teamAddPath, map[string]any{
		"team_name": tn,
		"members": []any{
			map[string]any{"user_id": author, "username": "Author", "is_active": true},
			map[string]any{"user_id": willOff, "username": "R1", "is_active": true},
			map[string]any{"user_id": other, "username": "R2", "is_active": true},
		},
	}, nil)
	require.Equal(t, http.StatusCreated, status, string(body))

	// Deactivate willOff and then create a PR -> willOff must NOT appear in assigned_reviewers
	status, body = post(t, base+usersSetActive, map[string]any{
		"user_id":   willOff,
		"is_active": false,
	}, map[string]string{"Authorization": "Bearer " + token})
	require.Equal(t, http.StatusOK, status, string(body))

	prID := "pr-" + tn
	status, body = post(t, base+prCreatePath, map[string]any{
		"pull_request_id":   prID,
		"pull_request_name": "after-off",
		"author_id":         author,
	}, map[string]string{"Authorization": "Bearer " + token})
	require.Equal(t, http.StatusCreated, status, string(body))

	var c map[string]any
	require.NoError(t, json.Unmarshal(body, &c))
	pr := asMap(t, c["pr"])
	rev := getArray(t, pr, "assigned_reviewers")
	require.False(t, containsString(rev, willOff), "deactivated user must not be assigned")
}

func TestUsers_GetReview_StillShowsAssignedPRs_AfterDeactivation(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()
	token := loginAsDefaultAdmin(t)

	tn := uniqueID("e2e-users-still-visible")
	author := "a1-" + tn
	reviewer := "r1-" + tn

	// Seed a team: author + one active reviewer
	status, body := post(t, base+teamAddPath, map[string]any{
		"team_name": tn,
		"members": []any{
			map[string]any{"user_id": author, "username": "Author", "is_active": true},
			map[string]any{"user_id": reviewer, "username": "R1", "is_active": true},
		},
	}, nil)
	require.Equal(t, http.StatusCreated, status, string(body))

	// Create PR while reviewer is active (so they are assigned)
	prID := "pr-" + tn
	status, body = post(t, base+prCreatePath, map[string]any{
		"pull_request_id":   prID,
		"pull_request_name": "before-off",
		"author_id":         author,
	}, map[string]string{"Authorization": "Bearer " + token})
	require.Equal(t, http.StatusCreated, status, string(body))

	// Deactivate reviewer
	status, body = post(t, base+usersSetActive, map[string]any{
		"user_id":   reviewer,
		"is_active": false,
	}, map[string]string{"Authorization": "Bearer " + token})
	require.Equal(t, http.StatusOK, status, string(body))

	// /users/getReview must still show the already assigned PR
	q := url.Values{}
	q.Set("user_id", reviewer)
	status, body = get(t, base+usersGetReview+"?"+q.Encode())
	require.Equal(t, http.StatusOK, status, string(body))

	var resp map[string]any
	require.NoError(t, json.Unmarshal(body, &resp))
	prs := getArray(t, resp, "pull_requests")

	found := false
	for _, it := range prs {
		pr := asMap(t, it)
		if getString(t, pr, "pull_request_id") == prID {
			found = true
			break
		}
	}
	require.True(t, found, "expected PR %s in user's review list after deactivation", prID)
}
