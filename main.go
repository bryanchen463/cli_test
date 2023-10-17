package main

import (
	fasthttpcli "bryanchen463/cli_test/http_cli/fast_http"
	"bryanchen463/cli_test/http_cli/gentleman"
	gohttpclient "bryanchen463/cli_test/http_cli/go-http-client"
	goretryablenhttpcli "bryanchen463/cli_test/http_cli/go-retryablenhttp"
	grequestscli "bryanchen463/cli_test/http_cli/grequests"
	heimdallcli "bryanchen463/cli_test/http_cli/heimdall"
	"bryanchen463/cli_test/http_cli/http"
	"bryanchen463/cli_test/utils"
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}

func main() {
	testRequest()
	fmt.Println(Result())
	return
	benchFn(func() error {
		return http.GET(utils.Url)
	}, 10000, "http_10000")
	benchFn(func() error {
		return http.POST(utils.EchoUrl, "hello")
	}, 10000, "http_post_10000")
	benchFn(func() error {
		return gentleman.Get(utils.Url)
	}, 10000, "gentleman_10000")
	benchFn(func() error {
		return gentleman.Post(utils.EchoUrl, "hello")
	}, 10000, "gentleman_post_10000")
	benchFn(func() error {
		return fasthttpcli.Get(utils.Url)
	}, 10000, "fasthttp_10000")
	benchFn(func() error {
		return fasthttpcli.Post(utils.EchoUrl, "hello")
	}, 10000, "fasthttp_post_10000")
	benchFn(func() error {
		return gohttpclient.Get(utils.Url)
	}, 10000, "gohttpclient_10000")
	benchFn(func() error {
		return gohttpclient.Post(utils.EchoUrl, "hello")
	}, 10000, "gohttpclient_post_10000")

	benchFn(func() error {
		return goretryablenhttpcli.Get(utils.Url)
	}, 10000, "goretryablenhttpcli_10000")
	benchFn(func() error {
		return goretryablenhttpcli.POST(utils.EchoUrl, "hello")
	}, 10000, "goretryablenhttpcli_post_10000")
	fmt.Println(Result())
}

type fn func(url string) error
type postFn func(url string, payload string) error

func wrapperPost(f func(url string, payload string) error) fn {

	return func(url string) error {
		return f(url, "hello")
	}
}

func testRequest() {
	getFuncs := []fn{http.GET, fasthttpcli.Get, gohttpclient.Get, goretryablenhttpcli.Get, grequestscli.Get, heimdallcli.Get}
	postFuncs := []postFn{http.POST, fasthttpcli.Post, gohttpclient.Post, goretryablenhttpcli.POST, grequestscli.Post, heimdallcli.Post}
	times := []int{5000, 10000, 100000}
	// times := []int{1000}
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
		fmt.Printf("-------------------------GET %d--------------------------\n", t)
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
			fmt.Printf("-------------------------POST %d %d--------------------------\n", t, len(payload))
			fmt.Println(Result())
		}
	}
}
