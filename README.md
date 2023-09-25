# Block-Indexer

Utilize an indexer to scan the historical block data into database and listen the latest block in parallel. Then, start an API service to fetch the stored block data rapidly from database without needing to go through blockchain.

## Deployment

1. Start MySQL server
2. Create config: `.env` (Can copy the schema from `.env-template`)
3. Install required external packages: `go get ./...`
4. Migrate database: `go run db/migrate.go` in db directory
5. Scan block data into database: `go run cmd/indexer/indexer.go`
6. Start API service for fetching stored block data: `go run cmd/app/app.go`

## Usage

### Indexer service

Scan block data into database

#### Parameter Adjustment

```go
// parameters for scanning historical blocks
start := 0 // start block number
end := 33585920 // end block number
workerCountOld := 5 // worker count to scan blocks in parallel for historical blocks
batchInsertCountOld := 10 // total block data count will be insert into DB every I/O for historical blocks

// parameters for keeping scanning the latest blocks
workerCountNew := 1 // worker count to scan blocks in parallel for new blocks
batchInsertCountNew := 1 // total block data count will be insert into DB every I/O for new blocks
fetchInterval := 5 * time.Second // periodically fetching interval
```

#### Running

```shell
go run cmd/indexer/indexer.go
```

### API service

Fetch block data stored in database

#### API

[GET] /blocks?limit=n

- get a limited count of block data, excluding transaction hashes
- `n` is the count of block

[GET] /blocks/:id

- get the specified block data, including its transaction hashes
- `:id` is block number

[GET] /transaction/:txHash

- get the specified transactions data, including its event logs
- `:txHash` is transaction hash

#### Running

```shell
go run cmd/app/app.go
```
