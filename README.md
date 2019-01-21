# go-craigslist

A Craigslist Slack bot written in Go.

## Configuration

`craig` takes a JSON config file as its first positional argument. The config file has this format:

```json
{
  "slackEndpoint": "http://hooks.slack.com/incoming/webhook/url/here",
  "searchTerms": [
    "wardrobe",
    "desk",
    "mirror"
  ]
}
```
