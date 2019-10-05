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
	assert.EqualValues(t, 10, FolderQueue.Len())

	for {
		item := FolderQueue.Pop()
		if item == nil {
			break
		}
	}
}

func Test_CreateFileLink(t *testing.T) {
	fmt.Println("Running Test_CreateFileLink")
	var conf Configuration
	conf.oauthToken = "blub"
	conf.listFunc = returnList

	var one putio.File
	id := int64(0)
	one.ID = id
	one.Name = "Folder" + strconv.FormatInt(id, 10)
	one.ContentType = "application/x-directory"
	CreateFileLink(conf, one)

	time.Sleep(10 * time.Millisecond)
	assert.False(t, LinkMap.IsEmpty())
	assert.Equal(t, 1, LinkMap.Count())
	LinkMap.Remove(strconv.FormatInt(id, 10))
	assert.True(t, LinkMap.IsEmpty())
}

func Test_AddLinks(t *testing.T) {
	fmt.Println("Running Test_AddLinks")
	var conf Configuration
	conf.oauthToken = "blub"
	conf.listFunc = returnList
	AddLinks(conf)

	time.Sleep(1 * time.Millisecond)

	assert.EqualValues(t, 10, LinkMap.Count())
	assert.Equal(t, int64(0), FolderQueue.Len())

	for {
		item := FolderQueue.Pop()
		if item == nil {
			break
		}
	}
}

func Test_CreateConfiguration(t *testing.T) {
	fmt.Println("Running Test_CreateConfiguration")
	conf := CreateConfiguration("blub", "local")
	assert.NotNil(t, conf)
	assert.NotNil(t, conf.listFunc)

	assert.Equal(t, "http://local:6800/jsonrpc", conf.aria2)
}
