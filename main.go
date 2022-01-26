package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func main() {

	// Global config map data
	if len(os.Getenv("ENV_CONFIG_MAP")) > 0 {
		_ = godotenv.Load(os.Getenv("ENV_CONFIG_MAP"))
	}

	// Load local .env files
	_ = godotenv.Load()

	// Create the DB instance
	gormDb, exception := setupDb()

	// We had a DB exception?
	if exception != nil {
		fmt.Printf("%s", exception.Error())
		return
	}

	// Always close the DB
	defer closeConnection(gormDb)

	// Build our container
	container := buildContainer(gormDb)

	// Reference our echo instance and create it early
	e := echo.New()

	// Get the API calls
	createEndpoints(e, container)

	e.Logger.Fatal(e.Start(":3000"))
}

func closeConnection(gormDb *gorm.DB) {
	log.Print("Closed database connection")

	// Get the DB connection
	sqlDb, _ := gormDb.DB()

	// Close the DB connection
	_ = sqlDb.Close()
}

func setupDb() (*gorm.DB, error) {

	// Load local .env files
	dbHost := os.Getenv("DB_HOST")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbDatabase := os.Getenv("DB_DATABASE")
	dbPort := os.Getenv("DB_PORT")
	dbCaCert := os.Getenv("DB_CA_CERT")

	// Make sure we have the right values
	if dbHost == "" {
		exception := errors.New("Missing env DB_HOST")
		return nil, exception
	}

	if dbUsername == "" {
		exception := errors.New("Missing env DB_USERNAME")
		return nil, exception
	}

	if dbPassword == "" {
		exception := errors.New("Missing env DB_PASSWORD")
		return nil, exception
	}

	if dbDatabase == "" {
		exception := errors.New("Missing env DB_DATABASE")
		return nil, exception
	}

	// Are we running in test mode?
	isTestMode := os.Getenv("IS_TEST") == "true"

	// Standard DSN for all replicas
	dbType := "mysql"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&timeout=2s", dbUsername, dbPassword, dbHost, dbPort, dbDatabase)

	// Are we in test mode?
	if isTestMode {
		dsn = ":memory:"
		dbType = "sqlite3"
	}

	// Debug connection line
	fmt.Println(fmt.Sprintf("[DB] Connecting %s to: %s with %s@%s", dbType, dbDatabase, dbUsername, dbHost), nil)

	// If we're not in test mode, run under SSL
	if dbCaCert != "" {

		// Append our custom TLS handler
		if !isTestMode {
			dsn += "&tls=custom"
		}

		// Create the root cert pool
		rootCertPool := x509.NewCertPool()
		var exception error
		var pem []byte

		// Read in our certificate file
		if isTestMode {
			pem = []byte(dbCaCert)
		} else {
			pem, exception = ioutil.ReadFile(dbCaCert)
		}

		// We had an error, handle it
		if exception != nil {
			return nil, exception
		}

		// Append the PEM certificate
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			// Create an error object
			exception := errors.New("Failed to append PEM.")
			return nil, exception
		}

		// Register our TLS configuration
		_ = mysql.RegisterTLSConfig("custom", &tls.Config{
			ServerName: dbHost,
			RootCAs:    rootCertPool,
		})
	}

	var dbConnection gorm.Dialector

	switch dbType {
	case "sqlite3":
		dbConnection = sqlite.Open(dsn)
	case "mysql":
		dbConnection = gormMysql.Open(dsn)
	}
	gormConfig := &gorm.Config{}

	// Debug mode
	if os.Getenv("IS_DEBUG") == "true" {
		gormConfig.Logger = gormLogger.Default.LogMode(gormLogger.Info)
	}

	// Create a Gorm-DB instance
	gormDb, exception := gorm.Open(dbConnection, gormConfig)

	// We had a DB exception?
	if exception != nil {
		return nil, exception
	}

	// Preload by default
	if os.Getenv("IS_PRELOAD") == "true" {
		gormDb.Set("gorm:auto_preload", true)
	}

	return gormDb, nil
}
