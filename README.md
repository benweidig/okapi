# Okapi [![Build Status](https://travis-ci.org/benweidig/okapi.svg?branch=master)](https://travis-ci.org/benweidig/okapi)

A database migration library for Go.

## Why?

Because I can. Seriously... While I was working on another project I realized I would need some kind of database migration functionality.
After looking around I found [darwin](https://github.com/GuiaBolso/darwin) but it didn't fit my needs perfectly so I decided to write one
myself.

## Features

- Batteries included (at least for stdlib sql, dialects only sqlite, mysql and postgres)
- Checksums for changesets
- Easy extendable
- Bring-your-own-logging

## Get it

```
go get -u github.com/benweidig/okapi
```

## Use it

1. Write `okapi.Changeset`s
2. Create `okapi.Driver`
3. Create `okapi`
3. Migrate
4. ...
5. Profit!

```
    // Step 1: Write changesets
    changesets := []okapi.Changeset{
        {
            ID:      "001",
            Comment: "create first table",
            Script: `
                CREATE TABLE my_table (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    name VARCHAR(64) NOT NULL
                );`,
        },
    }

    // Step 2: Create driver
    driver := MyCustomerDriver{}
    ... or skip this with the SqlDriver, see below
    
    // Step 3: Create okapi
    o, err := okapi.New(myCustomDriver, changesets)

    ... or simpler with SqlDriver
    driver := 
	o, err := okapi.WithSQLDriver(
		conn.DB,
		"sqlite3",
        changesets
    )

    // Step 4: Migrate
    err = o.Migrate()
```

## Changesets

| Field       | Type   | Size* | Description                                                                               |
| ----------- | ------ | ----- | ----------------------------------------------------------------------------------------- |
| ID          | string | 255   | An unique identifier of a changeset                                                       |
| SkipOnError | bool   | n/a   | Failed changesets will be skipped (instead of failing the whole migration) and not re-ran |
| Script      | string | n/a   | SQL script, the actual changeset                                                          |
| Comment     | string | 255   | A comment for easier debugging while reading the data on the server, can be empty         |

*=Max Size for `.WithSqlDriver(...)`

## Uniqueness

The uniqunes/immutability of a changeset if defined by 2 fields:
- `ID`: Must be unique over all changesets
- `Script`: The script must not change, except for whitespace.

## Caveats

- Changesets are script-based, not easily usable for different DBs
- Only sqlite3, mysql and postrgres built-in so far
- Not battle-tested 
- Whitespace ignorance isn't perfect, e.g. `SOME_SQL ( COLUMN )` != `SOME_SQL (COLUMN)`.


## Logging

To not enforce a specific logger a channel can be provided that will be receiving `okapi.Info` structs:

```
ch := make(chan okapi.Info)
go func() {
    for {
        log := <-ch
        ... your logging here
    }
}()
okapi.Notifier(ch)
```


## License

MIT. See [LICENSE](LICENSE).
