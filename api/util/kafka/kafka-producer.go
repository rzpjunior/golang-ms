package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"git.edenfarm.id/cuxs/env"
	"github.com/segmentio/kafka-go"
)

// Produce
func Produce(ctx context.Context, r interface{}, Topic string) (err error) {
	// initialize a counter
	i := 0

	// intialize the writer with the broker addresses, and the topic
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{env.GetString("BROKER_ADDRESS", "0.0.0.0:9092")},
		Topic:        Topic,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		BatchSize:    1,
		BatchTimeout: 10 * time.Millisecond,
	})
	defer w.Close()
	// each kafka message has a key and value. The key is used
	// to decide which partition (and consequently, which broker)
	// the message gets published on
	jsonReq, _ := json.Marshal(r)
	err = w.WriteMessages(ctx, kafka.Message{
		Key: []byte(strconv.Itoa(i)),
		// create an arbitrary message payload for the value

		Value: []byte(jsonReq),
	})
	if err != nil {
		return err
	}

	// log a confirmation once the message is written
	fmt.Println("writes:", i)
	i++
	// sleep for a second
	time.Sleep(time.Second)

	return err
}
