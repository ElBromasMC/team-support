# go-webserver-template

Team Support Peru webpage

## Local environment

### Prerequisites

* Docker

### .env file example

```
USER_UID="1000" # It must match your current user UID
ENV=development
DB_NAME=team-support
DB_PASSWORD=qwerty$321
SESSION_KEY=qwerty$321
PORT=8080
REL=1
SMTP_HOSTNAME=mail.example.com
SMTP_USER=user
SMTP_PASS=pass
WEBSERVER_HOSTNAME=www.domain.tld
IZIPAY_STOREID=id
IZIPAY_APIKEY=apikey
COMPANY_EMAIL=company@domain.tld
BOOK_EMAIL=book@domain.tld
CAPTCHA_SITE_KEY=
CAPTCHA_SECRET_KEY=
```

### Live reload

```shell
$ bin/run-dev.sh
```

## Docker

### Prerequisites

* [Traefik](https://doc.traefik.io/traefik/getting-started/quick-start/)

### Docker compose .env file example

```
ENV=production
DB_NAME=<database_name>
DB_PASSWORD=<password>
SESSION_KEY=mysecretkey
PORT=8080
REL=1
SMTP_HOSTNAME=mail.example.com
SMTP_USER=<username>
SMTP_PASS=<password>
WEBSERVER_HOSTNAME=www.domain.tld
IZIPAY_STOREID=
IZIPAY_APIKEY=
COMPANY_EMAIL=
BOOK_EMAIL=
CAPTCHA_SITE_KEY=
CAPTCHA_SECRET_KEY=
```

