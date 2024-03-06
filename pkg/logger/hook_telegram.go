package logger

import (
	"fmt"
	"github.com/kurneo/go-template/config"
	"github.com/kurneo/go-template/pkg/support/slices"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type telegramClient struct {
	botToken string
	chatID   int64
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

func (c *telegramClient) flush() {
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
	mention   map[logrus.Level][]string
	formatter logrus.Formatter
	mutex     sync.Mutex
}

func (h *telegramHook) SetLevel(level logrus.Level) *telegramHook {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.level = level
	return h
}

func (h *telegramHook) SetMention(m map[logrus.Level][]string) *telegramHook {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.mention = m
	return h
}

func (h *telegramHook) MentionOn(level logrus.Level, users ...string) *telegramHook {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for _, lvl := range logrus.AllLevels {
		if lvl > level {
			continue
		}

		if _, ok := h.mention[lvl]; !ok {
			h.mention[lvl] = make([]string, 0)
		}
		h.mention[lvl] = append(h.mention[lvl], users...)
	}
	return h
}

func (h *telegramHook) SetFormatter(formatter logrus.Formatter) *telegramHook {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.formatter = formatter
	return h
}

func (h *telegramHook) Fire(entry *logrus.Entry) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	var users []string
	for _, user := range h.mention[entry.Level] {
		users = append(users, "@"+user)
	}

	if len(users) > 0 {
		time, level, msg := entry.Time, entry.Level, entry.Message
		entry = entry.WithField("mention", strings.Join(users, ", "))
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

func addTelegramHook(cfg config.LogHook, l *logrus.Logger) {
	c := &telegramClient{
		botToken: cfg.Get("bot_token").(string),
		chatID:   int64(cfg.Get("chat_id").(int)),
		endpoint: cfg.Get("endpoint").(string),
		queue:    make(chan string, 128),
		cancel:   make(chan struct{}),
	}
	go c.flush()

	h := &telegramHook{
		client:    c,
		level:     getLogLevel(cfg.Get("level").(string)),
		mention:   make(map[logrus.Level][]string),
		formatter: formatter,
	}
	h.SetLevel(getLogLevel(cfg.Get("level").(string)))
	mentions := make(map[logrus.Level][]string, 0)
	for k, v := range cfg.Get("mentions").(config.LogHook) {
		mentions[getLogLevel(k)] = slices.Map(v.([]interface{}), func(v interface{}) string {
			return v.(string)
		})
	}
	h.SetMention(mentions)
	l.AddHook(h)
}
