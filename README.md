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

