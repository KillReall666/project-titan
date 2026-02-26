package config

import (
	"time"
)

type Config struct {
	RequestCount    int
	GoroutinesCount int
	URL             string
	AuthURL         string
	Username        string
	Password        string
	Timeout         time.Duration
	Verbose         bool
}

const (
	reqCount = 10000                          // Общее количество запросов
	gCount   = 200                            // Количество параллельных горутин (concurrency)
	url      = "http://localhost:8080/packet" // URL эндпоинта Gateway
	authURL  = "http://localhost:8080/login"  // "URL для получения токена (/login)
	username = "kirill@mail.ru"               // Имя пользователя для логина
	password = "test-pass"                    // Пароль для логина
	timeout  = 5 * time.Second                //Таймаут на каждый запрос
	verbose  = false                          // Выводить ошибки каждого запроса
)

func New() *Config {
	cfg := &Config{
		RequestCount:    reqCount,
		GoroutinesCount: gCount,
		URL:             url,
		AuthURL:         authURL,
		Username:        username,
		Password:        password,
		Timeout:         timeout,
		Verbose:         verbose,
	}

	return cfg

}
