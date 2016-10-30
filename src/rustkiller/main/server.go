package main

import (
	"net/http"
	"os"
	"io"
	"log"
	"strconv"
	"github.com/gorilla/mux"
	"rustkiller/workers"
	"hash/fnv"
	"time"
	"math/rand"
)

var workQ = make(chan *workers.Job)
const FINAL_IMAGE = "/Users/kpaltsev/test/rustkiller/upload/result/image.png";

func main() {
	var numWorkers int;
	var err error;

	if len(os.Args) < 2 {
		numWorkers = 5
	} else {
		numWorkers, err = strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}

		if numWorkers <= 0 {
			log.Fatal("Number of workers should be more than 0")
		}
	}

	workers.Spawn(workQ, numWorkers)

	router := mux.NewRouter();
	router.HandleFunc("/upload/test", uploadTest)
	router.HandleFunc("/v1/upload/preset1/{id:[0-9]+}", uploadPreset)
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":4444", nil))
}

func uploadPreset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageId := vars["id"]

	hash := hash(imageId + time.Now().Format(time.UnixDate) + string(rand.Int()))
	filename := strconv.FormatUint(uint64(hash), 10)
	filepath := "/Users/kpaltsev/test/rustkiller/upload/" + filename + ".png";

	file, err := os.OpenFile(filepath, os.O_WRONLY | os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	resizeJob := workers.Job{T: workers.TYPE_RESIZE, Filepath: filepath}
	workQ <- &resizeJob

	uploadJob := workers.Job{T: workers.TYPE_UPLOAD, Filepath: filepath}
	workQ <- &uploadJob
}

func uploadTest(w http.ResponseWriter, r *http.Request) {
	filename := FINAL_IMAGE
	file, err := os.OpenFile(filename, os.O_WRONLY | os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}