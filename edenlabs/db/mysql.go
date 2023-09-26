package db

import (
	"fmt"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/orm"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type DBMysqlOption struct {
	Host                 string
	Port                 int
	Username             string
	Password             string
	Name                 string
	AdditionalParameters string
	MaxOpenConns         int
	MaxIdleConns         int
	ConnMaxLifetime      time.Duration
}

func NewMysqlDatabase(serviceName string, lgr *logrus.Logger, alias string, option DBMysqlOption) (err error) {
	// set orm debug is true
	orm.Debug = true
	orm.DebugLog = orm.NewLog(serviceName, lgr)

	// set orm time loc is UTC
	orm.DefaultTimeLoc = time.UTC

	// set orm connection
	source := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", option.Username, option.Password, option.Host, option.Port, option.Name, option.AdditionalParameters)
	err = orm.RegisterDataBase(alias, "mysql", source)
	if err != nil {
		return
	}

	// set option properties
	orm.SetMaxIdleConns(alias, option.MaxIdleConns)
	orm.SetMaxOpenConns(alias, option.MaxOpenConns)
	orm.ConnMaxLifetime(option.ConnMaxLifetime)

	return
}
