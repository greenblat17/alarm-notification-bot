package eventconsumer

import (
	"log"
	"time"

	"github.com/greenblat17/alarm-notification-bot/internal/events"
)

type EventConsumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) *EventConsumer {
	return &EventConsumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *EventConsumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERROR]: consumer: %w", err)

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err)

			continue
		}
	}
}

func (c *EventConsumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("[ERROR] cannot handle event: %w", err)

			continue
		}
	}

	return nil
}
