package main

import (
	"github.com/astaxie/beego"
	_ "testvm/conf"
	_ "testvm/routers"
)

func main() {
	beego.Run()
}
