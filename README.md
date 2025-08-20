# Feature Flags Service

Сервис для хранения и управления **динамическими конфигурациями (feature flags)**.  
Реализован на **Go**, использует:
- [Fiber](https://github.com/gofiber/fiber) + [huma.rocks](https://huma.rocks) — HTTP API
- [PostgreSQL](https://www.postgresql.org/) — хранение данных
- [golang-lru](https://github.com/hashicorp/golang-lru) — кэш в памяти
- [Goose](https://github.com/pressly/goose) — миграции базы данных
- [Makefile](Makefile) — удобные команды для миграций и запуска

---

## 🚀 Возможности

- Хранение переменных (feature flags) в Postgres
- In-memory кэш (LRU + TTL 15 минут) для ускорения доступа
- REST API с двумя ручками:
  1. **GET /var/{var_name}** — получить значение переменной (с кэшем)  
     ```json
     {
       "key": "first-var",
       "value": 0.4
     }
     ```
  2. **POST /var/set** — установить/обновить значение переменной (инвалидация кэша)  
     **Request:**
     ```json
     {
       "key": "asdqwe",
       "value": "12345"
     }
     ```
     **Response:**
     ```json
     {
       "message": "var successfully updated"
     }
     ```

---

## 🛠️ Установка и запуск

### 1. Клонировать репозиторий
```bash
git clone https://github.com/d1mk9/feature-flags.git
cd feature-flags
```

### 2. Настроить Postgres и окружение
Создай файл `.env` в корне:
```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
```

А параметры подключения к БД указываются в `conf/config.yaml`:
```yaml
postgres:
  host: localhost
  port: 5432
  db: featuredb
```

### 3. Применить миграции
Используется `goose`, а также `make` для удобства:
```bash
make migrate-up
```

### 4. Запустить сервис
```bash
make run
```

Сервис будет доступен по адресу:  
👉 http://localhost:8080

---

## 📂 Структура проекта

```
.
├── cmd/app/main.go        # Точка входа
├── pkg/
│   ├── config/            # Загрузка конфигурации (.env + config.yaml)
│   ├── handlers/          # HTTP-хендлеры
│   ├── http/              # Сервер + маршрутизация
│   ├── models/            # Reform-модели
│   ├── repository/        # Работа с Postgres через Reform/Bob
│   ├── service/           # Бизнес-логика + кэш
│   └── storage/           # Инициализация Postgres
├── migrations/            # Goose-миграции
├── conf/config.yaml       # Конфигурация БД
├── Makefile               # Утилитные команды (migrate/run)
└── README.md
```

---

## ✅ Проверка работы

### Установить переменную
```bash
curl -X POST http://localhost:8080/var/set   -H "Content-Type: application/json"   -d '{"key":"first-var","value":0.4}'
```

### Получить переменную
```bash
curl http://localhost:8080/var/first-var
```

Ответ:
```json
{
  "key": "first-var",
  "value": 0.4
}
```

---

## 📌 TODO
- [ ] Добавить тесты

---

## 📝 Лицензия
MIT
