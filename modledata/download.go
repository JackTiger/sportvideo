package modledata

import (
    //"sort"
    "time"
	"gopkg.in/mgo.v2/bson"
    "github.com/ginuerzh/sportvideo/common/errors"
   	"github.com/ginuerzh/sportvideo/common/youlog"
)

//SearchQuery SearchQuery
type SearchQuery struct {
    Keyword            string      `bson:"keyword"`
    PageToken            string      `bson:"pagetoken"`
}

//DownloadStateData DownloadStateData
type DownloadStateData struct {
    PageToken       int64
    PageCount       int
}

//DownloadStateRetData DownloadStateRetData
type DownloadStateRetData struct {
    FileName            string      `bson:"filename"`
    ImageURL            string      `bson:"imageurl"`
    FileSize            int64       `bson:"fileSize"`
    VideoID             string      `bson:"videoid"`
    DownloadURL         string      `bson:"downloadurl"`
    DownloadItag        int         `bson:"downloaditag"`
    DownloadResolution  string      `bson:"downloadresolution"`
    DownloadExt         string      `bson:"downloadext"`
    Title               string      `bson:"title"`
    Author              string      `bson:"author"`
    Duration            int64       `bson:"duration"`
    PublishTime         int64       `bson:"publishtime"`
    Createtime          int64       `bson:"createtime"`
}

type createtimeslice []DownloadStateRetData
 
func (a createtimeslice) Len() int {
    return len(a)
}
func (a createtimeslice) Swap(i, j int) {
    a[i], a[j] = a[j], a[i]
}
func (a createtimeslice) Less(i, j int) bool {
    return a[j].Createtime < a[i].Createtime
}

//GetSearchQueryHandler GetSearchQueryHandler
func GetSearchQueryHandler() (retData SearchQuery, e * errors.Error) {
	youlog.Info("GetSearchQueryHandler from database begin")   
    err := FindOne(KeywordCollection, nil, nil, nil, &retData)
    
    if err!=nil{
        youlog.Warnning("Get search query from database failed")
		e = errors.NewError(errors.DbError, "Get search query from database failed")
		return
    }

	youlog.Info("GetSearchQueryHandler from database finished")
    return retData, nil
}

//UpsertSearchQueryHandler UpsertSearchQueryHandler
func UpsertSearchQueryHandler(retData SearchQuery) (e * errors.Error) {
	youlog.Info("UpsertSearchQueryHandler from database begin")

    change :=bson.M{
        "$set":
        bson.M{
            "keyword": retData.Keyword,
            "pagetoken": retData.PageToken,
        },
    }    
    
    err := Upsert(KeywordCollection, nil, change)
    if err!=nil{
	    youlog.Warnning("Insert new record to search query from database failed")
 		e = errors.NewError(errors.DbError, "Insert new record to search query from database failed")
		return e
    }

	youlog.Info("UpsertSearchQueryHandler from database finished")
    return nil
}

//GetDownloadListDataHandler GetDownloadListDataHandler
func GetDownloadListDataHandler(data DownloadStateData) (retDataList []DownloadStateRetData, e * errors.Error) {
	youlog.Info("GetDownloadListDataHandler from database begin")
    query := bson.M{
		"createtime": bson.M{
				"$lt": data.PageToken,
			},
	}
    
    err := FindOneLimit(DownloadCollection, query, data.PageCount, nil, []string{"-createtime"}, &retDataList)
    if err != nil{
		youlog.Warnning("Find download list from database error!!!")
		e = errors.NewError(errors.DbError, "Find download list from database error!!!")
		return retDataList, e
    }
    
    //sort.Sort(createtimeslice(retDataList))
	youlog.Info("GetDownloadListDataHandler from database finished")
    return retDataList, nil
}

//GetVideoInfobyVideoIDHandler GetVideoInfobyVideoIDHandler
func GetVideoInfobyVideoIDHandler(videoid string) (retData DownloadStateRetData, e * errors.Error) {
	youlog.Info("GetVideoInfobyVideoIDHandler from database begin")
    query := bson.M{
        "videoid": videoid,
    }
        
    err := FindOne(DownloadCollection, query, nil, nil, &retData)
    
    if err!=nil{
        youlog.Warnning("Get video info from database failed")
		e = errors.NewError(errors.DbError, "Get video info from database failed")
		return
    }

	youlog.Info("GetVideoInfobyVideoIDHandler from database finished")
    return retData, nil
}

//UpsertVideoInfoHandler UpsertVideoInfoHandler
func UpsertVideoInfoHandler(data DownloadStateRetData) (e * errors.Error) {
	youlog.Info("UpsertVideoInfoHandler from database begin")

    if(data.Createtime == 0){
        data.Createtime = time.Now().Unix()
    }
    
    query := bson.M{
        "videoid": data.VideoID,
        "downloaditag": data.DownloadItag,
    }
    
    change :=bson.M{
        "$set":
        bson.M{
            "filename": data.FileName,
            "imageurl": data.ImageURL,
            "fileSize": data.FileSize,
            "videoid": data.VideoID,
            "downloadurl": data.DownloadURL,
            "downloaditag": data.DownloadItag,
            "downloadresolution": data.DownloadResolution,
            "downloadext": data.DownloadExt,
            "title": data.Title,
            "author": data.Author,
            "duration": data.Duration,
            "publishtime": data.PublishTime,
            "createtime": data.Createtime,
        },
    }    
    
    err := Upsert(DownloadCollection, query, change)
    if err!=nil{
	    youlog.Warnning("Insert new record to download task list from database failed")
 		e = errors.NewError(errors.DbError, "Insert new record to download task list from database failed")
		return e
    }

	youlog.Info("UpsertVideoInfoHandler from database finished")
    return nil
}