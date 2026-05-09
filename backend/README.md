

## DB

```sh
brew services start postgresql@14
```

```sql
psql --host=localhost --dbname=go_auth_v1 
```

```sh
goose up # picks GOOSE vars from .env
```

## Test flows

See Postman `go-auth-v1` for existing requests.

1. Create a user using `POST /signup`
Check they're visible in `GET /users`

2. `POST /login` for the user, save the JWT token & refresh token

3. try `GET /users/:id`, it should work until the JWT is not expired
```
Authorization: Bearer <jwt>
```

4. Once JWT expires - try `/refresh-token` to get a new JWT token

Get the new JWT
```
Authorization: Bearer <refresh-token>
```

Use on subsequent requests, similar to step 3.
