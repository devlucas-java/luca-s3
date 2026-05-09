# Luca S3

Luca S3 es un proyecto en Go para la gestión de videos con autenticación basada en roles, utilizando MinIO para el almacenamiento de objetos, Cassandra como base de datos y Redis para caché.

## Funcionalidades

- **Autenticación por Rol**: Sistema de autenticación con diferentes roles de usuario.
- **Subida de Videos**: Los usuarios pueden subir videos a través de MinIO.
- **Gateway con Nginx**: Implementación de gateway usando Nginx.
- **Manifest HLS**: Recuperación del manifest.m3u8 generado para los videos.
- **Arquitectura Limpia**: Estructura organizada con capas de dominio, aplicación e infraestructura.

## Estructura de Carpetas

```
luca-s3/
├── cmd/
│   └── server/
│       └── main.go                 # Punto de entrada de la aplicación
├── configs/
│   └── config.go                   # Configuración de la aplicación
├── internal/
│   ├── application/
│   │   └── service/
│   │       └── auth_service.go     # Lógica de negocio para autenticación
│   ├── delivery/
│   │   ├── hls/                    # Handlers para HLS
│   │   └── http/
│   │       ├── dto/                # Data Transfer Objects
│   │       │   ├── dauth/          # DTOs de autenticación
│   │       │   └── duser/          # DTOs de usuario
│   │       ├── handler/            # HTTP Handlers
│   │       │   └── auth_handler.go
│   │       ├── middleware/         # Middlewares HTTP
│   │       │   ├── auth_middleware.go
│   │       │   └── role_middleware.go
│   │       └── response/           # Utilidades de respuesta
│   ├── domain/
│   │   ├── entity/                 # Entidades de dominio
│   │   │   ├── metadata.go
│   │   │   └── user.go
│   │   ├── enums/                  # Enums
│   │   │   ├── role.go
│   │   │   └── video_type.go
│   │   └── errors/                 # Definiciones de error
│   ├── infrastructure/
│   │   ├── cassandra/              # Conexión con Cassandra
│   │   │   ├── session.go
│   │   │   └── migrations/         # Scripts CQL
│   │   ├── database/               # Repositorios de base de datos
│   │   └── repository/             # Interfaces de repositorio
├── pkg/                            # Paquetes utilitarios
│   ├── id/
│   ├── logger/
│   ├── password_encoder/
│   └── utils/
├── docs/                           # Documentación
│   └── luca-s3-diagram.drawio      # Diagrama del proyecto
├── docker-compose.yaml             # Configuración Docker
├── .env                            # Variables de entorno
└── README.md                       # Este archivo
```

## Diagrama de Arquitectura

El diagrama del proyecto se encuentra en:

```
docs/luca-s3-diagram.drawio
```

