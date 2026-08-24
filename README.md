# incld Go SDK

Trusted-server Go client for the incld v1 API. The module path is `github.com/incld-dev/incld-go`; resource types live in the `client` package.

```go
admin := client.NewClient(os.Getenv("INCLD_SECRET_KEY"))
incld := admin.WithScope("org_123", "user_123")

schedule, err := incld.Schedules.CreateWithOptions(client.CreateScheduleParams{
    Action:                 "sync_contacts",
    Recurrence:             recurrence,
    Timezone:               "Australia/Melbourne",
}, &client.RequestOptions{Context: ctx, IdempotencyKey: requestID})
```

The client exposes Actions, Schedules, Runs, Sessions, Approvals, Approval Policies, Audit, and Bulk operations. List calls return `Page[T]` with cursor metadata. API failures return typed `APIError` values containing status, stable code, safe message, optional validation fields, and request ID.

Public schedule commands reference an action identifier through `Action`; action database UUIDs are not part of the contract. Approval checks require resource type, resource ID, action, and requester identity on trusted server calls. Audit uses `/audit-events`; `Audit.Tombstone` provides trusted-server PII erasure with a reason-coded record. Bulk uses `/bulk-operations`.

## Configure the client

```go
admin := client.NewClient(os.Getenv("INCLD_SECRET_KEY"))
incld := admin.WithScope("org_123", "user_123")
incld.BaseURL = "https://api.incld.dev/v1"
incld.HTTPClient = &http.Client{Timeout: 10 * time.Second}
```

This is a trusted-server client; never send the secret key to a browser. The value returned by `NewClient` is a project administrator client. Use `WithScope(organizationID, userID)` for end-user request paths or an empty `userID` for organization-wide server work. The API then constrains direct ID lookups as well as lists and mutations. To use incld React packages in a Go application, implement a same-origin proxy that authenticates the user, requires their active organization, injects both identities, and allowlists operations.
`WithScope` requires a non-empty organization ID. Project-wide action definition remains an administrator operation and should run through the unscoped client during deployment or startup synchronization.

## Services and methods

| Service | Methods |
| --- | --- |
| `Actions` | `Define`, `List`, `Get` |
| `Schedules` | `Create`, `CreateWithOptions`, `List`, `Get`, `Update`, `Delete`, `Pause`, `Resume`, `ListRuns`, `ListEvents`, `Preview` |
| `Runs` | `Get`, `List`, `ListPage` |
| `Approvals` | `List`, `Get`, `Check`, `Create`, `CreateWithOptions`, `Update`, `Delete`, `Decide`, `DecideWithOptions`, `Approve`, `Reject`, `Cancel`, `Revoke`, `Events` |
| Approval policies | `ListPolicies`, `GetPolicy`, `CreatePolicy`, `UpdatePolicy`, `DeletePolicy` on `Approvals` |
| `Audit` | `List`, `Get`, `Create`, `Tombstone` |
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
    log.Printf("incld %d %s request=%s fields=%v",
        apiErr.StatusCode, apiErr.Code, apiErr.RequestID, apiErr.Fields)
}
```

## Verify signed delivery

```go
import "github.com/incld-dev/incld-go/webhook"

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
