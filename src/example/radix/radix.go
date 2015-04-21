// example program

package main

import (
	"fmt"
	"github.com/fzzy/radix/extra/pool"
	"github.com/fzzy/radix/extra/pubsub"
	"github.com/fzzy/radix/redis"
	log "github.com/mgutz/logxi/v1"
	"os"
	"runtime"
	"time"

	// debug
	"flag"
	//"net/http"
	_ "net/http/pprof"
	"runtime/pprof"
)

func errHndlr(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			//log.Debug(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	/**
	c, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
	errHndlr(err)
	defer c.Close()

	*/
	loger := log.New("radis")
	loger.SetLevel(6)

	loger.Info("log XI ")

	runtime.GOMAXPROCS(4)
	pool, err := pool.NewPool("tcp", "localhost:6379", 10)
	errHndlr(err)

	conns := make([]*redis.Client, 5)
	for i := range conns {
		if conns[i], err = pool.Get(); err != nil {
			fmt.Println("error")
		}
	}

	for i := range conns {
		go radixClient(conns[i], i)
	}

	/*
		for i := range conns {
			pool.Put(conns[i])
		}
	*/
	defer pool.Empty()

	for {
	}

	// select database

}

func radixClient(c *redis.Client, i int) error {

	runTimeNow := time.Now()

	fmt.Println("########################", i, runTimeNow, "####################\n")

	loger := log.New("radis")
	loger.SetLevel(6)
	loger.Info("radis client", "number", i, "time", runTimeNow)

	r := c.Cmd("select", 8)
	errHndlr(r.Err)

	r = c.Cmd("flushdb")
	errHndlr(r.Err)

	r = c.Cmd("echo", "Hello world!")
	errHndlr(r.Err)

	s, err := r.Str()
	errHndlr(err)
	fmt.Println("echo:", s)

	//* Strings
	r = c.Cmd("set", "mykey0", "myval0")
	errHndlr(r.Err)

	s, err = c.Cmd("get", "mykey0").Str()
	errHndlr(err)
	fmt.Println("mykey0:", s)

	myhash := map[string]string{
		"mykey1": "myval1",
		"mykey2": "myval2",
		"mykey3": "myval3",
	}

	// Alternatively:
	// c.Cmd("mset", "mykey1", "myval1", "mykey2", "myval2", "mykey3", "myval3")
	r = c.Cmd("mset", myhash)
	errHndlr(r.Err)

	ls, err := c.Cmd("mget", "mykey1", "mykey2", "mykey3").List()
	errHndlr(err)
	fmt.Println("mykeys values:", ls)

	//* List handling
	mylist := []string{"foo", "bar", "qux"}

	// Alternativaly:
	// c.Cmd("rpush", "mylist", "foo", "bar", "qux")
	r = c.Cmd("rpush", "mylist", mylist)
	errHndlr(r.Err)

	mylist, err = c.Cmd("lrange", "mylist", 0, -1).List()
	errHndlr(err)

	fmt.Println("mylist:", mylist)

	//* Hash handling

	// Alternatively:
	// c.Cmd("hmset", "myhash", ""mykey1", "myval1", "mykey2", "myval2", "mykey3", "myval3")
	r = c.Cmd("hmset", "myhash", myhash)
	errHndlr(r.Err)

	myhash, err = c.Cmd("hgetall", "myhash").Hash()
	errHndlr(err)

	fmt.Println("myhash:", myhash)

	//* Pipelining
	c.Append("set", "multikey", "multival")
	c.Append("get", "multikey")

	c.GetReply()     // set
	r = c.GetReply() // get
	errHndlr(r.Err)

	s, err = r.Str()
	errHndlr(err)
	fmt.Println("multikey:", s)

	//* Publish/Subscribe

	// Subscribe
	c2, err := redis.Dial("tcp", "localhost:6379")
	errHndlr(err)
	defer c2.Close()
	psc := pubsub.NewSubClient(c2)
	psr := psc.Subscribe("queue1", "queue2")

	// Publish
	c.Cmd("publish", "queue1", "ohai")

	// Receive publish
	psr = psc.Receive() //Blocks until reply is received or timeout is tripped
	if !psr.Timeout() {
		fmt.Println("publish:", psr.Message)
	} else {
		fmt.Println("error: sub timedout")
		return nil
	}

	// Unsubscribe
	psc.Unsubscribe("queue1", "queue2") //Unsubscribe before issuing any other commands with c

	c.Close()

	return nil

}
