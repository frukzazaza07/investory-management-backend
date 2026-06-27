# Inventory Management API - Endpoint Test Script
# Usage: .\scripts\test_endpoints.ps1

$BASE_URL = "http://localhost:3000"
$POS_API_KEY = "dev-pos-api-key-2026"
$UNIQUE = Get-Date -Format 'yyyyMMddHHmmss'
$PASS = 0
$FAIL = 0

function Pass($label) {
    Write-Host "  [PASS] $label" -ForegroundColor Green
    $script:PASS++
}

function Fail($label, $detail) {
    Write-Host "  [FAIL] $label" -ForegroundColor Red
    if ($detail) { Write-Host "         $detail" -ForegroundColor DarkRed }
    $script:FAIL++
}

function Section($title) {
    Write-Host "`n=== $title ===" -ForegroundColor Cyan
}

function Invoke($method, $url, $body, $headers) {
    try {
        $params = @{ Method = $method; Uri = $url; ContentType = "application/json"; ErrorAction = "Stop" }
        if ($body)    { $params.Body    = ($body | ConvertTo-Json -Depth 10) }
        if ($headers) { $params.Headers = $headers }
        return Invoke-RestMethod @params
    } catch {
        $statusCode = $_.Exception.Response.StatusCode.value__
        $errBody    = $null
        try { $errBody = $_.ErrorDetails.Message | ConvertFrom-Json } catch {}
        return [PSCustomObject]@{ _error = $_.Exception.Message; _status = $statusCode; _body = $errBody }
    }
}

# ─────────────────────────────
Section "AUTH"
# ─────────────────────────────

$resp = Invoke "POST" "$BASE_URL/auth/login" @{ email = "admin@example.com"; password = "admin123" }
if ($resp.data.token) {
    Pass "Login as admin"
    $TOKEN = $resp.data.token
} else {
    Fail "Login as admin" ($resp | ConvertTo-Json)
    exit 1
}

$AUTH       = @{ Authorization = "Bearer $TOKEN" }
$POS_HDRS   = @{ "X-API-Key" = $POS_API_KEY }

# ─────────────────────────────
Section "SUPPLIERS"
# ─────────────────────────────

$resp = Invoke "GET" "$BASE_URL/api/v1/suppliers" $null $AUTH
if ($resp.status -eq "success") { Pass "List suppliers" } else { Fail "List suppliers" }

$resp = Invoke "POST" "$BASE_URL/api/v1/suppliers" @{
    name = "Test Supplier $UNIQUE"; contact_name = "Tester"; phone = "0800000001"
    email = "test$UNIQUE@supplier.com"; address = "1 Test Rd"
} $AUTH
if ($resp.data.id) {
    Pass "Create supplier"
    $SUPPLIER_ID = $resp.data.id
} else { Fail "Create supplier" ($resp.data | ConvertTo-Json) }

$resp = Invoke "GET" "$BASE_URL/api/v1/suppliers/$SUPPLIER_ID" $null $AUTH
if ($resp.data.id -eq $SUPPLIER_ID) { Pass "Get supplier by ID" } else { Fail "Get supplier" }

$resp = Invoke "PUT" "$BASE_URL/api/v1/suppliers/$SUPPLIER_ID" @{
    name = "Updated Supplier $UNIQUE"; contact_name = "Updated"; phone = "0800000002"
    email = "upd$UNIQUE@supplier.com"; address = "2 Updated Rd"
} $AUTH
if ($resp.data.name -like "Updated Supplier*") { Pass "Update supplier" } else { Fail "Update supplier" ($resp | ConvertTo-Json) }

# ─────────────────────────────
Section "INVENTORY ITEMS"
# ─────────────────────────────

$resp = Invoke "GET" "$BASE_URL/api/v1/inventory/items" $null $AUTH
if ($resp.status -eq "success") { Pass "List inventory items" } else { Fail "List inventory items" }

