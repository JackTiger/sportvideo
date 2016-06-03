package ytdv

import (
    "bytes"
	"encoding/xml"
	"fmt"
    "time"
    "net/url"
	"io/ioutil"
	"net/http"
	"strings"
    "strconv"
    "regexp"
    "github.com/PuerkitoBio/goquery"
)

const suggestBaseURL = "https://suggestqueries.google.com/complete/search?client=&output=toolbar&ds=yt&hl=en&q="
const searchQueryURL = "https://www.youtube.com/results?search_query=%s&page=%d&filters=video"
const webPageURLPre = "https://www.youtube.com"

// PreviewInfo contains the info of a search item
type PreviewInfo struct {
	// The video ID
	ID string `json:"id"`
	// The video title
	Title string `json:"title"`
    // Author of the video
	Author string `json:"author"`
    // Thumbnail url of the video
    ThumbnailURL string `json:"thumbnailurl"`
    // Web page url of the video
    WebPageURL string `json:"webpageurl"`
    // The date the video was published
	DatePublished string `json:"datePublished"`
    // View count of the video
    ViewCount int `json:"viewcount"`
	// Duration of the video
	Duration time.Duration
}

// SearchList contains the info of search result 
type SearchList struct {
	// The search item
    PreInfoList []PreviewInfo `json:"preinfolist"`

    // Page of the list
    Page int `json:"page"`
}

// GetSuggestionListFromQuery fetches suggest search query from key
func GetSuggestionListFromQuery(query string) ([]string, error) {
    u := fmt.Sprintf("%s%s", suggestBaseURL, query)
    fmt.Println(u)
	resp, err := http.Get(u)
    if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Invalid status code: %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return getSuggestionListFromXML(body), nil
}

func GetSearchListFromQuery(query string, page int) (*SearchList, error) {
    u := fmt.Sprintf(searchQueryURL, query, page)
    fmt.Println(u)
    
    httpclient := &http.Client{}
	req, err := http.NewRequest("POST", u, nil)
    req.Header.Add("Accept-Language", "en")
	resp, err := httpclient.Do(req)

    if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Invalid status code: %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return getSearchListFromHTML(body, page)
}

func isLiveStream(contentSelection *goquery.Selection) bool {
    elem := contentSelection.Find("span[class*=\"yt-badge-live\"]").First().Nodes
    
    if elem == nil {
        if contentSelection.Find("span[class*=\"video-time\"]").First().Nodes == nil {
            fmt.Println("Video is livestream!!!")
            return true
        }
        
        fmt.Println("Video is livestream!!!")
    }
    
    return elem != nil
}

func getSearchListFromHTML(html []byte, page int) (*SearchList, error) {
    doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, err
	}

	searchList := &SearchList{}
    searchList.Page = page
    
    //fmt.Println(string(html))
    doc.Find("ol[class=\"item-section\"]").Children().Each(func(i int, contentSelection *goquery.Selection) {
        if contentSelection.Find("div[class*=\"yt-lockup-video\"]").First().Nodes != nil && !isLiveStream(contentSelection) {
            prviewInfo := PreviewInfo{}
            webPageUrl, bAttr := contentSelection.Find("h3").First().Find("a").First().Attr("href")
        
            if bAttr {
                if !strings.Contains(webPageUrl, webPageURLPre) {
                    webPageUrl = fmt.Sprintf("%s%s", webPageURLPre, webPageUrl)
                }
            
                prviewInfo.WebPageURL = webPageUrl
            
                if len(prviewInfo.WebPageURL) > 0 {
                    u, _ := url.ParseRequestURI(prviewInfo.WebPageURL)
                    videoID := u.Query().Get("v")
	    
                    if len(videoID) == 0 {
		                fmt.Println("Invalid youtube url, no video id")
	                } else {
                        prviewInfo.ID = videoID
                    }
                }
            }
            
            if len(prviewInfo.ID) > 0 {
                prviewInfo.Title = contentSelection.Find("h3").First().Find("a").First().Text()
                prviewInfo.Duration = parseDurationString(contentSelection.Find("span[class=\"video-time\"]").First().Text())
                prviewInfo.Author = contentSelection.Find("div[class=\"yt-lockup-byline\"]").First().Find("a").First().Text()
                prviewInfo.DatePublished = contentSelection.Find("div[class=\"yt-lockup-meta\"]").First().Find("li").First().Text()
        
                viewCount := contentSelection.Find("div[class=\"yt-lockup-meta\"]").First().Find("li").Get(1).FirstChild.Data
                regexNum := regexp.MustCompile("([^\\d])")
	            viewCount = regexNum.ReplaceAllString(viewCount, "")
                
                nViewCount, err := strconv.Atoi(viewCount)
        
                if err == nil {
                    prviewInfo.ViewCount = nViewCount
                }
        
                element := contentSelection.Find("div[class=\"yt-thumb video-thumb\"]").First().Find("img").First()
                thunmbURL, bSrc := element.Attr("src")
        
                if bSrc {
                    if strings.Contains(thunmbURL, ".gif") {
                        thunmbURL, _ = element.Attr("data-thumb")
                    }
                }

                prviewInfo.ThumbnailURL = thunmbURL
                searchList.PreInfoList = append(searchList.PreInfoList, prviewInfo)
            } else {
                fmt.Println("Video Id is empty, not support!!!")
            }
        } else {
            fmt.Println("Not video type!!!")
        }
	})
    
    return searchList, err
}

func getSuggestionListFromXML(content []byte) ([]string) {
    inputReader := strings.NewReader(string(content))
    decoder := xml.NewDecoder(inputReader)
    suggestList := []string{}
    
    for t, err := decoder.Token(); err == nil; t, err = decoder.Token() {
        switch token := t.(type) {
            case xml.StartElement:
                name := token.Name.Local
                if name == "suggestion" {
                    for _, attr := range token.Attr {
                        attrName := attr.Name.Local
                        if attrName == "data" {
                            attrValue := attr.Value
                            suggestList = append(suggestList, attrValue)
                            break
                        }
                    } 
                }
             default:
        } 
    }
    
    return suggestList
}

func parseDurationString(input string) time.Duration {
    splitInput := strings.Split(input, ":")
    days := "0"
    hours := "0"
    minutes := "0"
    var seconds string

    switch(len(splitInput)) {
       case 4:
            days = splitInput[0];
            hours = splitInput[1];
            minutes = splitInput[2];
            seconds = splitInput[3];
            break;
        case 3:
            hours = splitInput[0];
            minutes = splitInput[1];
            seconds = splitInput[2];
            break;
        case 2:
            minutes = splitInput[0];
            seconds = splitInput[1];
            break;
        case 1:
            seconds = splitInput[0];
            break;
        default:
            fmt.Println("Error duration string with unknown format: " + input)
    }
    
    nDay, _ := strconv.Atoi(days)
    nHour, _ := strconv.Atoi(hours)
    nMin, _ := strconv.Atoi(minutes)
    nSec, _ := strconv.Atoi(seconds)
    
    return time.Duration(nDay * 24 * 3600 + nHour * 3600 + nMin * 60 + nSec)
}