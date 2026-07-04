# Frontend Integration Guide

A practical guide for connecting a React + TypeScript frontend to this Inventory Management API.

**Base URL:** `http://localhost:3000`  
**Auth:** JWT Bearer for dashboard routes · API Key for POS routes

---

## Table of Contents

1. [Setup](#1-setup)
2. [TypeScript Types](#2-typescript-types)
3. [Language (i18n)](#3-language-i18n)
4. [Authentication](#4-authentication)
5. [Hooks — Suppliers](#5-hooks--suppliers)
6. [Hooks — Inventory Items](#6-hooks--inventory-items)
7. [Hooks — Products & Barcode Scan](#7-hooks--products--barcode-scan)
8. [Hooks — Purchase Orders](#8-hooks--purchase-orders)
9. [Hooks — Webhooks](#9-hooks--webhooks)
10. [POS Hooks (API Key)](#10-pos-hooks-api-key)
11. [Error Handling](#11-error-handling)

---

## 1. Setup

### Install dependencies

```bash
npm install axios @tanstack/react-query
```

### Environment variables

```env
# .env.local
VITE_API_BASE_URL=http://localhost:3000
VITE_POS_API_KEY=your-pos-api-key-here
```

### Axios instances

```ts
// src/lib/api.ts
import axios from "axios";

const lang = localStorage.getItem("lang") ?? "en"; // "th" or "en"

// Dashboard API — JWT auth
export const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  headers: { "Accept-Language": lang },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) config.headers.Authorization = `Bearer ${token}`;
  // Reflect language changes at runtime
  config.headers["Accept-Language"] = localStorage.getItem("lang") ?? "en";
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

// POS API — API Key auth
export const posApi = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  headers: {
    "X-API-Key": import.meta.env.VITE_POS_API_KEY,
    "Accept-Language": lang,
  },
});
```

### React Query provider

```tsx
// src/main.tsx
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

const queryClient = new QueryClient({
  defaultOptions: { queries: { staleTime: 30_000, retry: 1 } },
});

ReactDOM.createRoot(document.getElementById("root")!).render(
  <QueryClientProvider client={queryClient}>
    <App />
  </QueryClientProvider>
);
```

---

## 2. TypeScript Types

```ts
// src/types/index.ts

export type ApiResponse<T = unknown> = {
  status: "success" | "error";
  message: string;
  data?: T;
};

export type Paginated<T> = {
  items: T[];
  total: number;
  page: number;
  limit: number;
};

export type Supplier = {
  id: string;
  name: string;
  contact_name: string;
  phone: string;
  email: string;
  address: string;
  created_at: string;
  updated_at: string;
};

// Units are admin-managed (CRUD via /api/v1/inventory/units), not a fixed
// union — fetch the current list with useInventoryUnits() rather than hardcoding it.
export type InventoryUnit = {
  id: string;
  code: string;
  name: string;
  created_at: string;
  updated_at: string;
};

export type InventoryItem = {
  id: string;
  sku: string;
  name: string;
  description: string;
  unit: string; // an InventoryUnit.code
  quantity_in_stock: number;
  min_quantity: number;
  cost_per_unit: number;
  created_at: string;
  updated_at: string;
};

export type StockTransaction = {
  id: string;
  inventory_item_id: string;
  transaction_type: "IN" | "OUT" | "ADJUSTMENT_ADD" | "ADJUSTMENT_REMOVE";
  quantity: number;
  quantity_before: number;
  quantity_after: number;
  reference_type: string;
  reference_id: string;
  note: string;
  created_at: string;
};

export type BOMItem = {
  id: string;
  product_id: string;
  inventory_item_id: string;
  quantity_required: number;
  inventory_item?: InventoryItem;
};

export type Product = {
  id: string;
  pos_product_id: string;
  name: string;
  sku: string;
  barcode: string;
  is_active: boolean;
  bom?: BOMItem[];
  created_at: string;
  updated_at: string;
};

export type PurchaseOrderItem = {
  id: string;
  purchase_order_id: string;
  inventory_item_id: string;
  quantity_ordered: number;
  quantity_received: number;
  cost_per_unit: number;
  inventory_item?: InventoryItem;
};

export type PurchaseOrder = {
  id: string;
  po_number: string;
  supplier_id: string;
  status: "DRAFT" | "ORDERED" | "PARTIALLY_RECEIVED" | "RECEIVED" | "CANCELLED";
  ordered_at: string | null;
  expected_at: string | null;
  received_at: string | null;
  notes: string;
  supplier?: Supplier;
  items?: PurchaseOrderItem[];
  created_at: string;
  updated_at: string;
};

export type Webhook = {
  id: string;
  name: string;
  url: string;
  events: ("STOCK_UPDATED" | "STOCK_LOW" | "STOCK_OUT")[];
  is_active: boolean;
  created_at: string;
  updated_at: string;
};

export type WebhookLog = {
  id: string;
  webhook_subscription_id: string;
  event: string;
  status_code: number;
  success: boolean;
  attempted_at: string;
};
```

---

## 3. Language (i18n)

The API returns `message` fields in **Thai (th)** or **English (en)** based on the `Accept-Language` header. Supported languages:

| Code | Language |
|------|----------|
| `en` | English (default) |
| `th` | ภาษาไทย |

### Language context + hook

```ts
// src/lib/lang.ts
export type Lang = "en" | "th";

export function getLang(): Lang {
  return (localStorage.getItem("lang") as Lang) ?? "en";
}

export function setLang(lang: Lang) {
  localStorage.setItem("lang", lang);
  // Reload so the axios instance picks up the new header
  window.location.reload();
}
```

```tsx
// src/components/LanguageSwitcher.tsx
import { getLang, setLang } from "@/lib/lang";

export function LanguageSwitcher() {
  const current = getLang();
  return (
    <div>
      <button
        onClick={() => setLang("th")}
        style={{ fontWeight: current === "th" ? "bold" : "normal" }}
      >
        ไทย
      </button>
      <button
        onClick={() => setLang("en")}
        style={{ fontWeight: current === "en" ? "bold" : "normal" }}
      >
        EN
      </button>
    </div>
  );
}
```

> **Note:** Only the `message` field in the response envelope is translated. `data` fields (product names, SKUs, etc.) are stored as-is in the database.

---

## 4. Authentication

```ts
// src/hooks/useAuth.ts
import { useMutation } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { ApiResponse } from "@/types";

export function useLogin() {
  return useMutation({
    mutationFn: async (body: { email: string; password: string }) => {
      const { data } = await api.post<ApiResponse<{ token: string }>>(
        "/auth/login",
        body
      );
      return data.data!.token;
    },
    onSuccess: (token) => {
      localStorage.setItem("token", token);
    },
  });
}

export function useRegister() {
  return useMutation({
    mutationFn: (body: { email: string; password: string }) =>
      api.post<ApiResponse>("/auth/register", body),
  });
}

export function useMe() {
  return useQuery({
    queryKey: ["me"],
    queryFn: async () => {
      const { data } = await api.get<
        ApiResponse<{ user_id: number; email: string; role: "admin" | "staff" }>
      >("/api/me");
      return data.data!;
    },
  });
}

// Gate admin-only UI (e.g. the unit management screen) on this.
// The API also enforces it server-side — this is UX only, not a security boundary.
export function useIsAdmin() {
  const { data } = useMe();
  return data?.role === "admin";
}

export function logout() {
  localStorage.removeItem("token");
  window.location.href = "/login";
}
```

---

## 4. Hooks — Suppliers

```ts
// src/hooks/useSuppliers.ts
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { ApiResponse, Paginated, Supplier } from "@/types";

const KEYS = {
  list: (page: number, search: string) => ["suppliers", page, search],
  detail: (id: string) => ["suppliers", id],
};

export function useSuppliers(page = 1, search = "") {
  return useQuery({
    queryKey: KEYS.list(page, search),
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<Paginated<Supplier>>>(
        "/api/v1/suppliers",
        { params: { page, limit: 20, search } }
      );
      return data.data!;
    },
  });
}

export function useSupplier(id: string) {
  return useQuery({
    queryKey: KEYS.detail(id),
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<Supplier>>(
        `/api/v1/suppliers/${id}`
      );
      return data.data!;
    },
    enabled: !!id,
  });
}

type SupplierBody = {
  name: string;
  contact_name?: string;
  phone?: string;
  email?: string;
  address?: string;
};

export function useCreateSupplier() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: SupplierBody) =>
      api.post<ApiResponse<Supplier>>("/api/v1/suppliers", body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["suppliers"] }),
  });
}

export function useUpdateSupplier(id: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: SupplierBody) =>
      api.put<ApiResponse<Supplier>>(`/api/v1/suppliers/${id}`, body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["suppliers"] }),
  });
}

export function useDeleteSupplier() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.delete(`/api/v1/suppliers/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["suppliers"] }),
  });
}
```

---

## 5. Hooks — Inventory Items

```ts
// src/hooks/useInventory.ts
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { ApiResponse, Paginated, InventoryItem, InventoryUnit, StockTransaction } from "@/types";

const KEYS = {
  list: (page: number, search: string) => ["inventory", page, search],
  detail: (id: string) => ["inventory", id],
  transactions: (id: string, page: number) => ["inventory", id, "transactions", page],
  units: ["inventory", "units"],
};

// Units for the item form's <select> — GET /api/v1/inventory/units.
// Admin-managed but changes rarely, so cache it for the session.
export function useInventoryUnits() {
  return useQuery({
    queryKey: KEYS.units,
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<InventoryUnit[]>>("/api/v1/inventory/units");
      return data.data!;
    },
    staleTime: 5 * 60 * 1000,
  });
}

// Admin-only — server returns 403 for non-admins (see useIsAdmin in useAuth.ts).
export function useCreateInventoryUnit() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: { code: string; name: string }) =>
      api.post<ApiResponse<InventoryUnit>>("/api/v1/inventory/units", body),
    onSuccess: () => qc.invalidateQueries({ queryKey: KEYS.units }),
  });
}

export function useUpdateInventoryUnit(id: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: { code: string; name: string }) =>
      api.put<ApiResponse<InventoryUnit>>(`/api/v1/inventory/units/${id}`, body),
    onSuccess: () => qc.invalidateQueries({ queryKey: KEYS.units }),
  });
}

