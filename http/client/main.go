package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/d15johro/examensarbete-vt18/osmdecoder/pbconv/pb"
	"github.com/golang/protobuf/proto"
)

type metrics struct {
	id                      int
	accessDuration          time.Duration
	responseDuration        time.Duration
	serializationDuration   time.Duration
	deserializationDuration time.Duration
	dataSize                int
}

func main() {
	doPB()
}

func doPB() {
	c := http.Client{}
	for i := 0; i < 10; i++ {
		startAccessClock := time.Now()
		url := fmt.Sprintf("http://localhost:8080/%d", i)
		resp, err := c.Get(url)
		if err != nil {
			log.Fatalln(err)
		}
		id, _ := strconv.Atoi(resp.Header.Get("id"))
		serializationDuration, _ := time.ParseDuration(resp.Header.Get("serializationDuration"))
		m := metrics{
			id:                    id,
			responseDuration:      time.Since(startAccessClock),
			serializationDuration: serializationDuration,
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Fatalln("got status code:", resp.StatusCode)
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		m.dataSize = len(data)
		startDeserializationClock := time.Now()
		var pbOSM pb.OSM
		if err := proto.Unmarshal(data, &pbOSM); err != nil {
			log.Fatalln(err)
		}
		m.deserializationDuration = time.Since(startDeserializationClock)
		m.accessDuration = time.Since(startAccessClock)
		if err != nil {
			log.Fatalln(err)
		}
		m.log()
	}
}

func (m *metrics) log() {
	log.Println("id:", m.id)
	log.Println("dataSize:", m.dataSize)
	log.Println("accessDuration", m.accessDuration)
	log.Println("responseDuration:", m.responseDuration)
	log.Println("serializationDuration:", m.serializationDuration)
	log.Println("deserializationDuration:", m.deserializationDuration)
}
