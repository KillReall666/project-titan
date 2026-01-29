DDoS-Test Microservices Project
Это микросервисный проект на Go, предназначенный для симуляции высокой нагрузки (например, "DDoS"-тестирования) на HTTP-сервер. 
Проект демонстрирует устойчивость системы под нагрузкой, асинхронную обработку событий через Kafka, сбор метрик через Prometheus и визуализацию в Grafana,
а также распределённую трассировку с OpenTelemetry + Zipkin.
Цель — создать простую, но масштабируемую систему, где:

Внешний трафик (симулируемый через инструменты вроде hey или vegeta) поступает на вход.
Пакеты обрабатываются асинхронно через Kafka.
Метрики (RPS, latency, errors, packets/sec, CPU/memory) собираются и отображаются в Grafana.
Tracing помогает отлаживать задержки и ошибки между сервисами.

Общая архитектура
[Внешний трафик (hey/vegeta)]
↓ (HTTP POST /packet)
[API Gateway] — rate limiting, auth, forwarding
↓ (HTTP POST /receive)
[Packet Receiver] — валидация, produce в Kafka (raw_packets)
↓ (Kafka topic: raw_packets)
[Traffic Analyzer] — consume, анализ (counts, anomalies), optional DB, produce alerts
↓ (Prometheus scrape /metrics)
[Metrics Service] — агрегация кастомных метрик, экспорт в Prometheus
↓
[Prometheus] ← scrape всех сервисов
↓
[Grafana] — дашборды (RPS, latency, CPU, packets/sec)
↓
[Zipkin] — UI для distributed traces (OpenTelemetry)

Синхронное общение: HTTP (Gin) между Gateway → Receiver.
Асинхронное общение: Kafka (топик raw_packets для сырых пакетов).
Observability:
Prometheus — сбор метрик (стандартные + кастомные).
Grafana — визуализация.
OpenTelemetry + Zipkin — tracing (spans в handlers и Kafka).
Logging — slog (structured JSON в stdout).


Проект рассчитан на высокую нагрузку: graceful shutdown, timeouts, rate limiting, backpressure через Kafka.

Сервисы:

1. API Gateway
   Порт: 8080
   Ответственность: Единая точка входа для внешнего трафика.
   Что делает:

Принимает POST /packet (JSON или raw body с "пакетами").
Проверяет rate limiting (golang.org/x/time/rate или Gin middleware).
Простая аутентификация (header Authorization).
Форвардит запрос в Packet Receiver по HTTP (POST /receive).
Экспортирует /metrics и /health.
Добавляет OTel spans для tracing.

Основные endpoints:

POST /packet — основной вход.
GET /health — liveness/readiness.
GET /metrics — Prometheus scrape.

2. Packet Receiver
Порт: 8081
Ответственность: Приём и первичная обработка пакетов.
Что делает:

Принимает HTTP-запросы от Gateway.
Валидирует данные (например, размер body > 0).
Производит сообщение в Kafka (топик raw_packets) — JSON {id, payload, timestamp}.
Добавляет OTel spans вокруг HTTP-handler и Kafka-produce.
Экспортирует /metrics и /health.

Основные endpoints:

POST /receive — internal (от Gateway).
GET /health, GET /metrics.

3. Traffic Analyzer
   Порт: 8082 (опционально, если нужен HTTP stats)
   Ответственность: Асинхронный анализ трафика.
   Что делает:

Consume из Kafka (топик raw_packets, consumer group analyzer-group).
Парсит сообщения, считает статистику (по IP, packets/sec, anomalies — если > threshold).
Хранит in-memory (sync.Map) или в PostgreSQL (pgx) — counts, alerts.
При аномалиях — produce в Kafka (опционально analyzed_events).
Добавляет OTel spans в consumer loop.
Экспортирует /metrics (custom: analysis_latency, anomalies_detected).

Общение: Только Kafka consumer (loop ReadMessage).
Опционально: GET /stats — текущие counts.

4. Metrics Service
   Порт: 9090
   Ответственность: Централизованный экспорт метрик.
   Что делает:

Собирает/агрегирует кастомные метрики (если нужно).
Экспортирует Prometheus endpoint /metrics.
Может consume из Kafka (если метрики асинхронные).
Добавляет OTel spans (если HTTP).

Основные endpoints:

GET /metrics — Prometheus scrape (основной).

Observability

Prometheus: scrape всех сервисов (config: targets localhost:8080, 8081, 8082, 9090).
Grafana: datasource Prometheus, дашборды для:
System: CPU, memory (runtime.MemStats).
App: requests_total, latency_histogram, errors_total, packets_per_second.

Zipkin: UI localhost:9411 — traces по service-name.
Logging: slog JSON в stdout (уровни info/error, trace_id).

Убить контейнер и volume БД: docker compose down -v
Собрать контейнер:  docker compose up --build
