# React Frontend Integration Guide

## Base URL

```
http://localhost:3000
```

---

## Response Envelope

Every endpoint returns the same shape:

```ts
type ApiResponse<T = unknown> = {
  status: "success" | "error";
  message: string;
  data?: T;
};
```

Paginated list responses nest inside `data`:

```ts
type Paginated<T> = {
  items: T[];
  total: number;
  page: number;
  limit: number;
};
```

---

## Authentication

### JWT — Dashboard routes

Login returns a JWT. Send it as a `Bearer` token on every `/api/*` request.

```ts
// lib/api.ts
import axios from "axios";

const api = axios.create({ baseURL: "http://localhost:3000" });

api.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});

api.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem("token");
      window.location.href = "/login";
    }
    return Promise.reject(err);
  }
);

export default api;
```

### API Key — POS routes

POS routes use the `X-API-Key` header instead of JWT. Keep the key in an env var.

```ts
const posApi = axios.create({
  baseURL: "http://localhost:3000",
  headers: { "X-API-Key": import.meta.env.VITE_POS_API_KEY },
});
```

---

## Error Handling

```ts
async function safeFetch<T>(fn: () => Promise<{ data: ApiResponse<T> }>) {
  try {
    const { data } = await fn();
    return { data: data.data, error: null };
  } catch (err: any) {
    const message = err.response?.data?.message ?? "Something went wrong";
    return { data: null, error: message };
  }
}
```

---

## Auth Endpoints

### POST `/auth/register`

```ts
type RegisterRequest = { email: string; password: string };
// Response: ApiResponse<null>  (201)

await api.post("/auth/register", { email, password });
```

### POST `/auth/login`

```ts
type LoginRequest  = { email: string; password: string };
type LoginData     = { token: string };

const { data } = await api.post<ApiResponse<LoginData>>("/auth/login", { email, password });
localStorage.setItem("token", data.data!.token);
```

### GET `/api/me`

```ts
type MeData = { user_id: number; email: string };

const { data } = await api.get<ApiResponse<MeData>>("/api/me");
```

---

## Suppliers

```ts
type Supplier = {
  id: string;
  name: string;
  contact_name: string;
  phone: string;
  email: string;
  address: string;
  created_at: string;
  updated_at: string;
};
```

| Method | Path | Notes |
|--------|------|-------|
| GET | `/api/v1/suppliers` | `?page=1&limit=20&search=` |
| POST | `/api/v1/suppliers` | body: `{ name*, contact_name, phone, email, address }` |
| GET | `/api/v1/suppliers/:id` | |
| PUT | `/api/v1/suppliers/:id` | same body as POST |
| DELETE | `/api/v1/suppliers/:id` | soft-delete |

```ts
// List
const { data } = await api.get<ApiResponse<Paginated<Supplier>>>("/api/v1/suppliers", {
  params: { page, limit, search },
});

// Create
await api.post("/api/v1/suppliers", { name, contact_name, phone, email, address });

// Update
await api.put(`/api/v1/suppliers/${id}`, { name, contact_name, phone, email, address });

// Delete
await api.delete(`/api/v1/suppliers/${id}`);
```

---

## Inventory Items

```ts
type InventoryItem = {
  id: string;
  sku: string;
  name: string;
  description: string;
  unit: string;
  quantity: number;
  min_quantity: number;
  cost_per_unit: number;
  created_at: string;
  updated_at: string;
};
```

| Method | Path | Notes |
|--------|------|-------|
| GET | `/api/v1/inventory/items` | `?page=1&limit=20&search=` |
| POST | `/api/v1/inventory/items` | body below |
| GET | `/api/v1/inventory/items/:id` | |
| PUT | `/api/v1/inventory/items/:id` | same body as POST |
| DELETE | `/api/v1/inventory/items/:id` | soft-delete |
| POST | `/api/v1/inventory/items/:id/adjust` | adjust stock |
| GET | `/api/v1/inventory/items/:id/transactions` | `?page&limit` |

