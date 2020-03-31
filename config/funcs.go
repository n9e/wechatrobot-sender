package config

import (
	"fmt"
	"os"
	"time"

	"github.com/n9e/wechatrobot-sender/corp"
	"github.com/toolkits/pkg/logger"
)

// InitLogger init logger toolkits
func InitLogger() {
	c := Get().Logger

	lb, err := logger.NewFileBackend(c.Dir)
	if err != nil {
		fmt.Println("cannot init logger:", err)
		os.Exit(1)
	}

	lb.SetRotateByHour(true)
	lb.SetKeepHours(c.KeepHours)

	logger.SetLogging(c.Level, lb)
}

func Test(args []string) {

	if len(args) == 0 {
		fmt.Println("im not given")
		os.Exit(1)
	}

	for i := 0; i < len(args); i++ {
		err := corp.Send(corp.Message{
			ToUser:  args[i],
			MsgType: "text",
			Text:    corp.Content{Content: fmt.Sprintf("test message from n9e at %v", time.Now())},
		})

		if err != nil {
			fmt.Printf("send to %s fail: %v\n", args[i], err)
		} else {
			fmt.Printf("send to %s succ\n", args[i])
		}
	}
}
