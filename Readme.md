# Ethereum Rewards API

This project is a RESTful API that interacts with the Ethereum blockchain to provide information about block rewards and sync committee duties for a given slot. It is implemented in Go using the [Gin web framework](https://gin-gonic.com/).

---

## Features

### Endpoints

1. **GET /blockreward/{slot}**
   - Retrieves information about the block reward for a given slot.
   - **Parameters:**
     - `slot` (integer): The slot number in the Ethereum blockchain.
   - **Response:**
     ```json
     {
       "status": "vanilla" | "relay",
       "reward": "<reward_in_gwei>"
     }
     ```

2. **GET /syncduties/{slot}**
   - Retrieves a list of validators with sync committee duties for a given slot.
   - **Parameters:**
     - `slot` (integer): The slot number in the Ethereum blockchain.
   - **Response:**
     ```json
     {
       "validators": ["<validator_index1>", "<validator_index2>", ...]
     }
     ```

---

## Design Choices and Frameworks

### Frameworks Used

- **Gin Web Framework**
  - Selected for its lightweight and fast performance, making it ideal for building high-performance RESTful APIs.

- **godotenv**
  - Chosen to simplify configuration management by loading environment variables from a `.env` file, ensuring flexibility and separation of concerns.

### Services Layer

- **Consensus Service:**
  - Encapsulates logic for interacting with the Ethereum consensus layer, including fetching beacon chain data like blocks, headers, and sync committee duties.

- **Execution Service:**
  - Provides abstraction for communicating with the Ethereum execution layer using JSON-RPC, enabling streamlined data retrieval for execution blocks.

### Docker

- Containerization using Docker ensures that the application can run consistently across various environments without dependency conflicts.

### Error Handling

- Developed custom utility functions for centralized error handling, ensuring meaningful and user-friendly HTTP responses in case of failures.

### Environment Variables

- Configured the QuickNode endpoint using the `QUICKNODE_ENDPOINT` environment variable, promoting secure and dynamic configuration.

---

## How to Build and Run

### Prerequisites

- Docker installed
- (Optional) Go installed for local development

### Running with Docker

1. **Pull the Docker Image:**
   ```bash
   docker pull ashhar001/eth-rewards-api:latest
   ```

2. **Run the Container:**
   ```bash
   docker run -d -p 8080:8080 -e QUICKNODE_ENDPOINT=<your_quicknode_endpoint> ashhar001/eth-rewards-api:latest
   ```

3. **Access the API:**
   The API will be accessible at `http://localhost:8080`.

### Running Locally

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/ashhar001/eth-rewards-api.git
   cd eth-rewards-api
   ```

2. **Set Up Environment Variables:**
   Create a `.env` file in the project root with the following content:
   ```env
   QUICKNODE_ENDPOINT=<your_quicknode_endpoint>
   ```

3. **Install Dependencies:**
   ```bash
   go mod download
   ```

4. **Run the Application:**
   ```bash
   go run cmd/main.go
   ```

5. **Access the API:**
   The API will be accessible at `http://localhost:8080`.

---

## Example API Calls

### Get Block Reward

#### Request
```bash
curl http://localhost:8080/blockreward/10590951
```

#### Response
```json
{
  "status": "relay",
  "reward": "20850"
}
```

#### Request
```bash
curl http://localhost:8080/blockreward/10589928
```

#### Response
```json
{
  "status": "vanilla",
  "reward": "20850"
}
```

### Get Sync Committee Duties

#### Request
```bash
curl http://localhost:8080/syncduties/10590951
```

#### Response
```json
{
  "validators": [
    "ValidatorIndex1",
    "ValidatorIndex2",
    "..."
  ]
}
```

---

