/*
 * 获取当前IDC的webrtc的负载信息
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

type OfferLoadInfoController struct {
	beego.Controller
	bdinfo loadinfo.LoadInfo
}

func (obic *OfferLoadInfoController) Get() {
	host := conf.GlobalConfig().RDS_Host
	port := conf.GlobalConfig().RDS_Port
	rdskeyname := conf.GlobalConfig().RDS_Keyname

	rdsm := new(redismanager.RedisManager)
	if ok, _ := rdsm.Connect(host, port); !ok {
		obic.Ctx.WriteString("Fail: redis connect fail!\n")
		logging.GetLogger().WithFields(logrus.Fields{
			"package" : "controllers",
			"file" : "offer_svrloadinfo.go",
		}).Infoln("redis connect fail!")
		return
	}
	defer rdsm.Close()
	ret, err := rdsm.GetAll(rdskeyname)
	if err != nil {
		obic.Ctx.WriteString("Fail: redis getall fail!\n")
		logging.GetLogger().WithFields(logrus.Fields{
			"package" : "controllers",
			"file" : "offer_svrloadinfo.go",
		}).Infoln("redis getall fail!")
		return
	}
	respInfo := new(loadinfo.ResponseInfo)
	for k, v := range ret {
		tmpStream := make(loadinfo.StreamsAmt)
		tmpStream = tmpStream.FromString(v)
		if tmpStream == nil {
			logging.GetLogger().WithFields(logrus.Fields{
				"package" : "controllers",
				"file" : "offer_svrloadinfo.go",
			}).Infoln("StreamAmt convert from string failed!")
			return
		}
		tmp := new(loadinfo.StreamInfo)
		tmp.Create(tmpStream)
		respInfo.Total.Add(tmp)
		detailinfo := new(loadinfo.DetailInfo)
		detailinfo.Streaminfo = *tmp
		detailinfo.Ipport = k
		respInfo.Detail = append(respInfo.Detail, *detailinfo)
	}
	obic.Data["json"] = respInfo
	obic.ServeJSON()
	logging.GetLogger().WithFields(logrus.Fields{
		"package" : "controllers",
		"file" : "offer_svrloadinfo.go",
	}).Infoln("return loadinfo with join response success!")
}

func (obic *OfferLoadInfoController) Post() {
	obic.Get()
}