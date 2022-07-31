Manual Migrations

These are sql queries, which will be run directly on the database.
First, run migrations on local server/stage server.
Once that's confirmed, take the ddl dump and put the ddl sql here as .sql extension.
Run those in prod along with proper care, as some ddl queries can block the entire production database.