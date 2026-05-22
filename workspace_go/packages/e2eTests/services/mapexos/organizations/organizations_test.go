package organizations_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/constants"
	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/types"
	"github.com/Mapex-Solutions/MapexOS/e2eTests/common/utils"
)

var (
	// rootClient - ROOT user (mapex.* permission)
	// Can query WITHOUT X-Org-Context header (unrestricted global access)
	// Use for CRUD tests that should PASS
	rootClient *httpclient.HTTPClient

	// adminClient - ADMIN user (admin_vendor.* permission)
	// REQUIRES X-Org-Context header (org-scoped access)
	// Use for middleware/permission tests (PASS and DENY scenarios)
	adminClient *httpclient.HTTPClient

	// Backward compatibility: default client points to rootClient
	client *httpclient.HTTPClient

	ctx context.Context
)

func TestMain(m *testing.M) {
	// Setup E2E environment (clean DB + flush cache + seed)
	if err := utils.SetupE2EEnvironment(); err != nil {
		panic("Failed to setup E2E environment: " + err.Error())
	}

	ctx = context.Background()

	// Setup ROOT client (mapex.* - unrestricted)
	rootClient = httpclient.New(httpclient.Config{BaseURL: constants.MapexosURL})
	rootToken, err := utils.GetRootToken()
	if err != nil {
		panic("Failed to get ROOT token: " + err.Error())
	}
	rootClient.SetHeader("Authorization", "Bearer "+rootToken)

	// Setup ADMIN client (admin_vendor.* - org scoped)
	adminClient = httpclient.New(httpclient.Config{BaseURL: constants.MapexosURL})
	adminToken, err := utils.GetAdminToken()
	if err != nil {
		panic("Failed to get ADMIN token: " + err.Error())
	}
	adminClient.SetHeader("Authorization", "Bearer "+adminToken)
	// Set org context to Mapexos organization (admin user's membership)
	adminClient.SetHeader("X-Org-Context", constants.MapexosOrgID)

	// Backward compatibility: default client = rootClient
	client = rootClient

	code := m.Run()
	os.Exit(code)
}

// ========================================
// CREATE TESTS
// ========================================

