package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const telegramEndpoint = "https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s"

type telegramClient struct {
	botToken string
	chatID   string
	endpoint string
	queue    chan string
	cancel   chan struct{}
}

func (c *telegramClient) SendMsg(text string) {
	c.queue <- text
}

func (c *telegramClient) Cancel() {
	c.cancel <- struct{}{}
}

func (c *telegramClient) background() {
	for {
		select {
		case txt := <-c.queue:
			query := fmt.Sprintf(c.endpoint, url.QueryEscape(c.botToken), c.chatID, url.QueryEscape(txt))
			resp, err := http.Get(query)
			if err != nil {
				log.Println("Error sending message:", err)
				continue
			}
			if resp.StatusCode != http.StatusOK {
				log.Println("Response status code:", resp.Status)
				continue
			}
		case <-c.cancel:
			return
		}
	}
}

// Hook
type telegramHook struct {
	client    *telegramClient
	level     logrus.Level
	mention   []string
	formatter logrus.Formatter
	mutex     sync.Mutex
}

func (h *telegramHook) Fire(entry *logrus.Entry) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if len(h.mention) > 0 {
		time, level, msg := entry.Time, entry.Level, entry.Message
		entry = entry.WithField("mention", strings.Join(h.mention, ", "))
		entry.Time, entry.Level, entry.Message = time, level, msg
	}

	buf, err := h.formatter.Format(entry)
	if err != nil {
		return err
	}

	h.client.SendMsg(string(buf))
	return nil
}

func (h *telegramHook) Levels() []logrus.Level {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	var out []logrus.Level
	for _, level := range logrus.AllLevels {
		if level <= h.level {
			out = append(out, level)
		}
	}

	return out
}

func addTelegramHook(l *logrus.Logger) {
	c := &telegramClient{
		botToken: viper.GetString("LOG_HOOK_TELE_BOT_TOKEN"),
		chatID:   viper.GetString("LOG_HOOK_TELE_CHAT_ID"),
		endpoint: telegramEndpoint,
		queue:    make(chan string, 128),
		cancel:   make(chan struct{}),
	}

	go c.background()

	h := &telegramHook{
		client:    c,
		level:     getLogLevel(viper.GetString("LOG_HOOK_TELE_LEVEL")),
		mention:   strings.Split(viper.GetString("LOG_HOOK_TELE_MENTIONS"), ","),
		formatter: formatter,
	}

	l.AddHook(h)
}
