# API Throttling Demo

This project implements a simple HTTP API with **rate limiting (throttling)** on two endpoints: `/foo` and `/bar`. Each endpoint uses a different rate-limiting algorithm and storage strategy.

---

## 📌 Features

- Two GET endpoints: `/foo` and `/bar`
- Clients are authenticated via the `Authorization: Bearer <client-id>` header.
- Client details and their configurable rate limits for each endpoint are stored persistently in a **MySQL database**.
  For example, two defined clients in the system are:
  - `client-1` with 2 rate limit
  - `client-2` with 3 rate limit
- The API dynamically loads the rate limits per client from the MySQL database to enforce throttling.
- This design allows easy updates of client rate limits without redeploying the application.
- Two rate-limiting algorithms and storage strategies:
  - `/bar`: Sliding Log algorithm with **in-memory storage**
  - `/foo`: Fixed Window algorithm with **persistent storage (Redis)**
- Dockerized and ready for deployment
- Basic tests included

---

## ⚙️ Configuration

- The default rate limiting window **duration** is **1 minute** (`1m`).
- You can change this duration by setting the environment variable `THROTTLE_WINDOW_DURATION` using a [Go duration string format](https://pkg.go.dev/time#ParseDuration).

## 🚀 Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) (v1.23+)
- [Docker](https://www.docker.com/)
- [Docker compose](https://docs.docker.com/compose/install/) (for local test)
- (Optional) Redis (included via Docker Compose)
- (Optional) MySql (included via Docker Compose)


## 🐳 Building Docker Image

You can build the Docker image of the API using the provided Makefile target:


    $  make build_docker

## Run Locally

    $  docker-compose up -d

## ☁️ Deployment

The API is deployed on an AWS server and accessible at: [http://16.171.60.125:8888](http://16.171.60.125:8888)

## 🔎 Example Test Request

Test the `/foo` endpoint with a client authorization header:

    $  curl -X GET "http://16.171.60.125:8888/foo" -H "Authorization: Bearer client-1"