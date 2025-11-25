# Go Microservices E-Commerce Platform

<img width="1113" height="628" alt="Microservices Architecture Diagram" src="https://github.com/user-attachments/assets/4a06160c-e7b1-4a88-95c7-78d538c85e68" />

A production-ready, scalable e-commerce backend built with **Go microservices architecture**, featuring gRPC inter-service communication, MongoDB persistence, JWT authentication, and Docker containerization.

## 🏗️ Architecture Overview

This project implements a **microservices-based architecture** with the following services:

- **Auth Service** - User authentication and authorization with JWT tokens
- **Product Service** - Product catalog management and inventory
- **Cart Service** - Shopping cart operations and management
- **Gateway Service** - API Gateway for routing and service orchestration

### Communication Patterns

- **gRPC** - High-performance inter-service communication
- **Protocol Buffers** - Efficient data serialization
- **REST API** - External client communication via Gateway

## 🚀 Features

### Authentication Service (Port: 8080)
- ✅ User registration with password hashing (bcrypt)
- ✅ JWT-based authentication
- ✅ Token validation and refresh
- ✅ MongoDB user persistence
- ✅ gRPC service for token validation

### Product Service (Port: 8081)
- ✅ Product CRUD operations
- ✅ Inventory management
- ✅ Product search and filtering
- ✅ MongoDB product catalog
- ✅ gRPC endpoints for internal communication

### Cart Service (Port: 8082)
- ✅ Add/Remove items from cart
- ✅ Update cart quantities
- ✅ Cart persistence
- ✅ User-specific cart management
- ✅ Integration with Product service via gRPC

### Gateway Service
- ✅ Unified API endpoint
- ✅ Request routing to microservices
- ✅ Load balancing
- ✅ Authentication middleware

## 🛠️ Tech Stack

### Core Technologies
- **Language**: Go 1.25.3
- **Database**: MongoDB
- **Communication**: gRPC, Protocol Buffers
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Containerization**: Docker, Docker Compose

### Key Dependencies
```go
- google.golang.org/grpc v1.76.0          // gRPC framework
- go.mongodb.org/mongo-driver v1.17.6     // MongoDB driver
- github.com/golang-jwt/jwt/v5 v5.3.0     // JWT authentication
- golang.org/x/crypto v0.40.0             // Password hashing
- github.com/joho/godotenv v1.5.1         // Environment management
```

## 📁 Project Structure

```
go-microservice/
├── auth/                    # Authentication microservice
│   ├── cmd/                 # Application entry point
│   ├── internal/
│   │   ├── api/            # HTTP/gRPC handlers
│   │   ├── config/         # Configuration management
│   │   ├── helper/         # Utility functions
│   │   ├── intra/          # Inter-service communication
│   │   ├── model/          # Data models
│   │   ├── repository/     # Database layer
│   │   └── service/        # Business logic
│   ├── Dockerfile
│   ├── go.mod
│   └── .env
│
├── product/                 # Product microservice
│   ├── cmd/
│   ├── internal/
│   │   ├── api/
│   │   ├── config/
│   │   ├── intra/
│   │   ├── model/
│   │   ├── repository/
│   │   └── service/
│   ├── Dockerfile
│   └── go.mod
│
├── cart/                    # Cart microservice
│   ├── cmd/
│   ├── internal/
│   │   ├── api/
│   │   ├── config/
│   │   ├── errors/
│   │   ├── intra/
│   │   ├── middleware/
│   │   ├── model/
│   │   ├── pkg/
│   │   ├── repository/
│   │   └── service/
│   ├── Dockerfile
│   └── go.mod
│
├── gateway/                 # API Gateway
│   ├── cmd/
│   ├── internal/
│   ├── Dockerfile
│   └── go.mod
│
├── proto/                   # Protocol Buffer definitions
│   ├── product.proto
│   ├── showProduct.proto
│   └── validateToken.proto
│
├── pb/                      # Generated protobuf code
├── docker-compose.yaml      # Multi-container orchestration
└── README.md
```

## 🔧 Installation & Setup

### Prerequisites
- Go 1.25.3 or higher
- Docker & Docker Compose
- MongoDB (or use Docker)
- Protocol Buffer Compiler (protoc)

### 1. Clone the Repository
```bash
git clone https://github.com/sachinggsingh/go-microservice.git
cd go-microservice
```

### 2. Environment Configuration

