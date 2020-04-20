package epic

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	key string
}

func NewEpicClient(key string) *Client {
	return &Client{
		key: key,
	}
}

func (ec *Client) GetAvailableDates() ([]string, error) {
	resp, err := http.Get("https://api.nasa.gov/EPIC/api/natural/available?api_key=" + ec.key)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	b, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return []string{}, err2
	}

	var arr []string
	err = json.Unmarshal(b, &arr)
	if err != nil {
		return []string{}, err
	}

	return arr, nil
}

type ImageInfo struct {
	Image string `json:"image"`
	Date  string `json:"date"`
}

func (ec *Client) GetInfoByDate(date string) ([]ImageInfo, error) {
	_, err := time.Parse("2006-01-02", date)

	if err != nil {
		return []ImageInfo{}, err
	}

	resp, err2 := http.Get("https://api.nasa.gov/EPIC/api/natural/date/" + date + "?api_key=" + ec.key)
	if err2 != nil {
		return []ImageInfo{}, err2
	}
	defer resp.Body.Close()

	b, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		return []ImageInfo{}, err3
	}

	var res []ImageInfo
	err = json.Unmarshal(b, &res)
	if err != nil {
		return []ImageInfo{}, err
	}

	return res, nil
}

func (ec *Client) SaveImage(date string, name string, path string)  error {
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return err
	}

	_ = os.Mkdir(path + "/"+ date, 0777)
	resp, err2 := http.Get("https://api.nasa.gov/EPIC/archive/natural/" +
		strings.Replace(date, "-", "/", -1) + "/png/" +
		name + ".png" + "?api_key=" + ec.key)
	if err2 != nil {
		return err2
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("Api error, status code: " + strconv.Itoa(resp.StatusCode))
	}


	file, err3 := os.Create(path + "/"+ date + "/" + name + ".png")
	if err3 != nil {
		return err3
	}

	_, err4 := io.Copy(file, resp.Body)
	if err4 != nil {
		return err4
	}
	defer file.Close()


	return nil
}
