/*
 * 给客户端返回可用的服务器节点
 */

package controllers

import (
	"github.com/astaxie/beego"
	"testvm/conf"
	"testvm/models/loadinfo"
	"testvm/models/redismanager"
	"testvm/models/logging"
	"github.com/Sirupsen/logrus"
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
		logging.GetLogger().WithFields(logrus.Fields{
			"package" : "controllers",
			"file" : "retrieve_servers.go",
		}).Infoln("redis connect fail!")
		return
	}
	defer rdsm.Close()
	ret, err := rdsm.GetAll(rdskeyname)
	if err != nil {
		rsc.Ctx.WriteString("Fail: redis getall fail!\n")
		logging.GetLogger().WithFields(logrus.Fields{
			"package" : "controllers",
			"file" : "retrieve_servers.go",
		}).Infoln("redis getall fail!")
		return
	}
	totalStreamInfo := new(loadinfo.StreamInfo)
	for _, v := range ret {
		tmpStream := make(loadinfo.StreamsAmt)
		tmpStream = tmpStream.FromString(v)
		if tmpStream == nil {
			rsc.Ctx.WriteString("Fail: StreamAmt convert from string fail!\n")
			logging.GetLogger().WithFields(logrus.Fields{
				"package" : "controllers",
				"file" : "retrieve_servers.go",
			}).Infoln("StreamAmt convert from string fail!")
			return
		}
		tmp := new(loadinfo.StreamInfo)
		tmp.Create(tmpStream)
		totalStreamInfo.Add(tmp)
	}
	totalSubStream := totalStreamInfo.GetTotalSub()
	if totalSubStream > limit {
		rsc.SetStatus(0)
		rsc.Ctx.WriteString("no usalbe lbaddress: substreaming more than limit\n")
		logging.GetLogger().WithFields(logrus.Fields{
			"package" : "controllers",
			"file" : "retrieve_servers.go",
		}).Infoln("no usalbe lbaddress: substreaming more than limit!")
	} else {
		rsc.SetStatus(1)
		rsc.Ctx.WriteString(lbaddr)
		logging.GetLogger().WithFields(logrus.Fields{
			"package" : "controllers",
			"file" : "retrieve_servers.go",
			"lbaddr" : lbaddr,
		}).Infoln("return lbaddress success!")
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
