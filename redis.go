package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	redigo "github.com/garyburd/redigo/redis"
)

var (
	redisAddress string = "127.0.0.1:6379"
	redisConn    *redigo.Pool
)

const prizePoolKey = "prizepool"

type KuponOVO struct {
	Name       string `json:"name"`
	Value      int64  `json:"value"`
	Percentage int64  `json:"percentage"`
}

func initRedis() {
	redisConn = &redigo.Pool{
		MaxIdle:     100,
		MaxActive:   100,
		IdleTimeout: 300 * time.Second,
		Wait:        true,
		Dial: func() (redigo.Conn, error) {
			return redigo.Dial("tcp", redisAddress)
		},
	}
	// do PING, if failed then Fatal
	if err := pingRedis(); err != nil {
		log.Fatal(err)
	}
}
func pingRedis() error {
	c := redisConn.Get()
	defer c.Close()
	if _, err := c.Do("PING"); err != nil {
		return err
	}
	return nil
}

// sample command to do a command to redis
func sampleDoCommand(command string) {
	c := redisConn.Get()
	defer c.Close()
	reply, err := c.Do(command)
	fmt.Println(reply, err)
}
func GetConn() (redigo.Conn, error) {
	conn, err := redisConn.Dial()
	if err != nil {
		return nil, err
	}
	return conn, err
}
func redisAddKey(param KuponOVO) error {
	conn, err := GetConn()
	if err != nil {
		return err
	}
	keyMember := normalizeParamName(param.Name)
	conn.Send("MULTI")
	conn.Send("ZADD", prizePoolKey, param.Percentage/100, keyMember)
	conn.Send("HMSET", keyMember, getRedisHMSETKeys(param))
	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}
	return nil
}
func normalizeParamName(key string) string {
	return strings.Replace(key, " ", "_", -1)
}
func getRedisHMSETKeys(param KuponOVO) []string {
	return []string{
		"name",
		param.Name,
		"value",
		strconv.FormatInt(param.Value, 10),
		"percentage",
		strconv.FormatInt(param.Percentage, 10),
	}
}
func redisGetRand() (KuponOVO, error) {
	conn, err := GetConn()
	if err != nil {
		return KuponOVO{}, err
	}

	randNum := rand.Intn(100)
	key, err := conn.Do("ZRANGEBYSCORE", prizePoolKey, randNum, "+inf LIMIT 0 1")
	if err != nil {
		return KuponOVO{}, err
	}

	result, err := conn.Do("HMGET", key)
	kuponOVO := result.(KuponOVO)
	return kuponOVO, nil
}
