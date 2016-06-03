// ytbmsg
package ytb

import (
	"log"
)

type SearchReq struct {
	KeyWord string `json:"keyword"`
	PageToken string `json:"next_token"`
    PageCount int64 `json:"pagecount"`
}

type SearchRsp struct {
    KeyWord string `json:"keyword"`
    PageToken string `json:"next_token"`
    SearchIDS []string `json:"video_ids"`
}

func (y *YtbService) SearchFromKeyword(searchReq SearchReq) (errRsp error, searchRsp SearchRsp) {
    search := y.svc.Search.List("id")
    search = search.Q(searchReq.KeyWord)
    search = search.MaxResults(searchReq.PageCount)
    search = search.PageToken(searchReq.PageToken)
    search = search.Type("video")
    search = search.Order("viewCount")
    
    results, err := search.Do()
	if err != nil {
		log.Println("could not retrieve video list:" + err.Error())
		return err, searchRsp
	}
    
    searchRsp.KeyWord = searchReq.KeyWord
    searchRsp.PageToken = results.NextPageToken
    
    var ids []string
	for _, v := range results.Items {
		ids = append(ids, v.Id.VideoId)
	}

    searchRsp.SearchIDS = ids
    return nil, searchRsp
}

//DownloadVideo DownloadVideo
func DownloadVideo(videoID string) {
    //Check downlad before from database
    videoInfo, _ := GetVideoInfoFromDatabase(videoID)
    
    if len(videoInfo.VideoID) > 0 {
        log.Println("This video is downloaded before, videoID is " + videoID)
    } else {
        startDownload(videoID)
    }
}