// Fails with 400 if the unit is still referenced by an inventory item.
export function useDeleteInventoryUnit() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.delete(`/api/v1/inventory/units/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: KEYS.units }),
  });
}

export function useInventoryItems(page = 1, search = "") {
  return useQuery({
    queryKey: KEYS.list(page, search),
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<Paginated<InventoryItem>>>(
        "/api/v1/inventory/items",
        { params: { page, limit: 20, search } }
      );
      return data.data!;
    },
  });
}

export function useInventoryItem(id: string) {
  return useQuery({
    queryKey: KEYS.detail(id),
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<InventoryItem>>(
        `/api/v1/inventory/items/${id}`
      );
      return data.data!;
    },
    enabled: !!id,
  });
}

export function useInventoryTransactions(id: string, page = 1) {
  return useQuery({
    queryKey: KEYS.transactions(id, page),
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<Paginated<StockTransaction>>>(
        `/api/v1/inventory/items/${id}/transactions`,
        { params: { page, limit: 20 } }
      );
      return data.data!;
    },
    enabled: !!id,
  });
}

type InventoryBody = {
  sku: string;
  name: string;
  unit: string;
  description?: string;
  min_quantity?: number;
  cost_per_unit?: number;
};

export function useCreateInventoryItem() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: InventoryBody) =>
      api.post<ApiResponse<InventoryItem>>("/api/v1/inventory/items", body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["inventory"] }),
  });
}

export function useUpdateInventoryItem(id: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: InventoryBody) =>
      api.put<ApiResponse<InventoryItem>>(`/api/v1/inventory/items/${id}`, body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["inventory"] }),
  });
}

export function useAdjustStock() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      id,
      quantity,
      is_add,
      note,
    }: {
      id: string;
      quantity: number;
      is_add: boolean;
      note?: string;
    }) =>
      api.post(`/api/v1/inventory/items/${id}/adjust`, { quantity, is_add, note }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["inventory"] }),
  });
}

export function useDeleteInventoryItem() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.delete(`/api/v1/inventory/items/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["inventory"] }),
  });
}
```

### Unit select in the item form

```tsx
function UnitSelect({ value, onChange }: { value: string; onChange: (v: string) => void }) {
  const { data: units, isLoading } = useInventoryUnits();

  return (
    <select value={value} onChange={(e) => onChange(e.target.value)} disabled={isLoading}>
      <option value="" disabled>Select unit</option>
      {units?.map((u) => (
        <option key={u.id} value={u.code}>{u.name} ({u.code})</option>
      ))}
    </select>
  );
}
```

`POST`/`PUT /api/v1/inventory/items` sends `unit` as the unit's `code` (e.g. `"kg"`, not the ID), and rejects any code not in the current list with a 400 (`err.unit_invalid`) — always populate the field from `useInventoryUnits()` rather than a free-text input.

### Admin unit management screen

Units are CRUD-managed by admins only — `POST`/`PUT`/`DELETE /api/v1/inventory/units` return 403 for non-admin users (checked via the `role` claim on the JWT, exposed at `/api/me`). Gate the screen client-side with `useIsAdmin()`, but remember the server enforces this independently:

```tsx
function UnitManagementPage() {
  const isAdmin = useIsAdmin();
  const { data: units } = useInventoryUnits();
  const createUnit = useCreateInventoryUnit();
  const deleteUnit = useDeleteInventoryUnit();

  if (!isAdmin) return <p>Admins only.</p>;

  return (
    <ul>
      {units?.map((u) => (
        <li key={u.id}>
          {u.name} ({u.code})
          <button onClick={() => deleteUnit.mutate(u.id)}>Delete</button>
        </li>
      ))}
    </ul>
  );
}
```

---

## 6. Hooks — Products & Barcode Scan

### Standard CRUD hooks

```ts
// src/hooks/useProducts.ts
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { ApiResponse, Paginated, Product, BOMItem } from "@/types";

