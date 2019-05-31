package auth

import (
	"wave/internal/cookies"
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

func (db db) SignUp(username, password string) (cookie string, uid int64, err error) {
	// open transaction
	tx, err := db.Begin()
	if err != nil {
		return "", 0, err
	}

	// try to sign up
	if err := tx.QueryRow(`
		INSERT INTO users(name, password)
		VALUES ($1, $2)
		RETURNING uid
	`, username, password).Scan(&uid); err != nil {
		return "", 0, err
	}

	// create a new session
	cookie = cookies.GenerateCookie()
	if _, err := tx.Exec(`
		INSERT INTO sessions(uid, cookie)
		VALUES ($1, $2)
	`, uid, cookie); err != nil {
		return "", 0, err
	}
	return cookie, uid, tx.Commit()
}

func (db db) LogIn(username, password string) (cookie string, uid int64, err error) {
	// open transaction
	tx, err := db.Begin()
	if err != nil {
		return "", 0, err
	}

	// try to authorise
	if err := tx.QueryRow(`
		SELECT uid 
			FROM users
			WHERE name=$1 and password=$2
	`, username, password).Scan(&uid); err != nil {
		return "", 0, nil
	}

	// create a new session
	cookie = cookies.GenerateCookie()
	if _, err := tx.Exec(`
		INSERT INTO sessions(uid, cookie)
		VALUES ($1, $2)
	`, uid, cookie); err != nil {
		return "", 0, err
	}
	return cookie, uid, tx.Commit()
}

func (db db) LogOut(cookie string) (err error) {
	_, err = db.Exec(`
		DELETE FROM sessions
		WHERE cookie=$1
	`, cookie)
	return err
}

func (db db) GetSession(cookie string) (uid int64, err error) {
	if err := db.QueryRow(`
		SELECT uid
			FROM sessions
			WHERE cookie=$1
	`, cookie).Scan(&uid); err != nil {
		return 0, err
	}
	return uid, nil
}
