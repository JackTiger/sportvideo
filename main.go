package main

import (
	"flag"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-gonic/gin"
    "github.com/ginuerzh/sportvideo/common/youlog"
	"net/http"
	"time"
    "path/filepath"
	"strings"
    "os"
)

var (
    staticDir     string
    dirStatic     string
)

func init() {
    flag.StringVar(&staticDir, "s", "static", "The file share server directory")
	flag.Parse()
    
    dirStatic = getCurrentDirectory() + "/" + staticDir
   
   if !isDirExists(dirStatic) {
       os.Mkdir(dirStatic, 0777)
   }
}

//CORSMiddleware CORSMiddleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept=Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-with")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
		} else {
			c.Next()
		}
	}
}

func isDirExists(path string) bool {
    fi, err := os.Stat(path)
 
    if err != nil {
        return os.IsExist(err)
    }
        
    return fi.IsDir()
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
        youlog.Warnning(err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func main() {
	router := gin.Default()
    router.Use(CORSMiddleware())
    router.Static("/static", dirStatic)
	router.GET("/getvideoslist", getVideoListHandler)
	router.GET("/getvideoinfo", getVideoInfoHandler)
	
    startSearchTimer()
    
	gracehttp.Serve(
		&http.Server{
			Addr:         ":9090",
			Handler:      router,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		})
}