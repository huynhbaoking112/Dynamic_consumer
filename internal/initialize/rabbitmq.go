package initialize

import (
	"fmt"

	"event_service/global"
	"event_service/internal/common"

	"github.com/rabbitmq/amqp091-go"
)

func InitRabbitMQ() {
	fmt.Println("Initializing RabbitMQ connection...")

	cfg := global.Config.RabbitMQ

	connStr := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		cfg.User, cfg.Password, cfg.Host, cfg.Port)

	conn, err := amqp091.Dial(connStr)
	if err != nil {
		panic(fmt.Errorf("%w: %v", common.ErrRabbitConnection, err))
	}

	// Create channel activity log queue
	ch, err := conn.Channel()
	if err != nil {
		panic(fmt.Errorf("%w: %v", common.ErrRabbitChannel, err))
	}

	// Declare queue activity log queue
	q, err := ch.QueueDeclare(
		cfg.ActivityLogQueue, // name
		true,                 // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		panic(fmt.Errorf("failed to declare queue: %v", err))
	}

	err = ch.QueueBind(
		q.Name,                    // queue name
		cfg.ActivityLogBindingKey, // routing key
		cfg.IAMExchange,           // exchange
		false,
		nil,
	)
	if err != nil {
		panic(fmt.Errorf("failed to bind queue: %v", err))
	}

	global.RabbitMQ = conn
	global.ActivityLogRabbitCh = ch

	fmt.Printf("RabbitMQ connected successfully")
}
