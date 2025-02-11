package repository

import (
	"log"

	"github.com/PosokhovVadim/stawberry/internal/repository/model"

	"github.com/PosokhovVadim/stawberry/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.GetDBConnString()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate schemas
	err = db.AutoMigrate(
		&model.User{},
		&model.Store{},
		&model.Product{},
		&model.Offer{},
		&model.Notification{},
		&model.Session{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	return db
}
