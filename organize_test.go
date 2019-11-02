package main

import (
	"fmt"
	"github.com/putdotio/go-putio/putio"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
)

func Test_OrganizeDownload(t *testing.T) {
	fmt.Println("Running Test_OrganizeDownload")

	var conf Configuration
	conf.oauthToken = "blub"
	conf.listFunc = returnList

	var one putio.File
	id := int64(0)
	one.ID = id
	one.Name = "testdata1"
	one.ContentType = "application/x-directory"
	CreateFolder(conf, one)
	assert.EqualValues(t, LinkMap.Count(), 1)

	var two putio.File
	two.ParentID = id
	id = int64(1)
	two.ID = id
	two.Name = "output.json"
	CreateFileLink(conf, two)
	assert.EqualValues(t, LinkMap.Count(), 2)

	putio_object, success := LinkMap.Get(strconv.FormatInt(id, 10))
	assert.True(t, success)
	Organize(putio_object.(PutioObject))
	//assert.NoError(t, err)
	clearLinkMap(t)
}

func Test_CompareFiles(t *testing.T) {
	fmt.Println("Running Test_CompareFiles")

	file, err := os.Stat("testdata/output.json")
	assert.NoError(t, err)
	var putio putio.File
	putio.Name = "output.json"
	putio.ID = 23
	putio.CRC32 = "16ec90d5"
	putio.Size = file.Size()
	output := CompareFiles("testdata/output.json", putio)
	assert.True(t, output)
}
