/*
 * Redis管理器，提供连接、关闭、写入和读取接口
 */

package redismanager

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "testvm/conf"
	. "testvm/models/loadinfo"
	"testvm/models/logging"
	"github.com/Sirupsen/logrus"
)

type RedisManager struct {
	connect redis.Conn
}

func (rm *RedisManager) Connect(host string, port uint) (bool, error) {
	var err error
	rm.connect, err = redis.Dial("tcp", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		logging.GetLogger().WithFields(logrus.Fields{
			"package" : "redismanager",
			"file" : "redismanager.go",
			"fail with error" : err,
		}).Errorln("redis connect fail!")
		return false, err
	}
	return true, nil
}

func (rm *RedisManager) Add(loadinfo LoadInfo, rdskeyname string) (bool, error) {
	for k, v := range loadinfo {
		str_v := v.String()
		_, err := rm.connect.Do("HSET", rdskeyname, k, str_v)
		if err != nil {
			logging.GetLogger().WithFields(logrus.Fields{
				"package" : "redismanager",
				"file" : "redismanager.go",
				"fail with error" : err,
			}).Errorln("redis hset fail!")
			return false, err
		}
	}
	return true, nil
}

func (rm *RedisManager) GetAll(rdskeyname string) (map[string]string, error) {
	ret, err := redis.StringMap(rm.connect.Do("HGETALL", rdskeyname))
	if err != nil {
		logging.GetLogger().WithFields(logrus.Fields{
			"package" : "redismanager",
			"file" : "redismanager.go",
			"fail with error" : err,
		}).Errorln("redis hgetall fail!")
		return nil, err
	} else {
		return ret, nil
	}
}

func (rm *RedisManager) Close() {
	rm.connect.Close()
}
