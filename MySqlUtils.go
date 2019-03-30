/**
Operation mysql Utils, based on go-mysql driver, tcp protocol, pass sql statement can operate mysql
With database connection pool, not included by default, please configure in the configuration file
Configure the number of database connections in the configuration file, not to configure negative and non-numeric, configuration 0 means not to use the database connection pool
The default value is: database: db1, ip: 127.0.0.1, host: 3306, user: root, password: root, number of database connections: 0
Author:zhanghaiqi
Date:2019.3.30 23:12:55

操作mysql的Utils，基于go-mysql驱动，tcp协议，传sql语句即可操作mysql
带有数据库连接池，默认不带有，请在配置文件中配置
在配置文件中配置数据库连接数，不可配置负数和非数字，配置0表示不使用数据库连接池
默认值为：数据库：db1,ip：127.0.0.1,host:3306,user:root,password:root,数据库连接数:0
作者：张海奇
日期：2019.3.30 23:12:55
 */
package haiqi

import (
	_ "../github.com/go-sql-driver/mysql"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

/**
Error handling without transaction

不带事务的错误处理
 */
func errs(e error)  {
	if e!=nil{
		panic(e)
	}
}
/**
Initialize connection pool method

调用初始化连接池方法
 */
func init()  {
	fmt.Println("Init mysqlPool!")
	InitDBPool()
}

/**
Initialize database connection pool

初始化数据库连接池
 */
func InitDBPool()  {
	sqlNumString := DefaultValue(GetFileValue("sql.num"),"0")
	sqlNum, e := strconv.Atoi(sqlNumString)
	if e!=nil {
		panic("配置文件中配置的连接数不为数字！")
	}
	if sqlNum<0{
		panic("配置文件中的连接应为正数！")
	}
	pool := sqlDBPool{}
	pool.sqlDB =make(chan  *sql.DB ,sqlNum)
	for i:=0; i<sqlNum ;i++  {
		dbName :=DefaultValue(GetFileValue("database.name"),"db1")
		user :=DefaultValue(GetFileValue("database.user"),"root")
		password :=DefaultValue(GetFileValue("database.password"),"root")
		host :=DefaultValue(GetFileValue("database.host"),"3306")
		ip :=DefaultValue(GetFileValue("database.ip"),"127.0.0.1")
		dbSql := user+":"+password+"@tcp("+ip+":"+host+")/"+dbName+"?charset=utf8";
		db, e := sql.Open("mysql",dbSql)
		errs(e)
		pool.sqlDB<-&(*db)
	}
	sqlDbPool = pool
}
/**
Consider the connection pool as a global variable

将连接池当做全局变量
 */
var sqlDbPool  sqlDBPool

/**
Get the connection, the parameter timeout is the timeout period for getting the connection.

获取连接,timeout为获取连接的超时时间
 */
func (p *sqlDBPool)getDB(timeout time.Duration)  (*sql.DB,error){
	select {
	case  res:=<-p.sqlDB:
		fmt.Println(res,reflect.TypeOf(res))
		return res,nil
	case <-time.After(timeout):
		return nil,errors.New("获取连接超时！请检查超时时间")
	default:
		dbname :=DefaultValue(GetFileValue("database.name"),"db1")
		user :=DefaultValue(GetFileValue("database.user"),"root")
		password :=DefaultValue(GetFileValue("database.password"),"root")
		host :=DefaultValue(GetFileValue("database.host"),"3306")
		ip :=DefaultValue(GetFileValue("database.ip"),"127.0.0.1")
		dbSql := user+":"+password+"@tcp("+ip+":"+host+")/"+dbname+"?charset=utf8";
		db, e := sql.Open("mysql",dbSql)
		return db,e
	}
}

/**
database connection pool structure

数据库连接池结构体
 */
type sqlDBPool struct {
	sqlDB chan *sql.DB
}


/**
Get a connection

获取连接
 */
func getDb() *sql.DB {
	db,err := sqlDbPool.getDB(5)
	errs(err)
	return db
}

/**
Return connection

归还连接
 */
func closeDb(sqldb *sql.DB)  {
	select{
	case sqlDbPool.sqlDB<- sqldb :
		return
	default:
		sqldb.Close()
	}
}

/**
Query encapsulation of MySQL
The first parameter is an SQL statement, and the latter parameter is a placeholder if it has one.
Return value: The result of the query, the map form in the slice and the result is string, each index in the slice is a row found, and the map is the data found.

mysql的查询封装
第一个参数为sql语句，后面参数若有则为占位符
返回值：查询出的结果，切片中套map形式且结果都为string，切片中每个索引为查出的一行，map为查出的数据
 */
func Query(query string,where ...string)[]map[string]string {
	db := getDb()
	defer closeDb(db)
	wheresql := make([]interface{}, len(where))
	for i:= range where {
		wheresql[i] = &where[i]
	}
	rows, e1 := db.Query(query,wheresql...)
	errs(e1)
	colums, e2 := rows.Columns()
	errs(e2)
	values := make([]sql.RawBytes, len(colums))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	mapString := []map[string]string{}
	var count = 0
	for rows.Next(){
		e1 = rows.Scan(scanArgs...)
		errs(e1)
		rowMap := make(map[string]string)
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col != nil {
				value = string(col)
				rowMap[colums[i]] = value
			}
		}
		mapString=append(mapString,rowMap)
		count++
	}
	return mapString;
}
/**
Error handling with transactions

带事务的错误处理
 */
