package main

import (
	"net/http"
	"fmt"
	"os"
	"io"
	"log"
	"strconv"
	"github.com/gorilla/mux"
)

var workQ = make(chan int)

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
	}

	router := mux.NewRouter();
	http.HandleFunc("/something", simpleHandle)
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)

	return;

	for i := 0; i < numWorkers; i++ {
		go worker(i, workQ)
	}

	http.HandleFunc("/v1/upload/preset1/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Println("coming")
		workQ <- 345345
		fmt.Println("processing")
	});
	http.HandleFunc("/upload/test", func(w http.ResponseWriter, r *http.Request) {
		filename := "E:\\11PROJECTS\\RustKiller\\upload\\result\\image.png"

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
	});

	log.Fatal(http.ListenAndServe(":8080", nil));
}

func worker(workerIdx int, c chan int) {
	fmt.Println("inside worker")
	for i := range c {
		fmt.Printf("Worker %d %d\n", workerIdx, i)
	}
}