package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// Register without Authorization header -> 401 UNAUTHORIZED
func TestAdmin_Register_Unauthorized(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()
	newID := uniqueID("e2e-unauthed")

	status, body := post(t, base+registerPath, map[string]any{
		"id":       newID,
		"password": "pass",
	}, nil)

	require.Equal(t, http.StatusUnauthorized, status, "register without token must be 401")

	var er map[string]any
	require.NoError(t, json.Unmarshal(body, &er))
	errObj := asMap(t, er["error"])
	require.Equal(t, "UNAUTHORIZED", getString(t, errObj, "code"))
}

// Authorized register, then login with the new admin
func TestAdmin_Register_Then_Login_Success(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()
	token := loginAsDefaultAdmin(t)

	newID := uniqueID("e2e-admin")
	newPass := "pass-123"

	// authorized register
	status, body := post(t, base+registerPath, map[string]any{
		"id":       newID,
		"password": newPass,
	}, map[string]string{"Authorization": "Bearer " + token})
	require.Equal(t, http.StatusCreated, status, "register should return 201, got: %s", string(body))

	var rr map[string]any
	require.NoError(t, json.Unmarshal(body, &rr))
	require.Equal(t, newID, getString(t, rr, "id"))

	// login with the new admin
	status, body = post(t, base+loginPath, map[string]any{
		"id":       newID,
		"password": newPass,
	}, nil)
	require.Equal(t, http.StatusOK, status, "login should return 200, got: %s", string(body))

	var lr map[string]any
	require.NoError(t, json.Unmarshal(body, &lr))
	require.NotEmpty(t, getString(t, lr, "token"))
}

// Duplicate register -> 409 ADMIN_EXISTS
func TestAdmin_Register_Duplicate_Returns409(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()
	token := loginAsDefaultAdmin(t)

	newID := uniqueID("e2e-dup")
	newPass := "pass-dup"

	// first register
	status, body := post(t, base+registerPath, map[string]any{
		"id":       newID,
		"password": newPass,
	}, map[string]string{"Authorization": "Bearer " + token})
	require.Equal(t, http.StatusCreated, status, "first register should be 201, got: %s", string(body))

	// duplicate
	status, body = post(t, base+registerPath, map[string]any{
		"id":       newID,
		"password": newPass,
	}, map[string]string{"Authorization": "Bearer " + token})
	require.Equal(t, http.StatusConflict, status, "duplicate admin must be 409")

	var er map[string]any
	require.NoError(t, json.Unmarshal(body, &er))
	errObj := asMap(t, er["error"])
	require.Equal(t, "ADMIN_EXISTS", getString(t, errObj, "code"))
}

// Wrong password -> 401 INVALID_CREDENTIALS
func TestAdmin_Login_InvalidPassword(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()
	token := loginAsDefaultAdmin(t)
	require.NotEmpty(t, token)

	// create a new admin first
	newID := uniqueID("e2e-badpass")
	newPass := "right-pass"
	status, body := post(t, base+registerPath, map[string]any{
		"id":       newID,
		"password": newPass,
	}, map[string]string{"Authorization": "Bearer " + token})
	require.Equal(t, http.StatusCreated, status, "register should return 201, got: %s", string(body))

	// wrong password
	status, body = post(t, base+loginPath, map[string]any{
		"id":       newID,
		"password": "WRONG",
	}, nil)
	require.Equal(t, http.StatusUnauthorized, status, "invalid creds should be 401")

	var er map[string]any
	require.NoError(t, json.Unmarshal(body, &er))
	errObj := asMap(t, er["error"])
	require.Equal(t, "INVALID_CREDENTIALS", getString(t, errObj, "code"))
}

// Empty fields -> 400 BAD_REQUEST (register)
func TestAdmin_Register_ValidationErrors(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()
	token := loginAsDefaultAdmin(t)

	// empty id
	status, body := post(t, base+registerPath, map[string]any{
		"id":       "",
		"password": "x",
	}, map[string]string{"Authorization": "Bearer " + token})
	require.Equal(t, http.StatusBadRequest, status)

	var er map[string]any
	require.NoError(t, json.Unmarshal(body, &er))
	errObj := asMap(t, er["error"])
	require.Equal(t, "BAD_REQUEST", getString(t, errObj, "code"))

	// empty password
	status, body = post(t, base+registerPath, map[string]any{
		"id":       "someone",
		"password": "",
	}, map[string]string{"Authorization": "Bearer " + token})
	require.Equal(t, http.StatusBadRequest, status)

	require.NoError(t, json.Unmarshal(body, &er))
	errObj = asMap(t, er["error"])
	require.Equal(t, "BAD_REQUEST", getString(t, errObj, "code"))
}

// Empty fields -> 400 BAD_REQUEST (login)
func TestAdmin_Login_ValidationErrors(t *testing.T) {
	t.Parallel()

	base := mustGetAppURL()

	// empty id
	status, body := post(t, base+loginPath, map[string]any{
		"id":       "",
		"password": "x",
	}, nil)
	require.Equal(t, http.StatusBadRequest, status)

	var er map[string]any
	require.NoError(t, json.Unmarshal(body, &er))
	errObj := asMap(t, er["error"])
	require.Equal(t, "BAD_REQUEST", getString(t, errObj, "code"))

	// empty password
	status, body = post(t, base+loginPath, map[string]any{
		"id":       "someone",
		"password": "",
	}, nil)
	require.Equal(t, http.StatusBadRequest, status)

	require.NoError(t, json.Unmarshal(body, &er))
	errObj = asMap(t, er["error"])
	require.Equal(t, "BAD_REQUEST", getString(t, errObj, "code"))
}
