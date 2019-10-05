package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

type AriaResult struct {
	Bitfield        int
	CompletedLength int
	Connections     int
	Dir             string
	DownloadSpeed   int
	Files           string
	Gid             string
	NumPieces       int
	PieceLength     int
	Status          string
	TotalLength     int
	UploadLength    int
	UploadSpeed     int
}

func MockHTTPServer(returnValue string) *gin.Engine {
	r := gin.Default()
	r.POST("/jsonrpc", func(ctx *gin.Context) {
		switch returnValue {
		case "ERROR":
			resp := AriaResponse{Id: "qwer", Jsonrpc: "2.0", Error: AriaError{Code: 1, Message: "No such download for GID#e5b086a949391fac"}}
			ctx.JSON(http.StatusBadRequest, resp)
		case "DOWNLOAD":
			resp := AriaResponse{Id: "qwer", Jsonrpc: "2.0", Result: "e5b086a949391fac"}
			ctx.JSON(http.StatusOK, resp)
		}
	})
	return r
}

func Test_TellStatus(t *testing.T) {
	fmt.Println("Running Test_TellStatus")
	StatusRequest := TellStatus("e5b086a949391fac")
	test := "{\"jsonrpc\":\"2.0\",\"id\":\"qwer\",\"method\":\"aria2.tellStatus\",\"params\":[\"e5b086a949391fac\"]}"

	var jsonData []byte
	jsonData, err := json.Marshal(StatusRequest)
	assert.NoError(t, err)
	assert.Equal(t, test, string(jsonData))
}

func Test_AddUri(t *testing.T) {
	fmt.Println("Running Test_AddUri")
	AddUriRequest := AddURI("https://blub/file")
	test := "{\"jsonrpc\":\"2.0\",\"id\":\"qwer\",\"method\":\"aria2.addUri\",\"params\":[[\"https://blub/file\"]]}"

	var jsonData []byte
	jsonData, err := json.Marshal(AddUriRequest)
	assert.NoError(t, err)
	assert.Equal(t, test, string(jsonData))
}

func Test_SendStatus(t *testing.T) {
	fmt.Println("Running Test_SendStatus")
	config := Configuration{aria2: "http://localhost:6800/jsonrpc"}
	server := http.Server{
		Addr:    "localhost:6800",
		Handler: MockHTTPServer("ERROR"),
	}
	go server.ListenAndServe()

	time.Sleep(10 * time.Millisecond)
	var v PutioObject
	v.StatusRequest = TellStatus("e5b086a949391fac")
	req, err := SendTellStatus("", v, config)
	assert.NoError(t, err)
	assert.NotNil(t, req)
	assert.NotNil(t, req.Error)

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err = enc.Encode(&v.StatusRequest)
	assert.NoError(t, err)
	err = server.Shutdown(context.Background())
	assert.NoError(t, err)
}

func Test_SendAddUri(t *testing.T) {
	fmt.Println("Running Test_SendAddUri")
	config := Configuration{aria2: "http://localhost:6800/jsonrpc"}
	server := http.Server{
		Addr:    "localhost:6800",
		Handler: MockHTTPServer("DOWNLOAD"),
	}
	go server.ListenAndServe()

	time.Sleep(10 * time.Millisecond)
	var v PutioObject
	id := "10"
	v.AddRequest = AddURI("http://blub/file")
	v.ID = int64(10)
	req, err := SendAddUri(id, v, config)
	assert.NoError(t, err)
	assert.NotNil(t, req)
	assert.NotNil(t, req.Error)

	assert.EqualValues(t, 1, LinkMap.Count())
	val, success := LinkMap.Pop(id)
	assert.True(t, success)
	v = val.(PutioObject)

	assert.Equal(t, TellStatus("e5b086a949391fac"), v.StatusRequest)

	err = server.Shutdown(context.Background())
	assert.NoError(t, err)
}

func Test_AddDownloads(t *testing.T) {
	fmt.Println("Running Test_AddDownloads")
	addingDownloads = true
	config := Configuration{aria2: "http://localhost:6800/jsonrpc"}
	server := http.Server{
		Addr:    "localhost:6800",
		Handler: MockHTTPServer("DOWNLOAD"),
	}
	go server.ListenAndServe()
	go AddDownloads(config)

	time.Sleep(10 * time.Millisecond)
	var v PutioObject
	id := "10"
	v.AddRequest = AddURI("http://blub/file")
	v.ID = int64(10)

	LinkMap.Set(id, v)
	err := FileQueue.Push(id)
	assert.NoError(t, err)
	time.Sleep(10 * time.Millisecond)

	assert.EqualValues(t, 1, LinkMap.Count())
	val, success := LinkMap.Pop(id)
	assert.True(t, success)
	v = val.(PutioObject)

	assert.Equal(t, TellStatus("e5b086a949391fac"), v.StatusRequest)

	err = server.Shutdown(context.Background())
	assert.NoError(t, err)
	addingDownloads = false
}
