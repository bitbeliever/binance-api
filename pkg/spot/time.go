package spot

import (
	"context"
	"log"
	"time"
)

// ServerTime 获取服务器时间 用于同步
func ServerTime() {
	t, err := NewClient().NewServerTimeService().Do(context.Background())
	if err != nil {
		log.Println(err)
	}

	log.Println(t, time.UnixMilli(t), time.Unix(t, 0))
}
