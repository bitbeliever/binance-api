package spot

import (
	"context"
	"log"
	"time"
)

func ServerTime() {
	t, err := NewClient().NewServerTimeService().Do(context.Background())
	if err != nil {
		log.Println(err)
	}

	log.Println(t, time.UnixMilli(t), time.Unix(t, 0))
}
