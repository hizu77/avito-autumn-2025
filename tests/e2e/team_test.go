package tests

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

// Creates a team with members and then fetches it back.
func TestTeam_Add_And_Get_Success(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()
	tn := uniqueID("e2e-team")

	// Create team (2 members)
	status, body := post(t, base+teamAddPath, map[string]any{
		"team_name": tn,
		"members": []any{
			map[string]any{"user_id": "u1-" + tn, "username": "Alice", "is_active": true},
			map[string]any{"user_id": "u2-" + tn, "username": "Bob", "is_active": true},
		},
	}, nil)
	require.Equal(t, http.StatusCreated, status, "add team should return 201, got: %s", string(body))

	// Validate envelope + payload fields
	var addResp map[string]any
	require.NoError(t, json.Unmarshal(body, &addResp))
	teamObj := asMap(t, addResp["team"])
	require.Equal(t, tn, getString(t, teamObj, "team_name"))
	members := getArray(t, teamObj, "members")
	require.Len(t, members, 2)

	// Fetch the team
	q := url.Values{}
	q.Set("team_name", tn)
	status, body = get(t, base+teamGetPath+"?"+q.Encode())
	require.Equal(t, http.StatusOK, status, "get team should return 200, got: %s", string(body))

	// /team/get returns the Team (no envelope)
	var got map[string]any
	require.NoError(t, json.Unmarshal(body, &got))
	require.Equal(t, tn, getString(t, got, "team_name"))
	gotMembers := getArray(t, got, "members")
	require.Len(t, gotMembers, 2)

	// Ensure both user IDs exist
	gotIDs := map[string]bool{}
	for _, it := range gotMembers {
		m := asMap(t, it)
		gotIDs[getString(t, m, "user_id")] = true
	}
	require.True(t, gotIDs["u1-"+tn])
	require.True(t, gotIDs["u2-"+tn])
}

// Creates a team and tries to create it again.
func TestTeam_Add_Duplicate(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()
	tn := uniqueID("e2e-team-dup")

	req := map[string]any{
		"team_name": tn,
		"members": []any{
			map[string]any{"user_id": "du1-" + tn, "username": "User1", "is_active": true},
		},
	}

	// First create -> 201
	status, body := post(t, base+teamAddPath, req, nil)
	require.Equal(t, http.StatusCreated, status, "first add must be 201, got: %s", string(body))

	// Duplicate -> 400
	status, body = post(t, base+teamAddPath, req, nil)
	require.Equal(t, http.StatusBadRequest, status, "duplicate should be 400 or 409")

	var er map[string]any
	require.NoError(t, json.Unmarshal(body, &er))
	errObj := asMap(t, er["error"])
	require.Equal(t, "TEAM_EXISTS", getString(t, errObj, "code"))
}

// Fetch a non-existing team -> 404 NOT_FOUND.
func TestTeam_Get_NotFound(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()
	tn := "no-such-" + uniqueID("team")

	q := url.Values{}
	q.Set("team_name", tn)
	status, body := get(t, base+teamGetPath+"?"+q.Encode())
	require.Equal(t, http.StatusNotFound, status, string(body))

	var er map[string]any
	require.NoError(t, json.Unmarshal(body, &er))
	errObj := asMap(t, er["error"])
	require.Equal(t, "NOT_FOUND", getString(t, errObj, "code"))
}

// Input validation: empty team_name or empty members -> 400 BAD_REQUEST.
func TestTeam_Add_ValidationErrors(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()

	// Empty name
	status, body := post(t, base+teamAddPath, map[string]any{
		"team_name": "",
		"members": []any{
			map[string]any{"user_id": "x1", "username": "X", "is_active": true},
		},
	}, nil)
	require.Equal(t, http.StatusBadRequest, status)

	var er map[string]any
	require.NoError(t, json.Unmarshal(body, &er))
	errObj := asMap(t, er["error"])
	require.Equal(t, "BAD_REQUEST", getString(t, errObj, "code"))

	// Empty members
	status, body = post(t, base+teamAddPath, map[string]any{
		"team_name": uniqueID("e2e-nomembers"),
		"members":   []any{},
	}, nil)
	require.Equal(t, http.StatusBadRequest, status)
	require.NoError(t, json.Unmarshal(body, &er))
	errObj = asMap(t, er["error"])
	require.Equal(t, "BAD_REQUEST", getString(t, errObj, "code"))
}

// Deduplication: if multiple members share the same user_id in the request,
// only one should end up saved/returned
func TestTeam_Add_Deduplicates_SameUserID(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()
	tn := uniqueID("e2e-team-dedup")
	dup := "u-dup-" + tn

	status, body := post(t, base+teamAddPath, map[string]any{
		"team_name": tn,
		"members": []any{
			map[string]any{"user_id": dup, "username": "Alice-1", "is_active": true},
			map[string]any{"user_id": dup, "username": "Alice-2", "is_active": false},
			map[string]any{"user_id": "u2-" + tn, "username": "Bob", "is_active": true},
		},
	}, nil)
	require.Equal(t, http.StatusCreated, status, "add team should return 201, got: %s", string(body))

	var addResp map[string]any
	require.NoError(t, json.Unmarshal(body, &addResp))
	teamObj := asMap(t, addResp["team"])
	members := getArray(t, teamObj, "members")

	// Count unique IDs in the "add" response
	idCnt := map[string]int{}
	for _, it := range members {
		m := asMap(t, it)
		idCnt[getString(t, m, "user_id")]++
	}
	require.Equal(t, 2, len(idCnt), "expected exactly 2 unique users after dedup")
	require.Equal(t, 1, idCnt[dup], "duplicate user_id must appear only once")

	// Verify via /team/get as well
	q := url.Values{}
	q.Set("team_name", tn)
	status, body = get(t, base+teamGetPath+"?"+q.Encode())
	require.Equal(t, http.StatusOK, status, "get team should return 200, got: %s", string(body))

	var got map[string]any
	require.NoError(t, json.Unmarshal(body, &got))
	gotMembers := getArray(t, got, "members")

	idCnt = map[string]int{}
	for _, it := range gotMembers {
		m := asMap(t, it)
		idCnt[getString(t, m, "user_id")]++
	}
	require.Equal(t, 2, len(idCnt))
	require.Equal(t, 1, idCnt[dup])
}