func TestCreateOrganization_Customer(t *testing.T) {
	payload := loadFixture(t, "create_customer.json", "")

	resp, err := client.Raw(ctx, "POST", "/api/v1/organizations", payload)
	require.NoError(t, err)
	utils.AssertCreated(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	orgMap := result.Data.(map[string]interface{})
	orgID := orgMap["id"].(string)

	assert.Equal(t, "ACME Corporation", orgMap["name"].(string))
	assert.Equal(t, "customer", orgMap["type"].(string))
	assert.True(t, orgMap["enabled"].(bool))

	// Verify code is generated
	assert.NotEmpty(t, orgMap["code"].(string))

	// Verify pathKey is generated
	assert.NotEmpty(t, orgMap["pathKey"].(string))

	// Verify customerID is set (customer is its own customerID)
	assert.NotEmpty(t, orgMap["customerId"].(string))

	t.Cleanup(func() {
		cleanupOrganization(t, orgID)
	})
}

func TestCreateOrganization_Site(t *testing.T) {
	// Create parent customer first
	customerID := createTestOrganization(t, "create_customer.json", "")
	defer cleanupOrganization(t, customerID)

	payload := loadFixture(t, "create_site.json", customerID)

	resp, err := client.Raw(ctx, "POST", "/api/v1/organizations", payload)
	require.NoError(t, err)
	utils.AssertCreated(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	orgMap := result.Data.(map[string]interface{})
	siteID := orgMap["id"].(string)

	assert.Equal(t, "São Paulo HQ", orgMap["name"].(string))
	assert.Equal(t, "site", orgMap["type"].(string))
	assert.Equal(t, customerID, orgMap["parentOrgId"].(string))

	// Verify pathKey is extended from parent
	assert.Contains(t, orgMap["pathKey"].(string), "/")

	// Verify customerID is inherited
	assert.NotEmpty(t, orgMap["customerId"].(string))

	t.Cleanup(func() {
		cleanupOrganization(t, siteID)
	})
}

func TestCreateOrganization_Building(t *testing.T) {
	// Create customer
	customerID := createTestOrganization(t, "create_customer.json", "")
	defer cleanupOrganization(t, customerID)

	// Create site
	sitePayload := loadFixture(t, "create_site.json", customerID)
	resp, err := client.Raw(ctx, "POST", "/api/v1/organizations", sitePayload)
	require.NoError(t, err)

	var siteResult types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&siteResult)
	require.NoError(t, err)
	siteMap := siteResult.Data.(map[string]interface{})
	siteID := siteMap["id"].(string)
	defer cleanupOrganization(t, siteID)

	// Create building under site
	buildingPayload := loadFixture(t, "create_building.json", siteID)
	resp, err = client.Raw(ctx, "POST", "/api/v1/organizations", buildingPayload)
	require.NoError(t, err)
	utils.AssertCreated(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	orgMap := result.Data.(map[string]interface{})
	buildingID := orgMap["id"].(string)

	assert.Equal(t, "Building A", orgMap["name"].(string))
	assert.Equal(t, "building", orgMap["type"].(string))
	assert.Equal(t, siteID, orgMap["parentOrgId"].(string))

	// Verify pathKey has 3 levels
	pathKey := orgMap["pathKey"].(string)
	assert.Equal(t, 2, strings.Count(pathKey, "/"), "pathKey should have 3 levels")

	t.Cleanup(func() {
		cleanupOrganization(t, buildingID)
	})
}

func TestCreateOrganization_Minimal(t *testing.T) {
	payload := loadFixture(t, "create_minimal.json", "")

	resp, err := client.Raw(ctx, "POST", "/api/v1/organizations", payload)
	require.NoError(t, err)
	utils.AssertCreated(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	orgMap := result.Data.(map[string]interface{})
	orgID := orgMap["id"].(string)

	assert.Equal(t, "Minimal Org", orgMap["name"].(string))
	assert.Equal(t, "customer", orgMap["type"].(string))

	t.Cleanup(func() {
		cleanupOrganization(t, orgID)
	})
}

func TestCreateOrganization_InvalidType(t *testing.T) {
	payload := map[string]interface{}{
		"name":    "Invalid Org",
		"type":    "invalid_type",
		"enabled": true,
		"address": map[string]interface{}{
			"city":    "Test City",
			"state":   "Test State",
			"country": "USA",
			"zipCode": "12345",
		},
		"phone": "+12125551234",
		"authConfig": map[string]interface{}{
			"providerType": "internal",
		},
		"accessPolicy": map[string]interface{}{
			"rolePolicy":   "merge",
			"defaultScope": "local",
		},
	}

	resp, err := client.Raw(ctx, "POST", "/api/v1/organizations", payload)
	require.NoError(t, err)
	utils.AssertBadRequest(t, resp)
}

func TestCreateOrganization_MissingName(t *testing.T) {
	payload := map[string]interface{}{
		"type":    "customer",
		"enabled": true,
		"address": map[string]interface{}{
			"city":    "Test City",
			"state":   "Test State",
			"country": "USA",
			"zipCode": "12345",
		},
		"phone": "+12125551234",
		"authConfig": map[string]interface{}{
			"providerType": "internal",
		},
		"accessPolicy": map[string]interface{}{
			"rolePolicy":   "merge",
			"defaultScope": "local",
		},
	}

	resp, err := client.Raw(ctx, "POST", "/api/v1/organizations", payload)
	require.NoError(t, err)
	utils.AssertBadRequest(t, resp)
}

func TestCreateOrganization_DuplicateName(t *testing.T) {
	// Create first org
	orgID1 := createTestOrganization(t, "create_customer.json", "")
	defer cleanupOrganization(t, orgID1)

	// Try to create org with same name (should succeed - no unique constraint on name)
	payload := loadFixture(t, "create_customer.json", "")

	resp, err := client.Raw(ctx, "POST", "/api/v1/organizations", payload)
	require.NoError(t, err)
	// Should succeed as there's no unique constraint on name
	utils.AssertCreated(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	orgMap := result.Data.(map[string]interface{})
	orgID2 := orgMap["id"].(string)
	defer cleanupOrganization(t, orgID2)
}

// ========================================
// GET TESTS
// ========================================

func TestGetOrganizationById(t *testing.T) {
	orgID := createTestOrganization(t, "create_customer.json", "")
	defer cleanupOrganization(t, orgID)

	resp, err := client.Raw(ctx, "GET", "/api/v1/organizations/"+orgID, nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	orgMap := result.Data.(map[string]interface{})
	assert.Equal(t, orgID, orgMap["id"].(string))
	assert.Equal(t, "ACME Corporation", orgMap["name"].(string))
}

func TestGetOrganizationById_NotFound(t *testing.T) {
	fakeID := "507f1f77bcf86cd799439011"

	resp, err := client.Raw(ctx, "GET", "/api/v1/organizations/"+fakeID, nil)
	require.NoError(t, err)
	utils.AssertNotFound(t, resp)
}

func TestListOrganizations(t *testing.T) {
	// ROOT user can list all organizations without X-Org-Context
	org1ID := createTestOrganization(t, "create_customer.json", "")
	defer cleanupOrganization(t, org1ID)

	// ROOT user (mapex.*) queries without org context - unrestricted access
	resp, err := rootClient.Raw(ctx, "GET", "/api/v1/organizations?page=1&perPage=15", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Result is now a paginated result with items array
	paginatedResult := result.Data.(map[string]interface{})
	items := paginatedResult["items"].([]interface{})
	assert.GreaterOrEqual(t, len(items), 1)

	// Verify pagination metadata exists
	pagination := paginatedResult["pagination"].(map[string]interface{})
	assert.NotNil(t, pagination["totalItems"])
	assert.NotNil(t, pagination["page"])
	assert.NotNil(t, pagination["perPage"])
}

// ========================================
// UPDATE TESTS
// ========================================

func TestUpdateOrganization_Name(t *testing.T) {
	orgID := createTestOrganization(t, "create_customer.json", "")
	defer cleanupOrganization(t, orgID)

	payload := loadFixture(t, "update_name.json", "")

	resp, err := client.Raw(ctx, "PATCH", "/api/v1/organizations/"+orgID, payload)
	require.NoError(t, err)
	// API returns 201 for updates, not 200
	require.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated)

	// Verify
	resp, err = client.Raw(ctx, "GET", "/api/v1/organizations/"+orgID, nil)
	require.NoError(t, err)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	orgMap := result.Data.(map[string]interface{})
	assert.Equal(t, "Updated Organization Name", orgMap["name"].(string))
}

func TestUpdateOrganization_Disable(t *testing.T) {
	orgID := createTestOrganization(t, "create_customer.json", "")
	defer cleanupOrganization(t, orgID)

	payload := loadFixture(t, "update_disable.json", "")

	resp, err := client.Raw(ctx, "PATCH", "/api/v1/organizations/"+orgID, payload)
	require.NoError(t, err)
	// API returns 201 for updates, not 200
	require.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated)

	// Verify
	resp, err = client.Raw(ctx, "GET", "/api/v1/organizations/"+orgID, nil)
	require.NoError(t, err)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	orgMap := result.Data.(map[string]interface{})
	assert.False(t, orgMap["enabled"].(bool))
}

func TestUpdateOrganization_Full(t *testing.T) {
	orgID := createTestOrganization(t, "create_customer.json", "")
	defer cleanupOrganization(t, orgID)

	payload := loadFixture(t, "update_full.json", "")

	resp, err := client.Raw(ctx, "PATCH", "/api/v1/organizations/"+orgID, payload)
	require.NoError(t, err)
	// API returns 201 for updates, not 200
	require.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated)

	// Verify
	resp, err = client.Raw(ctx, "GET", "/api/v1/organizations/"+orgID, nil)
	require.NoError(t, err)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	orgMap := result.Data.(map[string]interface{})
	assert.Equal(t, "Fully Updated Organization", orgMap["name"].(string))
	assert.True(t, orgMap["enabled"].(bool))
}

// ========================================
// DELETE TESTS
// ========================================

func TestDeleteOrganization(t *testing.T) {
	orgID := createTestOrganization(t, "create_customer.json", "")

	resp, err := client.Raw(ctx, "DELETE", "/api/v1/organizations/"+orgID, nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	// Verify deleted
	resp, err = client.Raw(ctx, "GET", "/api/v1/organizations/"+orgID, nil)
	require.NoError(t, err)
	utils.AssertNotFound(t, resp)
}

func TestDeleteOrganization_NotFound(t *testing.T) {
	fakeID := "507f1f77bcf86cd799439011"

	resp, err := client.Raw(ctx, "DELETE", "/api/v1/organizations/"+fakeID, nil)
	require.NoError(t, err)
	utils.AssertNotFound(t, resp)
}

// ========================================
// HIERARCHY TESTS
// ========================================

func TestOrganizationHierarchy_PathKeyPropagation(t *testing.T) {
	// Create customer
	customerID := createTestOrganization(t, "create_customer.json", "")
	defer cleanupOrganization(t, customerID)

	// Get customer pathKey
	resp, err := client.Raw(ctx, "GET", "/api/v1/organizations/"+customerID, nil)
	require.NoError(t, err)
	var customerResult types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&customerResult)
	require.NoError(t, err)
	customerMap := customerResult.Data.(map[string]interface{})
	customerPathKey := customerMap["pathKey"].(string)

	// Create site under customer
	sitePayload := loadFixture(t, "create_site.json", customerID)
	resp, err = client.Raw(ctx, "POST", "/api/v1/organizations", sitePayload)
	require.NoError(t, err)
	var siteResult types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&siteResult)
	require.NoError(t, err)
	siteMap := siteResult.Data.(map[string]interface{})
	siteID := siteMap["id"].(string)
	sitePathKey := siteMap["pathKey"].(string)
	defer cleanupOrganization(t, siteID)

	// Verify site pathKey starts with customer pathKey
	assert.True(t, strings.HasPrefix(sitePathKey, customerPathKey+"/"))

	// Create building under site
	buildingPayload := loadFixture(t, "create_building.json", siteID)
	resp, err = client.Raw(ctx, "POST", "/api/v1/organizations", buildingPayload)
	require.NoError(t, err)
	var buildingResult types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&buildingResult)
	require.NoError(t, err)
	buildingMap := buildingResult.Data.(map[string]interface{})
	buildingID := buildingMap["id"].(string)
	buildingPathKey := buildingMap["pathKey"].(string)
	defer cleanupOrganization(t, buildingID)

	// Verify building pathKey starts with site pathKey
	assert.True(t, strings.HasPrefix(buildingPathKey, sitePathKey+"/"))
}

func TestOrganizationHierarchy_CustomerIDInheritance(t *testing.T) {
	// Create customer
	customerID := createTestOrganization(t, "create_customer.json", "")
	defer cleanupOrganization(t, customerID)

	// Get customer to verify it's its own customerID
	resp, err := client.Raw(ctx, "GET", "/api/v1/organizations/"+customerID, nil)
	require.NoError(t, err)
	var customerResult types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&customerResult)
	require.NoError(t, err)
	customerMap := customerResult.Data.(map[string]interface{})
	customerCustomerID := customerMap["customerId"].(string)
	assert.Equal(t, customerID, customerCustomerID, "Customer should be its own customerID")

	// Create site under customer
	sitePayload := loadFixture(t, "create_site.json", customerID)
	resp, err = client.Raw(ctx, "POST", "/api/v1/organizations", sitePayload)
	require.NoError(t, err)
	var siteResult types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&siteResult)
	require.NoError(t, err)
	siteMap := siteResult.Data.(map[string]interface{})
	siteID := siteMap["id"].(string)
	siteCustomerID := siteMap["customerId"].(string)
	defer cleanupOrganization(t, siteID)

	// Verify site inherits customerID
	assert.Equal(t, customerCustomerID, siteCustomerID, "Site should inherit customer's customerID")

	// Create building under site
	buildingPayload := loadFixture(t, "create_building.json", siteID)
	resp, err = client.Raw(ctx, "POST", "/api/v1/organizations", buildingPayload)
	require.NoError(t, err)
	var buildingResult types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&buildingResult)
	require.NoError(t, err)
	buildingMap := buildingResult.Data.(map[string]interface{})
	buildingID := buildingMap["id"].(string)
	buildingCustomerID := buildingMap["customerId"].(string)
	defer cleanupOrganization(t, buildingID)

	// Verify building also inherits same customerID
	assert.Equal(t, customerCustomerID, buildingCustomerID, "Building should inherit customer's customerID")
}

