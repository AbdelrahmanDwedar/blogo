# Clean Architecture Documentation 🏗️

This document describes the clean architecture implementation of Blogo API.

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Project Structure](#project-structure)
- [Layers](#layers)
- [Dependency Flow](#dependency-flow)
- [Key Principles](#key-principles)
- [Adding New Features](#adding-new-features)

## Architecture Overview

Blogo follows **Clean Architecture** principles, separating the codebase into distinct layers with clear responsibilities and dependencies flowing inward.

```
┌─────────────────────────────────────────────────────────────┐
│                      External Interfaces                     │
│                    (HTTP, CLI, gRPC, etc.)                  │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                   Delivery Layer (HTTP Handlers)             │
│         - User Handler, Blog Handler                         │
│         - Request/Response transformation                    │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                    Use Case Layer                            │
│         - Business Logic                                     │
│         - User Use Case, Blog Use Case                       │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                   Domain Layer (Core)                        │
│         - Entities (User, Blog)                              │
│         - Repository Interfaces                              │
│         - Business Rules                                     │
└──────────────────────▲──────────────────────────────────────┘
                       │
┌──────────────────────┴──────────────────────────────────────┐
│                Infrastructure Layer                          │
│         - Database Implementation                            │
│         - Cache Implementation                               │
│         - External Services                                  │
└─────────────────────────────────────────────────────────────┘
```

## Project Structure

```
blogo/
├── cmd/                                # Application entry points
│   ├── api/
│   │   └── main.go                     # Main API server
│   └── seed/
│       └── main.go                     # Database seeding utility
│
├── internal/                           # Private application code
│   ├── domain/                         # Core business layer (innermost)
│   │   ├── entity/                     # Business entities
│   │   │   ├── user.go                 # User entity
│   │   │   ├── blog.go                 # Blog entity
│   │   │   └── errors.go               # Domain errors
│   │   └── repository/                 # Repository interfaces
│   │       ├── user_repository.go      # User repository interface
│   │       ├── blog_repository.go      # Blog repository interface
│   │       └── cache_repository.go     # Cache repository interface
│   │
│   ├── usecase/                        # Application business rules
│   │   ├── user_usecase.go             # User business logic
│   │   └── blog_usecase.go             # Blog business logic
│   │
│   ├── delivery/                       # Interface adapters
│   │   └── http/                       # HTTP delivery layer
│   │       ├── handler.go              # Main handler
│   │       ├── user_handler.go         # User HTTP handlers
│   │       └── blog_handler.go         # Blog HTTP handlers
│   │
│   └── infrastructure/                 # External interfaces & frameworks
│       ├── database/                   # Database implementations
│       │   ├── postgres.go             # PostgreSQL connection
│       │   ├── user_repository.go      # User repository implementation
│       │   └── blog_repository.go      # Blog repository implementation
│       └── cache/                      # Cache implementations
│           └── redis.go                # Redis cache implementation
│
├── pkg/                                # Public libraries (reusable)
│   ├── auth/                           # Authentication utilities
│   │   └── jwt.go                      # JWT token handling
│   └── response/                       # HTTP response helpers
│       └── json.go                     # JSON response utilities
│
├── scripts/                            # Utility scripts
└── config/                             # Configuration (if needed)
```

## Layers

### 1. Domain Layer (`internal/domain/`)

**Purpose**: Contains the core business entities and business rules. This is the heart of the application.

**Responsibilities**:
- Define business entities (User, Blog)
- Define repository interfaces (contracts)
- Business rule validation
- Domain-specific errors

**Dependencies**: NONE - This layer has no dependencies on any other layer.

**Key Files**:
- `entity/user.go` - User entity with validation
- `entity/blog.go` - Blog entity with validation
- `entity/errors.go` - Domain-specific errors
- `repository/user_repository.go` - User repository interface
- `repository/blog_repository.go` - Blog repository interface

**Example**:
```go
// Entity with business rules
type User struct {
    ID          int64
    Username    string
    Email       string
    // ...
}

func (u *User) Validate() error {
    if u.Username == "" {
        return ErrInvalidUsername
    }
    return nil
}
```

### 2. Use Case Layer (`internal/usecase/`)

**Purpose**: Contains application-specific business logic. Orchestrates the flow of data between layers.

**Responsibilities**:
- Coordinate business operations
- Implement application-specific logic
- Cache management
- Transaction coordination

**Dependencies**: Depends on Domain layer (entities and repository interfaces)

**Key Files**:
- `user_usecase.go` - User-related business operations
- `blog_usecase.go` - Blog-related business operations

**Example**:
```go
type UserUseCase struct {
    userRepo  repository.UserRepository
    cacheRepo repository.CacheRepository
}

func (uc *UserUseCase) CreateUser(...) (*entity.User, string, error) {
    // 1. Create entity
    user := entity.NewUser(...)
    
    // 2. Validate
    if err := user.Validate(); err != nil {
        return nil, "", err
    }
    
    // 3. Save to database
    if err := uc.userRepo.Create(user); err != nil {
        return nil, "", err
    }
    
    // 4. Cache it
    uc.cacheRepo.SetUser(user, 15*time.Minute)
    
    // 5. Generate token
    token, _ := auth.GenerateToken(...)
    
    return user, token, nil
}
```

### 3. Delivery Layer (`internal/delivery/http/`)

**Purpose**: Handles HTTP requests and responses. Transforms external data to internal formats.

**Responsibilities**:
- Parse HTTP requests
- Validate request format
- Call use cases
- Format responses
- Handle HTTP-specific concerns (status codes, headers)

**Dependencies**: Depends on Use Case and Domain layers

**Key Files**:
- `handler.go` - Main handler setup
- `user_handler.go` - User HTTP endpoints
- `blog_handler.go` - Blog HTTP endpoints

**Example**:
```go
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    // 1. Parse request
    var req struct {
        Username string `json:"username"`
        // ...
    }
    json.NewDecoder(r.Body).Decode(&req)
    
    // 2. Call use case
    user, token, err := h.userUC.CreateUser(req.Username, ...)
    
    // 3. Handle response
    if err != nil {
        response.Error(w, http.StatusBadRequest, err.Error())
        return
    }
    
    response.Created(w, map[string]interface{}{
        "user": user,
        "token": token,
    })
}
```

### 4. Infrastructure Layer (`internal/infrastructure/`)

**Purpose**: Implements interfaces defined in the domain layer. Handles external concerns.

**Responsibilities**:
- Database operations
- Cache operations
- External API calls
- File system operations

**Dependencies**: Depends on Domain layer (implements repository interfaces)

**Key Files**:
- `database/postgres.go` - PostgreSQL connection
- `database/user_repository.go` - User repository implementation
- `database/blog_repository.go` - Blog repository implementation
- `cache/redis.go` - Redis cache implementation

**Example**:
```go
type UserRepository struct {
    db *PostgresDB
}

// Implements domain.repository.UserRepository interface
func (r *UserRepository) Create(user *entity.User) error {
    // Database-specific implementation
    return r.db.Client.QueryRow(`
        INSERT INTO users (...) VALUES (...)
        RETURNING id
    `, ...).Scan(&user.ID)
}
```

### 5. Package Layer (`pkg/`)

**Purpose**: Reusable utilities and libraries that can be used across the application or even in other projects.

**Responsibilities**:
- Common utilities
- Helper functions
- Shared types

**Dependencies**: Should have minimal dependencies

**Key Files**:
- `auth/jwt.go` - JWT token utilities
- `response/json.go` - JSON response helpers

## Dependency Flow

The **Dependency Rule** states: **Dependencies only point inward**.

```
Infrastructure  ──►  Domain  ◄──  Use Case  ◄──  Delivery
      │                                              │
      └──────────────────────────────────────────────┘
```

- **Domain** has no dependencies
- **Use Case** depends on Domain
- **Delivery** depends on Use Case
- **Infrastructure** depends on Domain (implements interfaces)

This means:
- Domain entities don't know about databases
- Use cases don't know about HTTP
- Database code implements domain interfaces

## Key Principles

### 1. Dependency Inversion Principle (DIP)

High-level modules don't depend on low-level modules. Both depend on abstractions.

```go
// Domain defines the interface
type UserRepository interface {
    Create(user *User) error
}

// Use case depends on the interface
type UserUseCase struct {
    repo UserRepository  // Interface, not concrete type
}

// Infrastructure implements the interface
type PostgresUserRepository struct {
    db *sql.DB
}

func (r *PostgresUserRepository) Create(user *User) error {
    // Implementation
}
```

### 2. Interface Segregation Principle (ISP)

Interfaces are defined by the consumer, not the implementer.

```go
// Domain layer defines what it needs
type UserRepository interface {
    Create(user *User) error
    GetByID(id int64) (*User, error)
}

// Infrastructure provides the implementation
```

### 3. Single Responsibility Principle (SRP)

Each layer has a single, well-defined responsibility:
- **Entities**: Business rules
- **Use Cases**: Application logic
- **Handlers**: HTTP concerns
- **Repositories**: Data access

### 4. Testability

Clean architecture makes testing easier:

```go
// Mock repository for testing
type MockUserRepository struct {
    users []*entity.User
}

func (m *MockUserRepository) Create(user *entity.User) error {
    m.users = append(m.users, user)
    return nil
}

// Test use case with mock
func TestCreateUser(t *testing.T) {
    mockRepo := &MockUserRepository{}
    useCase := NewUserUseCase(mockRepo, nil)
    
    user, _, err := useCase.CreateUser("test", "test@test.com", "Test")
    
    assert.NoError(t, err)
    assert.Equal(t, 1, len(mockRepo.users))
}
```

## Adding New Features

### Example: Adding a Comment Feature

#### 1. Domain Layer

Create entity and repository interface:

```go
// internal/domain/entity/comment.go
type Comment struct {
    ID        int64
    BlogID    int64
    UserID    int64
    Content   string
    CreatedAt time.Time
}

func (c *Comment) Validate() error {
    if c.Content == "" {
        return ErrInvalidContent
    }
    return nil
}

// internal/domain/repository/comment_repository.go
type CommentRepository interface {
    Create(comment *Comment) error
    GetByBlogID(blogID int64) ([]*Comment, error)
    Delete(id, userID int64) error
}
```

#### 2. Use Case Layer

Create business logic:

```go
// internal/usecase/comment_usecase.go
type CommentUseCase struct {
    commentRepo repository.CommentRepository
    blogRepo    repository.BlogRepository
}

func (uc *CommentUseCase) CreateComment(blogID, userID int64, content string) (*entity.Comment, error) {
    // Validate blog exists
    _, err := uc.blogRepo.GetByID(blogID)
    if err != nil {
        return nil, entity.ErrBlogNotFound
    }
    
    comment := entity.NewComment(blogID, userID, content)
    if err := comment.Validate(); err != nil {
        return nil, err
    }
    
    if err := uc.commentRepo.Create(comment); err != nil {
        return nil, err
    }
    
    return comment, nil
}
```

#### 3. Infrastructure Layer

Implement repository:

```go
// internal/infrastructure/database/comment_repository.go
type CommentRepository struct {
    db *PostgresDB
}

func (r *CommentRepository) Create(comment *entity.Comment) error {
    return r.db.Client.QueryRow(`
        INSERT INTO comments (blog_id, user_id, content, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `, comment.BlogID, comment.UserID, comment.Content, comment.CreatedAt).Scan(&comment.ID)
}
```

#### 4. Delivery Layer

Create HTTP handler:

```go
// internal/delivery/http/comment_handler.go
type CommentHandler struct {
    commentUC *usecase.CommentUseCase
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
    claims, _ := auth.GetUserFromContext(r)
    blogID, _ := getIDFromPath(r, "blog_id")
    
    var req struct {
        Content string `json:"content"`
    }
    json.NewDecoder(r.Body).Decode(&req)
    
    comment, err := h.commentUC.CreateComment(blogID, claims.UserID, req.Content)
    if err != nil {
        response.Error(w, http.StatusBadRequest, err.Error())
        return
    }
    
    response.Created(w, comment)
}
```

#### 5. Wire it up in main.go

```go
// cmd/api/main.go
commentRepo := database.NewCommentRepository(db)
commentUC := usecase.NewCommentUseCase(commentRepo, blogRepo)
commentHandler := deliveryHttp.NewCommentHandler(commentUC)

r.HandleFunc("/api/b/{blog_id}/comments", auth.AuthMiddleware(commentHandler.CreateComment)).Methods("POST")
```

## Benefits of This Architecture

### 1. **Testability**
- Easy to mock dependencies
- Test business logic without database
- Test handlers without starting server

### 2. **Maintainability**
- Clear separation of concerns
- Easy to find and update code
- Changes in one layer don't affect others

### 3. **Flexibility**
- Swap databases without changing business logic
- Add new delivery mechanisms (gRPC, CLI)
- Replace cache implementation easily

### 4. **Scalability**
- Independent layer scaling
- Microservices extraction
- Team organization by layers

### 5. **Independence**
- Framework independent
- Database independent
- UI independent
- Testable without external dependencies

## Best Practices

1. **Keep domain layer pure** - No external dependencies
2. **Use interfaces** - Define contracts, not implementations
3. **Test each layer** - Unit tests for domain, integration tests for infrastructure
4. **Keep use cases thin** - Orchestration only, no complex logic
5. **One entity per file** - Easy to navigate
6. **Group by feature** - Not by layer (can scale to feature-based structure)
7. **Error handling** - Domain errors in domain layer, infrastructure errors wrapped

## Common Patterns

### Repository Pattern
Abstracts data access logic

### Use Case Pattern
Encapsulates application-specific business logic

### Dependency Injection
Inject dependencies through constructors

### Interface Segregation
Small, focused interfaces

## Further Reading

- [The Clean Architecture by Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)

---

**Remember**: The goal of clean architecture is to create a system that is:
- **Independent of frameworks**
- **Testable**
- **Independent of UI**
- **Independent of database**
- **Independent of any external agency**

Keep the dependencies pointing inward, and you'll have a maintainable, scalable application! 🎯


