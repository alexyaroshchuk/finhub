# A data pipeline from a continuous stream of finhub

Every second you receive several values from FinHub and process them concurrently. For example, if in the first second you got the value `1, 2, 3, 4, 5` then the result of the work will be 3, if the window size is greater than `5`.

If the window size is `3`, then you will process the last `3` digits, that is `3,4,5`, respectively, the value `4` will be stored in the database. 

It turns out that our database will be updated with each new record received from finhub.

You can setup window size (WINDOW) variable in env.

### Database:

Struct of database:
id|avg_value|symbol|created_at
Example:
```
1|1215.65|BINANCE:ETHUSDT|2023-01-01 00:00:00.0000
2|16540.39|BINANCE:BTCUSDT|2023-01-01 00:00:00.0000
```

### Before start:

You should register for a free account at Finnhub to get an API Key.
You should set up variables in .env files (required variables presented in .env_example file). 

### Run app: 

make run

### Run app after changes:

make build

### Test app: 

make test
