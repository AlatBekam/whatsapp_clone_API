# WhatsApp Clone API

A Golang-based REST API that replicates core WhatsApp messaging functionality.

## Prerequisites

- Go 1.19 or higher
- Git

## Installation

1. Clone the repository:

```bash
git clone <repository-url>
cd whatsapp_clone_API
```

2. Install dependencies:

```bash
go mod download
```

## Activating the Server

1. Navigate to the project directory:

```bash
cd /C:/Belajar/Golang/whatsapp_clone_API
```

2. Run the server:

```bash
go run main.go
```

Or build and run:

```bash
go build
./whatsapp_clone_API
```

3. The server will start on `http://localhost:8080` (or your configured port)

## API Endpoints

- `GET /health` - Server health check
- `POST /messages` - Send a message
- `GET /messages/:id` - Retrieve message

## Environment Variables

Create a `.env` file in the root directory:

```
PORT=8080
DATABASE_URL=your_database_url
```

## License

MIT
