// main_test.go
package main_test
 
import (
	"os"
	"testing"
	"api"
	"api/models"

	"net/http"
	"net/http/httptest"
	"bytes"
	"encoding/json"
 )

var a main.App

func TestMain(m *testing.M) {
	a.Initialize()

	code := m.Run()
	os.Exit(code)
}

func Reset() {
    a.Initialize()
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
			t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func addSampleData(t *testing.T) ([]models.Article) {
	var jsonStr = []byte(`
	[
		{
			"title": "title 1",
			"date" : "2016-09-22",
			"body" : "body 1",
			"tags" : ["health", "test", "test2"]
		},
		{
			"id": "2",
			"title": "title 2",
			"date" : "2016-09-22",
			"body" : "body 2",
			"tags" : ["health", "fitness", "science"]
		},
		{
			"title": "title 3",
			"date" : "2016-09-23",
			"body" : "body 3",
			"tags" : ["health", "test3", "test4"]
		}
	]`)
	req, _ := http.NewRequest("POST", "/articles", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var articles []models.Article 
	json.Unmarshal(response.Body.Bytes(), &articles)

	return articles
}

// Test cases 

func TestInvalidArticle(t *testing.T) {
	Reset()
	req, _ := http.NewRequest("GET", "/articles/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "{}\n" {
			t.Errorf("Expected an empty object. Got `%s`", body)
	}

	// string id
	req, _ = http.NewRequest("GET", "/articles/string", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestPostArticles(t *testing.T) {
	Reset()
	articles := addSampleData(t)

	if len(articles) != 3 {
		t.Errorf("Expected number of article to be 3. Got '%v'", len(articles))
	}
	if articles[0].Id != 1 {
			t.Errorf("Expected id of first article to be 1. Got '%v'", articles[0].Id)
	}
	if articles[1].Id != 2 {
			t.Errorf("Expected id of second article to be 2. Got '%v'", articles[1].Id)
	}
	if articles[2].Id != 3 {
			t.Errorf("Expected id of second article to be 3. Got '%v'", articles[2].Id)
	}

	// Try query the inserted article
	req, _ := http.NewRequest("GET", "/articles/2", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	var article models.Article 
	json.Unmarshal(response.Body.Bytes(), &article)

	if article.Id != 2 {
		t.Errorf("Expected the returned article has Id 2. Got `%v`", article.Id)
	}
	if article.Title != "title 2" {
			t.Errorf("Expected the returned article has title `title 2`. Got `%s`", article.Title)
	}
}

func TestPostEmptyArticles(t *testing.T) {
	Reset()
	var jsonStr = []byte(`[]`)
	req, _ := http.NewRequest("POST", "/articles", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var articles []models.Article 
	json.Unmarshal(response.Body.Bytes(), &articles)

	if len(articles) != 0 {
		t.Errorf("Expected number of article to be 0. Got '%v'", len(articles))
	}
}


func TestGetTags(t *testing.T) {
	Reset()
	_ = addSampleData(t)

	// Try query the inserted article
	req, _ := http.NewRequest("GET", "/tags/health/20160922", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	var tags models.Tag 
	json.Unmarshal(response.Body.Bytes(), &tags)

	if tags.Tag != "health" {
		t.Errorf("Expected the returned object has tag `health`. Got `%s`", tags.Tag)
	}
	if tags.Count != 2 {
			t.Errorf("Expected the returned object has count `2`. Got `%v`", tags.Count)
	}
	if len(tags.Articles) != 2 {
			t.Errorf("Expected the returned object has articles len `2`. Got `%v`", len(tags.Articles))
	}
	if len(tags.RelatedTags) != 4 {
			t.Errorf("Expected the returned object has relatedTags len `4`. Got `%v`", len(tags.RelatedTags))
	}
}

func TestGetEmptyTags(t *testing.T) {
	Reset()
	_ = addSampleData(t)

	// Try query the inserted article
	req, _ := http.NewRequest("GET", "/tags/invalid/20160922", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	var tags models.Tag 
	json.Unmarshal(response.Body.Bytes(), &tags)

	if tags.Tag != "invalid" {
		t.Errorf("Expected the returned object has tag `invalid`. Got `%s`", tags.Tag)
	}
	if tags.Count != 0 {
			t.Errorf("Expected the returned object has count `0`. Got `%v`", tags.Count)
	}
	if len(tags.Articles) != 0 {
			t.Errorf("Expected the returned object has articles len `0`. Got `%v`", len(tags.Articles))
	}
	if len(tags.RelatedTags) != 0 {
			t.Errorf("Expected the returned object has relatedTags len `0`. Got `%v`", len(tags.RelatedTags))
	}

	// string date
	req, _ = http.NewRequest("GET", "/tags/invalid/invalid", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}