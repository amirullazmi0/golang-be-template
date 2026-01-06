#!/bin/bash

echo "🚀 Starting Kratify Backend Monitoring Setup..."

# Check if docker-compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose not found. Please install docker-compose first."
    exit 1
fi

# Create logs directory
echo "📁 Creating logs directory..."
mkdir -p logs

# Create .env if not exists
if [ ! -f .env ]; then
    echo "📝 Creating .env file from .env.example..."
    cp .env.example .env
    echo "⚠️  Please update .env file with your configurations"
fi

# Start monitoring stack
echo "🐳 Starting Grafana, Loki, and Promtail..."
docker-compose up -d

# Wait for services to be ready
echo "⏳ Waiting for services to start..."
sleep 10

# Check services status
echo "📊 Checking services status..."
docker-compose ps

echo ""
echo "✅ Setup complete!"
echo ""
echo "📍 Access points:"
echo "   - Grafana: http://localhost:3000 (admin/admin)"
echo "   - Loki API: http://localhost:3100"
echo ""
echo "🔧 Next steps:"
echo "   1. Update your .env file"
echo "   2. Run: go run main.go"
echo "   3. Open Grafana and explore logs"
echo ""
echo "📚 Documentation: MONITORING.md"
