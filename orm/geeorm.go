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

type TxFunc func(session2 *session.Session) (interface{}, error)

func (engine *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := engine.NewSession()
	// 启用事务，对 tx 赋值
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		// 有异常，回滚事务
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p)
		} else if err != nil {
			_ = s.Rollback()
		} else {
			err = s.Commit()
		}
	}()
	// 调用转入的 函数 f
	return f(s)
}
