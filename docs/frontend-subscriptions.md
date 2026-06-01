# Frontend Integration — Subscriptions

This document describes how the frontend should integrate with the subscription / payment flow exposed by `nnp-stream-backend`.

Payments are processed by **Shwary** mobile money (DRC, USD). The backend creates a pending subscription, asks Shwary to initiate the payment on the user's phone, and the user confirms the transaction by entering their mobile-money PIN on their device. The final status is delivered asynchronously by Shwary to the backend webhook, which then activates (or fails) the subscription.

The frontend never talks to Shwary directly.

---

## Flow overview

```
┌──────────┐  1. POST /subscriptions          ┌─────────┐  2. Initiate payment  ┌────────┐
│ Frontend │ ───────────────────────────────▶ │ Backend │ ───────────────────▶  │ Shwary │
└──────────┘                                  └─────────┘                       └────────┘
     ▲   │                                         ▲                                │
     │   │ 5. GET /subscriptions/:id (poll)        │ 4. POST /shwary-web-hook       │
     │   ▼                                         └────────────────────────────────┘
     │  (or Supabase realtime)                            (final status)
     │
     │ 3. 202 Accepted — { subscription: pending, transaction: pending }

  Watching a video:
┌──────────┐  POST /subscriptions/access      ┌─────────┐
│ Frontend │ ───────────────────────────────▶ │ Backend │  → { has_access, reason, subscription? }
└──────────┘                                  └─────────┘
```

1. User picks a plan and enters their DRC mobile-money phone number.
2. Frontend calls `POST /subscriptions`.
3. Backend responds **`202 Accepted`** with a `pending` subscription and a `pending` Shwary transaction.
4. The user receives a USSD / push prompt on their phone and enters their PIN.
5. Shwary calls the backend webhook with the final status. Backend flips the subscription to `active` (or `failed` / `canceled`).
6. Frontend polls the subscription row (or listens via Supabase realtime) until the status changes.

---

## Endpoint reference

### `POST /subscriptions`

Initiates a payment and creates a pending subscription.

**Request body**

| Field          | Type    | Required | Notes                                                                                |
| -------------- | ------- | -------- | ------------------------------------------------------------------------------------ |
| `user_id`      | string  | yes      | The authenticated user's id (Clerk / Supabase user id used throughout the platform). |
| `plan_id`      | string  | yes      | The id of an existing plan in the `plans` table.                                     |
| `phone_number` | string  | yes      | DRC mobile number in E.164 format. **Must start with `+243`**.                       |
| `amount`       | number  | yes      | Amount in **USD**. Must be `> 0`.                                                    |

Example:

```json
{
  "user_id": "user_2abc...",
  "plan_id": "plan_basic_monthly",
  "phone_number": "+243812345678",
  "amount": 4.99
}
```

**Successful response — `202 Accepted`**

```json
{
  "subscription": {
    "id": "sub_uuid",
    "created_at": "2026-05-28T10:12:00Z",
    "user_id": "user_2abc...",
    "plan_id": "plan_basic_monthly",
    "status": "pending"
  },
  "transaction": {
    "id": "shwary_tx_uuid",
    "amount": 4.99,
    "currency": "USD",
    "status": "pending",
    "recipientPhoneNumber": "+243812345678",
    "referenceId": "sub_uuid",
    "isSandbox": true,
    "createdAt": "2026-05-28T10:12:00Z",
    "updatedAt": "2026-05-28T10:12:00Z"
  }
}
```

Persist `subscription.id` on the client — it is the handle used to poll for final status.

**Error responses**

| Status | When                                              | Body shape              |
| ------ | ------------------------------------------------- | ----------------------- |
| `400`  | Missing/invalid field, phone not `+243…` E.164    | `{ "error": "..." }`    |
| `400`  | Plan id does not exist                            | `{ "error": "plan not found" }` |
| `500`  | Could not create the pending subscription row     | `{ "error": "failed to create subscription" }` |
| `502`  | Shwary refused the payment request                | `{ "error": "payment provider error" }` |