```ts
// Create / Update
const body = {
  sku: "RAW-001",       // required
  name: "Sugar 1kg",   // required
  unit: "kg",          // required
  description: "",
  min_quantity: 10,
  cost_per_unit: 5.5,
};

// Stock adjustment
await api.post(`/api/v1/inventory/items/${id}/adjust`, {
  quantity: 50,    // must be > 0
  is_add: true,   // true = add, false = deduct
  note: "restock",
});

// Transaction history
const { data } = await api.get<ApiResponse<Paginated<StockTransaction>>>(
  `/api/v1/inventory/items/${id}/transactions`,
  { params: { page, limit } }
);
```

```ts
type StockTransaction = {
  id: string;
  inventory_item_id: string;
  type: string;         // "MANUAL_ADD" | "MANUAL_DEDUCT" | "PO_RECEIVE" | "POS_SALE"
  quantity: number;
  note: string;
  created_at: string;
};
```

---

## Products

```ts
type Product = {
  id: string;
  pos_product_id: string;
  name: string;
  sku: string;
  barcode: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
};

type BOMItem = {
  inventory_item_id: string;
  quantity: number;
  inventory_item?: InventoryItem;
};
```

| Method | Path | Notes |
|--------|------|-------|
| GET | `/api/v1/products` | `?page=1&limit=20&search=` |
| POST | `/api/v1/products` | body below |
| GET | `/api/v1/products/barcode/:barcode` | lookup by barcode scan |
| GET | `/api/v1/products/:id` | |
| PUT | `/api/v1/products/:id` | same body as POST |
| DELETE | `/api/v1/products/:id` | soft-delete |
| GET | `/api/v1/products/:id/bom` | returns `BOMItem[]` |
| PUT | `/api/v1/products/:id/bom` | full replace of BOM |

```ts
// Create
await api.post("/api/v1/products", {
  pos_product_id: "POS-SKU-001",  // required — must match your POS system
  name: "Latte",                  // required
  sku: "LATTE-001",
  barcode: "1234567890128",       // optional — used for barcode scanning
  is_active: true,
});

// Lookup by barcode (e.g. after scanning a product)
const { data } = await api.get<ApiResponse<Product>>(
  `/api/v1/products/barcode/1234567890128`
);

// Update BOM (full replace)
await api.put(`/api/v1/products/${id}/bom`, {
  items: [
    { inventory_item_id: "uuid-of-milk", quantity: 0.2 },
    { inventory_item_id: "uuid-of-coffee", quantity: 0.018 },
  ],
});
```

---

## Purchase Orders

```ts
type PurchaseOrder = {
  id: string;
  po_number: string;
  supplier_id: string;
  status: "DRAFT" | "ORDERED" | "PARTIALLY_RECEIVED" | "RECEIVED" | "CANCELLED";
  ordered_at: string | null;
  expected_at: string | null;
  received_at: string | null;
  notes: string;
  created_by: number | null;
  supplier?: Supplier;
  items?: PurchaseOrderItem[];
  created_at: string;
  updated_at: string;
};

type PurchaseOrderItem = {
  id: string;
  purchase_order_id: string;
  inventory_item_id: string;
  quantity_ordered: number;
  quantity_received: number;
  cost_per_unit: number;
  inventory_item?: InventoryItem;
};
```

| Method | Path | Notes |
|--------|------|-------|
| GET | `/api/v1/purchase-orders` | `?page&limit&status=DRAFT` |
| POST | `/api/v1/purchase-orders` | create |
| GET | `/api/v1/purchase-orders/:id` | includes items + supplier |
| PUT | `/api/v1/purchase-orders/:id` | update notes / status / expected_at |
| POST | `/api/v1/purchase-orders/:id/receive` | mark goods received |
| POST | `/api/v1/purchase-orders/:id/cancel` | cancel |

