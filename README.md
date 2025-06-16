# Extraction Manager

## Overview
The `Extraction Manager` is a production-ready Go application designed to process large JSON input files, extract specific fields, and write the results to multiple CSV output files. It implements efficient parallel processing, monitoring, and proper error handling.

## Features
- Efficient parallel processing with configurable worker pools
- Structured logging with different log levels
- Prometheus metrics for monitoring
- Health check endpoints
- Graceful shutdown handling
- Configuration validation
- Docker support for easy deployment
- Comprehensive error handling and recovery
- Automatic file rotation with configurable line limits
- Production-grade security practices

## Prerequisites
- Go 1.22 or later
- Docker and Docker Compose (for containerized deployment)
- Make (optional, for using Makefile commands)

## Project Structure
```
.
├── cmd/                    # Application entry points
├── config/                 # Configuration files
├── internal/              # Private application code
├── pkg/                   # Public libraries
│   ├── logger/           # Structured logging
│   ├── metrics/          # Prometheus metrics
│   └── health/           # Health checks
├── Dockerfile            # Container definition
├── docker-compose.yml    # Container orchestration
└── prometheus.yml        # Prometheus configuration
```

## Installation

### Local Development
1. Clone the repository:
```bash
git clone git@github.com:Shmulik27/assignment.git
cd assignment
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the application:
```bash
go build -o bin/assignment ./cmd/data_extraction
```

### Docker Deployment
1. Build and run using Docker Compose:
```bash
docker-compose up --build
```

## Configuration
The application can be configured through environment variables or the `config/app_configuration.json` file:

```json
{
  "inputFileName": "spins_input.json",
  "outputFileName": "output.csv",
  "numWorkers": 4,
  "linesPerFile": 2500,
  "linesChannelSize": 100,
  "resultsChannelSize": 100
}
```

Environment Variables:
- `LOG_LEVEL`: Logging level (debug, info, warn, error)
- `LOG_FORMAT`: Log format (json, text)
- `LOG_PATH`: Path to log file (optional)
- `CONFIG_FILE`: Path to configuration file

## Running the Application

### Local Run
```bash
go run ./cmd/data_extraction
```

### Docker Run
```bash
docker-compose up
```

## Monitoring and Observability
The application exposes several endpoints for monitoring:

- Health Check: `http://localhost:8080/health`
- Metrics: `http://localhost:8081/metrics`
- Prometheus: `http://localhost:9090`
- Grafana: `http://localhost:3000`

## Data Files
- Input: `spins_input.json` (not included in repository)
- Output: `output-*.csv` files (generated during processing)

Note: The input file is not included in the repository due to its size. You'll need to provide your own input file with the correct structure.

## Testing
Run the test suite:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## Production Deployment
1. Ensure all environment variables are properly set
2. Use Docker Compose for deployment:
```bash
docker-compose -f docker-compose.prod.yml up -d
```

3. Monitor the application:
   - Check health endpoint
   - Monitor metrics in Prometheus
   - Set up alerts in Grafana

## Security Considerations
- SSH key authentication for Git
- Proper file permissions
- Environment variable for sensitive data
- Container security best practices
- Regular dependency updates

## Contributing
1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License
[Your chosen license]

## Support
For issues and feature requests, please create an issue in the GitHub repository.