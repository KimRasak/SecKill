package test

import ("github.com/gavv/httpexpect"
	"net/http"
	"net/http/httptest"
)
func main()  {
	// create http.Handler
	handler := FruitsHandler()

	// run server using httptest
	server := httptest.NewServer(handler)
	defer server.Close()

	// create httpexpect instance
	e := httpexpect.New(t, server.URL)

	// is it working?
	e.GET("/fruits").
		Expect().
		Status(http.StatusOK).JSON().Array().Empty()
}