package routers

import (
	"github.com/astaxie/beego"
	"testvm/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/testvm/retrievesvrs", &controllers.RetrieveServersController{})
	beego.Router("/testvm/loadinfo/push", &controllers.WriteLoadInfoController{})
	beego.Router("/testvm/loadinfo/get", &controllers.OfferLoadInfoController{})
}
