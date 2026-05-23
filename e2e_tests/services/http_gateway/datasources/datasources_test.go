package datasources_test

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
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
	rootClient  *httpclient.HTTPClient
	adminClient *httpclient.HTTPClient
	client      *httpclient.HTTPClient
	ctx         context.Context
)

func TestMain(m *testing.M) {
	if err := utils.SetupE2EEnvironment(); err != nil {
		panic("Failed to setup E2E environment: " + err.Error())
	}

	ctx = context.Background()

	rootClient = httpclient.New(httpclient.Config{BaseURL: constants.GatewayURL})
	rootToken, err := utils.GetRootToken()
	if err != nil {
		panic("Failed to get ROOT token: " + err.Error())
	}
	rootClient.SetHeader("Authorization", "Bearer "+rootToken)

	adminClient = httpclient.New(httpclient.Config{BaseURL: constants.GatewayURL})
	adminToken, err := utils.GetAdminToken()
	if err != nil {
		panic("Failed to get ADMIN token: " + err.Error())
	}
	adminClient.SetHeader("Authorization", "Bearer "+adminToken)
	adminClient.SetHeader("X-Org-Context", constants.MapexosOrgID)

	client = rootClient

	code := m.Run()

	os.Exit(code)
}

func TestCreateDataSource(t *testing.T) {
	payload := loadFixture(t, "create_datasource_http.json")

	resp, err := client.Raw(ctx, "POST", "/api/v1/data_sources", payload)
	require.NoError(t, err)
	utils.AssertCreated(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.NotNil(t, result.Data)

	dataSourceMap := result.Data.(map[string]interface{})
	dataSourceID := dataSourceMap["id"].(string)
	assert.NotEmpty(t, dataSourceID)

	assert.Equal(t, "HTTP DataSource Test", dataSourceMap["name"].(string))
	assert.Equal(t, true, dataSourceMap["enabled"].(bool))
	assert.Equal(t, "pull", dataSourceMap["mode"].(string))
	assert.Equal(t, "http", dataSourceMap["protocol"].(string))

	t.Cleanup(func() {
		cleanupDataSource(t, dataSourceID)
	})
}

func TestGetDataSourceById(t *testing.T) {
	dataSourceID := createTestDataSource(t, "create_datasource_http.json")
	t.Cleanup(func() {
		cleanupDataSource(t, dataSourceID)
	})

	resp, err := client.Raw(ctx, "GET", "/api/v1/data_sources/"+dataSourceID, nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	dataSourceMap := result.Data.(map[string]interface{})
	assert.Equal(t, dataSourceID, dataSourceMap["id"].(string))
	assert.Equal(t, "HTTP DataSource Test", dataSourceMap["name"].(string))
}

func TestUpdateDataSource(t *testing.T) {
	dataSourceID := createTestDataSource(t, "create_datasource_http.json")
	t.Cleanup(func() {
		cleanupDataSource(t, dataSourceID)
	})

	updatePayload := loadFixture(t, "update_datasource.json")

	resp, err := client.Raw(ctx, "PATCH", "/api/v1/data_sources/"+dataSourceID, updatePayload)
	require.NoError(t, err)
	utils.AssertCreated(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	dataSourceMap := result.Data.(map[string]interface{})
	assert.Equal(t, "HTTP DataSource Test Updated", dataSourceMap["name"].(string))
	assert.Equal(t, "Updated description for testing", dataSourceMap["description"].(string))
}

func TestDeleteDataSource(t *testing.T) {
	dataSourceID := createTestDataSource(t, "create_datasource_http.json")

	resp, err := client.Raw(ctx, "DELETE", "/api/v1/data_sources/"+dataSourceID, nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	getResp, err := client.Raw(ctx, "GET", "/api/v1/data_sources/"+dataSourceID, nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, getResp.StatusCode)
}

func TestListDataSources_BasicPagination(t *testing.T) {
	resp, err := client.Raw(ctx, "GET", "/api/v1/data_sources?page=1&perPage=10", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	paginatedResult := result.Data.(map[string]interface{})

	if pagination, ok := paginatedResult["pagination"].(map[string]interface{}); ok {
		assert.NotNil(t, pagination["page"])
		assert.NotNil(t, pagination["perPage"])
		assert.Equal(t, float64(1), pagination["page"].(float64))
		assert.Equal(t, float64(10), pagination["perPage"].(float64))
	}

	if items, ok := paginatedResult["items"].([]interface{}); ok {
		t.Logf("Found %d datasources in the system", len(items))
	}
}

func TestListDataSources_FilterByName(t *testing.T) {
	dataSourceID1 := createTestDataSource(t, "create_datasource_http.json")
	dataSourceID2 := createTestDataSource(t, "create_datasource_mqtt.json")
	t.Cleanup(func() {
		cleanupDataSource(t, dataSourceID1)
		cleanupDataSource(t, dataSourceID2)
	})

	resp, err := client.Raw(ctx, "GET", "/api/v1/data_sources?name=HTTP", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	paginatedResult := result.Data.(map[string]interface{})
	items := paginatedResult["items"].([]interface{})

	// The service performs a case-insensitive partial match — assert the
	// same so saga-created datasources named with lowercase "http" are
	// still considered valid matches.
	for _, item := range items {
		dataSourceMap := item.(map[string]interface{})
		name := dataSourceMap["name"].(string)
		assert.Contains(t, strings.ToLower(name), "http")
	}
}

func TestListDataSources_FilterByEnabled(t *testing.T) {
	dataSourceID1 := createTestDataSource(t, "create_datasource_http.json")
	dataSourceID2 := createTestDataSource(t, "create_datasource_mqtt.json")
	t.Cleanup(func() {
		cleanupDataSource(t, dataSourceID1)
		cleanupDataSource(t, dataSourceID2)
	})

	resp, err := client.Raw(ctx, "GET", "/api/v1/data_sources?enabled=true", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	paginatedResult := result.Data.(map[string]interface{})
	items := paginatedResult["items"].([]interface{})

	for _, item := range items {
		dataSourceMap := item.(map[string]interface{})
		assert.Equal(t, true, dataSourceMap["enabled"].(bool))
	}
}

func TestListDataSources_FilterByMode(t *testing.T) {
	dataSourceID1 := createTestDataSource(t, "create_datasource_http.json")
	dataSourceID2 := createTestDataSource(t, "create_datasource_mqtt.json")
	t.Cleanup(func() {
		cleanupDataSource(t, dataSourceID1)
		cleanupDataSource(t, dataSourceID2)
	})

	resp, err := client.Raw(ctx, "GET", "/api/v1/data_sources?mode=pull", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	paginatedResult := result.Data.(map[string]interface{})
	items := paginatedResult["items"].([]interface{})

	for _, item := range items {
		dataSourceMap := item.(map[string]interface{})
		if mode, ok := dataSourceMap["mode"].(string); ok {
			assert.Equal(t, "pull", mode)
		}
	}
}

func TestListDataSources_FilterByProtocol(t *testing.T) {
	dataSourceID1 := createTestDataSource(t, "create_datasource_http.json")
	dataSourceID2 := createTestDataSource(t, "create_datasource_mqtt.json")
	t.Cleanup(func() {
		cleanupDataSource(t, dataSourceID1)
		cleanupDataSource(t, dataSourceID2)
	})

	resp, err := client.Raw(ctx, "GET", "/api/v1/data_sources?protocol=http", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	paginatedResult := result.Data.(map[string]interface{})
	items := paginatedResult["items"].([]interface{})

	for _, item := range items {
		dataSourceMap := item.(map[string]interface{})
		if protocol, ok := dataSourceMap["protocol"].(string); ok {
			assert.Equal(t, "http", protocol)
		}
	}
}

func TestListDataSources_MultipleFilters(t *testing.T) {
	dataSourceID1 := createTestDataSource(t, "create_datasource_http.json")
	dataSourceID2 := createTestDataSource(t, "create_datasource_mqtt.json")
	t.Cleanup(func() {
		cleanupDataSource(t, dataSourceID1)
		cleanupDataSource(t, dataSourceID2)
	})

	resp, err := client.Raw(ctx, "GET", "/api/v1/data_sources?name=HTTP&enabled=true&mode=pull", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	paginatedResult := result.Data.(map[string]interface{})
	items := paginatedResult["items"].([]interface{})

	for _, item := range items {
		dataSourceMap := item.(map[string]interface{})
		assert.Contains(t, dataSourceMap["name"].(string), "HTTP")
		assert.Equal(t, true, dataSourceMap["enabled"].(bool))
		assert.Equal(t, "pull", dataSourceMap["mode"].(string))
	}
}

func TestListDataSources_Projection(t *testing.T) {
	dataSourceID := createTestDataSource(t, "create_datasource_http.json")
	t.Cleanup(func() {
		cleanupDataSource(t, dataSourceID)
	})

	resp, err := client.Raw(ctx, "GET", "/api/v1/data_sources?projection=name,enabled", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	paginatedResult := result.Data.(map[string]interface{})
	items := paginatedResult["items"].([]interface{})

	if len(items) > 0 {
		dataSourceMap := items[0].(map[string]interface{})
		assert.NotNil(t, dataSourceMap["id"])
		assert.NotNil(t, dataSourceMap["name"])
		assert.NotNil(t, dataSourceMap["enabled"])
	}
}

func TestListDataSources_WithOrgContext(t *testing.T) {
	dataSourceID := createTestDataSource(t, "create_datasource_http.json")
	t.Cleanup(func() {
		cleanupDataSource(t, dataSourceID)
	})

	resp, err := adminClient.Raw(ctx, "GET", "/api/v1/data_sources", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	paginatedResult := result.Data.(map[string]interface{})
	items := paginatedResult["items"].([]interface{})

	for _, item := range items {
		dataSourceMap := item.(map[string]interface{})
		if orgId, ok := dataSourceMap["orgId"].(string); ok {
			assert.Equal(t, constants.MapexosOrgID, orgId)
		}
	}
}

func TestListDataSources_RootUser(t *testing.T) {
	dataSourceID := createTestDataSource(t, "create_datasource_http.json")
	t.Cleanup(func() {
		cleanupDataSource(t, dataSourceID)
	})

	resp, err := rootClient.Raw(ctx, "GET", "/api/v1/data_sources", nil)
	require.NoError(t, err)
	utils.AssertOK(t, resp)

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	paginatedResult := result.Data.(map[string]interface{})
	assert.NotNil(t, paginatedResult["items"])
}

func loadFixture(t *testing.T, filename string) map[string]interface{} {
	data, err := os.ReadFile("fixtures/" + filename)
	require.NoError(t, err, "Failed to load fixture: "+filename)

	var payload map[string]interface{}
	err = json.Unmarshal(data, &payload)
	require.NoError(t, err, "Failed to parse fixture: "+filename)

	return payload
}

func createTestDataSource(t *testing.T, fixtureFile string) string {
	payload := loadFixture(t, fixtureFile)

	resp, err := client.Raw(ctx, "POST", "/api/v1/data_sources", payload)
	require.NoError(t, err, "Failed to create test datasource")
	require.Equal(t, http.StatusCreated, resp.StatusCode, "Expected 201 Created")

	var result types.StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err, "Failed to parse create response")

	dataSourceMap := result.Data.(map[string]interface{})
	dataSourceID := dataSourceMap["id"].(string)
	require.NotEmpty(t, dataSourceID, "DataSource ID should not be empty")

	return dataSourceID
}

func cleanupDataSource(t *testing.T, dataSourceID string) {
	resp, err := client.Raw(ctx, "DELETE", "/api/v1/data_sources/"+dataSourceID, nil)
	if err != nil {
		t.Logf("Failed to cleanup datasource %s: %v", dataSourceID, err)
		return
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		t.Logf("Unexpected status code during cleanup: %d", resp.StatusCode)
	}
}
