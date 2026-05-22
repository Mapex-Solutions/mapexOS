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

// InitTieredCache registers MinIO client and TieredCache for templates.
func InitTieredCache(c *dig.Container) {
	serviceName, _ := config.GetStringValue("service_name")

	// Initialize S3 Client for Templates (L2 Cache - TieredCache)
	minioEndpoint, _ := config.GetStringValue("minio_endpoint")
	minioAccessKey, _ := config.GetStringValue("minio_access_key")
	minioSecretKey, _ := config.GetStringValue("minio_secret_key")
	minioUseSSL := false
	if val := config.GetConfigValue("minio_use_ssl"); val != nil {
		if b, ok := val.(bool); ok {
			minioUseSSL = b
		}
	}
	minioRegion, _ := config.GetStringValue("minio_region")
	minioTemplatesBucket, _ := config.GetStringValue("minio_templates_bucket")

	c.Provide(func() *minioModel.MinIOClient {
		mc, err := minioModel.New(minioModel.Config{
			Endpoint:        minioEndpoint,
			AccessKeyID:     minioAccessKey,
			SecretAccessKey: minioSecretKey,
			UseSSL:          minioUseSSL,
			Region:          minioRegion,
			BucketName:      minioTemplatesBucket,
			KeyPrefix:       "templates/", // Must match assets service write prefix
		})
		if err != nil {
			logger.Panic("Failed to connect to S3/MinIO: " + err.Error())
		}
		logger.Info("[APP:BOOTSTRAP] S3 Templates client initialized (bucket: " + minioTemplatesBucket + ")")
		return mc
	}, container.Name("templates"))

	// Initialize TieredCache for Templates (L0 + L1 + L2 + Fallback)
	c.Provide(func(params struct {
		container.In
		MinIOClient *minioModel.MinIOClient `name:"templates"`
	}) common.TieredCache {
		cacheL0MaxSize, _ := config.GetIntValue("cache_l0_max_size")
		cacheL0MaxItems, _ := config.GetIntValue("cache_l0_max_items")
		cacheL0TTL, _ := config.GetIntValue("cache_l0_ttl_seconds")
		cacheL1BaseDir, _ := config.GetStringValue("cache_l1_dir")
		cacheL1Dir := cacheL1BaseDir + "/" + serviceName + "/templates"
		cacheL1MaxSize, _ := config.GetIntValue("cache_l1_max_size")
		cacheL1TTL, _ := config.GetIntValue("cache_l1_ttl_seconds")

		// Fallback configuration - calls Assets Service when L2 misses
		assetsURL, _ := config.GetStringValue("assets_url")
		internalAPIKey, _ := config.GetStringValue("internal_api_key")

		tc, err := tieredCacheModel.New(tieredCacheModel.Config{
			EnableL0:     true,
			L0MaxSize:    int64(cacheL0MaxSize),
			L0MaxItems:   int64(cacheL0MaxItems),
			L0DefaultTTL: time.Duration(cacheL0TTL) * time.Second,

			EnableL1:     true,
			L1Dir:        cacheL1Dir,
			L1MaxSize:    int64(cacheL1MaxSize),
			L1DefaultTTL: time.Duration(cacheL1TTL) * time.Second,
			KeyPrefix:    "template:",

			// L2 loader to fetch from S3
			// Key format: {templateId} → S3 path: {templateId}.json
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
			FallbackEndpoint: "/internal/templates",
		})
		if err != nil {
			logger.Panic("Failed to initialize TieredCache (templates): " + err.Error())
		}

		logger.Info("[APP:BOOTSTRAP] TieredCache (templates) initialized (L0+L1+L2(S3)+Fallback: " + assetsURL + ")")
		return tc
	}, container.Name("templates"))
}