const KEYS = {
  list: (page: number, search: string) => ["products", page, search],
  detail: (id: string) => ["products", id],
  barcode: (barcode: string) => ["products", "barcode", barcode],
  bom: (id: string) => ["products", id, "bom"],
};

export function useProducts(page = 1, search = "") {
  return useQuery({
    queryKey: KEYS.list(page, search),
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<Paginated<Product>>>(
        "/api/v1/products",
        { params: { page, limit: 20, search } }
      );
      return data.data!;
    },
  });
}

export function useProduct(id: string) {
  return useQuery({
    queryKey: KEYS.detail(id),
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<Product>>(
        `/api/v1/products/${id}`
      );
      return data.data!;
    },
    enabled: !!id,
  });
}

export function useProductBOM(id: string) {
  return useQuery({
    queryKey: KEYS.bom(id),
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<BOMItem[]>>(
        `/api/v1/products/${id}/bom`
      );
      return data.data!;
    },
    enabled: !!id,
  });
}

type ProductBody = {
  pos_product_id: string;
  name: string;
  sku?: string;
  barcode?: string;
  is_active?: boolean;
};

export function useCreateProduct() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: ProductBody) =>
      api.post<ApiResponse<Product>>("/api/v1/products", body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["products"] }),
  });
}

export function useUpdateProduct(id: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: ProductBody) =>
      api.put<ApiResponse<Product>>(`/api/v1/products/${id}`, body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["products"] }),
  });
}

