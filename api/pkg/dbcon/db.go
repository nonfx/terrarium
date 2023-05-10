package dbcon

import (
	"fmt"
	"time"

	"github.com/cldcvr/terrarium/api/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connect(dsn string, maxRetries int, retryInterval time.Duration) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			return db, nil
		}

		time.Sleep(retryInterval)
	}

	return nil, fmt.Errorf("failed to connect to database after %d retries: %v", maxRetries, err)

}

func Connect(host, user, password, dbName string, port int, sslMode bool) (*gorm.DB, error) {
	var sslModeStr string
	switch sslMode {
	case true:
		sslModeStr = "enable"
	case false:
		sslModeStr = "disable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", host, user, password, dbName, port, sslModeStr)
	return connect(dsn, 30, time.Second)
}

func ConnectFromEnv() (*gorm.DB, error) {
	return Connect(
		config.DBHost(),
		config.DBUser(),
		config.DBPassword(),
		config.DBName(),
		config.DBPort(),
		config.DBSSLMode(),
	)
}
