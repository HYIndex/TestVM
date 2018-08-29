package controllers

import (
	"github.com/astaxie/beego"
	"testvm/conf"
	"testvm/models/loadinfo"
	"testvm/models/redismanager"
)

type RetrieveServersController struct {
	beego.Controller
	status int
}

func (rsc *RetrieveServersController) Get() {
	//从Redis获取负载信息与配置文件设置的门限值比较
	host := conf.GlobalConfig().RDS_Host
	port := conf.GlobalConfig().RDS_Port
	rdskeyname := conf.GlobalConfig().RDS_Keyname
	limit := conf.GlobalConfig().RS_Limit
	lbaddr := conf.GlobalConfig().RS_Slbaddr

	rdsm := new(redismanager.RedisManager)
	if ok, _ := rdsm.Connect(host, port); !ok {
		rsc.Ctx.WriteString("Fail: redis connect fail!\n")
		return
	}
	defer rdsm.Close()
	ret, err := rdsm.GetAll(rdskeyname)
	if err != nil {
		rsc.Ctx.WriteString("Fail: redis getall fail!\n")
		return
	}
	totalStreamInfo := new(loadinfo.StreamInfo)
	for _, v := range ret {
		tmpStream := make(loadinfo.StreamsAmt)
		tmpStream = tmpStream.FromString(v)
		if tmpStream == nil {
			rsc.Ctx.WriteString("Fail: convert from string fail!\n")
			return
		}
		tmp := new(loadinfo.StreamInfo)
		tmp.Create(tmpStream)
		totalStreamInfo.Add(tmp)
	}
	totalSubStream := totalStreamInfo.GetTotalSub()
	if totalSubStream > limit {
		rsc.Ctx.WriteString("no usalbe lbaddress: substreaming more than limit\n")
	} else {
		rsc.Ctx.WriteString(lbaddr + "\n")
	}
}

func (rsc *RetrieveServersController) Post() {
	rsc.Get()
}

func (rsc *RetrieveServersController) Status() int {
	return rsc.status
}

func (rsc *RetrieveServersController) SetStatus(s int) {
	rsc.status = s
}