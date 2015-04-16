/*
  这是一个 go goroutine 是否异常退出的侦测示例程序

  goroutine 定义了一个匿名函数，在异常退出时该匿名函数将向一个 mess channel 发送一个信息

  tsingson 2015/03/11
*/
package main

import (
	"log"
	"math/rand"
	"runtime"
	"time"
)

type message struct {
	normal bool                   //true means exit normal, otherwise
	state  map[string]interface{} //goroutine state
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	mess := make(chan message, 100)
	for i := 0; i < 100000; i++ {
		go worker(mess)
	}
	supervisor(mess)
}

func worker(mess chan<- message) {

	// 异常退出时调用 defer 定义的匿名函数，该函数向一个 mess channel 发送一个异常信息
	defer func() {
		exit_message := message{state: make(map[string]interface{})}
		i := recover()
		if i != nil {
			exit_message.normal = false
		} else {
			exit_message.normal = true
		}
		mess <- exit_message
	}()
	now := time.Now()
	seed := now.UnixNano()
	rand.Seed(seed)
	num := rand.Int63()

	if num%2 != 0 {
		panic("not evening")
	} else {
		runtime.Goexit()
	}
}

// mess channel 用于侦测 worker 是否异常那家出
func supervisor(mess <-chan message) {
	for i := 0; i < 100; i++ {
		m := <-mess
		switch m.normal {
		case true:
			log.Println(time.Now, "exit normal, nothing serious!")
		case false:
			log.Println(time.Now, "exit abnormal, something went wrong")
		}

	}
}
