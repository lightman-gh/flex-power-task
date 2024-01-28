# flex-power-task

# Task 1

## How to deploy
```
    make docker-build
    docker-compose up -d
``` 

## How to test 
```
# Create user 
curl --location 'http://127.0.0.1:8081/users' \
--header 'Content-Type: application/json' \
--data '{
    "username": "username",
    "password": "password"
}'

# Create new trade 
curl --location 'http://127.0.0.1:8081/trades' \
--header 'Content-Type: application/json' \
--header 'Accept: application/json' \
--header 'Authorization: Basic bGlnaHRtYW46cGFzc3dvcmQ=' \
--data '{
  "id": "9f7f95c5-624a-4171-b2e2-44800995201b",
  "price": 5,
  "quantity": 32,
  "direction": "sell",
  "delivery_day": "2023-03-25",
  "delivery_hour": 12,
  "trader_id": "a94ae747-07e7-42cb-b6ac-c90c7e7d8e7a",
  "execution_time": "2024-01-31T15:04:05Z"
}'

# Get list 

curl --location 'http://127.0.0.1:8081/trades?trader_id=anton&delivery_day=2023-03-26' \
--header 'Accept: application/json' \
--header 'Authorization: Basic bGlnaHRtYW46cGFzc3dvcmQ=' \
--data ''
```