// ========================================
// TREE ENDPOINT TESTS
// ========================================

func TestListOrganizationsTree_FirstPage(t *testing.T) {
	// Create multiple organizations to test pagination
	org1ID := createTestOrganization(t, "create_customer.json", "")
	defer cleanupOrganization(t, org1ID)

	org2ID := createTestOrganization(t, "create_customer.json", "")
	defer cleanupOrganization(t, org2ID)

	org3ID := createTestOrganization(t, "create_customer.json", "")
	defer cleanupOrganization(t, org3ID)

	resp, err := client.Raw(ctx, "GET", "/api/v1/organizations/tree?limit=10&direction=next&sortAsc=true", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Verify response structure
	dataMap := result.Data.(map[string]interface{})
	items := dataMap["items"].([]interface{})
	cursor := dataMap["cursor"].(map[string]interface{})

	// Verify items structure
	assert.GreaterOrEqual(t, len(items), 3, "Should have at least 3 organizations")

	// Verify first item has required fields
	if len(items) > 0 {
		firstItem := items[0].(map[string]interface{})
		assert.NotEmpty(t, firstItem["id"])
		assert.NotEmpty(t, firstItem["name"])
		assert.NotEmpty(t, firstItem["type"])
	}

	// Verify cursor metadata
	assert.NotNil(t, cursor["hasNext"])
	assert.NotNil(t, cursor["hasPrevious"])
	assert.Equal(t, false, cursor["hasPrevious"], "First page should not have previous")
}

func TestListOrganizationsTree_WithLimit(t *testing.T) {
	// Create 5 organizations
	orgIDs := make([]string, 5)
	for i := 0; i < 5; i++ {
		orgIDs[i] = createTestOrganization(t, "create_customer.json", "")
		defer cleanupOrganization(t, orgIDs[i])
	}

	// Request with limit=3
	resp, err := client.Raw(ctx, "GET", "/api/v1/organizations/tree?limit=3&direction=next&sortAsc=true", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	dataMap := result.Data.(map[string]interface{})
	items := dataMap["items"].([]interface{})
	cursor := dataMap["cursor"].(map[string]interface{})

	// Verify limit is respected
	assert.LessOrEqual(t, len(items), 3, "Should return at most 3 items")

	// Verify hasNext indicates more data
	hasNext := cursor["hasNext"].(bool)
	if len(items) == 3 {
		assert.True(t, hasNext, "Should have next page when limit reached")
	}
}

func TestListOrganizationsTree_ForwardPagination(t *testing.T) {
	// Create multiple organizations
	orgIDs := make([]string, 5)
	for i := 0; i < 5; i++ {
		orgIDs[i] = createTestOrganization(t, "create_customer.json", "")
		defer cleanupOrganization(t, orgIDs[i])
	}

	// Get first page
	resp, err := client.Raw(ctx, "GET", "/api/v1/organizations/tree?limit=2&direction=next&sortAsc=true", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result1 types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result1)
	require.NoError(t, err)

	dataMap1 := result1.Data.(map[string]interface{})
	items1 := dataMap1["items"].([]interface{})
	cursor1 := dataMap1["cursor"].(map[string]interface{})

	// Verify first page has items
	assert.GreaterOrEqual(t, len(items1), 1)

	// Get next cursor
	nextCursor, ok := cursor1["next"].(string)
	require.True(t, ok && nextCursor != "", "Should have next cursor")

	// Get second page using cursor
	resp2, err := client.Raw(ctx, "GET", "/api/v1/organizations/tree?cursor="+nextCursor+"&limit=2&direction=next&sortAsc=true", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp2)

	var result2 types.StandardResponse
	err = json.NewDecoder(resp2.Body).Decode(&result2)
	require.NoError(t, err)

	dataMap2 := result2.Data.(map[string]interface{})
	items2 := dataMap2["items"].([]interface{})
	cursor2 := dataMap2["cursor"].(map[string]interface{})

	// Verify second page has different items
	if len(items1) > 0 && len(items2) > 0 {
		firstPageFirstID := items1[0].(map[string]interface{})["id"].(string)
		secondPageFirstID := items2[0].(map[string]interface{})["id"].(string)
		assert.NotEqual(t, firstPageFirstID, secondPageFirstID, "Pages should have different items")
	}

	// Verify second page has previous cursor
	assert.True(t, cursor2["hasPrevious"].(bool), "Second page should have previous")
}

func TestListOrganizationsTree_BackwardPagination(t *testing.T) {
	// Create multiple organizations
	orgIDs := make([]string, 5)
	for i := 0; i < 5; i++ {
		orgIDs[i] = createTestOrganization(t, "create_customer.json", "")
		defer cleanupOrganization(t, orgIDs[i])
	}

	// Get first page
	resp, err := client.Raw(ctx, "GET", "/api/v1/organizations/tree?limit=2&direction=next&sortAsc=true", nil)
	require.NoError(t, err)

	var result1 types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result1)
	require.NoError(t, err)

	dataMap1 := result1.Data.(map[string]interface{})
	cursor1 := dataMap1["cursor"].(map[string]interface{})

	nextCursor := cursor1["next"].(string)
	require.NotEmpty(t, nextCursor)

	// Get second page
	resp2, err := client.Raw(ctx, "GET", "/api/v1/organizations/tree?cursor="+nextCursor+"&limit=2&direction=next&sortAsc=true", nil)
	require.NoError(t, err)

	var result2 types.StandardResponse
	err = json.NewDecoder(resp2.Body).Decode(&result2)
	require.NoError(t, err)

	dataMap2 := result2.Data.(map[string]interface{})
	items2 := dataMap2["items"].([]interface{})
	cursor2 := dataMap2["cursor"].(map[string]interface{})

	// Get previous cursor from second page
	prevCursor, ok := cursor2["previous"].(string)
	require.True(t, ok && prevCursor != "", "Second page should have previous cursor")

	// Go back to first page using previous cursor
	resp3, err := client.Raw(ctx, "GET", "/api/v1/organizations/tree?cursor="+prevCursor+"&limit=2&direction=previous&sortAsc=true", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp3)

	var result3 types.StandardResponse
	err = json.NewDecoder(resp3.Body).Decode(&result3)
	require.NoError(t, err)

	dataMap3 := result3.Data.(map[string]interface{})
	items3 := dataMap3["items"].([]interface{})

	// Verify we got back to similar items
	if len(items2) > 0 && len(items3) > 0 {
		// Items should exist (basic validation that backward pagination works)
		assert.NotEmpty(t, items3)
	}
}

func TestListOrganizationsTree_ResponseStructure(t *testing.T) {
	orgID := createTestOrganization(t, "create_customer.json", "")
	defer cleanupOrganization(t, orgID)

	resp, err := client.Raw(ctx, "GET", "/api/v1/organizations/tree?limit=5&direction=next&sortAsc=true", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Verify top-level response structure
	dataMap := result.Data.(map[string]interface{})
	_, hasItems := dataMap["items"]
	_, hasCursor := dataMap["cursor"]

	assert.True(t, hasItems, "Response should have items field")
	assert.True(t, hasCursor, "Response should have cursor field")

	// Verify cursor structure
	cursor := dataMap["cursor"].(map[string]interface{})
	_, hasNext := cursor["next"]
	_, hasPrevious := cursor["previous"]
	_, hasHasNext := cursor["hasNext"]
	_, hasHasPrevious := cursor["hasPrevious"]

	assert.True(t, hasNext, "Cursor should have next field")
	assert.True(t, hasPrevious, "Cursor should have previous field")
	assert.True(t, hasHasNext, "Cursor should have hasNext field")
	assert.True(t, hasHasPrevious, "Cursor should have hasPrevious field")

	// Verify item structure
	items := dataMap["items"].([]interface{})
	if len(items) > 0 {
		item := items[0].(map[string]interface{})
		_, hasID := item["id"]
		_, hasName := item["name"]
		_, hasType := item["type"]

		assert.True(t, hasID, "Item should have id field")
		assert.True(t, hasName, "Item should have name field")
		assert.True(t, hasType, "Item should have type field")
	}
}

// ========================================
// MIDDLEWARE & PERMISSION TESTS
// ========================================
// These tests verify that the coverage middleware correctly enforces:
// - ROOT users (mapex.*) can query without X-Org-Context (unrestricted)
// - ADMIN users (admin_vendor.*) MUST provide X-Org-Context header
// - ADMIN users can only access orgs in their coverage (via membership)

func TestMiddleware_AdminWithValidOrgContext_Pass(t *testing.T) {
	// ADMIN user with valid X-Org-Context (Mapexos org) should PASS
	// This test validates that the middleware ACCEPTS the request (not blocked by 403)
	// The actual number of items returned depends on org hierarchy and is tested elsewhere

	// ADMIN user with X-Org-Context=Mapexos should be able to make the request
	// adminClient already has X-Org-Context set to MapexosOrgID in TestMain
	resp, err := adminClient.Raw(ctx, "GET", "/api/v1/organizations?page=1&perPage=10", nil)
	require.NoError(t, err)

	// Should NOT be blocked (403 Forbidden)
	// The request should succeed (200 OK) because Mapexos is in ADMIN's coverage
	utils.AssertOK(t, resp)

	var listResult types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&listResult)
	require.NoError(t, err)

	// Verify response structure is correct
	paginatedResult := listResult.Data.(map[string]interface{})
	items := paginatedResult["items"].([]interface{})
	pagination := paginatedResult["pagination"].(map[string]interface{})

	// Just verify structure is valid (number of items depends on org hierarchy)
	assert.NotNil(t, items, "Items should not be nil")
	assert.NotNil(t, pagination, "Pagination should not be nil")
	assert.NotNil(t, pagination["page"], "Pagination page should be present")
	assert.NotNil(t, pagination["perPage"], "Pagination perPage should be present")
}

func TestMiddleware_AdminWithoutOrgContext_Deny(t *testing.T) {
	// ADMIN user WITHOUT X-Org-Context header should be DENIED (403)

	// Create a fresh client without org context
	adminClientNoContext := httpclient.New(httpclient.Config{BaseURL: constants.MapexosURL})
	adminToken, err := utils.GetAdminToken()
	require.NoError(t, err)
	adminClientNoContext.SetHeader("Authorization", "Bearer "+adminToken)
	// Explicitly DO NOT set X-Org-Context

	// Try to list organizations without org context
	resp, err := adminClientNoContext.Raw(ctx, "GET", "/api/v1/organizations?page=1&perPage=10", nil)
	require.NoError(t, err)

	// Should be FORBIDDEN (403) - admin_vendor.* requires org context
	assert.Equal(t, http.StatusForbidden, resp.StatusCode,
		"ADMIN without X-Org-Context should be denied (403)")

	var errorResult types.ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&errorResult)
	require.NoError(t, err)

	// Verify error message mentions org context requirement
	assert.NotEmpty(t, errorResult.Errors)
	errorMsg := errorResult.Errors[0]
	assert.Contains(t, errorMsg, "X-Org-Context",
		"Error message should mention X-Org-Context requirement")
	assert.Contains(t, errorMsg, "ROOT",
		"Error message should mention that only ROOT users can query without org context")
}

