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

func clearLinkMap(t *testing.T) {
	for id := range LinkMap.Items() {
		LinkMap.Remove(id)
	}
	assert.True(t, LinkMap.IsEmpty())
}

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
	clearLinkMap(t)
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
	assert.False(t, LinkMap.IsEmpty())
	assert.Equal(t, 1, LinkMap.Count(), LinkMap.Items())

	var two putio.File
	two.ParentID = id
	id = int64(1)
	two.ID = id
	two.Name = "File" + strconv.FormatInt(id, 10)
	CreateFileLink(conf, two)

	assert.False(t, LinkMap.IsEmpty())
	assert.Equal(t, 2, LinkMap.Count(), LinkMap.Items())
	clearLinkMap(t)
}

func Test_CreateFolder(t *testing.T) {
	fmt.Println("Running Test_CreateFolder")
	var conf Configuration
	conf.oauthToken = "blub"
	conf.listFunc = returnList

	var one putio.File
	id := int64(0)
	one.ID = id
	one.Name = "Folder" + strconv.FormatInt(id, 10)
	one.ContentType = "application/x-directory"
	CreateFolder(conf, one)
	assert.EqualValues(t, LinkMap.Count(), 1)
	CreateFolder(conf, one)
	assert.EqualValues(t, LinkMap.Count(), 1)
	clearLinkMap(t)
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
	clearLinkMap(t)
}

func Test_CreateConfiguration(t *testing.T) {
	fmt.Println("Running Test_CreateConfiguration")
	conf := CreateConfiguration("blub", "local")
	assert.NotNil(t, conf)
	assert.NotNil(t, conf.listFunc)

	assert.Equal(t, "http://local:6800/jsonrpc", conf.aria2)
}
