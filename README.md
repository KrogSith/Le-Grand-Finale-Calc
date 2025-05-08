# Calc
# Golang calculator

Данная программа - веб-калькулятор: пользователь отправляет арифметическое выражение по HTTP и получает в ответ его результат, в противном случае - ошибку.

## Структура проекта
```bash
.
├── cmd
│   └── main.go
├── internal
│   ├── application
│   │   ├── application.go
│   │   └── handlers.go
│   ├── calculation
│   │   ├── calculation_test.go
│   │   └── calculation.go
│   └── modules
│       ├── objects.go
│       └── stack.go
├── go.mod
└── go.sum
```


## Запуск:
Для запуска введите в терминал команду: ```go run ./cmd/main.go```

Адрес сервера: ```http://localhost:8080```

Запустив программу, воспользуйтесь сервисом Postman и отошлите запрос вида:
```bash
{
    "expression": "выражение, которое вы хотите посчитать"
}
```

## Ограничения:
Калькулятор сломается, если выражение начинается с отрицательного числа


## Пример использования:
Адрес: ```http://localhost:8080/api/v1/calculate```

Запрос:
```bash
{
    "expression": "2+2"
}
```
Ответ:
```bash
{
  "id": "47fe0a5c-2beb-4f7b-9704-e9ccb145683c"
}
```


Адрес: ```http://localhost:8080/api/v1/calculate```

Запрос:
```bash
{
    "expression": "124-(253878*351753)/1251"
}
```
Ответ:
```bash
{
    "id": 81f7fde0-23f7-4747-9ec6-40a2cb456038
}
```

Адрес: ```http://localhost:8080/api/v1/calculate```

Запрос:
```bash
{
    "expression": "qwerty"
}
```
Ответ:
```bash
{
  "error": "Invalid expression"
}
```


Адрес: ```http://localhost:8080/api/v1/expressions```

Запрос:
```bash

```
Ответ:
```bash
{
  "expressions": [
    {
      "id": "47fe0a5c-2beb-4f7b-9704-e9ccb145683c",
      "expression": "2+2",
      "status": "OK",
      "result": 4
    },
    {
      "id": "81f7fde0-23f7-4747-9ec6-40a2cb456038",
      "expression": "124-(253878*351753)/1251",
      "status": "OK",
      "result": -71384646.69064748
    },
    {
      "id": "6c145245-90f0-4e14-ae71-ad5714977b35",
      "expression": "qwerty",
      "status": "Invalid expression",
      "result": 0
    }
  ]
}
```


Адрес: ```http://localhost:8080/api/v1/expressions/47fe0a5c-2beb-4f7b-9704-e9ccb145683c```

Запрос:
```bash

```
Ответ:
```bash
{
  "id": "47fe0a5c-2beb-4f7b-9704-e9ccb145683c",
  "expression": "2+2",
  "status": "OK",
  "result": 4
}
```