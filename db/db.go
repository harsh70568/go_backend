package db

import (
	"fmt"
	"go_edtech_backend/models"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	/* Load Config */
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to read env file: %v", err)
	}

	host := viper.GetString("DB_HOST")
	user := viper.GetString("DB_USER")
	password := viper.GetString("DB_PASSWORD")
	dbname := viper.GetString("DB_NAME")
	port := viper.GetString("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	DB = db

	DB.AutoMigrate(&models.User{})
	fmt.Println("Connected to database succesfully")
}

func GetPort() string {
	ServerPort := viper.GetString("PORT")
	return ServerPort
}

func GetJWTSecret() []byte {
	token := []byte(viper.GetString("JWT_SECRET"))
	return token
}
