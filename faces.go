//before running create "images" dirrectory
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	URL     = "http://uifaces.com/api/v1/random"
	DIR     = "images"
	DEBUG   = false
	WORKERS = 2

	mx sync.Mutex
)

func main() {
	DEBUG = true
	for i := 0; i < WORKERS; i++ {
		go work()
	}
	for {
	}
}

func work() {
	user := Random()
	mx.Lock()
	if !checkFile(user) {
		save(user)
	}
	mx.Unlock()
}

func get(url string) []byte {
	r, err := http.Get(url)
	if err != nil {
		return []byte{}
	}
	defer r.Body.Close()
	b, _ := ioutil.ReadAll(r.Body)
	return b
}

func checkFile(u Response) bool {
	if _, err := os.Stat(u.ImagePath()); os.IsNotExist(err) {
		DBG("File %s not exist", u.ImagePath())
		return false
	}
	return true
}

func save(u Response) {
	DBG("Start get image")
	im := get(u.ImageUrls.Epic)
	DBG("End get image")
	if len(im) != 0 {
		err := ioutil.WriteFile(u.ImagePath(), im, 0777)
		if err != nil {
			panic(err)
		}
	}
}

func Random() Response {
	DBG("Start fetch random user")
	bts := get(URL)
	DBG("End fetch random user")
	var res Response
	json.Unmarshal(bts, &res)
	return res
}

type Response struct {
	Username  string    `json:"username"`
	ImageUrls ImageUrls `json:"image_urls"`
}

func (r Response) ImagePath() string {
	return DIR + "/" + r.Username + ".jpg"
}

type ImageUrls struct {
	Epic   string `json:"epic"`
	Bigger string `json:"bigger"`
	Normal string `json:"normal"`
	Mini   string `json:"mini"`
}

func DBG(format string, args ...interface{}) {
	if DEBUG {
		log.Printf(format, args...)
	}
}
