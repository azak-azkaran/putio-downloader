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
)

type AriaResponse struct {
	Id      string
	Jsonrpc string
	Result  AriaResult
	Error   AriaError
}

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

//u'bitfield': u'0000000000',
//             u'completedLength': u'901120',
//             u'connections': u'1',
//             u'dir': u'/downloads',
//             u'downloadSpeed': u'15158',
//             u'files': [{u'index': u'1',
//                         u'length': u'34896138',
//                         u'completedLength': u'34896138',
//                         u'path': u'/downloads/file',
//                         u'selected': u'true',
//                         u'uris': [{u'status': u'used',
//                                    u'uri': u'http://example.org/file'}]}],
//             u'gid': u'2089b05ecca3d829',
//             u'numPieces': u'34',
//             u'pieceLength': u'1048576',
//             u'status': u'active',
//             u'totalLength': u'34896138',
//             u'uploadLength': u'0',
//             u'uploadSpeed': u'0'}}

type AriaError struct {
	Code    int
	Message string
}

func MockHTTPServer() *gin.Engine {
	r := gin.Default()
	r.POST("/jsonrpc", func(ctx *gin.Context) {
		resp := AriaResponse{Id: "qwer", Jsonrpc: "2.0", Error: AriaError{Code: 1, Message: "No such download for GID#e5b086a949391fac"}}
		ctx.JSON(http.StatusBadRequest, resp)
	})
	return r
}

func Test_TellStatus(t *testing.T) {

	StatusRequest := TellStatus("e5b086a949391fac")
	test := "{\"jsonrpc\":\"2.0\",\"id\":\"qwer\",\"method\":\"aria2.tellStatus\",\"params\":[\"e5b086a949391fac\"]}"

	var jsonData []byte
	jsonData, err := json.Marshal(StatusRequest)
	assert.NoError(t, err)
	assert.Equal(t, test, string(jsonData))
}

func Test_SendStatus(t *testing.T) {
	config := Configuration{aria2: "http://localhost:6800/jsonrpc"}
	server := http.Server{
		Addr:    "localhost:6800",
		Handler: MockHTTPServer(),
	}
	go server.ListenAndServe()

	var v PutioObject
	v.StatusRequest = TellStatus("e5b086a949391fac")
	req := SendTellStatus("", v, config)
	assert.NotNil(t, req)

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(&v.StatusRequest)
	fmt.Println("Request:", string(buf.Bytes()))

	err := server.Shutdown(context.Background())
	assert.NoError(t, err)
}