$resp = Invoke "POST" "$BASE_URL/api/v1/inventory/items" @{
    sku = "TEST-ITEM-$UNIQUE"; name = "Test Item $UNIQUE"
    description = "A test item"; unit = "piece"; min_quantity = 5; cost_per_unit = 10.5
} $AUTH
if ($resp.data.id) {
    Pass "Create inventory item"
    $ITEM_ID = $resp.data.id
} else { Fail "Create inventory item" ($resp._body | ConvertTo-Json) }

$resp = Invoke "GET" "$BASE_URL/api/v1/inventory/items/$ITEM_ID" $null $AUTH
if ($resp.data.id -eq $ITEM_ID) { Pass "Get inventory item by ID" } else { Fail "Get inventory item" }

$resp = Invoke "PUT" "$BASE_URL/api/v1/inventory/items/$ITEM_ID" @{
    sku = "TEST-ITEM-$UNIQUE"; name = "Test Item Updated"
    description = "Updated"; unit = "piece"; min_quantity = 10; cost_per_unit = 12.0
} $AUTH
if ($resp.data.name -eq "Test Item Updated") { Pass "Update inventory item" } else { Fail "Update inventory item" ($resp | ConvertTo-Json) }

$resp = Invoke "POST" "$BASE_URL/api/v1/inventory/items/$ITEM_ID/adjust" @{
    quantity = 100; is_add = $true; note = "Initial stock in"
} $AUTH
if ($resp.data.quantity_in_stock -eq 100) { Pass "Adjust stock (add 100)" } else { Fail "Adjust stock add" ($resp | ConvertTo-Json) }

$resp = Invoke "POST" "$BASE_URL/api/v1/inventory/items/$ITEM_ID/adjust" @{
    quantity = 20; is_add = $false; note = "Remove 20"
} $AUTH
if ($resp.data.quantity_in_stock -eq 80) { Pass "Adjust stock (remove 20)" } else { Fail "Adjust stock remove" ($resp | ConvertTo-Json) }

$resp = Invoke "GET" "$BASE_URL/api/v1/inventory/items/$ITEM_ID/transactions" $null $AUTH
if ($resp.data.total -ge 2) { Pass "Get stock transactions (2+ records)" } else { Fail "Get stock transactions" ($resp | ConvertTo-Json) }

$resp = Invoke "GET" "$BASE_URL/api/v1/inventory/items?search=Coffee" $null $AUTH
if ($resp.data.total -ge 1) { Pass "Search inventory items" } else { Fail "Search inventory items" }

# ─────────────────────────────
Section "PRODUCTS"
# ─────────────────────────────

$resp = Invoke "GET" "$BASE_URL/api/v1/products" $null $AUTH
if ($resp.data.total -ge 3) { Pass "List products (3+ seeded)" } else { Fail "List products" ($resp.data.total) }

$resp = Invoke "POST" "$BASE_URL/api/v1/products" @{
    pos_product_id = "pos-test-drink-$UNIQUE"; name = "Test Drink $UNIQUE"
    sku = "BEV-TEST-$UNIQUE"; is_active = $true
} $AUTH
if ($resp.data.id) {
    Pass "Create product"
    $PRODUCT_ID = $resp.data.id
} else { Fail "Create product" ($resp._body | ConvertTo-Json) }

$resp = Invoke "GET" "$BASE_URL/api/v1/products/$PRODUCT_ID" $null $AUTH
if ($resp.data.id -eq $PRODUCT_ID) { Pass "Get product by ID" } else { Fail "Get product" }

$resp = Invoke "PUT" "$BASE_URL/api/v1/products/$PRODUCT_ID" @{
    pos_product_id = "pos-test-drink-$UNIQUE"; name = "Test Drink Updated"
    sku = "BEV-TEST-$UNIQUE"; is_active = $true
} $AUTH
if ($resp.data.name -eq "Test Drink Updated") { Pass "Update product" } else { Fail "Update product" ($resp | ConvertTo-Json) }

$invResp   = Invoke "GET" "$BASE_URL/api/v1/inventory/items?search=Coffee" $null $AUTH
$COFFEE_ID = $invResp.data.items[0].id

