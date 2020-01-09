# apiServer

Sample APP to obtain information from domain.

## Tecnologies

- DB
- Golang
- VueJS (no CLI)
- Bootstrap Vue

## API ENDPOINTS

### All Posts

- Path : `/log`
- Method: `GET`
- Response: `200`

### Create Post

- Path : `/newSearch`
- Method: `POST`
- Fields: `domain`
- Response: `201`

## Required Packages

- Dependency management
  - [dep](https://github.com/golang/dep)
- Database
  - [CockRoachDB - pq driver](https://github.com/lib/pq)
- Routing
  - [chi](https://github.com/go-chi/chi)

## Ports:

- [DB] 27257 & 25258
- [Go] 8005

## Quick Run Project

1. First clone the repo then go to apiServer folder. After that run .apiserver. Make sure you have installed and configured "go" and "dep" in your machine.
   FrontEnd is located on public/index

```
git clone https://github.com/luciagirasoles/apiServer.git

cd apiServer

dep ensure

go build . && .apiServer

```

2. Run CockRoachDB on apiServer file.

```
 apiServer/cockroach start --insecure --http-addr=localhost:26258 --background
```

Copy what is in servers.sql inside cockroackDB console and quit with "\q"

```
cockroach sql --insecure
```

3. Run index.html, it contains all required for styles and methods.
