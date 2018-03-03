package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/d15johro/examensarbete-vt18/osmdecoder/fbsconv/fbs"
	"github.com/d15johro/examensarbete-vt18/osmdecoder/pbconv/pb"
	"github.com/golang/protobuf/proto"
)

var serializationFormat = flag.String("sf", "fbs", "Serialization format")

type metrics struct {
	id                  int
	accessTime          float64
	responseTime        float64
	serializationTime   float64
	deserializationTime float64
	dataSize            int
	filename            string
}

func init() {
	flag.Parse()
}

func main() {
	m := metrics{filename: "./" + *serializationFormat + ".txt"}
	if err := m.setup(); err != nil {
		log.Fatalln(err)
	}
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
		// validate response and get :
		if resp.StatusCode != http.StatusOK {
			log.Fatalln(http.StatusText(resp.StatusCode))
		}
		id, err := strconv.Atoi(resp.Header.Get("id"))
		if err != nil {
			log.Fatalln(err)
		}
		m.id = id
		serializationDuration, err := time.ParseDuration(resp.Header.Get("serializationDuration"))
		if err != nil {
			log.Fatalln(err)
		}
		m.serializationTime = serializationDuration.Seconds() * 1000
		// read response data:
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		m.responseTime = time.Since(startAccessClock).Seconds() * 1000
		m.dataSize = len(data)
		// deserialize data:
		startDeserializationClock := time.Now()
		var osm fbs.OSM // change type depending on serialization format being used
		if err := deserialize(data, &osm); err != nil {
			log.Fatalln(err)
		}
		m.deserializationTime = time.Since(startDeserializationClock).Seconds() * 1000
		m.accessTime = time.Since(startAccessClock).Seconds() * 1000
		m.log(*serializationFormat + ".txt")
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

func (m *metrics) log(filename string) {
	s := fmt.Sprintf("%d,%f,%f,%f,%f,%d\n",
		m.id, m.accessTime,
		m.responseTime,
		m.serializationTime,
		m.deserializationTime,
		m.dataSize)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	if _, err = file.WriteString(s); err != nil {
		log.Fatalln(err)
	}
}

func (m *metrics) setup() error {
	_, err := os.Stat(m.filename)
	if err != nil {
		if !os.IsNotExist(err) { // error even though file exists
			if err := os.Remove(m.filename); err != nil {
				return err
			}
		}
	}
	_, err = os.Create(m.filename)
	return err

}
