package i18n

import "github.com/gofiber/fiber/v3"

const (
	EN = "en"
	TH = "th"
)

var messages = map[string]map[string]string{
	// Auth
	"auth.invalid_request":  {EN: "invalid request", TH: "คำขอไม่ถูกต้อง"},
	"auth.register_success": {EN: "registered successfully", TH: "ลงทะเบียนสำเร็จ"},
	"auth.login_success":    {EN: "login successfully", TH: "เข้าสู่ระบบสำเร็จ"},

	// Suppliers
	"supplier.list":   {EN: "suppliers retrieved", TH: "ดึงข้อมูลผู้จัดจำหน่ายสำเร็จ"},
	"supplier.get":    {EN: "supplier retrieved", TH: "ดึงข้อมูลผู้จัดจำหน่ายสำเร็จ"},
	"supplier.create": {EN: "supplier created", TH: "สร้างผู้จัดจำหน่ายสำเร็จ"},
	"supplier.update": {EN: "supplier updated", TH: "อัปเดตผู้จัดจำหน่ายสำเร็จ"},
	"supplier.delete": {EN: "supplier deleted", TH: "ลบผู้จัดจำหน่ายสำเร็จ"},

	// Inventory
	"inventory.list":         {EN: "inventory items retrieved", TH: "ดึงข้อมูลสินค้าคงคลังสำเร็จ"},
	"inventory.get":          {EN: "inventory item retrieved", TH: "ดึงข้อมูลสินค้าคงคลังสำเร็จ"},
	"inventory.create":       {EN: "inventory item created", TH: "สร้างสินค้าคงคลังสำเร็จ"},
	"inventory.update":       {EN: "inventory item updated", TH: "อัปเดตสินค้าคงคลังสำเร็จ"},
	"inventory.delete":       {EN: "inventory item deleted", TH: "ลบสินค้าคงคลังสำเร็จ"},
	"inventory.adjust":       {EN: "stock adjusted", TH: "ปรับสต็อกสำเร็จ"},
	"inventory.transactions": {EN: "transactions retrieved", TH: "ดึงประวัติธุรกรรมสำเร็จ"},

	// Products
	"product.list":       {EN: "products retrieved", TH: "ดึงข้อมูลสินค้าสำเร็จ"},
	"product.get":        {EN: "product retrieved", TH: "ดึงข้อมูลสินค้าสำเร็จ"},
	"product.create":     {EN: "product created", TH: "สร้างสินค้าสำเร็จ"},
	"product.update":     {EN: "product updated", TH: "อัปเดตสินค้าสำเร็จ"},
	"product.delete":     {EN: "product deleted", TH: "ลบสินค้าสำเร็จ"},
	"product.bom_get":    {EN: "BOM retrieved", TH: "ดึงข้อมูลสูตรสำเร็จ"},
	"product.bom_update": {EN: "BOM updated", TH: "อัปเดตสูตรสำเร็จ"},

	// Purchase Orders
	"po.list":   {EN: "purchase orders retrieved", TH: "ดึงข้อมูลใบสั่งซื้อสำเร็จ"},
	"po.get":    {EN: "purchase order retrieved", TH: "ดึงข้อมูลใบสั่งซื้อสำเร็จ"},
	"po.create": {EN: "purchase order created", TH: "สร้างใบสั่งซื้อสำเร็จ"},
	"po.update": {EN: "purchase order updated", TH: "อัปเดตใบสั่งซื้อสำเร็จ"},
	"po.receive": {EN: "goods received", TH: "รับสินค้าสำเร็จ"},
	"po.cancel": {EN: "purchase order cancelled", TH: "ยกเลิกใบสั่งซื้อสำเร็จ"},

	// Webhooks
	"webhook.list":      {EN: "webhooks retrieved", TH: "ดึงข้อมูล webhook สำเร็จ"},
	"webhook.get":       {EN: "webhook retrieved", TH: "ดึงข้อมูล webhook สำเร็จ"},
	"webhook.not_found": {EN: "webhook not found", TH: "ไม่พบ webhook"},
	"webhook.create":    {EN: "webhook created", TH: "สร้าง webhook สำเร็จ"},
	"webhook.update":    {EN: "webhook updated", TH: "อัปเดต webhook สำเร็จ"},
	"webhook.delete":    {EN: "webhook deleted", TH: "ลบ webhook สำเร็จ"},
	"webhook.test_sent": {EN: "test event sent", TH: "ส่ง test event สำเร็จ"},
	"webhook.logs":      {EN: "webhook logs retrieved", TH: "ดึงประวัติ webhook สำเร็จ"},

	// POS
	"pos.stock_deducted": {EN: "stock deducted", TH: "ตัดสต็อกสำเร็จ"},
	"pos.stock_levels":   {EN: "stock levels retrieved", TH: "ดึงข้อมูลสต็อกสำเร็จ"},
	"pos.availability":   {EN: "availability checked", TH: "ตรวจสอบความพร้อมขายสำเร็จ"},

	// Validation errors
	"err.invalid_body":               {EN: "invalid request body", TH: "รูปแบบคำขอไม่ถูกต้อง"},
	"err.invalid_request":            {EN: "invalid request", TH: "คำขอไม่ถูกต้อง"},
	"err.name_required":              {EN: "name is required", TH: "กรุณาระบุชื่อ"},
	"err.sku_name_unit_required":     {EN: "sku, name and unit are required", TH: "กรุณาระบุ SKU ชื่อ และหน่วย"},
	"err.pos_product_id_name_required": {EN: "pos_product_id and name are required", TH: "กรุณาระบุ pos_product_id และชื่อ"},
	"err.supplier_id_required":       {EN: "supplier_id is required", TH: "กรุณาระบุ supplier_id"},
	"err.items_required":             {EN: "items are required", TH: "กรุณาระบุรายการสินค้า"},
	"err.name_url_secret_required":   {EN: "name, url and secret are required", TH: "กรุณาระบุชื่อ URL และ secret"},
	"err.pos_order_id_required":      {EN: "pos_order_id is required", TH: "กรุณาระบุ pos_order_id"},
	"err.quantity_positive":          {EN: "quantity must be greater than 0", TH: "จำนวนต้องมากกว่า 0"},
}

// T returns the translated message for the given language, falling back to EN.
func T(lang, key string) string {
	if m, ok := messages[key]; ok {
		if v, ok := m[lang]; ok && v != "" {
			return v
		}
		return m[EN]
	}
	return key
}

// Lang reads the language set by the language middleware from the request context.
func Lang(c fiber.Ctx) string {
	if lang, ok := c.Locals("lang").(string); ok && lang != "" {
		return lang
	}
	return EN
}
