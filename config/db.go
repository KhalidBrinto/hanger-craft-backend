package config

import (
	"backend/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type DatabaseConfiguration struct {
	Driver   string
	Dbname   string
	Username string
	Password string
	Host     string
	Port     string
	LogMode  bool
}

func ConnectDatabase() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	log.Println("Attempting to connect to db")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Error),
		SkipDefaultTransaction: true,
		// TranslateError:         true,
		// DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Println(err.Error())
		panic("Failed to Connect Database !")
	} else {
		log.Println("Database Connected Successfully !")
	}
	log.Println("Attempting to migrate")
	db.AutoMigrate(
		models.Brand{},
		models.BrandImage{},
		models.CartItem{},
		models.Category{},
		models.CategoryImage{},
		models.ContentImage{},
		models.Coupon{},
		models.CouponUsageHistory{},
		models.Inventory{},
		models.Order{},
		models.OrderItem{},
		models.Payment{},
		models.PaymentOption{},
		models.Product{},
		models.ProductImage{},
		models.Review{},
		models.ShippingAddress{},
		models.ShoppingCart{},
		models.ShippingOptions{},
		models.User{},
		models.ProductAttribute{},
		models.WishList{},
	)
	log.Println("Finished migration")
	DB = db
}