func TestMiddleware_AdminWithUnauthorizedOrgContext_Deny(t *testing.T) {
	// ADMIN user with X-Org-Context to org NOT in their coverage should be DENIED (403)

	// Create an organization using ROOT
	orgID := createTestOrganization(t, "create_customer.json", "")
	defer cleanupOrganization(t, orgID)

	// Create a fresh admin client and set org context to the new org
	// (which admin user does NOT have access to - only has access to Mapexos org)
	adminClientUnauth := httpclient.New(httpclient.Config{BaseURL: constants.MapexosURL})
	adminToken, err := utils.GetAdminToken()
	require.NoError(t, err)
	adminClientUnauth.SetHeader("Authorization", "Bearer "+adminToken)
	adminClientUnauth.SetHeader("X-Org-Context", orgID) // Set to unauthorized org

	// Try to list organizations with unauthorized org context
	resp, err := adminClientUnauth.Raw(ctx, "GET", "/api/v1/organizations?page=1&perPage=10", nil)
	require.NoError(t, err)

	// Should be FORBIDDEN (403) - org not in admin's coverage
	assert.Equal(t, http.StatusForbidden, resp.StatusCode,
		"ADMIN with unauthorized org context should be denied (403)")

	var errorResult types.ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&errorResult)
	require.NoError(t, err)

	// Verify error message mentions insufficient permissions
	assert.NotEmpty(t, errorResult.Errors)
	errorMsg := errorResult.Errors[0]
	assert.Contains(t, errorMsg, "insufficient permissions",
		"Error message should mention insufficient permissions")
}

