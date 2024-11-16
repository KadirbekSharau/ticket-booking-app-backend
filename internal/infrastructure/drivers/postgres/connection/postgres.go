package postgres

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"ticket-booking-app-backend/internal/helpers"
	"ticket-booking-app-backend/internal/infrastructure/drivers/postgres/models"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	dbUsernameKey = "DB_USERNAME"
	dbPasswordKey = "DB_PASSWORD"
	dbHostKey     = "DB_HOST"
	dbPortKey     = "DB_PORT"
	dbNameKey     = "DB_NAME"
)

var (
	dbInstance *Database
	once       sync.Once
)

type Database struct {
	Conn *gorm.DB
}

type DbCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Dbname   string `json:"dbname"`
	Port     string `json:"port"`
}

func NewDatabase() *Database {
	username, _ := helpers.GetEnv(dbUsernameKey)
	password, _ := helpers.GetEnv(dbPasswordKey)
	host, _ := helpers.GetEnv(dbHostKey)
	port, _ := helpers.GetEnv(dbPortKey)
	dbname, _ := helpers.GetEnv(dbNameKey)

	once.Do(func() {
		dbCredentials := DbCredentials{
			Username: username,
			Password: password,
			Host:     host,
			Port:     port,
			Dbname:   dbname,
		}

		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
			dbCredentials.Host, dbCredentials.Username, dbCredentials.Password, dbCredentials.Dbname, dbCredentials.Port,
		)

		var db *gorm.DB
		var err error

		for i := 0; i < 5; i++ {
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			})
			if err == nil {
				break
			}
			fmt.Printf("failed to connect to database: %v\n", err)
			time.Sleep(10 * time.Second)
		}

		if err != nil {
			logrus.Fatalf("failed to connect database: %v", err)
		}

		if err = execSqlFromFile(db, "internal/infrastructure/drivers/postgres/setup/setup.sql"); err != nil {
			logrus.Fatalf("failed to execute setup sql file: %v", err)
		}

		// Auto-migrate the database schema
		err = db.AutoMigrate(
			&models.User{},
			&models.Event{},
			&models.Ticket{},
			&models.Payment{},
		)
		if err != nil {
			logrus.Fatalf("failed to auto-migrate database: %v", err)
		}

		dbInstance = &Database{Conn: db}
		logrus.Info("Database connection established and migrated")
	})

	return dbInstance
}

func execSqlFromFile(db *gorm.DB, filePath string) error {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("could not determine absolute path: %v", err)
	}

	content, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("could not read sql file: %v", err)
	}

	err = db.Exec(string(content)).Error
	if err != nil {
		return fmt.Errorf("could not execute sql file: %v", err)
	}

	return nil
}
