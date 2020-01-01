package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

var pool *redis.Pool

func initPool(address string,maxIdle,maxActive int, idleTimeout time.Duration)  {
	pool = &redis.Pool{
		Dial: func() (conn redis.Conn, e error) {
			return redis.Dial("tcp", address)
		},
		TestOnBorrow:    nil,
		MaxIdle:         maxIdle,
		MaxActive:       maxActive,
		IdleTimeout:     idleTimeout,
		Wait:            false,
		MaxConnLifetime: 0,
	}
}