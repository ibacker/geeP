package orm

import (
	"database/sql"
	"orm/dialect"
	"orm/log"
	"orm/session"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
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

	dia, ok := dialect.GetDialect(driver)
	if !ok {
		log.Error("dialect %s not found", driver)
		return
	}
	e = &Engine{db: db, dialect: dia}
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
	return session.New(engine.db, engine.dialect)
}
