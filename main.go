package main

import (
	"context"
	"fmt"
	"github.com/vwency/intern-task/unused_sub_pub"
	"time"
)

func main() {
	sp := subpub.NewSubPub()

	// Подписчик 1
	sub1 := sp.Subscribe("topic1", func(msg interface{}) {
		fmt.Printf("Subscriber 1 received: %v\n", msg)
	})

	// Подписчик 2 (медленный)
	sub2 := sp.Subscribe("topic1", func(msg interface{}) {
		time.Sleep(2 * time.Second)
		fmt.Printf("Subscriber 2 received: %v\n", msg)
	})

	// Публикация сообщений
	for i := 0; i < 5; i++ {
		err := sp.Publish("topic1", fmt.Sprintf("message %d", i))
		if err != nil {
			fmt.Println("Publish error:", err)
		}
	}

	// Увеличиваем время ожидания для медленного подписчика
	time.Sleep(3 * time.Second) // было 1 секунду, увеличиваем до 3

	// Сначала закрываем publisher
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sp.Close(ctx); err != nil {
		fmt.Println("Close error:", err)
	}

	// Затем отписываемся (Unsubscribe после Close безопасен)
	sub1.Unsubscribe()
	sub2.Unsubscribe()
}
