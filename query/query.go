package query

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type query struct {
	conn   *sql.DB
	where  map[string]interface{}
	fields []string
	order  map[string]string
	group  []string
	save   map[string]interface{}
}

//mysql 配置结构
type DbConfig struct {
	Host    string
	Port    int
	User    string
	Passwd  string
	Dbname  string
	Charset string
}

// 获取一个配置结构
func NewDbConfig(h string, p int, u string, pwd string, db string, c string) *DbConfig {
	reutrn & DbConfig{
		Host:    h,
		Port:    p,
		User:    u,
		Passwd:  pwd,
		Dbname:  db,
		Charset: c,
	}
}

//获取一个新的query结构
func NewQuery(conn *sql.DB) *query {
	return &query{
		conn:   conn,
		where:  make(map[string]interface{}),
		fields: make([]string, 0, 10),
		order:  make(map[string]string),
		group:  make([]string, 0, 1),
		save:   make(map[string]interface{}),
	}
}

//根据配置，连接db并返回一个query结构
func NewQueryWithConfig(conf *DbConfig) (*query, error) {
	db, err := dbConnection(conf)
	if err != nil {
		return nil, err
	}
	return &query{
		conn:   db,
		where:  make(map[string]interface{}),
		fields: make([]string, 0, 10),
		order:  make(map[string]string),
		group:  make([]string, 0, 1),
		save:   make(map[string]interface{}),
	}

}

//mysql连接
func dbConnection(c *DbConfig) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", c.User, c.Passwd, c.Host, c.Port, c.Dbname, c.Charset))
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
