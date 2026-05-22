package bootstrap

import (
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"go.uber.org/dig"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	chManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/clickhouse/manager"
)

// InitClickHouse registers ClickHouse manager and raw connection in DIG container.
func InitClickHouse(c *dig.Container) {
	clickhouseCfg := config.GetClickHouseConfig()
	monitorInterval, _ := config.GetIntValue("clickhouse_monitor_interval")
	if monitorInterval == 0 {
		monitorInterval = 30 // Default 30 seconds
	}

	c.Provide(func() *chManager.ClickHouseManager {
		ch, err := chManager.New(chManager.Config{
			Host:            clickhouseCfg.Host,
			Port:            clickhouseCfg.Port,
			Database:        clickhouseCfg.Database,
			Username:        clickhouseCfg.Username,
			Password:        clickhouseCfg.Password,
			MaxOpenConns:    20,
			MaxIdleConns:    10,
			EnableMonitor:   true,
			MonitorInterval: time.Duration(monitorInterval) * time.Second,
		})
		if err != nil {
			logger.Panic(err.Error())
		}
		return ch
	})

	// Provide the raw ClickHouse connection for repositories
	c.Provide(func(ch *chManager.ClickHouseManager) driver.Conn {
		return ch.GetConn()
	})
}
