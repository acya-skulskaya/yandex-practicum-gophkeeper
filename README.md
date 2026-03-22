# Менеджер паролей GophKeeper
GophKeeper представляет собой клиент-серверную систему, позволяющую пользователю надёжно и безопасно хранить логины, пароли, бинарные данные и прочую приватную информацию.

# Функции

- клиент-серверное взаимодействие по gRPC, данные передаются по зашифрованному каналу
- данные хранятся в БД PostgreSQL; данные секретов хранятся в бд в зашифрованном виде; файлы с бинарными данными хранятся на диске в зашифрованном виде
- регистрация и аутентификация пользователей: JWT-токены для авторизации, в запросах передаются в gRPC metadata в ключе `authorization`; на стороне клиента токен сохраняется в [go keyring](https://github.com/zalando/go-keyring).
- клиентское приложение на основе [Cobra](https://github.com/spf13/cobra)

# Настройка и запуск

Клиент и сервер берут насторйки из окружения и/или файла .env

Для разработки можно сгенерировать файлы self-signed сертификатов, а так же отключить проверку подлинности сертификатов настройкой `GOPHKEEPER_TLS_SKIP_VERIFY`.
Так же поддерживается авто-выпуск сертификатов, настройка `GOPHKEEPER_TLS_AUTOCERT`

Для запуска сервера необходимо создать БД PostgreSQL и прописать настройки подключения к ней в переменной окружения `GOPHKEEPER_DB_DSN`. При первом запуске будут выполнены миграции.

По умолчанию клиент и сервер настроены на работу сервера на адресе `GOPHKEEPER_GRPC_SERVER_ADDRESS=":50051"`

## Настройка и запуск сервера

Пример файла с настройками для сервера `.env.example.server`

Запуск командой `make run-server`

## Настройка клиента

Пример файла с настройками для клиента `.env.example.client`

## Использование Gophkeeper
Все команды поддерживают как указание параметров сразу, так и интерактивную работу. Все аргументы, не указанные при вызове команты, будут запрошены дополнительно

### Регистрация
```shell
# register --login=NEW_USER_LOGIN --password=NEW_USER_PASSWORD

$ go run ./cmd/gophkeeper-client/main.go register
login is required, enter login: user1
password is required, enter password: 
trying to register as user1
successfully registered and logged in as user1
```

### Авторизация зарегистрированного пользователя
```shell
# login --login=YOUR_LOGIN --password=YOUR_PASSWORD

$ go run ./cmd/gophkeeper-client/main.go login
login is required, enter login: user3
password is required, enter password: 
trying to log in as user3
successfully logged in as user3
```

### Выход
```shell
# logout

$ go run ./cmd/gophkeeper-client/main.go logout
successfully logged out
```

### Инфомация о приложении
```shell
# version

$ go run ./cmd/gophkeeper-client/main.go version
BUILD INFO: 
VERSION: 0.0.1
DATE: 2026-03-21_16:25:48
COMMIT HASH: 463abae
```

### Создание секретов
```shell
# create 
#	--name=SECRET_NAME 
#	--type=SECRET_TYPE 
#	--path=PATH_TO_BINARY_FILE 
#	--metadata=METADATA 
#	--login=LOGIN_TO_SAVE 
#	--password=PASSWORD_TO_SAVE 
#	--bc-number=BANK_CARD_NUMBER 
#	--bc-cvv=BANK_CARD_CVV 
#	--bc-expiry=BANK_CARD_EXPIRY_DATE 
#	--bc-holder-name=BANK_CARD_HOLDER_NAME

$ go run ./cmd/gophkeeper-client/main.go create
name is required, enter name: github-creds
type is required, enter one of types: binary, text, bank_card, login_password: login_password
enter metadata or leave empty: website=github.com
login is required, enter login: acyaskulskaya
password is required, enter password: 
secret was saved with ID=27

$ go run ./cmd/gophkeeper-client/main.go create
name is required, enter name: a-poem
type is required, enter one of types: binary, text, bank_card, login_password: text
enter metadata or leave empty: 
secret text is required, enter text: Ty_ne_pej_iz_unitaza,/Tam_bacilly_i_zaraza!/Ty_snachala_vodu_slej,/Penu_sduj,_potom_uzh_pej!
secret was saved with ID=28

$ go run ./cmd/gophkeeper-client/main.go create
name is required, enter name: cert      
type is required, enter one of types: binary, text, bank_card, login_password: binary
enter metadata or leave empty: needs-key
file path to the binary file is required, enter path: ./resources/ssl/cert.pem
secret was saved with ID=29

$ go run ./cmd/gophkeeper-client/main.go create
name is required, enter name: kreditka
type is required, enter one of types: binary, text, bank_card, login_password: bank_card
enter metadata or leave empty: sber
card number is required, enter card number: 1234123412341234
card holder name is required, enter holder name: GOPH_GOPHEROVICH_GOPHEROV
CVV is required, enter CVV: 666
card expiry date is required, enter expiry date: 03/26
secret was saved with ID=30
```

### Обновление секретов
```shell
# update 
#	--id=SECRET_ID
#	--name=SECRET_NAME
#	--path=PATH_TO_BINARY_FILE 
#	--metadata=METADATA 
#	--login=LOGIN_TO_SAVE 
#	--password=PASSWORD_TO_SAVE 
#	--bc-number=BANK_CARD_NUMBER 
#	--bc-cvv=BANK_CARD_CVV 
#	--bc-expiry=BANK_CARD_EXPIRY_DATE 
#	--bc-holder-name=BANK_CARD_HOLDER_NAME

$ go run ./cmd/gophkeeper-client/main.go update
secret id is required, enter id: 29
enter a new secret name if you want to update it or leave empty: 
enter new metadata or leave empty: 
enter new path to the binary file: ./resources/ssl/cert.pem
secret 29 was updated
```

### Просмотр секрета
```shell
# get --id=SECRET_ID --version-id=SECRET_VERSION_ID

$ go run ./cmd/gophkeeper-client/main.go get
secret id is required, enter id: 29
secret version id is empty, you can specify version id, write 0 to get latest version data or leave empty to get a list of versions: 
———————————————————————————————————————————————————————————————————————————————————————————————————————————————————
ID    | TYPE            | NAME                 | CREATED                   | UPDATED                   | NUM.VER.
–––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––
29    | binary          | cert                 | 2026-03-21 15:44:07       | 2026-03-21 15:50:24       | 2
———————————————————————————————————————————————————————————————————————————————————————————————————————————————————
VERSIONS:
——————————————————————————————
ID    | CREATED
––––––––––––––––––––––––––––––
32    | 2026-03-21 15:50:24
30    | 2026-03-21 15:44:07
——————————————————————————————

$ go run ./cmd/gophkeeper-client/main.go get
secret id is required, enter id: 29
secret version id is empty, you can specify version id, write 0 to get latest version data or leave empty to get a list of versions: 0
———————————————————————————————————————————————————————————————————————————————————————————————————————————————————
ID    | TYPE            | NAME                 | CREATED                   | UPDATED                   | NUM.VER.
–––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––
29    | binary          | cert                 | 2026-03-21 15:44:07       | 2026-03-21 15:50:24       | 1
———————————————————————————————————————————————————————————————————————————————————————————————————————————————————
VERSION ID: 32
VERSION CREATE DATE: 2026-03-21 15:50:24
ORIGINAL FILE PATH: ./resources/ssl/cert.pem
FILE NAME: cert.pem
FILE SIZE: 2191
METADATA: needs-key
```

### Скачивание бинарного файла
```shell
# download --id=SECRET_ID --version-id=SECRET_VERSION_ID --path=LOCAL_PATH_TO_SAVE_FILE

$ go run ./cmd/gophkeeper-client/main.go download
secret id is required, enter id: 29
secret id is empty, you can specify it or leave empty: 
specify path to save file ot leave empty to save file to current working directory: ./
saved file to ./secret-29-version-32-cert.pem
```

### Список секретов
```shell
# list

$ go run ./cmd/gophkeeper-client/main.go list
———————————————————————————————————————————————————————————————————————————————————————————————————————————————————
ID    | TYPE            | NAME                 | CREATED                   | UPDATED                   | NUM.VER.
–––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––––
30    | bank_card       | kreditka             | 2026-03-21 15:47:08       | 2026-03-21 15:47:08       | 1
29    | binary          | cert                 | 2026-03-21 15:44:07       | 2026-03-21 15:59:17       | 3
28    | text            | a-poem               | 2026-03-21 15:40:13       | 2026-03-21 15:40:13       | 1
27    | login_password  | github-creds         | 2026-03-21 15:32:48       | 2026-03-21 15:32:48       | 1
———————————————————————————————————————————————————————————————————————————————————————————————————————————————————
```

### Удаление секрета
```shell
# delete --id=SECRET_ID --version-id=SECRET_VERSION_ID

$ go run ./cmd/gophkeeper-client/main.go delete
secret id is required, enter id: 29
version id is empty, you can specify it or leave empty: 
secret 29 was deleted
```

## Возможности `Makefile`

- `golangci-lint`, `golangci-lint-format` - запуск линтера и форматирование результата
- `build-gophkeeper` - сборка сервера
- `build-gophkeeper-client` - сборка клиента
- `build-gophkeeper-proto` - сборка proto-файлов
- `compile-all` - компиляция приложений для всех платформ
- `run-server` - запуск сервера



## Генерация тестовых сертификатов
Для генерации сертификатов для тестирования выполните команду

`openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes`

## Генерация моков

1. Установить mockery https://vektra.github.io/mockery/latest/installation/
2. Сгенерировать конфиг `$ mockery init github.com/vektra/mockery/v3/internal/fixtures`
3. Добавить в конфиг пакеты, для которых нужно сгенерировать моки

```yaml 
  packages:
    github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/server/repository/user:
      interfaces:
        UserRepositoryInterface:
```

4. Сгенерировать моки `$ mockery`