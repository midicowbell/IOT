# Smart Warehouse System

Весовой мониторинг и автоматизация остатков.

Бэкенд-сервис на Go для автоматизации складского учета на основе данных от умных полок с датчиками веса. Система использует **PostgreSQL** для надежного хранения данных и предоставляет REST API для управления полками и товарами.

---

## Функционал

* **Мониторинг веса** — Прием данных от IoT-датчиков веса через REST API.
* **Расчет остатков** — Автоматическое вычисление количества товара на основе текущего веса нетто и веса единицы товара.
* **Оповещение о дефиците** — Отслеживание минимально допустимого порога веса (`MinWeight`) и выставление статуса `REFILL`.
* **Персистентность данных** — Все данные хранятся в PostgreSQL с логированием изменений веса.
* **Многослойная архитектура** — Handlers → Service → Storage с использованием интерфейсов для абстракции хранилища.

---

## Стек технологий

* **Язык:** Go 1.25
* **Роутер:** `github.com/gorilla/mux`
* **База данных:** PostgreSQL
* **Driver:** `github.com/jackc/pgx/v5`
* **Архитектура:** Многослойная (Handlers → Service → Storage)

---

## Требования

* Go 1.25+
* PostgreSQL 12+
* Docker (опционально)

---

## Установка и запуск

### 1. Клонирование репозитория
```bash
git clone https://github.com/midicowbell/IOT.git
cd IOT
```

### 2. Установка зависимостей
```bash
go mod download
```

### 3. Подготовка базы данных

Создайте базу данных и таблицы:

```sql
CREATE DATABASE warehouse_iot;

\c warehouse_iot;

CREATE TABLE shelves (
    shelf_id INTEGER PRIMARY KEY,
    product_id INTEGER,
    current_weight FLOAT NOT NULL,
    status VARCHAR(50) DEFAULT 'EMPTY',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
    product_id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    weight_grams FLOAT NOT NULL,
    min_weight_grams FLOAT NOT NULL
);

CREATE TABLE weight_logs (
    log_id SERIAL PRIMARY KEY,
    shelf_id INTEGER NOT NULL REFERENCES shelves(shelf_id),
    raw_weight FLOAT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE shelves ADD CONSTRAINT fk_product 
    FOREIGN KEY (product_id) REFERENCES products(product_id);
```

### 4. Запуск сервиса

```bash
go run main.go
```

Сервис запустится на `http://localhost:8080`

---

## Конфигурация

Строка подключения к БД устанавливается в `main.go`:

```go
url := "postgres://postgres:12345@localhost:5432/warehouse_iot"
```

Формат: `postgres://username:password@host:port/database`

---

## API Endpoints

### Управление полками

#### Получить статус всех полок
```
GET /api/shelves
```

**Ответ:**
```json
[
    {
        "id": 1,
        "product": {
            "id": 1,
            "name": "Товар A",
            "weight": 100.0,
            "min_weight": 500.0
        },
        "quantity": 5,
        "current_weight": 500.0,
        "status": "OK",
        "updated_at": "2026-06-01T21:09:10Z"
    }
]
```

#### Добавить новую полку
```
POST /api/shelves
Content-Type: application/json

{
    "id": 1,
    "current_weight": 0.0,
    "status": "EMPTY"
}
```

**Ответ (201 Created):**
```json
{
    "id": 1,
    "product": null,
    "quantity": 0,
    "current_weight": 0.0,
    "status": "EMPTY"
}
```

#### Удалить полку
```
DELETE /api/shelves/{id}
```

#### Обновить вес на полке (от датчика)
```
PATCH /api/shelves/weight
Content-Type: application/json

{
    "id": 1,
    "weight": 750.5
}
```

**Ответ:**
```json
{
    "answer": "STATUS: Вес в норме",
    "time": "2026-06-01T21:09:10Z"
}
```

Если вес ниже `MinWeight`, вернет:
```json
{
    "answer": "STATUS: Необходим рефил",
    "time": "2026-06-01T21:09:10Z"
}
```

---

### Управление товарами

#### Привязать товар к полке
```
POST /api/products
Content-Type: application/json

{
    "product_id": 1,
    "shelf_id": 1,
    "name": "Товар A",
    "weight": 100.0,
    "quantity": 5,
    "min_weight": 500.0
}
```

#### Получить информацию о товаре на полке
```
GET /api/products/{shelf_id}
```

**Ответ:**
```json
{
    "id": 1,
    "name": "Товар A",
    "weight": 100.0,
    "min_weight": 500.0
}
```

#### Удалить товар с полки
```
DELETE /api/products/{shelf_id}
```

---

## Коды ошибок

| Код | Описание |
|-----|----------|
| 200 | OK - Успешно |
| 201 | Created - Ресурс создан |
| 400 | Bad Request - Некорректный запрос |
| 404 | Not Found - Ресурс не найден |
| 500 | Internal Server Error - Ошибка сервера |

---

## Структура проекта

```
IOT/
├── handlers/          # HTTP обработчики
├── service/           # Бизнес-логика
├── storage/           # Слой доступа к данным (PostgreSQL)
├── models/            # Модели данных
├── server/            # HTTP сервер
├── main.go            # Точка входа
├── go.mod
├── go.sum
└── README.md
```

---

## Логирование

Система логирует все изменения веса в таблице `weight_logs` для аудита и анализа.

---

## Статусы полок

| Статус | Описание |
|--------|---------|
| EMPTY | Полка пуста (нет товара) |
| OK | Вес в норме |
| REFILL | Требуется пополнение (вес ниже минимума) |

---

## Обработка ошибок

API возвращает ошибки в формате:
```json
{
    "error": "Описание ошибки",
    "timestamp": "2026-06-01T21:09:10Z"
}
```

---
