# TASK.md — Luca S3

> Plataforma de gestión y streaming de videos con autenticación por roles, almacenamiento en MinIO, base de datos Cassandra, caché Redis, gateway Nginx y procesamiento HLS con FFmpeg.

---

## Leyenda / Legend

| Símbolo | Significado |
|---------|-------------|
| ✅ | Implementado y funcional / Implemented and working |
| 🔧 | Parcialmente implementado / Partially implemented |
| ❌ | Pendiente / Not yet implemented |

---

## 1. Infraestructura / Infrastructure

| # | Tarea / Task | Estado |
|---|---|---|
| 1.1 | Cassandra — conexión con autenticación y keyspace | ✅ |
| 1.2 | Cassandra — bootstrap session (sin keyspace) para migraciones | ✅ |
| 1.3 | Cassandra — runner de migraciones `.cql` ordenadas | ✅ |
| 1.4 | Cassandra — migración `001`: keyspace `luca_s3` | ✅ |
| 1.5 | Cassandra — migración `002`: tabla `users` con índices | ✅ |
| 1.6 | Cassandra — migración `003`: tabla `metadata` con índices | ✅ |
| 1.7 | Docker Compose — servicios: Cassandra, MinIO, Redis, App, Nginx | ✅ |
| 1.8 | Docker Compose — healthchecks en todos los servicios | ✅ |
| 1.9 | Docker Compose — `depends_on` con `condition: service_healthy` | ✅ |
| 1.10 | Nginx — configuración básica como gateway | 🔧 archivo creado, rutas pendientes |
| 1.11 | MinIO — cliente Go inicializado y conectado | ❌ |
| 1.12 | MinIO — creación automática del bucket al iniciar | ❌ |
| 1.13 | Redis — cliente Go inicializado y conectado | ❌ |
| 1.14 | Config — carga desde `.default.env` con validación | ✅ |

---

## 2. Dominio / Domain

| # | Tarea / Task | Estado |
|---|---|---|
| 2.1 | Entidad `User` con `gocql.UUID`, roles y timestamps | ✅ |
| 2.2 | Entidad `MetaData` con `gocql.UUID`, status y timestamps | ✅ |
| 2.3 | Enum `Role` — ADMIN, USER | ✅ |
| 2.4 | Enum `Status` — FAILED, PENDING, UPLOADED, NOUPLOADED | ✅ |
| 2.5 | Enum `VideoType` — normal, 360 | ✅ |
| 2.6 | Errores de dominio tipados (`AppError`) con código y status HTTP | ✅ |
| 2.7 | Todos los constructores de error: NotFound, Conflict, Unauthorized, etc. | ✅ |

---

## 3. Repositorios / Repositories

| # | Tarea / Task | Estado |
|---|---|---|
| 3.1 | Interface `UserRepository` con `gocql.UUID` | ✅ |
| 3.2 | Interface `MetaDataRepository` | ✅ |
| 3.3 | `UserDB` — Create, Save, Updates, DeleteByID | ✅ |
| 3.4 | `UserDB` — FindByID, FindByEmail, FindByUsername, FindByEmailOrUsername | ✅ |
| 3.5 | `UserDB` — paginación con cursor Cassandra (`PageState`) | ✅ |
| 3.6 | `MetaDataDB` — Create, Updates, DeleteByID | ✅ |
| 3.7 | `MetaDataDB` — FindByID, FindAllByUserID, FindByObjectKey, FindAllByStatus | ✅ |
| 3.8 | `MetaDataDB` — paginación con cursor Cassandra en todos los listados | ✅ |
| 3.9 | `pkg/pagination` — tipos `Page` y `PagedResult[T]` genéricos | ✅ |

---

## 4. Seguridad / Security

