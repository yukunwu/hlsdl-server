package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/canhlinh/hlsdl"
)

type Queue struct {
	Data        map[string]*DownloadInfo
	Workers     int
	DownloadDir string
}

var (
	queueIns *Queue
	once     = sync.Once{}
	taskChan = make(chan *DownloadInfo, 1)
)

func GetInstance() *Queue {
	once.Do(func() {
		queueIns = new(Queue)
		queueIns.Data = make(map[string]*DownloadInfo)
	})
	return queueIns
}
func (q *Queue) Run(workers int, downloadDir string) {
	q.Workers = workers
	q.DownloadDir = downloadDir
	defer close(taskChan)
	defer Recover()
	for info := range taskChan {
		q.Data[info.UUID] = info
		go q.handleDownload(info.UUID)
	}
}
func (q *Queue) Add(info *DownloadInfo) string {
	var (
		uuid = fmt.Sprintf("%d", time.Now().UnixNano())
	)
	info.UUID = uuid
	taskChan <- info
	return uuid
}

func (q *Queue) handleDownload(uuid string) {
	var info = q.Data[uuid]
	defer Recover()
	var headers = map[string]string{}
	for _, header := range info.Headers {
		h := strings.SplitN(header, ":", 2)
		if len(h) == 2 {
			headers[strings.TrimSpace(h[0])] = strings.TrimSpace(h[1])
		}
	}
	hlsDL := hlsdl.New(info.Url, headers, q.DownloadDir, info.FileName, q.Workers, false)
	var (
		finished = make(chan bool)
	)
	go func() {
		defer func() {
			Recover()
			finished <- true
		}()
		fp, err := hlsDL.Download()
		if err != nil {
			info.Error = err.Error()
			log.Println("Download:", err)
		} else {
			info.FilePath = fp
		}
	}()

	ticker := time.NewTicker(time.Millisecond * 500)
	for {
		select {
		case <-finished:
			ticker.Stop()
			info.Finished = true
			return
		case <-ticker.C:
			info.Progress = hlsDL.GetProgress()
		}
	}
}
func (q *Queue) GetInfo(uuid string) (error, *DownloadInfo) {
	info, ok := q.Data[uuid]
	if !ok {
		return errors.New(fmt.Sprintf("Download of uuid %s is not exists", uuid)), nil
	}
	return nil, info
}
