# ImgBB Go Service 🚀

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![Docker](https://img.shields.io/badge/Docker-Supported-2496ED?style=flat&logo=docker)
![License](https://img.shields.io/badge/License-MIT-green?style=flat)

A lightweight, blazing-fast Go microservice designed to act as a bridge for uploading images directly to the [ImgBB API](https://api.imgbb.com/).

## ✨ Features

- **Blazing Fast**: Written in pure Go, taking advantage of Go's incredible concurrency and minimal overhead.
- **Docker Ready**: Includes a multi-stage Dockerfile that builds an ultra-lightweight image (based on `scratch`).
- **Flexible API Keys**: Configure a global API key for anonymous client uploads, or allow clients to provide their own API key per request.
- **Plug & Play**: Zero dependencies required on the host machine. Just run it.

## 🚀 Getting Started

### Prerequisites
- [Docker](https://www.docker.com/) (Recommended)
- [Go 1.21+](https://go.dev/) (If running locally without Docker)
- An ImgBB API Key (Get it for free at [api.imgbb.com](https://api.imgbb.com/))

### Running with Docker (Recommended)

1. Clone the repository:
   ```bash
   git clone https://github.com/SantiiRepair/imgbb.git
   cd imgbb
   ```

2. Configure your API key (Optional):
   Create a `.env` file in the root directory and add your ImgBB API key:
   ```env
   IMGBB_API_KEY=your_api_key_here
   ```

3. Start the service:
   ```bash
   docker-compose up -d --build
   ```

The service will be available at `http://localhost:8080`.

### Running Locally

```bash
# Export your API key
export IMGBB_API_KEY="your_api_key_here"

# Run the server
go run main.go
```

## 📖 API Reference

### Upload Image
Uploads an image to ImgBB and returns the public URL.

**Endpoint:** `POST /upload`
**Content-Type:** `multipart/form-data`

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `image` | `file` | **Yes** | The image file to be uploaded. |
| `api_key` | `string` | No* | Your ImgBB API key. *(Required only if `IMGBB_API_KEY` is not set on the server)* |

#### Example Request (cURL)
```bash
curl -X POST -F "image=@my_photo.jpg" http://localhost:8080/upload
```

#### Example Response
```json
{
  "url": "https://i.ibb.co/example/my_photo.jpg"
}
```

## 📜 License

This project is open-source and available under the MIT License.
