package db

import (
	"context"
	"strconv"

	"github.com/go-pg/pg/v9"

	"github.com/imtanmoy/authz/config"
)

var DB *pg.DB

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	//fmt.Println(q.FormattedQuery())
	return nil
}

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
	db.AddQueryHook(dbLogger{})
	DB = db
	return nil
}

func Shutdown() error {
	err := DB.Close()
	return err
}
