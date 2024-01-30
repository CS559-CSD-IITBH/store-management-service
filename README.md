# Store Management Service

It is a Go-based service that provides functionalities related to store CRUD and item CRUD for a merchant.

## Features

- **Integration with MongoDB:** Utilizes [MongoDB](https://github.com/mongodb/mongo) to store store details, and the list of all items with their respective details. 
- **Gin Web Framework:** Uses the [Gin](https://github.com/gin-gonic/gin) web framework for handling HTTP requests and responses.

## Prerequisites

Before running the service, make sure you have the following dependencies installed:

- Go (version 1.20 or higher)
- Docker 

## Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/store-management-service.git
   cd store-management-service

3. Add a `.env` file in the root directory with following fields:
  
   ```
   MONGO_URL=<add your mongo url hosted on cloud>
   MONGO_DB_NAME=<add name of db in the mongo instance>
   MONGO_COLLECTION_NAME=<add name of collection in the mongo instance>
   PORT=<add host port>
   SESSION_KEY=<private encryption key for sessions>
   FRONTEND_URL=<url of frontend service / reverse proxy>
   ```

4. Build the docker image for the **api-service**. Run the following command in the root directory.

   ```
   docker build -t store-management-service-api:latest .
   ```

5. Start the services using **docker-compose**. Note that it is expected your *Mongo* is hosted on cloud.

   ```
   sudo docker-compose compose build && sudo docker-compose up
   ```