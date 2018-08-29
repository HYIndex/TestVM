package conf

import (
	"fmt"
	"github.com/astaxie/beego"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Config struct {
	RS_Slbaddr  string //LB地址
	RS_Limit    uint   //限制人数
	RDS_Host    string
	RDS_Port    uint
	RDS_Keyname string
}

var (
	globalConfig Config
	configLock   = new(sync.RWMutex)
)

func init() {
	//加载配置
	if !loadConfig() {
		//
		os.Exit(1)
	}
	reload()
}

func GlobalConfig() Config {
	configLock.RLock()
	defer configLock.RUnlock()
	return globalConfig
}

func loadConfig() bool {
	err := beego.LoadAppConfig("ini", "conf/config.ini")
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
			fmt.Println("reload config!")
			if !loadConfig() {
				//
				os.Exit(1)
			}
		}
	}()
}
