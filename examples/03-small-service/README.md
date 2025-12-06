# Пример 3: Small Service with Internal Packages

## Что делает этот пример?

HTTP API для работы с пользователями — создание, получение списка. Хранение в памяти (in-memory).

## Когда использовать?

- Нужна модульность и переиспользование кода
- Требуется инкапсуляция внутренней логики
- Работает команда разработчиков
- API/веб-сервис с несколькими доменными областями

## Структура

```
03-small-service/
├── main.go                    # Точка входа, wiring
├── internal/                  # Приватный код
│   ├── config/
│   │   └── config.go          # Конфигурация
│   └── user/                  # Доменная область: пользователи
│       ├── user.go            # Модель и бизнес-логика
│       ├── handler.go         # HTTP-хендлеры
│       └── storage.go         # Работа с данными
├── go.mod
└── README.md
```

## Почему именно так?

- **`/internal`** — Go-компилятор гарантирует, что код из `internal/` не может быть импортирован другими проектами.
- **Группировка по функциональности (не по типам)** — пакет `user/` содержит ВСЁ, что связано с пользователями: модель, хендлеры, хранилище.

```
internal/user/          # Всё о пользователях (ok)
internal/models/        # Модели чего? (not ok)
internal/handlers/      # Обработчики чего?(not ok)
```

- **`main.go` в корне** — точка входа остаётся простой, только инициализация зависимостей. Вся логика в `internal/`.
- **Плоская иерархия** — только 2 уровня вложенности (`internal/user/`). Не создаём глубокие иерархии типа `internal/domain/entities/user/models/`.

## Как запустить?

```bash
# Запустить сервер
go run main.go

# В другом терминале:

# Получить всех пользователей
curl http://localhost:8080/users

# Создать пользователя
curl -X POST http://localhost:8080/users/create \
  -H "Content-Type: application/json" \
  -d '{"name":"Charlie","email":"charlie@example.com"}'
```

## Ключевые концепции

### Package-Oriented Design
Каждый пакет = законченная функциональность с чётким API.

```go
package user

type User struct { ... }
type Handler struct { ... }
type Storage struct { ... }
```

### Dependency Injection через конструкторы
```go
// main.go
storage := user.NewStorage()
handler := user.NewHandler(storage)
```

### Нет преждевременных интерфейсов
`Storage` — конкретная реализация, не интерфейс. Интерфейсы создаются на стороне потребителя (consumer-side interfaces).

### Именование без stuttering
```go
user.Handler      // читаемо.
user.UserHandler  // визуальный мусор.
```
