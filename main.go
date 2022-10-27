package main

import (
	"fmt"
	"web-flash/orm"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	engine, _ := orm.NewEngine("sqlite3", "gee.db")
	defer engine.Close()
	s := engine.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	count, _ := result.RowsAffected()
	fmt.Printf("Exec success, %d affected\n", count)
}

//import (
//	"fmt"
//	"log"
//	"net/http"
//	"web-flash/cache"
//)
//
//var db = map[string]string{
//	"Tom":  "630",
//	"Jack": "589",
//	"Sam":  "567",
//}
//
//func main() {
//	cache.NewGroup("scores", 2<<10, cache.GetterFunc(
//		func(key string) ([]byte, error) {
//			log.Println("[SlowDB] search key", key)
//			if v, ok := db[key]; ok {
//				return []byte(v), nil
//			}
//			return nil, fmt.Errorf("%s not exist", key)
//		},
//	))
//
//	addr := ":9999"
//	peers := cache.NewHTTPPool(addr)
//	log.Println("geecache is running at", addr)
//	log.Fatal(http.ListenAndServe(addr, peers))
//}

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
