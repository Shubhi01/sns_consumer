package main

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"httpclient"
	"models"
	"smtp"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

var wg = sync.WaitGroup{}

var logger = logrus.NewEntry(logrus.New())

var client = httpclient.NewHTTPClient("", "", &http.Transport{}, logger)

func main() {

	wg.Add(3)
	go consumeFromTopic("HTTP", &wg)
	go consumeFromTopic("SMS", &wg)
	go consumeFromTopic("EMAIL", &wg)
	wg.Wait()
}

func consumeFromTopic(topic string, wg *sync.WaitGroup) {

	defer wg.Done()

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"10.46.143.17:9092"},
		GroupID: "consumer-group-id",
		//	Partition: 0,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	defer r.Close()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}

		logger.Infof("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		var action models.Action
		err = json.Unmarshal(m.Value, &action)
		if err != nil {
			logger.Error(err)
			continue
		}

		switch topic {
		case "HTTP":
			var actionSpec models.HTTPActionSpec
			err := json.Unmarshal([]byte(action.ActionSpec), &actionSpec)
			if err != nil {
				logger.Error(err)
				continue
			}
			_, _, err = client.Do(context.Background(), actionSpec.Method, actionSpec.URL, []byte(actionSpec.Body))
			if err != nil {
				logger.Error(err)
				continue
			}
		case "EMAIL":
			var actionSpec models.EmailActionSpec
			err := json.Unmarshal([]byte(action.ActionSpec), &actionSpec)
			if err != nil {
				logger.Error(err)
				continue
			}
			smtp.Send(&actionSpec)

		default:
			logger.Infof("unexpected topic: %s", topic)
			continue

		}

	}

}
