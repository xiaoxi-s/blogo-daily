package handlers

import (
	"blogoproducer/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

var newsRefreshInterval = 1 * time.Second
var lastUpdatedTime time.Time
var channelAmqp *amqp.Channel

func init() {
	amqpConnection, err := amqp.Dial(os.Getenv("RABBITMQ_URI"))
	if err != nil {
		log.Fatal(err)
	}
	lastUpdatedTime = time.Now()
	channelAmqp, _ = amqpConnection.Channel()
}

func ParseHandler(c *gin.Context) {
	if time.Now().Before(lastUpdatedTime.Add(newsRefreshInterval)) {
		c.JSON(http.StatusNotModified, gin.H{"message": "do not update daily"})
		return
	}
	var request models.RssFeedRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	data, _ := json.Marshal(request)
	err := channelAmqp.Publish("", os.Getenv("RABBITMQ_QUEUE"), false, false, amqp.Publishing{ContentType: "application/json", Body: []byte(data)})
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Error while publishing to RabbitMQ"})
		return
	}
	c.JSON(http.StatusOK, map[string]string{"message": "success"})
}
