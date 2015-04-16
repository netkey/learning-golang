package main

import (
	"fmt"
	"github.com/wangtuanjie/ip17mon"
)

func checkError(err error) {
	if err != nil {
		fmt.Println("err:", err)
	}
}

func init() {

	ipDataFile := "./17monipdb.dat"
	err := ip17mon.Init(ipDataFile)

	checkError(err)
}

func main() {
	loc, err := ip17mon.Find("116.228.111.18")
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println(loc)
}
