package main

import (
	"github.com/astaxie/beego"
	_ "testvm/routers"
	_ "testvm/conf"
)

func main() {
	beego.Run()
}
