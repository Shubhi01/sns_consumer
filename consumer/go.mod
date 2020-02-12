module consumer

go 1.12

require (
	github.com/segmentio/kafka-go v0.3.5 // indirect
	github.com/sirupsen/logrus v1.4.2 // indirect
	httpclient v0.0.0-00010101000000-000000000000 // indirect
	models v0.0.0-00010101000000-000000000000 // indirect
	smtp v0.0.0-00010101000000-000000000000 // indirect
)

replace httpclient => ../httpclient

replace models => ../models

replace smtp => ../smtp
