Umbilical Choir: Release Manager
--------------------------------

### `/poll` endpoint
The endpoint for the children to poll for new updates and sending their geo area and number of the children they have  
Sample input:
```json
{
  "geographic_area": {
    "type": "Polygon",
    "coordinates": [
      [
        [
          13.34138389963175,
          52.49855383364354
        ],
        [
          13.474766810586402,
          52.49855383364354
        ],
        [
          13.474766810586402,
          52.557371936926614
        ],
        [
          13.34138389963175,
          52.557371936926614
        ],
        [
          13.34138389963175,
          52.49855383364354
        ]
      ]
    ]
  },
  "number_of_children": 10,
  "id": "7cb606ee-fde1-4b2c-bffc-20f558fc2867"
}
```
Sample output:
```json
{
	"id": "7cb606ee-fde1-4b2c-bffc-20f558fc2867",
  "new_release": "<releaseID>"
}
```
`new_release` will be empty string if there is no release for the client (child) and will contain the ID of the (first) new release if there is one.

### `/release?childID=<ID>&releaseID=<releaseID>` endpoint
The endpoint for the children to download a specific release (usually given by the /poll endpoint).
releaseID and (child) ID should be passed as GET parameters. The release will be served as text.

### `/release/functions/{release_id}` endpoint
This endpoint allows children to download the functions associated with a specific release.
The `release_id` parameter specifies the release for which the functions are being requested.
It returns a `fns.zip` file containing the functions of the specified release.
Each function is stored in a separate folder as defined in the release file.
A sample `fns.zip` file (name matters) includes `fns/sample_f1` and `fns/sample_f2` functions in tinyfaas format.

### `/result` endpoint
The endpoint for the children to send the summary results of a release.
Sample input:
```json
{
    "id": "7cb606ee-fde1-4b2c-bffc-20f558fc2867",
    "release_id": 0,
    "release_summary":  {
      "ProxyTimes":
        {"Median":8.24175, "Minimum":7.03625, "Maximum":21.592958},
      "F1TimesSummary":
        {"Median":7.912958, "Minimum":6.979334, "Maximum":19.588959},
      "F2TimesSummary":
        {"Median":8.200792, "Minimum":7.104208, "Maximum":21.5505}, 
      "F1ErrRate":0, "F2ErrRate":0}
}
```

### `/end_stage` endpoint
The endpoint for the children to poll for an end of a stage signal, for the stages of type `WaitForSignal`.
Sample input:
```json
{
    "id": "7cb606ee-fde1-4b2c-bffc-20f558fc2867",
    "strategy_id": 21,
    "stage_name": "Canary Sieve Function"
}
```
Sample output:
```json
{
    "end_stage": true
}
```

