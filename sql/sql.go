package sql

// TODO: start database

import (
	"os"

	"github.com/anonyindian/giga/config"
	"github.com/anonyindian/logger"
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

var SESSION *gorm.DB

func Load(l *logger.Logger) {
	l = l.Create("DATABASE")
	defer l.ChangeLevel(logger.LevelMain).Println("LOADED")
	conn, err := pq.ParseURL(config.ValueOf.DatabaseURI)
	if err != nil {
		l.ChangeLevel(logger.LevelError).Println("failed to parse postgres URL:", err.Error())
		os.Exit(1)
	}
	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 glogger.Default.LogMode(glogger.Silent),
	})
	if err != nil {
		l.Println("failed to connect to DB:", err.Error())
		os.Exit(1)
	}
	SESSION = db
	dB, _ := db.DB()
	dB.SetMaxOpenConns(100)

	// Create tables if they don't exist
	SESSION.AutoMigrate()
}
