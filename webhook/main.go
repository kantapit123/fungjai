package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/v7/linebot"
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
	Hash_userId        string             `json:"hash_userId"`
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
	// ------------godotenv
	// --loads .env file >> work on local only if you want to run on docker comment it!
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	// ------------

	// ------------Initialize linebot client
	client := &http.Client{}
	bot, err := linebot.New(os.Getenv("CHANNEL_SECRET"), os.Getenv("CHANNEL_ACCESS_TOKEN"), linebot.WithHTTPClient(client))
	if err != nil {
		log.Fatal("Line bot client ERROR: ", err)
	}

	// ------------Initialize a new MinIO client object
	// ctx := context.Background()
	// endpoint := os.Getenv("MINIO_ENDPOINT")                 //os.Getenv("MINIO_ENDPOINT")
	// accessKeyID := os.Getenv("MINIO_ACCESS_KEY_ID")         //os.Getenv("MINIO_ACCESS_KEY_ID")
	// secretAccessKey := os.Getenv("MINIO_SECRET_ACCESS_KEY") //os.Getenv("MINIO_SECRET_ACCESS_KEY")
	// useSSL := false                                         // Change to true if you want to use SSL
	// minioClient, err := minio.New(endpoint, &minio.Options{
	// 	Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	// 	Secure: useSSL,
	// })
	// if err != nil {
	// 	log.Println("Error to Initialize a new MinIO client")
	// 	fmt.Println(err)
	// }

	// ------------Make a new bucket.
	// bucketName := os.Getenv("MINIO_BUCKET_NAME")     //os.Getenv("MINIO_BUCKET_NAME")
	// location := os.Getenv("MINIO_LOCATION")          //os.Getenv("MINIO_LOCATION")
	// folderName := os.Getenv("MINIO_FOLDER_NAME")     //os.Getenv("MINIO_FOLDER_NAME")
	// responseType := os.Getenv("MINIO_RESPONSE_TYPE") //os.Getenv("MINIO_RESPONSE_TYPE")

	// err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	// if err != nil {
	// 	// ------------Check to see if we already own this bucket (which happens if you run this twice)
	// 	exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
	// 	if errBucketExists == nil && exists {
	// 		log.Printf("We already own %s\n", bucketName)
	// 	} else {
	// 		log.Println("Error to Create Bucket")
	// 		log.Fatalln(err)
	// 	}
	// } else {
	// 	log.Printf("Successfully created %s\n", bucketName)
	// }
	log.Print("--------------------> MinIO Start succeed")

	// ------------Initialize kafka producer
	err = InitKafka()
	if err != nil {
		log.Fatal("Kafka producer ERROR: ", err)
	}
	log.Println("Kafka Init succeed")
	// ------------

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
							"url": "https://img.pic.in.th/46ce813d404f5489e.png",
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
							"url": "https://img.pic.in.th/2591a2942c9f8ac3c.png",
							"size": "full",
							"aspectMode": "cover",
							"aspectRatio": "21:16",
							"gravity": "center",
							"action": {
								"type": "postback",
								"label": "mood2",
								"data": "question=1&value=2",
								"displayText": "2"
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
							"url": "https://img.pic.in.th/36de74b811eb44ff4.png",
							"gravity": "center",
							"size": "full",
							"aspectRatio": "21:16",
							"aspectMode": "cover",
							"action": {
								"type": "postback",
								"label": "mood3",
								"data": "question=1&value=3",
								"displayText": "3"
							},
							"position": "relative"
							},
							{
							"type": "image",
							"url": "https://img.pic.in.th/1584dd766dcbd8f0c.png",
							"gravity": "center",
							"size": "full",
							"aspectRatio": "21:16",
							"aspectMode": "cover",
							"action": {
								"type": "postback",
								"label": "mood1",
								"data": "question=1&value=1",
								"displayText": "1"
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
	// ------------Initilaze Echo web servers
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	e.POST("/webhook", func(c echo.Context) error {
		events_bot, err := bot.ParseRequest(c.Request())
		log.Println(c.Request().Header.Get("x-line-signature"))

		if err != nil {
			log.Fatal("Event bot ERROR: ", err)
		}
		// ------------kafka sent data method
		topics := os.Getenv("KAFKA_TOPIC_NAME")
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
		// -------------
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
					_, err := json.Marshal(&RawLineMessage{
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
					if err != nil {
						log.Panicln(err)
					}

					messageResponse := fmt.Sprintf("ข้อความของคุณคือ💌 " + message.Text)
					replyMessage(messageResponse, bot, event.ReplyToken)
					if message.Text == "/mood" {
						log.Print("mood")
						PushFlexMessage(json_flex1, bot, event.Source.UserID)
					}
				}
			}
			if event.Type == linebot.EventTypePostback {
				message := event
				PostbackLineMessage, err := JSONMarshal(&PostbackLineMessage{
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
					ReplyToken:  event.ReplyToken,
					Mode:        string(event.Mode),
					Hash_userId: string(strconv.FormatUint(uint64(hash(event.Source.UserID)), 10)),
				})

				display_string := string(event.Postback.Data[len(event.Postback.Data)-1:])
				// messageResponse := fmt.Sprintf("ขอบคุณที่แชร์ความรู้สึกกับเรานะ ^^")
				// if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(messageResponse)).Do(); err != nil {
				// 	log.Print(err)
				// }
				var display_mood_string string
				var messageResponse string
				if display_string == "1" {
					display_mood_string = "แย่มาก"
					messageResponse = fmt.Sprintf("ขอบคุณที่แชร์ความรู้สึกกับเรานะ💓\nคำตอบของคุณคือ " + display_mood_string + "\nมีอะไรบอกเราได้นะเป็นห่วงนะ 🥺")
				} else if display_string == "2" {
					display_mood_string = "แย่"
					messageResponse = fmt.Sprintf("ขอบคุณที่แชร์ความรู้สึกกับเรานะ💓\nคำตอบของคุณคือ " + display_mood_string + "\nไม่เป็นไรนะสู้ ๆ 🤗")
				} else if display_string == "3" {
					display_mood_string = "ดี"
					messageResponse = fmt.Sprintf("ขอบคุณที่แชร์ความรู้สึกกับเรานะ💓\nคำตอบของคุณคือ " + display_mood_string + "\nเยี่ยมไปเลย พรุ่งนี้ต้องดีกว่านี้แน่นอน ✌🏻")
				} else if display_string == "4" {
					display_mood_string = "ดีมาก"
					messageResponse = fmt.Sprintf("ขอบคุณที่แชร์ความรู้สึกกับเรานะ💓\nคำตอบของคุณคือ " + display_mood_string + "\nขอให้ดีแบบนี้ต่อไปน้าา ✨")
				}

				replyMessage(messageResponse, bot, event.ReplyToken)

				producerErr := Produce(topics, string(PostbackLineMessage))
				if producerErr != nil {
					log.Print(err)
					panic("errorkafka")
				} else {
					messageResponse := fmt.Sprintf("Produced [%s] successfully", message.Postback.Data)
					log.Print(messageResponse)
				}

				// -------------- Put data direct to Minio
				// currentTime := time.Now().UTC()
				// date := fmt.Sprintf("%d-%02d-%02d", currentTime.Year(), currentTime.Month(), currentTime.Day())
				// timestamp_string := strconv.Itoa(int(event.Timestamp.UnixNano()))
				// uid_hash_string := strconv.FormatUint(uint64(hash(event.Source.UserID)), 10)
				// log.Println("Add " + timestamp_string + "-" + uid_hash_string + ".json to MinIO Bucket ->> " + bucketName)

				// objectName := fmt.Sprintf(folderName + "/" + responseType + "/" + date + "/" + timestamp_string + "-" + uid_hash_string + ".json")
				// _, err = minioClient.PutObject(context.Background(), bucketName, objectName, bytes.NewReader(PostbackLineMessage), int64(len(PostbackLineMessage)), minio.PutObjectOptions{})
				// if err != nil {
				// 	log.Println(err)
				// }
				// log.Println("JSON data put into MinIO successfully")
				// ---------------------
			}

		}

		log.Println("-->> message success")
		return c.String(http.StatusOK, "ok!!")

	})

	// ------------Broadcast Message everytime go run main.go
	// BroadcastFlexMessage(json_flex1, bot)
	// ------------Schedule main.go
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
		linebot.NewFlexMessage("What are you feeling today.", container),
	).Do(); err != nil {
		return err
	}
	return nil
}

