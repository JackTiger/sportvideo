// ytbcommon
package ytb

import (
	"net/http"
    "net"
    "errors"
	"io/ioutil"
    "path/filepath"
    "encoding/json"
	"strings"
    "os"
    "fmt"
)

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

// Fetch function
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

func isDirExists(path string) bool {
    fi, err := os.Stat(path)
 
    if err != nil {
        return os.IsExist(err)
    } else {
        return fi.IsDir()
    }
 
    panic("not reached")
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func generateDownloadFilePath() string {
    fileRoot := getCurrentDirectory() + "/static/videos"
    
    if !isDirExists(fileRoot) {
        os.Mkdir(fileRoot, 0777)
    }
    
    return fileRoot
}

func downloadThumbnailFile(url string, fileName string, fileDir string) (string, error){
    bytes, err := fetch(url)
    
    if err != nil {
        return "", err
    }
    
    dirResolutionPath :=  getCurrentDirectory() + "/static" + fileDir
    
    if !isDirExists(dirResolutionPath) {
        os.Mkdir(dirResolutionPath, 0777)
    }
    
    filePath := dirResolutionPath + "/" + fileName + ".png"
    err = ioutil.WriteFile(filePath, bytes, 0644)
    
    if err != nil {
        return "", err
    }
    
    filePath = "http://" + getLocalAddr() + ":9090/static" + fileDir + "/" + fileName + ".png"
    return filePath, nil
}

func getThumbnailURLPath(fileName string) string{
    return "http://" + getLocalAddr() + ":9090/static" + "/ThumbnailImages/" + fileName + ".png"
}

func getDownloadURLPath(fileName string) string{
    return "http://" + getLocalAddr() + ":9090/static/videos/" + fileName
}

func getLocalAddr() string { //Get ip
    /*conn, err := net.Dial("udp", "baidu.com:80")
    if err != nil {
        fmt.Println(err.Error())
        return "Erorr"
    }
    defer conn.Close()
    return strings.Split(conn.LocalAddr().String(), ":")[0]*/
    
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
    
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
			}
		}
	}
    
    return ""
}
