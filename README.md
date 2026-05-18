# Enterprise LMS (Learning Management System)

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Author](https://img.shields.io/badge/Author-Azhar%20Faturohman%20Ahidin-orange)](https://github.com/azharf99)

A robust, enterprise-grade Learning Management System (LMS) backend built with Go. This project provides a comprehensive suite of features for managing courses, enrollments, lessons, quizzes, and exams with a focus on security and scalability.

## 🚀 Features

- **User Management**: Secure authentication with JWT, role-based access control (Admin, Instructor, Student).
- **Course Management**: Create and manage courses with hierarchical structures (Courses -> Modules -> Lessons).
- **Enrollment System**: Track student progress and course access.
- **Assessment Engine**: 
  - Flexible Quiz and Exam systems.
  - Support for multiple question types.
  - Automated scoring and attempt tracking.
- **Analytics**: Performance tracking for exams and user progress.
- **Security First**: 
  - Global Security Headers.
  - Rate Limiting to prevent brute-force attacks.
  - CORS configuration.
  - BCrypt password hashing.
- **Modern Tech Stack**: Clean architecture using Repository and Usecase patterns.

## 🛠 Tech Stack

- **Language**: [Go (Golang)](https://golang.org/)
- **Web Framework**: [Gin Gonic](https://gin-gonic.com/)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: [PostgreSQL](https://www.postgresql.org/)
- **Authentication**: JWT (JSON Web Tokens)
- **Containerization**: [Docker](https://www.docker.com/) & Docker Compose

## 📋 Prerequisites

- Go 1.26 or higher
- PostgreSQL
- Docker & Docker Compose (optional but recommended)

## ⚙️ Installation & Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/azharf99/enterprise-lms.git
   cd enterprise-lms
   ```

2. **Environment Configuration**:
   Create a `.env` file in the root directory and configure your variables:
   ```env
   DB_HOST=localhost
   DB_USER=postgres
   DB_PASSWORD=yourpassword
   DB_NAME=lms_db
   DB_PORT=5432
   JWT_SECRET=your_secret_key
   ADMIN_USERNAME=admin
   ADMIN_PASSWORD=admin123
   ```

3. **Run with Docker Compose**:
   ```bash
   docker-compose up -d
   ```

4. **Run Manually**:
   ```bash
   go mod download
   go run cmd/api/main.go
   ```

The API will be available at `http://localhost:8080`.

## 📂 Project Structure

- `cmd/api`: Entry point for the application.
- `internal/config`: Database and environment configurations.
- `internal/delivery/http`: HTTP handlers and routing logic.
- `internal/delivery/http/middleware`: Auth, role, and security middlewares.
- `internal/domain`: Core domain models and entities.
- `internal/repository`: Data access layer (Postgres/GORM).
- `internal/usecase`: Business logic layer.
- `pkg/utils`: Common utility functions.

## 📄 License

This project is licensed under the **Apache License 2.0**.

**Attribution Requirement**: 
If you use this project, you **MUST** give credit to the original author:
**Azhar Faturohman Ahidin** ([github.com/azharf99](https://github.com/azharf99))

See the [LICENSE](LICENSE) file for more details.

---
Developed with ❤️ by [Azhar Faturohman Ahidin](https://github.com/azharf99)