| # | Tarea / Task | Estado |
|---|---|---|
| 4.1 | `JWTService` — generación de token HS256 con claims (user_id, email, roles) | ✅ |
| 4.2 | `JWTService` — validación y extracción de claims | ✅ |
| 4.3 | `AuthMiddleware` — valida JWT, extrae usuario del DB, inyecta en contexto | ✅ |
| 4.4 | `RoleMiddleware` — verifica roles requeridos, retorna 403 si no tiene acceso | ✅ |
| 4.5 | `IdempotencyMiddleware` — prevención de requests duplicados con Redis | 🔧 stub vacío |
| 4.6 | `password_encoder` — bcrypt encode y match | ✅ |
| 4.7 | Rate limiting middleware | ❌ |
| 4.8 | CORS middleware | ❌ |

---

## 5. Autenticación / Authentication

| # | Tarea / Task | Estado |
|---|---|---|
| 5.1 | `POST /auth/login` — login con email o username, retorna JWT | ✅ |
| 5.2 | `POST /auth/register` — registro, hash de password, retorna JWT | ✅ |
| 5.3 | `PUT /auth/password` — cambio de password (requiere auth) | ✅ |
| 5.4 | DTOs: `LoginRequest`, `RegisterDTO`, `UpdatePasswordRequest` con validación | ✅ |
| 5.5 | DTO: `JWTResponse` con token + datos del usuario | ✅ |
| 5.6 | `UserMapper` — entity ↔ DTO | ✅ |
| 5.7 | `POST /auth/refresh` — renovación de token JWT | ❌ |
| 5.8 | `POST /auth/logout` — invalidación de token (blacklist en Redis) | ❌ |

---

## 6. Usuarios / Users

| # | Tarea / Task | Estado |
|---|---|---|
| 6.1 | `GET /users/me` — perfil del usuario autenticado | ❌ |
| 6.2 | `PUT /users/me` — actualización de perfil (nombre, username) | ❌ |
| 6.3 | `DELETE /users/me` — eliminación de cuenta | ❌ |
| 6.4 | `GET /users` — listado paginado (solo ADMIN) | ❌ |
| 6.5 | `UserService` — lógica de negocio para gestión de usuarios | ❌ |
| 6.6 | DTOs de usuario para update y respuesta | 🔧 archivos existen, lógica pendiente |

---

## 7. Videos — Upload y Almacenamiento / Upload & Storage

| # | Tarea / Task | Estado |
|---|---|---|
| 7.1 | Cliente MinIO inicializado en el arranque | ❌ |
| 7.2 | `VideoService` — lógica de negocio para subida de videos | ❌ |
| 7.3 | `POST /videos/upload` — subida de video a MinIO | ❌ |
| 7.4 | Generación de URL pre-firmada para subida directa desde el cliente | ❌ |
| 7.5 | Guardar `MetaData` en Cassandra tras subida exitosa | ❌ |
| 7.6 | `GET /videos` — listado paginado de videos del usuario | ❌ |
| 7.7 | `GET /videos/{id}` — detalle de un video | ❌ |
| 7.8 | `DELETE /videos/{id}` — eliminar video de MinIO y metadata de Cassandra | ❌ |
| 7.9 | `MetaDataService` — lógica de negocio para gestión de metadata | ❌ |
| 7.10 | DTOs para video upload, respuesta y metadata | ❌ |

---

## 8. Streaming HLS / HLS Streaming

| # | Tarea / Task | Estado |
|---|---|---|
| 8.1 | Worker FFmpeg — transcodificación de video a HLS (`.m3u8` + `.ts`) | ❌ |
| 8.2 | Subida de segmentos HLS a MinIO tras transcodificación | ❌ |
| 8.3 | Actualización de `MetaData.ManifestKey` y `Status` tras procesamiento | ❌ |
| 8.4 | `GET /videos/{id}/manifest` — retorna manifest `.m3u8` | ❌ |
| 8.5 | `GET /videos/{id}/segments/{file}` — sirve segmentos `.ts` | ❌ |
| 8.6 | Soporte para video 360° (enum `VideoType.360`) | ❌ |
| 8.7 | Integración del worker FFmpeg en Docker Compose | ❌ |

