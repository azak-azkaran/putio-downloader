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
	assert.True(t, FolderQueue.Len() == 10)

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
	assert.True(t, FolderQueue.Len() == 10)
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

	assert.NotEqual(t, LinkMap.Count(), 0)

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
	assert.NotNil(t, conf)
	assert.NotNil(t, conf.listFunc)
}
