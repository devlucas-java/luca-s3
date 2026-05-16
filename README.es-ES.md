# Luca-S3 Transcode Worker

Un worker de transcodificación de vídeo de alto rendimiento que convierte vídeos al formato HLS (HTTP Live Streaming) con streaming adaptativo de bitrate. Construido con Go y diseñado para funcionar perfectamente con almacenamiento de objetos MinIO.

## 🎯 Visión General

Luca-S3 es un microservicio especializado enfocado exclusivamente en transcodificación de vídeo. No es un sistema completo de gestión de vídeos - es un **worker** que:

- ✅ Transcodifica vídeos a formato HLS con múltiples resoluciones
- ✅ Genera segmentos de 2 segundos para streaming fluido
- ✅ Crea miniaturas para cada resolución
- ✅ Soporta vídeos 360° con metadatos espaciales
- ✅ Reanuda trabajos de transcodificación interrumpidos
- ✅ Valida restricciones de resolución (no hace upscaling)

**Lo que NO hace:**
- ❌ Subida de vídeo (usa MinIO directamente)
- ❌ Gestión de metadatos de vídeo (usa tu propia base de datos)
- ❌ Autenticación de usuario (implementa en tu API gateway)
- ❌ Reproducción de vídeo (usa reproductor HLS en tu frontend)

## 🏗️ Arquitectura

```
┌─────────────┐      ┌──────────────┐      ┌─────────────┐
│   Cliente   │─────▶│  API gRPC    │─────▶│   Worker    │
│  (Tu App)   │      │ (Puerto 50051)│      │  (FFmpeg)   │
└─────────────┘      └──────────────┘      └─────────────┘
                            │                      │
                            ▼                      ▼
                     ┌──────────────┐      ┌─────────────┐
                     │    Redis     │      │    MinIO    │
                     │(Estado Jobs) │      │(Almacenamiento)│
                     └──────────────┘      └─────────────┘
```

### Componentes

- **Servidor gRPC**: Expone API de transcodificación
- **FFmpeg**: Procesa los vídeos
- **MinIO**: Almacenamiento de objetos (volumen compartido con worker)
- **Redis**: Gestión de estado de los trabajos
- **Almacenamiento Temporal**: `/tmp/ffmpeg` para procesamiento (efímero)

## 📁 Estructura del Proyecto

```
luca-s3/
├── cmd/
│   └── server/
│       └── main.go                 # Punto de entrada de la aplicación
├── internal/
│   ├── application/
│   │   └── service/
│   │       ├── transcode_service.go          # Lógica principal de transcodificación
│   │       ├── hls_transcoder.go             # Generación de HLS
│   │       ├── video_analyzer.go             # Extracción de metadatos del vídeo
│   │       └── segment_thumbnail_generator.go # Generación de miniaturas
│   ├── delivery/
│   │   └── grpc/
│   │       ├── handler.go          # Handlers gRPC
│   │       ├── server.go           # Configuración del servidor gRPC
│   │       └── pb/                 # Archivos protobuf generados
│   ├── domain/
│   │   ├── enums/
│   │   │   ├── extension.go        # Formatos soportados (MP4, MOV, WEBM)
│   │   │   ├── job_status.go       # Estados del trabajo
│   │   │   ├── resolution.go       # Resoluciones soportadas
│   │   │   └── video_type.go       # Vídeo normal/360°
│   │   └── model/
│   │       └── job.go              # Modelo TranscodeJob
│   └── infrastructure/
│       ├── minio/
│       │   ├── client.go           # Cliente MinIO
│       │   └── client_test.go      # Tests MinIO
│       └── redis/
│           ├── client.go           # Cliente Redis
│           ├── client_test.go      # Tests Redis
│           ├── job_repository.go   # Persistencia de trabajos
│           └── job_repository_test.go
├── pkg/
│   └── logger/
│       └── logger.go               # Utilidad de logging
├── proto/
│   └── video.proto                 # Definición del servicio gRPC
├── configs/
│   └── config.go                   # Gestión de configuración
├── docker-compose.yaml             # Orquestación Docker
├── Dockerfile                      # Contenedor del worker
├── .default.env                    # Variables de entorno
└── README.md
```