---

## 9. Caché / Cache (Redis)

| # | Tarea / Task | Estado |
|---|---|---|
| 9.1 | Cliente Redis inicializado en el arranque | ❌ |
| 9.2 | Caché de sesiones / tokens JWT (blacklist para logout) | ❌ |
| 9.3 | `IdempotencyMiddleware` — caché de requests con Redis | 🔧 stub vacío |
| 9.4 | Caché de metadata de videos frecuentemente accedidos | ❌ |
| 9.5 | TTL configurable por tipo de caché | ❌ |

---

## 10. Gateway Nginx / Nginx Gateway

| # | Tarea / Task | Estado |
|---|---|---|
| 10.1 | Archivo `nginx.conf` creado | 🔧 básico |
| 10.2 | Proxy reverso hacia el servicio `app` en `/api` | ❌ |
| 10.3 | Proxy hacia MinIO para acceso a objetos en `/storage` | ❌ |
| 10.4 | Servir segmentos HLS directamente desde MinIO vía Nginx | ❌ |
| 10.5 | Headers de seguridad (CORS, X-Frame-Options, etc.) | ❌ |
| 10.6 | Nginx integrado y funcional en Docker Compose | 🔧 servicio definido, config incompleta |

---

## 11. gRPC (Opcional / Optional)

| # | Tarea / Task | Estado |
|---|---|---|
| 11.1 | Definición de `.proto` para servicios de video y metadata | ❌ |
| 11.2 | Generación de código Go desde `.proto` | ❌ |
| 11.3 | Servidor gRPC para comunicación interna (worker FFmpeg ↔ API) | ❌ |
| 11.4 | Integración gRPC en Docker Compose | ❌ |

---

## 12. Dashboard (Frontend Simple / Simple Frontend)

| # | Tarea / Task | Estado |
|---|---|---|
| 12.1 | Página de login | ❌ |
| 12.2 | Página de registro | ❌ |
| 12.3 | Listado de videos del usuario con paginación | ❌ |
| 12.4 | Reproductor de video HLS (usando `hls.js` o similar) | ❌ |
| 12.5 | Formulario de subida de video | ❌ |
| 12.6 | Indicador de estado de procesamiento (PENDING → UPLOADED) | ❌ |
| 12.7 | Panel de administración (solo ADMIN) — listado de usuarios | ❌ |

---

## 13. Calidad / Quality

| # | Tarea / Task | Estado |
|---|---|---|
| 13.1 | Tests unitarios — `password_encoder` | ✅ |
| 13.2 | Tests unitarios — `AuthService` (Login, Register, UpdatePassword) | ❌ |
| 13.3 | Tests unitarios — `UserDB` con mock de sesión Cassandra | ❌ |
| 13.4 | Tests unitarios — `MetaDataDB` | ❌ |
| 13.5 | Tests de integración — endpoints de autenticación | ❌ |
| 13.6 | Tests de integración — upload y streaming de video | ❌ |
| 13.7 | Endpoint `GET /health` — estado de todos los servicios | ❌ |
| 13.8 | Structured logging con niveles (DEBUG, INFO, ERROR) | ✅ |

---

## Progreso General / Overall Progress

```
Infraestructura   ████████░░  75%
Dominio           ██████████  100%
Repositorios      ██████████  100%
Seguridad         ████████░░  75%
Autenticación     ████████░░  80%
Usuarios          ██░░░░░░░░  15%
Videos/Storage    ░░░░░░░░░░   0%
HLS/FFmpeg        ░░░░░░░░░░   0%
Redis/Cache       ░░░░░░░░░░   5%
Nginx Gateway     ██░░░░░░░░  15%
gRPC              ░░░░░░░░░░   0%
Dashboard         ░░░░░░░░░░   0%
Calidad/Tests     ██░░░░░░░░  15%
```

---

> Última actualización / Last updated: 2026-05-09
