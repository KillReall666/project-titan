package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"titan/internal/observability/logger"

	"github.com/google/uuid"

	"titan/internal/ddos-simulator/config"
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
	//ожидаем инициализации остальных сервисов
	flag.Parse()

	token, err := s.Login()
	if err != nil {
		return err //todo: подумать что тут
	}

	logger.Logger.Info("login success, ready to attack...")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	sem := make(chan struct{}, s.cfg.GoroutinesCount)
	success := 0
	failed := 0
	mu := sync.Mutex{}
	startTime := time.Now()
loop:
	for i := 0; i < s.cfg.RequestCount; i++ {
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
				if s.cfg.Verbose {
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
		fmt.Printf("Средний RPS на горутину: %.2f\n", rps/float64(s.cfg.GoroutinesCount))
	}

	return nil
}

func (s *DdosSimulatorService) Login() (string, error) {
	payload := map[string]string{
		"user_mail": s.cfg.Username,
		"password":  s.cfg.Password,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", s.cfg.AuthURL, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: s.cfg.Timeout}
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

func (s *DdosSimulatorService) Attack(token string) error {
	packet := map[string]interface{}{
		"packet_id": uuid.New().String(),
		"data":      fmt.Sprintf("test-packet-%d", rand.Intn(999999)),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"size":      rand.Intn(1024) + 128, // имитация разного размера пакета
	}

	body, _ := json.Marshal(packet)

	req, _ := http.NewRequest("POST", s.cfg.URL, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: s.cfg.Timeout}
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
