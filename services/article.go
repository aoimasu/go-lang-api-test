package services

import (
	"encoding/json"
	"time"
	"sort"
	"math"
	"io/ioutil"
	"strconv"
	"net/http"
	"github.com/gorilla/mux"
	"api/models"
	"api/utils"
)

// - MARK: Types

type ArticleService struct {
	Router *mux.Router
	Articles []models.Article
}

// - MARK: Initialization

func (s *ArticleService) Initialize() {
	s.Articles = []models.Article{}

	s.initializeRoutes()
}

// Initialize routes for Articles service
func (s *ArticleService) initializeRoutes() {
	s.Router.HandleFunc("/articles", s.storeArticles).Methods("POST")
	s.Router.HandleFunc("/articles/{id:[0-9]+}", s.returnSingleArticle).Methods("GET")
	s.Router.HandleFunc("/tags/{name}/{date:[0-9]{8}}", s.returnFilteredArticles).Methods("GET")
}

// - MARK: Handle functions

// Store submitted articles to service store
func (s *ArticleService) storeArticles(w http.ResponseWriter, r *http.Request) {  
	reqBody, _ := ioutil.ReadAll(r.Body)
	var articles []models.Article 
	json.Unmarshal(reqBody, &articles)

	if err := json.Unmarshal(reqBody, &articles); err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(&models.Error{Code: 400, Message: err.Error()})
		return
	}

	for i := 0; i < len(articles); i++ {
		article := &articles[i]
		article.Id = s.upsertArticle(*article)
	}

	json.NewEncoder(w).Encode(articles)
}

// Retrieve single article with the provided id
func (s *ArticleService) returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, err := strconv.ParseInt(vars["id"], 0, 8)

	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(&models.Error{Code: 400, Message: err.Error()})
		return
	}

	for _, article := range s.Articles {
			if article.Id == key {
					json.NewEncoder(w).Encode(article)
					return
			}
	}

	json.NewEncoder(w).Encode(&models.Article{})
}

// Retrieve Tags from name and date
func (s *ArticleService) returnFilteredArticles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	dateString := vars["date"]

	date, err := time.Parse(utils.YYYYMMDD_FORMAT, dateString)

	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(&models.Error{Code: 400, Message: err.Error()})
		return
	}

	articleIds := []string{}
	relatedTags := []string{}
	// sort by date desc
	articles := s.Articles
	sort.SliceStable(articles, func(i, j int) bool {
    return articles[i].Date > articles[j].Date
	})
	// filter articles
	for _, article := range articles {
			if utils.IndexOf(name, article.Tags) > -1 && article.Date == date.Format(utils.YYYY_MM_DD_FORMAT) {
				articleIds = append(articleIds, strconv.FormatInt(article.Id, 10))
				relatedTags = append(relatedTags, article.Tags...)
			}
	}

	tagCount := len(articleIds)
	maxNoOfArticle := int(math.Min(10, float64(tagCount)))
	// Return result
	json.NewEncoder(w).Encode(&models.Tag{
		Tag: name,
		Count: tagCount,
		Articles: articleIds[:maxNoOfArticle],
		RelatedTags: utils.Unique(relatedTags, []string{name}),
	})
}

// - MARK: helper functions

// Insert new article with generated Id to local store
func (s *ArticleService)insertNewArticle(article models.Article) int64 {
	if article.Id == 0 {
		article.Id = s.getNextId()
	}
	
	s.Articles = append(s.Articles, article)
	return article.Id
}

// Update/Insert Article into local store
func (s *ArticleService)upsertArticle(newArticle models.Article) int64 {
	found := false
	newId := newArticle.Id
	for index, article := range s.Articles {
		if newArticle.Id == article.Id {
				s.Articles[index] = newArticle
				found = true
				break
		}
	}

	if (!found) {
		newId = s.insertNewArticle(newArticle)
	}

	return newId
}



// Calculate next Id in local Articles store
func (s *ArticleService)getNextId() int64 {
	var nextId int64 = 0
	for _, article := range s.Articles {
		if nextId < article.Id {
			nextId = article.Id
		}
	}

	return nextId + 1
}

