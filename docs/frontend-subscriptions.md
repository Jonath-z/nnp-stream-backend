# Frontend Integration — Subscriptions

This document describes how the frontend should integrate with the subscription / payment flow exposed by `nnp-stream-backend`.

Payments are processed by **Shwary** mobile money (DRC, USD). The backend creates a pending subscription, asks Shwary to initiate the payment on the user's phone, and the user confirms the transaction by entering their mobile-money PIN on their device. The final status is delivered asynchronously by Shwary to the backend webhook, which then activates (or fails) the subscription.

The frontend never talks to Shwary directly.

---

## Flow overview

```
┌──────────┐  1. POST /subscriptions   ┌─────────┐  2. Initiate payment  ┌────────┐
│ Frontend │ ────────────────────────▶ │ Backend │ ───────────────────▶  │ Shwary │
└──────────┘                           └─────────┘                       └────────┘
     ▲                                      ▲                                │
     │ 5. Poll subscription status          │ 4. POST /shwary-web-hook       │
     │ (or subscribe to Supabase realtime)  └────────────────────────────────┘
     │                                              (final status)
     │
     │ 3. 202 Accepted — { subscription: pending, transaction: pending }
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

### Option B — Polling

If realtime is not wired up, poll the subscription row directly from Supabase (or via your own backend read endpoint):

```ts
async function pollSubscription(id: string, { intervalMs = 3000, timeoutMs = 120_000 } = {}) {
  const deadline = Date.now() + timeoutMs;
  while (Date.now() < deadline) {
    const { data, error } = await supabase
      .from("subscriptions")
      .select("*")
      .eq("id", id)
      .single();
    if (error) throw error;
    if (data.status !== "pending") return data;
    await new Promise((r) => setTimeout(r, intervalMs));
  }
  throw new Error("subscription_timeout");
}
```

Use a 2–4 second interval and a 2-minute hard timeout. Shwary mobile-money confirmations typically arrive in under 60 seconds.

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
```

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
