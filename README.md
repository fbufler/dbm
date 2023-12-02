# dbm - the database monitoring tool

Ever had a problem to Identify your current replica master? Or just wanting to grep some metrics or uptimes about your database?
This small tool might help you!

## Usage

### Configuration

`dbm` reads in a yaml configuration that can be stored on multiple paths (`/etc/dbm/config.yaml`, `./config.yaml`, `$HOME/.dbm/config.yaml`)
The configuration is used for all commands.

```yaml
logfile: test.log # your logfile (if not set none will be created)
debug: false # enable this for debug logs 
test_timeout: 5 # maximum runtime of a test until it will be canceled
test_interval: 30 # interval when to run the next test
databases: # your database configurations
  - host: localhost
    port: 5432
    username: postgres
    password: postgres # <- I know this is not nice yet, I will try to provide another way for configuration soon
    database: postgres
    use_ssl: true
    SSLCertPath: /some/path
    SSLKeyPath: /some/path
    SSLRootCertPath: /some/path
    connection_timeout: 5
  - host: localhost
    port: 5433
    username: postgres
    password: postgres
    database: postgres
    connection_timeout: 5

```

### Locally

#### Setup

There is an easy way to setup your databases to support this tool.
As it does some read/write testing you need to have a table created.
Right now the table is hardcoded as `test`.
In order to set it up, just configure your databases and execute `dbm setup`

#### Local

If you want to make an initial test locally, please use `dbm local`.
This will just run the database tester and provide you with information about their behaving.

#### Serve

This command serves the results of the tester as json.
```json
{
    "results":  {
        "localhost:5432/postgres":  {
            "database":"localhost:5432/postgres",
            "connectable":true,
            "connection_time":81034,
            "writable":false,
            "write_time":0,
            "readable":false,
            "read_time":0,
            "timestamp":"2023-12-02T23:38:57.552428198+01:00"
        },
        "localhost:5433/postgres":  {
            "database":"localhost:5433/postgres",
            "connectable":true,
            "connection_time":248551,
            "writable":false,
            "write_time":0,
            "readable":false,
            "read_time":0,
            "timestamp":"2023-12-02T23:38:57.552428752+01:00"
        }
    }
}
```