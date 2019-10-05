package main

import (
	"context"
	"strconv"
	"strings"

	"github.com/DamnWidget/goqueue"
	"github.com/putdotio/go-putio/putio"
	"time"
)

var FolderQueue = goqueue.New()

func AddObject(value PutioObject) {
	id := strconv.FormatInt(value.ID, 10)
	if value.IsDir {
		Info.Printf("Found Folder: %s%s/%s%s", yellow, value.Foldername, value.Name, reset)
	} else {
		Info.Printf("Adding: %s%s%s: %s%s%s\t%s", cyan, id, reset, yellow, value.Foldername, reset, value.Name)
		FileQueue.Push(id)
	}
	LinkMap.SetIfAbsent(id, value)
}

func CreateFileLink(conf Configuration, value putio.File) {
	if LinkMap.Has(strconv.FormatInt(value.ID, 10)) {
		return
	}
	var current PutioObject
	var builder strings.Builder
	builder.WriteString("https://api.put.io/v2/files/")
	builder.WriteString(strconv.FormatInt(value.ID, 10))
	builder.WriteString("/download?oauth_token=")
	builder.WriteString(conf.oauthToken)
	currentlink := builder.String()
	current.ID = value.ID
	current.Link = currentlink
	current.Name = value.Name
	current.file = value
	current.IsDir = false
	parent, success := LinkMap.Get(strconv.FormatInt(value.ParentID, 10))
	if success {
		current.Foldername = parent.(PutioObject).Foldername + "/" + parent.(PutioObject).Name
	} else {
		current.Foldername = "/"
	}

	current.AddRequest = AddURI(currentlink)
	current.StatusRequest = TellStatus(currentlink)
	AddObject(current)
}

func CreateFolder(conf Configuration, value putio.File) {
	if LinkMap.Has(strconv.FormatInt(value.ID, 10)) {
		return
	}

	var current PutioObject
	current.ID = value.ID
	current.Name = value.Name
	current.file = value
	current.IsDir = true

	parent, success := LinkMap.Get(strconv.FormatInt(value.ParentID, 10))
	if success {
		current.Foldername = parent.(PutioObject).Foldername + "/" + parent.(PutioObject).Name
	}
	AddObject(current)
}

func AddFolders(dir int64, conf Configuration) {
	Info.Printf("Checking folder: %s%s%s", cyan, strconv.FormatInt(dir, 10), reset)
	list, _, err := conf.listFunc(context.Background(), dir)
	if err != nil {
		Error.Println("error folder:", err)
	}

	for _, value := range list {
		err = FolderQueue.Push(value)
		if err != nil {
			Error.Println("Error while pushing to Queue:", err)
		}
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
		} else {
			file := value.(putio.File)
			if file.IsDir() {
				CreateFolder(conf, file)
				AddFolders(file.ID, conf)
			} else {
				CreateFileLink(conf, file)
			}
		}
	}
}
