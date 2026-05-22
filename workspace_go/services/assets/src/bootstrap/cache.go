package bootstrap

import (
	"context"

	"go.uber.org/dig"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	minioModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/minio"
	tieredCacheModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/tieredcache"
)

// InitTieredCache registers MinIO clients (assets + templates) and TieredCache for templates.
func InitTieredCache(c *dig.Container) {
	minioEndpoint, _ := config.GetStringValue("minio_endpoint")
	minioAccessKey, _ := config.GetStringValue("minio_access_key")
	minioSecretKey, _ := config.GetStringValue("minio_secret_key")
	minioUseSSL := config.GetConfigValue("minio_use_ssl").(bool)
	minioRegion, _ := config.GetStringValue("minio_region")

	// Initialize MinIO Client for Assets (L2 Cache - TieredCache)
	// Used for: Asset read model storage for consuming services (Router, JS-Executor)
	// Key format: assets/{assetUUID}.json
	minioAssetsBucket, _ := config.GetStringValue("minio_assets_bucket")

	c.Provide(func() *minioModel.MinIOClient {
		mc, err := minioModel.New(minioModel.Config{
			Endpoint:        minioEndpoint,
			AccessKeyID:     minioAccessKey,
			SecretAccessKey: minioSecretKey,
			UseSSL:          minioUseSSL,
			Region:          minioRegion,
			BucketName:      minioAssetsBucket,
			KeyPrefix:       "assets/",
		})
		if err != nil {
			logger.Panic("[INFRA:MinIO] Failed to connect: " + err.Error())
		}
		logger.Info("[INFRA:Cache] MinIO Assets client initialized (bucket: " + minioAssetsBucket + ")")
		return mc
	}, container.Name("assets"))

	// Initialize MinIO Client for Asset Auth Projection (L2 cache for broker).
	// Slim auth-only payload keyed flat by assetUUID — the broker plugin reads
	// `{assetUUID}.json` on every CONNECT lookup. Separate bucket so the
	// broker's hot path is decoupled from the full read-model bucket.
	minioAssetAuthBucket, _ := config.GetStringValue("minio_asset_auth_bucket")

	c.Provide(func() *minioModel.MinIOClient {
		mc, err := minioModel.New(minioModel.Config{
			Endpoint:        minioEndpoint,
			AccessKeyID:     minioAccessKey,
			SecretAccessKey: minioSecretKey,
			UseSSL:          minioUseSSL,
			Region:          minioRegion,
			BucketName:      minioAssetAuthBucket,
		})
		if err != nil {
			logger.Panic("[INFRA:MinIO] Failed to connect (asset-auth): " + err.Error())
		}
		logger.Info("[INFRA:Cache] MinIO AssetAuth client initialized (bucket: " + minioAssetAuthBucket + ")")
		return mc
	}, container.Name("asset-auth"))

	// Initialize MinIO Client for Asset Templates (L2 Cache - TieredCache)
	// Used for: Template scripts storage for consuming services (JS-Executor)
	// Key format: templates/{templateId}/scripts.json
	minioTemplatesBucket, _ := config.GetStringValue("minio_templates_bucket")

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
			logger.Panic("[INFRA:MinIO] Failed to connect (templates): " + err.Error())
		}
		logger.Info("[INFRA:Cache] MinIO Templates client initialized (bucket: " + minioTemplatesBucket + ")")
		return mc
	}, container.Name("templates"))

	// Initialize TieredCache for Asset Templates (L1 + L2)
	// L0 (RAM): Disabled - scripts are large and rarely change
	// L1 (Disk): Enabled - fast local cache for scripts
	// L2 (MinIO): Source of truth for scripts
	c.Provide(func(params struct {
		container.In
		MinIOClient *minioModel.MinIOClient `name:"templates"`
	}) common.TieredCache {
		tieredCacheL1BaseDir, _ := config.GetStringValue("cache_l1_dir")
		serviceName, _ := config.GetStringValue("service_name")
		tieredCacheL1Dir := tieredCacheL1BaseDir + "/" + serviceName // Include service name to avoid conflicts

		tc, err := tieredCacheModel.New(tieredCacheModel.Config{
			EnableL0:  false, // RAM disabled - scripts are large
			EnableL1:  true,  // Disk enabled - fast local cache
			L1Dir:     tieredCacheL1Dir + "/templates",
			KeyPrefix: "template:",

			// L2 loader to fetch from MinIO
			EnableL2: true,
			L2Loader: func(ctx context.Context, key string) ([]byte, error) {
				result, err := params.MinIOClient.Get(ctx, key+".json")
				if err != nil {
					return nil, err
				}
				return result.Data, nil
			},
		})
		if err != nil {
			logger.Panic("[INFRA:TieredCache] Failed to initialize (templates): " + err.Error())
		}

		logger.Info("[INFRA:Cache] TieredCache (templates) initialized (L0: false, L1: true)")
		return tc
	}, container.Name("templates"))
}