$resp = Invoke "PUT" "$BASE_URL/api/v1/products/$PRODUCT_ID/bom" @{
    items = @(
        @{ inventory_item_id = $COFFEE_ID; quantity_required = 15 },
        @{ inventory_item_id = $ITEM_ID;   quantity_required = 1  }
    )
} $AUTH
if ($resp.data.Count -eq 2) { Pass "Update product BOM (2 items)" } else { Fail "Update product BOM" ($resp | ConvertTo-Json) }

$resp = Invoke "GET" "$BASE_URL/api/v1/products/$PRODUCT_ID/bom" $null $AUTH
if ($resp.data.Count -eq 2) { Pass "Get product BOM" } else { Fail "Get product BOM" ($resp | ConvertTo-Json) }

$resp = Invoke "GET" "$BASE_URL/api/v1/products/prod-americano-001/bom" $null $AUTH
if ($resp.data.Count -ge 3) { Pass "Seeded Americano BOM (3+ items)" } else { Fail "Seeded Americano BOM" ($resp | ConvertTo-Json) }

# ─────────────────────────────
Section "PURCHASE ORDERS"
# ─────────────────────────────

$resp = Invoke "GET" "$BASE_URL/api/v1/purchase-orders" $null $AUTH
if ($resp.status -eq "success") { Pass "List purchase orders" } else { Fail "List purchase orders" }

# Use seeded supplier & inventory items for the PO
$COFFEE_ITEM_ID = $COFFEE_ID

$resp = Invoke "POST" "$BASE_URL/api/v1/purchase-orders" @{
    supplier_id = "sup-coffee-world-0001"
    notes       = "Test PO $UNIQUE"
    items       = @(
        @{ inventory_item_id = "inv-coffee-beans-001"; quantity_ordered = 2000; cost_per_unit = 0.45 },
        @{ inventory_item_id = "inv-milk-001";         quantity_ordered = 5000; cost_per_unit = 0.035 }
    )
} $AUTH
if ($resp.data.id) {
    Pass "Create purchase order"
    $PO_ID       = $resp.data.id
    $PO_ITEM_ID  = $resp.data.items[0].id
    $PO_ITEM2_ID = $resp.data.items[1].id
} else { Fail "Create purchase order" ($resp._body | ConvertTo-Json) }

$resp = Invoke "GET" "$BASE_URL/api/v1/purchase-orders/$PO_ID" $null $AUTH
if ($resp.data.status -eq "DRAFT") { Pass "Get purchase order (DRAFT status)" } else { Fail "Get purchase order" ($resp.data.status) }

$resp = Invoke "PUT" "$BASE_URL/api/v1/purchase-orders/$PO_ID" @{
    notes = "Updated notes $UNIQUE"; status = "ORDERED"
} $AUTH
if ($resp.data.status -eq "ORDERED") { Pass "Update PO status to ORDERED" } else { Fail "Update purchase order" ($resp | ConvertTo-Json) }

$stockBefore = (Invoke "GET" "$BASE_URL/api/v1/inventory/items/inv-coffee-beans-001" $null $AUTH).data.quantity_in_stock

$resp = Invoke "POST" "$BASE_URL/api/v1/purchase-orders/$PO_ID/receive" @{
    note  = "Full shipment"
    items = @(
        @{ purchase_order_item_id = $PO_ITEM_ID;  quantity_received = 2000 },
        @{ purchase_order_item_id = $PO_ITEM2_ID; quantity_received = 5000 }
    )
} $AUTH
if ($resp.data.status -eq "RECEIVED") { Pass "Receive PO (fully received)" } else { Fail "Receive PO" ($resp | ConvertTo-Json) }

$stockAfter = (Invoke "GET" "$BASE_URL/api/v1/inventory/items/inv-coffee-beans-001" $null $AUTH).data.quantity_in_stock
if ($stockAfter -eq ($stockBefore + 2000)) {
    Pass "Stock increased by 2000 after PO receive"
} else {
    Fail "Stock increase after receive" "Before=$stockBefore After=$stockAfter expected=$($stockBefore+2000)"
}

