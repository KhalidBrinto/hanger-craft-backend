package config

import (
	"log"

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
	// dbUser := os.Getenv("DB_USER")
	// dbPassword := os.Getenv("DB_PASS")
	// dbHost := os.Getenv("DB_HOST")
	// dbName := os.Getenv("DB_NAME")
	// dbPort := os.Getenv("DB_PORT")

	log.Println("Attempting to connect to db")
	// dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", dbHost, dbUser, dbPassword, dbName, dbPort)
	dsn := "host=localhost user=postgres password=root dbname=hanger-craft port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Error),
		SkipDefaultTransaction: true,
		TranslateError:         true,
		// DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Println(err.Error())
		panic("Failed to Connect Database !")
	} else {
		log.Println("Database Connected Successfully !")
	}
	// log.Println("Attempting to migrate")
	// db.AutoMigrate(models.CartItem{}, models.Category{}, models.Inventory{}, models.Order{}, models.OrderItem{}, models.Payment{}, models.Product{}, models.Review{}, models.ShippingAddress{}, models.ShoppingCart{}, models.User{}, models.ProductAttribute{})
	// log.Println("Finished migration")
	DB = db
}
