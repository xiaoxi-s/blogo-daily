package handlers

import (
	"blogoproducer/models"
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var ctx context.Context
var newsRefreshInterval = 1 * time.Hour
var lastUpdatedTime time.Time

func init() {
	ctx = context.Background()
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
}

func Dispatch(url string) int {
	if strings.Contains(url, "google") {
		return 0
	}
	return -1
}

func GetFeedFromGoogleNews(url string) ([]models.Entry, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 ( Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	byteValue, _ := ioutil.ReadAll(resp.Body)
	var googleNewsFeed models.GoogleNewsFeed
	err = xml.Unmarshal(byteValue, &googleNewsFeed)
	if err != nil {
		return nil, err
	}
	fmt.Print("rsstag: ", googleNewsFeed.XMLName)
	fmt.Println("channel: ", googleNewsFeed.Items[0])
	fmt.Println("Size of News entires: ", len(googleNewsFeed.Items))

	newsEntries := make([]models.Entry, len(googleNewsFeed.Items))
	for ind, googleNewsEntry := range googleNewsFeed.Items {
		newsEntries[ind].Title = googleNewsEntry.Title
		newsEntries[ind].Link = googleNewsEntry.Link
		newsEntries[ind].Description = googleNewsEntry.Description
	}

	return newsEntries, err
}

func GetFeedEntries(url string) ([]models.Entry, error) {
	dispatch := Dispatch(url)

	var entries []models.Entry
	var err error
	if dispatch == 0 {
		entries, err = GetFeedFromGoogleNews(url)
	}

	return entries, err
}

func ParseHandler(c *gin.Context) {

	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("news")
	if time.Now().Before(lastUpdatedTime.Add(newsRefreshInterval)) {
		cur, err := collection.Find(ctx, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		entries := make([]models.Entry, 0)
		var entry models.Entry
		for cur.Next(ctx) {
			cur.Decode(&entry)
			entries = append(entries, entry)
		}
		c.JSON(http.StatusOK, entries)
	}

	var request models.RssFeedRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	entries, err := GetFeedEntries(request.Url)
	var maxNumOfNews int

	if len(entries) > 5 {
		maxNumOfNews = 5
	} else {
		maxNumOfNews = len(entries)
	}

	for _, entry := range entries[0:maxNumOfNews] {
		collection.InsertOne(ctx, bson.M{
			"title":       entry.Title,
			"description": entry.Description,
			"link":        entry.Link,
		})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Error while parsing the rss feed"})
		return
	}
	c.JSON(http.StatusOK, entries)
}
