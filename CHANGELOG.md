# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-01-01

### Added
- Initial release of Blogo API
- User management system
  - Create users
  - Update user profiles
  - View user profiles with stats
- Authentication system
  - JWT-based authentication
  - Token generation and validation
  - Protected endpoints
- Blog management
  - Create, read, update, delete blog posts
  - List blogs with pagination
  - View blogs by author
- Social features
  - Follow/unfollow users
  - Like/unlike blog posts
  - View followers and following lists
  - View blog likes
- Database integration
  - PostgreSQL database support
  - Automatic table creation
  - Database connection pooling
  - Proper indexing for performance
- Redis caching
  - User profile caching
  - Blog post caching
  - Automatic cache invalidation
  - Graceful fallback if Redis unavailable
- RESTful API
  - Clean and intuitive endpoints
  - Consistent error handling
  - JSON responses
  - Pagination support
- Documentation
  - Comprehensive README
  - Quick start guide
  - Docker setup guide
  - API examples
  - Contributing guidelines
- Development tools
  - Makefile for common tasks
  - Docker and Docker Compose support
  - Database seeding script
  - Environment configuration example
  - API testing examples

### Database Schema
- Users table with profile information
- Blogs table with author relationship
- Followers table for user relationships
- Likes table for blog interactions
- Proper foreign keys and constraints
- Optimized indexes

### API Endpoints

#### User Endpoints
- `POST /api/u/new` - Create new user
- `GET /api/u/{id}` - Get user by ID
- `POST /api/u/{id}` - Follow/unfollow user
- `POST /api/u/{id}/manage` - Update user profile
- `GET /api/u/{id}/following` - Get users following
- `GET /api/u/{id}/follows` - Get user followers

#### Blog Endpoints
- `GET /api/b` - Get all blogs
- `POST /api/b/new` - Create new blog
- `GET /api/b/{id}` - Get blog by ID
- `POST /api/b/{id}` - Like/unlike blog
- `POST /api/b/{id}/edit` - Update blog
- `POST /api/b/{id}/delete` - Delete blog
- `GET /api/b/{id}/likes` - Get blog likes

### Security
- JWT token-based authentication
- Parameterized SQL queries
- Password environment variables
- Secure token validation
- User authorization checks

### Performance
- Redis caching for frequently accessed data
- Database connection pooling
- Efficient query optimization
- Proper database indexing
- Pagination for large datasets

### Development Experience
- Simple project structure
- Easy local setup
- Docker support for all services
- Comprehensive documentation
- Example API requests
- Database seeding for testing

## [Unreleased]

### Planned Features
- OAuth2 integration (Google, GitHub)
- Search functionality for blogs
- Image upload support
- Comments on blog posts
- Tags/categories for blogs
- Email notifications
- Rate limiting
- User verification system
- Comprehensive test suite
- API documentation with Swagger

---

## Version History

### Version 1.0.0
- First stable release
- All core features implemented
- Production-ready API
- Complete documentation

---

## Release Notes Format

Each release should include:

### Added
- New features

### Changed
- Changes in existing functionality

### Deprecated
- Soon-to-be removed features

### Removed
- Removed features

### Fixed
- Bug fixes

### Security
- Vulnerability fixes


