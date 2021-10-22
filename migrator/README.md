
## Postgres Database Migration

1. Start postgres server
```
make postgres-up
```
2. Create database
```
make create-db
```
3. Up migrations
```
make migrate-up
```
4. Down migrations
```
make migrate-down
```
5. Drop database
```
make drop-db
```
6. Stop postgres server
```
make postgres-down
```
