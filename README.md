# INCLD Go SDK

Trusted-server Go client for the INCLD v1 API. The module path is `github.com/incld-oss/incld-go`; resource types live in the `client` package.

```go
incld := client.NewClient(os.Getenv("INCLD_SECRET_KEY"))

schedule, err := incld.Schedules.CreateWithOptions(client.CreateScheduleParams{
    ExternalUserID: "user_123",
    Action:         "sync_contacts",
    Recurrence:     recurrence,
    Timezone:       "Australia/Melbourne",
}, &client.RequestOptions{Context: ctx, IdempotencyKey: requestID})
```

The client exposes Actions, Schedules, Runs, Sessions, Approvals, Approval Policies, Audit, and Bulk operations. List calls return `Page[T]` with cursor metadata. API failures return typed `APIError` values containing status, stable code, safe message, optional validation fields, and request ID.

Public schedule commands reference an action identifier through `Action`; action database UUIDs are not part of the contract. Approval checks require resource type, resource ID, action, and requester identity on trusted server calls. Audit uses `/audit-events`, and Bulk uses `/bulk-operations`.

## Configure the client

```go
incld := client.NewClient(os.Getenv("INCLD_SECRET_KEY"))
incld.BaseURL = "https://api.incld.dev/v1"
incld.HTTPClient = &http.Client{Timeout: 10 * time.Second}
```

This is a trusted-server client; never send the secret key to a browser. To use INCLD React packages in a Go application, implement a same-origin proxy that authenticates your user, scopes reads, injects actor identity, and allowlists operations.

## Services and methods

| Service | Methods |
| --- | --- |
| `Actions` | `Define`, `List`, `Get` |
| `Schedules` | `Create`, `CreateWithOptions`, `List`, `Get`, `Update`, `Delete`, `Pause`, `Resume`, `ListRuns`, `ListEvents`, `Preview` |
| `Runs` | `Get`, `List`, `ListPage` |
| `Approvals` | `List`, `Get`, `Check`, `Create`, `CreateWithOptions`, `Update`, `Delete`, `Decide`, `DecideWithOptions`, `Approve`, `Reject`, `Cancel`, `Revoke`, `Events` |
| Approval policies | `ListPolicies`, `GetPolicy`, `CreatePolicy`, `UpdatePolicy`, `DeletePolicy` on `Approvals` |
| `Audit` | `List`, `Get`, `Create` |
| `Bulk` | `List`, `Get`, `Create`, `CreateWithOptions`, `Chunks`, `Events`, `Cancel` |
| `Sessions` | `Create` |

## Context, idempotency, and errors

Use `RequestOptions` where an `*WithOptions` method is available:

```go
operation, err := incld.Bulk.CreateWithOptions(
    &client.CreateBulkOperationParams{
        Action: "sync_contacts",
        Items: items,
        ChunkSize: 250,
    },
    &client.RequestOptions{
        Context: ctx,
        IdempotencyKey: "contact-import:" + importID,
    },
)
```

Use stable idempotency keys for retried creates and decisions. Collection methods return `Page[T]`; pass `Meta.NextCursor` to the next list call unchanged.

```go
var apiErr *client.APIError
if errors.As(err, &apiErr) {
    log.Printf("INCLD %d %s request=%s fields=%v",
        apiErr.StatusCode, apiErr.Code, apiErr.RequestID, apiErr.Fields)
}
```

## Verify signed delivery

```go
import "github.com/incld-oss/incld-go/webhook"

func handleIncld(w http.ResponseWriter, r *http.Request) {
    event, err := webhook.VerifyAndParse(r, os.Getenv("INCLD_WEBHOOK_SECRET"))
    if err != nil {
        http.Error(w, "invalid webhook", http.StatusUnauthorized)
        return
    }
    if err := dispatch(event); err != nil {
        http.Error(w, "delivery failed", http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}
```

`VerifyAndParse` reads the exact request body and `Incld-Signature`. For raw bytes, use `VerifyWebhookSignature(payload, signature, secret, toleranceSeconds)`; a tolerance of zero defaults to 300 seconds. Deduplicate application side effects using the event/delivery identity before returning success.
