package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/robfig/cron/v3"
)

type PostbackLineMessage struct {
	Type               string             `json:"type"`
	Postback           Postback           `json:"postback"`
	WebhookEventID     string             `json:"webhookEventId"`
	RawDeliveryContext RawDeliveryContext `json:"deliveryContext"`
	Timestamp          int64              `json:"timestamp"`
	RawSource          RawSource          `json:"source"`
	ReplyToken         string             `json:"replyToken"`
	Mode               string             `json:"mode"`
}
type Postback struct {
	Data string `json:"data"`
}

type RawLineMessage struct {
	Type               string             `json:"type"`
	RawMessage         RawMessage         `json:"message"`
	WebhookEventID     string             `json:"webhookEventId"`
	RawDeliveryContext RawDeliveryContext `json:"deliveryContext"`
	Timestamp          int64              `json:"timestamp"`
	RawSource          RawSource          `json:"source"`
	ReplyToken         string             `json:"replyToken"`
	Mode               string             `json:"mode"`
}
type RawMessage struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Text string `json:"text"`
}
type RawDeliveryContext struct {
	IsRedelivery bool `json:"isRedelivery"`
}
type RawSource struct {
	Type   string `json:"type"`
	UserID string `json:"userId"`
}

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

	// Initialize a new MinIO client object
	ctx := context.Background()
	endpoint := "localhost:9000"
	accessKeyID := "lmenrLntC9lTc9zf"
	secretAccessKey := "z2t4IpmIsWL15ArBIvCbqcQbtKCjLZsX"
	useSSL := false // Change to true if you want to use SSL
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		fmt.Println(err)
	}

	// Make a new bucket.
	bucketName := "fungjai"
	location := "ap-southeast-1"
	folderName := "responses"
	responseType := "mood"

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}
	log.Print("--------------------> MinIO Start succeed")
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
		for _, event := range events_bot {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					// messageJson, err := json.Marshal(&producedMessage{
					// 	UserID:    event.Source.UserID,
					// 	Timestamp: event.Timestamp.UnixNano(),
					// 	MessageID: message.ID,
					// 	Message:   message.Text,
					// })
					RawmessageJson, err := json.Marshal(&RawLineMessage{
						Type: string(event.Type),
						RawMessage: RawMessage{
							Type: string(linebot.MessageTypeText),
							ID:   message.ID,
							Text: message.Text,
						},
						WebhookEventID: event.WebhookEventID,
						RawDeliveryContext: RawDeliveryContext{
							IsRedelivery: event.DeliveryContext.IsRedelivery,
						},
						Timestamp: event.Timestamp.UnixNano(),
						RawSource: RawSource{
							Type:   string(event.Source.Type),
							UserID: event.Source.UserID,
						},
						ReplyToken: event.ReplyToken,
						Mode:       string(event.Mode),
					})

					messageResponse := fmt.Sprintf("messageResponse is " + message.Text)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(messageResponse)).Do(); err != nil {
						log.Print(err)
					}
					currentTime := time.Now().UTC()
					date := fmt.Sprintf("%d-%02d-%02d", currentTime.Year(), currentTime.Month(), currentTime.Day())
					timestamp_string := strconv.Itoa(int(event.Timestamp.UnixNano()))
					uid_hash_string := strconv.FormatUint(uint64(hash(event.Source.UserID)), 10)
					log.Println("Add " + timestamp_string + "-" + uid_hash_string + ".json to MinIO Bucket ->> " + bucketName)

					objectName := fmt.Sprintf(folderName + "/" + responseType + "/" + date + "/" + timestamp_string + "-" + uid_hash_string + ".json")
					_, err = minioClient.PutObject(context.Background(), bucketName, objectName, bytes.NewReader(RawmessageJson), int64(len(RawmessageJson)), minio.PutObjectOptions{})
					if err != nil {
						log.Println(err)
					}
					log.Println("JSON data put into MinIO successfully")

					// producerErr := producer.Produce(topics, string(messageJson))
					// if producerErr != nil {
					// 	log.Print(err)
					// }

					//messageResponse := fmt.Sprintf("Produced [%s] successfully", message.Text)
				}
			}
			if event.Type == linebot.EventTypePostback {
				message := event
				PostbackLineMessage, err := json.Marshal(&PostbackLineMessage{
					Type: string(message.Type),
					Postback: Postback{
						Data: event.Postback.Data,
					},
					WebhookEventID:     event.WebhookEventID,
					RawDeliveryContext: RawDeliveryContext(event.DeliveryContext),
					Timestamp:          event.Timestamp.UnixNano(),
					RawSource: RawSource{
						Type:   string(event.Source.Type),
						UserID: event.Source.UserID,
					},
					ReplyToken: event.ReplyToken,
					Mode:       string(event.Mode),
				})
				messageResponse := fmt.Sprintf("The Data postback from flex message is " + event.Postback.Data)
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(messageResponse)).Do(); err != nil {
					log.Print(err)
				}
				currentTime := time.Now().UTC()
				date := fmt.Sprintf("%d-%02d-%02d", currentTime.Year(), currentTime.Month(), currentTime.Day())
				timestamp_string := strconv.Itoa(int(event.Timestamp.UnixNano()))
				uid_hash_string := strconv.FormatUint(uint64(hash(event.Source.UserID)), 10)
				log.Println("Add " + timestamp_string + "-" + uid_hash_string + ".json to MinIO Bucket ->> " + bucketName)

				objectName := fmt.Sprintf(folderName + "/" + responseType + "/" + date + "/" + timestamp_string + "-" + uid_hash_string + ".json")
				_, err = minioClient.PutObject(context.Background(), bucketName, objectName, bytes.NewReader(PostbackLineMessage), int64(len(PostbackLineMessage)), minio.PutObjectOptions{})
				if err != nil {
					log.Println(err)
				}
				log.Println("JSON data put into MinIO successfully")
			}

		}

		log.Println("-->> message success")
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
								"type": "postback",
								"label": "mood4",
								"data": "question=1&value=4",
								"displayText": "4"
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
								"label": "mood2",
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
								"label": "mood3",
								"text": "3"
							},
							"position": "relative"
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
								"label": "mood1",
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
	BroadcastFlexMessage(json_flex1, bot)
	runCronJobs(json_flex1, bot)
	e.Logger.Fatal(e.Start(":1323"))
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
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

func runCronJobs(jsonData []byte, botClient *linebot.Client) {
	s := cron.New()

	// Broadcast At minute 30 past hour 8, 11, and 16 on every day-of-week from Monday through Friday.
	s.AddFunc("30 8,11,16 * * 1-5", func() {
		BroadcastFlexMessage(jsonData, botClient)
	})

	s.Start()
}
