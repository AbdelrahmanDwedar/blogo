#!/bin/bash

echo "🚀 Blogo Setup Script"
echo "====================="
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.19 or higher."
    exit 1
fi

echo "✅ Go found: $(go version)"

# Check if PostgreSQL is available
if ! command -v psql &> /dev/null; then
    echo "⚠️  PostgreSQL client not found. Make sure PostgreSQL is installed and accessible."
else
    echo "✅ PostgreSQL client found"
fi

# Check if Redis is available
if ! command -v redis-cli &> /dev/null; then
    echo "⚠️  Redis not found. Redis is optional but recommended for caching."
else
    echo "✅ Redis found"
fi

echo ""
echo "📦 Installing Go dependencies..."
go mod download

if [ $? -eq 0 ]; then
    echo "✅ Dependencies installed successfully"
else
    echo "❌ Failed to install dependencies"
    exit 1
fi

echo ""
echo "📝 Setting up environment file..."

if [ ! -f .env ]; then
    if [ -f env.example ]; then
        cp env.example .env
        echo "✅ Created .env file from env.example"
        echo "⚠️  Please edit .env and update the configuration values!"
    else
        echo "❌ env.example not found"
        exit 1
    fi
else
    echo "⚠️  .env file already exists, skipping..."
fi

echo ""
echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Edit .env file with your database credentials"
echo "2. Create a PostgreSQL database: createdb blogo"
echo "3. (Optional) Start Redis: redis-server"
echo "4. Run the application: go run cmd/api/main.go"
echo ""
echo "For more information, see README.md"

