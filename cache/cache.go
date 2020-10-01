package cache

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func init() {
	_, err := os.Stat("cache.json")
	if os.IsNotExist(err) {
		_, _ = os.Create("cache.json")
	}
}

type JsonCache struct {
	Caches []cacheData `json:"caches"`
}

type cacheData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

//retrieveCacheJson using keys, or not exists, create a key pair
func RetrieveCache(key string) ([]byte, bool) {
	fBody, err := ioutil.ReadFile("cache.json")
	if err != nil {
		log.Println(err)
	}

	var cache JsonCache
	_ = json.Unmarshal(fBody, &cache)

	for _, item := range cache.Caches {
		if item.Key == key {
			data, _ := base64.StdEncoding.DecodeString(item.Value)
			return data, true
		}
	}

	return nil, false
}

func Put(key string, value []byte) bool {
	fBody, err := ioutil.ReadFile("cache.json")
	if err != nil {
		log.Println(err)
	}

	var cache JsonCache

	_ = json.Unmarshal(fBody, &cache)

	cache.Caches = append(cache.Caches, cacheData{Key: key, Value: base64.StdEncoding.EncodeToString(value)})

	js, err := json.Marshal(cache)
	if err != nil {
		return false
	}

	err = ioutil.WriteFile("cache.json", js, os.ModePerm)
	if err != nil {
		return false
	}

	return true
}