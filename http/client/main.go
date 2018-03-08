package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/d15johro/examensarbete-vt18/expcodec"
	"github.com/d15johro/examensarbete-vt18/metrics"
)

var (
	serializationFormat = flag.String("sf", "pb", "Serialization format")
)

func init() {
	flag.Parse()
}

func main() {
	m := metrics.New()
	m.Filepath = "./http_" + *serializationFormat + ".txt"
	if err := m.Setup(); err != nil {
		log.Fatalln(err)
	}
	c := http.Client{}
	for i := 0; i < 10; i++ { // experimental
		log.Println(i)
		startAccessClock := time.Now()
		startResponseClock := time.Now()
		// Send GET request to server:
		url := fmt.Sprintf("http://localhost:8080/%d", i)
		resp, err := c.Get(url)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
		// Validate response:
		if resp.StatusCode != http.StatusOK {
			log.Fatalln(http.StatusText(resp.StatusCode))
		}
		// Collect metrics from response header:
		id, err := strconv.Atoi(resp.Header.Get("id"))
		if err != nil {
			log.Fatalln(err)
		}
		m.ID = uint32(id)
		serializationDuration, err := time.ParseDuration(resp.Header.Get("serializationDuration"))
		if err != nil {
			log.Fatalln(err)
		}
		m.SerializationTime = serializationDuration.Seconds() * 1000
		structuringDuration, err := time.ParseDuration(resp.Header.Get("structuringDuration"))
		if err != nil {
			log.Fatalln(err)
		}
		m.StructuringTime = structuringDuration.Seconds() * 1000
		// Read response data:
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		m.ResponseTime = time.Since(startResponseClock).Seconds() * 1000
		m.DataSize = len(data)
		// Deserialize data:
		startDeserializationClock := time.Now()
		switch *serializationFormat {
		case "pb":
			if err := expcodec.DeserializePB(data); err != nil {
				log.Fatalln(err)
			}
		case "fbs":
			if err := expcodec.DeserializeFbs(data); err != nil {
				log.Fatalln(err)
			}
		default:
			log.Fatalln("serialization format not supported")
		}
		m.DeserializationTime = time.Since(startDeserializationClock).Seconds() * 1000
		m.AccessTime = time.Since(startAccessClock).Seconds() * 1000

		m.Log()
	}
}
