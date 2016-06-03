package main

import (
	"github.com/gin-gonic/gin"
	ytb "github.com/ginuerzh/sportvideo/ytbd"
	log "github.com/ginuerzh/sportvideo/common/youlog"
	"time"
	"net/http"
)

var keywords = []string{"跑步", "跑步训练", "跑步达人", "跑步基础"}

//SearchReq SearchReq
type SearchReq struct {
	PageToken int64 `form:"pagetoken"`
	PageCount int `form:"pagecount" binding:"required"`
}

//SearchItem SearchItem
type SearchItem struct {
	VideoID string `json:"videoid"`
    ThumbnailURL string `json:"preview_url"`
    Title string `json:"title"`
    Author string `json:"author"`
    DatePublished int64 `json:"datepublished"`
    Duration int64 `json:"duration"`
}

//SearchList SearchList
type SearchList struct {
	PageToken int64 `json:"pagetoken"`
	VideoList []SearchItem `json:"videolist"`
}

//VideoInfo VideoInfo
type VideoInfo struct {
	VideoID string `json:"videoid"`
    ThumbnailURL string `json:"preview_url"`
	DownloadURL string `json:"download_url"`
    Title string `json:"title"`
	Description string `json:"desc"`
    Author string `json:"author"`
    DatePublished int64 `json:"datepublished"`
    Duration int64 `json:"duration"`
	ViewCount int `json:"viewcount"`
	LikeCount int `json:"likecount"`
	DislikeCount int `json:"dislikecount"`
}

func getVideoListHandler(c *gin.Context) {
	searchReq := &SearchReq{}
	errExist := c.Bind(searchReq)
	
	if errExist != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Bad Para!!!"})
		return
	}
	
	log.Info(c.Request.Method + ", para is " + string(ytb.Encode(searchReq)))
	videoList, err := ytb.GetDownloadListFromDatabase(searchReq.PageToken, searchReq.PageCount)
	
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Get Video List Error!"})
		return
	}
	
	searchList := SearchList{}
	
	for _, item := range videoList {
		searchItem := SearchItem{}
		searchItem.VideoID = item.VideoID
		searchItem.Title = item.Title
		searchItem.ThumbnailURL = item.ImageURL
		searchItem.Duration = item.Duration
		
		searchItem.DatePublished = item.PublishTime
		searchItem.Author = item.Author
		
		if item.Createtime < searchList.PageToken || searchList.PageToken == 0 {
			searchList.PageToken = item.Createtime
		}
		
		searchList.VideoList = append(searchList.VideoList, searchItem)
	}
	
	log.Info("response_data:" + string(ytb.Encode(searchList)))
	c.JSON(http.StatusOK, gin.H{"response_data": string(ytb.Encode(searchList))})
}

func getVideoInfoHandler(c *gin.Context) {
	videoID := c.Query("videoid")
	log.Info(c.Request.Method + ", para is " + videoID)
	
	if len(videoID) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "Bad Para!!!"})
		return
	}

	videoInfo, err := ytb.GetVideoInfoFromDatabase(videoID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Get VideoInfo Error!"})
		return
	}
	
	vidBaseInfo := ytb.GetVideoInfoBaseFromID(videoID)
	
	if vidBaseInfo == nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Get VideoInfo Error!"})
		return
	}
	
	vidInfo := VideoInfo{}
	vidInfo.Author = videoInfo.Author
	vidInfo.DatePublished = videoInfo.PublishTime
	vidInfo.DownloadURL = videoInfo.DownloadURL
	vidInfo.Duration = videoInfo.Duration
	vidInfo.ThumbnailURL = videoInfo.ImageURL
	vidInfo.Title = videoInfo.Title
	vidInfo.VideoID = videoInfo.VideoID
	vidInfo.Description = vidBaseInfo.Description
	vidInfo.DislikeCount = vidBaseInfo.DislikeCount
	vidInfo.LikeCount = vidBaseInfo.LikeCount
	vidInfo.ViewCount = vidBaseInfo.ViewCount
	
	log.Info("response_data:" + string(ytb.Encode(vidInfo)))
	c.JSON(http.StatusOK, gin.H{"response_data": string(ytb.Encode(vidInfo))})
}

func startSearchTimer() {
	startTimer(search)
}

func startTimer(f func()) {
    go func() {
        for {
            f()
            now := time.Now()
			
            // 计算下一个零点
            next := now.Add(time.Hour * 24)
            next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
            t := time.NewTimer(next.Sub(now))
            <-t.C
        }
    }()
}

func search(){
	/*client := ytb.NewTube()
	err, searchRsp := client.SearchFromKeyword(ytb.SearchReq {
		KeyWord:"跑步",
		PageToken:"CBkQAA",
		PageCount:5,
	})
	
	if err == nil {
		log.Info(string(ytb.Encode(searchRsp)))
	}*/
	
	keyword := keywords[0]
	pageToken := ""
	//Check From Database
	searchQuery, _ := ytb.GetSearchQueryFromDatabase()
	
	if len(searchQuery.Keyword) > 0 {
		keyword = searchQuery.Keyword
		pageToken = searchQuery.PageToken
		
		log.Info("Use local database search query, keyword is " + keyword + ", pagetoken is " + pageToken)
	}
	
	client := ytb.NewTube()
	
	searchReq := ytb.SearchReq {
		KeyWord:keyword,
		PageToken:pageToken,
		PageCount:5,
	}
	
	log.Info("Auto to begin search and download list, para is " + string(ytb.Encode(searchReq)))
	err, searchRsp := client.SearchFromKeyword(searchReq)
	
	if err != nil {
		//Network is not available, do nothing
		log.Info("Search youtube error, network not available, can not connect to youtube google api!!!")
	} else {
		if len(searchRsp.SearchIDS) > 0 {
			//Begin to download
			for _, v := range searchRsp.SearchIDS {
				go ytb.DownloadVideo(v)
			}
			
			ytb.InsertSearchQueryToDatabase(searchRsp.KeyWord, searchRsp.PageToken)
		} else {
			//Get Next Search Keywords
			log.Info("Search result is empty, begin search new key word next day!!!")
			
			for index, v := range keywords {
				if v == keyword {
					if index < len(keywords) - 1 {
						keyword = keywords[index + 1]
						pageToken = ""
						ytb.InsertSearchQueryToDatabase(keyword, pageToken)
					} else {
						log.Info("Search not support, all keyword is used up!!!")
					}
					
					break
				}
			}
		}
	}
}