# Cancel test
$resp2 = Invoke "POST" "$BASE_URL/api/v1/purchase-orders" @{
    supplier_id = "sup-coffee-world-0001"
    items = @(@{ inventory_item_id = "inv-coffee-beans-001"; quantity_ordered = 100; cost_per_unit = 0.5 })
} $AUTH
$PO_CANCEL_ID = $resp2.data.id

$resp = Invoke "POST" "$BASE_URL/api/v1/purchase-orders/$PO_CANCEL_ID/cancel" $null $AUTH
if ($resp.data.status -eq "CANCELLED") { Pass "Cancel purchase order" } else { Fail "Cancel purchase order" ($resp | ConvertTo-Json) }

$resp = Invoke "GET" "$BASE_URL/api/v1/purchase-orders?status=RECEIVED" $null $AUTH
if ($resp.data.total -ge 1) { Pass "Filter POs by RECEIVED status" } else { Fail "Filter POs by status" }

# ─────────────────────────────
Section "POS INTEGRATION"
# ─────────────────────────────

# No API key → 401
$resp = Invoke "GET" "$BASE_URL/api/v1/pos/stock/levels" $null $null
if ($resp._status -eq 401) { Pass "Reject request without API key (401)" } else { Fail "API key protection" }

# Stock levels
$resp = Invoke "GET" "$BASE_URL/api/v1/pos/stock/levels" $null $POS_HDRS
if ($resp.status -eq "success" -and $resp.data.Count -ge 5) {
    Pass "Get stock levels (5+ items)"
} else { Fail "Get stock levels" ($resp | ConvertTo-Json) }

# Availability
$resp = Invoke "GET" "$BASE_URL/api/v1/pos/products/pos-latte/availability?quantity=2" $null $POS_HDRS
if ($resp.data.is_available -eq $true) { Pass "Check Latte availability (qty=2) - available" } else { Fail "Check Latte availability" ($resp | ConvertTo-Json) }

# Deduct stock
$coffeeBefore = (Invoke "GET" "$BASE_URL/api/v1/inventory/items/inv-coffee-beans-001" $null $AUTH).data.quantity_in_stock

$ORDER_ID = "POS-ORDER-$UNIQUE"
$resp = Invoke "POST" "$BASE_URL/api/v1/pos/stock/deduct" @{
    pos_order_id = $ORDER_ID
    items = @(
        @{ pos_product_id = "pos-americano"; quantity = 1 },
        @{ pos_product_id = "pos-latte";     quantity = 2 }
    )
} $POS_HDRS
if ($resp.data.status -eq "processed") { Pass "Deduct stock (1×Americano + 2×Latte)" } else { Fail "Deduct stock" ($resp | ConvertTo-Json) }

# Americano: 18g coffee. Latte×2: 36g coffee. Total: 54g
$coffeeAfter = (Invoke "GET" "$BASE_URL/api/v1/inventory/items/inv-coffee-beans-001" $null $AUTH).data.quantity_in_stock
if ($coffeeAfter -eq ($coffeeBefore - 54)) {
    Pass "Coffee deducted correctly by 54g"
} else {
    Fail "Coffee deduction" "Before=$coffeeBefore After=$coffeeAfter expected=$($coffeeBefore-54)"
}

# Idempotency
$resp = Invoke "POST" "$BASE_URL/api/v1/pos/stock/deduct" @{
    pos_order_id = $ORDER_ID
    items = @(@{ pos_product_id = "pos-americano"; quantity = 1 })
} $POS_HDRS
if ($resp.data.status -eq "already_processed") { Pass "Idempotency: same order ID → already_processed" } else { Fail "Idempotency" ($resp | ConvertTo-Json) }

