# 常用的工具类和数据库操作工具类
1.基于go-mysql驱动开发的工具类，只用传递sql语句即可操作mysql数据库，可配置连接池参数（MySqlUtils）

2.自己开发整理的常用工具类（Utilts）

3.基于redigo驱动开发的工具类，操作redis数据库，可配置连接池参数（redisUtils）

4.操作配置文件的工具类（FileReadUtils）



1.MySqlUtils:(注：使用该工具类需安装go-mysql驱动，地址:https://github.com/go-sql-driver/mysql)。

操作mysql的Utils，基于go-mysql驱动，tcp协议，传sql语句即可操作mysql
带有数据库连接池，默认不带有，请在配置文件中配置
在配置文件中配置数据库连接数，不可配置负数和非数字，配置0表示不使用数据库连接池
默认值为：数据库：db1,ip：127.0.0.1,host:3306,user:root,password:root,数据库连接数:0

方法
（1）Query(query string,where ...string)[]map[string]string 

mysql的查询封装，第一个参数为sql语句，后面参数若有则为占位符
返回值：查询出的结果，切片中套map形式且结果都为string，切片中每个索引为查出的一行，map为查出的数据

（2） Update(sqlUpdata string,where ...string) 

增删改单条，不带事务
参数：sqlUpdate:SQL语句，where 条件

（3） Updates(sql [][]string) 

带事务，增删改多条，参数：切片中套切片，第一个切片记录多少条语句，第2个切片第一个索引处为sql语句，后面为条件
参数 eg:[[update xcx_userStatus set url = 'a' where id = 1 ] [insert into xcx_userStatus (id) values (?)    6]]
 
 
 2.Utilts
（1）StringSlice(str string) []string

将字符串转化为字符串切片返回

（2）GetStringIndex(str string,indexStr string) int

获取字符串中含有某个字符的第一个索引。
参数：str 整个字符串。
     indexStr 整个字符串中要包含的字符。
返回：第一个索引，无则返回-1。

（3） GetStringLastIndex(str string,indexStr string) int

获取字符串中含有某个字符的最后一个索引。
参数：str 整个字符串。
     indexStr 整个字符串中要包含的字符。
返回：第一个索引，无则返回-1。

（4） GetStringIndexForNum(str string,indexStr string,indexS int) int 

获取字符串中含有某个字符的第自定义个索引。
参数：str 整个字符串。
     indexstr 整个字符串中要包含的字符。
     indexs  字符串中要包含的第几个位置字符。
返回：第indexs个索引，无则返回-1。

（5）DefaultValue(value string,defaultBalue string)string

模拟三元表达式，value不为""则返回本身，否则返回默认值



3.redisUtils(注：使用该工具类需安装redigo驱动，地址:https://github.com/gomodule/redigo)

（1）SetRedis(key string,value string)
设置redis的值，string形式

（2） GetRedis(key string) string
根据redis的键获取值，参数key为键

（3）SetRedisTime(key string,value string,time int)
设置redis的值及其失效时间，参数time为几秒，传1就是1秒后删除
 
（4）ExistsRedisKey(key string) bool 
根据key判断redis中是否存在值

（5）DeleteRedis(key string)
根据键删除redis数据

（6）SetRedisLpush(key string,value string) 
redis设置数据，list形式，从左加入list

（7）SetRedisRpush(key string,value string) 
redis设置数据，list形式，从右加入list

（8） DeleteRedisLpop(key string) 
redis的list从左开始删元素

（9）DeleteRedisRpop(key string)  
redis的list从右开始删元素,key为list名
 
 （10）GetRedisLists(key string,start int,end int)[][]string 
 获取redis的list 值,strat为list开始时索引,end为list结束时索引,
返回值为string切片中包含string切片

（11）GetRedisList(key string,start int,end int)[]string 
获取redis的list形式名为key的值，返回值为string切片

（12） SetRedisSadd(key string,value string)  
redis设置数据，set形式

（13）SetRedisCustomize(typeString string,key string,value string) 
自定义给Redis存的语句，
参数：typeString：redis执行的语句 ，key:redis存入的key，value:redis存入的value

（14）GetRedisZadd(key string,value string)
redis设置数据，sortedset形式

4.FileReadUtils

(1)const (
	FN="/src/haiqi/application.properties"
)
配置文件所在目录，现应在/src/haiqi/config/application.properties

(2)GetFileValue(str string) string
读取配置文件，参数str为配置文件的key，返回值为对应的value，无则返回“”
 
 









