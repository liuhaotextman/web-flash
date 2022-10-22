package cache

import (
	"fmt"
	"log"
	"testing"
)

var db = map[string]string{
	"tom":  "630",
	"jack": "589",
	"sam":  "567",
}

func TestGetGroup(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	gee := NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[slow db] search key", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; ok {
					loadCounts[key] = 0
				}
				loadCounts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		},
	))

	for k, v := range db {
		if view, err := gee.Get(k); err != nil || view.String() != v {
			t.Fatal("failed to get value of tom")
		}
		if _, err := gee.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}

		if view, err := gee.Get("unknown"); err == nil {
			t.Fatalf("the value of unknow should be empty, but %s got", view)
		}
	}
}
