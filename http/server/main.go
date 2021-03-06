package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/d15johro/examensarbete-vt18/expcodec"
	"github.com/d15johro/examensarbete-vt18/fs"
	"github.com/d15johro/examensarbete-vt18/osmdecoder/pbconv"
	flatbuffers "github.com/google/flatbuffers/go"

	"github.com/d15johro/examensarbete-vt18/osmdecoder"
)

var (
	addr                = flag.String("addr", ":8080", "the address to host the server")
	serializationFormat = flag.String("sf", "", "Serialization format")
	mapsDir             = "../../data/maps/"
	numberOfFiles       uint32
)

type handler struct{}

func init() {
	flag.Parse()
	// sf flag is required:
	if *serializationFormat == "" {
		log.Fatal("The following flags must be provided:\nflag\t\tvalue\t\tmeaning\nsf\t\tpb or fb\tThe serializationformat to use")
	}
	var err error
	numberOfFiles, err = fs.FileCount(mapsDir)
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
		WriteTimeout:   1000 * time.Second,         // we are writing large files to client...
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
	filename := mapsDir + "map" + fmt.Sprintf("%d", uint32(id)%numberOfFiles) + ".osm"
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	x, err := osmdecoder.Decode(file)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	file.Close()
	// get original file size:
	fileinfo, err := os.Stat(filename)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	originalDataSize := fileinfo.Size()
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
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		structuringDuration = time.Since(startStructuringClock)
		startSerializationClock = time.Now()
		data, err = expcodec.SerializePB(osm)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	case "fbs":
		// Since flatbuffers already stores data in a "serialized" form, serialization basically
		// means getting a pointer to the internal storage. Therefore, unlike pb where we start the
		// serialize clock after the pb.OSM object has been structured, full build/serialize cycle
		// is measured.
		startSerializationClock = time.Now()
		builder := flatbuffers.NewBuilder(0)
		data, err = expcodec.SerializeFBS(builder, x)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		// Structuring time is not meaused in fbs since its the same as serialization time. We send
		// the default value of time.Duration back to the client.
	default:
		log.Fatalln("serialization format not supported")
	}
	serializationDuration := time.Since(startSerializationClock)
	// Write header and data to response:
	w.Header().Add("id", segs[0])
	w.Header().Add("originalDataSize", fmt.Sprint(originalDataSize))
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
