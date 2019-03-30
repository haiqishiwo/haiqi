/**
The Operation Package Utils of redis, Based on the Development of redigo Driver
With database connection pool, not by default, configure it in the configuration file
Configure the number of database connections in the configuration file, not negative and non-numeric. Configure 0 means that database connection pool is not used.
The default values are: ip: 127.0.0.1, host: 3306, number of database pool connections: 0
Author:zhanghaiqi
Date:2019.3.30 23:12:55

redis的操作封装Utils，基于redigo驱动开发
带有数据库连接池，默认不带有，请在配置文件中配置
在配置文件中配置数据库连接数，不可配置负数和非数字，配置0表示不使用数据库连接池
默认值为：ip：127.0.0.1,host:3306,数据库池连接数:0
作者：张海奇
日期：2019.3.30 23:12:55
 */
package haiqi

import (
	"../github.com/gomodule/redigo/redis"
	"errors"
	"fmt"
	"strconv"
	"time"
)

func init()  {
	fmt.Println("Init redisPool!")
	InitRedisDBPool()
}

/**
Redis error handling method, if the error is redigo: nil returned, return ""

redis的错误处理方式，若错误为redigo: nil returned则返回""
 */
func errRedis(e error) string {
	if e!=nil{
		sprintf := fmt.Sprintf("%s", e)
		if sprintf == "redigo: nil returned"{
			return ""
		}
		panic(e)
		return "2"
	}
	return "1"
}


/**
Redis database connection pool structure

Redis数据库连接池结构体
 */
type redisDBPool struct {
	redisConn chan *redis.Conn
}
/**
Treat the connection pool as a global variable

将连接池当做全局变量
 */
var redisDbPool redisDBPool


/**
Initialize the Redis database connection pool

初始化Redis数据库连接池
 */
func InitRedisDBPool()  {
	RedisNumString := DefaultValue(GetFileValue("redis.num"),"0")
	RedisNum, e := strconv.Atoi(RedisNumString)
	if e!=nil {
		panic("配置文件中配置的连接数不为数字！")
	}
	if RedisNum<0{
		panic("配置文件中的连接应为正数！")
	}
	pool := redisDBPool{}
	pool.redisConn =make(chan  *redis.Conn ,RedisNum)
	for i:=0;i<RedisNum ;i++  {
		ip := DefaultValue(GetFileValue("redis.ip"),"127.0.0.1")
		host := DefaultValue(GetFileValue("redis.host"),"6379")
		conn, e := redis.Dial("tcp", ip+":"+host)
		errs(e)
		pool.redisConn<-&conn
	}
	redisDbPool = pool
}

/**
Get the connection, the parameter timeout is the timeout period for getting the connection.

获取连接，参数timeout为获取连接的超时时间
 */
func (p *redisDBPool)getRedisDB(timeout time.Duration)  (*redis.Conn,error){
	select {
	case  res:=<-p.redisConn:
		return res,nil
	case <-time.After(timeout):
		return nil,errors.New("获取连接超时！请检查超时时间")
	default:
		ip := DefaultValue(GetFileValue("redis.ip"),"127.0.0.1")
		host := DefaultValue(GetFileValue("redis.host"),"6379")
		conn, e := redis.Dial("tcp", ip+":"+host)
		return &conn,e
	}
}

/**
Get a Redis connection

获取Redis连接
 */
func getRedisDb() *redis.Conn {
	conn, e := redisDbPool.getRedisDB(5)
	errs(e)
	return conn
}

/**
Return connection

归还连接
 */
func closeRedisConn(redisDb *redis.Conn)  {
	select{
	case redisDbPool.redisConn<- redisDb :
		return
	default:
		(*redisDb).Close()
	}
}



/**
Set the value of redis, string form

设置redis的值，string形式
 */
func SetRedis(key string,value string)  {
	conn := getRedisDb()
	defer closeRedisConn(conn)
	_ , err := (*conn).Do("SET", key, value)
	errs(err)
}

/**
Get the value according to the redis key, with the parameter key as the key

根据redis的键获取值，参数key为键
 */
func GetRedis(key string) string {
	conn := getRedisDb()
	defer closeRedisConn(conn)
	s, e1 := redis.String((*conn).Do("GET", key))
	i := errRedis(e1)
	if i==""{
		s = i
	}
	return s
}
/**
Set the value of redis and its expiration time. The parameter time is a few seconds. Pass 1 is deleted after 1 second.

设置redis的值及其失效时间，参数time为几秒，传1就是1秒后删除
 */
func SetRedisTime(key string,value string,time int){
	conn := getRedisDb()
	defer closeRedisConn(conn)
	_, e1 := (*conn).Do("SET", key,value, "EX", time)
	errs(e1)
}
/**
Judging whether there is a value in redis based on key

根据key判断redis中是否存在值
 */
func ExistsRedisKey(key string) bool {
	conn := getRedisDb()
	defer closeRedisConn(conn)
	b, e1 := redis.Bool((*conn).Do("EXISTS", key))
	errs(e1)
	return b
}
/**
Delete redis data by key

根据键删除redis数据
 */
