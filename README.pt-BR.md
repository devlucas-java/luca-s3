# Luca-S3 Transcode Worker

Um worker de transcodificação de vídeo de alta performance que converte vídeos para o formato HLS (HTTP Live Streaming) com streaming adaptativo de bitrate. Construído com Go e projetado para funcionar perfeitamente com armazenamento de objetos MinIO.

## 🎯 Visão Geral

Luca-S3 é um microsserviço especializado focado exclusivamente em transcodificação de vídeo. Não é um sistema completo de gerenciamento de vídeos - é um **worker** que:

- ✅ Transcodifica vídeos para formato HLS com múltiplas resoluções
- ✅ Gera segmentos de 2 segundos para streaming suave
- ✅ Cria thumbnails para cada resolução
- ✅ Suporta vídeos 360° com metadados espaciais
- ✅ Retoma jobs de transcodificação interrompidos
- ✅ Valida restrições de resolução (não faz upscale)

**O que ele NÃO faz:**
- ❌ Upload de vídeo (use MinIO diretamente)
- ❌ Gerenciamento de metadados de vídeo (use seu próprio banco de dados)
- ❌ Autenticação de usuário (implemente no seu API gateway)
- ❌ Reprodução de vídeo (use player HLS no seu frontend)

## 🏗️ Arquitetura

```
┌─────────────┐      ┌──────────────┐      ┌─────────────┐
│   Cliente   │─────▶│  API gRPC    │─────▶│   Worker    │
│  (Sua App)  │      │ (Porta 50051)│      │  (FFmpeg)   │
└─────────────┘      └──────────────┘      └─────────────┘
                            │                      │
                            ▼                      ▼
                     ┌──────────────┐      ┌─────────────┐
                     │    Redis     │      │    MinIO    │
                     │(Estado Jobs) │      │(Armazenamento)│
                     └──────────────┘      └─────────────┘
```

### Componentes

- **Servidor gRPC**: Expõe API de transcodificação
- **FFmpeg**: Processa os vídeos
- **MinIO**: Armazenamento de objetos (volume compartilhado com worker)
- **Redis**: Gerenciamento de estado dos jobs
- **Armazenamento Temporário**: `/tmp/ffmpeg` para processamento (efêmero)

## 📁 Estrutura do Projeto

```
luca-s3/
├── cmd/
│   └── server/
│       └── main.go                 # Ponto de entrada da aplicação
├── internal/
│   ├── application/
│   │   └── service/
│   │       ├── transcode_service.go          # Lógica principal de transcodificação
│   │       ├── hls_transcoder.go             # Geração de HLS
│   │       ├── video_analyzer.go             # Extração de metadados do vídeo
│   │       └── segment_thumbnail_generator.go # Geração de thumbnails
│   ├── delivery/
│   │   └── grpc/
│   │       ├── handler.go          # Handlers gRPC
│   │       ├── server.go           # Configuração do servidor gRPC
│   │       └── pb/                 # Arquivos protobuf gerados
│   ├── domain/
│   │   ├── enums/
│   │   │   ├── extension.go        # Formatos suportados (MP4, MOV, WEBM)
│   │   │   ├── job_status.go       # Estados do job
│   │   │   ├── resolution.go       # Resoluções suportadas
│   │   │   └── video_type.go       # Vídeo normal/360°
│   │   └── model/
│   │       └── job.go              # Modelo TranscodeJob
│   └── infrastructure/
│       ├── minio/
│       │   ├── client.go           # Cliente MinIO
│       │   └── client_test.go      # Testes MinIO
│       └── redis/
│           ├── client.go           # Cliente Redis
│           ├── client_test.go      # Testes Redis
│           ├── job_repository.go   # Persistência de jobs
│           └── job_repository_test.go
├── pkg/
│   └── logger/
│       └── logger.go               # Utilitário de logging
├── proto/
│   └── video.proto                 # Definição do serviço gRPC
├── configs/
│   └── config.go                   # Gerenciamento de configuração
├── docker-compose.yaml             # Orquestração Docker
├── Dockerfile                      # Container do worker
├── .default.env                    # Variáveis de ambiente
└── README.md
```

## 🚀 Início Rápido

### Pré-requisitos

- Docker & Docker Compose
- Go 1.23+ (para desenvolvimento local)
- FFmpeg (para desenvolvimento local)

### Executando com Docker

```bash
# Iniciar todos os serviços
docker-compose up -d

# Verificar logs
docker-compose logs -f worker

# Parar serviços
docker-compose down
```

Serviços estarão disponíveis em:
- **Worker gRPC**: `localhost:50051`
- **Console MinIO**: `http://localhost:9001` (username/password)
- **API MinIO**: `localhost:9000`
- **Redis**: `localhost:6379`

### Desenvolvimento Local

```bash
# Instalar dependências
go mod download

# Executar testes
go test ./internal/infrastructure/...

# Gerar arquivos proto
protoc --go_out=internal/delivery/grpc/pb \
       --go-grpc_out=internal/delivery/grpc/pb \
       --go_opt=paths=source_relative \
       --go-grpc_opt=paths=source_relative \
       proto/video.proto

# Executar worker
go run cmd/server/main.go
```

## 🔧 Configuração

Variáveis de ambiente (`.default.env`):

```env
# Configuração Redis
REDIS_ADDR=redis:6379
REDIS_PASSWORD=
REDIS_DB=0

# Configuração MinIO
MINIO_ENDPOINT=minio:9000
MINIO_ACCESS_KEY=username
MINIO_SECRET_KEY=password
MINIO_USE_SSL=false

# Configuração do Servidor
SERVER_PORT=50051

# Configuração FFmpeg
FFMPEG_WORK_DIR=/tmp/ffmpeg
```

