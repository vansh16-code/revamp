# Go Backend Boilerplate

A clean and minimal Go backend boilerplate with Gin, GORM, and PostgreSQL (Neon DB).

## Features

- **Gin** - Fast HTTP web framework
- **GORM** - ORM with auto-migration
- **PostgreSQL** - Neon DB cloud database
- **Docker** - Containerized deployment
- **RESTful API** - Clean API structure
- **Environment Variables** - Secure configuration

## Project Structure

```
.
├── config/          # Database configuration
├── handlers/        # Request handlers
├── models/          # Data models
├── routes/          # API routes
├── main.go          # Application entry point
├── .env             # Environment variables
├── Dockerfile       # Docker configuration
└── docker-compose.yml
```

## Getting Started

### Prerequisites

- Go 1.23+
- Docker & Docker Compose (optional)
- Neon DB account

### Installation

1. Clone the repository
```bash
git clone <your-repo>
cd proj
```

2. Install dependencies
```bash
go mod download
```

3. Configure environment variables
```bash
# Create .env file
DATABASE_URL='your-neon-db-connection-string'
```

4. Run the application
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## Docker Deployment

```bash
# Build and run
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop
docker-compose down
```

## API Endpoints

### Users
- `GET /api/users` - Get all users

### Posts
- `POST /api/posts` - Create a post
- `GET /api/posts` - Get all posts
- `GET /api/posts/:id` - Get a post by ID
- `PUT /api/posts/:id` - Update a post
- `DELETE /api/posts/:id` - Delete a post

## Example Requests

### Create a Post
```bash
curl -X POST http://localhost:8080/api/posts \
  -H "Content-Type: application/json" \
  -d '{"title":"My First Post","content":"Hello World","user_id":1}'
```

### Get All Posts
```bash
curl http://localhost:8080/api/posts
```

### Update a Post
```bash
curl -X PUT http://localhost:8080/api/posts/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated Title","content":"Updated content"}'
```

### Delete a Post
```bash
curl -X DELETE http://localhost:8080/api/posts/1
```

## Adding New Models

1. Create model in `models/`
```go
type Product struct {
    gorm.Model
    Name  string `json:"name" gorm:"not null"`
    Price float64 `json:"price"`
}
```

2. Add to auto-migration in `main.go`
```go
config.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Product{})
```

3. Create handlers in `handlers/`
4. Register routes in `routes/routes.go`

## Database Models

### User
- `id` - Primary key (auto-generated)
- `name` - User name
- `created_at` - Timestamp
- `updated_at` - Timestamp
- `deleted_at` - Soft delete timestamp

### Post
- `id` - Primary key (auto-generated)
- `title` - Post title
- `content` - Post content
- `user_id` - Foreign key to User
- `created_at` - Timestamp
- `updated_at` - Timestamp
- `deleted_at` - Soft delete timestamp

## GORM Model Convention

All models embed `gorm.Model` which automatically provides:
- `ID` - Auto-incrementing primary key
- `CreatedAt` - Auto-set on creation
- `UpdatedAt` - Auto-updated on save
- `DeletedAt` - Enables soft deletes

## Environment Variables

| Variable | Description |
|----------|-------------|
| `DATABASE_URL` | PostgreSQL connection string (Neon DB) |

## License

MIT
