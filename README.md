# 🚀 GoChat - Real-time Distributed Chat Application

A production-ready, scalable real-time chat application with a stunning interface, built with Go, React 19, WebSocket, Kafka, Redis, and PostgreSQL.

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![React](https://img.shields.io/badge/React-20232A?style=for-the-badge&logo=react&logoColor=61DAFB)](https://react.dev/)
[![TypeScript](https://img.shields.io/badge/TypeScript-007ACC?style=for-the-badge&logo=typescript&logoColor=white)](https://www.typescriptlang.org/)
[![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white)](https://redis.io/)
[![Kafka](https://img.shields.io/badge/Apache_Kafka-231F20?style=for-the-badge&logo=apache-kafka&logoColor=white)](https://kafka.apache.org/)

## ✨ Features

### 🎨 Frontend
- **Modern UI/UX**: Premium glassmorphism design with dark/light themes
- **Real-time Chat**: Instant messaging with WebSocket
- **Typing Indicators**: See when others are typing
- **Online Presence**: Real-time user status
- **Message History**: Load previous conversations
- **Multiple Rooms**: Create and join chat rooms
- **Media Support**: File sharing (images, videos, audio)
- **Responsive Design**: Works on mobile, tablet, and desktop
- **Animations**: Smooth transitions and effects
- **Mention System**: @mentions with suggestions and highlights
- **SuperChat Integration**: AI-powered chatbot with Gemini API

### ⚡ Backend
- **Distributed Architecture**: Multiple server instances
- **Load Balancing**: Nginx with least connections algorithm
- **WebSocket**: Real-time bidirectional communication
- **Message Persistence**: Kafka + PostgreSQL
- **Pub/Sub**: Redis for real-time message distribution
- **JWT Authentication**: Secure user authentication
- **Rate Limiting**: Protection against abuse
- **Health Checks**: Monitoring endpoints
- **Graceful Shutdown**: Clean service termination

### 🏗️ Infrastructure
- **Docker Compose**: Easy deployment
- **Horizontal Scaling**: Add more servers easily
- **Message Queue**: Kafka for reliable message delivery
- **Caching**: Redis for performance
- **Reverse Proxy**: Nginx for routing and load balancing
- **Environment Configuration**: Comprehensive environment management

## 🚀 Quick Start

### Prerequisites
- Docker 20.10+
- Docker Compose 2.0+
- 4GB+ RAM
- Node.js (for development)

### 1. Clone Repository
```bash
git clone https://github.com/yourusername/gochat.git
cd gochat
```

### 2. Environment Setup
```bash
# Copy sample environment
cp sample.env .env

# Configure your environment variables
# Update POSTGRES_PASSWORD, JWT_SECRET, GEMINI_API_KEY, etc.
```

### 3. Start Services
```bash
# Production deployment
docker-compose up -d

# Development mode with hot-reload
docker-compose -f docker-compose.dev.yml up -d
```

### 4. Access Application
Open [http://localhost](http://localhost) in your browser

That's it! 🎉

## 📁 Project Structure

```
gochat/
├── frontend/                 # React 19 + TypeScript + Vite
│   ├── public/               # Static assets
│   ├── src/
│   │   ├── components/      # Reusable UI components
│   │   ├── contexts/        # React contexts (WebSocket)
│   │   ├── pages/           # Page components
│   │   ├── services/        # API services
│   │   ├── types/           # TypeScript definitions
│   │   └── utils/           # Utility functions
│   ├── package.json         # Node dependencies
│   ├── tsconfig.json        # TypeScript configuration
│   └── Dockerfile           # Production build
│
├── src/                      # Go backend
│   ├── main.go              # Entry point
│   ├── auth/                # Authentication handlers and middleware
│   ├── chat/                # Chat business logic
│   ├── database/            # Database connections and migrations
│   ├── cache/               # Redis cache implementation
│   ├── kafka/               # Kafka producers/consumers
│   ├── middleware/          # Global middleware
│   ├── models/              # Data models
│   ├── services/            # Business logic
│   └── utils/               # Utility functions
│
├── kafka-consumer/          # Kafka consumer service
│   ├── main.go              # Message persistence
│   └── Dockerfile
│
├── go.mod                   # Go module definition
├── go.sum                   # Go dependency checksums
│
├── docker-compose.yml       # Production compose
├── docker-compose.dev.yml   # Development compose
├── nginx.conf               # Main nginx configuration
├── .env                     # Environment variables
├── .gitignore               # Git ignore rules
└── README.md                # This file
```

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        Nginx (Port 80)                                  │
│                  Reverse Proxy & Load Balancer                          │
└─────────────────────────────────────────────────────────────────────────┘
         │                    │                    │
         ▼                    ▼                    ▼
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│  Frontend    │    │  Backend API │    │  WebSocket   │
│ (React 19)   │    │ Load Balanced│    │Load Balanced │
└──────────────┘    └──────────────┘    └──────────────┘
         │                    │                    │
         └────────────────────┼────────────────────┘
                              │
         ┌───────────────────┼───────────────────┐
         ▼                   ▼                   ▼
    ┌────────┐          ┌────────┐          ┌────────┐
    │  App1  │          │  App2  │          │  App3  │
    │ (Asia) │          │(M.East)│          │(Europe)│
    └────────┘          └────────┘          └────────┘
         │                   │                   │
         └───────────────────┼───────────────────┘
                              │
         ┌───────────────────┼───────────────────┐
         ▼                   ▼                   ▼
    ┌──────────┐      ┌──────────┐      ┌──────────┐
    │PostgreSQL│      │  Redis   │      │  Kafka   │
    │   DB     │      │ Pub/Sub  │      │ Queue    │
    └──────────┘      └──────────┘      └──────────┘
         │                                    │
         └────────────────┐                   │
                          ▼                   ▼
                   ┌──────────────────────────────┐
                   │       Kafka Consumer         │
                   │  (Message Persistence)       │
                   └──────────────────────────────┘
```

## 🛠️ Tech Stack

### Frontend
- **React 19**: Latest React with concurrent features
- **TypeScript**: Type-safe development
- **Vite**: Lightning-fast build tool
- **Tailwind CSS**: Utility-first CSS framework
- **Lucide React**: Beautiful icons
- **WebSocket API**: Real-time communication
- **React Router**: Client-side routing

### Backend
- **Go 1.23+**: High-performance backend
- **Fiber**: Fast HTTP framework
- **Gorilla WebSocket**: WebSocket implementation
- **JWT-Go**: JWT authentication
- **GORM**: ORM for database operations
- **Go-Redis**: Redis client
- **Segment Kafka**: Kafka client

### Infrastructure
- **PostgreSQL 16**: Primary database
- **Redis 7**: Pub/Sub and caching
- **Apache Kafka 3.8**: Message queue
- **Nginx 1.27**: Reverse proxy and load balancer
- **Docker**: Containerization
- **Docker Compose**: Multi-container orchestration

## 🔧 Development

### Local Development (No Docker)

**Terminal 1 - Dependencies:**
```bash
docker-compose up -d postgres redis kafka zookeeper
```

**Terminal 2 - Backend:**
```bash
cd .
go run src/main.go
```

**Terminal 3 - Frontend:**
```bash
cd frontend
npm install
npm run dev
```

### Development with Docker (Hot Reload)

```bash
docker-compose -f docker-compose.dev.yml up -d
```

- Frontend: [http://localhost:5173](http://localhost:5173) (hot reload)
- Backend: [http://localhost:8080](http://localhost:8080)

### Build Configuration
```bash
# Build backend
go build -o main src/main.go

# Build frontend
cd frontend
npm install
npm run build

# Build all containers
docker-compose build
```

## 🐳 Docker Commands

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Rebuild services
docker-compose build --no-cache

# Scale backend
docker-compose up -d --scale app1=3

# Check status
docker-compose ps

# Restart specific service
docker-compose restart nginx
```

## 🔐 Environment Variables

Key variables in `.env`:

```env
# Frontend
VITE_API_URL=http://nginx:80        # API base URL
VITE_WS_HOST=localhost               # WebSocket host
VITE_WS_PORT=80                      # WebSocket port

# Backend
SERVER_PORT=8080                     # Server port
JWT_SECRET=your-super-secret-key    # JWT secret key
SERVER_NAME=local                    # Server name for clustering

# Database
POSTGRES_HOST=postgres               # PostgreSQL host
POSTGRES_PORT=5432                   # PostgreSQL port
POSTGRES_DATABASE=chat_db            # Database name
POSTGRES_USER=postgres               # Database user
POSTGRES_PASSWORD=your-password      # Database password

# Redis
REDIS_HOST=redis                     # Redis host
REDIS_PORT=6379                      # Redis port

# Kafka
KAFKA_HOST=kafka                     # Kafka host
KAFKA_PORT=9092                      # Kafka port
KAFKA_TOPIC=chat_messages            # Kafka topic
KAFKA_GROUP_ID=chat_consumer_group   # Kafka consumer group

# Cloudinary (for file uploads)
CLOUDINARY_CLOUD_NAME=your_cloud     # Cloudinary cloud name
CLOUDINARY_API_KEY=your_api_key      # Cloudinary API key
CLOUDINARY_API_SECRET=your_secret    # Cloudinary API secret
CLOUDINARY_UPLOAD_PRESET=preset_name # Cloudinary preset

# Gemini API (for SuperChat AI)
GEMINI_API_KEY=your_gemini_key       # Google Gemini API key

# Nginx
NGINX_PORT=80                        # Nginx port
NGINX_ENV=development               # Environment (production/development)
```

## 🚀 Deployment

### Production Deployment

**Update .env for production:**
```env
# Security
JWT_SECRET=production-jwt-secret-key
POSTGRES_PASSWORD=strong-production-password

# Production settings
NGINX_ENV=production
ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
```

**Build and deploy:**
```bash
# Build production images
docker-compose build --no-cache

# Start services
docker-compose up -d

# Check health
docker-compose ps
curl http://localhost/api/health
```

### Cloud Platforms
- **AWS**: EC2, ECS, or Elastic Beanstalk
- **Google Cloud**: Compute Engine or Cloud Run
- **DigitalOcean**: Droplets or App Platform
- **Azure**: Container Instances or App Service

## 📊 Monitoring

### Health Checks
```bash
# Frontend health
curl http://localhost/health

# API health
curl http://localhost/api/health

# Database health
curl http://localhost/api/health/database

# Database metrics
curl http://localhost/api/metrics/database
```

### Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f frontend
docker-compose logs -f app1

# Save logs
docker-compose logs > logs.txt
```

### Metrics
```bash
# Container stats
docker stats

# Disk usage
docker system df

# Network info
docker network inspect chat-network
```

## 🧪 Testing

### Frontend Tests
```bash
cd frontend
npm run test
```

### Backend Tests
```bash
go test ./...
```

### Integration Tests
```bash
# Start services
docker-compose up -d

# Run tests
./run-integration-tests.sh
```

## 🐛 Troubleshooting

### Common Issues

**Port 80 already in use:**
```bash
# Change NGINX_PORT in .env
NGINX_PORT=8080

# Restart
docker-compose down
docker-compose up -d
```

**WebSocket not connecting:**
```bash
# Check nginx logs
docker-compose logs nginx

# Check backend logs
docker-compose logs app1

# Verify WebSocket upgrade
curl -i -N -H "Connection: Upgrade" -H "Upgrade: websocket" http://localhost/ws/chat
```

**Database connection failed:**
```bash
# Check PostgreSQL
docker-compose exec postgres pg_isready

# Check logs
docker-compose logs postgres

# Restart database
docker-compose restart postgres
```

**Redis connection issues:**
```bash
# Test Redis connection
docker-compose exec redis redis-cli ping

# Check logs
docker-compose logs redis
```

**Kafka startup issues:**
```bash
# Check Kafka status
docker-compose exec kafka kafka-topics --bootstrap-server localhost:9092 --list

# Check logs
docker-compose logs kafka
```

### Performance Optimization
- Monitor CPU/memory usage with `docker stats`
- Adjust container resources in docker-compose files
- Use read replicas for PostgreSQL in production
- Enable Redis persistence for production
- Monitor Kafka partition distribution

## 🤝 Contributing

We welcome contributions! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests if applicable
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Development Guidelines
- Follow Go naming conventions and standards
- Write TypeScript with proper typing
- Use React best practices and hooks
- Maintain clean Docker configurations
- Update documentation when needed
- Write tests for new features

### Issue Templates
- Bug reports: Include steps to reproduce, expected vs actual, environment
- Feature requests: Include use cases and implementation ideas
- Security issues: Report privately first

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Go team for the excellent language
- React team for the innovative framework
- Fiber framework developers
- Tailwind CSS for the utility-first approach
- All open-source contributors whose work powers this project

## 🗺️ Roadmap

- [ ] End-to-end encryption
- [ ] File sharing improvements
- [ ] Voice/video calls
- [ ] Mobile apps (React Native)
- [ ] Desktop apps (Electron)
- [ ] Message reactions
- [ ] User profiles
- [ ] Private messaging
- [ ] Group chats
- [ ] Message search
- [ ] Notifications
- [ ] Themes customization
- [ ] Mentions with autocomplete
- [ ] Message threading
- [ ] Typing indicators enhancement

## ⭐ Star History

If you find this project useful, please consider giving it a star!

## 📞 Support

- **Documentation**: Check the docs folder
- **Issues**: GitHub Issues
- **Discussions**: GitHub Discussions
- **Email**: support@gochat.example.com
- **Community**: Join our Discord server

---

**Built with ❤️ using Go, React, and modern web technologies**

[⬆ Back to Top](#-gochat---real-time-distributed-chat-application)