package test

import (
	"github.com/tidwall/gjson"
	"github.com/ChengjinWu/gojson"
	gojson22 "github.com/widuu/gojson"
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

const (
	testCount = 1000000
	memCount = 10

)

func getJsonData() []byte {
	b, err := ioutil.ReadFile("data.json")
	if err != nil {
		fmt.Print(err)
	}
	return b
}

func Test_Gojson(t *testing.T) {
	data := getJsonData()
	urls := make([]string,memCount)
	startTime := time.Now()
	object,_ := gojson.FromBytes(data)
	for i := 0; i < testCount; i++ {
		urls[i%memCount] = object.GetJsonObject("owner").GetJsonObject("received_events_url").GetString()
	}
	t.Log(time.Now().Sub(startTime))
	t.Log(urls[:10])
}

func Test_Gjson(t *testing.T) {
	data := string(getJsonData())
	urls := make([]string,memCount)
	startTime := time.Now()
	for i := 0; i < testCount; i++ {
		urls[i%memCount] = gjson.Get(data, "owner.received_events_url").Str
	}
	t.Log(time.Now().Sub(startTime))
	t.Log(urls[:10])
}

func Test_Gojson22(t *testing.T) {
	data := string(getJsonData())
	urls := make([]string,memCount)
	startTime := time.Now()
	for i := 0; i < testCount; i++ {
		urls[i%memCount] = gojson22.Json(data).Get("owner").Get("received_events_url").Tostring()
	}
	t.Log(time.Now().Sub(startTime))
	t.Log(urls[:10])
}

//
//func Test_Coding_Gojson(t *testing.T) {
//	data := getJsonData()
//	urls := make([]string,memCount)
//	startTime := time.Now()
//	for i := 0; i < testCount; i++ {
//		object,_ := gojson.FromBytes(data)
//		urls[i%memCount] = object.GetCoding("chengjin")
//	}
//	t.Log(time.Now().Sub(startTime))
//	t.Log(urls[:1])
//}