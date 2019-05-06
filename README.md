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
* `CRAIG_S3_BUCKET` - S3 bucket to use for the database and config file

### Configuration file

The configuration file must be called `config.json` and placed in the current
working directory.

Example configuration:

```json
{
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