Create `.env` files for each service:

**auth/.env**
```env
MONGO_URI=mongodb://localhost:27017
DB_NAME=ecommerce
JWT_SECRET=your-secret-key-here
PORT=8080
```

**product/.env**
```env
MONGO_URI=mongodb://localhost:27017
DB_NAME=ecommerce
PORT=8081
```

**cart/.env**
```env
MONGO_URI=mongodb://localhost:27017
DB_NAME=ecommerce
PORT=8082
```

### 3. Generate Protocol Buffers

```bash
# Install protoc-gen-go and protoc-gen-go-grpc
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate Go code from proto files
protoc --go_out=. --go-grpc_out=. proto/*.proto
```

### 4. Run with Docker Compose

```bash
# Build and start all services
docker-compose up --build

# Run in detached mode
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down
```

### 5. Run Locally (Development)

```bash
# Terminal 1 - Auth Service
cd auth
go mod download
go run cmd/main.go

# Terminal 2 - Product Service
cd product
go mod download
go run cmd/main.go

# Terminal 3 - Cart Service
cd cart
go mod download
go run cmd/main.go

# Terminal 4 - Gateway Service
cd gateway
go mod download
go run cmd/main.go
```

## 📡 API Documentation

### Authentication Endpoints

#### Register User
```http
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword",
  "name": "John Doe"
}
```

#### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword"
}

Response:
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {...}
}
```

### Product Endpoints

#### Get All Products
```http
GET /api/products
Authorization: Bearer <token>
```

#### Get Product by ID
```http
GET /api/products/:id
Authorization: Bearer <token>
```

#### Create Product
```http
POST /api/products
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Product Name",
  "description": "Product Description",
  "price": 99.99,
  "stock": 100
}
```

### Cart Endpoints

#### Get Cart
```http
GET /api/cart
Authorization: Bearer <token>
```

#### Add to Cart
```http
POST /api/cart/add
Authorization: Bearer <token>
Content-Type: application/json

{
  "product_id": "product_id_here",
  "quantity": 2
}
```

#### Update Cart Item
```http
PUT /api/cart/update
Authorization: Bearer <token>
Content-Type: application/json

{
  "product_id": "product_id_here",
  "quantity": 5
}
```

#### Remove from Cart
```http
DELETE /api/cart/remove/:product_id
Authorization: Bearer <token>
```

## 🔒 Security Features

- **Password Hashing**: bcrypt with salt rounds
- **JWT Authentication**: Secure token-based auth
- **Environment Variables**: Sensitive data protection
- **Input Validation**: Request payload validation
- **CORS**: Cross-origin resource sharing configuration

## 🧪 Testing

```bash
# Run tests for all services
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific service
cd auth
go test ./internal/...
```

## 🐳 Docker Configuration

Each microservice includes:
- **Multi-stage builds** for optimized image size
- **Non-root user** for security
- **.dockerignore** for efficient builds
- **Health checks** for container monitoring

## 📊 Monitoring & Logging

- Structured logging with contextual information
- Error tracking and handling
- Service health endpoints
- Request/Response logging

## 🚦 Service Ports

| Service | Port | Protocol |
|---------|------|----------|
| Auth    | 8080 | HTTP/gRPC |
| Product | 8081 | HTTP/gRPC |
| Cart    | 8082 | HTTP/gRPC |
| Gateway | TBD  | HTTP      |

## 🔄 Inter-Service Communication

Services communicate via **gRPC** for:
- Token validation (Auth → Cart/Product)
- Product details (Product → Cart)
- High-performance data exchange

## 🛣️ Roadmap

- [ ] Order Service implementation
- [ ] Payment gateway integration
- [ ] Redis caching layer
- [ ] Message queue (RabbitMQ/Kafka)
- [ ] Service mesh (Istio)
- [ ] Kubernetes deployment
- [ ] CI/CD pipeline
- [ ] API rate limiting
- [ ] Distributed tracing (Jaeger)
- [ ] Metrics collection (Prometheus/Grafana)

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 👨‍💻 Author

**Sachin Singh**
- GitHub: [@sachinggsingh](https://github.com/sachinggsingh)

## 🙏 Acknowledgments

- Go community for excellent libraries
- gRPC team for the framework
- MongoDB team for the driver
- Docker for containerization

---

**Built with ❤️ using Go and Microservices Architecture**
