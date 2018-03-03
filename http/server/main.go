package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/flatbuffers/go"

	"github.com/d15johro/examensarbete-vt18/osmdecoder/fbsconv"
	"github.com/d15johro/examensarbete-vt18/osmdecoder/pbconv"
	"github.com/golang/protobuf/proto"

	"github.com/d15johro/examensarbete-vt18/osmdecoder"
)

var (
	addr                = flag.String("addr", ":8080", "the address to host the server")
	serializationFormat = flag.String("sf", "fbs", "Serialization format")
)

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
	log.Println("running http server on", *addr)
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
	x, err := osmdecoder.DecodeFile(file)
	if err != nil {
		log.Println("write:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// serialize:
	var (
		data                    []byte
		startSerializationClock time.Time
	)
	switch *serializationFormat {
	case "pb":
		startSerializationClock = time.Now()
		osm, err := pbconv.Make(x)
		if err != nil {
			log.Println("write:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		data, err = proto.Marshal(osm)
		if err != nil {
			log.Println("write:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	case "fbs":
		startSerializationClock = time.Now()
		builder := flatbuffers.NewBuilder(0)
		err = fbsconv.Build(builder, x)
		if err != nil {
			log.Println("write:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		data = builder.Bytes[builder.Head():]
	default:
		log.Fatalln("serialization format not supported")
	}
	serializationDuration := time.Since(startSerializationClock)
	// write header
	w.Header().Add("id", segs[0])
	w.Header().Add("serializationDuration", serializationDuration.String())
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		log.Println("write:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func pathSegments(p string) []string {
	return strings.Split(strings.Trim(p, "/"), "/")
}
