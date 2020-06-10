package database

import (
	"fmt"
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type DB struct {
	Host       string
	Port       string
	UserName   string
	Password   string
	DBName     string
	Connection *gorm.DB
}

func NewDB() *DB {
	conf := config.GetInstance()
	return newDB(&DB{
		Host:     conf.DB.Host,
		Port:     conf.DB.Port,
		UserName: conf.DB.UserName,
		Password: conf.DB.Password,
		DBName:   conf.DB.DBName,
	})
}

func newDB(d *DB) *DB {
	param := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		d.Host, d.Port, d.UserName, d.DBName, d.Password)
	db, err := gorm.Open("postgres", param)
	if err != nil {
		panic(err)
	}
	d.Connection = db
	return d
}

func (db *DB) Begin() *gorm.DB {
	return db.Connection.Begin()
}

func (db *DB) Connect() *gorm.DB {
	return db.Connection
}
