package main

import (
	"bytes"
	"encoding/json"
	"github.com/DamnWidget/goqueue"
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

func SendTellStatus(id string, v PutioObject, conf Configuration) (AriaResponse, error) {
	resp, err := send(v.StatusRequest, conf)
	if err != nil {
		Error.Println("Error")
		return resp, err
	}
	//Info.Printf("%s%s%s", yellow, resp, reset)
	return resp, nil
}

func SendAddUri(id string, v PutioObject, conf Configuration) (AriaResponse, error) {
	resp, err := send(v.AddRequest, conf)
	if err != nil {
		return resp, err
	}
	v.StatusRequest = TellStatus(resp.Result)
	Info.Printf("%s%s%s", yellow, v.StatusRequest, reset)
	LinkMap.Set(id, v)
	return resp, nil
}

func send(value interface{}, conf Configuration) (AriaResponse, error) {
	var r AriaResponse
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	err := encoder.Encode(&value)
	if err != nil {
		Error.Println("Error encoding json request:", err)
		return r, err
	}
	aria := http.Client{}
	Info.Printf("Sending %s to %s%s%s", value, yellow, conf.aria2, reset)
	resp, err := aria.Post(conf.aria2, "application/json", &buf)
	//if err != nil || resp.StatusCode != http.StatusOK {
	if err != nil {
		Error.Println("Error sending json request:", err)
		return r, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&r)
	if err != nil {
		Error.Println("Error while decoding json response:", err)
		return r, err
	}
	//Info.Printf("Response: %s%s%s", magenta, r, reset)
	return r, nil
}

func AddDownloads(conf Configuration) {
	for addingDownloads {
		id := FileQueue.Pop()
		if id != nil {
			file, success := LinkMap.Get(id.(string))
			if success {
				SendAddUri(id.(string), file.(PutioObject), conf)
			}
			StatusQueue.Push(id)
		}
	}
}

func CheckStatus(conf Configuration) {
	for checkStatus {
		id := StatusQueue.Pop()
		if id != nil {
			file, success := LinkMap.Get(id.(string))
			if success {
				SendTellStatus(id.(string), file.(PutioObject), conf)
			}
		}
	}
}
