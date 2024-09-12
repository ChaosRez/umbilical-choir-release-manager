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
  "new_release": "/release"
}
```
`new_release` will be empty string if there is no release and will contain the address of the new release if there is one.

### `/release` endpoint
The endpoint for the children to download the new release. 
No input is required. The release will be served as text.

### `/release/functions/{release_id}` endpoint
The endpoint for the children to download the functions of the new release.
The `release_id` is the id of the release that the child want to download the functions of.
It will return a zip file containing the functions of the release.
Each function is in a folder that is defined in the release file.

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

### Note
Sample function:
there is a `fns.zip` that contains `fns/sample_f1` and `fns/sample_f2` functions in tinyfaas format which is just for test.