export function useUpdateProductBOM(id: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (items: { inventory_item_id: string; quantity_required: number }[]) =>
      api.put<ApiResponse<BOMItem[]>>(`/api/v1/products/${id}/bom`, { items }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["products", id] }),
  });
}

export function useDeleteProduct() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.delete(`/api/v1/products/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["products"] }),
  });
}
```

### Barcode scan hook

```ts
// src/hooks/useProductByBarcode.ts
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { ApiResponse, Product } from "@/types";

/**
 * Looks up a product by its barcode. Pass an empty string or null to disable.
 *
 * Usage:
 *   const { data, isLoading, isError } = useProductByBarcode(scannedBarcode);
 */
export function useProductByBarcode(barcode: string | null) {
  return useQuery({
    queryKey: ["products", "barcode", barcode],
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<Product>>(
        `/api/v1/products/barcode/${barcode}`
      );
      return data.data!;
    },
    enabled: !!barcode,
    retry: false, // 404 should surface immediately, not retry
  });
}
```

### Barcode scanner component example

```tsx
// src/components/BarcodeScanner.tsx
import { useState } from "react";
import { useProductByBarcode } from "@/hooks/useProductByBarcode";

export function BarcodeScanner() {
  const [input, setInput] = useState("");
  const [barcode, setBarcode] = useState<string | null>(null);

  const { data: product, isLoading, isError } = useProductByBarcode(barcode);

  const handleScan = (e: React.FormEvent) => {
    e.preventDefault();
    setBarcode(input.trim() || null);
  };

  return (
    <div>
      <form onSubmit={handleScan}>
        <input
          autoFocus
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="Scan or type barcode…"
        />
        <button type="submit">Look up</button>
      </form>

      {isLoading && <p>Looking up…</p>}
      {isError && <p>Product not found for barcode: {barcode}</p>}

      {product && (
        <div>
          <h2>{product.name}</h2>
          <p>SKU: {product.sku}</p>
          <p>Barcode: {product.barcode}</p>
          <h3>Ingredients (BOM)</h3>
          <ul>
            {product.bom?.map((item) => (
              <li key={item.id}>
                {item.inventory_item?.name} — need {item.quantity_required}{" "}
                {item.inventory_item?.unit}, have{" "}
                {item.inventory_item?.quantity_in_stock}
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
}
```

> **Hardware scanner tip:** Most USB/Bluetooth barcode scanners act as a keyboard and append `Enter` after the code. Set `autoFocus` on the input and the form's `onSubmit` handles it automatically — no extra event listener needed.

---

## 7. Hooks — Purchase Orders

```ts
// src/hooks/usePurchaseOrders.ts
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { ApiResponse, Paginated, PurchaseOrder } from "@/types";

const KEYS = {
  list: (page: number, status?: string) => ["purchase-orders", page, status],
  detail: (id: string) => ["purchase-orders", id],
};

export function usePurchaseOrders(page = 1, status?: string) {
  return useQuery({
    queryKey: KEYS.list(page, status),
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<Paginated<PurchaseOrder>>>(
        "/api/v1/purchase-orders",
        { params: { page, limit: 20, status } }
      );
      return data.data!;
    },
  });
}

export function usePurchaseOrder(id: string) {
  return useQuery({
    queryKey: KEYS.detail(id),
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<PurchaseOrder>>(
        `/api/v1/purchase-orders/${id}`
      );
      return data.data!;
    },
    enabled: !!id,
  });
}

type CreatePOBody = {
  supplier_id: string;
  notes?: string;
  expected_at?: string;
  items: { inventory_item_id: string; quantity_ordered: number; cost_per_unit: number }[];
};

export function useCreatePurchaseOrder() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: CreatePOBody) =>
      api.post<ApiResponse<PurchaseOrder>>("/api/v1/purchase-orders", body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["purchase-orders"] }),
  });
}

export function useUpdatePurchaseOrder(id: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: { status?: string; notes?: string; expected_at?: string }) =>
      api.put<ApiResponse<PurchaseOrder>>(`/api/v1/purchase-orders/${id}`, body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["purchase-orders"] }),
  });
}

export function useReceivePurchaseOrder(id: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: {
      note?: string;
      items: { purchase_order_item_id: string; quantity_received: number }[];
    }) => api.post(`/api/v1/purchase-orders/${id}/receive`, body),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["purchase-orders"] });
      qc.invalidateQueries({ queryKey: ["inventory"] });
    },
  });
}

export function useCancelPurchaseOrder(id: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: () => api.post(`/api/v1/purchase-orders/${id}/cancel`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["purchase-orders"] }),
  });
}
```

---

## 8. Hooks — Webhooks

```ts
// src/hooks/useWebhooks.ts
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { ApiResponse, Webhook, WebhookLog } from "@/types";

export function useWebhooks() {
  return useQuery({
    queryKey: ["webhooks"],
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<Webhook[]>>("/api/v1/webhooks");
      return data.data!;
    },
  });
}

export function useWebhook(id: string) {
  return useQuery({
    queryKey: ["webhooks", id],
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<Webhook>>(`/api/v1/webhooks/${id}`);
      return data.data!;
    },
    enabled: !!id,
  });
}

export function useWebhookLogs(id: string) {
  return useQuery({
    queryKey: ["webhooks", id, "logs"],
    queryFn: async () => {
      const { data } = await api.get<ApiResponse<WebhookLog[]>>(
        `/api/v1/webhooks/${id}/logs`
      );
      return data.data!;
    },
    enabled: !!id,
  });
}

type WebhookBody = {
  name: string;
  url: string;
  secret?: string;
  events?: ("STOCK_UPDATED" | "STOCK_LOW" | "STOCK_OUT")[];
  is_active?: boolean;
};

export function useCreateWebhook() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: WebhookBody) =>
      api.post<ApiResponse<Webhook>>("/api/v1/webhooks", body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["webhooks"] }),
  });
}

export function useUpdateWebhook(id: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: Partial<WebhookBody>) =>
      api.put<ApiResponse<Webhook>>(`/api/v1/webhooks/${id}`, body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["webhooks"] }),
  });
}

export function useDeleteWebhook() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.delete(`/api/v1/webhooks/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["webhooks"] }),
  });
}

