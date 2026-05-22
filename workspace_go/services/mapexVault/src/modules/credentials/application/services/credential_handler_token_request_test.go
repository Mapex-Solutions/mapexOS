package services

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"mapexVault/src/modules/credentials/application/di"
	"mapexVault/src/modules/credentials/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

func newTestServiceForTokenRequest() *CredentialService {
	return &CredentialService{
		deps: di.CredentialServiceDependenciesInjection{
			Publisher:       &mockPublisher{},
			ScheduleManager: &mockScheduleManager{},
		},
	}
}

func TestExecuteTokenRequest_JSONBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var parsed map[string]interface{}
		json.Unmarshal(body, &parsed)
		if parsed["username"] != "test" {
			t.Fatalf("expected username=test, got %v", parsed["username"])
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": "new_token",
			"expires_in":   float64(3600),
		})
	}))
	defer server.Close()

	svc := newTestServiceForTokenRequest()
	cred := &entities.Credential{ID: model.NewObjectID(), Type: entities.CredentialOAuth2}
	config := &entities.TokenRequestConfig{
		Url:             server.URL,
		ContentType:     "application/json",
		Body:            map[string]interface{}{"username": "test"},
		AccessTokenPath: "access_token",
	}

	resp, err := svc.executeTokenRequest(cred, config, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.AccessToken != "new_token" {
		t.Fatalf("expected access_token=new_token, got %s", resp.AccessToken)
	}
	if resp.ExpiresAt == nil {
		t.Fatal("expected non-nil ExpiresAt from expires_in")
	}
}

func TestExecuteTokenRequest_FormURLEncodedBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Fatalf("expected form-urlencoded content type, got %s", r.Header.Get("Content-Type"))
		}
		body, _ := io.ReadAll(r.Body)
		if len(body) == 0 {
			t.Fatal("expected non-empty body")
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": "form_token",
		})
	}))
	defer server.Close()

	svc := newTestServiceForTokenRequest()
	cred := &entities.Credential{ID: model.NewObjectID(), Type: entities.CredentialOAuth2}
	config := &entities.TokenRequestConfig{
		Url:         server.URL,
		ContentType: "application/x-www-form-urlencoded",
		Body: map[string]interface{}{
			"grant_type":    "refresh_token",
			"refresh_token": "{{credential.refreshToken}}",
		},
		AccessTokenPath: "access_token",
	}
	data := map[string]interface{}{"refreshToken": "my_refresh_token"}

	resp, err := svc.executeTokenRequest(cred, config, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.AccessToken != "form_token" {
		t.Fatalf("expected access_token=form_token, got %s", resp.AccessToken)
	}
}

func TestExecuteTokenRequest_CustomHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test_token" {
			t.Fatalf("expected Authorization=Bearer test_token, got %s", auth)
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": "header_token",
		})
	}))
	defer server.Close()

	svc := newTestServiceForTokenRequest()
	cred := &entities.Credential{ID: model.NewObjectID(), Type: entities.CredentialOAuth2}
	config := &entities.TokenRequestConfig{
		Url:             server.URL,
		Headers:         map[string]string{"Authorization": "Bearer {{credential.accessToken}}"},
		AccessTokenPath: "access_token",
	}
	data := map[string]interface{}{"accessToken": "test_token"}

	resp, err := svc.executeTokenRequest(cred, config, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.AccessToken != "header_token" {
		t.Fatalf("expected access_token=header_token, got %s", resp.AccessToken)
	}
}

func TestExecuteTokenRequest_QueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientId := r.URL.Query().Get("client_id")
		if clientId != "abc" {
			t.Fatalf("expected client_id=abc, got %s", clientId)
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": "query_token",
		})
	}))
	defer server.Close()

	svc := newTestServiceForTokenRequest()
	cred := &entities.Credential{ID: model.NewObjectID(), Type: entities.CredentialOAuth2}
	config := &entities.TokenRequestConfig{
		Method:          "GET",
		Url:             server.URL,
		QueryParams:     map[string]string{"client_id": "{{credential.clientId}}"},
		AccessTokenPath: "access_token",
	}
	data := map[string]interface{}{"clientId": "abc"}

	resp, err := svc.executeTokenRequest(cred, config, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.AccessToken != "query_token" {
		t.Fatalf("expected access_token=query_token, got %s", resp.AccessToken)
	}
}

func TestExecuteTokenRequest_NestedResponsePath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"tokens": map[string]interface{}{
					"access": "nested_token",
				},
			},
		})
	}))
	defer server.Close()

	svc := newTestServiceForTokenRequest()
	cred := &entities.Credential{ID: model.NewObjectID(), Type: entities.CredentialOAuth2}
	config := &entities.TokenRequestConfig{
		Url:             server.URL,
		AccessTokenPath: "data.tokens.access",
	}

	resp, err := svc.executeTokenRequest(cred, config, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.AccessToken != "nested_token" {
		t.Fatalf("expected access_token=nested_token, got %s", resp.AccessToken)
	}
}
