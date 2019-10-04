package main

import (
	"context"
	"fmt"
	"github.com/putdotio/go-putio/putio"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func returnList(ctx context.Context, id int64) ([]putio.File, putio.File, error) {
	var one putio.File

	var children [10]putio.File

	for i := 0; i < 10; i++ {
		one.ID = int64(i)
		one.Name = "Folder" + strconv.FormatInt(int64(i), 10)
		children[i] = one
	}

	return children[0:10], one, nil
}

func Test_AddFolders(t *testing.T) {
	fmt.Println("Running Test_AddFolders")

	var conf Configuration
	conf.listFunc = returnList

	AddFolders(0, conf)

	time.Sleep(1 * time.Millisecond)
	if FolderQueue.Len() != 10 {
		t.Error("Queue does not have the right size: ", FolderQueue.Len())
	}

	for {
		item := FolderQueue.Pop()
		if item == nil {
			break
		}
	}
}

func Test_CreateLink(t *testing.T) {
	fmt.Println("Running Test_CreateLink")
	var conf Configuration
	conf.oauthToken = "blub"
	conf.listFunc = returnList

	var one putio.File
	one.ID = int64(0)
	one.Name = "Folder" + strconv.FormatInt(int64(0), 10)
	one.ContentType = "application/x-directory"
	CreateLink(conf, one)

	time.Sleep(10 * time.Millisecond)
	if FolderQueue.Len() != 10 {
		t.Error("Queue does not have the right size: ", FolderQueue.Len())
	}
	for {
		item := FolderQueue.Pop()
		if item == nil {
			break
		}
	}

}

func Test_AddLinks(t *testing.T) {
	fmt.Println("Running Test_AddLinks")
	var conf Configuration
	conf.oauthToken = "blub"
	conf.listFunc = returnList
	AddLinks(conf)

	time.Sleep(1 * time.Millisecond)

	if LinkMap.Count() == 0 {
		t.Error("Link Map is empty")
	}

	for {
		item := FolderQueue.Pop()
		if item == nil {
			break
		}
	}
}

func Test_CreateConfiguration(t *testing.T) {
	fmt.Println("Running Test_CreateConfiguration")
	conf := CreateConfiguration("blub")
	if conf.client == nil {
		t.Error("Configuration was not created proberly")
	}

	if conf.listFunc == nil {
		t.Error("Configuration is missing the list function")
	}
}