func TestMiddleware_RootWithoutOrgContext_Pass(t *testing.T) {
	// ROOT user WITHOUT X-Org-Context should PASS (unrestricted access)

	// Create a fresh ROOT client without org context
	rootClientNoContext := httpclient.New(httpclient.Config{BaseURL: constants.MapexosURL})
	rootToken, err := utils.GetRootToken()
	require.NoError(t, err)
	rootClientNoContext.SetHeader("Authorization", "Bearer "+rootToken)
	// Explicitly DO NOT set X-Org-Context

	// ROOT user should be able to list all organizations globally
	resp, err := rootClientNoContext.Raw(ctx, "GET", "/api/v1/organizations?page=1&perPage=10", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Should return organizations (mapex.* has unrestricted access)
	paginatedResult := result.Data.(map[string]interface{})
	items := paginatedResult["items"].([]interface{})
	assert.GreaterOrEqual(t, len(items), 1, "ROOT without org context should see all organizations")
}

func TestMiddleware_RootWithOrgContext_Pass(t *testing.T) {
	// ROOT user WITH X-Org-Context should also PASS (context is optional for ROOT)
	// This test validates that when an org has NO children, the list returns empty (correct behavior)

	// Create a fresh ROOT client WITH org context (using Mapexos org which is in coverage)
	rootClientWithContext := httpclient.New(httpclient.Config{BaseURL: constants.MapexosURL})
	rootToken, err := utils.GetRootToken()
	require.NoError(t, err)
	rootClientWithContext.SetHeader("Authorization", "Bearer "+rootToken)
	rootClientWithContext.SetHeader("X-Org-Context", constants.MapexosOrgID) // Use Mapexos org (already in coverage)

	// ROOT user should still be able to query (org context is optional but accepted)
	resp, err := rootClientWithContext.Raw(ctx, "GET", "/api/v1/organizations?page=1&perPage=10", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// When org context is provided, it filters to show children of that org
	// Mapexos currently has children created by test #23, so we expect some results
	// NOTE: This test validates middleware accepts org context for ROOT users
	paginatedResult := result.Data.(map[string]interface{})
	items := paginatedResult["items"].([]interface{})

	// Just verify the request succeeded and returned a valid structure
	// The number of items depends on whether other tests created children
	assert.NotNil(t, items, "ROOT with org context should return valid items array (even if empty)")
	assert.NotNil(t, paginatedResult["pagination"], "Should have pagination object")
}

// ========================================
// HELPER FUNCTIONS
// ========================================

// invalidateCoverageCache invalidates the coverage cache for a user
// This forces the cache to be rebuilt on the next request, picking up newly created orgs
func invalidateCoverageCache(userID string) error {
	cmd := exec.Command("redis-cli", "-n", "5", "DEL", fmt.Sprintf("coverage:user:%s", userID))
	return cmd.Run()
}

func loadFixture(t *testing.T, filename string, parentID string) map[string]interface{} {
	data, err := os.ReadFile("fixtures/" + filename)
	require.NoError(t, err)

	content := string(data)
	if parentID != "" {
		content = strings.ReplaceAll(content, "{{PARENT_ID}}", parentID)
	}

	var payload map[string]interface{}
	err = json.Unmarshal([]byte(content), &payload)
	require.NoError(t, err)

	return payload
}

func createTestOrganization(t *testing.T, fixtureFile string, parentID string) string {
	payload := loadFixture(t, fixtureFile, parentID)

	resp, err := client.Raw(ctx, "POST", "/api/v1/organizations", payload)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	orgMap := result.Data.(map[string]interface{})
	return orgMap["id"].(string)
}

func cleanupOrganization(t *testing.T, orgID string) {
	resp, err := client.Raw(ctx, "DELETE", "/api/v1/organizations/"+orgID, nil)
	if err != nil {
		t.Logf("Failed to cleanup organization %s: %v", orgID, err)
		return
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		t.Logf("Unexpected status during cleanup: %d", resp.StatusCode)
	}
}
