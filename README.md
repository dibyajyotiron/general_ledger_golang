Peoject depends on .env files, in production if used with ecs, make sure to create dotenv and store it inside s3.

This below part should be with the ecs task definition ->
```
"environmentFiles": [
  {
    "value": "arn:aws:s3:::s3_bucket_name/envfile_object_name.env",
    "type": "s3"
  }
]
```

To run the server, 
  1. create .env file ->
      ```
      APP_ENV = xx
      RUN_MODE = xx 
      DB_TYPE = xx
      DB_USER = xx
      DB_PASSWORD = xx
      DB_HOST = 127.0.0.1
      DB_PORT = 5432
      DB_NAME = xx
      DB_TABLE_PREFIX = xx
      DB_SSL_MODE = disable
      JWT_SECRET = xxxx
      ```

  2. install dependencies -> `go mod vendor && go mod tidy`
  3. run `./start_local_server.sh`

The server should run and have auto reload.

Notes:
1. Book Create/update method will create a book if the name of the book doesn't exist else it will update the book.
2. It is the ledger client's responsibility to maintain uniqueness of the book. 
3. To ensure uniqueness of the books for a given account holder, ledger client should create debit/credit books based on uuid-v1. 