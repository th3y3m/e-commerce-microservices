[![CodeQL](https://github.com/th3y3m/e-commerce-microservices/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/th3y3m/e-commerce-microservices/actions/workflows/github-code-scanning/codeql)
![MIT License](https://img.shields.io/badge/License-MIT-yellow.svg)

# 🎉 Welcome to E-commerce Microservices Project 🎉

## 🌟 Overview
This project is an e-commerce platform built using microservices architecture. It leverages various technologies such as PostgreSQL, Redis, RabbitMQ, and Docker to ensure scalability, reliability, and performance.

## 🚀 Features
- **Microservices Architecture**: Each service is designed to be independent and scalable.
- **PostgreSQL**: Used as the primary database for storing user and order information.
- **Redis Cache**: Implemented for caching frequently accessed data to improve performance.
- **Email Verification and Order Notifications**: Utilizes `net/smtp` for sending verification emails and order notifications.
- **Social Login**: Supports login via Google and Facebook using `github.com/markbates/goth`.
- **REST API**: Provides a RESTful API for interaction with the services.
- **Service Communication**: Uses RabbitMQ for inter-service communication.
- **Payment Integration**: Integrated with MoMo and VNPay for payment processing.
- **Docker**: Configured to run all services in isolated containers.

## 🛠️ Getting Started

### Prerequisites
- Docker
- Docker Compose
- PostgreSQL
- Redis
- RabbitMQ

### Installation
1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/ecommerce-microservices.git
    cd ecommerce-microservices
    ```

2. Set up environment variables:
    ```sh
    cp .env.example .env
    ```

3. Update the `.env` file with your configuration details.

4. Build and run the services using Docker Compose:
    ```sh
    docker-compose up --build
    ```

### Usage
- Access the API at `http://localhost:8000/api`.
- Use the provided endpoints to interact with the services.

## 📦 Services

### User Service
- Manages user registration, authentication, and profile management.
- Supports social login via Google and Facebook.

### Order Service
- Handles order creation, updates, and tracking.
- Sends order notifications via email.

### Payment Service
- Integrates with MoMo and VNPay for payment processing.

### Notification Service
- Sends email notifications for account verification and order updates.

## 🔄 Communication
- Services communicate with each other using RabbitMQ.

## ⚡ Caching
- Redis is used to cache frequently accessed data to improve performance.

## 🗄️ Database
- PostgreSQL is used as the primary database for storing user and order information.

## 🚀 Deployment
- The project is containerized using Docker, making it easy to deploy and scale.

## 📜 License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments
- [Mark Bates](https://github.com/markbates/goth) for the `goth` library.
- [MoMo](https://momo.vn) and [VNPay](https://vnpay.vn) for payment integration.

#### 📧 Connect with me via: truongtanhuy3006@gmail.com

##### © 2024 th3y3m
