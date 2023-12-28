package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/aoticombr/golang/m3u8"
	"github.com/google/uuid"
)

func TestUp_Down(t *testing.T) {
	fmt.Println("Teste")
	down := m3u8.NewM3u8()
	byt, err := down.GetVideoByte("https://b-vz-22fba020-9c7.tv.pandavideo.com.br/c1808569-132d-4bcc-980b-9a1e8bef9698/playlist.m3u8")
	if err != nil {
		fmt.Println(err)
	}
	newUUID := uuid.New()
	timeString := time.Now().Format("2006-01-02-15-04-05")
	Name := timeString + newUUID.String() + "playlist.mkv"

	m3u8.SaveByteToFile(Name, byt)
}
