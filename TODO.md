# Blogo Project Roadmap & TODO 🗺️

This document tracks planned features, enhancements, and improvements for the Blogo API.

## Priority Legend

- 🔴 High Priority
- 🟡 Medium Priority
- 🟢 Low Priority
- ✅ Completed

---

## High Priority Features 🔴

### OAuth2 Integration


- [ ] Google OAuth2 authentication
- [ ] GitHub OAuth2 authentication
- [ ] OAuth2 token management
- [ ] Social login endpoints

**Why**: Modern authentication expectations, easier user onboarding


### Search Functionality

- [ ] Full-text search for blog posts
- [ ] Search by title and content
- [ ] Search filters (author, date, tags)
- [ ] Search result pagination
- [ ] Search result highlighting

**Why**: Essential for content discovery


### Rate Limiting

- [ ] Rate limiting middleware
- [ ] Per-IP rate limits

**Why**: API protection and fair usage

---


## Medium Priority Features 🟡

### Comments System

- [ ] Add comments on blog posts
- [ ] Comment creation endpoint
- [ ] Get comments for a blog
- [ ] Edit/delete own comments
- [ ] Nested comments (replies)
- [ ] Comment likes


**Why**: Engagement and community building

**Implementation Notes**:

- Follow clean architecture pattern

- Add `Comment` entity in domain layer
- Create `CommentRepository` interface
- Implement in infrastructure layer

### Tags/Categories

- [ ] Add tags/categories for blogs
- [ ] Tag creation and management
- [ ] Get blogs by tag
- [ ] Tag search and autocomplete
- [ ] Popular tags endpoint
- [ ] Category hierarchy

**Why**: Content organization and discoverability

### Image Upload

- [ ] Image upload support
- [ ] Profile image upload
- [ ] Blog cover image upload
- [ ] Image compression
- [ ] Multiple image formats support
- [ ] CDN integration for images

**Why**: Rich content and better UX

### User Verification

- [ ] Email verification system
- [ ] Verification email sending
- [ ] Verification token management
- [ ] Resend verification email
- [ ] Verified badge on profiles


**Why**: Trust and security

---

## Low Priority Features 🟢

### Email Notifications

- [ ] Welcome email on signup
- [ ] New follower notification
- [ ] Blog like notification
- [ ] Comment notification
- [ ] Email preferences management
- [ ] Email templates

**Why**: User engagement

### API Documentation

- [ ] Swagger/OpenAPI documentation

**Why**: Better developer experience

---

## Testing & Quality 🧪

### Test Suite

- [ ] Unit tests for domain entities
- [ ] Unit tests for use cases
- [ ] Integration tests for repositories
- [ ] HTTP handler tests
- [ ] End-to-end API tests
- [ ] Test coverage reporting
- [ ] CI/CD integration

**Target**: 80%+ code coverage

---

## Infrastructure & DevOps 🚀

### Monitoring & Logging

- [ ] Structured logging
- [ ] Log levels configuration
- [ ] Request/response logging
- [ ] Error tracking (Sentry, etc.)
- [ ] Metrics collection (Prometheus)
- [ ] Performance monitoring
- [ ] Health check endpoints

### Deployment

- [ ] CI/CD pipeline (GitHub Actions)
- [ ] Automated testing in CI
- [ ] Automated deployment
- [ ] Database migrations

---

## Advanced Features 💡

### Analytics

- [ ] Blog view tracking
- [ ] User engagement metrics
- [ ] Popular content dashboard
- [ ] User retention metrics

### Content Management

- [ ] Draft blog posts
- [ ] Scheduled publishing
- [ ] Blog post revisions
- [ ] Blog post templates

### Admin Features

- [ ] Admin dashboard
- [ ] User management
- [ ] Content moderation

---

## Technical Debt & Refactoring 🔧

- [ ] Add request validation middleware
- [ ] Standardize error responses
- [ ] Add request ID tracing
- [ ] Improve error messages
- [ ] Code documentation improvements

---
