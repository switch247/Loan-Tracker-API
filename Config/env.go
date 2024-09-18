package Config

import (
	"log"
	"os"
)

var Port = ":8080"
var BASE_URL string

// Global variable to store the Env variables
var JwtSecret = []byte("your_jwt_secret")
var MONGO_CONNECTION_STRING string
var SERVICE_ID string
var TEMPLATE_ID string
var PUBLIC_KEY string
var Cloud_api_key string
var Cloud_api_secret string
var Data_Base_Name string

func Envinit() {

	JwtSecretKey := os.Getenv("JWT_SECRETE_KEY")
	if JwtSecretKey != "" {
		JwtSecret = []byte(JwtSecretKey)
	} else {
		JwtSecret = []byte("JwtSecretKey")
		log.Fatal("JWT secret key not configured")
	}
	// Read MONGO_CONNECTION_STRING from environment
	MONGO_CONNECTION_STRING = os.Getenv("MONGO_CONNECTION_STRING")
	if MONGO_CONNECTION_STRING == "" {
		MONGO_CONNECTION_STRING = "tst"
		log.Fatal("MONGO_CONNECTION_STRING is not set")
	}

	// Read PORT from environment
	Port = os.Getenv("PORT")
	if Port == "" {
		Port = ":8080"
	}
	BASE_URL = "http://localhost" + Port

	// serviceid
	SERVICE_ID = os.Getenv("SERVICE_ID")
	if SERVICE_ID == "" {
		log.Fatal("SERVICE_ID is not set")
	}
	// templateid
	TEMPLATE_ID = os.Getenv("TEMPLATE_ID")
	if TEMPLATE_ID == "" {
		log.Fatal("TEMPLATE_ID is not set")
	}
	// publickey
	PUBLIC_KEY = os.Getenv("PUBLIC_KEY")
	if PUBLIC_KEY == "" {
		log.Fatal("PUBLIC_KEY is not set")
	}

}
