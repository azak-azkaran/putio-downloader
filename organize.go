package main

import (
	"errors"
	"github.com/DamnWidget/goqueue"
	"github.com/putdotio/go-putio/putio"
	"hash/crc32"
	"io"
	"os"
	"strconv"
)

var CompletedQueue = goqueue.New()
var checkCompleted = false

func CheckCompleted() {
	for checkCompleted {
		id := CompletedQueue.Pop()

		if id != nil {
			file, success := LinkMap.Get(id.(string))
			obj := file.(PutioObject)
			if success {
				Info.Printf("recieved: %s%s%s", blue, obj.file.CRC32, reset)
			}
		}
	}
}

func Organize(obj PutioObject) error {
	completeFilepath := obj.Foldername + obj.Name
	completeFolderpath := obj.Foldername
	compare := CompareFiles(completeFilepath, obj.file)
	if compare {
		//CreateFolder(completeFolderpath)
		newfolder, err := os.Stat(completeFolderpath)
		if err != nil {
			Error.Println("Error Folder missing will not move file: ", err)
			return err
		} else {
			Info.Println("Folder created: ", newfolder)
			//if len(putFile.Folder) != 0 && newfolder.IsDir() {
			//	utils.Info.Println("Moving to: ", completeFolderpath+"/"+putFile.Name)
			//	err := os.Rename(completeFilepath, completeFolderpath+"/"+putFile.Name)
			//	if err != nil {
			//		utils.Error.Fatalln("Error while moving File: ", putFile.Name, "\n", err)
			//		return
			//	}
			//}
			return nil
		}
	}
	return errors.New("Online file and offline file differ")
}

func CompareFiles(path string, file putio.File) bool {
	offlineFile, err := os.Open(path)
	if err != nil {
		Error.Println("Error while reading file: ", file.Name, "\tError: ", err)
		return false
	}
	hash := crc32.NewIEEE()
	//Copy the file in the interface
	if _, err := io.Copy(hash, offlineFile); err != nil {
		Error.Println("Error while copying file for CRC : ", file.Name, "\tError: ", err)
		return false
	}

	err = offlineFile.Close()
	if err != nil {
		Error.Println("Error while closing file: ", file.Name, "\tError: ", err)
		return false
	}
	//Generate the hash
	hashInBytes := hash.Sum32()

	//Encode the hash to a string
	crc := int64(hashInBytes)

	fileCrc, err := strconv.ParseInt(file.CRC32, 16, 64)
	if err != nil {
		Error.Println("Error while converting CRC from Putio: ", err)
		return false
	}
	if fileCrc != crc {
		Info.Println("CRC values for: ", file.Name, " are different", "\nOnline CRC: ", strconv.FormatInt(fileCrc, 16), "\nOffline CRC: ", strconv.FormatInt(crc, 16))
		return false
	}
	return true
}
