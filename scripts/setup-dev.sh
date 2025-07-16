#!/bin/bash

# BezBase Development Setup Script

set -e

echo "🚀 Setting up BezBase development environment..."

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Create environment files if they don't exist
if [ ! -f backend/.env ]; then
    echo "📝 Creating backend environment file..."
    cp backend/.env.example backend/.env
fi

if [ ! -f frontend/.env ]; then
    echo "📝 Creating frontend environment file..."
    cp frontend/.env.example frontend/.env
fi

# Start the development environment
echo "🐳 Starting development environment..."
docker-compose -f docker-compose.dev.yml up -d

# Wait for database to be ready
echo "⏳ Waiting for database to be ready..."
sleep 10

# Check if services are running
echo "🔍 Checking service status..."
docker-compose -f docker-compose.dev.yml ps

echo "✅ Development environment setup complete!"
echo ""
echo "🌐 Services available at:"
echo "  - Backend API: http://localhost:8080"
echo "  - Database: localhost:5432"
echo ""
echo "📚 Next steps:"
echo "  1. Install frontend dependencies: cd frontend && npm install"
echo "  2. Start frontend development server: npm start"
echo "  3. Frontend will be available at: http://localhost:3000"
echo ""
echo "🛑 To stop the development environment:"
echo "  docker-compose -f docker-compose.dev.yml down"

