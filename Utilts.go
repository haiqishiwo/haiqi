/**
Organize and develop Utlis for various uses
Author:zhanghaiqi
Date:2019.3.30 23:12:55


整理和开发各种使用的Utlis
作者：张海奇
日期：2019.3.30 23:12:55
*/

package haiqi

/**
Convert a string to a string slice back

将字符串转化为字符串切片返回
 */
func StringSlice(str string) []string {
	result := []string{}
	for _,v := range  str {
		result = append(result,string(v))
	}
	return  result;
}
/**
Get the first index of a character in a string
Parameters: str the entire string
			indexStr the characters to be included in the entire string
Returns: the first index, if not, returns -1

获取字符串中含有某个字符的第一个索引
参数：str 整个字符串
     indexStr 整个字符串中要包含的字符
返回：第一个索引，无则返回-1
 */
func GetStringIndex(str string,indexStr string) int {
	strings := []string{}
	//遍历字符串并将其存入string切片中
	for _,k :=range str{
		strings = append(strings, string(k))
	}
	index :=[]int{}
	//将每一个符合条件的字符索引存入int切片中
	for i:=0;i<len(strings);i++ {
		if indexStr == strings[i]{
			index = append(index,i)
		}
	}
	//返回正确索引值
	if len(index)>0{
		return index[0]
	}
	return -1
}
/**
Get the last index of a string containing a character
Parameters: str the entire string
			indexStr the characters to be included in the entire string
Returns: the first index, if not, returns -1

获取字符串中含有某个字符的最后一个索引
参数：str 整个字符串
     indexStr 整个字符串中要包含的字符
返回：第一个索引，无则返回-1
 */
func GetStringLastIndex(str string,indexStr string) int {
	strings := []string{}
	//遍历字符串并将其存入string切片中
	for _,k :=range str{
		strings = append(strings, string(k))
	}
	index :=[]int{}
	//将每一个符合条件的字符索引存入int切片中
	for i:=0;i<len(strings);i++ {
		if indexStr == strings[i]{
			index = append(index,i)
		}
	}
	//返回正确索引值
	if len(index)>0{
		ints := index[len(index)-1:]
		return ints[0]
	}
	return -1
}
/**
Get the first custom index containing a character in a string
Parameters: str The entire string
			indexStr The characters to be included in the entire string
			indexS The first position characters to be included in the string
Returns: indexS index, none returns -1

获取字符串中含有某个字符的第自定义个索引
参数：str 整个字符串
     indexstr 整个字符串中要包含的字符
	 indexs  字符串中要包含的第几个位置字符
返回：第indexs个索引，无则返回-1
 */
func GetStringIndexForNum(str string,indexStr string,indexS int) int {
	strings := []string{}
	//遍历字符串并将其存入string切片中
	for _,k :=range str{
		strings = append(strings, string(k))
	}
	index :=[]int{}
	//将每一个符合条件的字符索引存入int切片中
	for i:=0;i<len(strings);i++ {
		if indexStr == strings[i]{
			index = append(index,i)
		}
	}
	//返回正确索引值
	if len(index)>0 && len(index)>=indexS{
		ints := index[indexS-1:indexS]
		return ints[0]
	}
	return -1
}


/**
Simulate a ternary expression, returning itself if value is not "", otherwise returning the default value

模拟三元表达式，value不为""则返回本身，否则返回默认值
 */
func DefaultValue(value string,defaultBalue string)string  {
	if value==""{
		return defaultBalue
	}
	return value
}