Puedes abrirlo con [draw.io](https://app.diagrams.net/) (online) o con la extensión **Draw.io Integration** en VS Code. Muestra el flujo completo entre el cliente, el gateway Nginx, los handlers HTTP, los servicios, MinIO, Cassandra y Redis.

## Cómo Funciona

1. **Autenticación**: Los usuarios inician sesión y reciben tokens JWT basados en sus roles.
2. **Subida de Video**:
   - El usuario envía una solicitud al handler.
   - El handler valida la solicitud y llama al service.
   - El service interactúa con MinIO para la subida y retorna URLs firmadas.
3. **Recuperación de Video**: Endpoint para obtener el manifest.m3u8 del video procesado.
4. **Base de Datos**: Cassandra almacena usuarios y metadatos de los videos.

## Cómo Levantar el Proyecto

### Requisitos Previos

- Docker y Docker Compose
- Go 1.26+

### Pasos

1. **Clonar el repositorio**:
   ```bash
   git clone https://github.com/devlucas-java/luca-s3.git
   cd luca-s3
   ```

2. **Levantar los servicios con Docker Compose**:
   ```bash
   docker-compose up -d
   ```
   Esto iniciará Cassandra, MinIO, Redis y el worker FFmpeg.

3. **Configurar el .env**:
   - Edita el archivo `.default.env` en la raíz con tus configuraciones (ya viene con valores de prueba).

4. **Ejecutar la aplicación**:
   ```bash
   go run cmd/server/main.go
   ```

5. **Acceder a la aplicación**:
   - API: http://localhost:8080
   - MinIO Console: http://localhost:9001 (username/password)

## Uso

### Endpoints Principales

- `POST /auth/login` — Inicio de sesión
- `POST /auth/register` — Registro de usuario
- `POST /videos/upload` — Subida de video (requiere autenticación)
- `GET /videos/{id}/manifest` — Obtener manifest HLS (requiere autenticación)

### Ejemplo de Subida

1. Inicia sesión para obtener el token JWT.
2. Envía el video via POST a `/videos/upload` con el token en el header `Authorization: Bearer <token>`.
3. Recibe las URLs firmadas de MinIO para la subida directa.

## Desarrollo

- Usa `go mod tidy` para gestionar dependencias.
- Ejecuta los tests con `go test ./...`.
- Para debugging, usa las configuraciones de VS Code o las herramientas integradas.

## Contribución

1. Haz un fork del proyecto
2. Crea una branch para tu feature
3. Haz commit de tus cambios
4. Push a la branch
5. Abre un Pull Request

---

# Luca S3 — English

Luca S3 is a Go project for video management with role-based authentication, using MinIO for object storage, Cassandra as the database, and Redis for caching.

## Features

- **Role-Based Authentication**: Authentication system with different user roles.
- **Video Upload**: Users can upload videos through MinIO.
- **Nginx Gateway**: Gateway implementation using Nginx.
- **HLS Manifest**: Retrieval of the generated manifest.m3u8 for videos.
- **Clean Architecture**: Organized structure with domain, application, and infrastructure layers.

## Folder Structure

```
luca-s3/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── configs/
│   └── config.go                   # Application configuration
├── internal/
│   ├── application/
│   │   └── service/
│   │       └── auth_service.go     # Business logic for authentication
│   ├── delivery/
│   │   ├── hls/                    # HLS handlers
│   │   └── http/
│   │       ├── dto/                # Data Transfer Objects
│   │       │   ├── dauth/          # Authentication DTOs
│   │       │   └── duser/          # User DTOs
│   │       ├── handler/            # HTTP Handlers
│   │       │   └── auth_handler.go
│   │       ├── middleware/         # HTTP Middlewares
│   │       │   ├── auth_middleware.go
│   │       │   └── role_middleware.go
│   │       └── response/           # Response utilities
│   ├── domain/
│   │   ├── entity/                 # Domain entities
│   │   │   ├── metadata.go
│   │   │   └── user.go
│   │   ├── enums/                  # Enums
│   │   │   ├── role.go
│   │   │   └── video_type.go
│   │   └── errors/                 # Error definitions
│   ├── infrastructure/
│   │   ├── cassandra/              # Cassandra connection
│   │   │   ├── session.go
│   │   │   └── migrations/         # CQL scripts
│   │   ├── database/               # Database repositories
│   │   └── repository/             # Repository interfaces
├── pkg/                            # Utility packages
│   ├── id/
│   ├── logger/
│   ├── password_encoder/
│   └── utils/
├── docs/                           # Documentation
│   └── luca-s3-diagram.drawio      # Project diagram
├── docker-compose.yaml             # Docker configuration
├── .env                            # Environment variables
└── README.md                       # This file
```

## Architecture Diagram

The project diagram is located at:

```
docs/luca-s3-diagram.drawio
```

Open it with [draw.io](https://app.diagrams.net/) (online) or the **Draw.io Integration** extension in VS Code. It shows the full flow between the client, Nginx gateway, HTTP handlers, services, MinIO, Cassandra, and Redis.

## How It Works

1. **Authentication**: Users log in and receive JWT tokens based on their roles.
2. **Video Upload**:
   - User sends a request to the handler.
   - Handler validates the request and calls the service.
   - Service interacts with MinIO for the upload and returns signed URLs.
3. **Video Retrieval**: Endpoint to get the manifest.m3u8 of the processed video.
4. **Database**: Cassandra stores users and video metadata.

## How to Run the Project

### Prerequisites

- Docker and Docker Compose
- Go 1.26+

### Steps

1. **Clone the repository**:
   ```bash
   git clone https://github.com/devlucas-java/luca-s3.git
   cd luca-s3
   ```

2. **Start services with Docker Compose**:
   ```bash
   docker-compose up -d
   ```
   This will start Cassandra, MinIO, Redis, and the FFmpeg worker.

3. **Configure the .env**:
   - Edit the `.default.env` file at the root with your settings (already filled with test values).

4. **Run the application**:
   ```bash
   go run cmd/server/main.go
   ```

5. **Access the application**:
   - API: http://localhost:8080
   - MinIO Console: http://localhost:9001 (username/password)

## Usage

### Main Endpoints

- `POST /auth/login` — User login
- `POST /auth/register` — User registration
- `POST /videos/upload` — Video upload (requires authentication)
- `GET /videos/{id}/manifest` — Get HLS manifest (requires authentication)

### Upload Example

1. Log in to obtain the JWT token.
2. Send the video via POST to `/videos/upload` with the token in the `Authorization: Bearer <token>` header.
3. Receive signed MinIO URLs for direct upload.

## Development

- Use `go mod tidy` to manage dependencies.
- Run tests with `go test ./...`.
- For debugging, use VS Code configurations or integrated tools.

## Contributing

1. Fork the project
2. Create a branch for your feature
3. Commit your changes
4. Push to the branch
5. Open a Pull Request