export function useTestWebhook(id: string) {
  return useMutation({
    mutationFn: () => api.post(`/api/v1/webhooks/${id}/test`),
  });
}
```

---

## 9. POS Hooks (API Key)

```ts
// src/hooks/usePos.ts
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { posApi } from "@/lib/api";
import type { ApiResponse } from "@/types";

type StockLevel = {
  inventory_item_id: string;
  sku: string;
  name: string;
  unit: string;
  quantity_in_stock: number;
  min_quantity: number;
  unit_cost: number;
  is_low: boolean;
  is_out: boolean;
};

type AvailabilityResult = {
  pos_product_id: string;
  name: string;
  is_available: boolean;
  details: {
    inventory_item_id: string;
    sku: string;
    name: string;
    required: number;
    available: number;
    unit_cost: number;
    is_sufficient: boolean;
  }[];
};

type StockDeduction = {
  inventory_item_id: string;
  sku: string;
  name: string;
  quantity_deducted: number;
  quantity_remaining: number;
};

type CostBreakdownItem = {
  pos_product_id: string;
  unit_cost: number;
  total_cost: number;
};

type DeductStockResult = {
  pos_order_id: string;
  status: "processed" | "already_processed";
  deductions: StockDeduction[];
  /** One entry per distinct pos_product_id, present on every "processed" response (unit_cost may be 0 if a product has no BOM cost data yet). Omitted on "already_processed" replays — the cost was already applied to the original order. */
  cost_breakdown?: CostBreakdownItem[];
};