## 📡 Referência da API

### TranscodeVideo

Inicia transcodificação HLS para um vídeo no MinIO.

```protobuf
rpc TranscodeVideo(TranscodeVideoRequest) returns (TranscodeVideoResponse);

message TranscodeVideoRequest {
  string video_id = 1;                    // ID do vídeo no MinIO
  string original_path = 2;               // Caminho: videos/{id}.mp4
  repeated Resolution resolutions = 3;    // Resoluções desejadas
}
```

**Exemplo:**
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

Consulta o status do job de transcodificação.

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

Obtém URL pré-assinada para o master playlist HLS.

```protobuf
rpc GetHLSManifest(GetHLSManifestRequest) returns (GetHLSManifestResponse);

message GetHLSManifestRequest {
  string video_id = 1;
  int32 expires_in = 2;  // segundos, padrão 3600
}
```

### DeleteVideo

Deleta vídeo e todos os arquivos relacionados do MinIO.

```protobuf
rpc DeleteVideo(DeleteVideoRequest) returns (DeleteVideoResponse);
```

## 🎬 Como Funciona

### 1. Upload de Vídeo (Sua Responsabilidade)

Faça upload dos vídeos diretamente para o MinIO:

```bash
# Usando MinIO CLI
mc cp video.mp4 myminio/videos/abc123.mp4

# Ou use o SDK do MinIO na sua aplicação
```

### 2. Requisição de Transcodificação

Chame o worker via gRPC:

```go
job, err := client.TranscodeVideo(ctx, &pb.TranscodeVideoRequest{
    VideoId:      "abc123",
    OriginalPath: "videos/abc123.mp4",
    Resolutions:  []pb.Resolution{pb.Resolution_RESOLUTION_720P},
})
```

### 3. Processamento

O worker:
1. Baixa o vídeo do MinIO para `/tmp/ffmpeg/{job_id}/`
2. Analisa o vídeo (resolução, duração)
3. Filtra resoluções inválidas (sem upscaling)
4. Verifica MinIO por segmentos existentes (capacidade de retomar)
5. Gera segmentos HLS (2s cada) com FFmpeg
6. Faz upload dos segmentos para MinIO: `videos/{id}/hls/{resolution}/`
7. Gera thumbnails
8. Cria master playlist
9. Limpa arquivos temporários

### 4. Reprodução (Sua Responsabilidade)

Use a URL do manifest HLS no seu player de vídeo:

```javascript
const player = new Hls();
player.loadSource('https://minio.example.com/videos/abc123/hls/master.m3u8');
player.attachMedia(videoElement);
```

## 🎯 Formatos Suportados

### Formatos de Entrada
- **MP4** (H.264/H.265)
- **MOV** (QuickTime)
- **WEBM** (VP8/VP9)

### Resoluções de Saída
- **4K** (3840x2160)
- **1440p** (2560x1440)
- **1080p** (1920x1080)
- **720p** (1280x720)
- **480p** (854x480)
- **360p** (640x360)

**Nota:** O worker filtra automaticamente resoluções maiores que o vídeo original.

## 🔄 Capacidade de Retomar

O worker pode retomar jobs interrompidos:

1. Verifica MinIO por segmentos existentes: `videos/{id}/hls/{res}/segment_*.ts`
2. Encontra o número do último segmento
3. Continua a partir de `last_segment + 1`
4. Atualiza progresso do job no Redis

Isso evita reprocessar segmentos já completados.

## 🧪 Testes

```bash
# Executar todos os testes
go test ./...

# Executar testes de infraestrutura (requer Redis & MinIO)
go test ./internal/infrastructure/...

# Executar com cobertura
go test -cover ./...

# Executar teste específico
go test -v ./internal/infrastructure/minio -run TestMinIOClient
```

## 🐳 Volumes Docker

- **luca_s3_data**: Compartilhado entre MinIO e worker (persistente)
- **/tmp/ffmpeg**: Diretório de processamento temporário (efêmero)

O volume compartilhado permite que o worker acesse arquivos do MinIO diretamente sem transferência de rede.

## 📊 Monitoramento

### Estados do Job

- **pending**: Job criado, aguardando início
- **processing**: Transcodificando atualmente
- **done**: Completado com sucesso
- **failed**: Ocorreu um erro

### Chaves Redis

Jobs são armazenados no Redis com TTL de 7 dias:
```
transcode:job:{job_id}
```

## 🔒 Considerações de Segurança

1. **Sem Autenticação**: Implemente autenticação no seu API gateway
2. **Acesso MinIO**: Use políticas IAM para restringir acesso ao bucket
3. **Rede**: Execute o worker em rede privada, exponha apenas porta gRPC
4. **Secrets**: Use Docker secrets ou vault para credenciais

## 🚧 Limitações

- Sem rate limiting integrado (implemente no seu API gateway)
- Sem priorização de jobs (processamento FIFO)
- Sem processamento distribuído (instância única do worker)
- Sem notificações webhook (implemente na sua aplicação)

## 📝 Licença

Licença MIT - Veja arquivo LICENSE para detalhes

## 🤝 Contribuindo

1. Faça fork do repositório
2. Crie uma branch de feature
3. Adicione testes para nova funcionalidade
4. Submeta um pull request

## 📞 Suporte

Para problemas e questões:
- Abra uma issue no GitHub
- Verifique issues existentes para soluções

---

**Construído com ❤️ usando Go, FFmpeg, MinIO e Redis**
