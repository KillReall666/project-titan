package config

import (
	"time"
)

type Config struct {
	RequestCount   int
	GorutinesCount int
	URL            string
	AuthURL        string
	Username       string
	Password       string
	Timeout        time.Duration
	verbose        bool
}

const (
	reqCount = 10000                          // общее количество запросов
	gCount   = 200                            // количество параллельных горутин (concurrency)
	url      = "http://localhost:8080/packet" // URL эндпоинта Gateway
	authURL  = "http://localhost:8083/login"  // "URL для получения токена (/login)
	username = "admin"                        // имя пользователя для логина
	password = "supersecret123"               // пароль для логина
	timeout  = 5 * time.Second                //таймаут на каждый запрос
	verbose  = false                          // "выводить ошибки каждого запроса")
)

func New() *Config {
	cfg := &Config{
		RequestCount:   reqCount,
		GorutinesCount: gCount,
		URL:            url,
		AuthURL:        authURL,
		Username:       username,
		Password:       password,
		Timeout:        timeout,
		verbose:        verbose,
	}

	return cfg

}
