# Endpoints

This document outlines the endpoints that Craig provides.

## GET /searches

List search criteria for Craig.

### Response

The response structure is as follows:

```json
{
  "searches": [
    {
      "id": <id>,
      "region": <region>,
      "category": <category>,
      "terms": [<term>, ...],
      "nh": [<neighborhood>, ...]
    },
    { ... },
    ...
  ]
}
```

* `id`: The Craig-assigned ID of the search criteria.
* `region`: The domain name prefix that represents the region in which to
  perform the search, e.g. `"sfbay"`
* `category`: The search category retrieved from URL GET parameters from
  Craigslist.
* `terms`: A list of search strings (each joined with OR)
* `nh`: A list of neighborhood ids, provided as integers.

## POST /searches

Add a new search criterion.

### Request Payload

The payload for the search criterion is as follows:

```json
{
  "region": <region>,
  "category": <category>,
  "terms": [<term>, ...],
  "nh": [<neighborhood>, ...]
}
```

* `region`: The domain name prefix that represents the region in which to
  perform the search, e.g. `"sfbay"`
* `category`: The search category retrieved from URL GET parameters from
  Craigslist.
* `terms`: A list of search strings (each joined with OR)
* `nh`: A list of neighborhood ids, provided as integers.

### Response

#### 204

The response was successfully processed.

## DELETE /searches/:search_id

Remove a search criterion

## Notes

There is no remove endpoint. To modify a search criterion, remove it and create
a new one.
