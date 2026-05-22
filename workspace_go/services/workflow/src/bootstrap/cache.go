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

// InitTieredCache registers MinIO client and TieredCache instances for workflow definitions, plugins, and instances.
func InitTieredCache(c *dig.Container) {
	minioEndpoint, _ := config.GetStringValue("minio_endpoint")
	minioAccessKey, _ := config.GetStringValue("minio_access_key")
	minioSecretKey, _ := config.GetStringValue("minio_secret_key")
	minioUseSSL := config.GetConfigValue("minio_use_ssl").(bool)
	minioRegion, _ := config.GetStringValue("minio_region")

	// Initialize MinIO Client for Definitions (L2 Cache - TieredCache)
	// Used for: Workflow definition storage for consuming services
	// Key format: definitions/{orgId}/{definitionId}.json
	minioDefinitionsBucket, _ := config.GetStringValue("minio_definitions_bucket")

	c.Provide(func() *minioModel.MinIOClient {
		mc, err := minioModel.New(minioModel.Config{
			Endpoint:        minioEndpoint,
			AccessKeyID:     minioAccessKey,
			SecretAccessKey: minioSecretKey,
			UseSSL:          minioUseSSL,
			Region:          minioRegion,
			BucketName:      minioDefinitionsBucket,
		})
		if err != nil {
			logger.Panic("Failed to connect to MinIO: " + err.Error())
		}
		logger.Info("[APP:BOOTSTRAP] MinIO Definitions client initialized (bucket: " + minioDefinitionsBucket + ")")
		return mc
	}, container.Name("definitions"))

	// Initialize TieredCache for Workflow Definitions (L0 + L1 + L2)
	// L0 (RAM): Enabled - definitions are small JSON, frequently accessed during execution
	// L1 (Disk): Enabled - fast local cache for definitions
	// L2 (MinIO): Source of truth for definitions
	c.Provide(func(params struct {
		container.In
		MinIOClient *minioModel.MinIOClient `name:"definitions"`
	}) common.TieredCache {
		cacheL0MaxSize, _ := config.GetIntValue("cache_l0_max_size")
		cacheL0MaxItems, _ := config.GetIntValue("cache_l0_max_items")
		cacheL0TTL, _ := config.GetIntValue("cache_l0_ttl_seconds")
		
		cacheL1BaseDir, _ := config.GetStringValue("cache_l1_dir")
		serviceName, _ := config.GetStringValue("service_name")
		cacheL1Dir := cacheL1BaseDir + "/" + serviceName
		cacheL1MaxSize, _ := config.GetIntValue("cache_l1_max_size")
		cacheL1TTL, _ := config.GetIntValue("cache_l1_ttl_seconds")

		tc, err := tieredCacheModel.New(tieredCacheModel.Config{

			// Cache L0
			EnableL0:     true,
			L0MaxSize:    int64(cacheL0MaxSize),
			L0MaxItems:   int64(cacheL0MaxItems),
			L0DefaultTTL: time.Duration(cacheL0TTL) * time.Second,

			// Cache L1
			EnableL1:     true,
			L1Dir:        cacheL1Dir + "/definitions",
			L1MaxSize:    int64(cacheL1MaxSize),
			L1DefaultTTL: time.Duration(cacheL1TTL) * time.Second,

			// L2 loader to fetch from MinIO
			// Key format: {orgId}/{definitionId} -> MinIO path: definitions/{orgId}/{definitionId}.json
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
			logger.Panic("Failed to initialize TieredCache (definitions): " + err.Error())
		}

		logger.Info("[APP:BOOTSTRAP] TieredCache (definitions) initialized (L0: true, L1: true, L2: MinIO)")
		return tc
	}, container.Name("definitions"))

	// Initialize TieredCache for Plugin Manifests (L0 + L1 only, NO L2)
	// L0 (RAM): Enabled - plugin manifests are small JSON (~5KB), frequently accessed
	// L1 (Disk): Enabled - fast local cache for manifests
	// L2 (MinIO): Disabled - MongoDB is source of truth for plugins
	c.Provide(func() common.TieredCache {
		pluginsL0MaxSize, _ := config.GetIntValue("plugins_cache_l0_max_size")
		pluginsL0MaxItems, _ := config.GetIntValue("plugins_cache_l0_max_items")
		pluginsL0TTL, _ := config.GetIntValue("plugins_cache_l0_ttl_seconds")

		cacheL1BaseDir, _ := config.GetStringValue("cache_l1_dir")
		serviceName, _ := config.GetStringValue("service_name")
		cacheL1Dir := cacheL1BaseDir + "/" + serviceName
		pluginsL1MaxSize, _ := config.GetIntValue("plugins_cache_l1_max_size")
		pluginsL1TTL, _ := config.GetIntValue("plugins_cache_l1_ttl_seconds")

		tc, err := tieredCacheModel.New(tieredCacheModel.Config{

			// Cache L0
			EnableL0:     true,
			L0MaxSize:    int64(pluginsL0MaxSize),
			L0MaxItems:   int64(pluginsL0MaxItems),
			L0DefaultTTL: time.Duration(pluginsL0TTL) * time.Second,

			// Cache L1
			EnableL1:     true,
			L1Dir:        cacheL1Dir + "/plugins",
			L1MaxSize:    int64(pluginsL1MaxSize),
			L1DefaultTTL: time.Duration(pluginsL1TTL) * time.Second,

			// No L2 — MongoDB is source of truth for plugins
			EnableL2: false,
		})
		if err != nil {
			logger.Panic("Failed to initialize TieredCache (plugins): " + err.Error())
		}

		logger.Info("[APP:BOOTSTRAP] TieredCache (plugins) initialized (L0: true, L1: true, L2: disabled)")
		return tc
	}, container.Name("plugins"))

	// Initialize TieredCache for Workflow Instances (L0 + L1 only, NO L2)
	// L0 (RAM): Enabled - instance configs are small (~1KB), read on every execution
	// L1 (Disk): Enabled - fast local cache for instance configs
	// L2 (MinIO): Disabled - MongoDB is source of truth for instances
	c.Provide(func() common.TieredCache {
		instancesL0MaxSize, _ := config.GetIntValue("instances_cache_l0_max_size")
		instancesL0MaxItems, _ := config.GetIntValue("instances_cache_l0_max_items")
		instancesL0TTL, _ := config.GetIntValue("instances_cache_l0_ttl_seconds")

		cacheL1BaseDir, _ := config.GetStringValue("cache_l1_dir")
		serviceName, _ := config.GetStringValue("service_name")
		cacheL1Dir := cacheL1BaseDir + "/" + serviceName
		instancesL1MaxSize, _ := config.GetIntValue("instances_cache_l1_max_size")
		instancesL1TTL, _ := config.GetIntValue("instances_cache_l1_ttl_seconds")

		tc, err := tieredCacheModel.New(tieredCacheModel.Config{

			// Cache L0
			EnableL0:     true,
			L0MaxSize:    int64(instancesL0MaxSize),
			L0MaxItems:   int64(instancesL0MaxItems),
			L0DefaultTTL: time.Duration(instancesL0TTL) * time.Second,

			// Cache L1
			EnableL1:     true,
			L1Dir:        cacheL1Dir + "/instances",
			L1MaxSize:    int64(instancesL1MaxSize),
			L1DefaultTTL: time.Duration(instancesL1TTL) * time.Second,

			// No L2 — MongoDB is source of truth for instances
			EnableL2: false,
		})
		if err != nil {
			logger.Panic("Failed to initialize TieredCache (instances): " + err.Error())
		}

		logger.Info("[APP:BOOTSTRAP] TieredCache (instances) initialized (L0: true, L1: true, L2: disabled)")
		return tc
	}, container.Name("instances"))
}
