package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/robfig/cron/v3"
)

type LineMessage struct {
	Destination string `json:"destination"`
	Events      []struct {
		ReplyToken string `json:"replyToken"`
		Type       string `json:"type"`
		Timestamp  int64  `json:"timestamp"`
		Source     struct {
			Type   string `json:"type"`
			UserID string `json:"userId"`
		} `json:"source"`
		Message struct {
			ID   string `json:"id"`
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"message"`
	} `json:"events"`
}

type ReplyMessageGo struct {
	ReplyToken string `json:"replyToken"`
	Messages   []Text `json:"messages"`
}

type BroadcastMessageTexts struct {
	Messages []Text `json:"messages"`
}

type Text struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ProFile struct {
	UserID        string `json:"userId"`
	DisplayName   string `json:"displayName"`
	PictureURL    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
}

type producedMessage struct {
	Timestamp int64  `json:"timestamp"`
	UserID    string `json:"userId"`
	MessageID string `json:"msgId"`
	Message   string `json:"message"`
}

func main() {

	// loads .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize linebot client
	client := &http.Client{}
	bot, err := linebot.New(os.Getenv("CHANNEL_SECRET"), os.Getenv("CHANNEL_ACCESS_TOKEN"), linebot.WithHTTPClient(client))
	if err != nil {
		log.Fatal("Line bot client ERROR: ", err)
	}

	//minio config
	ctx := context.Background()
	//endpoint := "5817-2405-9800-b860-b117-e4d4-d3f9-7321-59a.ap.ngrok.io"
	host := "localhost:9000"
	accessKeyID := "lmenrLntC9lTc9zf"
	secretAccessKey := "z2t4IpmIsWL15ArBIvCbqcQbtKCjLZsX"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("----------------------------------> minioClient is now setup")
	// log.Printf("%#v\n", minioClient) // minioClient is now setup

	// Make a new bucket called test1.
	bucketName := "test1kantapit2"
	location := "ap-southeast-1"

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Printf("Fail to Create %s\n", bucketName)
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	// Initialize kafka producer
	// err = producer.InitKafka()
	// if err != nil {
	// 	log.Fatal("Kafka producer ERROR: ", err)
	// }

	// Initilaze Echo web servers
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	e.POST("/webhook", func(c echo.Context) error {
		events_bot, err := bot.ParseRequest(c.Request())
		if err != nil {
			log.Fatal("Event bot ERROR: ", err)
		}
		// topics := "user-messages"

		// for _, event := range events_bot {
		// 	if event.Type == linebot.EventTypeMessage {
		// 		switch message := event.Message.(type) {
		// 		case *linebot.TextMessage:
		// 			messageJson, _ := json.Marshal(&producedMessage{
		// 				UserID:    event.Source.UserID,
		// 				Timestamp: event.Timestamp.Unix(),
		// 				MessageID: message.ID,
		// 				Message:   message.Text,
		// 			})
		// 			producerErr := producer.Produce(topics, string(messageJson))
		// 			if producerErr != nil {
		// 				log.Print(err)
		// 				panic("errorkafka")
		// 			} else {
		// 				messageResponse := fmt.Sprintf("Produced [%s] successfully", message.Text)
		// 				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(messageResponse)).Do(); err != nil {
		// 					log.Print(err)
		// 				}
		// 			}
		// 		}
		// 	}
		// }

		// topics := "user-messages"
		for _, event := range events_bot {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					// messageJson, _ := json.Marshal(&producedMessage{
					// 	UserID:    event.Source.UserID,
					// 	Timestamp: event.Timestamp.Unix(),
					// 	MessageID: message.ID,
					// 	Message:   message.Text,
					// })
					// fmt.Println(messageJson)
					messageResponse := "messageResponse is " + message.Text
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(messageResponse)).Do(); err != nil {
						log.Print("webhook error")
						log.Print(err)
					}
					// producerErr := producer.Produce(topics, string(messageJson))
					// if producerErr != nil {
					// 	log.Print(err)
					// }

					//messageResponse := fmt.Sprintf("Produced [%s] successfully", message.Text)
				}
			}
		}

		log.Println("%% message success")
		return c.String(http.StatusOK, "ok!!")

	})

	//Broadcast Message everytime go run main.go
	// text := Text{
	// 	Type: "text",
	// 	Text: "Server is ready !!!",
	// }

	// message := BroadcastMessageTexts{
	// 	Messages: []Text{
	// 		text,
	// 	},
	// }

	// BroadcastMessageLine(message)
	json_flex1 := []byte(`
		{
			"type": "bubble",
			"header": {
				"type": "box",
				"layout": "vertical",
				"contents": [
				{
					"type": "image",
					"url": "https://sv1.picz.in.th/images/2023/02/24/LMiinI.png",
					"size": "full",
					"aspectMode": "cover",
					"position": "relative"
				}
				],
				"background": {
				"type": "linearGradient",
				"angle": "0deg",
				"startColor": "#E8E0D5",
				"endColor": "#ffffff"
				},
				"position": "relative",
				"paddingAll": "0%"
			},
			"body": {
				"type": "box",
				"layout": "vertical",
				"contents": [
				{
					"type": "box",
					"layout": "horizontal",
					"contents": [
					{
						"type": "box",
						"layout": "vertical",
						"contents": [
						{
							"type": "image",
							"url": "https://sv1.picz.in.th/images/2023/02/24/LMigQZ.png",
							"size": "full",
							"aspectMode": "cover",
							"aspectRatio": "21:16",
							"gravity": "center",
							"action": {
							"type": "message",
							"label": "action",
							"text": "4"
							}
						},
						{
							"type": "image",
							"url": "https://sv1.picz.in.th/images/2023/02/24/LMiTmz.png",
							"size": "full",
							"aspectMode": "cover",
							"aspectRatio": "21:16",
							"gravity": "center",
							"action": {
							"type": "message",
							"label": "action",
							"text": "2"
							}
						}
						],
						"flex": 1
					},
					{
						"type": "box",
						"layout": "vertical",
						"contents": [
						{
							"type": "image",
							"url": "https://sv1.picz.in.th/images/2023/02/24/LMxaAa.png",
							"gravity": "center",
							"size": "full",
							"aspectRatio": "21:16",
							"aspectMode": "cover",
							"action": {
							"type": "message",
							"label": "action",
							"text": "3"
							}
						},
						{
							"type": "image",
							"url": "https://sv1.picz.in.th/images/2023/02/24/LMi1bR.png",
							"gravity": "center",
							"size": "full",
							"aspectRatio": "21:16",
							"aspectMode": "cover",
							"action": {
							"type": "message",
							"label": "action",
							"text": "1"
							}
						}
						]
					}
					]
				}
				],
				"paddingAll": "0px"
			}
			}
		`)

	//run the code under every day
	runCronJobs(json_flex1, bot)
	e.Logger.Fatal(e.Start(":1323"))
}

func BroadcastFlexMessage(jsonData []byte, bot *linebot.Client) error {

	container, err := linebot.UnmarshalFlexMessageJSON(jsonData)
	if err != nil {
		log.Printf("failed to unmarshal Flex Message: %v", err)
	}

	if _, err = bot.BroadcastMessage(
		linebot.NewFlexMessage("You have new flex message", container),
	).Do(); err != nil {
		return err
	}
	return nil
}

func replyMessageLine(Message ReplyMessageGo) error {
	value, _ := json.Marshal(Message)

	url := "https://api.line.me/v2/bot/message/reply"

	var jsonStr = []byte(value)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("CHANNEL_ACCESS_TOKEN"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))

	return err
}

func BroadcastMessageLine(Message BroadcastMessageTexts) error {
	value, _ := json.Marshal(Message)

	url := "https://api.line.me/v2/bot/message/broadcast"

	var jsonStr = []byte(value)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("CHANNEL_ACCESS_TOKEN"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))

	return err
}

func getProfile(userId string) string {

	url := "https://api.line.me/v2/bot/profile/" + userId

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("CHANNEL_ACCESS_TOKEN"))

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var profile ProFile
	if err := json.Unmarshal(body, &profile); err != nil {
		log.Println("%% err \n")
	}
	log.Println(profile.DisplayName)
	return profile.DisplayName

}

func hello(name string) {
	message := fmt.Sprintf("Hi, %v", name)
	fmt.Println(message)
}

func runCronJobs(jsonData []byte, botClient *linebot.Client) {
	s := cron.New()

	// Broadcast At minute 30 past hour 8, 11, and 16 on every day-of-week from Monday through Friday.
	s.AddFunc("30 8,11,16 * * 1-5", func() {
		BroadcastFlexMessage(jsonData, botClient)
	})

	s.Start()
}
