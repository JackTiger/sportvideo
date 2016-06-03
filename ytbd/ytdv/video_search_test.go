package ytdv

import (
	"testing"
    "fmt"
)

func TestGetSuggestionListFromQuery(t *testing.T) {
    suggestList, err := GetSuggestionListFromQuery("NSQ")
    fmt.Println(suggestList)
    
    if err != nil {
        t.Fatal(err)
        fmt.Println("GetSuggestionListFromQuery Error")
    }
}
