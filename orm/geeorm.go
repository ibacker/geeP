package orm

import (
	"database/sql"
	"orm/log"
	"orm/session"
)

type Engine struct {
	db *sql.DB
}

func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}

	// db.Ping 测试数据库连接是否有效
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}

	e = &Engine{db: db}
	log.Info("Connected to ", driver, " database")
	return
}

func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error("Error closing database: ", err)
	}
	log.Info("Closed database")
}

func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db)
}
