
## Project Documentation

### Overview
This project provides an API for registring zos nodes using the Go Gin framework and PostgreSQL database. 
The API supports operations like registring, listing, and updating farms and nodes, as well as reporting uptime and consumption data for nodes.

### Endpoint Descriptions

#### Farms Endpoints
1. **GET /farms/** - List all farms, or use FarmFilter to list specific set of farms.
2. **GET /farms/:farm_id** - Get a specific farm by ID.
3. **POST /farms/** - Create a new farm.
4. **PATCH /farms/** - Update an existing farm.

#### Nodes Endpoints
1. **GET /nodes/** - List all nodes, or use NodeFilter to list specific set of nodes.
2. **GET /nodes/:node_id** - Get a specific node by ID.
3. **POST /nodes/** - Register a new node.
4. **POST /nodes/:node_id/uptime** - Report uptime for a specific node.
5. **POST /nodes/:node_id/consumption** - Report consumption for a specific node.

## Setup Instructions

1. **Start PostgreSQL:**
   ```bash
   make postgres
   ```
2. **Run the Server:**
   ```bash
   make run
   ```
3. **Stop PostgreSQL:**
   ```bash
   make stop-postgres
   ```

### Swagger Documentation
Once the server is running, Swagger documentation can be accessed at:
```
http://<domain>:<port>/swagger/index.html
```
Replace `<domain>` and `<port>` with the appropriate values.

### How to Use the Server
1. Use a tool like Postman or cURL to interact with the API.
2. Refer to the Swagger documentation for detailed information about request parameters and response structures.