A `502` means the subscription row was created and immediately marked `failed`. The user can retry by calling `POST /subscriptions` again.

---

### `GET /subscriptions/:subscriptionID`

Returns the current state of a subscription. Use this to poll after `POST /subscriptions` until the row transitions out of `pending`.

**Path parameters**

| Param            | Type   | Notes                                            |
| ---------------- | ------ | ------------------------------------------------ |
| `subscriptionID` | string | The `subscription.id` returned by `POST /subscriptions`. |

**Successful response — `200 OK`**

```json
{
  "id": "sub_uuid",
  "created_at": "2026-05-28T10:12:00Z",
  "user_id": "<supabase user id>",
  "plan_id": "plan_basic_monthly",
  "status": "active",
  "started_at": "2026-05-28T10:12:48Z",
  "expires_at": "2026-06-28T10:12:48Z"
}
```

`status` is one of `pending`, `active`, `failed`, `canceled`. `started_at` and `expires_at` are only set once the subscription becomes `active` (and `expires_at` is omitted entirely for `life_time` plans).

**Error responses**

| Status | When                                  | Body shape                          |
| ------ | ------------------------------------- | ----------------------------------- |
| `400`  | Empty / missing `subscriptionID`      | `{ "error": "subscriptionID is required" }` |
| `404`  | No subscription matches that id       | `{ "error": "subscription not found" }`     |

This endpoint is **read-only** and safe to call on a tight interval.

---

### `POST /subscriptions/access`

Checks whether a user has access to a given video. The backend looks up the plans that gate the video (via the `video_plans` join table), then checks the user's subscriptions against those plans.

Call this before rendering the video player. If `has_access` is `false`, show the paywall / upsell instead of requesting a Mux playback URL.

**Request body**

| Field      | Type   | Required | Notes                                                                            |
| ---------- | ------ | -------- | -------------------------------------------------------------------------------- |
| `user_id`  | string | yes      | The Clerk user id (same value used in `POST /subscriptions`).                    |
| `video_id` | string | yes      | The id of the video the user is trying to watch. Matches `video_plans.video_id`. |

Example:

```json
{
  "user_id": "user_2abc...",
  "video_id": "vid_intro_2026"
}
```

**Successful response — `200 OK`**

```json
{
  "has_access": true,
  "reason": "active",
  "subscription": {
    "id": "sub_uuid",
    "user_id": "<supabase user id>",
    "plan_id": "plan_basic_monthly",
    "status": "active",
    "started_at": "2026-05-28T10:12:48Z",
    "expires_at": "2026-06-28T10:12:48Z"
  }
}
```

`reason` is one of:

| Reason            | `has_access` | Meaning                                                                 |
| ----------------- | ------------ | ----------------------------------------------------------------------- |
| `free_video`      | `true`       | No plan gates this video — anyone can watch.                            |
| `active`          | `true`       | User has an active, unexpired subscription to a plan that includes this video. |
| `pending`         | `false`      | User has a pending subscription — the payment hasn't completed yet. Show "payment processing" instead of the paywall. |
| `expired`         | `false`      | Subscription exists but is `failed`, `canceled`, or past `expires_at`.  |
| `no_subscription` | `false`      | User has never subscribed to any gating plan. Show the paywall.         |

The `subscription` field is included whenever the user has any record for a gating plan (active, pending, or expired) — useful for showing tailored messaging. It is omitted for `free_video` and `no_subscription`.

**Error responses**

| Status | When                                       | Body shape                          |
| ------ | ------------------------------------------ | ----------------------------------- |
| `400`  | Missing field, or Clerk user does not exist | `{ "error": "..." }`                |
| `500`  | Backend failed to resolve the answer       | `{ "error": "..." }`                |

---

## Subscription lifecycle

A subscription transitions through these statuses (string values stored in the row):

| Status     | Meaning                                                      |
| ---------- | ------------------------------------------------------------ |
| `pending`  | Awaiting Shwary callback. User must confirm on their phone.  |
| `active`   | Payment completed. `started_at` is set; `expires_at` is set for monthly/yearly plans (unset for `life_time`). |
| `failed`   | Shwary reported failure, or the payment provider call failed. |
| `canceled` | Shwary reported cancellation (user declined / timed out).    |

