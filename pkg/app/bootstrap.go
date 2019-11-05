package app

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var (
	// EnvFile Path
	EnvFile string
	//Port to run the service on
	Port string
	//JwtSecret Signing secret
	JwtSecret string
	//OAuthIssuer Address for OAuth provider
	OAuthIssuer string
	//ClientID 0Auth Client ID
	ClientID string
	//ClientSecret OAuth Client Secret
	ClientSecret string
	//RedirectURL OAuth Redirect URL
	RedirectURL string
	//DBType possible options are postgres, pg or sqlite
	DBType string
	//DBHost Host address of the database
	DBHost string
	//DBPort Port on which the database is running
	DBPort int
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
	h, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/")
	viper.AddConfigPath(h + "/.go-auth/")

	setDefaultValues()

	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalln("Config file not found")
		} else {
			log.Fatalln(err)
		}
	}

	for _, k := range viper.AllKeys() {
		if viper.IsSet(k) == false {
			log.Fatalf("Config Value required: %q \n", k)
		}
	}

	populateValues()
}

func setDefaultValues() {
	viper.SetDefault("PG_DB_SSL_MODE", "disable")
	viper.SetDefault("APP_PORT", ":3000")
	viper.SetDefault("JWT_SECRET", "MySuperSecret")
	viper.SetDefault("OAUTH_AUDIENCE", "go-auth")
}

func populateValues() {
	Port = ":" + strconv.Itoa(viper.GetInt("APP_PORT"))
	JwtSecret = viper.GetString("JWT_SECRET")
	OAuthIssuer = viper.GetString("OAUTH_ISSUER")
	ClientID = viper.GetString("OAUTH_CLIENT_ID")
	ClientSecret = viper.GetString("OAUTH_CLIENT_SECRET")
	RedirectURL = viper.GetString("OAUTH_CALLBACK_URL")
	DBType = viper.GetString("DB_TYPE")
	DBHost = viper.GetString("DB_HOST")
	DBPort = viper.GetInt("DB_PORT")
	DBUser = viper.GetString("DB_USER")
	DBPass = viper.GetString("DB_PASS")
	DBSSLMode = viper.GetString("PG_DB_SSL_MODE")

}
