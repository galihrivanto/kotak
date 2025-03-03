<div align="center">
  <h1>Kotak</h1>
  
  Kotak is simple temporary email service
</div>

## **Setup**

### Requirements

- [Go](https://golang.org/dl/)
- [Node](https://nodejs.org/en/download/)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Bootstrap

```bash
./bootstrap
```

sometimes, we need to manually stop previously running docker container

```bash
./docker/down
```

### Build & Run Server

```bash
go build -o kotak
./kotak
```

### Build & Run Frontend

```bash
npm --prefix frontend run dev
```

### Database

Example config for database

#### PostgreSQL

```yaml
database:
  driver: postgres
  host: localhost
  port: 5432
  username: kotak
  password: kotak
  database: kotak
```

#### MySQL

```yaml
database:
  driver: mysql
  host: localhost
  port: 3306
  username: kotak
  password: kotak
  database: kotak
```

#### SQLite

```yaml
database:
  driver: sqlite
  database: kotak
```

### Example Config

Or you can copy from example config

```bash
cp config.yaml.example_sqlite config.yaml
```

or

```bash
cp config.yaml.example_postgres config.yaml
```




