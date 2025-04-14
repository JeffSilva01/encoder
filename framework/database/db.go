package database

import (
	"encoder/domain"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	Db            *gorm.DB
	Dsn           string
	DsnTest       string
	DbType        string
	DbTypeTest    string
	Debug         bool
	AutoMigrateDb bool
	Env           string
}

func NewDb() *Database {
	return &Database{}
}

func NewDbTest() *gorm.DB {
	dbInstance := NewDb()
	dbInstance.Env = "test"
	dbInstance.DbTypeTest = "sqlite3"
	dbInstance.DsnTest = ":memory:"
	dbInstance.AutoMigrateDb = true
	dbInstance.Debug = true

	connection, err := dbInstance.Connect()
	if err != nil {
		log.Fatalf("Test db error: %v", err)
	}

	return connection
}

func (d *Database) Connect() (*gorm.DB, error) {
	var err error

	logConfig := logger.Config{
		LogLevel: logger.Silent,
		Colorful: false,
	}

	if d.Debug {
		logConfig.LogLevel = logger.Info
		logConfig.Colorful = true
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logConfig,
	)

	openConfig := &gorm.Config{
		Logger: newLogger,
	}

	if d.Env != "test" {
		d.Db, err = gorm.Open(postgres.Open(d.Dsn), openConfig)
	} else {
		d.Db, err = gorm.Open(sqlite.Open(d.DsnTest), openConfig)
	}

	if err != nil {
		return nil, err
	}

	if d.AutoMigrateDb {
		d.Db.AutoMigrate(&domain.Video{}, &domain.Job{})
	}

	return d.Db, nil
}
