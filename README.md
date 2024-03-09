# go-webserver-template
Team Support Peru webpage

### Prerequisites
* Go
* Node and npm
* PostgreSQL
* [Air](https://github.com/cosmtrek/air#installation)
* [Templ](https://templ.guide/quick-start/installation)
* inotify-tools

### Initialize the required tables
```shell
$ psql -d <database_name> -U <username> -f ./db/init.sql
```

### .env file example
```
PORT=8080
DATABASE_URL=postgres://<username>:<password>@localhost:5432/<database_name>
```

### Then
```shell
$ npm install
```

### Live reload
```shell
$ ENV=development make live
```
