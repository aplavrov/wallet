# Wallet Service

Сервис для управления балансом кошельков.

Поддерживает:
- создание кошелька
- пополнение
- списание
- проверку баланса

## Запуск приложения

```bash
docker-compose up --build
````

Будут запущены:

* postgres — основная база данных (миграция применяется в коде)
* postgres-test — база для интеграционных тестов
* app — приложение

Сервис будет доступен по адресу:

```
http://localhost:9000
```

## API

### POST /api/v1/wallet

Выполнение операции пополнения или списания.

Пример запроса:

```json
{
  "walletId": "uuid",
  "operationType": "DEPOSIT",
  "amount": 100
}
```

Возможные значения `operationType` это "DEPOSIT" или "WITHDRAW".

### GET /api/v1/wallets/{walletId}

Получение текущего баланса кошелька.

### POST /api/v1/wallets

Создает кошелек и возвращает его UUID, добавлено для удобства тестирования.

## Тестирование

### Unit-тесты
Запускаются через Makefile:
```bash
make test-unit
```

### Интеграционные тесты
Сначала необходимо применить миграции, затем запустить через Makefile:
```bash
make test-migration-up
make test-integration
```

### Ручное тестирование
Примеры curl-запросов:

Создание кошелька
```bash
curl -X POST http://localhost:9000/api/v1/wallets
```
Пополнение кошелька
```bash
curl -X POST http://localhost:9000/api/v1/wallet -H "Content-Type: application/json" -d '{"walletId": "11111111-1111-1111-1111-111111111111", "operationType": "DEPOSIT", "amount": 1000}'
```
Снятие с кошелька
```bash
curl -X POST http://localhost:9000/api/v1/wallet -H "Content-Type: application/json" -d '{"walletId": "11111111-1111-1111-1111-111111111111", "operationType": "WITHDRAW", "amount": 100}'
```
Баланс кошелька
```bash
curl -X GET http://localhost:9000/api/v1/wallets/11111111-1111-1111-1111-111111111111
```
