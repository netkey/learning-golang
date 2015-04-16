package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

func main() {

	// 从配置文件获取redis配置并连接

	host := "127.0.0.1:6379"
	rs, err := redis.Dial("tcp", host)

	defer rs.Close()

	if err != nil {
		fmt.Println(err)
		fmt.Println("redis connect error")
		return
	}

	//rs.Do("SELECT", db)

	key := "aaa"
	value := "bbb"
	// 操作redis时调用Do方法，第一个参数传入操作名称（字符串），然后根据不同操作传入key、value、数字等
	// 返回2个参数，第一个为操作标识，成功则为1，失败则为0；第二个为错误信息
	n, err := rs.Do("SETNX", key, value)
	// 若操作失败则返回
	if err != nil {
		fmt.Println(err)
		return
	}
	// 返回的n的类型是int64的，所以得将1或0转换成为int64类型的再比较
	if n == int64(1) {
		// 设置过期时间为24小时
		n, _ := rs.Do("EXPIRE", key, 24*3600)
		if n == int64(1) {
			fmt.Println("success")
		}
	} else if n == int64(0) {
		fmt.Println("the key has already existed")
	}

	newkey := "aaa"

	newvalue, err1 := redis.String(rs.Do("GET", newkey))
	if err1 != nil {
		fmt.Println("fail")
	}
	fmt.Println(newvalue)

	// 存json数据
	key1 := "aaa"
	imap := map[string]string{"key1": "111", "key2": "222"}
	// 将map转换成json数据
	value1, _ := json.Marshal(imap)
	// 存入redis
	n, err2 := rs.Do("SETNX", key1, value1)
	if err2 != nil {
		fmt.Println(err)
	}
	if n == int64(1) {
		fmt.Println("success")
	}

	// 取json数据
	// 先声明imap用来装数据
	var imap2 map[string]string
	key3 := "aaa"
	// json数据在go中是[]byte类型，所以此处用redis.Bytes转换
	value3, err3 := redis.Bytes(rs.Do("GET", key3))
	if err3 != nil {
		fmt.Println(err)
	}
	// 将json解析成map类型
	errShal := json.Unmarshal(value3, &imap2)
	if errShal != nil {
		fmt.Println(err)
	}
	fmt.Println(imap["key1"])
	fmt.Println(imap["key2"])

}
