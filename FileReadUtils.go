/**
Read the Utils of the configuration file, pay attention to the configuration file storage path
Author:zhanghaiqi
Date:2019.3.30 23:12:55

读取配置文件的Utils，注意配置文件存储路径
作者：张海奇
日期：2019.3.30 23:12:55
 */
package haiqi

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)
/**The directory where the configuration file is located should now be in /src/haiqi/config/application.properties
配置文件所在目录，现应在/src/haiqi/config/application.properties**/
const (
	FN="/src/haiqi/config/application.properties"
)
/**
Read the configuration file. The parameter str is the key of the configuration file. The return value is the corresponding value. If not, it returns "".

读取配置文件，参数str为配置文件的key，返回值为对应的value，无则返回“”
 */
func GetFileValue(str string) string{
	s, _:= filepath.Abs(filepath.Dir(""))
	f, err := os.Open(s+FN)
	defer f.Close()
	if err != nil {
		os.Exit(-1)
	}
	buf := bufio.NewReader(f)
	for {
		line, _ := buf.ReadString('\n')
		if line ==""{
			return ""
		}
		line = strings.TrimSpace(line)
		strs:=strings.Split(line,"=")
		if strings.TrimSpace(strs[0])==str{
			return strings.TrimSpace(strs[1])
		}
	}
	return ""
}
