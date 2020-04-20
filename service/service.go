package service

import (
	"context"
	"fmt"
	"go-test/pkg/epic"
	"log"
	"strconv"
	"strings"
	"time"
)

type DownloadService struct {
	epicClient *epic.Client
	guard chan struct{}
}

func NewDownloadService(key string, nrGoRoutines int) *DownloadService {
	return &DownloadService{
		epic.NewEpicClient(key),
		make(chan struct{}, nrGoRoutines),
	}
}

type ResponseStruct struct {
	data string
}

func (ds *DownloadService) DownloadNatural(qtx context.Context, in *RequestString) (*ResponseString, error) {
	_, err := time.Parse("2006-01-02", in.Date)

	if err != nil {
		return &ResponseString{
			Message: "Incorrect date format",
			Error:   true,
		}, nil
	}

	availableDates, err2 := ds.epicClient.GetAvailableDates()
	if err2 != nil {
		log.Println(err2)
		return &ResponseString{
			Message: "Nasa Epic api error",
			Error:   true,
		}, nil
	}

	found := find(availableDates, in.Date)
	if !found {
		return &ResponseString{
			Message: "The date not found in available dates",
			Error:   true,
		}, nil
	}

	imagesInfo, err3 := ds.epicClient.GetInfoByDate(in.Date)
	if err3 != nil {
		log.Println(err3)
		return &ResponseString{
			Message: "Nasa Epic api error",
			Error:   true,
		}, nil
	}

	select {
	case ds.guard <- struct{}{}:
		fmt.Println("Not full")
		go func() {
			ds.downloadAllImages(imagesInfo)
			<-ds.guard
		}()
	default:
		fmt.Println("Is full")
		return &ResponseString{
			Message: "Too many downloads at the same time, please wait",
			Error:   false,
		}, nil
	}

	return &ResponseString{
		Message: "The images from date " + in.Date + " started to download",
		Error:   false,
	}, nil
}

func (ds *DownloadService) downloadAllImages(imagesInfo []epic.ImageInfo) {
	count := 0
	for _, item := range imagesInfo {
		dateTime := strings.Split(item.Date, " ")
		err := ds.epicClient.SaveImage(dateTime[0], item.Image, "./downloads")
		if err != nil {
			log.Println("Save image error: ", err)
		} else {
			count++
		}
	}
	log.Println(strconv.Itoa(count) + " images from " + strconv.Itoa(len(imagesInfo)) +
		" were successfully downloaded")
}

func find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
