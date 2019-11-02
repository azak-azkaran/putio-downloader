package main

import (
	"context"
	"github.com/orcaman/concurrent-map"
	"github.com/putdotio/go-putio/putio"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

var Info = log.New(os.Stdout,
	"",
	log.Lshortfile)

var Error = log.New(os.Stderr,
	"",
	log.Lshortfile)
var LinkMap = cmap.New()

type PutioObject struct {
	ID            int64
	Link          string
	Name          string
	Foldername    string
	file          putio.File
	IsDir         bool
	Status        int
	AddRequest    AddUriRequest
	StatusRequest TellStatusRequest
}

type Configuration struct {
	oauthToken  string
	oauthClient *http.Client
	client      *putio.Client
	listFunc    func(ctx context.Context, id int64) (children []putio.File, parent putio.File, err error)
	aria2       string
}

func CreateConfiguration(oauthToken string, aria2host string) Configuration {
	var conf Configuration
	if len(oauthToken) > 0 {
		conf.oauthToken = oauthToken
		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: conf.oauthToken})
		conf.oauthClient = oauth2.NewClient(context.TODO(), tokenSource)
		conf.client = putio.NewClient(conf.oauthClient)
		conf.listFunc = conf.client.Files.List

		Info.Printf("Using oauth Token: %s%s%s", red, conf.oauthToken, reset)
	} else {
		Error.Println("No Token found")
		panic("No Token found")
	}

	if len(aria2host) > 0 {
		Info.Printf("Using host: %s%s%s", red, aria2host, reset)
		if strings.Contains(aria2host, ":") {
			conf.aria2 = "http://" + aria2host + "/jsonrpc"
		} else {
			conf.aria2 = "http://" + aria2host + ":6800/jsonrpc"
		}
	} else {
		Error.Println("No host found, using localhost instead")
		conf.aria2 = "http://localhost:6800/jsonrpc"
	}
	return conf
}
func main() {

	Info.Println("Starting putio-downloader")
	conf := CreateConfiguration(os.Getenv("PUTIO_TOKEN"), os.Getenv("ARIA2_HOST"))

	Info.Println("Starting")
	addingDownloads = true
	checkStatus = true
	checkCompleted = true
	for {
		go AddLinks(conf)
		Info.Println("Running check for new Files")
		go AddDownloads(conf)
		go CheckStatus(conf)
		go CheckCompleted()
		Info.Println("Starting Download")
		time.Sleep(5 * time.Second)
		Info.Println("still alive")
	}
}
