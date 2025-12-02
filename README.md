# Blogo 📝

A fast, simple, and full-featured blogging REST API built with Go, PostgreSQL, and Redis caching.

## Features ✨

- **User Management**: Create users, follow/unfollow, view profiles
- **Blog Management**: Create, read, update, delete blog posts
- **Social Features**: Like posts, follow users, view followers/following
- **JWT Authentication**: Secure authentication using JWT tokens
- **Redis Caching**: Fast response times with intelligent caching
- **PostgreSQL Database**: Reliable and scalable data storage
- **RESTful API**: Clean and intuitive API design
- **No Over-engineering**: Simple, maintainable codebase without unnecessary frameworks

## Tech Stack 🛠️

- **Language**: Go 1.19+
- **Architecture**: Clean Architecture
- **Database**: PostgreSQL
- **Cache**: Redis
- **Router**: Gorilla Mux
- **Authentication**: JWT (golang-jwt/jwt)

## Prerequisites 📋

- Go 1.19 or higher
- PostgreSQL 12+
- Redis 6+ (optional, but recommended for caching)

## Installation 🚀

### 1. Clone the repository

```bash
git clone <your-repo-url>
cd blogo
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Setup PostgreSQL

Create a PostgreSQL database:

```sql
CREATE DATABASE blogo;
```

### 4. Setup Redis (Optional)

Install and start Redis:

```bash
# On Ubuntu/Debian
sudo apt-get install redis-server
sudo systemctl start redis

# On macOS
brew install redis
brew services start redis

# On Windows
# Download from https://redis.io/download
```

### 5. Configure environment variables

Copy the example environment file and update it with your settings:

```bash
cp env.example .env
```

Edit `.env` with your database credentials:

```env
PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=blogo

REDIS_ADDR=localhost:6379
REDIS_PASSWORD=

JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
```

### 6. Run the application

```bash
go run cmd/api/main.go
# or use Make
make run
```

The API will be available at `http://localhost:8080`

## API Documentation 📚

### Authentication

Most endpoints require authentication. Include the JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

### User Endpoints

#### Create New User
```http
POST /api/u/new
Content-Type: application/json

{
  "username": "johndoe",
  "email": "john@example.com",
  "display_name": "John Doe"
}
```

**Response:**
```json
{
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "display_name": "John Doe",
    "bio": "",
    "profile_image": "",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Get User by ID
```http
GET /api/u/{id}
```

**Response:**
```json
{
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "display_name": "John Doe",
    "bio": "Software developer",
    "profile_image": "https://example.com/image.jpg",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "stats": {
    "followers_count": 10,
    "following_count": 5,
    "blogs_count": 3
  }
}
```

#### Update User Profile (Authenticated)
```http
POST /api/u/{id}/manage
Authorization: Bearer <token>
Content-Type: application/json

