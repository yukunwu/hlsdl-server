package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DownloadInfo struct {
	UUID     string   `json:"-"`
	Url      string   `json:"url"`
	FileName string   `json:"filename,omitempty"`
	Progress float64  `json:"progress,omitempty"`
	Headers  []string `json:"headers,omitempty"`
	Record   bool     `json:"record,omitempty"`
	FilePath string   `json:"filepath,omitempty"`
	Finished bool     `json:"finished"`
	Error    string   `json:"error,omitempty"`
}

func download(c *gin.Context) {
	var info = DownloadInfo{}
	if err := c.ShouldBind(&info); err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}
	uuid := GetInstance().Add(&info)
	c.JSON(http.StatusOK, gin.H{
		"uuid": uuid,
	})
}

func query(c *gin.Context) {
	var uuid = c.DefaultQuery("uuid", "")
	if len(uuid) == 0 {
		c.JSON(http.StatusOK, "uuid is no valid")
		return
	}
	err, info := GetInstance().GetInfo(uuid)
	if err != nil {
		c.JSON(http.StatusOK, err.Error())
	} else {
		c.JSON(http.StatusOK, info)
	}
}
