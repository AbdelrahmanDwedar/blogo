# Contributing to Blogo 🤝

Thank you for considering contributing to Blogo! This document provides guidelines and instructions for contributing.

## Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Focus on what is best for the community
- Show empathy towards other community members

## How Can I Contribute?

### Reporting Bugs 🐛

Before creating bug reports, please check existing issues to avoid duplicates.

**When submitting a bug report, include:**

- Clear and descriptive title
- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- Environment details (OS, Go version, etc.)
- Relevant logs or error messages

### Suggesting Enhancements 💡

Enhancement suggestions are tracked as GitHub issues.

**When suggesting an enhancement:**

- Use a clear and descriptive title
- Provide a detailed description of the suggested enhancement
- Explain why this enhancement would be useful
- Include examples if applicable

### Pull Requests 🔀

1. Fork the repository
2. Create a new branch: `git checkout -b feature/my-new-feature`
3. Make your changes
4. Write or update tests if applicable
5. Ensure code passes linting
6. Commit your changes: `git commit -m 'Add some feature'`
7. Push to the branch: `git push origin feature/my-new-feature`
8. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.19+
- PostgreSQL
- Redis (optional)
- Make (optional, but recommended)

### Setup Steps

1. **Clone the repository**

```bash
git clone <repository-url>
cd blogo
```

2. **Install dependencies**

```bash
make install
# or
go mod download
```

3. **Setup environment**

```bash
cp env.example .env
# Edit .env with your settings
```

4. **Create database**

```bash
make db-create
# or
createdb blogo
```

5. **Seed database (optional)**

```bash
make seed
```

6. **Run the application**

```bash
make run
# or
go run .
```

## Coding Guidelines

### Go Style Guide

Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

**Key points:**

- Use `gofmt` for formatting
- Follow Go naming conventions
- Write meaningful comments
- Keep functions focused and small
- Handle errors appropriately
- Use meaningful variable names

### Project Structure

```
blogo/
├── auth/           # Authentication and JWT handling
├── cache/          # Redis caching layer
├── tables/         # Database models and operations
├── scripts/        # Utility scripts
├── api.go          # API handlers
└── main.go         # Application entry point
```

### Best Practices

1. **Error Handling**
   ```go
   if err != nil {
       return fmt.Errorf("descriptive message: %w", err)
   }
   ```

2. **Context Usage**
   ```go
   func (s *Service) DoSomething(ctx context.Context) error {
       // Use context for cancellation and timeouts
   }
   ```

3. **Database Operations**
   - Always use parameterized queries
   - Lock appropriately (RLock for reads, Lock for writes)
   - Close rows/statements after use

4. **API Responses**
   ```go
   respondWithJSON(w, http.StatusOK, data)
   respondWithError(w, http.StatusBadRequest, "message")
   ```

## Testing

### Writing Tests

```go
func TestSomething(t *testing.T) {
    // Setup
    // Execute
    // Assert
    // Cleanup
}
```

### Running Tests

```bash
make test
# or
go test -v ./...
```

## Documentation

- Update README.md for user-facing changes
- Add comments for exported functions and types
- Update API documentation in README.md for endpoint changes
- Add examples in api-examples.http for new endpoints

## Commit Messages

Use clear and meaningful commit messages:

```
type: subject

body (optional)

footer (optional)
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```
feat: add user profile image upload

fix: resolve database connection leak

docs: update API documentation for blog endpoints
```

## Pull Request Process

1. **Update Documentation**
   - Update README.md if needed
   - Add/update API examples
   - Update QUICKSTART.md if setup process changes

2. **Code Quality**
   - Run `go fmt ./...`
   - Fix any linter warnings
   - Ensure tests pass

3. **Describe Your Changes**
   - What does this PR do?
   - Why is this change needed?
   - Any breaking changes?
   - Screenshots (if UI changes)

4. **Review Process**
   - Address review comments
   - Keep discussions focused and professional
   - Update PR based on feedback

## Features We'd Love to See

### High Priority
- [ ] OAuth2 integration (Google, GitHub)
- [ ] Search functionality for blogs
- [ ] Rate limiting middleware
- [ ] Comprehensive test coverage

### Medium Priority
- [ ] Image upload support
- [ ] Comments on blog posts
- [ ] Tags/categories for blogs
- [ ] Email notifications

### Low Priority
- [ ] User verification system
- [ ] Admin dashboard
- [ ] Blog analytics
- [ ] Export/Import functionality

## Getting Help

- Open an issue for bugs or questions
- Check existing issues and documentation
- Join discussions on open issues

## Recognition

Contributors will be:
- Listed in the project README
- Credited in release notes
- Mentioned in project documentation

## License

By contributing, you agree that your contributions will be licensed under the same license as the project (MIT License).

---

Thank you for contributing to Blogo! 🎉


