package app

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var (
	// EnvFile Path
	EnvFile string
	//Port to run the service on
	Port string
	//JwtSecret Signing secret
	JwtSecret string
	//ClientID 0Auth Client ID
	ClientID string
	//ClientSecret 0Auth Client Secret
	ClientSecret string
	//RedirectURL 0Auth Redirect URL
	RedirectURL string
	//DBType possible options are postgres, pg or sqlite
	DBType string
	//DBHost Host address of the database
	DBHost string
	//DBPort Port on which the database is running
	DBPort string
	//DBUser Username to connect to the database
	DBUser string
	//DBPass Password to connect to the databse
	DBPass string
	//DBSSLMode possible options are enable/disable.
	DBSSLMode string
)

// IDTokenClaims struct is GORM model to write to the database
// The values for this model are sourced from the response
// payload of the OAuth provider and are marshalled into this struct
type IDTokenClaims struct {
	gorm.Model
	ISS           string `gorm:"type:varchar(1000)" json:"iss"`
	Email         string `gorm:"type:varchar(1000);unique;not null" json:"email"`
	EmailVerified bool   `gorm:"type:boolean" json:"email_verified"`
	Name          string `gorm:"type:varchar(1000)" json:"name"`
	Picture       string `gorm:"type:varchar(1000)" json:"picture"`
	GivenName     string `gorm:"type:varchar(1000)" json:"given_name"`
	FamilyName    string `gorm:"type:varchar(1000)" json:"family_name"`
	Locale        string `gorm:"type:varchar(1000)" json:"locale"`
	Iat           int    `gorm:"type:integer" json:"iat"`
	Exp           int    `gorm:"type:integer" jsong:"exp"`
}

// Init Sources & Initializes configuration
// required for starting the service.
func init() {
	fmt.Println("Populating environment variables...")
	requiredEnvVariables := [10]string{"APP_PORT",
		"JWT_SECRET",
		"OAUTH_CLIENT_ID",
		"OAUTH_CLIENT_SECRET",
		"DB_TYPE",
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASS",
		"DB_SSL_MODE",
	}

	fmt.Println("Sourcing config from environment variables or local file .env")
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())

	}

	Port = os.Getenv("APP_PORT")
	JwtSecret = os.Getenv("JWT_SECRET")
	ClientID = os.Getenv("OAUTH_CLIENT_ID")
	ClientSecret = os.Getenv("OAUTH_CLIENT_SECRET")
	RedirectURL = os.Getenv("OAUTH_CALLBACK_URL")
	DBType = os.Getenv("DB_TYPE")
	DBHost = os.Getenv("DB_HOST")
	DBPort = os.Getenv("DB_PORT")
	DBUser = os.Getenv("DB_USER")
	DBPass = os.Getenv("DB_PASS")
	DBSSLMode = os.Getenv("DB_SSL_MODE")
	fmt.Println(ClientID)
	fmt.Println(ClientSecret)
	fmt.Println(RedirectURL)

	for _, v := range requiredEnvVariables {
		ev := os.Getenv(v)
		if len(ev) < 1 {
			log.Fatalf("app: Missing Env variable: %s", v)
		}
	}
}
