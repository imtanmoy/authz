package db

import (
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authz/config"
	"strconv"
)

var DB *pg.DB

func InitDB() error {
	db := pg.Connect(&pg.Options{
		User:     config.Conf.DB.USERNAME,
		Password: config.Conf.DB.PASSWORD,
		Database: config.Conf.DB.DBNAME,
		Addr:     config.Conf.DB.HOST + ":" + strconv.Itoa(config.Conf.DB.PORT),
	})
	var n int
	_, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
	if err != nil {
		return err
	}
	DB = db
	return nil
}

func Shutdown() error {
	err := DB.Close()
	return err
}
