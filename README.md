
<img src="https://img.shields.io/badge/Golang-007d9c?style=for-the-badge&logo=Go&logoColor=ffffff"/>  <img src="https://img.shields.io/badge/++-369?style=for-the-badge&logo=C&logoColor=ffffff"/>  <img src="https://img.shields.io/badge/gRpc-399?style=for-the-badge&logo=&logoColor=000000"/> 





# Микросервис по анализу текста

## Проект анализа текста клиента
Проект предназначен для анализирования текстов по следующим параметрам:
1. Сложность чтения
2. Уровень воды в тексте
3. Настроение текста

## Установка 

1. Склонируйте репозиторий проекта с помощью команды:
```bash
git clone https://github.com/NeZorinEgor/microservice_1.0
```

2. Установите зависимости
```bash
go get github.com/go-sql-driver/mysql
go get github.com/gorilla/mux
go get github.com/ledongthuc/pdf
go get google.golang.org/grpc
```

## Использование

1. Перейдите в директорию проекта:
```bash
cd go_microservice
```

2. Создайте таблицу в базе данных
```sql
CREATE TABLE `states` (
  `id` int UNSIGNED NOT NULL,
  `title` varchar(50) NOT NULL,
  `reading` int UNSIGNED NOT NULL,
  `water` int UNSIGNED NOT NULL,
  `mood` varchar(20) NOT NULL
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```

4. Пропишите в директории проекта команду

```go
go run main.go
```
_____

## Запуск и сборка analysis_service

### О сервисе

Чтобы использовать сервис, нужнно подключить саму логику анализа, для этого нужно запустить его.

### Сборка

- *Если вы не хотите собирать логику проекта, и ее модернезировать перейдите к пункту запуск*

Следуйте следующим инструкциям в консоли:
```bash
git submodule update --init --recursive
mkdir analysis_microservice/build
cd analysis_microservice/build
cmake ..
```
Либо же собирайте `CMake` в любой другой IDE которая поддерживает это, не забудтье обновить сабмодули

### Запуск

Вы можете запустить ваше собранное приложение, либо брать последний `release` из данного репозитория.
Чтобы развернуть сервер все что вам нужно это IP и порт. Пример для `localhost` тестирования:
`./analysis_service 127.0.0.1 1111`

### Связь между go и возможно другими сервисами

Чтобы подключить сервиз анализа к своему проекту вам нужно вязть `.proto` из вложенных папок и сгенерировать
`gRPC` клиент для своего языка, после можно с легкостью использовать запросы и ответы сервиса. Описание протокола общения

```c++
service TextAnalysService  {
  rpc getResult (SettingsTextPB) returns (ResultParsingPB) {} // Основаное взаимодействие с сервисом
}

message SettingsTextPB {
  string text = 1; // Исходный текст на вход
}

message ResultParsingPB {

  enum Mood { //Настроение текста
    sad = 0;
    happy = 1;
    lovely = 2;
    terrible = 3;
    boring = 4;
  }

  int32 water_value = 1; // Уровень воды в текстк
  Mood mood = 2; 
  int32 hard_reading = 3; // Сложность чтения текста
}
```

### Зависимые библеотеки 

- `gRPC` и `gProtobuf` используются как сабмодули, также `protoc` нужной версии лежит в утилитах
- `sqlite` используется более развернутая библеотека для с++, также сабмодулем, для хэширования предыдщуших результатов. 

## Руководство пользователя

![Альтернативный текст](https://github.com/NeZorinEgor/microservice_1.0/blob/main/static/screen.gif?raw=true)

##
