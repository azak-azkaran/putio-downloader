package main

import (
	"bytes"
	"encoding/json"
	"github.com/DamnWidget/goqueue"
	"net/http"
	"net/http/httputil"
)

var FileQueue = goqueue.New()

type Response struct {
	ID      string
	Jsonrpc string
	Result  string
	Error   string
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

func SendTellStatus(id string, v PutioObject, conf Configuration) Response {
	resp, err := send(v.StatusRequest, conf)
	if err != nil {
		Error.Println("Error")
		return resp
	}
	Info.Printf("%s%s%s", yellow, resp, reset)
	return resp
}

func SendAddUri(id string, v PutioObject, conf Configuration) {
	resp, err := send(v.AddRequest, conf)
	if err != nil {
		return
	}
	Info.Printf("%s%s%s", yellow, resp, reset)
	v.StatusRequest = TellStatus(resp.Result)
	LinkMap.Set(id, v)
}

func send(value interface{}, conf Configuration) (Response, error) {
	var r Response
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(&value)
	aria := http.Client{}
	Info.Printf("Sending %s to %s%s%s", value, yellow, conf.aria2, reset)
	resp, err := aria.Post(conf.aria2, "application/json", &buf)
	//if err != nil || resp.StatusCode != http.StatusOK {
	if err != nil {
		Error.Println("Error sending json request:", err)
		return r, err
	}
	//decoder := json.NewDecoder(resp.Body)
	//err = decoder.Decode(&r)
	//if err != nil {
	//	Error.Println("Error while decoding json response:", err)
	//	return r, err
	//}
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		Error.Println("Error while dumping reponse:", err)
	}
	Info.Printf("Response: %s%s%s", magenta, string(dump), reset)
	defer resp.Body.Close()
	return r, nil
}

func AddDownloads(conf Configuration) {
	for {
		id := FileQueue.Pop()
		if id != nil {
			file, success := LinkMap.Pop(id.(string))
			if success {
				SendAddUri(id.(string), file.(PutioObject), conf)
			}
			FileQueue.Push(id)
		}
	}
}

func CheckStatus(conf Configuration) {
	for {
		id := FileQueue.Pop()
		if id != nil {
			file, success := LinkMap.Pop(id.(string))
			if success {
				SendTellStatus(id.(string), file.(PutioObject), conf)
			}
		}
	}
}
