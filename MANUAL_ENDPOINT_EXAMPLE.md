## Example of manual json file to put in var/config dir

```json
[
  {
    "name": "get_all_brokers",
    "path": "http://localhost/v1/query-filters/assets/brokers",
    "verb": "GET",
    "requestHeader": [
      {
        "parameter": "Accept",
        "value": "application/json"
      }
    ],
    "requestPath": [],
    "requestBody": "",
    "results": [
      {
        "code": 200
      },
      {
        "code": 401
      },
      {
        "code": 403
      },
      {
        "code": 500
      }
    ]
  }
]
```