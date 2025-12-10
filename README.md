A very light-weight double-entry accounting ledger built for Go performance. Postgres is the database of choice because ledger workflows depend on transactions and relational guarantees; balances should never exist for missing books and transactional integrity is critical.

## Features
- BookId `1` is the company cashbook (recommended for traceability).
- Environment substitution supports `${VAR_NAME}` and JSON-shaped env values (`map[string]string`, `map[string]bool`, nested maps).
- Double-entry with per-book balances by asset and operation.
- Asset-agnostic (stocks, crypto, etc.).
- Skip balance rollups for specific books via `EXCLUDED_BALANCE_BOOK_IDS` (comma-separated). Default is to store all.
- Concurrent-safe operations without loading state into memory.
- DB constraint prevents negative `OVERALL` balances (BookId 1 excluded).
- Operation-level grouping (LIMIT_ORDER, MARKET_ORDER, DEPOSIT, WITHDRAW, TRADE, etc.).
- Metadata supports trade types (INTRA-DAY, QUARTERLY, etc.).
- No session/transaction advisory locks to maximize throughput.

> If `operationType` is not provided when fetching balances, `OVERALL` is used.

## API Docs
Import `collections/postman.json` into Postman for examples.

## Project Structure (high level)
- `cmd/http`, `cmd/grpc`, `cmd/combined`: server entrypoints.
- `internal/transport/http`, `internal/transport/grpc`: HTTP and gRPC handlers/routers.
- `internal/service/*`: domain services (books, operations) using repositories and DTOs.
- `internal/repository/*`: GORM repositories.
- `internal/config`, `internal/database`, `internal/logger`, `internal/util`, `internal/app`: shared utilities and setup.
- `domain/`, `dto/`, `models/`: shared types and persistence models.

## Configuration
- Viper loads configs from `internal/config/` based on `APP_ENV` (e.g., `local.yaml`).
- `.env` is used by default in `local/localhost`, or when `DOT_ENV=enable` in other envs.
- Env overrides use the `APP_` prefix automatically (`APP_SERVER_HTTPPORT`, `APP_DB_HOST`, etc.). Common aliases like `APP_PORT`/`HTTP_PORT`/`GRPC_PORT` are bound for convenience.
- YAML placeholders like `${VAR}` are expanded from the process environment before parsing.
- JSON env values remain supported for nested maps (see `SERVICE_TOKEN_WHITELIST` examples below).

Example `.env`:
```
APP_ENV=local
DOT_ENV=enable
RUN_MODE=release
DB_TYPE=postgres
DB_USER=postgres
DB_PASSWORD=admin
DB_HOST=127.0.0.1
DB_PORT=5432
DB_NAME=golang_ledger
DB_TABLE_PREFIX=golang_ledger_
DB_SSL_MODE=disable
JWT_SECRET=xxxx
EXCLUDED_BALANCE_BOOK_IDS=1,2,3
SERVICE_TOKEN_WHITELIST={"user_module":{"read":"abc","write":"cde"}}
```

## Running Locally
1. Install tools: `protoc`, `protoc-gen-go`, `protoc-gen-go-grpc`, `make`.
2. Install deps: `go mod tidy`.
3. Generate protos: `make ledger`.
4. Start servers:
   - HTTP: `go run cmd/http/main.go`
   - gRPC: `go run cmd/grpc/main.go`
   - Combined: `./run.sh` (loads `.env`, regenerates protos, runs combined with auto-reload)

## Deployment Notes
- Build steps must generate proto Go code (`make ledger`) before starting servers.
- Dockerfile intentionally omits proto generation; handle it in your build pipeline.
- For ECS + dotenv, store the env file in S3 and reference via `environmentFiles`.

## Operational Notes
1. Book create/update is idempotent by name: creates if missing, otherwise updates.
2. Client must ensure book uniqueness (e.g., uuid-v1 for debit/credit books per user).

### Book Types (suggested)
- CashBook `1`: main company book (can be negative to reflect spending).
- RevenueBook `2`: income; should not go negative.
- ThirdPartyVendorBook `3`: payables to vendors.
- ExpenseBook `4`: liabilities/expenses (e.g., BookId 1 -> BookId 4 for purchases).
- AssetBook `5`: company-owned assets.
- TDSBook `6`: deducted TDS to be remitted.
- IncomeTaxBook `7`: accrued income tax.

Formula: `Total Asset = Ⲉ(Liability Books) + Ⲉ(Equity Books)`

## TODO
1. Test cases. (Integration test added, modification required)
2. ~~BookId validation while creating operation.~~ (Done)
3. ~~Better Config, using yaml and viper.~~ (Done)
4. Example Ledger Client implementation to manage the ledger of a crypto trading org.
5. Customisable bookIds, based on type (asset or liability).
6. Reserve top 100 bookIds for company books. Migration to partition the balances table, such that below 100 ids should get in a specific partition, remaining should be partitioned based on hash.
7. Better file naming, code cleanup.
8. ~~Grpc support.~~ (Done)
