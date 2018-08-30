package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"runtime"
	"path/filepath"
	_ "testvm/routers"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
		"bytes"
)

const (
	TEST_URL1 = "/testvm/retrievesvrs"
	TEST_BODY = ""
	TEST_URL2 = "/testvm/loadinfo/push"
	TEST_BODY2 = "APS=10&APC=10&ASS=20&ASC=20&VPS=10&VPC=10&VSS=20&VSC=20&SRC=192.168.0.0:1238"
	TEST_URL3 = "/testvm/loadinfo/get"
	)

func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, "../.." + string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

// TestBeego is a sample to run an endpoint test
func TestRetrievesvrs(t *testing.T) {
	r, _ := http.NewRequest("POST", TEST_URL2, bytes.NewReader([]byte(TEST_BODY2)))
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("testing", "TestBeego", "Status Code:", w.Code, "ResponseMessage:", w.Body.String())

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
	})
}

// 压力测试
func BenchmarkRetrievesvrs(t *testing.B) {
	for i := 0; i < t.N; i++ {
		r, _ := http.NewRequest("POST", TEST_URL2, bytes.NewReader([]byte(TEST_BODY2)))
		w := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		beego.Trace("testing", "TestBeego", "Status Code:", w.Code, "ResponseMessage:", w.Body.String())
	}
}
