//
// Created by kenenbek on 16.12.17.
//
package src

import (
	lib "gosan"
	"gosan/anomaly"
	"io/ioutil"
	"os"
)

type (
	FileInfo struct {
		*lib.File
		metric     *FileMetric
		serverName string

		clientContr,
		contrCache,
		cacheJbod,
		jbodCache,
		cacheContr,
		contrClient string
	}

	FileMetric struct {
		timeStart,
		timeEnd,
		sumDeltaTime float64
	}
)

func NewFileInfo(file *lib.File, con *Controller, jbod *SANJBODController) *FileInfo {
	var rwPrefix string

	filename := file.Filename
	serverName := con.NamingProps.Name
	jbodName := jbod.NamingProps.Name

	switch file.RequestType {
	case lib.WRITE:
		rwPrefix = "_w"
	case lib.READ:
		rwPrefix = "_r"
	}
	f := &FileInfo{
		File:       file,
		serverName: serverName,

		clientContr: serverName + "_" + filename + rwPrefix,
		cacheJbod:   jbodName + "_" + filename + rwPrefix,

		cacheContr:  serverName + "_" + filename + "_ACK" + rwPrefix,
		contrClient: "Client_" + filename + "_ACK" + rwPrefix,

		metric: &FileMetric{},
	}
	return f
}

type (
	LogJson struct {
		Timestamp         float64            `json:"timestamp"`
		Ambience          *AtmosphereControl `json:"ambience"`
		StorageComponents []TraceAble        `json:"storage_components"`
		StState           *DiskState         `json:"storage_state"`
	}
	NamingProps struct {
		Type  string `json:"type"`
		Name  string `json:"-"`
		ID    string `json:"id"`
		Owner string `json:"owner"`
	}
	CommonProps struct {
		Status  anomaly.ComponentStatus `json:"health"`
		DevTemp float64         `json:"dev_temp"`
		Uptime  float64         `json:"uptime"`
	}
	TraceAble interface {
		ResetDiffValues()
	}
)

func NewNamingProps(name, typeID, id string) *NamingProps {
	return &NamingProps{
		Name: name,
		Type: typeID,
		ID:   id,
	}
}

func ParseFileAndUnmarshal(filename string) []byte {

	jsonFile, err := os.Open(filename)
	if err != nil {
		panic("File error")
	}

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic("Read error")
	}

	return bytes
}