## 🚀 Inicio Rápido

### Requisitos Previos

- Docker & Docker Compose
- Go 1.23+ (para desarrollo local)
- FFmpeg (para desarrollo local)

### Ejecutando con Docker

```bash
# Iniciar todos los servicios
docker-compose up -d

# Verificar logs
docker-compose logs -f worker

# Detener servicios
docker-compose down
```

Los servicios estarán disponibles en:
- **Worker gRPC**: `localhost:50051`
- **Consola MinIO**: `http://localhost:9001` (username/password)
- **API MinIO**: `localhost:9000`
- **Redis**: `localhost:6379`

### Desarrollo Local

```bash
# Instalar dependencias
go mod download

# Ejecutar tests
go test ./internal/infrastructure/...

# Generar archivos proto
protoc --go_out=internal/delivery/grpc/pb \
       --go-grpc_out=internal/delivery/grpc/pb \
       --go_opt=paths=source_relative \
       --go-grpc_opt=paths=source_relative \
       proto/video.proto

# Ejecutar worker
go run cmd/server/main.go
```

## 🔧 Configuración

Variables de entorno (`.default.env`):

```env
# Configuración Redis
REDIS_ADDR=redis:6379
REDIS_PASSWORD=
REDIS_DB=0

# Configuración MinIO
MINIO_ENDPOINT=minio:9000
MINIO_ACCESS_KEY=username
MINIO_SECRET_KEY=password
MINIO_USE_SSL=false

# Configuración del Servidor
SERVER_PORT=50051

# Configuración FFmpeg
FFMPEG_WORK_DIR=/tmp/ffmpeg
```

## 📡 Referencia de la API

### TranscodeVideo

Inicia transcodificación HLS para un vídeo en MinIO.

```protobuf
rpc TranscodeVideo(TranscodeVideoRequest) returns (TranscodeVideoResponse);

message TranscodeVideoRequest {
  string video_id = 1;                    // ID del vídeo en MinIO
  string original_path = 2;               // Ruta: videos/{id}.mp4
  repeated Resolution resolutions = 3;    // Resoluciones deseadas
}
```

**Ejemplo:**
```go
req := &pb.TranscodeVideoRequest{
    VideoId:      "abc123",
    OriginalPath: "videos/abc123.mp4",
    Resolutions:  []pb.Resolution{
        pb.Resolution_RESOLUTION_1080P,
        pb.Resolution_RESOLUTION_720P,
        pb.Resolution_RESOLUTION_480P,
    },
}
```

### GetJobStatus

Consulta el estado del trabajo de transcodificación.

```protobuf
rpc GetJobStatus(JobStatusRequest) returns (JobStatusResponse);

message JobStatusResponse {
  string job_id = 1;
  string video_id = 2;
  JobStatus status = 3;
  repeated ResolutionProgress progress = 4;
  int32 original_width = 5;
  int32 original_height = 6;
  double duration_seconds = 7;
}
```

### GetHLSManifest

Obtiene URL prefirmada para el master playlist HLS.

```protobuf
rpc GetHLSManifest(GetHLSManifestRequest) returns (GetHLSManifestResponse);

message GetHLSManifestRequest {
  string video_id = 1;
  int32 expires_in = 2;  // segundos, por defecto 3600
}
```

### DeleteVideo

Elimina vídeo y todos los archivos relacionados de MinIO.

```protobuf
rpc DeleteVideo(DeleteVideoRequest) returns (DeleteVideoResponse);
```

## 🎬 Cómo Funciona

### 1. Subida de Vídeo (Tu Responsabilidad)

Sube los vídeos directamente a MinIO:

```bash
# Usando MinIO CLI
mc cp video.mp4 myminio/videos/abc123.mp4

# O usa el SDK de MinIO en tu aplicación
```

### 2. Solicitud de Transcodificación

Llama al worker vía gRPC:

