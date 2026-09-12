# golang-back — API de Biblioteca

API REST en Go para la gestión de libros, construida con **arquitectura limpia**
(handler → service → store), configuración por variables de entorno,
migraciones embebidas, logging estructurado, apagado ordenado (*graceful
shutdown*) y una suite de tests. Pensada como una base sólida y realista para
un backend en producción, no solo un ejemplo de CRUD.

## Características

- ✅ Arquitectura en capas: `transport/rest` → `service` → `store`, con
  interfaces definidas del lado del consumidor (bajo acoplamiento y fácil de testear).
- ✅ Router [chi](https://github.com/go-chi/chi) con middlewares de logging,
  recuperación de *panics*, CORS y timeouts.
- ✅ Configuración 100% por variables de entorno (con `.env` opcional para desarrollo).
- ✅ Migraciones SQL embebidas (`embed.FS`), versionadas y aplicadas automáticamente al arrancar.
- ✅ Paginación y búsqueda (`?q=`, `?page=`, `?page_size=`) en el listado de libros.
- ✅ Validación de datos de entrada con errores por campo (HTTP 422).
- ✅ Manejo de errores consistente vía `errors.Is`/`errors.As`, sin filtrar detalles internos al cliente.
- ✅ Logging estructurado con `log/slog` (JSON).
- ✅ Endpoints de salud `/healthz` (liveness) y `/readyz` (readiness, verifica la BD).
- ✅ Apagado ordenado ante `SIGINT`/`SIGTERM`.
- ✅ Tests unitarios de `service` y `transport/rest` con dobles de prueba (sin base de datos real).
- ✅ Dockerfile multi-stage + `docker-compose.yml` (API + PostgreSQL) para levantar todo con un comando.

## Estructura del proyecto

```
.
├── cmd/
│   └── api/
│       └── main.go            # Punto de entrada: wiring, migraciones, arranque y apagado
├── internal/
│   ├── config/                # Carga y valida configuración desde variables de entorno
│   ├── model/                 # Entidades de dominio (Book) y reglas de validación
│   ├── service/                # Casos de uso / lógica de negocio
│   ├── store/                  # Persistencia en PostgreSQL + migraciones embebidas
│   │   └── migrations/
│   └── transport/
│       └── rest/               # Router, middlewares y handlers HTTP
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── .env.example
├── go.mod
└── go.sum
```

Cada capa depende únicamente de interfaces declaradas por el paquete que las
consume (`service.BookStore`, `rest.BookService`), lo que permite reemplazar
la base de datos o probar cada capa de forma aislada sin mocks externos.

## Requisitos

- Go 1.26+
- PostgreSQL 13+ (o Docker, ver más abajo)

## Puesta en marcha

### Opción A: Docker Compose (recomendado)

Levanta la API y PostgreSQL con un solo comando:

```bash
make docker-up
```

La API quedará disponible en `http://localhost:8080`.

### Opción B: Local

1. Copia `.env.example` a `.env` y ajusta las credenciales de tu PostgreSQL local:

   ```bash
   cp .env.example .env
   ```

2. Ejecuta la aplicación:

   ```bash
   make run
   # o bien: go run ./cmd/api
   ```

Al arrancar, la aplicación crea automáticamente las tablas necesarias
(`books`, `schema_migrations`) si no existen.

## Variables de entorno

| Variable                  | Descripción                                  | Default        |
|----------------------------|-----------------------------------------------|----------------|
| `SERVER_PORT`               | Puerto HTTP                                   | `8080`         |
| `SERVER_READ_TIMEOUT`       | Timeout de lectura de peticiones              | `5s`           |
| `SERVER_WRITE_TIMEOUT`      | Timeout de escritura de respuestas            | `10s`          |
| `SERVER_IDLE_TIMEOUT`       | Timeout de conexiones keep-alive              | `60s`          |
| `SERVER_SHUTDOWN_TIMEOUT`   | Tiempo máximo para el apagado ordenado        | `10s`          |
| `DB_HOST`                   | Host de PostgreSQL                            | `localhost`    |
| `DB_PORT`                   | Puerto de PostgreSQL                          | `5432`         |
| `DB_USER`                   | Usuario de PostgreSQL                         | `postgres`     |
| `DB_PASSWORD`               | Password de PostgreSQL                        | *(vacío)*      |
| `DB_NAME`                   | Nombre de la base de datos                    | `golang_back`  |
| `DB_SSLMODE`                | Modo SSL (`disable`, `require`, ...)          | `disable`      |
| `DB_MAX_OPEN_CONNS`         | Máximo de conexiones abiertas al pool         | `25`           |
| `DB_MAX_IDLE_CONNS`         | Máximo de conexiones inactivas en el pool     | `25`           |
| `DB_CONN_MAX_LIFETIME`      | Tiempo de vida máximo de una conexión         | `5m`           |

## Endpoints

Todas las respuestas son JSON. Los errores siguen el formato:

```json
{ "error": { "message": "descripción del error", "fields": { "title": "el título es requerido" } } }
```

| Método   | Ruta                     | Descripción                                      |
|----------|--------------------------|---------------------------------------------------|
| `GET`    | `/healthz`               | Liveness check                                     |
| `GET`    | `/readyz`                | Readiness check (verifica conexión a la BD)        |
| `GET`    | `/api/v1/books`          | Lista libros (`?q=`, `?page=`, `?page_size=`)      |
| `POST`   | `/api/v1/books`          | Crea un libro                                      |
| `GET`    | `/api/v1/books/{id}`     | Obtiene un libro por ID                            |
| `PUT`    | `/api/v1/books/{id}`     | Actualiza un libro por ID                          |
| `DELETE` | `/api/v1/books/{id}`     | Elimina un libro por ID                            |

### Ejemplo: crear un libro

```bash
curl -X POST http://localhost:8080/api/v1/books \
  -H "Content-Type: application/json" \
  -d '{"title":"Cien años de soledad","author":"Gabriel García Márquez","year":1967}'
```

### Ejemplo: listar con búsqueda y paginación

```bash
curl "http://localhost:8080/api/v1/books?q=garcia&page=1&page_size=10"
```

Respuesta:

```json
{
  "items": [ { "id": 1, "title": "Cien años de soledad", "author": "Gabriel García Márquez", "year": 1967, "created_at": "...", "updated_at": "..." } ],
  "page": 1,
  "page_size": 10,
  "total_items": 1,
  "total_pages": 1
}
```

## Desarrollo

```bash
make run     # Ejecuta la API localmente
make test    # Corre la suite de tests con cobertura
make lint    # go vet + gofmt
make tidy    # Sincroniza go.mod/go.sum
```

## Migraciones

Las migraciones viven en `internal/store/migrations/*.sql`, se embeben en el
binario en tiempo de compilación (`//go:embed`) y se aplican de forma
idempotente al arrancar, registrando cada versión aplicada en la tabla
`schema_migrations`. Para añadir un cambio de esquema, agrega un nuevo
archivo `000N_descripcion.sql` en esa carpeta.
