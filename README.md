# Blogo

Blog REST API made using Golang

## TODO

1. Add dependencies
    - github.com/joho/godotenv
    - github.com/gorilla/mux
    - github.com/lib/pq
    - golang.org/x/oauth2
    - github.com/golang-jwt/jwt/v5
    - github.com/auth0/go-jwt-middleware/v2
1. Make API
    - `/api/u`
        -  `/api/u/new` (POST)
        -  `/api/u/{id}`
        -  `/api/u/{id}` (POST: for following)
        -  `/api/u/{id}/manage` (POST)
        -  `/api/u/{id}/following`
        -  `/api/u/{id}/follows`
    - `/api/b`
        - `/api/b/new` (POST)
        - `/api/b/{id}`
        - `/api/b/{id}` (POST: for liking)
        - `/api/b/{id}/edit` (POST)
        - `/api/b/{id}/likes`
1. Make tables/models
    - Users (use OAuth and JWT)
    - Blogs
