package main

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/DamnWidget/goqueue"
	"github.com/orcaman/concurrent-map"
	"github.com/putdotio/go-putio/putio"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

type Configuration struct {
	oauthToken  string
	oauthClient *http.Client
	client      *putio.Client
	listFunc    func(ctx context.Context, id int64) (children []putio.File, parent putio.File, err error)
}

type PutioObject struct {
	ID   int64
	Link string
	Name string
	file putio.File
}

var FolderQueue = goqueue.New()
var LinkMap = cmap.New()

func AddObject(value PutioObject) {
	log.Println("Adding: ", value.ID, ": ", value.Name, "\t", value.Link)
	LinkMap.Set(strconv.FormatInt(value.ID, 10), value)
}

func CreateConfiguration(oauthToken string) Configuration {
	var conf Configuration
	if len(oauthToken) > 0 {
		conf.oauthToken = oauthToken
		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: conf.oauthToken})
		conf.oauthClient = oauth2.NewClient(context.TODO(), tokenSource)
		conf.client = putio.NewClient(conf.oauthClient)
		conf.listFunc = conf.client.Files.List
		log.Println("Using oauth Token: ", conf.oauthToken)
	} else {
		log.Fatalln("No Token found")
		panic("No Token found")
	}
	return conf
}

func CreateLink(conf Configuration, value putio.File) {
	if value.IsDir() {
		go AddFolders(value.ID, conf)
	} else {
		go func(token string, value putio.File) {
			var current PutioObject
			var builder strings.Builder
			builder.WriteString("https://api.put.io/v2/files/")
			builder.WriteString(strconv.FormatInt(value.ID, 10))
			builder.WriteString("/download?oauth_token=")
			builder.WriteString(token)
			currentlink := builder.String()
			current.ID = value.ID
			current.Link = currentlink
			current.Name = value.Name
			current.file = value
			AddObject(current)
		}(conf.oauthToken, value)
	}
}

func AddFolders(dir int64, conf Configuration) {
	log.Println("Checking folder: ", strconv.FormatInt(dir, 10))
	list, _, err := conf.listFunc(context.Background(), dir)
	if err != nil {
		log.Println("error folder:", err)
	}

	for _, value := range list {
		FolderQueue.Push(value)
	}
}

func AddLinks(conf Configuration) {
	AddFolders(0, conf)

	for {
		value := FolderQueue.Pop()
		if value == nil {
			time.Sleep(1 * time.Millisecond)
			if FolderQueue.Len() == 0 {
				break
			}
		}
		CreateLink(conf, value.(putio.File))
	}
}
