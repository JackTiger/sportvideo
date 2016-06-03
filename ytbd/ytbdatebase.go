// ytbdatebase
package ytb

import (
    "time"
    "fmt"
    "github.com/ginuerzh/sportvideo/modledata"
    log "github.com/ginuerzh/sportvideo/common/youlog"
)

func GetDownloadListFromDatabase(pageToken int64, pageCount int) ([]modledata.DownloadStateRetData, error) {
    downloadStateData := modledata.DownloadStateData {
        PageCount:pageCount,
        PageToken:pageToken,
    }
    
    if downloadStateData.PageToken == 0 {
        downloadStateData.PageToken = time.Now().Unix()
    }
    
    log.Info("GetDownloadListFromDatabase, para is " + string(Encode(downloadStateData)))
    videoList, err := modledata.GetDownloadListDataHandler(downloadStateData)
    
    if err != nil {
        log.Warnning(fmt.Sprintf("GetDownloadListFromDatabase error is %v: ", err))
        return nil, err
    }
    
    log.Info("GetDownloadListFromDatabase successfully, data is " + string(Encode(videoList)))
    return videoList, nil
}

func GetSearchQueryFromDatabase() (modledata.SearchQuery, error) {
    searchQuery, err := modledata.GetSearchQueryHandler()
    
    if err != nil {
        log.Warnning(fmt.Sprintf("getSearchQueryFromDatabase error is %v: ", err))
        return searchQuery, err
    }
    
    log.Info("getSearchQueryFromDatabase successfully, data is " + string(Encode(searchQuery)))
    return searchQuery, nil
}

func InsertSearchQueryToDatabase(keyword string, pagetoken string) error {
    searchQuery := modledata.SearchQuery {
        Keyword:keyword,
        PageToken:pagetoken,
    }
    
    log.Info("insertSearchQueryToDatabase, database is " + string(Encode(searchQuery)))
    err := modledata.UpsertSearchQueryHandler(searchQuery)
    
    if err != nil {
        log.Warnning(fmt.Sprintf("insertSearchQueryToDatabase error is %v: ", err))
    } else {
        log.Info("insertSearchQueryToDatabase successfully!!!")
    }
    
    return err
}

func GetVideoInfoFromDatabase(videoID string) (modledata.DownloadStateRetData, error) {
    videoInfo, err := modledata.GetVideoInfobyVideoIDHandler(videoID)
    
    if err != nil {
        log.Warnning(fmt.Sprintf("getVideoInfoFromDatabase error is %v: ", err))
        return videoInfo, err
    }
    
    log.Info("getVideoInfoFromDatabase successfully, data is " + string(Encode(videoInfo)))
    return videoInfo, nil
}

func InsertVideoInfoToDatabase(videoInfo modledata.DownloadStateRetData) error {
    log.Info("insertVideoInfoToDatabase, database is " + string(Encode(videoInfo)))
    err := modledata.UpsertVideoInfoHandler(videoInfo)
    
    if err != nil {
        log.Warnning(fmt.Sprintf("insertVideoInfoToDatabase error is %v: ", err))
    } else {
        log.Info("insertVideoInfoToDatabase successfully!!!")
    }
    
    return err
}