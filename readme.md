# Library App Search Service

Distributed Library System — Search Service

---

## Overview

The Search Service is a dedicated search microservice responsible for providing fast and efficient chapter search capabilities within the distributed library platform.

Built with Go, the service exposes a REST API that allows clients to search for chapters using OpenSearch as its search engine and Redis as a caching layer.

It handles search query validation, pagination, result retrieval, and response generation while keeping search operations independent from the transactional services.

The service acts as the read-oriented search layer of the library ecosystem, providing a dedicated interface for retrieving indexed chapter data without directly accessing the Catalog Service database.

---

## Service Boundaries

### Provides

- Chapter search API
- Full-text search through OpenSearch
- Redis-based search caching
- Search query validation
- Paginated search results
- Search result metadata
- HTTP request handling
- Health and readiness checks

### Does Not Handle

- Book catalog transaction management
- Chapter creation and updates
- Search index maintenance
- Book borrowing and return operations
- Inventory management
- User authentication
- User account management
- Notification delivery
- Analytical aggregation

---

## Badges

<!-- Code Quality & Tests -->

[![Tests](https://github.com/damouu/library-app-search/actions/workflows/tests.yml/badge.svg?branch=test)](https://github.com/damouu/library-app-search/actions/workflows/tests.yml)

<!-- Coverage -->

[![Codecov](https://codecov.io/gh/damouu/library-app-search/branch/test/graph/badge.svg)](https://codecov.io/gh/damouu/library-app-search)

<!-- Docker -->

[![Docker Build](https://github.com/damouu/library-app-search/actions/workflows/build-and-publish.yml/badge.svg)](https://github.com/damouu/library-app-search/actions/workflows/build-and-publish.yml)
[![Docker Image](https://img.shields.io/docker/v/damou/library-app-search?label=docker&logo=docker)](https://hub.docker.com/r/damou/library-app-search)
[![Docker Pulls](https://img.shields.io/docker/pulls/damou/library-app-search?logo=docker)](https://hub.docker.com/r/damou/library-app-search)

<!-- Git / Version -->

[![Git Tag](https://img.shields.io/github/v/tag/damouu/library-app-search?logo=github)](https://github.com/damouu/library-app-search/tags)

<!-- Technology -->

![Go](https://img.shields.io/badge/Go-1.27-00ADD8?logo=go)
![OpenSearch](https://img.shields.io/badge/OpenSearch-integrated-005EB8?logo=opensearch)
![Redis](https://img.shields.io/badge/Redis-integrated-red?logo=redis)

---

## Responsibilities

- Process chapter search requests
- Validate search query parameters
- Execute full-text searches through OpenSearch
- Cache search results using Redis
- Handle pagination and result limits
- Return structured search responses
- Maintain search API availability
- Expose liveness and readiness endpoints
- Handle graceful application shutdown
- Manage HTTP request lifecycle

---

## Technology Stack

| Category              | Technology             |
|-----------------------|------------------------|
| Runtime               | Go 1.27                |
| HTTP Router           | Chi v5                 |
| Search Engine         | OpenSearch             |
| Search Client         | opensearch-go/v4       |
| Cache                 | Redis                  |
| Redis Client          | go-redis/v9            |
| Configuration         | godotenv               |
| Identifiers           | Google UUID            |
| Testing               | Go testing / miniredis |
| Static Analysis       | `go vet`               |
| Code Coverage         | Go Cover / Codecov     |
| Containerization      | Docker                 |
| Development           | Air / Delve            |
| CI/CD                 | GitHub Actions         |

---

## API Endpoints

### Search Operations

#### Search Chapters

```http
GET /search?q={query}&page=1&size=10
```

Searches for chapters matching the provided query and returns paginated results.

### Query Parameters

| Parameter | Required | Default | Constraints |
|-----------|----------|---------|-------------|
| `q`       | Yes      | —       | Non-empty search query |
| `page`    | No       | `1`     | Integer between 1 and 1000 |
| `size`    | No       | `10`    | Integer between 1 and 100 |

### Example Request

```http
GET /search?q=dragon&page=1&size=10
```

### Response Structure

The API returns a JSON response containing the matching chapters and pagination metadata.

```json
{
  "items": [],
  "page": 1,
  "size": 10,
  "total": 0,
  "total_pages": 0,
  "has_next": false,
  "has_prev": false
}
```

### Response Fields

| Field         | Description |
|---------------|-------------|
| `items`       | Matching chapter results |
| `page`        | Current page number |
| `size`        | Number of results per page |
| `total`       | Total number of matching results |
| `total_pages` | Total number of available pages |
| `has_next`    | Indicates whether another page exists |
| `has_prev`    | Indicates whether a previous page exists |

---

## Search Processing Flow

Client Request

↓

HTTP Router (Chi)

↓

Search Handler

↓

Query Validation

↓

Search Application Layer

↓

Redis Cache / OpenSearch

↓

Paginated JSON Response

---

## Search Architecture

The Search Service separates HTTP request handling from search infrastructure through dedicated application and infrastructure components.

```text
Client
  │
  ▼
HTTP Router
  │
  ▼
Search Handler
  │
  ▼
Search Use Case
  │
  ├── Redis Cache
  │
  └── OpenSearch
          │
          ▼
    Search Results
          │
          ▼
    JSON Response
```

This separation allows search query handling and infrastructure integrations to evolve independently.

---

## Caching Strategy

Redis is used as the caching layer to improve search API responsiveness and reduce unnecessary search engine queries.

The caching layer complements OpenSearch by providing fast access to previously requested search results.

### Responsibilities

- Search result caching
- Reduced repeated search operations
- Lower search engine query overhead
- Improved response performance

---

## OpenSearch Integration

OpenSearch provides the search infrastructure used to retrieve indexed chapter information.

The service communicates with OpenSearch through the official Go client.

### Client

```text
github.com/opensearch-project/opensearch-go/v4
```

### Configuration

```env
OPENSEARCH_URL=
OPENSEARCH_USERNAME=
OPENSEARCH_PASSWORD=
```

The Search Service operates as a read-oriented search API, while the Search Indexer is responsible for asynchronously maintaining the searchable chapter documents.

---

## Redis Integration

Redis provides the caching infrastructure used by the Search Service.

### Client

```text
github.com/redis/go-redis/v9
```

### Configuration

```env
REDIS_URL=
```

Redis is initialized during application bootstrap and gracefully closed during application shutdown.

---

## Health Checks

The service exposes dedicated health endpoints to support container orchestration and service availability monitoring.

### Liveness

```http
GET /health/live
```

Returns:

```text
200 OK
```

when the application process is running.

### Readiness

```http
GET /health/ready
```

Returns:

```text
200 OK
```

when the required dependencies are available.

Returns:

```text
503 Service Unavailable
```

when a required dependency is unavailable.

---

## HTTP Middleware

The HTTP router uses Chi middleware to improve request handling and operational reliability.

### Enabled Middleware

- Request ID generation
- HTTP request logging
- Panic recovery

These middleware components provide request correlation, execution visibility, and protection against unexpected handler panics.

---

## Build Quality

The CI pipeline enforces:

- Automated test execution
- Race detector validation
- Static analysis with `go vet`
- Coverage generation
- Codecov reporting
- Test failure issue creation
- Pull request validation
- Docker image publication
- GitOps deployment manifest updates

---

## Configuration

The service uses environment-driven configuration.

### Environment Variables

```env
SERVER_PORT=

OPENSEARCH_URL=
OPENSEARCH_USERNAME=
OPENSEARCH_PASSWORD=

REDIS_URL=

OTEL_SERVICE_NAME=
OTEL_EXPORTER_OTLP_ENDPOINT=
```

Configuration is loaded through environment variables using `godotenv`.

---

## Deployment Architecture

The Search Service is distributed as a Docker container.

The production Docker image uses a multi-stage build process to separate compilation from the runtime environment.

### Production Image

```text
Build Stage
    │
    ▼
Go 1.27 Builder
    │
    ▼
Compiled Search API
    │
    ▼
Debian Bookworm Slim
    │
    ▼
Production Container
```

### Docker Image

```text
damou/library-app-search
```

The GitHub Actions pipeline automatically builds and publishes the Docker image and updates the corresponding Kubernetes deployment manifest in the GitOps repository.

---

## Architectural Role

The Search Service represents the search and retrieval layer of the distributed library system.

It provides a dedicated read interface for retrieving chapter information through OpenSearch while using Redis to accelerate repeated search operations.

By separating search queries from catalog transactions and index maintenance, the service contributes to a loosely coupled architecture in which search infrastructure can evolve and scale independently from the core business services.

The Search Indexer maintains the searchable data projection, while the Search Service exposes it to frontend clients through a dedicated REST API.

---

## License

Not specified.
