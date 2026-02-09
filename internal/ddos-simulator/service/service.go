package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"titan/common/logger"
	"titan/internal/ddos-simulator/config"
)

var (
	n        = flag.Int("n", 10000, "общее количество запросов")
	c        = flag.Int("c", 200, "количество параллельных горутин (concurrency)")
	url      = flag.String("url", "http://localhost:8080/packet", "URL эндпоинта Gateway")
	authURL  = flag.String("auth", "http://localhost:8083/login", "URL для получения токена (/login)")
	username = flag.String("u", "admin", "имя пользователя для логина")
	password = flag.String("p", "supersecret123", "пароль для логина")
	timeout  = flag.Duration("timeout", 5*time.Second, "таймаут на каждый запрос")
	verbose  = flag.Bool("v", false, "выводить ошибки каждого запроса")
)

type DdosSimulatorService struct {
	cfg *config.Config
}

func New(cfg *config.Config) *DdosSimulatorService {
	return &DdosSimulatorService{
		cfg: cfg,
	}
}

func (s *DdosSimulatorService) Setup() error {
	token, err := s.Login()
	if err != nil {
		return err //todo: подумать что тут
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	sem := make(chan struct{}, *c)
	success := 0
	failed := 0
	mu := sync.Mutex{}
	startTime := time.Now()
loop:
	for i := 0; i < *n; i++ {
		select {
		case <-stop:
			logger.Logger.Warn("Receive a stop signal, finishing...")
			break loop
		default:
		}
		wg.Add(1)
		sem <- struct{}{}

		go func(reqID int) {
			defer wg.Done()
			defer func() { <-sem }()

			err = s.Attack(token)
			if err != nil {
				mu.Lock()
				failed++
				mu.Unlock()
				if *verbose {
					logger.Logger.Info("[req %d] error; %v", reqID, err)
				}
			} else {
				mu.Lock()
				success++
				mu.Unlock()
			}
		}(i)

		// Небольшая задержка, чтобы не перегружать систему сразу (опционально)
		// time.Sleep(time.Millisecond * 1)
	}
	wg.Wait()

	duration := time.Since(startTime)
	rps := float64(success) / duration.Seconds()

	fmt.Println("\nРезультаты:")
	fmt.Printf("Успешно:     %d\n", success)
	fmt.Printf("С ошибками:  %d\n", failed)
	fmt.Printf("Всего:       %d\n", success+failed)
	fmt.Printf("Время:       %v\n", duration.Round(time.Millisecond))
	fmt.Printf("RPS:         %.2f\n", rps)
	if success > 0 {
		fmt.Printf("Средний RPS на горутину: %.2f\n", rps/float64(*c))
	}

}

func (s *DdosSimulatorService) Login() (string, error) {
	payload := map[string]string{
		"username": *username,
		"password": *password,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", *authURL, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: *timeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("запрос на /login: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("статус %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("не удалось распарсить ответ /login: %w", err)
	}

	if result.Token == "" {
		return "", errors.New("токен не вернулся в ответе")
	}

	return result.Token, nil
}

func (s *DdosSimulatorService) Attack() error {
	packet := map[string]interface{}{
		"packet_id": uuid.New().String(),
		"data":      fmt.Sprintf("test-packet-%d", rand.Intn(999999)),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"size":      rand.Intn(1024) + 128, // имитация разного размера пакета
	}

	body, _ := json.Marshal(packet)

	req, _ := http.NewRequest("POST", *url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: *timeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Читаем тело ответа, чтобы соединение могло быть переиспользовано
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("статус %d", resp.StatusCode)
	}

	return nil
}
