package services

import (
	"context"
	"testing"

	"mapexVault/src/modules/credentials/application/di"
	"mapexVault/src/modules/credentials/application/dtos"
	"mapexVault/src/modules/credentials/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
)

// capturingRepo records the filters passed to FindWithFilters and returns a
// preset page, so tests can assert the security scoping the service applies.
type capturingRepo struct {
	lastFilters model.Map
	page        *model.PaginatedResult[entities.Credential]
}

func (m *capturingRepo) Create(_ context.Context, _ *entities.Credential) (*entities.Credential, error) {
	return nil, nil
}
func (m *capturingRepo) FindById(_ context.Context, _ *string) (*entities.Credential, error) {
	return nil, nil
}
func (m *capturingRepo) FindByIdAndUpdate(_ context.Context, _ *string, _ map[string]any) (*entities.Credential, error) {
	return nil, nil
}
func (m *capturingRepo) DeleteById(_ context.Context, _ *string) error { return nil }
func (m *capturingRepo) FindWithFilters(_ context.Context, filters model.Map, _ *model.PaginationOpts, _ model.Map) (*model.PaginatedResult[entities.Credential], error) {
	m.lastFilters = filters
	if m.page != nil {
		return m.page, nil
	}
	return &model.PaginatedResult[entities.Credential]{}, nil
}
func (m *capturingRepo) FindActiveWithTokenExpiry(_ context.Context) ([]entities.Credential, error) {
	return nil, nil
}
func (m *capturingRepo) CountDocuments(_ context.Context, _ model.Map) (int64, error) {
	return 0, nil
}

func serviceWithRepo(repo *capturingRepo) *CredentialService {
	return &CredentialService{deps: di.CredentialServiceDependenciesInjection{CredentialRepo: repo}}
}

func orgContext() *reqCtx.RequestContext {
	org := model.NewObjectID().Hex()
	return &reqCtx.RequestContext{OrgContext: &org}
}

// TestGetCredentials_ExcludesTemplatesAndScopesOrg proves the external list
// hardcodes isTemplate:false and is scoped to the caller's organization.
func TestGetCredentials_ExcludesTemplatesAndScopesOrg(t *testing.T) {
	repo := &capturingRepo{}
	svc := serviceWithRepo(repo)

	if _, err := svc.GetCredentials(context.Background(), orgContext(), &dtos.CredentialQueryDTO{}); err != nil {
		t.Fatalf("GetCredentials: %v", err)
	}

	if v, ok := repo.lastFilters["isTemplate"]; !ok || v != false {
		t.Errorf("list filter must hardcode isTemplate=false, got %v (present=%v)", v, ok)
	}
	if _, ok := repo.lastFilters["orgId"]; !ok {
		t.Errorf("list filter must scope by org, filters=%v", repo.lastFilters)
	}
}

// TestGetCredentialById_NotFoundForTemplateOrCrossOrg proves get-by-id is
// scoped (template / wrong org yields an empty page) and returns not-found
// without leaking, while applying isTemplate:false + _id + org to the query.
func TestGetCredentialById_NotFoundForTemplateOrCrossOrg(t *testing.T) {
	repo := &capturingRepo{page: &model.PaginatedResult[entities.Credential]{}} // empty
	svc := serviceWithRepo(repo)

	_, err := svc.GetCredentialById(context.Background(), orgContext(), model.NewObjectID().Hex())
	if err == nil {
		t.Fatal("expected not-found for a template/cross-org id, got nil")
	}

	for _, key := range []string{"isTemplate", "_id", "orgId"} {
		if _, ok := repo.lastFilters[key]; !ok {
			t.Errorf("by-id filter missing %q, filters=%v", key, repo.lastFilters)
		}
	}
	if repo.lastFilters["isTemplate"] != false {
		t.Errorf("by-id filter must force isTemplate=false, got %v", repo.lastFilters["isTemplate"])
	}
}

// TestGetCredentialById_ReturnsOwnedCredential proves a matching, non-template,
// in-org credential is returned.
func TestGetCredentialById_ReturnsOwnedCredential(t *testing.T) {
	id := model.NewObjectID()
	repo := &capturingRepo{page: &model.PaginatedResult[entities.Credential]{
		Items: []entities.Credential{{ID: id, Name: "prod-key"}},
	}}
	svc := serviceWithRepo(repo)

	resp, err := svc.GetCredentialById(context.Background(), orgContext(), id.Hex())
	if err != nil {
		t.Fatalf("GetCredentialById: %v", err)
	}
	if resp == nil || resp.Name != "prod-key" {
		t.Fatalf("expected owned credential 'prod-key', got %+v", resp)
	}
}
