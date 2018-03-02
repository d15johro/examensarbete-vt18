package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/d15johro/examensarbete-vt18/osmdecoder/pbconv"
	"github.com/golang/protobuf/proto"

	"github.com/d15johro/examensarbete-vt18/osmdecoder"
)

var addr = flag.String("addr", ":8080", "the address to host the server")

type handler struct{}

func main() {
	flag.Parse()
	srv := &http.Server{
		Addr:           *addr,
		Handler:        &handler{},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes, // ~2 MB
	}
	log.Fatalln(srv.ListenAndServe())
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// validate http method:
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	// validate URL path:
	segs := pathSegments(r.URL.Path)
	log.Println(segs)
	if len(segs) != 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// validate and extract id from URL:
	id, err := strconv.Atoi(segs[0])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// decode .osm file depending on id:
	file := "../../testdata/test_data" + fmt.Sprintf("%d", id%6) + ".osm"
	OSM, err := osmdecoder.DecodeFile(file)
	if err != nil {
		log.Println("write:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// make pb.OSM
	startSerializationClock := time.Now()
	pbOSM, err := pbconv.Make(OSM)
	if err != nil {
		log.Println("write:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// serialize:
	data, err := proto.Marshal(pbOSM)
	if err != nil {
		log.Println("write:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	endDurationClock := time.Now()
	serializationDuration := endDurationClock.Sub(startSerializationClock)
	log.Println("serializationDuration", serializationDuration)
	// write header
	w.Header().Add("id", segs[0])
	w.Header().Add("serializationDuration", serializationDuration.String())
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		log.Fatalln(err)
	}
}

func pathSegments(p string) []string {
	return strings.Split(strings.Trim(p, "/"), "/")
}