func errsUpdate(e error , tx *sql.Tx)  {
	if e!=nil{
		tx.Rollback()
		panic(e)
	}
}

/**
Add, delete and amend items without transaction
Parameter: sqlUpdate: SQL statement, where condition

增删改单条，不带事务
参数：sqlUpdate:SQL语句，where 条件
 */
func Update(sqlUpdata string,where ...string)  {
	db := getDb()
	defer  closeDb(db)
	stmt, e := db.Prepare(sqlUpdata)
	errs(e)
	wheresql := make([]interface{}, len(where))
	for i:= range where {
		wheresql[i] = &where[i]
	}
	_, i := stmt.Exec(wheresql...)
	errs(i)
}

/**
Tape transaction
Add, delete and modify more than one, parameter: Slice in slice, the first slice records how many statements, the first index of the second slice is SQL statement, followed by conditions
Parameter eg: [[update xcx_userStatus set url ='a'where id = 1] [insert into xcx_userStatus (id) values (?) 6]]

带事务
增删改多条，参数：切片中套切片，第一个切片记录多少条语句，第2个切片第一个索引处为sql语句，后面为条件
参数 eg:[[update xcx_userStatus set url = 'a' where id = 1 ] [insert into xcx_userStatus (id) values (?)    6]]
 */
func Updates(sql [][]string)  {
	db := getDb()
	defer  closeDb(db)
	tx, i := db.Begin()
	errs(i)
	if len(sql)==0{
	fmt.Println("要增删改的数据为空")
		return
	}
	for i:=0;i<len(sql);i++{
		var sqlUpdata []string  =sql[i]
		if len(sqlUpdata)==0{
			fmt.Println("要增删改的数据为空")
			return
		}
		if len(sqlUpdata)==1{
			stmt, e := tx.Prepare(sqlUpdata[0])
			errsUpdate(e,tx)
			_, i := stmt.Exec()
			errsUpdate(i,tx)
		}else {
			stmt, e := tx.Prepare(sqlUpdata[0])
			stringSql := []string{}
			for i:=1;i<len(sqlUpdata) ;i++  {
				stringSql= append(stringSql, sqlUpdata[i])
			}
			errsUpdate(e,tx)
			wheresql := make([]interface{}, len(stringSql))
			for i:= range stringSql {
				wheresql[i] = &stringSql[i]
			}
			_, i := stmt.Exec(wheresql...)
			errsUpdate(i,tx)
		}
	}
	tx.Commit();
}