# Exceed stock
$resp = Invoke "POST" "$BASE_URL/api/v1/pos/stock/deduct" @{
    pos_order_id = "POS-FAIL-$UNIQUE"
    items = @(@{ pos_product_id = "pos-americano"; quantity = 99999 })
} $POS_HDRS
if ($resp._status -eq 400 -or $resp.status -eq "error") { Pass "Exceed stock → 400 error" } else { Fail "Exceed stock check" ($resp | ConvertTo-Json) }

# ─────────────────────────────
Section "WEBHOOKS"
# ─────────────────────────────

$resp = Invoke "POST" "$BASE_URL/api/v1/webhooks" @{
    name = "POS Hook $UNIQUE"; url = "http://localhost:9999/webhook"
    secret = "test-secret-xyz"; events = @("STOCK_UPDATED","STOCK_LOW","STOCK_OUT"); is_active = $true
} $AUTH
if ($resp.data.id) {
    Pass "Create webhook subscription"
    $WEBHOOK_ID = $resp.data.id
} else { Fail "Create webhook" ($resp | ConvertTo-Json) }

$resp = Invoke "GET" "$BASE_URL/api/v1/webhooks" $null $AUTH
if ($resp.data.Count -ge 1) { Pass "List webhooks" } else { Fail "List webhooks" }

$resp = Invoke "GET" "$BASE_URL/api/v1/webhooks/$WEBHOOK_ID" $null $AUTH
if ($resp.data.id -eq $WEBHOOK_ID) { Pass "Get webhook by ID" } else { Fail "Get webhook" }

$resp = Invoke "PUT" "$BASE_URL/api/v1/webhooks/$WEBHOOK_ID" @{
    name = "Updated Hook $UNIQUE"; url = "http://localhost:9999/webhook"
    secret = ""; events = @("STOCK_UPDATED"); is_active = $false
} $AUTH
if ($resp.data.is_active -eq $false) { Pass "Update webhook (deactivate)" } else { Fail "Update webhook" ($resp | ConvertTo-Json) }

# Test delivery (webhook server not running → delivery fails but API returns 200)
$resp = Invoke "POST" "$BASE_URL/api/v1/webhooks/$WEBHOOK_ID/test" $null $AUTH
if ($resp.status -eq "success") { Pass "Test webhook (delivery attempted, 200 returned)" } else { Fail "Test webhook" ($resp | ConvertTo-Json) }

$resp = Invoke "GET" "$BASE_URL/api/v1/webhooks/$WEBHOOK_ID/logs" $null $AUTH
if ($resp.status -eq "success") { Pass "Get webhook delivery logs" } else { Fail "Get webhook logs" }

$resp = Invoke "DELETE" "$BASE_URL/api/v1/webhooks/$WEBHOOK_ID" $null $AUTH
if ($resp.status -eq "success") { Pass "Delete webhook" } else { Fail "Delete webhook" ($resp | ConvertTo-Json) }

# ─────────────────────────────
Section "CLEANUP"
# ─────────────────────────────

$resp = Invoke "DELETE" "$BASE_URL/api/v1/inventory/items/$ITEM_ID" $null $AUTH
if ($resp.status -eq "success") { Pass "Delete test inventory item" } else { Fail "Delete inventory item" ($resp | ConvertTo-Json) }

$resp = Invoke "DELETE" "$BASE_URL/api/v1/products/$PRODUCT_ID" $null $AUTH
if ($resp.status -eq "success") { Pass "Delete test product" } else { Fail "Delete product" ($resp | ConvertTo-Json) }

$resp = Invoke "DELETE" "$BASE_URL/api/v1/suppliers/$SUPPLIER_ID" $null $AUTH
if ($resp.status -eq "success") { Pass "Delete test supplier" } else { Fail "Delete supplier" ($resp | ConvertTo-Json) }

# ─────────────────────────────
Write-Host "`n================================================" -ForegroundColor White
$color = if ($FAIL -eq 0) { "Green" } else { "Yellow" }
Write-Host "  RESULTS: $PASS passed, $FAIL failed" -ForegroundColor $color
Write-Host "================================================`n" -ForegroundColor White
if ($FAIL -gt 0) { exit 1 }
