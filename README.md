# go-webserver-template
Team Support Peru webpage

## Local environment

### Prerequisites
* Go
* Node and npm
* Docker compose
* [Air](https://github.com/cosmtrek/air#installation)
* [Templ](https://templ.guide/quick-start/installation)
* inotify-tools

### Install build dependencies
```shell
$ npm install
```

### .env file example
```
ENV=development
WEBSERVER_HOSTNAME=www.domain.tld
DB_NAME=
DB_PASSWORD=
DATABASE_URL=postgres://postgres:${DB_PASSWORD}@localhost:5432/${DB_NAME}
SESSION_KEY=mysecretkey
PORT=8080
REL=1
SMTP_HOSTNAME=mail.example.com
SMTP_USER=<username>
SMTP_PASS=<password>
IZIPAY_STOREID=
IZIPAY_APIKEY=
COMPANY_EMAIL=
```

### Initialize the required tables
```shell
$ docker compose -f docker-compose.dev.yml up -d
```

### Load env variables
```shell
$ set -a
$ source .env
$ set +a
```

### Live reload
```shell
$ make live
```

## Docker

### Prerequisites
* [Traefik](https://doc.traefik.io/traefik/getting-started/quick-start/)

### Docker compose .env file example
```
ENV=production
WEBSERVER_HOSTNAME=www.domain.tld
DB_NAME=<database_name>
DB_PASSWORD=<password>
SESSION_KEY=mysecretkey
PORT=8080
REL=1
SMTP_HOSTNAME=mail.example.com
SMTP_USER=<username>
SMTP_PASS=<password>
IZIPAY_STOREID=
IZIPAY_APIKEY=
COMPANY_EMAIL=
```
