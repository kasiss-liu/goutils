package query

import (
	"fmt"
	"testing"
)

func TestQuery(t *testing.T) {
	config := NewDbConfig("127.0.0.1", 3306, "root", "123qwe", "test", "utf8")
	query, err := NewQueryWithConfig(config)
	if err != nil {
		t.Fatal(err.Error())
	}
	//test queryOne
	user := query.Table("user").Where("id", 1, "=").Select([]string{"id", "name"}...).QueryOne()
	fmt.Println(user)

	//test query
	users := query.Table("user").Where("id", []int{1, 2}, "in").Select("id", "name", "age").GroupBy("name").OrderBy("id", "desc").Offset(0).Limit(1).Query()
	fmt.Println(users)
}

//测试开启事务下的增删改查
func TestTxExecute(t *testing.T) {
	config := NewDbConfig("127.0.0.1", 3306, "root", "123qwe", "test", "utf8")
	query, err := NewQueryWithConfig(config)
	if err != nil {
		t.Fatal(err.Error())
	}
	//开启事务
	err = query.BeginTransaction()
	if err != nil {
		t.Fatal(err.Error())
	}
	//写入一条新数据
	insertID := query.Table("user").Create(map[string]interface{}{"name": "john1", "age": 28})
	if insertID == 0 {
		t.Fatal(query.GetLastError())
	} else {
		fmt.Println("insertId:", insertID)
	}
	//查询插入的新数据
	user := query.Table("user").Where("id", insertID, "=").QueryOne()
	fmt.Println("user:", user)
	//修改插入的数据
	affectedRows := query.Table("user").Where("name", "john1", "=").Update(map[string]interface{}{"age": 29})
	if affectedRows == 0 {
		t.Error(query.GetLastError())
	} else {
		fmt.Println("affectedRows:", affectedRows)
	}
	//删除插入的数据
	deleteRows := query.Table("user").Where("id", insertID, "=").Delete()
	if deleteRows == 0 {
		t.Error(query.GetLastError())
	} else {
		fmt.Println("deleteRows:", deleteRows)
	}
	//重置query
	query.Reset()
	//关闭db链接
	query.Close()
}

//测试批量写入
func TestCreateBatch(t *testing.T) {
	config := NewDbConfig("127.0.0.1", 3306, "root", "123qwe", "test", "utf8")
	query, err := NewQueryWithConfig(config)
	if err != nil {
		t.Fatal(err.Error())
	}
	//开启事务
	err = query.BeginTransaction()
	if err != nil {
		t.Fatal(err.Error())
	}
	//写入一条新数据
	insertID := query.Table("user").CreateBatch([]map[string]interface{}{{"name": "john1", "age": 28}, {"name": "john2", "age": 29}})
	if insertID == 0 {
		t.Fatal(query.GetLastError())
	} else {
		fmt.Println("insertId:", insertID)
	}

	//打印sql
	fmt.Println("Last Sql:", query.GetLastSQL())

	//事务回滚
	query.Rollback()

	//关闭db链接
	query.Close()
}

func TestQueryRaw(t *testing.T) {
	config := NewDbConfig("127.0.0.1", 3306, "root", "123qwe", "test", "utf8")
	query, err := NewQueryWithConfig(config)
	if err != nil {
		t.Fatal(err.Error())
	}
	//开启事务
	query.BeginTransaction()
	//开启sql日志
	query.StartLogQuery()
	//正常查询
	user := query.QueryRaw("select * from user where id in (?,?)", []interface{}{1, 2}...)
	fmt.Println("user:", user)
	//错误查询
	users := query.QueryRaw("select * from users where id in (?,?)", []interface{}{1, 2}...)
	fmt.Println("users:", users)
	//获取最后一个错误信息
	fmt.Println(query.GetLastError())
	//获取最后一条执行的sql
	fmt.Println("last sql:", query.GetLastSQL())
	//获取所有错误
	errors := query.GetAllErrors()
	fmt.Println("errors:", errors)
	//清空错误日志
	query.ResetErrorLog()
	//停止sql日志
	query.StopLogQuery()
	//打印sql日志
	fmt.Println("sqls:", query.GetSQLs())
	//清空sql日志
	query.ResetQueryLog()
	//事务提交
	query.Commit()
}

func TestExecRaw(t *testing.T) {
	config := NewDbConfig("127.0.0.1", 3306, "root", "123qwe", "test", "utf8")
	query, err := NewQueryWithConfig(config)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = query.Ping()
	if err != nil {
		t.Fatal(err.Error())
	}
	//开启事务
	query.BeginTransaction()

	//正常查询
	result := query.ExecRaw("update `user` set `name` = ? where id = ?", "Ashi", "1")
	lastInsertID, _ := result.LastInsertId()
	rowsAffected, _ := result.RowsAffected()
	fmt.Println("LastInsertId:", lastInsertID, "RowsAffected:", rowsAffected)

	//获取最后一条执行的sql
	fmt.Println("last sql:", query.GetLastSQL())
	//获取所有错误
	errors := query.GetAllErrors()
	fmt.Println("errors:", errors)
	//清空错误日志
	query.ResetErrorLog()

	//事务提交
	query.Rollback()
}