export function useStockLevels() {
  return useQuery({
    queryKey: ["pos", "stock-levels"],
    queryFn: async () => {
      const { data } = await posApi.get<ApiResponse<StockLevel[]>>(
        "/api/v1/pos/stock/levels"
      );
      return data.data!;
    },
    refetchInterval: 5 * 60 * 1000, // sync every 5 minutes
  });
}

export function useProductAvailability(posProductId: string, quantity = 1) {
  return useQuery({
    queryKey: ["pos", "availability", posProductId, quantity],
    queryFn: async () => {
      const { data } = await posApi.get<ApiResponse<AvailabilityResult>>(
        `/api/v1/pos/products/${posProductId}/availability`,
        { params: { quantity } }
      );
      return data.data!;
    },
    enabled: !!posProductId,
  });
}

export function useDeductStock() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (body: {
      pos_order_id: string;
      items: { pos_product_id: string; quantity: number }[];
    }) => {
      const { data } = await posApi.post<ApiResponse<DeductStockResult>>(
        "/api/v1/pos/stock/deduct",
        body
      );
      return data.data!;
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ["pos", "stock-levels"] }),
  });
}
```

### Using `cost_breakdown` to compute order cost/profit

`cost_breakdown` is additive — keep sending orders exactly as before. When the
array is present and non-empty, trust it over any locally-stored `cost_price`
(a `0` `total_cost` for a given product means Inventory genuinely has no BOM
cost data for it yet — that's still authoritative, not a signal to fall
back). Only fall back to POS's own `cost_price` when the array itself is
missing or empty, e.g. on an `already_processed` replay.

```ts
const deductStock = useDeductStock();

const result = await deductStock.mutateAsync({
  pos_order_id: orderId,
  items: cart.map((i) => ({ pos_product_id: i.posProductId, quantity: i.qty })),
});

const costByProduct = new Map(
  result.cost_breakdown?.map((c) => [c.pos_product_id, c]) ?? []
);

const orderCost = cart.reduce((sum, item) => {
  const fromInventory = costByProduct.get(item.posProductId)?.total_cost;
  const fallback = item.costPrice * item.qty; // POS's own manual cost_price
  return sum + (fromInventory ?? fallback);
}, 0);

const profit = orderTotal - orderCost;
```

---

## 10. Error Handling

```ts
// src/lib/apiError.ts
import { AxiosError } from "axios";
import type { ApiResponse } from "@/types";

export function getApiError(err: unknown): string {
  if (err instanceof AxiosError) {
    return (
      (err.response?.data as ApiResponse)?.message ??
      err.message ??
      "Something went wrong"
    );
  }
  return "Something went wrong";
}
```

```tsx
// Usage in a component
const createProduct = useCreateProduct();

const handleSubmit = async (values: ProductBody) => {
  try {
    await createProduct.mutateAsync(values);
    toast.success("Product created");
  } catch (err) {
    toast.error(getApiError(err));
  }
};
```

### HTTP status reference

| Status | Meaning |
|--------|---------|
| `400` | Validation error or insufficient stock |
| `401` | Missing or expired token / wrong API key |
| `404` | Resource not found |
| `500` | Server error |