{
  "display_name": "John Smith",
  "bio": "Full-stack developer",
  "profile_image": "https://example.com/newimage.jpg"
}
```

#### Follow/Unfollow User (Authenticated)
```http
POST /api/u/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "action": "follow"
}
```

**Action can be:** `follow` or `unfollow`

#### Get User's Followers
```http
GET /api/u/{id}/follows?limit=20&offset=0
```

#### Get Users a User is Following
```http
GET /api/u/{id}/following?limit=20&offset=0
```

### Blog Endpoints

#### Get All Blogs
```http
GET /api/b?limit=20&offset=0
```

**Response:**
```json
{
  "blogs": [
    {
      "id": 1,
      "title": "My First Blog Post",
      "description": "An introduction to my blog",
      "body": "Full blog content here...",
      "author_id": 1,
      "author": {
        "id": 1,
        "username": "johndoe",
        "display_name": "John Doe",
        ...
      },
      "likes_count": 5,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "limit": 20,
  "offset": 0
}
```

#### Get Blog by ID
```http
GET /api/b/{id}
```

#### Create New Blog (Authenticated)
```http
POST /api/b/new
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "My Awesome Blog Post",
  "description": "A short description",
  "body": "The full content of the blog post..."
}
```

#### Update Blog (Authenticated)
```http
POST /api/b/{id}/edit
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Updated Title",
  "description": "Updated description",
  "body": "Updated content..."
}
```

**Note:** You can only edit your own blog posts.

#### Delete Blog (Authenticated)
```http
POST /api/b/{id}/delete
Authorization: Bearer <token>
```

**Note:** You can only delete your own blog posts.

#### Like/Unlike Blog (Authenticated)
```http
POST /api/b/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "action": "like"
}
```

**Action can be:** `like` or `unlike`

#### Get Blog Likes
```http
GET /api/b/{id}/likes?limit=20&offset=0
```

**Response:**
```json
{
  "likes": [
    {
      "id": 1,
      "username": "johndoe",
      "display_name": "John Doe",
      ...
    }
  ],
  "limit": 20,
  "offset": 0
}
```

### Pagination

All list endpoints support pagination using query parameters:

- `limit`: Number of items per page (default: 20, max: 100)
- `offset`: Number of items to skip (default: 0)

Example: `/api/b?limit=10&offset=20`

## Database Schema 🗄️

### Users Table
```sql
- id (SERIAL PRIMARY KEY)
- username (VARCHAR, UNIQUE)
- email (VARCHAR, UNIQUE)
- display_name (VARCHAR)
- bio (TEXT)
- profile_image (VARCHAR)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
```

### Blogs Table
```sql
- id (SERIAL PRIMARY KEY)
- title (VARCHAR)
- description (TEXT)
- body (TEXT)
- author_id (INTEGER, FK -> users.id)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
```

### Followers Table
```sql
- id (SERIAL PRIMARY KEY)
- follower_id (INTEGER, FK -> users.id)
- following_id (INTEGER, FK -> users.id)
- created_at (TIMESTAMP)
- UNIQUE(follower_id, following_id)
```

### Likes Table
```sql
- id (SERIAL PRIMARY KEY)
- blog_id (INTEGER, FK -> blogs.id)
- user_id (INTEGER, FK -> users.id)
- created_at (TIMESTAMP)
- UNIQUE(blog_id, user_id)
```

## Caching Strategy 📦

The API uses Redis for caching to improve performance:

- **User data**: Cached for 15 minutes
- **Blog posts**: Cached for 10 minutes
- **Automatic invalidation**: Cache is invalidated when data is updated

If Redis is unavailable, the API will continue to work without caching.

## Project Structure 📁

This project follows **Clean Architecture** principles for better maintainability, testability, and scalability.

```
blogo/
├── cmd/                    # Application entry points
│   ├── api/main.go        # Main API server
│   └── seed/main.go       # Database seeding
├── internal/              # Private application code
│   ├── domain/           # Core business layer
│   │   ├── entity/       # Business entities (User, Blog)
│   │   └── repository/   # Repository interfaces
│   ├── usecase/          # Business logic use cases
│   ├── delivery/http/    # HTTP handlers & routing
│   └── infrastructure/   # External implementations
│       ├── database/     # PostgreSQL implementation
│       └── cache/        # Redis implementation
├── pkg/                   # Public reusable packages
│   ├── auth/             # JWT authentication
│   └── response/         # HTTP response helpers
├── scripts/              # Utility scripts
├── go.mod                # Go module dependencies
├── ARCHITECTURE.md       # Clean architecture documentation
└── README.md             # This file
```

For detailed architecture documentation, see [ARCHITECTURE.md](ARCHITECTURE.md)

## Error Handling 🚨

All endpoints return errors in the following format:

```json
{
  "error": "Error message description"
}
```

Common HTTP status codes:
- `200 OK`: Success
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Authentication required or invalid token
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

## Security 🔒

- **JWT Authentication**: Secure token-based authentication
- **Password hashing**: Use bcrypt for password hashing in production
- **SQL Injection Prevention**: Parameterized queries
- **CORS**: Configure CORS headers for production use
- **Rate Limiting**: Consider adding rate limiting for production

## Performance ⚡

- **Redis Caching**: Reduces database load and improves response times
- **Database Indexing**: Optimized indexes on foreign keys and frequently queried fields
- **Connection Pooling**: Efficient database connection management
- **Pagination**: All list endpoints support pagination to handle large datasets

## Development 🔧

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
go build -o blogo cmd/api/main.go
./blogo
# or use Make
make build
./blogo
```

### Docker Support (Optional)

Create a `Dockerfile`:

```dockerfile
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o blogo

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/blogo .
EXPOSE 8080
CMD ["./blogo"]
```

Build and run:

```bash
docker build -t blogo .
docker run -p 8080:8080 --env-file .env blogo
```

## Contributing 🤝

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License 📄

This project is licensed under the MIT License - see the LICENSE file for details.

## Roadmap 🗺️

- [ ] Add OAuth2 integration (Google, GitHub)
- [ ] Implement search functionality
- [ ] Add image upload support
- [ ] Add comments on blog posts
- [ ] Add tags/categories for blogs
- [ ] Add user verification system
- [ ] Add email notifications
- [ ] Add rate limiting middleware
- [ ] Add comprehensive test suite
- [ ] Add API documentation with Swagger

## Support 💬

If you have any questions or issues, please open an issue on GitHub.

---

**Built with ❤️ using Go**
