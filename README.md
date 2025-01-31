# Umbilical Choir: Release Manager

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
The endpoint for the children to send the summary results of one (or more) stage.
Sample input:
```json
{
  "id": "<childID>",
  "release_id": 0,
  "stage_summaries": [
    {
      "ProxyTimes": {
        "Median": 8.24175,
        "Minimum": 7.03625,
        "Maximum": 21.592958
      },
      "F1TimesSummary": {
        "Median": 7.912958,
        "Minimum": 6.979334,
        "Maximum": 19.588959
      },
      "F2TimesSummary": {
        "Median": 8.200792,
        "Minimum": 7.104208,
        "Maximum": 21.5505
      },
      "F1ErrRate": 0,
      "F2ErrRate": 0,
      "status": "Completed",
      "next_stage": "<next Stage name that child is going to run. in case of a 'rollout' or 'rollback', it will be nil>" 
    }
  ]
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

## Status codes
Upon a child's registration within a Release Strategy, their status code is initialized to `Todo`,
and the associated Release Strategy stages are assigned an initial status of `Pending`.

When the child is notified of a release (by polling `/poll`), it will call on `/release` to get the instructions, and the release (strategy) status will be `Doing` and the first stage's status will set as `InProgress`.
Then, the child is expected to run each stage and send the results to `/result`, where the stage's status also updates to either `Completed`, `Failure`, or `Error` (sent by the child agent).
If the stage is `Completed`, the next stage will be set to `InProgress` and the child will be notified of the next stage.
Otherwise, if there is no NextStage (a `rollout` or `rollback` action happened), which means the child finished/stopped the release strategy, the stage's status will be checked:
If the stage is a `Failure` or `Error`, the child's release strategy status will be set to `Failed`.
Otherwise, the stage should be `Completed` and the release strategy status will be set to `Done`.
In case of a `WaitForSignal` stage, the child polls `/end_stage` if the stage's status is `ShouldEnd` (or later) to end the stage, regardless of a minimum run time and call count.

Here is a status description for `ReleaseStatus` and `StageStatus` (enums):
### ReleaseStatus
state of a specific release (strategy) for a child:
- **No**: The child should not get the release instruction.
- **Todo**: Marked to get the release.
- **Doing**: The child is notified of the release.
- **Done**: The child has completed all stages.
- **Failed**: The release has failed.

### StageStatus
current state of stages in a release strategy for a child:

- **Pending**: The stage is pending.
- **InProgress**: The child is notified, and the stage is in progress.
- **SuccessWaiting**: Only for `WaitForSignal` stage type. The child received enough calls and was successful, now waiting for the parent to end the stage. The child also sends a preliminary aggregated results. 
- **ShouldEnd**: Only for `WaitForSignal` stage type. The child polls for it on `/end_stage` to finish a stage.
- **Completed**: The stage result has been received and the stage is completed.
- **Failure**: The stage result has been received as a failure.
- **Error**: The stage result has been received as an error.

