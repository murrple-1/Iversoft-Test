# Iversoft-Test

## Build

### Requirements
* Go 1.9
    - Untested with other versions
    - Ensure your GO environment variables are set up
* MySQL/MariaSQL
    - Recommened DB/user setup
        - DB Name: `iversoft`
        - User: `iversoft_user`
        - Password: `password`
    - Ensure user is granted `SELECT, INSERT, UPDATE, DELETE`, at least
    - Run `mysql -u root -p iversoft < dump.sql`

### Build

#### Linux/MacOS
Run `./scripts/build.sh`. `./iversoft.out` should be created.

#### Windows
Not tested, but use `./scripts/build.sh` as a basis for the build steps.

### Run

#### Linux/MacOS
Run `./iversoft.out`. Alternatively, modify and run `./scripts/run.sh`. It contains a commented-out list of available environment variables, which you can edit to suit your setup.

Open your web browser to `http://localhost:<PORT>` (`PORT` defaults to `8080`).

#### Windows
Not tested, but use `./scripts/run.sh` as a basis for the run steps.

### API Docs

If you use [Postman](https://www.getpostman.com), you can import `./test/Iversoft.postman_collection.json` to play with the API.

`GET /api/user/<id>`

Returns a JSON object for the requested user.

`POST /api/user`

Creates a new user.

Body:
```
{
    "username": <string>,
    "email": <string>,
    "roleLabel": <string>,
    "address": {
        "address": <string|null>, // optional
        "city": <string|null>, // optional
        "province": <string|null>, // optional
        "country": <string|null>, // optional
        "postalCode": <string|null> // optional
    }
}
```

`PUT /api/user/<id>`

Edits an existing user.

Body:
```
{
    "email": <string>, // optional
    "roleLabel": <string>, // optional
    "address": {
        "address": <string|null>, // optional
        "city": <string|null>, // optional
        "province": <string|null>, // optional
        "country": <string|null>, // optional
        "postalCode": <string|null> // optional
    }
}
```

`DELETE /api/user/<id>`

Deletes a user

`GET /api/users?[count=<int>][&skip=<int>]`

Returns a paginated list of users (JSON array of objects).

`count`: similar to SQL `LIMIT`. Defaults to `50`, max `1000`.
`skip`: similar to SQL `OFFSET`. Defaults to `0`.