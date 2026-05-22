package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// resolveTokenExpirationFromResponse determines when a token expires.
// Priority: 1) JWT exp claim from access token, 2) ExpiresInPath from config, 3) expires_in at root, 4) nil
func resolveTokenExpirationFromResponse(response map[string]interface{}, accessToken string, expiresInPath string) *time.Time {
	// 1. Try JWT exp claim from extracted access token
	if accessToken != "" {
		if exp := decodeJWTExpiration(accessToken); exp != nil {
			return exp
		}
	}

	// 2. Try custom expiresInPath from config
	if expiresInPath != "" {
		if val := extractByPath(response, expiresInPath); val != "" {
			if seconds := parseExpiresIn(val); seconds > 0 {
				t := time.Now().Add(time.Duration(seconds) * time.Second)
				return &t
			}
		}
	}

	// 3. Fallback: expires_in at response root
	if expiresIn, ok := response["expires_in"]; ok {
		var seconds float64
		switch v := expiresIn.(type) {
		case float64:
			seconds = v
		case int:
			seconds = float64(v)
		case json.Number:
			f, _ := v.Float64()
			seconds = f
		}
		if seconds > 0 {
			t := time.Now().Add(time.Duration(seconds) * time.Second)
			return &t
		}
	}

	return nil
}

// parseExpiresIn parses a string value as a number of seconds.
func parseExpiresIn(val string) float64 {
	var f float64
	if _, err := fmt.Sscanf(val, "%f", &f); err != nil {
		return 0
	}
	return f
}

// decodeJWTExpiration decodes a JWT token (without verification) and reads the exp claim.
// Returns nil if the token is not a valid JWT or has no exp claim.
func decodeJWTExpiration(token string) *time.Time {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil // Not a JWT
	}

	// Decode payload (part 1)
	payload := parts[1]
	// Add padding if needed
	switch len(payload) % 4 {
	case 2:
		payload += "=="
	case 3:
		payload += "="
	}

	decoded, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return nil
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return nil
	}

	// Read exp claim
	if exp, ok := claims["exp"]; ok {
		var seconds float64
		switch v := exp.(type) {
		case float64:
			seconds = v
		case json.Number:
			f, _ := v.Float64()
			seconds = f
		}
		if seconds > 0 {
			t := time.Unix(int64(seconds), 0)
			return &t
		}
	}

	return nil
}

// extractByPath extracts a value from a nested map using dot-notation path.
// Example: extractByPath(map, "data.token") → map["data"]["token"]
func extractByPath(data map[string]interface{}, path string) string {
	if path == "" {
		return ""
	}

	parts := strings.Split(path, ".")
	current := interface{}(data)

	for _, part := range parts {
		m, ok := current.(map[string]interface{})
		if !ok {
			return ""
		}
		current = m[part]
	}

	if current == nil {
		return ""
	}

	switch v := current.(type) {
	case string:
		return v
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}
