# Wallet Service

Сервис для управления балансом кошельков.

Поддерживает:
- создание кошелька
- пополнение
- списание
- проверку баланса

## Запуск приложения

### 1. Поднять систему

```bash
docker-compose up --build
````

Будут запущены:

* postgres — основная база данных
* postgres-test — база для интеграционных тестов
* app — приложение

Сервис будет доступен по адресу:

```
http://localhost:9000
```

---

### 2. Применить миграции

Перед первым запуском необходимо применить миграции к основной базе:

```bash
make migration-up
```

Для тестовой базы:

```bash
make test-migration-up
```

---

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

Возможные значения `operationType`:

* DEPOSIT
* WITHDRAW

---

### GET /api/v1/wallets/{walletId}

Получение текущего баланса кошелька.

---

## Тестирование

### Unit-тесты

```bash
make test-unit
```

---

### Интеграционные тесты

1. Убедиться, что запущена тестовая база:

```bash
docker-compose up -d postgres-test
```

2. Применить миграции:

```bash
make test-migration-up
```

3. Запустить тесты:

```bash
make test-integration
```
