package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	teamAddPath    = "/team/add"
	teamGetPath    = "/team/get"
	prCreatePath   = "/pullRequest/create"
	prMergePath    = "/pullRequest/merge"
	prReassignPath = "/pullRequest/reassign"
	loginPath      = "/admins/login" // используется в loginAsDefaultAdmin()
	registerPath   = "/admins/register"
	usersSetActive = "/users/setIsActive"
	usersGetReview = "/users/getReview"
)

func mustGetAppURL() string {
	url := os.Getenv("APP_URL")
	if url == "" {
		panic("APP_URL environment variable not set")
	}

	return url
}

func getEnvDefault(key string, defaultVal string) string {
	env := os.Getenv(key)
	if env == "" {
		return defaultVal
	}

	return env
}

func loginAsDefaultAdmin(t *testing.T) string {
	t.Helper()
	base := mustGetAppURL()
	id := getEnvDefault("ADMIN_ID", "admin")
	pass := getEnvDefault("ADMIN_PASSWORD", "admin")

	status, body := post(t, base+loginPath, map[string]any{
		"id":       id,
		"password": pass,
	}, nil)
	require.Equal(t, http.StatusOK, status, "login should return 200, got: %s", string(body))

	var resp map[string]any
	require.NoError(t, json.Unmarshal(body, &resp))
	token := getString(t, resp, "token")
	require.NotEmpty(t, token)
	return token
}

func post(t *testing.T, url string, payload any, headers map[string]string) (int, []byte) {
	t.Helper()

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	cli := &http.Client{Timeout: 10 * time.Second}
	res, err := cli.Do(req)
	require.NoError(t, err)
	defer func() {
		_ = res.Body.Close()
	}()

	b, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	return res.StatusCode, b
}

func get(t *testing.T, rawURL string) (int, []byte) {
	t.Helper()

	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	require.NoError(t, err)

	cli := &http.Client{Timeout: 10 * time.Second}
	res, err := cli.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	return res.StatusCode, b
}

func uniqueID(prefix string) string {
	return prefix + "-" + strconv.FormatInt(time.Now().UnixNano(), 10)
}

func asMap(t *testing.T, v any) map[string]any {
	t.Helper()
	m, ok := v.(map[string]any)
	require.True(t, ok, "expected object, got %T", v)
	return m
}

func getString(t *testing.T, m map[string]any, key string) string {
	t.Helper()
	raw, ok := m[key]
	require.True(t, ok, "missing key %q", key)
	s, ok := raw.(string)
	require.True(t, ok, "key %q is not string (got %T)", key, raw)
	return s
}

func getBool(t *testing.T, m map[string]any, key string) bool {
	t.Helper()
	raw, ok := m[key]
	require.True(t, ok, "missing key %q", key)
	b, ok := raw.(bool)
	require.True(t, ok, "key %q is not bool (got %T)", key, raw)
	return b
}

func getArray(t *testing.T, m map[string]any, key string) []any {
	t.Helper()
	raw, ok := m[key]
	require.True(t, ok, "missing key %q", key)
	arr, ok := raw.([]any)
	require.True(t, ok, "key %q is not array (got %T)", key, raw)
	return arr
}

func containsString(arr []any, val string) bool {
	for _, it := range arr {
		if s, ok := it.(string); ok && s == val {
			return true
		}
	}
	return false
}
