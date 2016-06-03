// ytb
package ytb

import (
	"github.com/google/google-api-go-client/googleapi/transport"
	"google.golang.org/api/youtube/v3"
    "github.com/cavaliercoder/grab"
    log "github.com/ginuerzh/sportvideo/common/youlog"
    "github.com/ginuerzh/sportvideo/ytbd/ytdv"
    "github.com/ginuerzh/sportvideo/modledata"
    "fmt"
    "time"
    "strings"
	"net/http"
    "net/url"
)

const (
    timeoutCount = 200
    developerKey = "AIzaSyDTZj7tbRQscz584zuTAt_xQoIzuxyD9RQ"
)

type YtbService struct {
	svc *youtube.Service
}

func NewTube() YtbService {
	client := &http.Client{Transport: &transport.APIKey{Key: developerKey}}
	s, err := youtube.New(client)
	if err != nil {
		log.Fatal(err.Error())
	}
	return YtbService{svc: s}
}

func getThumbnailURL(videoID string, quality ytdv.ThumbnailQuality) *url.URL {
	u, _ := url.Parse(fmt.Sprintf("http://img.youtube.com/vi/%s/%s.jpg",
		videoID, quality))
	return u
}

func getThumbnailUrl(videoID string) (strPath string) {
    //check location database contain thumbnailUrl Path

    thumbnailUrl := getThumbnailURL(videoID, ytdv.ThumbnailQualityHigh)
    filePath, _ := downloadThumbnailFile(thumbnailUrl.String(), videoID, "/ThumbnailImages")
    
    return filePath
}

func GetVideoInfoBaseFromID(videoID string) *ytdv.VideoInfo {
    vid, err := ytdv.GetVideoInfoFromID(videoID)
    
    if err != nil {
        log.Warnning("Get VideoInfo Error!!!")
        return nil
    }
    
    return vid
}

func formatBytes(i int64) (result string) {
	switch {
	case i > (1024 * 1024 * 1024 * 1024):
		result = fmt.Sprintf("%.02f TB", float64(i)/1024/1024/1024/1024)
	case i > (1024 * 1024 * 1024):
		result = fmt.Sprintf("%.02f GB", float64(i)/1024/1024/1024)
	case i > (1024 * 1024):
		result = fmt.Sprintf("%.02f MB", float64(i)/1024/1024)
	case i > 1024:
		result = fmt.Sprintf("%.02f KB", float64(i)/1024)
	default:
		result = fmt.Sprintf("%d B", i)
	}
	result = strings.Trim(result, " ")
	return
}


func startDownload(videoID string) {
    vidInfo, err := ytdv.GetVideoInfoFromID(videoID)
    
    if err != nil {
        log.Warnning("Download error, get videoinfo fatal!!!")
		return
	}
    
    formats := ytdv.FormatList{}
   
    for _, v := range vidInfo.Formats {
        if len(v.Extension) == 0 || len(v.Resolution) == 0 || len(v.VideoEncoding) == 0 || len(v.AudioEncoding) == 0 || v.AudioBitrate == 0 || v.Extension == "webm" {
            continue
        }
        
        formats = append(formats, v)
    }
    
    var downloadFormat ytdv.Format
    var formatList = formats.Best(ytdv.FormatResolutionKey)
    
    if len(formatList) > 0 {
        downloadFormat = formatList[0]
    } else {
        log.Info(string(Encode(formats)))
        log.Warnning("Download error, not to support format to download!!!, videoID is " + videoID + ", Format list is " + string(Encode(vidInfo.Formats)))
        return
    }
    
    downloadURL, err := vidInfo.GetDownloadURL(downloadFormat)

	if err != nil {
        log.Warnning("Download error, get download url fatal!!!")
		return
	}
    
    downloadReq := modledata.DownloadStateRetData{
        FileName:videoID + "." + downloadFormat.Extension,
        VideoID:videoID,
        DownloadURL:getDownloadURLPath(videoID + "." + downloadFormat.Extension),
        DownloadItag:downloadFormat.Itag,
        DownloadResolution:downloadFormat.Resolution,
        DownloadExt:downloadFormat.Extension,
        Title:vidInfo.Title,
        Author:vidInfo.Author,
        Duration:int64(vidInfo.Duration.Seconds()),
        PublishTime:vidInfo.DatePublished.Unix(),
    }

    downloadYoutubeVideo(downloadURL.String(), downloadReq)
}

