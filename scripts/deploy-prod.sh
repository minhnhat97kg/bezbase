#!/bin/bash

# BezBase Production Deployment Script

set -e

echo "🚀 Deploying BezBase to production..."

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

# Check if production environment files exist
if [ ! -f backend/.env ]; then
    echo "❌ Backend .env file not found. Please create it from .env.example"
    exit 1
fi

if [ ! -f frontend/.env ]; then
    echo "❌ Frontend .env file not found. Please create it from .env.example"
    exit 1
fi

# Build and start production environment
echo "🏗️  Building production images..."
docker-compose build --no-cache

echo "🐳 Starting production environment..."
docker-compose up -d

# Wait for services to be ready
echo "⏳ Waiting for services to start..."
sleep 15

# Check if services are running
echo "🔍 Checking service status..."
docker-compose ps

# Test backend health
echo "🏥 Testing backend health..."
if curl -f http://localhost:8080/api/health > /dev/null 2>&1; then
    echo "✅ Backend is healthy"
else
    echo "❌ Backend health check failed"
    exit 1
fi

# Test frontend
echo "🌐 Testing frontend..."
if curl -f http://localhost:3000 > /dev/null 2>&1; then
    echo "✅ Frontend is accessible"
else
    echo "❌ Frontend is not accessible"
    exit 1
fi

echo "✅ Production deployment complete!"
echo ""
echo "🌐 Services available at:"
echo "  - Frontend: http://localhost:3000"
echo "  - Backend API: http://localhost:8080"
echo "  - Full application (via Nginx): http://localhost:80"
echo ""
echo "📊 View logs:"
echo "  docker-compose logs -f [service-name]"
echo ""
echo "🛑 To stop production environment:"
echo "  docker-compose down"