`expires_at` is computed from the plan's `billing_cycle`:

- `monthly` → `now + 1 month`
- `yearly` → `now + 1 year`
- `life_time` → no `expires_at` (subscription never expires)

The frontend should treat a subscription as entitled when `status === "active"` **and** (`expires_at` is null/unset OR `expires_at > now`).

---

## Waiting for the final status

The `POST /subscriptions` response only confirms the payment was *initiated*. Use one of these strategies to detect completion.

### Option A — Supabase realtime (recommended)

The `subscriptions` table is in Supabase. Subscribe to row-level updates filtered by the subscription id you got back:

```ts
import { createClient } from "@supabase/supabase-js";

const supabase = createClient(SUPABASE_URL, SUPABASE_ANON_KEY);

function watchSubscription(subscriptionId: string, onChange: (sub: Subscription) => void) {
  const channel = supabase
    .channel(`subscription:${subscriptionId}`)
    .on(
      "postgres_changes",
      {
        event: "UPDATE",
        schema: "public",
        table: "subscriptions",
        filter: `id=eq.${subscriptionId}`,
      },
      (payload) => onChange(payload.new as Subscription)
    )
    .subscribe();

  return () => supabase.removeChannel(channel);
}
```

Stop listening once `status` is one of `active`, `failed`, `canceled`.

### Option B — Polling the backend (recommended fallback)

If realtime is not wired up, poll the backend's `GET /subscriptions/:subscriptionID` endpoint. This is the right choice when you don't want the frontend to hold Supabase credentials, or when you want a single integration surface.

```ts
async function pollSubscription(
  id: string,
  { intervalMs = 3000, timeoutMs = 120_000 }: { intervalMs?: number; timeoutMs?: number } = {},
): Promise<Subscription> {
  const deadline = Date.now() + timeoutMs;
  while (Date.now() < deadline) {
    const res = await fetch(`${API_BASE_URL}/subscriptions/${id}`);
    if (!res.ok) throw new Error(`poll failed: ${res.status}`);
    const sub = (await res.json()) as Subscription;
    if (sub.status !== "pending") return sub;
    await new Promise((r) => setTimeout(r, intervalMs));
  }
  throw new Error("subscription_timeout");
}
```

Use a 2–4 second interval and a 2-minute hard timeout. Shwary mobile-money confirmations typically arrive in under 60 seconds.

### Option C — Polling Supabase directly

If the frontend already has a Supabase client and you'd rather skip the extra backend hop, query the row directly. Functionally equivalent to Option B:

```ts
const { data, error } = await supabase
  .from("user_subscriptions")
  .select("*")
  .eq("id", subscriptionId)
  .single();
```

Stop polling once `data.status !== "pending"`.

---

## Reference TypeScript client

