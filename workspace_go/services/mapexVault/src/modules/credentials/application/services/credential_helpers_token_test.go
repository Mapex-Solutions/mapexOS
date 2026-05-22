package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"testing"
	"time"
)

func buildJWT(exp int64) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload, _ := json.Marshal(map[string]interface{}{"exp": exp, "sub": "user123"})
	payloadB64 := base64.RawURLEncoding.EncodeToString(payload)
	return fmt.Sprintf("%s.%s.dummysignature", header, payloadB64)
}

func TestResolveTokenExpirationFromResponse_JWTExpTakesPriority(t *testing.T) {
	jwtExp := time.Now().Add(2 * time.Hour).Unix()
	token := buildJWT(jwtExp)

	resp := map[string]interface{}{
		"access_token": token,
		"expires_in":   float64(3600), // 1 hour — should be ignored
	}

	result := resolveTokenExpirationFromResponse(resp, token, "")
	if result == nil {
		t.Fatal("expected non-nil expiration")
	}

	// JWT exp should match, not expires_in
	diff := math.Abs(float64(result.Unix() - jwtExp))
	if diff > 1 {
		t.Fatalf("expected JWT exp %d, got %d (diff=%f)", jwtExp, result.Unix(), diff)
	}
}

func TestResolveTokenExpirationFromResponse_FallsBackToExpiresIn(t *testing.T) {
	resp := map[string]interface{}{
		"access_token": "not-a-jwt-token",
		"expires_in":   float64(3600),
	}

	now := time.Now()
	result := resolveTokenExpirationFromResponse(resp, "not-a-jwt-token", "")
	if result == nil {
		t.Fatal("expected non-nil expiration from expires_in fallback")
	}

	expected := now.Add(3600 * time.Second)
	diff := math.Abs(float64(result.Unix() - expected.Unix()))
	if diff > 2 {
		t.Fatalf("expected ~%d, got %d (diff=%f)", expected.Unix(), result.Unix(), diff)
	}
}

func TestResolveTokenExpirationFromResponse_NilWhenNoExpInfo(t *testing.T) {
	resp := map[string]interface{}{
		"token_type": "bearer",
	}

	result := resolveTokenExpirationFromResponse(resp, "", "")
	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}
