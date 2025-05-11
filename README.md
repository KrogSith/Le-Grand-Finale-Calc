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
│   │   ├── db.go
│   │   └── handlers.go
│   ├── calculation
│   │   ├── calculation_test.go
│   │   └── calculation.go
│   └── modules
│       └── stack.go
├── go.mod
└── go.sum
```


## Запуск:
Для запуска введите в терминал команду: ```go run ./cmd/main.go```

Для быстрого запуска в папке ```./cmd``` лежит ```cmd.exe```

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
Адрес: ```http://localhost:8080/api/v1/register```

Запрос:
```bash
{
    "name": "A",
    "password": "B"
}
```
Ответ:
```bash
{
  "message": "Registration succesful. Welcome, A!"
}
```


Адрес: ```http://localhost:8080/api/v1/register```

Запрос:
```bash
{
    "name": "A",
    "password": "abb"
}
```
Ответ:
```bash
{
  "error": "User with this name already exists"
}
```


Адрес: ```http://localhost:8080/api/v1/login```

Запрос:
```bash
{ 
    "name": "abba",
    "password": "1"
}
```
Ответ:
```bash
{
  "error": "User not found"
}
```


Адрес: ```http://localhost:8080/api/v1/login```

Запрос:
```bash
{ 
    "name": "A",
    "password": "B"
}
```
Ответ:
```bash
{
  "message": "Welcome, A!"
}
```


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
  "id": "2acabfe1-369d-4d8c-9047-333457912f9b"
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
  "id": "75ba65cb-b287-4c09-81c9-7d797bb1efa8"
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
  "error": "invalid expression"
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
      "id": "2acabfe1-369d-4d8c-9047-333457912f9b",
      "expression": "2+2",
      "userId": "09c2317c-fa53-4289-90f3-b646a8b7f6e4",
      "status": "OK",
      "result": 4
    },
    {
      "id": "75ba65cb-b287-4c09-81c9-7d797bb1efa8",
      "expression": "124-(253878*351753)/1251",
      "userId": "09c2317c-fa53-4289-90f3-b646a8b7f6e4",
      "status": "OK",
      "result": -71384646.69064748
    },
    {
      "id": "8c989a26-f0d3-4d1b-bc8b-9f8ff1ce8082",
      "expression": "qwerty",
      "userId": "09c2317c-fa53-4289-90f3-b646a8b7f6e4",
      "status": "invalid expression",
      "result": 0
    }
  ]
}
```


Адрес: ```http://localhost:8080/api/v1/expressions/75ba65cb-b287-4c09-81c9-7d797bb1efa8```

Запрос:
```bash

```
Ответ:
```bash
{
  "id": "75ba65cb-b287-4c09-81c9-7d797bb1efa8",
  "expression": "124-(253878*351753)/1251",
  "userId": "09c2317c-fa53-4289-90f3-b646a8b7f6e4",
  "status": "OK",
  "result": -71384646.69064748
}
```