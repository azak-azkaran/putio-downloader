package main

import (
	"bytes"
	"encoding/json"
	"github.com/DamnWidget/goqueue"
	"io/ioutil"
	"net/http"
)

var FileQueue = goqueue.New()
var StatusQueue = goqueue.New()
var addingDownloads = false
var checkStatus = false

type AriaResponse struct {
	Id      string
	Jsonrpc string
	Result  string
	Error   AriaError
}

type AriaStatusResponse struct {
	Id      string
	Jsonrpc string
	Result  AriaResult
	Error   AriaError
}
type AriaError struct {
	Code    int
	Message string
}

type AddUriRequest struct {
	Jsonrpc string     `json:"jsonrpc"`
	ID      string     `json:"id"`
	Method  string     `json:"method"`
	Params  [][]string `json:"params"`
}

type TellStatusRequest struct {
	Jsonrpc string   `json:"jsonrpc"`
	ID      string   `json:"id"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
}

type AriaResult struct {
	//	Bitfield        string
	CompletedLength string
	Connections     string
	Dir             string
	DownloadSpeed   string
	Files           []AriaFiles
	Gid             string
	NumPieces       string
	PieceLength     string
	Status          string
	TotalLength     string
	UploadLength    string
	UploadSpeed     string
}

type AriaFiles struct {
	Index           string
	Length          string
	CompletedLength string
	Path            string
	Selected        string
	Uris            []AriaUri
}

type AriaUri struct {
	Uri    string
	status string
}

func AddURI(link string) AddUriRequest {
	request := AddUriRequest{}
	request.Jsonrpc = "2.0"
	request.ID = "qwer"
	request.Method = "aria2.addUri"
	var nested []string
	nested = append(nested, link)
	request.Params = append(request.Params, nested)
	return request
}

func TellStatus(link string) TellStatusRequest {
	request := TellStatusRequest{}
	request.Jsonrpc = "2.0"
	request.ID = "qwer"
	request.Method = "aria2.tellStatus"
	request.Params = append(request.Params, link)
	return request
}

func SendTellStatus(id string, v PutioObject, conf Configuration) (*AriaStatusResponse, error) {
	dump, err := send(v.StatusRequest, conf)
	if err != nil {
		return nil, err
	}

	var resp AriaStatusResponse
	reader := bytes.NewReader(dump)
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&resp)
	if err != nil {
		Error.Printf("recieved: %s%s%s", yellow, string(dump), reset)
		Error.Println("Error while decoding json response:", err)
		return nil, err
	}
	return &resp, nil
}

func SendAddUri(id string, v PutioObject, conf Configuration) (*AriaResponse, error) {
	dump, err := send(v.AddRequest, conf)
	if err != nil {
		return nil, err
	}
	var resp AriaResponse
	reader := bytes.NewReader(dump)
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&resp)
	if err != nil {
		Error.Printf("recieved: %s%s%s", yellow, string(dump), reset)
		Error.Println("Error while decoding json response:", err)
		return nil, err
	}
	v.StatusRequest = TellStatus(resp.Result)
	LinkMap.Set(id, v)
	return &resp, nil
}

func send(value interface{}, conf Configuration) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	err := encoder.Encode(&value)
	if err != nil {
		Error.Println("Error encoding json request:", err)
		return nil, err
	}
	aria := http.Client{}
	//Info.Printf("Sending %s to %s%s%s", value, yellow, conf.aria2, reset)
	resp, err := aria.Post(conf.aria2, "application/json", &buf)
	//if err != nil || resp.StatusCode != http.StatusOK {
	if err != nil {
		Error.Println("Error sending json request:", err)
		return nil, err
	}
	defer resp.Body.Close()
	dump, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return dump, nil
}

func AddDownloads(conf Configuration) {
	var err error
	for addingDownloads {
		id := FileQueue.Pop()
		if id != nil {
			file, success := LinkMap.Get(id.(string))
			if success {
				_, err = SendAddUri(id.(string), file.(PutioObject), conf)
				if err != nil {
					Error.Println("Error sending AddUri request:", err)
				}
			}
			err = StatusQueue.Push(id)
			if err != nil {
				Error.Println("Error while Pushing ID:", id, "\n", err)
			}
		}
	}
}

func CheckStatus(conf Configuration) {
	var err error
	for checkStatus {
		id := StatusQueue.Pop()
		if id != nil {
			file, success := LinkMap.Get(id.(string))
			if success {
				_, err = SendTellStatus(id.(string), file.(PutioObject), conf)
				if err != nil {
					Error.Println("Error sending SendTellStatus request:", err)
				}
				//if resp.Result.Status == "complete" {
				err = CompletedQueue.Push(id)

				if err != nil {
					Error.Println("Error while Pushing ID:", id, "\n", err)
				}
				//} else {
				//Info.Printf("Status of Request %s: %s", id.(string), resp.Result.Status)
				//StatusQueue.Push(id)
				//}
			}
		}
	}
}
