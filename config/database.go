package config

import (
	"fmt"
	"log"
	"tusk/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 3306
	user     = "root"
	password = ""
	dbName   = "tusk"
)

func DatabaseConnection() *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		user, password, host, port, dbName,
	)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return database
}

// ✅ Tambahkan function migration
func RunMigrations(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Task{},
	)

	if err != nil {
		log.Fatal("❌ Migration failed:", err)
	}

	log.Println("✅ Database migrated successfully!")
}

func CreateOwnerAccount(db *gorm.DB) {
	hashedPasswordBytes, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	owner := models.User{
		Role:     "Admin",
		Name:     "Owner",
		Password: string(hashedPasswordBytes),
		Email:    "owner@go.id",
	}

	if db.Where("email=?", owner.Email).First(&owner).RowsAffected == 0 {
		db.Create(&owner)
		log.Println("✅ Owner account created!")
	} else {
		log.Println("ℹ️ Owner account already exists")
	}
}
