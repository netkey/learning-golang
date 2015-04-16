package main

import "fmt"
import "encoding/json"

func main() {
	var s S
	s.a = 5
	s.b[0] = 3.123
	s.b[1] = 111.11
	s.b[2] = 1234.123
	s.c = "hello"
	s.d[0] = 0x55

	j, _ := json.Marshal(s)
	fmt.Println(string(j))
}

type S struct {
	a int
	b [4]float32
	c string
	d [12]byte
}

func (this S) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"a": this.a,
		"b": this.b,
		"c": this.c,
		"d": this.d,
	})
}