func DeleteRedis(key string)  {
	conn := getRedisDb()
	defer closeRedisConn(conn)
	_ , e1 := (*conn).Do("DEL", key)
	errs(e1)
}


/**
Redis sets data in list form, adding list from left

redis设置数据，list形式，从左加入list
 */
func SetRedisLpush(key string,value string)  {
	conn := getRedisDb()
	defer closeRedisConn(conn)
	_, e1 := (*conn).Do("lpush", key, value)
	errs(e1)
}

/**
Redis sets the data in list form, adding list from right

redis设置数据，list形式，从右加入list
 */
func SetRedisRpush(key string,value string)  {
	conn := getRedisDb()
	defer closeRedisConn(conn)
	_, e1 := (*conn).Do("rpush", key, value)
	errs(e1)
}


/**
List of redis deletes elements from left

redis的list从左开始删元素
 */
func DeleteRedisLpop(key string)  {
	conn := getRedisDb()
	defer closeRedisConn(conn)
	_, e1 := (*conn).Do("lpop", key)
	errs(e1)
}
/**
The list of redis deletes the element from the right, and the key is the list name.

redis的list从右开始删元素,key为list名
 */
func DeleteRedisRpop(key string)  {
	conn := getRedisDb()
	defer closeRedisConn(conn)
	_, e1 := (*conn).Do("rpop", key)
	errs(e1)
}



/**
Get the list value of redis, the strat is the index at the beginning of the list, and the end is the index at the end of the list.
The return value is [][]string

获取redis的list 值,strat为list开始时索引,end为list结束时索引,
返回值为string切片中包含string切片
 */
func GetRedisLists(key string,start int,end int)[][]string {
	conn := getRedisDb()
	defer closeRedisConn(conn)
	result, e1 := redis.Values((*conn).Do("lrange", key, start,end))
	errs(e1)
	resultString := [][]string{}
	for _,v := range result{
		s := v.([]uint8)
		resultStrings := []string{}
		for _,v1:= range s{
		resultStrings=append(resultStrings,string(v1))
		}
		resultString=append(resultString,resultStrings)
	}
	return resultString
}
/**
Gets the list form of redis called key and the return value is string slice

获取redis的list形式名为key的值，返回值为string切片
 */
func GetRedisList(key string,start int,end int)[]string {
	conn := getRedisDb()
	defer closeRedisConn(conn)
	result, e1 := redis.Values((*conn).Do("lrange", key, start,end))
	errs(e1)
	resultString := []string{}
	for _,v := range result{
		s := v.([]uint8)
		for _,v1:= range s{
			resultString=append(resultString,string(v1))
		}
	}
	return resultString
}

/**
Redis set data, set form

redis设置数据，set形式
 */
func SetRedisSadd(key string,value string)  {
	conn := getRedisDb()
	defer closeRedisConn(conn)
	_, e1 := (*conn).Do("sadd", key, value)
	errs(e1)
}


/**
Customize the statement stored in Redis,
Parameter: statement executed by typeString: redis, key: redis stored in, value: redis stored in

自定义给Redis存的语句，
参数：typeString：redis执行的语句 ，key:redis存入的key，value:redis存入的value
 */
func SetRedisCustomize(typeString string,key string,value string)  {
	conn := getRedisDb()
	defer closeRedisConn(conn)
	_, e1 := (*conn).Do(typeString, key, value)
	errs(e1)
}

/**
Redis setting data, sortedset form

redis设置数据，sortedset形式
 */
func GetRedisZadd(key string,value string)  {
	conn := getRedisDb()
	defer closeRedisConn(conn)
	_, e1 :=  (*conn).Do("zadd", key, value)
	errs(e1)
}

///**
//redis删除，set形式
// */
//func DeleteRedisSadd(key string,value string)  {
//	ip := GetFileValue("redis.ip")
//	host := GetFileValue("redis.host")
//	conn, e := redis.Dial("tcp", ip+":"+host)
//	errs(e)
//	defer conn.Close()
//	_, e1 := conn.Do("sadd", key, value)
//	errs(e1)
//}
///**
//redis设置数据，hash形式
// */
//func GetRedisHash(key string,value string)  {
//	ip := GetFileValue("redis.ip")
//	host := GetFileValue("redis.host")
//	conn, e := redis.Dial("tcp", ip+":"+host)
//	errs(e)
//	defer conn.Close()
//	_, e1 := conn.Do("hset", key, value)
//	errs(e1)
//}
///**
//redis删除数据，hash形式
// */
//func DeleteRedisHash(key string)  {
//	ip := GetFileValue("redis.ip")
//	host := GetFileValue("redis.host")
//	conn, e := redis.Dial("tcp", ip+":"+host)
//	errs(e)
//	defer conn.Close()
//	_, e1 := conn.Do("hdel", key)
//	errs(e1)
//}



