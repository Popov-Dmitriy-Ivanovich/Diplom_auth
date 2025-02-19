package models

import (
	"os"
	"sync"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbConnection *gorm.DB = nil
var dbInitMutex sync.Mutex

func initDb() error {
	dbInitMutex.Lock()
	defer dbInitMutex.Unlock()

	if dbConnection != nil {
		return nil
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbName + " port=" + port
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(ALL_MODELS...)

	if os.Getenv("ADMIN_PAS") !=  "" {
		pHash, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("ADMIN_PAS")),bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		admin := User{
			Login: "admin",
			Password: pHash,
			AccessRights: AR_ALL,
		}
		if err := db.FirstOrCreate(&admin).Error; err != nil {
			panic(err)
		}
		if err := db.Save(&admin).Error; err != nil {
			panic(err)
		}
	}
	

	
	dbConnection = db
	return nil
}

func GetDb() *gorm.DB {
	if dbConnection == nil {
		if err := initDb(); err != nil {
			panic(err)
		}
		// Get generic database object sql.DB to use its functions
		sqlDB, err := dbConnection.DB()
		if err != nil {
			panic(err)
		}

		// SetMaxOpenConns sets the maximum number of open connections to the database.
		sqlDB.SetMaxOpenConns(20)
		//sqlDB.SetMaxIdleConns(2)
	}
	return dbConnection
}