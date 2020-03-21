# Craig

A Craigslist Slack bot written in Go. Craig searches for terms in neighborhoods and posts to Slack when he sees something interesting.

## Craig's Endpoints

### `GET /searches`

List the Searches currently configured with Craig.

### `POST /search`

Execute a search. Craig goes to Craigslist looking for the search terms, updates the db, and posts to slack.

### `GET /health`

Health check.

## Building

Download dependencies
```shell script
make deps
```

Run tests and build craig server
```sh
make test build
```
or just
```shell script
make
```

To cross-complile Craig, you will need to install musl-cross:

    brew install FiloSottile/musl-cross/musl-cross

## Running

Run craig server. Assumes there is a file `./.env` with secrets, or that secrets
already exist in environment. Assumes there is a config file `./dev.config.json`

```shell script
make run
```

## Deployment

Craig is deployed with Terraform.

Several variables are expected to be defined in the environment:

* `TF_VAR_slack_endpoint`
* `TF_VAR_aws_region`

Deploy Craig with the Make target:

```sh
make deploy
```

## Configuration

Craig uses environment variables and a configuration file.

### Environment variables

* `CRAIG_SLACK_ENDPOINT` - Slack endpoint to use for posting messages

### Configuration file

The configuration file must be specified via the `--config-file` flag when running craig (by default, craig will read a file called `config.json` in the current
working directory).

Example configuration:

```json
{
  "db_type": "json",
  "db_dir": "/tmp/craig",
  "region": "sfc",
  "searches": [
    {
      "category": "zip",
      "terms": [""],
      "nh": [3]
    },
    {
      "category": "ata",
      "terms": ["end table", "lamp", "mirror", "queen bed"]
    }
  ]
}
```
