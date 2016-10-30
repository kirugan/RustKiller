package workers

import (
	"fmt"
	"github.com/nfnt/resize"
	"image/png"
	"os"
	"path"
	"net/http"
)

const TYPE_RESIZE = 0;
const TYPE_SAVE = 1;
const TYPE_UPLOAD = 2;

type Job struct {
	// type
	T uint
	Filepath string
}

func Spawn(workQ chan *Job, total int) {
	for i := 0; i < total; i++ {
		go worker(i, workQ)
	}
}

func worker(workerIdx int, c chan *Job) {
	for job := range c {
		dispatch(*job)
	}
}

func dispatch(job Job) {
	switch (job.T) {
	case TYPE_RESIZE:
		doResize(job)
	case TYPE_SAVE:
		doSave(job)
	case TYPE_UPLOAD:
		doUpload(job)
	}
}

func doUpload(job Job) {
	buf, _ := os.Open(job.Filepath)
	_, err := http.Post("http://localhost:4444/upload/test", "image/png", buf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Upload ", job.Filepath)
}

func doSave(job Job) {
	// doesnt matter
}

func doResize(job Job) {
	file, err := os.Open(job.Filepath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// decode jpeg into image.Image
	img, err := png.Decode(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// respect only height
	m := resize.Resize(0, 250, img, resize.Bicubic)

	basename := path.Base(job.Filepath)
	out, err := os.Create(basename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer out.Close()

	// write new image to file
	png.Encode(out, m)
	fmt.Println("Resize done", job.Filepath)
}