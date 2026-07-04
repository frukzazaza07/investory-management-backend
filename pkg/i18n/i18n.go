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
	"inventory.unit.list":    {EN: "units retrieved", TH: "ดึงข้อมูลหน่วยสำเร็จ"},
	"inventory.unit.get":     {EN: "unit retrieved", TH: "ดึงข้อมูลหน่วยสำเร็จ"},
	"inventory.unit.create":  {EN: "unit created", TH: "สร้างหน่วยสำเร็จ"},
	"inventory.unit.update":  {EN: "unit updated", TH: "อัปเดตหน่วยสำเร็จ"},
	"inventory.unit.delete":  {EN: "unit deleted", TH: "ลบหน่วยสำเร็จ"},

	// Products
	"product.list":       {EN: "products retrieved", TH: "ดึงข้อมูลสินค้าสำเร็จ"},
	"product.get":        {EN: "product retrieved", TH: "ดึงข้อมูลสินค้าสำเร็จ"},
	"product.create":     {EN: "product created", TH: "สร้างสินค้าสำเร็จ"},
	"product.update":     {EN: "product updated", TH: "อัปเดตสินค้าสำเร็จ"},
	"product.delete":     {EN: "product deleted", TH: "ลบสินค้าสำเร็จ"},
	"product.bom_get":    {EN: "BOM retrieved", TH: "ดึงข้อมูลสูตรสำเร็จ"},
	"product.bom_update": {EN: "BOM updated", TH: "อัปเดตสูตรสำเร็จ"},

	// Purchase Orders
	"po.list":    {EN: "purchase orders retrieved", TH: "ดึงข้อมูลใบสั่งซื้อสำเร็จ"},
	"po.get":     {EN: "purchase order retrieved", TH: "ดึงข้อมูลใบสั่งซื้อสำเร็จ"},
	"po.create":  {EN: "purchase order created", TH: "สร้างใบสั่งซื้อสำเร็จ"},
	"po.update":  {EN: "purchase order updated", TH: "อัปเดตใบสั่งซื้อสำเร็จ"},
	"po.receive": {EN: "goods received", TH: "รับสินค้าสำเร็จ"},
	"po.cancel":  {EN: "purchase order cancelled", TH: "ยกเลิกใบสั่งซื้อสำเร็จ"},

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
	"err.invalid_body":                 {EN: "invalid request body", TH: "รูปแบบคำขอไม่ถูกต้อง"},
	"err.invalid_request":              {EN: "invalid request", TH: "คำขอไม่ถูกต้อง"},
	"err.name_required":                {EN: "name is required", TH: "กรุณาระบุชื่อ"},
	"err.sku_name_unit_required":       {EN: "sku, name and unit are required", TH: "กรุณาระบุ SKU ชื่อ และหน่วย"},
	"err.unit_invalid":                 {EN: "unit is invalid, see /api/v1/inventory/units for allowed values", TH: "หน่วยไม่ถูกต้อง ดูรายการที่ /api/v1/inventory/units"},
	"err.unit_code_name_required":      {EN: "code and name are required", TH: "กรุณาระบุ code และชื่อ"},
	"err.pos_product_id_name_required": {EN: "pos_product_id and name are required", TH: "กรุณาระบุ pos_product_id และชื่อ"},
	"err.supplier_id_required":         {EN: "supplier_id is required", TH: "กรุณาระบุ supplier_id"},
	"err.items_required":               {EN: "items are required", TH: "กรุณาระบุรายการสินค้า"},
	"err.name_url_secret_required":     {EN: "name, url and secret are required", TH: "กรุณาระบุชื่อ URL และ secret"},
	"err.pos_order_id_required":        {EN: "pos_order_id is required", TH: "กรุณาระบุ pos_order_id"},
	"err.quantity_positive":            {EN: "quantity must be greater than 0", TH: "จำนวนต้องมากกว่า 0"},

	// Api messages
	"api.success":                  {EN: "success", TH: "สำเร็จ"},
	"api.error":                    {EN: "error", TH: "เกิดข้อผิดพลาด"},
	"api.not_found":                {EN: "not found", TH: "ไม่พบข้อมูล"},
	"api.invalid_request":          {EN: "invalid request", TH: "คำขอไม่ถูกต้อง"},
	"api.unauthorized":             {EN: "unauthorized", TH: "ไม่ได้รับอนุญาต"},
	"api.forbidden":                {EN: "forbidden", TH: "ห้ามเข้าถึง"},
	"api.internal_error":           {EN: "internal server error", TH: "เกิดข้อผิดพลาดภายในเซิร์ฟเวอร์"},
	"api.validation_error":         {EN: "validation error", TH: "เกิดข้อผิดพลาดในการตรวจสอบข้อมูล"},
	"api.database_error":           {EN: "database error", TH: "เกิดข้อผิดพลาดในการเข้าถึงฐานข้อมูล"},
	"api.external_service_error":   {EN: "external service error", TH: "เกิดข้อผิดพลาดจากบริการภายนอก"},
	"api.rate_limit_exceeded":      {EN: "rate limit exceeded", TH: "เกินขีดจำกัดการใช้งาน"},
	"api.service_unavailable":      {EN: "service unavailable", TH: "บริการไม่พร้อมใช้งาน"},
	"api.timeout":                  {EN: "request timeout", TH: "หมดเวลาการร้องขอ"},
	"api.conflict":                 {EN: "conflict", TH: "เกิดความขัดแย้ง"},
	"api.not_implemented":          {EN: "not implemented", TH: "ยังไม่ได้ดำเนินการ"},
	"api.bad_gateway":              {EN: "bad gateway", TH: "เกตเวย์ไม่ถูกต้อง"},
	"api.gateway_timeout":          {EN: "gateway timeout", TH: "หมดเวลาการเชื่อมต่อเกตเวย์"},
	"api.unsupported_media_type":   {EN: "unsupported media type", TH: "ประเภทสื่อไม่รองรับ"},
	"api.too_many_requests":        {EN: "too many requests", TH: "คำขอมากเกินไป"},
	"api.precondition_failed":      {EN: "precondition failed", TH: "เงื่อนไขล่วงหน้าไม่สำเร็จ"},
	"api.payload_too_large":        {EN: "payload too large", TH: "ข้อมูลมากเกินไป"},
	"api.uri_too_long":             {EN: "URI too long", TH: "URI ยาวเกินไป"},
	"api.unsupported_http_version": {EN: "unsupported HTTP version", TH: "เวอร์ชัน HTTP ไม่รองรับ"},
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
