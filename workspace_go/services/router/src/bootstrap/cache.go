package bootstrap

import (
	"context"
	"time"

	"go.uber.org/dig"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	minioModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/minio"
	tieredCacheModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/tieredcache"
)

// InitTieredCache registers MinIO clients and TieredCaches for assets and templates.
// Both caches use the same L0/L1/L2/Fallback architecture with different MinIO buckets.
func InitTieredCache(c *dig.Container) {
	// Initialize MinIO Client for Assets (L2 Cache - TieredCache)
	minioEndpoint, _ := config.GetStringValue("minio_endpoint")
	minioAccessKey, _ := config.GetStringValue("minio_access_key")
	minioSecretKey, _ := config.GetStringValue("minio_secret_key")
	minioUseSSL := config.GetConfigValue("minio_use_ssl").(bool)
	minioRegion, _ := config.GetStringValue("minio_region")
	minioAssetsBucket, _ := config.GetStringValue("minio_assets_bucket")

	c.Provide(func() *minioModel.MinIOClient {
		mc, err := minioModel.New(minioModel.Config{
			Endpoint:        minioEndpoint,
			AccessKeyID:     minioAccessKey,
			SecretAccessKey: minioSecretKey,
			UseSSL:          minioUseSSL,
			Region:          minioRegion,
			BucketName:      minioAssetsBucket,
			KeyPrefix:       "assets/", // Must match assets service write prefix
		})
		if err != nil {
			logger.Panic("Failed to connect to MinIO: " + err.Error())
		}
		logger.Info("[APP:BOOTSTRAP] MinIO Assets client initialized (bucket: " + minioAssetsBucket + ")")
		return mc
	}, container.Name("assets"))

	// Initialize TieredCache for Assets (L0 + L1 + L2 + Fallback)
	c.Provide(func(params struct {
		container.In
		MinIOClient *minioModel.MinIOClient `name:"assets"`
	}) common.TieredCache {
		
		cacheL0MaxSize, _ := config.GetIntValue("cache_l0_max_size")
		cacheL0MaxItems, _ := config.GetIntValue("cache_l0_max_items")
		cacheL0TTL, _ := config.GetIntValue("cache_l0_ttl_seconds")
		
		cacheL1BaseDir, _ := config.GetStringValue("cache_l1_dir")
		serviceName, _ := config.GetStringValue("service_name")
		cacheL1Dir := cacheL1BaseDir + "/" + serviceName // Include service name to avoid conflicts
		cacheL1MaxSize, _ := config.GetIntValue("cache_l1_max_size")
		cacheL1TTL, _ := config.GetIntValue("cache_l1_ttl_seconds")

		// Fallback configuration - calls Assets Service when L2 misses
		assetsURL, _ := config.GetStringValue("assets_url")
		internalAPIKey, _ := config.GetStringValue("internal_api_key")
		fallbackTimeout, _ := config.GetIntValue("cache_fallback_timeout")

		tc, err := tieredCacheModel.New(tieredCacheModel.Config{
			EnableL0:     true,
			L0MaxSize:    int64(cacheL0MaxSize),
			L0MaxItems:   int64(cacheL0MaxItems),
			L0DefaultTTL: time.Duration(cacheL0TTL) * time.Second,
			EnableL1:     true,
			L1Dir:        cacheL1Dir + "/assets",
			L1MaxSize:    int64(cacheL1MaxSize),
			L1DefaultTTL: time.Duration(cacheL1TTL) * time.Second,

			// L2 loader to fetch from MinIO
			// Key format: {orgId}/{assetUUID} → MinIO path: {orgId}/{assetUUID}.json
			EnableL2: true,
			L2Loader: func(ctx context.Context, key string) ([]byte, error) {
				result, err := params.MinIOClient.Get(ctx, key+".json")
				if err != nil {
					return nil, err
				}
				return result.Data, nil
			},

			// Fallback configuration
			FallbackBaseURL:  assetsURL,
			FallbackAPIKey:   internalAPIKey,
			FallbackEndpoint: "/internal/assets",
			FallbackTimeout:  time.Duration(fallbackTimeout) * time.Second,
		})
		if err != nil {
			logger.Panic("Failed to initialize TieredCache (assets): " + err.Error())
		}

		logger.Info("[APP:BOOTSTRAP] TieredCache (assets) initialized (L0: true, L1: true, Fallback: " + assetsURL + ")")
		return tc
	}, container.Name("assets"))

	/**
	 * Template TieredCache (L0 + L1 + L2 MinIO + Fallback HTTP)
	 * Used by EventService to enrich events with template name and description.
	 * Cache key format: {templateOrgId}/{templateId}
	 * MinIO path: templates/{templateOrgId}/{templateId}.json
	 */

	minioTemplatesBucket, _ := config.GetStringValue("minio_templates_bucket")

	// MinIO Client for Templates (L2 Cache source)
	c.Provide(func() *minioModel.MinIOClient {
		mc, err := minioModel.New(minioModel.Config{
			Endpoint:        minioEndpoint,
			AccessKeyID:     minioAccessKey,
			SecretAccessKey: minioSecretKey,
			UseSSL:          minioUseSSL,
			Region:          minioRegion,
			BucketName:      minioTemplatesBucket,
			KeyPrefix:       "templates/",
		})
		if err != nil {
			logger.Panic("Failed to connect to MinIO (templates): " + err.Error())
		}
		logger.Info("[APP:BOOTSTRAP] MinIO Templates client initialized (bucket: " + minioTemplatesBucket + ")")
		return mc
	}, container.Name("templates"))

	// TieredCache for Templates (reuses same L0/L1 config as assets)
	c.Provide(func(params struct {
		container.In
		MinIOClient *minioModel.MinIOClient `name:"templates"`
	}) common.TieredCache {
		cacheL0MaxSize, _ := config.GetIntValue("cache_l0_max_size")
		cacheL0MaxItems, _ := config.GetIntValue("cache_l0_max_items")
		cacheL0TTL, _ := config.GetIntValue("cache_l0_ttl_seconds")

		cacheL1BaseDir, _ := config.GetStringValue("cache_l1_dir")
		serviceName, _ := config.GetStringValue("service_name")
		cacheL1Dir := cacheL1BaseDir + "/" + serviceName
		cacheL1MaxSize, _ := config.GetIntValue("cache_l1_max_size")
		cacheL1TTL, _ := config.GetIntValue("cache_l1_ttl_seconds")

		assetsURL, _ := config.GetStringValue("assets_url")
		internalAPIKey, _ := config.GetStringValue("internal_api_key")
		fallbackTimeout, _ := config.GetIntValue("cache_fallback_timeout")

		tc, err := tieredCacheModel.New(tieredCacheModel.Config{
			EnableL0:     true,
			L0MaxSize:    int64(cacheL0MaxSize),
			L0MaxItems:   int64(cacheL0MaxItems),
			L0DefaultTTL: time.Duration(cacheL0TTL) * time.Second,
			EnableL1:     true,
			L1Dir:        cacheL1Dir + "/templates",
			L1MaxSize:    int64(cacheL1MaxSize),
			L1DefaultTTL: time.Duration(cacheL1TTL) * time.Second,

			// L2 loader — fetches from MinIO
			// Key format: {templateOrgId}/{templateId} → MinIO path: {templateOrgId}/{templateId}.json
			EnableL2: true,
			L2Loader: func(ctx context.Context, key string) ([]byte, error) {
				result, err := params.MinIOClient.Get(ctx, key+".json")
				if err != nil {
					return nil, err
				}
				return result.Data, nil
			},

			// Fallback — calls Assets Service /internal/templates endpoint
			FallbackBaseURL:  assetsURL,
			FallbackAPIKey:   internalAPIKey,
			FallbackEndpoint: "/internal/templates",
			FallbackTimeout:  time.Duration(fallbackTimeout) * time.Second,
		})
		if err != nil {
			logger.Panic("Failed to initialize TieredCache (templates): " + err.Error())
		}

		logger.Info("[APP:BOOTSTRAP] TieredCache (templates) initialized (L0: true, L1: true, Fallback: " + assetsURL + ")")
		return tc
	}, container.Name("templates"))
}
