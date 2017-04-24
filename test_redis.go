package main

import (
	// "encoding/json"
	"github.com/garyburd/redigo/redis"
)

func main() {
	rs, _ := redis.Dial("tcp", "localhost:6379")
	rs.Do("SELECT", 1)
	defer rs.Close()
	key := "token123"
	value := "lalal"
	rs.Do("SET", key, value)
}
