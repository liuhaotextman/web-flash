package main

import (
	"fmt"
	"log"
	"net/http"
	"web-flash/cache"
)

//web框架
//import (
//	"fmt"
//	httpFlash "web-flash/http"
//	"web-flash/routing"
//)
//
//func main() {
//	route := routing.New()
//	v1 := route.Group("/v1")
//	v1.Use(func(context *httpFlash.Context) {
//		fmt.Println("log")
//	})
//	{
//		v1.GET("/hello", func(context *httpFlash.Context) {
//			context.Html(200, "hello route")
//		})
//	}
//
//
//	route.Run(":9090")
//}


var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	cache.NewGroup("scores", 2<<10, cache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		},
	))

	addr := ":9999"
	peers := cache.NewHTTPPool(addr)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}