```go
job, err := client.TranscodeVideo(ctx, &pb.TranscodeVideoRequest{
    VideoId:      "abc123",
    OriginalPath: "videos/abc123.mp4",
    Resolutions:  []pb.Resolution{pb.Resolution_RESOLUTION_720P},
})
```

### 3. Procesamiento

El worker:
1. Descarga el vídeo de MinIO a `/tmp/ffmpeg/{job_id}/`
2. Analiza el vídeo (resolución, duración)
3. Filtra resoluciones inválidas (sin upscaling)
4. Verifica MinIO por segmentos existentes (capacidad de reanudar)
5. Genera segmentos HLS (2s cada uno) con FFmpeg
6. Sube los segmentos a MinIO: `videos/{id}/hls/{resolution}/`
7. Genera miniaturas
8. Crea master playlist
9. Limpia archivos temporales

### 4. Reproducción (Tu Responsabilidad)

Usa la URL del manifest HLS en tu reproductor de vídeo:

```javascript
const player = new Hls();
player.loadSource('https://minio.example.com/videos/abc123/hls/master.m3u8');
player.attachMedia(videoElement);
```

## 🎯 Formatos Soportados

### Formatos de Entrada
- **MP4** (H.264/H.265)
- **MOV** (QuickTime)
- **WEBM** (VP8/VP9)

### Resoluciones de Salida
- **4K** (3840x2160)
- **1440p** (2560x1440)
- **1080p** (1920x1080)
- **720p** (1280x720)
- **480p** (854x480)
- **360p** (640x360)

**Nota:** El worker filtra automáticamente resoluciones mayores que el vídeo original.

## 🔄 Capacidad de Reanudar

El worker puede reanudar trabajos interrumpidos:

1. Verifica MinIO por segmentos existentes: `videos/{id}/hls/{res}/segment_*.ts`
2. Encuentra el número del último segmento
3. Continúa desde `last_segment + 1`
4. Actualiza progreso del trabajo en Redis

Esto evita reprocesar segmentos ya completados.

## 🧪 Tests

```bash
# Ejecutar todos los tests
go test ./...

# Ejecutar tests de infraestructura (requiere Redis & MinIO)
go test ./internal/infrastructure/...

# Ejecutar con cobertura
go test -cover ./...

# Ejecutar test específico
go test -v ./internal/infrastructure/minio -run TestMinIOClient
```

## 🐳 Volúmenes Docker

- **luca_s3_data**: Compartido entre MinIO y worker (persistente)
- **/tmp/ffmpeg**: Directorio de procesamiento temporal (efímero)

El volumen compartido permite que el worker acceda a archivos de MinIO directamente sin transferencia de red.

## 📊 Monitorización

### Estados del Trabajo

- **pending**: Trabajo creado, esperando inicio
- **processing**: Transcodificando actualmente
- **done**: Completado con éxito
- **failed**: Ocurrió un error

### Claves Redis

Los trabajos se almacenan en Redis con TTL de 7 días:
```
transcode:job:{job_id}
```

## 🔒 Consideraciones de Seguridad

1. **Sin Autenticación**: Implementa autenticación en tu API gateway
2. **Acceso MinIO**: Usa políticas IAM para restringir acceso al bucket
3. **Red**: Ejecuta el worker en red privada, expón solo puerto gRPC
4. **Secrets**: Usa Docker secrets o vault para credenciales

## 🚧 Limitaciones

- Sin rate limiting integrado (implementa en tu API gateway)
- Sin priorización de trabajos (procesamiento FIFO)
- Sin procesamiento distribuido (instancia única del worker)
- Sin notificaciones webhook (implementa en tu aplicación)

## 📝 Licencia

Licencia MIT - Ver archivo LICENSE para detalles

## 🤝 Contribuyendo

1. Haz fork del repositorio
2. Crea una rama de feature
3. Añade tests para nueva funcionalidad
4. Envía un pull request

## 📞 Soporte

Para problemas y cuestiones:
- Abre una issue en GitHub
- Verifica issues existentes para soluciones

---

**Construido con ❤️ usando Go, FFmpeg, MinIO y Redis**