func PushFlexMessage(FlexJson []byte, bot *linebot.Client, UserID string) error {

	container, err := linebot.UnmarshalFlexMessageJSON(FlexJson)
	if err != nil {
		log.Printf("failed to unmarshal Flex Message: %v", err)
	}

	if _, err = bot.PushMessage(UserID, linebot.NewFlexMessage("flex reply /mood", container)).Do(); err != nil {
		log.Println("Error to Reply")
		log.Print(err)
	}
	return nil
}

func replyMessage(Message string, bot *linebot.Client, replyToken string) error {

	_, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(Message)).Do()
	if err != nil {
		log.Print(err)
	}
	return nil
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
	//Use UTC timezone with convert loval time to UTC
	s := cron.New(cron.WithLocation(time.UTC))

	// Broadcast At minute 30 past hour 8, 11, and 16 on every day-of-week from Monday through Friday.
	s.AddFunc("30 8 * * 1-5", func() {
		BroadcastFlexMessage(jsonData, botClient)
	})
	// s.AddFunc("@hourly", func() {
	// 	BroadcastFlexMessage(jsonData, botClient)
	// })
	log.Print(s.Entries())

	s.Start()
}

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

var producer *kafka.Producer

func InitKafka() error {
	var err error
	producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": os.Getenv("BOOTSTRAP_SERVERS")})
	return err
}

func Produce(topics string, message string) error {
	deliveryChan := make(chan kafka.Event)
	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topics, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, deliveryChan)
	if err != nil {
		fmt.Printf("Produce failed: %v\n", err)
	}
	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	close(deliveryChan)
	return m.TopicPartition.Error
}