```ts
// Create PO
await api.post("/api/v1/purchase-orders", {
  supplier_id: "uuid",
  notes: "Urgent restock",
  expected_at: "2026-07-05T00:00:00Z",  // RFC3339, optional
  items: [
    {
      inventory_item_id: "uuid",
      quantity_ordered: 100,
      cost_per_unit: 4.5,
    },
  ],
});

// Receive goods (partial or full)
await api.post(`/api/v1/purchase-orders/${id}/receive`, {
  note: "Received at warehouse",
  items: [
    {
      purchase_order_item_id: "item-uuid",
      quantity_received: 100,
    },
  ],
});

// Update (e.g. change status to ORDERED)
await api.put(`/api/v1/purchase-orders/${id}`, {
  status: "ORDERED",
  notes: "Confirmed with supplier",
  expected_at: "2026-07-10T00:00:00Z",
});
```

---

## Webhooks

```ts
type Webhook = {
  id: string;
  name: string;
  url: string;
  events: string[];   // "STOCK_UPDATED" | "STOCK_LOW" | "STOCK_OUT"
  is_active: boolean;
};

type WebhookLog = {
  id: string;
  webhook_subscription_id: string;
  event: string;
  status_code: number;
  success: boolean;
  attempted_at: string;
};
```

| Method | Path | Notes |
|--------|------|-------|
| GET | `/api/v1/webhooks` | list all |
| POST | `/api/v1/webhooks` | create |
| GET | `/api/v1/webhooks/:id` | |
| PUT | `/api/v1/webhooks/:id` | update |
| DELETE | `/api/v1/webhooks/:id` | |
| POST | `/api/v1/webhooks/:id/test` | send test event |
| GET | `/api/v1/webhooks/:id/logs` | delivery history |

```ts
// Create
await api.post("/api/v1/webhooks", {
  name: "Stock Alerts",
  url: "https://yourapp.com/hooks/inventory",
  secret: "signing-secret",
  events: ["STOCK_LOW", "STOCK_OUT"],   // default: all three
  is_active: true,
});
```

---

## POS Endpoints (API Key auth)

These are called by your POS system, not the dashboard. Use `posApi` (see above).

### POST `/api/v1/pos/stock/deduct`

Idempotent — duplicate `pos_order_id` values are silently ignored.

```ts
await posApi.post("/api/v1/pos/stock/deduct", {
  pos_order_id: "ORDER-20260628-001",  // required, used for idempotency
  items: [
    { pos_product_id: "POS-SKU-001", quantity: 2 },
  ],
});
```

### GET `/api/v1/pos/stock/levels`

Returns current quantity for every inventory item.

```ts
const { data } = await posApi.get("/api/v1/pos/stock/levels");
```

### GET `/api/v1/pos/products/:pos_product_id/availability`

```ts
const { data } = await posApi.get(
  `/api/v1/pos/products/POS-SKU-001/availability`,
  { params: { quantity: 3 } }
);
// data.data: { available: boolean, ... }
```

---

## React Query Example

```tsx
// hooks/useInventory.ts
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import api from "@/lib/api";

export function useInventoryItems(page = 1, search = "") {
  return useQuery({
    queryKey: ["inventory", page, search],
    queryFn: async () => {
      const { data } = await api.get("/api/v1/inventory/items", {
        params: { page, limit: 20, search },
      });
      return data.data; // { items, total, page, limit }
    },
  });
}

export function useAdjustStock() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, quantity, isAdd, note }: {
      id: string; quantity: number; isAdd: boolean; note: string;
    }) =>
      api.post(`/api/v1/inventory/items/${id}/adjust`, {
        quantity,
        is_add: isAdd,
        note,
      }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["inventory"] }),
  });
}
```

---

## Environment Variables (Vite)

```env
# .env.local
VITE_API_BASE_URL=http://localhost:3000
VITE_POS_API_KEY=your-pos-api-key-here
```
