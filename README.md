# WhatsApp Clone API

Hello, this is the API endpoint for the Whatsapp_Clone project that me and two others friends created. On this Golang code server, we created several endpoints that we use.

However, for me (AlatBekam), there are several endpoints and features that I created myself, namely:

## API Endpoint

### 1. Status

For status, I created four status-related endpoints.

- `GET 'api/public/users/statuses'` is used to retrieve all available statuses on the server.

- `GET 'api/private/users/statuses/'` is used to retrieve statuses that have previously been viewed by the user.

- `POST 'api/private/users/status'` is used to create a status.

- `POST 'api/private/users/status/view'` is used to change the status of that status.

### 2. Channels

For channels, I created two channel-related endpoints:

- `GET 'api/private/channels'` is used to retrieve all channels.

- `POST 'api/public/channels'` is used to create a channel.

### 3. Users

For users, I used four user-related endpoints:

- `GET 'api/public/users'` is used to retrieve users

- `GET 'api/public/users/:id'` is used to retrieve users based on their ID

- `POST 'api/public/users'` is used to add users

- `PUT 'api/private/users'` is used to edit users

### 4. Login Handler

Regarding user login, I created an endpoint

- `POST 'api/public/login'`

used as a login handler. This endpoint will generate tokens based on the expiration time specified in `GenerateJWT.go`.

In addition, I also created an Authentication Middleware that handles permissions for all private endpoints.

## The Application

If you want to see the mobile application that i made for this API, click [tWhatsapp_clone](https://github.com/AlatBekam/whatsapp_clone)

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

sincerely, **AlatBekam**
