# HTTP/1.1 Server From Scratch (Go)

This project is a **from-scratch implementation of core parts of the HTTP/1.1 protocol in Go**, built without using Go’s `net/http` package or any third-party HTTP libraries.

The goal of this project is to deeply understand how HTTP works at a low level — from parsing raw TCP data to constructing valid HTTP responses.

---

## 📌 Project Overview

The server listens on a TCP socket and manually processes incoming HTTP/1.1 requests by implementing protocol features step by step.

Instead of relying on abstractions, every part of the request lifecycle is handled explicitly:
- Reading raw bytes from the network
- Parsing HTTP structures
- Handling different transfer mechanisms
- Writing binary responses back to the client

This project was inspired by the **“Learn HTTP in Go”** track from [boot.dev], with all logic implemented manually for learning purposes.

---

## ⚙️ Implemented Features

- **TCP-based HTTP server**
- **Request Line Parser**
  - Method, path, and HTTP version
- **Header Parser**
  - Case-insensitive headers
  - Proper CRLF handling
- **Request Body Handling**
  - Content-Length based bodies
- **Chunked Transfer Encoding**
  - Chunk size parsing
  - Chunk termination handling
- **Binary Response Writing**
  - Correct status line formatting
  - Headers + body serialization
- **Protocol-compliant HTTP/1.1 responses**

---

## 🧠 What This Project Focuses On

- Understanding the **HTTP/1.1 specification**
- Low-level **network programming** in Go
- Parsing text-based protocols safely and correctly
- Handling edge cases in request formatting
- Avoiding framework abstractions to learn fundamentals

---

## 🛠️ Tech Stack

- **Language:** Go (Golang)
- **Networking:** `net` package (TCP)
- **Protocol:** HTTP/1.1

---

## ▶️ How to Run

```bash
go run .
