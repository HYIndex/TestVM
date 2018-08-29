package redismanager

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "testvm/conf"
	. "testvm/models/loadinfo"
)

type RedisManager struct {
	connect redis.Conn
}

func (rm *RedisManager) Connect(host string, port uint) (bool, error) {
	var err error
	rm.connect, err = redis.Dial("tcp", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		fmt.Println("Connect to redis error! [Error]:", err)
		return false, err
	}
	return true, nil
	//连接成功
}

func (rm *RedisManager) Add(loadinfo LoadInfo, rdskeyname string) (bool, error) {
	for k, v := range loadinfo {
		str_v := v.String()
		_, err := rm.connect.Do("HSET", rdskeyname, k, str_v)
		if err != nil {
			//redis hset failed
			return false, err
		}
	}
	return true, nil
}

func (rm *RedisManager) GetAll(rdskeyname string) (map[string]string, error) {
	ret, err := redis.StringMap(rm.connect.Do("HGETALL", rdskeyname))
	if err != nil {
		//redis hset failed
		return nil, err
	} else {
		return ret, nil
	}
}

func (rm *RedisManager) Close() {
	rm.connect.Close()
}
