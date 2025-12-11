# Именование пакетов: Функциональность vs Типы

> "Хорошее имя — это документация"
> — **Rob Pike**

---

## Оглавление

- [Главное правило](#главное-правило)
- [Почему типы — плохо](#почему-типы--плохо)
- [Функциональный подход](#функциональный-подход)
- [Примеры из stdlib](#примеры-из-stdlib)
- [Антипаттерны](#антипаттерны)
- [Практические сценарии](#практические-сценарии)

---

## Главное правило

### Называй пакет по тому **ЧТО ОН ДЕЛАЕТ**, а не **ЧЕМ ОН ЯВЛЯЕТСЯ**

```go
// Правильно: описывает функциональность
package user       // управляет пользователями
package payment    // обрабатывает платежи
package email      // отправляет email'ы

// Неправильно: описывает тип
package models     // модели чего?
package handlers   // обработчики чего?
package services   // сервисы для чего?
```

### Тест "одного предложения"

Если не можешь описать пакет **одним предложением** без слов "содержит", "включает", "различные" — имя выбрано неправильно.

```
Правильно:
- user: "Управляет пользователями и аутентификацией"
- payment: "Обрабатывает платежные транзакции"
- storage: "Сохраняет и загружает данные из БД"

Неправильно:
- models: "Содержит различные модели" (какие? зачем?)
- utils: "Включает утилиты" (для чего?)
- common: "Общий код" (что общего?)
```

---

## Почему типы — плохо

### Проблема 1: Нарушение Single Responsibility Principle

```go
// Пакет по типу = grab-bag (мусорка)
package models

type User struct { ... }
type Order struct { ... }
type Payment struct { ... }
type Invoice struct { ... }
type Product struct { ... }
// ...ещё 20 моделей
```

**Проблемы:**
- Пакет не имеет чёткой ответственности
- Растёт до сотен файлов
- Невозможно понять границы
- Высокая связанность

---

### Проблема 2: Циклические зависимости

```go
// Группировка по типам → циклические зависимости
package models
type User struct { Orders []Order }

package services
type UserService struct {
    repo repositories.UserRepository  // models → repositories
}

package repositories
type UserRepository struct {
    // нужен models.User
}
// services → repositories → models → services (цикл!)
```

```go
// Функциональная группировка → нет циклов
package user
type User struct { ... }
type Repository interface { ... }
type Service struct { repo Repository }

package order
type Order struct {
    UserID string  // просто ID, без импорта user
}
```

---

### Проблема 3: Бессмысленные импорты

```go
// С группировкой по типам
import (
    "project/models"
    "project/handlers"
    "project/services"
)

// Что делает этот код? Непонятно из импортов...
func main() {
    user := models.User{...}
    service := services.UserService{...}
    handler := handlers.UserHandler{...}
}
```

```go
// С функциональной группировкой
import (
    "project/user"
    "project/payment"
    "project/email"
)

// Понятно: работаем с пользователями, платежами, email(сразу легче читается, да?)
func main() {
    u := user.User{...}
    svc := user.NewService(...)
    paymentSvc := payment.NewService(...)
}
```

---

## Функциональный подход

### Принцип: Один пакет = один Bounded Context

```go
// Функциональная группировка
internal/
  user/              # Bounded Context: управление пользователями
    user.go          # Entity
    repository.go    # Интерфейс репозитория
    service.go       # Бизнес-логика
    handler.go       # HTTP handlers
    postgres.go      # Реализация для Postgres

  payment/           # Bounded Context: платежи
    payment.go
    processor.go
    sbp.go        	 # Интеграция с СБП

  notification/      # Bounded Context: уведомления
    notification.go
    email.go
    sms.go
```

### Преимущества:

1. **Чёткие границы:** Каждый пакет — законченный модуль
2. **Низкая связанность:** Пакеты общаются через интерфейсы
3. **Легко масштабировать:** Новый функционал = новый пакет
4. **Понятно из импортов:** `import "project/payment"` → работа с платежами

---

### Как организовывать код внутри пакета

```go
// Правильно: всё относящееся к user в одном пакете
package user

// Доменная модель
type User struct {
    ID       string
    Email    string
    Password string
}

// Интерфейс на стороне потребителя
type Repository interface {
    GetByID(id string) (*User, error)
    Save(user *User) error
}

// Бизнес-логика
type Service struct {
    repo Repository
}

func NewService(repo Repository) *Service {
    return &Service{repo: repo}
}

func (s *Service) Register(email, password string) (*User, error) {
    // валидация
    // хеширование пароля
    // сохранение
}

// HTTP handler
type Handler struct {
    service *Service
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
    // парсинг request
    // вызов service
    // формирование response
}
```

**Важно:** Всё связанное с пользователями — в одном месте.

---

## Примеры из stdlib

Go stdlib — **лучший учебник** по именованию пакетов.

### Пример 1: net/*

```go
// Функциональная группировка
net/              # сетевое взаимодействие
  http/           # HTTP протокол
    client.go
    server.go
  smtp/           # SMTP протокол
  url/            # работа с URL
  rpc/            # RPC коммуникация
```

**Обрати внимание:**
- Нет `net/models/` или `net/handlers/`
- Каждый пакет = конкретный протокол/функциональность
- Импорт говорит сам за себя: `import "net/http"`

---

### Пример 2: encoding/*

```go
encoding/         # кодирование данных
  json/           # JSON формат
  xml/            # XML формат
  base64/         # Base64 кодирование
  csv/            # CSV формат
  gob/            # Go binary format
```

**Почему не так:**
```go
encoding/
    encoders/
      json_encoder.go
      xml_encoder.go
    decoders/
      json_decoder.go
```

**Потому что:**
- Группировка по типам (`encoders/decoders`) бессмысленна
- JSON-кодирование и декодирование — одна функциональность
- Пакет `encoding/json` — законченный модуль для работы с JSON

---

### Пример 3: database/sql

```go
database/
  sql/            # SQL взаимодействие
    driver/       # интерфейсы драйверов
```

**Обрати внимание:**
- Нет `database/repositories/` или `database/queries/`
- Пакет предоставляет функциональность работы с SQL
- Драйверы — отдельный sub-пакет (они расширяют функциональность)

---

## Антипаттерны

### Антипаттерн 1: Grab-Bag пакеты

```go
// Плохо: пакет-мусорка
package utils

func ValidateEmail(email string) bool { ... }
func FormatDate(date time.Time) string { ... }
func CalculateTax(amount float64) float64 { ... }
func SendEmail(to, subject, body string) error { ... }
```

**Проблемы:**
- Нет чёткой ответственности
- Растёт до бесконечности
- Высокая связанность (все импортируют `utils`)

**Решение:**
```go
// Хорошо: распредели по функциональным пакетам
package email
func Validate(email string) bool { ... }
func Send(to, subject, body string) error { ... }

package format
func Date(date time.Time) string { ... }

package tax
func Calculate(amount float64) float64 { ... }
```

---

### Антипаттерн 2: Технические слои

```go
// Плохо: горизонтальная архитектура (по типам)
project/
  models/
    user.go
    order.go
  repositories/
    user_repository.go
    order_repository.go
  services/
    user_service.go
    order_service.go
  handlers/
    user_handler.go
    order_handler.go
```

**Проблемы:**
- Чтобы понять `user`, нужно открыть 4 файла в 4 папках
- Изменение в `User` → правки в 4 местах
- Нет изоляции (все слои видят друг друга)

```go
// Хорошо: вертикальная архитектура (по функциональности)
project/
  internal/
    user/
      user.go          # модель + репозиторий + сервис + handler
    order/
      order.go
```

---

### Антипаттерн 3: Суффиксы/префиксы в названиях

```go
// Плохо
package usermodels
package userhandlers
package userservices

// Использование:
import "project/usermodels"
u := usermodels.User{...}  // избыточно
```

```go
// Хорошо
package user

// Использование:
import "project/user"
u := user.User{...}        // лаконично
svc := user.NewService(...)
```

**Правило:** Имя пакета уже даёт контекст, не дублируй его в типах.

---

### Антипаттерн 4: Слишком общие названия

```go
// Избегайте:
package data       // какие данные?
package core       // что является core?
package base       // базовое для чего?
package common     // общее для чего?
package helpers    // помощники в чём?
package utilities  // утилиты для чего?
package lib        // библиотека чего?
package app        // всё приложение?
package ...
```

**Задавай себе вопрос:** Если название пакета применимо к любому проекту — оно слишком общее.

---

## Практические сценарии

### Сценарий 1: REST API для интернет-магазина

```go
// Группировка по типам ("традиционный подход")
ecommerce/
  models/
    user.go
    product.go
    order.go
    payment.go
  controllers/
    user_controller.go
    product_controller.go
    order_controller.go
  services/
    user_service.go
    order_service.go
  repositories/
    user_repository.go
    order_repository.go
```

**Проблемы:**
- 4 файла для одной фичи (user)
- Сложно найти связанный код

```go
// Функциональная группировка (Go-way)
ecommerce/
  internal/
    user/           # Всё о пользователях
      user.go       # type User + валидация
      auth.go       # аутентификация
      repository.go # интерфейс
      postgres.go   # реализация
      handler.go    # HTTP API

    product/        # Всё о товарах
      product.go
      catalog.go    # управление каталогом
      search.go     # поиск товаров
      handler.go

    order/          # Всё о заказах
      order.go
      cart.go       # корзина
      checkout.go   # оформление
      handler.go

    payment/        # Всё о платежах
      payment.go
      webhook.go    # webhooks от платёжной системы
```

**Преимущества:**
- Всё связанное с заказами — в одной папке
- Легко найти код
- Модульность: можно извлечь `payment` в отдельный микросервис

---

### Сценарий 2: CLI утилита для девопса

```go
// Избыточная структура
devops-cli/
  cmd/
    commands/
      deploy_command.go
      build_command.go
  internal/
    models/
      config.go
    services/
      deploy_service.go
```

```go
// Простая структура
devops-cli/
  main.go
  deploy.go      # команда + логика деплоя
  build.go       # команда + логика сборки
  config.go      # конфигурация
```

**Почему проще:**
- Проект небольшой (< 1000 строк)
- Нет необходимости в модулях
- Легко понять структуру

---

### Сценарий 3: Микросервис с очередями

```go
// Функциональная группировка
notification-service/
  internal/
    notification/    # Ядро: создание уведомлений
      notification.go
      service.go

    email/           # Канал: email
      sender.go      # отправка
      template.go    # шаблоны
      smtp.go        # SMTP клиент

    sms/             # Канал: SMS
      sender.go
      exolve.go      # интеграция с Exolve

    push/            # Канал: Push notifications
      sender.go
      fcm.go         # Firebase Cloud Messaging

    queue/           # Очереди
      consumer.go    # потребитель из RabbitMQ
      publisher.go   # публикация в RabbitMQ
```

**Обрати внимание:**
- `email`, `sms`, `push` — функциональности (каналы доставки)
- НЕ `senders/email_sender.go` (группировка по типу)
- Каждый пакет = законченный модуль

---

## Именование: Запомни

### 1. Короткие имена

```go
// Правильно
package user
package http
package sql

// Слишком длинно
package usermanagement
package hypertext
package structuredquerylanguage
```

**Почему:** Имя пакета повторяется при каждом использовании.

```go
// С коротким именем
user.User{...}
user.NewService(...)

// С длинным именем
usermanagement.User{...}            // избыточно
usermanagement.NewService(...)
```

---

### 2. Существительные, а не глаголы

```go
// Правильно
package user       // сущность
package payment    // сущность
package storage    // место

// Неправильно
package managing   // чем управляет?
package handling   // что обрабатывает?
package processing // что обрабатывает?
```

---

### 3. Единственное число

```go
// Правильно
package user       // не users
package order      // не orders
package product    // не products

// Неправильно
package users
```

**Почему:** `user.User` читается лучше, чем `users.User`

**Исключение:** Если пакет содержит утилиты (как `strings`, `bytes` в stdlib)

---

### 4. Избегай повторения контекста

```go
// Плохо: избыточность
package user

type UserService struct { ... }   // "user" дважды
type UserRepository struct { ... }

// Хорошо
package user

type Service struct { ... }       // user.Service
type Repository struct { ... }    // user.Repository
```

---

### 5. Используйте domain language (термины из предметной области)

```go
// Правильно: термины из предметной области
package order      // заказ
package cart       // корзина
package checkout   // оформление
package shipment   // доставка

// Неправильно: технические термины
package entities
package aggregates
package valueobjects
```

---

## Запомни:

1. **Функциональность > Типы**
   - Группируй по тому, что делает код
   - Избегай технических слоёв

2. **Один пакет = один контекст**
   - Всё связанное — в одном месте
   - Чёткие границы

3. **Простые имена**
   - Короткие, существительные
   - Без избыточности

### Главное правило:

> "Если ты не можешь описать пакет одним предложением без слов 'содержит' или 'различные' — пересмотри имя."

---

## Будет полезно:

- [Package names](https://go.dev/blog/package-names) — официальный блог Go
- [Organizing Go code](https://go.dev/talks/2014/organizeio.slide) — Andrew Gerrand
- [Standard Package Layout](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1) — Ben Johnson

---

**Помни:** Хорошее именование — половина успеха в поддержке кода.
