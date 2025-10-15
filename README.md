# Gateway Service 🌐

A Go-based API gateway for the real-time forum project that handles authentication, user management, chat functionality, and WebSocket connections for real-time communication.

<div align="center">

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![WebSocket](https://img.shields.io/badge/WebSocket-Real--time-4A90E2?style=for-the-badge&logo=websocket&logoColor=white)
![Microservices](https://img.shields.io/badge/Architecture-Microservices-FF6B6B?style=for-the-badge&logo=microgenetics&logoColor=white)
![API Gateway](https://img.shields.io/badge/API-Gateway-FFA500?style=for-the-badge&logo=amazonapigateway&logoColor=white)
![REST API](https://img.shields.io/badge/REST-API-25D366?style=for-the-badge&logo=rest&logoColor=white)

</div>

## ✨ Features

- **Authentication & Authorization**: User registration, login, and session management
- **Real-time Communication**: WebSocket support for live chat and notifications
- **API Gateway**: Routes requests to different microservices (auth, post, chat, comment)
- **CORS Support**: Cross-origin resource sharing middleware
- **Session Management**: Cookie-based session handling
- **Microservices Architecture**: Communicates with multiple backend services

## 📋 API Endpoints

### HTTP Endpoints
- `POST /register` - User registration
- `POST /login` - User authentication
- `GET /getUserData` - Get user profile data
- `GET /authorized` - Check authorization status
- `GET /getAllChats` - Retrieve all chat conversations

### WebSocket Endpoint
- `WS /ws` - Real-time WebSocket connection for chat and live updates

## 🏗️ Architecture

The gateway communicates with the following microservices:
- **Auth API** (port 8081) - Authentication and user management
- **Post API** (port 8082) - Forum posts and content
- **Chat API** (port 8083) - Chat and messaging
- **Comment API** (port 8084) - Comments and discussions

## 🛠️ Prerequisites

- Go 1.20 or higher
- Docker (optional, for containerized deployment)

## 📦 Installation

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd gateway
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Run the application**
   ```bash
   go run main.go
   ```

The gateway will start on `http://localhost:8080`

### Docker Deployment

1. **Build the Docker image**
   ```bash
   docker build -t gateway-service .
   ```

2. **Run the container**
   ```bash
   docker run -p 8080:8080 gateway-service
   ```

## 🔧 Configuration

The service configuration is located in `config/constants.go`:

```go
const URLauth = "http://authapi:8081"    // Authentication service
const URLPost = "http://postapi:8082"    // Post service
const URLChat = "http://chatapi:8083"    // Chat service
const URLComment = "http://commentapi:8084" // Comment service
const Port = "8080"                      // Gateway port
```

## 📁 Project Structure

```
gateway/
├── main.go                 # Application entry point
├── Dockerfile             # Docker configuration
├── go.mod                 # Go module dependencies
├── go.sum                 # Dependency checksums
├── config/
│   └── constants.go       # Configuration constants
├── internals/
│   ├── handlers/          # HTTP request handlers
│   │   ├── authorized.go
│   │   ├── getChats.go
│   │   ├── getUserData.go
│   │   ├── login.go
│   │   ├── register.go
│   │   └── ws/           # WebSocket handlers
│   │       ├── chat.go
│   │       ├── comment.go
│   │       ├── getUser.go
│   │       ├── logout.go
│   │       ├── post.go
│   │       └── websocket.go
│   ├── middleware/        # HTTP middleware
│   │   ├── cors.go
│   │   └── session.go
│   └── tools/
│       └── utils.go       # Utility functions
├── model/                 # Data models
│   ├── comment.go
│   ├── message.go
│   ├── post.go
│   ├── request.go
│   └── user.go
└── script/               # Build and deployment scripts
    ├── init.sh
    └── push.sh
```

## 🔌 WebSocket Communication

The gateway supports real-time communication through WebSocket connections. Clients can:
- Send and receive chat messages
- Get live updates on posts and comments
- Receive user status notifications
- Handle real-time forum interactions

### WebSocket Message Format
```json
{
  "action": "message_type",
  "data": "message_payload"
}
```

## 🛡️ Security Features

- **CORS Middleware**: Configured to handle cross-origin requests
- **Session Management**: Cookie-based session validation
- **Authorization Checks**: Middleware for protected endpoints
- **Input Validation**: Request data validation and sanitization

## 🚦 Health Check

The service runs on port 8080 by default. You can verify it's running by making a request to any endpoint or checking the WebSocket connection.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 Dependencies

- [gorilla/websocket](https://pkg.go.dev/github.com/gorilla/websocket) - WebSocket implementation
- [google/uuid](https://pkg.go.dev/github.com/google/uuid) - UUID generation
- Go standard library for HTTP handling

---

**Note**: This gateway service is designed to work with other microservices in the real-time forum ecosystem. Make sure all dependent services are running for full functionality.

---

<div align="center">

**⭐ Star this repository if you found it helpful! ⭐**

Made with ❤️ from 🇸🇳

</div>