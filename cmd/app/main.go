package main

import (
	"feature-flags/pkg/config"
	httpapi "feature-flags/pkg/http"
	"feature-flags/pkg/service"
	"feature-flags/pkg/storage"
)

func main() {
	// Загружаем конфигурацию из переменных окружения
	cfg := config.MustLoad()

	// Инициализируем подключение к Postgres
	db := storage.MustInitPostgres(cfg.PostgresDSN)
	defer db.Close()

	// Создаём сервис работы с фичами (256 элементов кэша, TTL 15 минут)
	varsService := service.NewFeatureService(db, 256, 15)

	// Поднимаем HTTP API
	srv := httpapi.NewServer(cfg, varsService)
	srv.Run()
}
