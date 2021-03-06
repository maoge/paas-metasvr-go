package pool

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/maoge/paas-metasvr-go/pkg/utils"
)

type DbPool struct {
	DB              *sqlx.DB
	Addr            string        `json:"addr,omitempty"`
	Username        string        `json:"username,omitempty"`
	Password        string        `json:"password,omitempty"`
	DbName          string        `json:"dbname,omitempty"`
	DbType          string        `json:"dbtype,omitempty"`
	MaxOpenConns    int           `json:"max_open_conns,omitempty"`
	MaxIdleConns    int           `json:"max_idle_conns,omitempty"`
	ConnTimeout     int           `json:"conn_timeout,omitempty"`
	ReadTimeout     int           `json:"read_timeout,omitempty"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime,omitempty"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time,omitempty"`
}

func (pool *DbPool) Connect() bool {
	// postgresql: sql.Open("pgx","postgres://localhost:5432/postgres")
	// mysql: username:password@tcp(ip:port)/dbname?timeout=%ds&readTimeout=%ds&charset=utf8
	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=%ds&readTimeout=%ds",
		pool.Username, pool.Password, pool.Addr, pool.DbName, pool.ConnTimeout, pool.ReadTimeout)

	result := false
	db, err := sqlx.Connect(pool.DbType, connStr)
	if err == nil {
		info := fmt.Sprintf("database: %v connect OK", pool.Addr)
		utils.LOGGER.Info(info)

		db.SetMaxOpenConns(pool.MaxOpenConns)
		db.SetMaxIdleConns(pool.MaxIdleConns)
		db.SetConnMaxLifetime(pool.ConnMaxLifetime)
		db.SetConnMaxIdleTime(pool.ConnMaxIdleTime)

		if db.Ping() == nil {
			pool.DB = db
			result = true
		}
	} else {
		err := fmt.Sprintf("database connect fail, %v", err.Error())
		utils.LOGGER.Error(err)
	}

	return result
}

func (pool *DbPool) Release() {
	if pool.DB != nil {
		pool.DB.Close()
		pool.DB = nil
	}
}
