package processor

import (
	"fmt"
	"io/ioutil"
	"log"
	"main/utils"
	"os"
	"strings"
	"time"
)

type QueueLogger struct {
	queueRequest QueueRequest
	r            *os.File
	w            *os.File
	rescueStdout *os.File
	tempLog      string
}

func NewQueueLogger(queueRequest QueueRequest) *QueueLogger {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	log.New(os.Stdout, "", log.Ldate|log.Ltime)
	return &QueueLogger{queueRequest: queueRequest, r: r, w: w, rescueStdout: rescueStdout}
}

func (f *QueueLogger) Log(strInput string, forceUpload ...bool) {

	upload := false
	if len(forceUpload) > 0 {
		upload = forceUpload[0]
	}

	for _, s := range strings.Split(strInput, "\n") {
		if s != "" {
			log.Println(s)
			fmt.Println("[" + time.Now().UTC().Format("2006-01-01T01:00:00") + "]" + s)
		}
	}

	if upload == true {

		f.w.Close()
		out, _ := ioutil.ReadAll(f.r)
		f.tempLog += string(out)
		f.r, f.w, _ = os.Pipe()
		os.Stdout = f.w

		container := f.queueRequest.LogContainerName
		fileName := f.queueRequest.LogFileName

		_ = utils.PutBlob(f.queueRequest.LogStorageConnectionString, container, fileName, f.tempLog)
	}
}

func (f *QueueLogger) LogSave() {

	f.w.Close()
	out, _ := ioutil.ReadAll(f.r)
	os.Stdout = f.rescueStdout
	log.SetOutput(os.Stdout)

	container := f.queueRequest.LogContainerName
	fileName := f.queueRequest.LogFileName

	_ = utils.PutBlob(f.queueRequest.LogStorageConnectionString, container, fileName, f.tempLog+string(out))

}