func downloadYoutubeVideo(urlStr string, downloadReq modledata.DownloadStateRetData) {
	// create a custom client
	client := grab.NewClient()
	client.UserAgent = ""

	// create requests from command arguments
	reqs := make([]*grab.Request, 0)
    req, err := grab.NewRequest(urlStr)
    
    if err != nil {
        log.Warnning("Download error " + err.Error())
		return
	}
		
    req.Filename = generateDownloadFilePath() + "/" + downloadReq.FileName
    
	reqs = append(reqs, req)

	/*for _, item := range fileList {
		req, err := grab.NewRequest(item.url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}

		req.Filename = item.fileName
		reqs = append(reqs, req)
	}*/
    
	// start file downloads, 3 at a time
    log.Info(fmt.Sprintf(">>>>>>>>>>>>>>>>>>>Downloading file name is %s, file path is %s, download request is %s\n\n", downloadReq.FileName, req.Filename, string(Encode(downloadReq))))
	respch := client.DoBatch(3, reqs...)

	// start a ticker to update progress every 3s
	t := time.NewTicker(3000 * time.Millisecond)

	// monitor downloads
	completed := 0
	inProgress := 0
    reConnect := 0
    var transferred uint64
	responses := make([]*grab.Response, 0)
	for completed < len(reqs) {
		select {
		case resp := <-respch:
			// a new response has been received and has started downloading
			// (nil is received once, when the channel is closed by grab)
			if resp != nil {
				responses = append(responses, resp)
			}

		case <-t.C:
			// clear lines
			if inProgress > 0 {
                log.Info(fmt.Sprintf("\033[%dA\033[K", inProgress))
			}

			// update completed downloads
			for i, resp := range responses {
				if resp != nil && resp.IsComplete() {
					// print final result
					if resp.Error != nil {
                        log.Warnning(fmt.Sprintf("############# [Error] downloading %s: %v\n", resp.Request.URL(), resp.Error))
					} else {
                        log.Info(fmt.Sprintf("%s downloading completely, >>>>>>>>>>>>>>>>>>> %s / %s, speed is %.2f kb/s, total time is %s\n", resp.Filename, formatBytes(int64(resp.BytesTransferred())),
                                formatBytes(int64(resp.Size)), resp.AverageBytesPerSecond() / 1024, resp.Duration().String()))
                        downloadReq.FileSize = int64(resp.Size)
                        downloadReq.ImageURL = getThumbnailUrl(downloadReq.VideoID)
                        InsertVideoInfoToDatabase(downloadReq)
					}

					// mark completed
					responses[i] = nil
					completed++
				}
			}

			// update downloads in progress
			inProgress = 0
			for _, resp := range responses {
				if resp != nil {
                    if resp.Error != nil {
                        log.Warnning(fmt.Sprintf("############# [Error] downloading %s: %v\n", resp.Request.URL(), resp.Error))
                        client.CancelRequest(resp.Request)
                        
                        // mark completed
					    completed++
                        break
                    } else {
                        //Check network timeout if over 1 minute, cancel request.
                        if reConnect > timeoutCount {
                            reConnect = 0
                            transferred = 0
                            log.Warnning(fmt.Sprintf("############# [Error] downloading %s timeout, video id is %s, cancel download request", resp.Filename, downloadReq.VideoID))
                            client.CancelRequest(resp.Request)
                        
                            // mark completed
					        completed++
                            break
                        }
                        
                        if transferred < resp.BytesTransferred() {
                            reConnect = 0
                            transferred = resp.BytesTransferred()
                        } else {
                            reConnect++
                        }
                        
					    inProgress++
                        log.Info(fmt.Sprintf("Downloading %s >>>>>>>>>>>>>>>>>>> %s / %s progress (%.2f%%)\033[K, speed is %.2f kb/s, estimated time %s\n\n", resp.Filename, formatBytes(int64(resp.BytesTransferred())),
                                formatBytes(int64(resp.Size)), 100*resp.Progress(), resp.AverageBytesPerSecond() / 1024, resp.ETA().Sub(time.Now()).String()))
                    }
				}
			}
		}
	}

	t.Stop()
    log.Info(fmt.Sprintf("%s download request end.\n", downloadReq.FileName))
}