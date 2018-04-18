package query

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type query struct {
	conn   *sql.DB
	where  []where
	fields []string
	order  map[string]string
	group  []string
	save   []map[string]interface{}
	offset int
	limit  int
}

type where struct {
	field   string
	compare string
	value   interface{}
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
		offset: int,
		limit:  int,
	}
}

//根据配置，连接db并返回一个query结构
func NewQueryWithConfig(conf *DbConfig) (*query, error) {
	db, err := dbConnection(conf)
	if err != nil {
		return nil, err
	}
	return NewQuery(db), nil

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

//注入sql的where条件
func (q *query) Where(f string, v interface{}, c string) *query {
	var w where = where{
		field:   f,
		value:   v,
		compare: c,
	}
	q.where = append(q.where, w)
	return q
}

//设置要查询的字段
func (q *query) Select(fields []string) *query {
	for _, v := range fields {
		if q.inFields(v) {
			continue
		}
	}
	return q
}

//设置limit值
func (q *query) Limit(i ...int) *query {
	if len(i) == 1 {
		q.limit = i[0]
	}
	if len(i) == 2 {
		q.offset = i[0]
		q.limit = i[1]
	}
	return q
}

//设置offset
func (q *query) Offset(i int) *query {
	q.offset = i
	return q
}

//判断是否含有元素
func (q *query) inFields(field string) bool {
	for _, v := range q.fields {
		if v == field {
			return true
		}
	}
	return false
}
