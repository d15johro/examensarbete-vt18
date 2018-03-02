package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/d15johro/examensarbete-vt18/osmdecoder/fbsconv/fbs"
	"github.com/d15johro/examensarbete-vt18/osmdecoder/pbconv/pb"
	"github.com/golang/protobuf/proto"
)

type metrics struct {
	id                      int
	accessDuration          float64
	responseDuration        float64
	serializationDuration   float64
	deserializationDuration float64
	dataSize                int
}

func main() {
	c := http.Client{}
	for i := 0; i < 15; i++ {
		startAccessClock := time.Now()
		// send GET request to server:
		url := fmt.Sprintf("http://localhost:8080/%d", i)
		resp, err := c.Get(url)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
		// validate response:
		if resp.StatusCode != http.StatusOK {
			log.Fatalln(http.StatusText(resp.StatusCode))
		}
		id, err := strconv.Atoi(resp.Header.Get("id"))
		if err != nil {
			log.Fatalln(err)
		}
		serializationDuration, err := time.ParseDuration(resp.Header.Get("serializationDuration"))
		if err != nil {
			log.Fatalln(err)
		}
		// read response data:
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		responseDuration := time.Since(startAccessClock)
		dataSize := len(data)
		// deserialize data:
		startDeserializationClock := time.Now()
		var osm pb.OSM // change type depending on serialization format being used
		if err := deserialize(data, &osm); err != nil {
			log.Fatalln(err)
		}
		deserializationDuration := time.Since(startDeserializationClock)
		accessDuration := time.Since(startAccessClock)
		// collect metrics:
		m := metrics{
			id:                      id,
			accessDuration:          accessDuration.Seconds() * 1000,
			serializationDuration:   serializationDuration.Seconds() * 1000,
			deserializationDuration: deserializationDuration.Seconds() * 1000,
			responseDuration:        responseDuration.Seconds() * 1000,
			dataSize:                dataSize,
		}
		// log metrics
		m.log()
	}
}

func deserialize(data []byte, v interface{}) (err error) {
	if osm, ok := v.(*pb.OSM); ok {
		return proto.Unmarshal(data, osm)
	}
	if _, ok := v.(*fbs.OSM); ok {
		v = fbs.GetRootAsOSM(data, 0)
		return
	}
	err = fmt.Errorf("deserialize: type not supported")
	return
}

func (m *metrics) log() {
	log.Println("id:", m.id)
	log.Println("dataSize:", m.dataSize)
	log.Println("accessDuration", m.accessDuration)
	log.Println("responseDuration:", m.responseDuration)
	log.Println("serializationDuration:", m.serializationDuration)
	log.Println("deserializationDuration:", m.deserializationDuration)
}
