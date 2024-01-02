# MTL Voting API

A voting API backend.

# Starting the API

Start up the API and the database using docker compose.

```bash
    docker-compose up
```

This will start a MySQL database at port 3306 named and start the API at port 3000.

The database has schema named **mtl**.
The username is `mtl` and the password is `password`.

This database can access with `mysql://mtl:password@tcp(localhost:3306)/mtl`.

The API backend can be accessed at `http://localhost:3000`

# Credentials

The initial API user credentials are

 username: `test`
 password: `password`


username: `John`
password: `password`

username: `test_user`
password: `password`

# APIs

#### GET /
A health checking endpoint.

## login

### POST /login
Logging in using the credentials above in exchange for a JWK token for subsequent usage in other endpoints.

#### Request body

```json
{
  "username": "test", 
  "password": "password"
}
```

#### Response

```json
{
  "token": "<jwt_token>"
}
```

## vote_item

Allow to create /update /delete vote items. 
Endpoints require the token from post /login as bearer authentication header.

### GET /vote_items

Retrieve vote items sorted by number of voted in descending order.

#### Response

```json
{
 "data": [
  {
   "id": 1,
   "name": "Hamburger",
   "description": "description",
   "vote_count": 3
  },
  {
   "id": 2,
   "name": "Pizzar",
   "description": "description",
   "vote_count": 0
  }
 ]
}
```

### POST /vote_items

Create a vote item.

#### Request

```json
{
 "name": "<name>",
 "description": "<description>"
}
```

#### Response

```json
{
 "data": {
  "id": 3,
  "name": "<name>",
  "description": "<description>"
 }
}
```

### PATCH /vote_items/:id

Update a vote item by the given id.

#### Request

```json
{
 "name": "<name>",
 "description": "<description>"
}
```

#### Response

```json
{
 "data": {
   "success": true
 }
}
```

### DELETE /vote_items/:id

Delete a vote item by the given id. The item must have 0 vote to be deletable.

#### Request

#### Response

```json
{
 "data": {
   "success": true
 }
}
```

### POST /vote_items/reset

Delete every vote items, regardless of their voting status.

#### Request

#### Response

```json
{
 "data": {
   "success": true
 }
}
```

## vote

### POST /vote_items/:id/vote

Give a vote to the given item ID.

#### Response

```json
{
 "data": {
   "success": true
 }
}
```

### POST /vote_items/:id/reset

Clear votes from the vote item.

#### Response

```json
{
 "data": {
   "success": true
 }
}
```
