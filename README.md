# Cloud Notes - Collaborative Note-Taking Application

A real-time collaborative note-taking application built with React frontend and Go backend, featuring WebSocket-based live collaboration, user authentication, and note sharing capabilities.

## Features

- **User Authentication**: Secure signup/login with JWT tokens
- **Real-time Collaboration**: Multiple users can edit notes simultaneously
- **Note Management**: Create, edit, delete, and organize notes
- **Note Sharing**: Share notes with other users with read/write permissions
- **Responsive Design**: Works on desktop and mobile devices
- **Containerized Deployment**: Docker and Kubernetes ready

## Technology Stack

### Frontend
- React 19 with TypeScript
- React Router for navigation
- Axios for API calls
- WebSocket for real-time collaboration
- Vite for build tooling

### Backend
- Go with Gin framework
- PostgreSQL database
- GORM for database operations
- JWT for authentication
- WebSocket for real-time features
- Redis for session storage (optional)

### Infrastructure
- Docker for containerization
- Kubernetes for orchestration
- Nginx for load balancing and SSL termination

## Quick Start

### Prerequisites
- Node.js 18+
- Go 1.21+
- PostgreSQL 15+
- Docker (optional)

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd Cloud-Notes
   ```

2. **Start the backend**
   ```bash
   cd backend
   cp .env.example .env
   # Edit .env with your database credentials
   go mod download
   go run .
   ```

3. **Start the frontend**
   ```bash
   cd Frontend
   npm install
   npm run dev
   ```

4. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080/api

### Docker Deployment

1. **Using Docker Compose**
   ```bash
   docker-compose up -d
   ```

2. **Access the application**
   - Frontend: http://localhost:3000
   - Backend: http://localhost:8080

### Kubernetes Deployment

1. **Deploy to Kubernetes**
   ```bash
   kubectl apply -f k8s/
   ```

2. **Access via Ingress**
   - Add `cloud-notes.local` to your hosts file
   - Access: https://cloud-notes.local

## Project Structure

```
Cloud-Notes/
├── Frontend/                 # React frontend application
│   ├── src/
│   │   ├── components/      # React components
│   │   ├── context/         # React context providers
│   │   ├── services/        # API and WebSocket services
│   │   ├── types/           # TypeScript type definitions
│   │   └── ...
│   ├── Dockerfile
│   └── package.json
├── backend/                 # Go backend application
│   ├── controllers/         # HTTP handlers
│   ├── models/             # Database models
│   ├── middleware/         # HTTP middleware
│   ├── websocket/          # WebSocket handlers
│   ├── config/             # Configuration
│   ├── Dockerfile
│   └── go.mod
├── k8s/                    # Kubernetes manifests
├── docker-compose.yaml     # Docker Compose configuration
└── README.md
```

## API Documentation

### Authentication Endpoints
- `POST /api/auth/signup` - Create new user account
- `POST /api/auth/login` - User login
- `POST /api/auth/refresh` - Refresh access token

### Note Endpoints
- `GET /api/notes` - Get user's notes
- `POST /api/notes` - Create new note
- `GET /api/notes/:id` - Get specific note
- `PUT /api/notes/:id` - Update note
- `DELETE /api/notes/:id` - Delete note
- `POST /api/notes/:id/share` - Share note with user
- `GET /api/notes/:id/collaborators` - Get note collaborators

### WebSocket Endpoint
- `WS /ws/:noteId?user_id=:userId` - Real-time collaboration

## Real-time Collaboration

The application uses WebSocket connections to enable real-time collaborative editing:

1. **Operation-based Synchronization**: Changes are represented as operations (insert, delete, retain)
2. **Conflict Resolution**: Operations are applied in order with proper transformation
3. **User Presence**: Shows which users are currently editing a note
4. **Automatic Reconnection**: Handles connection drops gracefully

## Security Features

- **JWT Authentication**: Secure token-based authentication
- **Password Hashing**: Bcrypt for secure password storage
- **CORS Protection**: Configurable cross-origin request handling
- **Input Validation**: Server-side validation for all inputs
- **SQL Injection Prevention**: GORM provides safe query building

## Deployment Options

### Local Development
- Frontend dev server with hot reload
- Go backend with live reloading (using air)
- Local PostgreSQL database

### Docker
- Multi-stage builds for optimized images
- Docker Compose for local orchestration
- Separate containers for frontend, backend, and database

### Kubernetes
- Horizontal pod autoscaling
- Service discovery and load balancing
- Persistent storage for database
- Ingress with SSL termination

## Environment Variables

### Frontend (.env)
```
VITE_API_URL=http://localhost:8080/api
VITE_WS_URL=ws://localhost:8080
```

### Backend (.env)
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=collaborative_notes
JWT_SECRET=your-secret-key
PORT=8080
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
