# District Funding

A crowdfunding service for people who want to make their district better!

*This project is currently work in progress*

# Project structure
By now this project contains 3 microservices
- **Auth**: service for storing sensitive data and authenticating users
- **Campaign**: CRUD service for crowdfunding campaigns
- **Payment**: service for handling donations via [Yookassa]()

All services are using REST api to communicate.

# Running
1. Set environment variables shown below, or create .env file in root directory.
    ```dotenv
    JWT_SECRET=test

    AUTH_POSTGRES_PASSWORD=test
    AUTH_POSTGRES_HOST=auth-db

    CAMPAIGN_POSTGRES_PASSWORD=test
    CAMPAIGN_POSTGRES_HOST=campaign-db

    PAYMENT_POSTGRES_PASSWORD=test
    PAYMENT_POSTGRES_HOST=payment-db
    ```

2. Use docker compose to automatically make all three services and postgres instances.
    ```
    docker compose up
    ```