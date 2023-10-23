package main

import (
	fasthttpcli "bryanchen463/cli_test/http_cli/fast_http"
	gohttpclient "bryanchen463/cli_test/http_cli/go-http-client"
	gorestycli "bryanchen463/cli_test/http_cli/go-resty"
	goretryablenhttpcli "bryanchen463/cli_test/http_cli/go-retryablenhttp"
	grequestscli "bryanchen463/cli_test/http_cli/grequests"
	heimdallcli "bryanchen463/cli_test/http_cli/heimdall"
	"bryanchen463/cli_test/http_cli/http"
	pestercli "bryanchen463/cli_test/http_cli/pester"
	reqcli "bryanchen463/cli_test/http_cli/req"
	"bryanchen463/cli_test/utils"
	fasthttpwscli "bryanchen463/cli_test/websocket_cli/fasthttp_ws"
	fastwscli "bryanchen463/cli_test/websocket_cli/fastws"
	gobwascli "bryanchen463/cli_test/websocket_cli/gobwas"
	gorillawebsocketclient "bryanchen463/cli_test/websocket_cli/gorilla"
	nhooyrcli "bryanchen463/cli_test/websocket_cli/nhooyr"
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var rd = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateRandomString(length int) string {

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rd.Intn(len(charset))]
	}

	return string(b)
}

func main() {
	testRequest()
	// fmt.Println(Result())
	testWsEcho()
	// time.Sleep(time.Second * 10)
	testWsReceive()
}

type fn func(url string) error
type postFn func(url string, payload string) error

func wrapperPost(f func(url string, payload string) error) fn {

	return func(url string) error {
		return f(url, "hello")
	}
}

func testRequest() {
	getFuncs := []fn{http.GET, fasthttpcli.Get, gohttpclient.Get, goretryablenhttpcli.Get, grequestscli.Get, heimdallcli.Get, pestercli.Get, reqcli.Get, gorestycli.Get}
	postFuncs := []postFn{http.POST, fasthttpcli.Post, gohttpclient.Post, goretryablenhttpcli.POST, grequestscli.Post, heimdallcli.Post, pestercli.Post, reqcli.Post, gorestycli.Post}
	times := []int{5000, 10000, 100000}
	getFuncs = []fn{http.GET, fasthttpcli.Get, gohttpclient.Get, goretryablenhttpcli.Get, grequestscli.Get, pestercli.Get, reqcli.Get, gorestycli.Get}
	postFuncs = []postFn{http.POST, fasthttpcli.Post, gohttpclient.Post, goretryablenhttpcli.POST, grequestscli.Post, pestercli.Post, reqcli.Post, gorestycli.Post}
	// times := []int{100}
	paylaods := []string{}
	paylaodsLen := []int{100, 1024}
	for _, l := range paylaodsLen {
		randomString := generateRandomString(l)
		paylaods = append(paylaods, randomString)
	}
	for _, t := range times {
		for _, fn := range getFuncs {
			funcValue := reflect.ValueOf(fn)
			funcName := runtime.FuncForPC(funcValue.Pointer()).Name()
			name := fmt.Sprintf("%s_%d", funcName, t)
			benchFn(func() error {
				return fn(utils.Url)
			}, t, name)
		}
		fmt.Printf("#### GET%d次\n", t)
		fmt.Println(Result())
		for _, payload := range paylaods {
			for _, fn := range postFuncs {
				funcValue := reflect.ValueOf(fn)
				funcName := runtime.FuncForPC(funcValue.Pointer()).Name()
				name := fmt.Sprintf("%s_%d_%d", funcName, len(payload), t)
				benchFn(func() error {
					return fn(utils.EchoUrl, payload)
				}, t, name)
			}
			fmt.Printf("### POST %d次%d字节\n", t, len(payload))
			fmt.Println(Result())
		}
	}
}

type wsFn func(string) (int64, error)
type wsInit func(string) error
type wsClean func()

func testWsEcho() {
	initFuncs := []wsInit{gorillawebsocketclient.Init, fastwscli.Init, gobwascli.Init, nhooyrcli.Init, fasthttpwscli.Init}
	wsFuncs := []wsFn{gorillawebsocketclient.SendReacv, fastwscli.SendReacv, gobwascli.SendReacv, nhooyrcli.SendReacv, fasthttpwscli.SendReacv}
	cleanFuncs := []wsClean{gorillawebsocketclient.Clean, fastwscli.Clean, gobwascli.Clean, nhooyrcli.Clean, fasthttpwscli.Clean}
	// randomString := generateRandomString(5000)
	times := []int{5000, 10000, 100000}
	paylaods := []string{}
	paylaodsLen := []int{100, 1024}
	for _, l := range paylaodsLen {
		randomString := generateRandomString(l)
		paylaods = append(paylaods, randomString)
	}
	for _, t := range times {
		for _, payload := range paylaods {
			for i, fn := range wsFuncs {
				funcValue := reflect.ValueOf(fn)
				funcName := runtime.FuncForPC(funcValue.Pointer()).Name()
				initFuncs[i]("wss://127.0.0.1:8080/ws/echo")
				benchFn(func() error {
					_, err := fn(payload)
					return err
				}, t, funcName)
				cleanFuncs[i]()
			}
			fmt.Printf("#### echo%d字节数据%d次\n", len(payload), t)
			fmt.Println(Result())
		}
	}
}

type wsReciveFn func() (int64, error)

func testWsReceive() {
	initFuncs := []wsInit{gorillawebsocketclient.Init, fastwscli.Init, gobwascli.Init, nhooyrcli.Init, fasthttpwscli.Init}
	cleanFuncs := []wsClean{gorillawebsocketclient.Clean, fastwscli.Clean, gobwascli.Clean, nhooyrcli.Clean, fasthttpwscli.Clean}
	wsFuncs := []wsReciveFn{gorillawebsocketclient.Receive, fastwscli.Receive, gobwascli.Receive, nhooyrcli.Receive, fasthttpwscli.Receive}
	times := []int{5000, 10000, 100000}
	// paylaodsLen := []int{100, 1024}
	counts := []int{100, 1024}
	for _, t := range times {
		for _, n := range counts {
			for i, fn := range wsFuncs {
				funcValue := reflect.ValueOf(fn)
				funcName := runtime.FuncForPC(funcValue.Pointer()).Name()
				// name := fmt.Sprintf("%s_%d_%d", funcName, n, t)
				addr := fmt.Sprintf("%s/message/%d/%d", utils.WsUrl, n, t)
				err := initFuncs[i](addr)
				if err != nil {
					fmt.Println(err)
					return
				}
				benchFn(func() error {
					_, err := fn()
					return err
				}, t, funcName)
				cleanFuncs[i]()
			}
			fmt.Printf("#### 接收%d字节数据%d次\n", n, t)
			fmt.Println(Result())
		}
	}
}
