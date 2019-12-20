package testing

import (
	"SecKill/data"
	"SecKill/engine"
	"github.com/gavv/httpexpect"
	"net/http/httptest"
	"testing"
)

var E *httpexpect.Expect

func TestNormal(t *testing.T)  {
	// 启动服务器
	server := httptest.NewServer(engine.SeckillEngine())
	E = httpexpect.New(t, server.URL)
	defer data.Close()


}