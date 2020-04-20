package router

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
)

const (
	articlesURL         = "https://storage.googleapis.com/aller-structure-task/articles.json"
	contentMarketingURL = "https://storage.googleapis.com/aller-structure-task/contentmarketing.json"
)

var (
	articles         map[string]interface{}
	contentMarketing map[string]interface{}
	resultContent    []interface{}
)

//getArticlesHandler assembling resultContent map
func getArticlesHandler(c *gin.Context) {

	var wg sync.WaitGroup

	wg.Add(2)
	go getArticles(&wg)
	go getContentMarketing(&wg)

	wg.Wait()

	err := resultMapper()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, resultContent)
}

// ------------------------------------------------------------------------------------------------------------------ //

//getArticles get Articles from the third party source and populate articles map
func getArticles(wg *sync.WaitGroup) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", articlesURL, nil)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	defer func() {
		resp.Body.Close()
		wg.Done()
	}()

	err = json.NewDecoder(resp.Body).Decode(&articles)
	if err != nil {
		log.Println(err)
	}
}

//getContentMarketing get Content Marketing from the third party source and populate contentMarketing map
func getContentMarketing(wg *sync.WaitGroup) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", contentMarketingURL, nil)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	defer func() {
		resp.Body.Close()
		wg.Done()
	}()

	err = json.NewDecoder(resp.Body).Decode(&contentMarketing)
	if err != nil {
		log.Println(err)
	}
}

//resultMapper mapping resultContent according the pattern
func resultMapper() error {

	if respArticles, ok := articles["response"].(map[string]interface{}); ok {
		if articlesItems, ok := respArticles["items"].([]interface{}); ok {
			if respContentMarketing, ok := contentMarketing["response"].(map[string]interface{}); ok {
				if contentMarketingItems, ok := respContentMarketing["items"].([]interface{}); ok {
					resultContentCap := len(articlesItems) + 1 + len(contentMarketingItems) + len(articlesItems)/6
					for i := 0; i < resultContentCap; i++ {
						if i != 0 && i > 11 && i%6 == 0 || i == 6 {
							if len(contentMarketingItems) > 0 {
								// Shift
								itemCM := contentMarketingItems[0]
								contentMarketingItems = contentMarketingItems[1:]
								articlesItems = insert(articlesItems, itemCM, i-1)
							} else {
								articlesItems = insert(articlesItems, map[string]string{"type": "Ad"}, i-1)
							}
						}
					}
					resultContent = articlesItems
				} else {
					log.Println("Invalid Content Marketing items format")
					return errors.New("invalid Content Marketing items format")
				}
			} else {
				log.Println("Invalid Content Marketing response format")
				return errors.New("invalid Content Marketing response format")
			}
		} else {
			log.Println("invalid Articles items format")
			return errors.New("invalid Articles items format")
		}
	} else {
		log.Println("invalid Articles response format")
		return errors.New("invalid Articles response format")
	}
	return nil
}

// insert value in a slice at given index
func insert(a []interface{}, c interface{}, i int) []interface{} {
	return append(a[:i], append([]interface{}{c}, a[i:]...)...)
}
