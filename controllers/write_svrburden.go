package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"testvm/conf"
	"testvm/models/loadinfo"
	"testvm/models/redismanager"
)

type WriteLoadInfoController struct {
	beego.Controller
	bdinfo loadinfo.LoadInfo
}

func (wc *WriteLoadInfoController) Post() {
	response := ""
	msg := wc.Ctx.Input.RequestBody
	msg_str := string(msg)
	wc.bdinfo = make(loadinfo.LoadInfo)
	//解析body消息，转换为LoadInfo类型
	for {
		if !wc.bdinfo.Parse(msg_str) {
			response = "Fail: request body information illegal\n"
			break
		}
		rdsm := new(redismanager.RedisManager)
		rdskeyname := conf.GlobalConfig().RDS_Keyname
		host := conf.GlobalConfig().RDS_Host
		port := conf.GlobalConfig().RDS_Port
		ok, _ := rdsm.Connect(host, port)
		if !ok {
			response = "Fail: connect redis failed.\n"
			break
		}
		if ok, _ := rdsm.Add(wc.bdinfo, rdskeyname); !ok {
			response = "Fail: redis add loadinfo failed\n"
			break
		}
		response = "OK\n"
		rdsm.Close()
		break
	}
	fmt.Println(response)
	wc.Ctx.WriteString(response)
}
