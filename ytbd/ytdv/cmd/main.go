package main

import (
	//"github.com/youtubeservice/ytdv"
	"youtubedownload/modeldownload/ytbd/ytdv"
	"encoding/json"
	"fmt"
	"flag"
	"net/http"
    "errors"
	"io/ioutil"
    "path/filepath"
	"strings"
    "os"
)

var (
    key = flag.String("q", "NSQ", "Search Query Key")
)

func fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("Failed to fetch " + url)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Failed to parse " + url)
	}
	return body, nil
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func downloadThumbnailFile(url string, fileName string) (string, error){
    bytes, err := fetch(url)
    
    if err != nil {
        return "", err
    }
    
    filePath := getCurrentDirectory() + "/" + fileName + ".jpg"
    err = ioutil.WriteFile(filePath, bytes, 0644)
    
    if err != nil {
        return "", err
    }
    
    return filePath, nil
}

func Encode(voidHandle interface{}) []byte {
	buf, err := json.Marshal(voidHandle)
	if err != nil {
		return buf
	}

	return buf
}

func Decode(mesData []byte, voidHandle interface{}){
    err := json.Unmarshal(mesData, voidHandle)
    
	if err != nil {
		return
	}
}

func main() {
	flag.Parse()
	searList, _ := ytdv.GetSearchListFromQuery(*key, 1)
	
	if len(searList.PreInfoList) > 0 {
		searchItem := searList.PreInfoList[0]
		fmt.Println(string(Encode(searchItem)))
		downloadThumbnailFile(searchItem.ThumbnailURL, searchItem.ID)
		
		/*if len(searchItem.ID) > 0 {
			videoInfo, getErr := ytdv.GetVideoInfoFromID(searchItem.ID)
			
			if getErr == nil {
				fmt.Println(string(Encode(videoInfo)))
			}
		}*/
	}
}
