package platform

import (
	"wave/internal/service"

	"github.com/jackc/pgx"
)

type db struct {
	*pgx.ConnPool

	log  service.ILogger
	conf Config
}

func newDB(log service.ILogger, conf Config) *db {
	conn, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     conf.DB.Host,
			Port:     uint16(conf.DB.Port),
			User:     conf.DB.User,
			Password: conf.DB.Password,
			Database: conf.DB.Database,
		},
		MaxConnections: 50,
	})
	service.PanicIf(err)
	return &db{
		ConnPool: conn,
		conf:     conf,
		log:      log,
	}
}

// --------------------------|
