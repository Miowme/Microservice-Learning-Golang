package rabbitmq

import (
	"encoding/json"
	"log"
	"time"

	"payment-service/internal/model" 

	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

type PaymentMessage struct {
	PaymentID string  `json:"payment_id"`
	UserID    string  `json:"user_id"`
	MerchantID string  `json:"merchant_id"`
	Amount    float64 `json:"amount"`
	CreatedAt string  `json:"created_at"`
}

func StartConsumer(amqpURL string, db *gorm.DB) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"payment.queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Println("[Payment Service] Waiting for messages...")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var payMsg PaymentMessage
			if err := json.Unmarshal(d.Body, &payMsg); err != nil {
				log.Printf("Error decoding message: %v", err)
				d.Nack(false, false)
				continue
			}

			createdAt, err := time.Parse(time.RFC3339, payMsg.CreatedAt)
			if err == nil && time.Since(createdAt) > 10*time.Minute {
				log.Printf("[Expired] Payment %s is older than 10 minutes. Expiring.", payMsg.PaymentID)

				db.Save(&model.Payment{
					PaymentID: payMsg.PaymentID,
					UserID:    payMsg.UserID,
					MerchantID: payMsg.MerchantID,
					Amount:    payMsg.Amount,
					Status:    "expired",
					CreatedAt: createdAt,
				})

				d.Ack(false)
				continue
			}

			log.Printf("[Processing] Handling payment ID: %s for User: %s, Amount: %.2f", payMsg.PaymentID, payMsg.UserID, payMsg.Amount)

			db.Create(&model.Payment{
				PaymentID: payMsg.PaymentID,
				UserID:    payMsg.UserID,
				MerchantID: payMsg.MerchantID,
				Amount:    payMsg.Amount,
				Status:    "processing",
				CreatedAt: createdAt,
			})

			time.Sleep(3 * time.Second)

			db.Model(&model.Payment{}).Where("payment_id = ?", payMsg.PaymentID).Update("status", "completed")

			log.Printf("[Completed] Payment ID: %s successfully processed.", payMsg.PaymentID)

			d.Ack(false)
		}
	}()

	<-forever
}