package router

import (
	"investory-management-backend/internal/handler"
	"investory-management-backend/internal/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/swaggo/swag"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/swagger/doc.json", func(c fiber.Ctx) error {
		doc, err := swag.ReadDoc()
		if err != nil {
			return c.Status(500).SendString("swagger doc not found")
		}
		c.Set("Content-Type", "application/json")
		return c.SendString(doc)
	})

	app.Get("/docs", func(c fiber.Ctx) error {
		html := `<!DOCTYPE html>
<html>
<head><title>API Docs</title></head>
<body>
<script id="api-reference" data-url="/swagger/doc.json"></script>
<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`
		c.Set("Content-Type", "text/html")
		return c.SendString(html)
	})

	// Auth
	auth := app.Group("/auth")
	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)

	// POS Integration (API Key auth) — register BEFORE the JWT group to avoid prefix conflict
	pos := app.Group("/api/v1/pos", middleware.POSAPIKey)
	pos.Post("/stock/deduct", handler.DeductStock)
	pos.Get("/stock/levels", handler.GetStockLevels)
	pos.Get("/products/barcode/:barcode", handler.GetProductByBarcode)
	pos.Get("/products/:pos_product_id/availability", handler.CheckProductAvailability)

	// Protected API (JWT)
	api := app.Group("/api", middleware.Protected)
	api.Get("/me", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"user_id": c.Locals("userID"),
			"email":   c.Locals("email"),
		})
	})

	// Suppliers
	suppliers := api.Group("/v1/suppliers")
	suppliers.Get("/", handler.ListSuppliers)
	suppliers.Post("/", handler.CreateSupplier)
	suppliers.Get("/:id", handler.GetSupplier)
	suppliers.Put("/:id", handler.UpdateSupplier)
	suppliers.Delete("/:id", handler.DeleteSupplier)

	// Inventory Items
	inventory := api.Group("/v1/inventory/items")
	inventory.Get("/", handler.ListInventoryItems)
	inventory.Post("/", handler.CreateInventoryItem)
	inventory.Get("/:id", handler.GetInventoryItem)
	inventory.Put("/:id", handler.UpdateInventoryItem)
	inventory.Delete("/:id", handler.DeleteInventoryItem)
	inventory.Post("/:id/adjust", handler.AdjustInventoryStock)
	inventory.Get("/:id/transactions", handler.GetInventoryItemTransactions)

	// Products
	products := api.Group("/v1/products")
	products.Get("/", handler.ListProducts)
	products.Post("/", handler.CreateProduct)
	products.Get("/barcode/:barcode", handler.GetProductByBarcode)
	products.Get("/:id", handler.GetProduct)
	products.Put("/:id", handler.UpdateProduct)
	products.Delete("/:id", handler.DeleteProduct)
	products.Get("/:id/bom", handler.GetProductBOM)
	products.Put("/:id/bom", handler.UpdateProductBOM)

	// Purchase Orders
	po := api.Group("/v1/purchase-orders")
	po.Get("/", handler.ListPurchaseOrders)
	po.Post("/", handler.CreatePurchaseOrder)
	po.Get("/:id", handler.GetPurchaseOrder)
	po.Put("/:id", handler.UpdatePurchaseOrder)
	po.Post("/:id/receive", handler.ReceivePurchaseOrder)
	po.Post("/:id/cancel", handler.CancelPurchaseOrder)

	// Webhooks
	webhooks := api.Group("/v1/webhooks")
	webhooks.Get("/", handler.ListWebhooks)
	webhooks.Post("/", handler.CreateWebhook)
	webhooks.Get("/:id", handler.GetWebhook)
	webhooks.Put("/:id", handler.UpdateWebhook)
	webhooks.Delete("/:id", handler.DeleteWebhook)
	webhooks.Post("/:id/test", handler.TestWebhook)
	webhooks.Get("/:id/logs", handler.GetWebhookLogs)

}
