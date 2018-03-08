package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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
	serializationFormat = flag.String("sf", "pb", "Serialization format")
	mapsDir             = "../../data/maps/"
	numberOfFiles       uint32
)

type handler struct{}

func init() {
	flag.Parse()
	var err error
	numberOfFiles, err = fileCount(mapsDir)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	log.Println(numberOfFiles)
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
	// Validate http method:
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	// Validate URL path:
	segs := pathSegments(r.URL.Path)
	if len(segs) != 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// Extract and validate id from path segments:
	id, err := strconv.Atoi(segs[0])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// Decode .osm file depending on id:
	file := mapsDir + "map" + fmt.Sprintf("%d", uint32(id)%numberOfFiles) + ".osm"
	x, err := osmdecoder.DecodeFile(file)
	if err != nil {
		log.Println("write:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// Serialize object:
	var (
		data                    []byte
		startSerializationClock time.Time
		startStructuringClock   time.Time
		structuringDuration     time.Duration
	)
	switch *serializationFormat {
	case "pb":
		startStructuringClock = time.Now()
		osm, err := pbconv.Make(x)
		if err != nil {
			log.Println("write:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		structuringDuration = time.Since(startStructuringClock)
		startSerializationClock = time.Now()
		data, err = proto.Marshal(osm)
		if err != nil {
			log.Println("write:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	case "fbs":
		// Since flatbuffers already stores data in a "serialized" form, serialization basically
		// means getting a pointer to the internal storage. Therefore, unlike pb where we start the
		// serialize clock after the pb.OSM object has been structured, full build/serialize cycle
		// is measured.
		startStructuringClock = time.Now()
		startSerializationClock = time.Now()
		builder := flatbuffers.NewBuilder(0)
		err = fbsconv.Build(builder, x)
		if err != nil {
			log.Println("write:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		data = builder.Bytes[builder.Head():]
		structuringDuration = time.Since(startStructuringClock) // structuring duration will be the same as serialization duration
	default:
		log.Fatalln("serialization format not supported")
	}
	serializationDuration := time.Since(startSerializationClock)
	// Write header and data to response:
	w.Header().Add("id", segs[0])
	w.Header().Add("serializationDuration", serializationDuration.String())
	w.Header().Add("structuringDuration", structuringDuration.String())
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

func fileCount(dirpath string) (uint32, error) {
	i := 0
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return 0, err
	}
	for _, file := range files {
		if !file.IsDir() {
			i++
		}
	}
	return uint32(i), nil
}
