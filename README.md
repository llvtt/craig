# go-craigslist

A Craigslist Slack bot written in Go.

## Buildling

```sh
go get .
go build
```

## Configuration

Craig uses environment variables and a configuration file.

### Environment variables

* `CRAIG_SLACK_ENDPOINT` - Slack endpoint to use for posting messages

### Configuration file

The configuration file must be specified via the --config-file flag when running craig (by default, craig will read a file called `config.json` in the current
working directory).

Example configuration:

```json
{
  "db_type": "json",
  "db_file": "/path/to/db/file.json",
  "region": "sfc",
  "searches": [
    {
      "category": "zip",
      "terms": [""],
      "nh": [3]
    },
    {
      "category": "ata",
      "terms": ["lamp", "mirror", "queen bed"]
    }
  ]
}
```
