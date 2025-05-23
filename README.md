# YouTube Video Fetcher API

This is a Go-based API service that fetches and stores YouTube videos for a given search query. The service continuously fetches new videos in the background and provides a paginated API to access the stored videos.

## Features

- Background job to fetch latest YouTube videos
- Paginated API to access stored videos
- Support for multiple YouTube API keys
- PostgreSQL database for video storage
- RESTful API built with Gin framework

## Prerequisites

- Go 1.24 or higher
- PostgreSQL
- YouTube Data API v3 key(s)

## Setup

1. Clone the repository:
```bash
git clone https://github.com/yourusername/fampay-assignment.git
cd fampay-assignment
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file in the root directory with the following variables:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=youtube_videos
YOUTUBE_API_KEYS=key1,key2,key3
SEARCH_QUERY=your_search_query
```

4. Create a database named : `youtube_videos`
```sql
CREATE DATABASE youtube_videos

```
4. Run database migrations:
```bash
go run cmd/migrate/main.go
```

5. Start the server:
```bash
go run cmd/server/main.go
```

## API Endpoints

### GET /api/videos
Returns paginated list of videos sorted by published date in descending order.

Query Parameters:
- `page`: Page number (default: 1)
- `limit`: Number of items per page (default: 10)

## Project Structure

```
.
├── cmd/
│   ├── server/     # Main application entry point
│   └── migrate/    # Database migration scripts
├── internal/
│   ├── config/     # Configuration management
│   ├── database/   # Database models and queries
│   ├── handler/    # HTTP handlers
│   ├── service/    # Business logic
│   └── youtube/    # YouTube API client
├── migrations/     # SQL migration files
└── pkg/           # Shared packages
```

## License

MIT

