# Sustainability Tracker User Service

This service handles user registration, login, and profile management using Golang with Gin, PostgreSQL, and JWT authentication.

## Setup

1. **Environment Variables**

   Set the following environment variables:
   - `DATABASE_URL` (e.g., `postgres://user:password@localhost:5432/sustainability?sslmode=disable`)
   - `JWT_SECRET`
   - `PORT` (optional, default is 8080)

2. **Database Migration**

   To run the database migrations (which create the `users` table), use:

   ```sh
   make migrate