```ts
type SubscribeRequest = {
  user_id: string;
  plan_id: string;
  phone_number: string; // must start with +243
  amount: number;       // USD
};

type SubscriptionStatus = "pending" | "active" | "failed" | "canceled";

type Subscription = {
  id: string;
  created_at: string;
  user_id: string;
  plan_id: string;
  status: SubscriptionStatus;
  started_at?: string;
  expires_at?: string;
};

type ShwaryTransaction = {
  id: string;
  amount: number;
  currency: string;
  status: string;
  recipientPhoneNumber: string;
  referenceId: string;
  txHash?: string;
  failureReason?: string;
  completedAt?: string;
  isSandbox: boolean;
  createdAt: string;
  updatedAt: string;
};

type SubscribeResponse = {
  subscription: Subscription;
  transaction: ShwaryTransaction;
};

type VideoAccessReason =
  | "active"
  | "pending"
  | "expired"
  | "no_subscription"
  | "free_video";

type CheckVideoAccessResponse = {
  has_access: boolean;
  reason: VideoAccessReason;
  subscription?: Subscription;
};

export async function subscribeToPlan(body: SubscribeRequest): Promise<SubscribeResponse> {
  const res = await fetch(`${API_BASE_URL}/subscriptions`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    const { error } = (await res.json().catch(() => ({ error: res.statusText }))) as { error?: string };
    throw new Error(error ?? `subscribe failed: ${res.status}`);
  }
  return res.json();
}

export async function getSubscription(id: string): Promise<Subscription> {
  const res = await fetch(`${API_BASE_URL}/subscriptions/${id}`);
  if (!res.ok) {
    const { error } = (await res.json().catch(() => ({ error: res.statusText }))) as { error?: string };
    throw new Error(error ?? `get subscription failed: ${res.status}`);
  }
  return res.json();
}

export async function checkVideoAccess(
  userId: string,
  videoId: string,
): Promise<CheckVideoAccessResponse> {
  const res = await fetch(`${API_BASE_URL}/subscriptions/access`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ user_id: userId, video_id: videoId }),
  });
  if (!res.ok) {
    const { error } = (await res.json().catch(() => ({ error: res.statusText }))) as { error?: string };
    throw new Error(error ?? `access check failed: ${res.status}`);
  }
  return res.json();
}
```

---

## Gating the video player

Wrap the video page in an access check. The pattern is:

1. Call `POST /subscriptions/access` with the current Clerk user id and the video id.
2. Branch on `reason`:

```ts
const access = await checkVideoAccess(clerkUserId, videoId);

switch (access.reason) {
  case "active":
  case "free_video":
    // Render the Mux player.
    return <Player videoId={videoId} />;

  case "pending":
    // Don't show the paywall — the user already paid and is waiting.
    return <PaymentProcessingNotice subscriptionId={access.subscription!.id} />;

  case "expired":
    return <RenewalPrompt subscription={access.subscription!} />;

  case "no_subscription":
  default:
    return <Paywall videoId={videoId} />;
}
```

Re-run the check on:

- initial page load,
- successful return from the subscribe flow (i.e. once polling sees `status === "active"`),
- whenever the user's session changes.

Never cache `has_access: true` across sessions or beyond `subscription.expires_at`. The backend is the source of truth — recheck whenever the user reopens the app.

---

## UX guidance

- **Validate the phone number client-side** to E.164 `+243` before submitting. The backend rejects everything else with `400`, but catching it in the form is faster.
- **Show a "Check your phone" state** as soon as `POST /subscriptions` returns `202`. Display the masked phone number you submitted so the user knows where the prompt will arrive.
- **Disable the submit button while pending** to prevent duplicate payments. If the user wants to retry, only re-enable it after the subscription transitions to `failed` or `canceled`.
- **Surface the failure reason** when status is `failed`. The reason lives on the `payment_transactions` row (`failure_reason` column) — read it if you need to display it.
- **Do not assume `active` until you see the row flip.** The synchronous response is always `pending`; granting entitlement before the webhook fires will hand out access to unpaid users.
- **Idempotency.** The backend webhook is idempotent, but `POST /subscriptions` is **not** — every call creates a new pending subscription and a new Shwary payment. Guard the submit action accordingly.

---

## Listing plans

Plans are stored in the `plans` table in Supabase. There is no dedicated backend list endpoint yet — read them directly via the Supabase client:

```ts
const { data: plans } = await supabase
  .from("plans")
  .select("id, name, description, billing_cycle, created_at")
  .order("created_at", { ascending: true });
```

Each plan has a `billing_cycle` of `monthly`, `yearly`, or `life_time`. The frontend is responsible for displaying the price; the price the user is charged is whatever the frontend sends as `amount` in the subscribe request. Keep pricing in a single place (Supabase row, config file, etc.) so the displayed price and the submitted `amount` cannot drift apart.

---

## Local testing

Set `SHWARY_SANDBOX=true` on the backend to route payment initiation to Shwary's sandbox. Sandbox transactions still trigger the webhook and complete the subscription lifecycle end-to-end, so the frontend flow can be exercised without real money.
