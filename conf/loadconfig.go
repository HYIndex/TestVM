package conf

import (
	"github.com/astaxie/beego"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"github.com/Sirupsen/logrus"
)

type Config struct {
	RS_Slbaddr  string // LB地址
	RS_Limit    uint   // 限制人数
	RDS_Host    string // Redis服务主机ip
	RDS_Port    uint   // Redis服务进程端口
	RDS_Keyname string // Redis 保存负载信息的键名
	LOG_Level	string // 日志等级
}

var (
	globalConfig Config
	configLock   = new(sync.RWMutex)
	log = logrus.New()
)

func init() {
	//加载配置
	log.Out = os.Stdout
	log.Level = logrus.DebugLevel
	log.Formatter = new(logrus.JSONFormatter)

	if !loadConfig() {
		log.WithFields(logrus.Fields{
			"package" : "conf",
			"function" : "init",
		}).Fatalln("Loag config file failed!")
		os.Exit(1)
	}
	reload()
}

// 给其他包提供全局配置的接口
// 热加载会在更外的go线程中使用globalConfig，需加锁保护
func GlobalConfig() Config {
	configLock.RLock()
	defer configLock.RUnlock()
	return globalConfig
}

func loadConfig() bool {
	err := beego.LoadAppConfig("ini", "/usr/local/testvm/conf/config.ini")
	if err != nil {
		return false
	}
	var temp Config
	temp.RS_Slbaddr = beego.AppConfig.String("retrievesvr::slbaddr")
	tmplimit, _ := beego.AppConfig.Int("retrievesvr::limit")
	temp.RS_Limit = uint(tmplimit)
	temp.RDS_Host = beego.AppConfig.String("redis::host")
	tmpport, _ := beego.AppConfig.Int("redis::port")
	temp.RDS_Port = uint(tmpport)
	temp.RDS_Keyname = beego.AppConfig.String("redis::keyname")
	temp.LOG_Level = beego.AppConfig.String("log::level")
	configLock.Lock()
	globalConfig = temp
	configLock.Unlock()
	return true
}

//配置文件热加载
func reload() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGUSR1)
	go func() {
		for {
			<-s
			if !loadConfig() {
				log.WithFields(logrus.Fields{
					"package" : "conf",
					"function" : "reload",
				}).Fatalln("Loag config file failed!")
				os.Exit(1)
			}
			log.WithFields(logrus.Fields{
				"package" : "conf",
			}).Infoln("Config file reloaded!")
		}
	}()
}
