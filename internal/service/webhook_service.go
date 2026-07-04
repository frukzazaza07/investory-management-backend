package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"investory-management-backend/internal/models"
	"investory-management-backend/internal/repository"
)

func ListWebhooks() ([]models.WebhookSubscription, error) {
	return repository.GetWebhookSubscriptions()
}

func GetWebhook(id string) (*models.WebhookSubscription, error) {
	return repository.GetWebhookSubscriptionByID(id)
}

func CreateWebhook(name, url, secret string, events []string, isActive bool) (*models.WebhookSubscription, error) {
	eventsJSON, _ := json.Marshal(events)
	sub := &models.WebhookSubscription{
		Name:     name,
		URL:      url,
		Secret:   secret,
		Events:   string(eventsJSON),
		IsActive: isActive,
	}
	if err := repository.CreateWebhookSubscription(sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func UpdateWebhook(id, name, url, secret string, events []string, isActive bool) (*models.WebhookSubscription, error) {
	sub, err := repository.GetWebhookSubscriptionByID(id)
	if err != nil {
		return nil, err
	}
	sub.Name = name
	sub.URL = url
	if secret != "" {
		sub.Secret = secret
	}
	eventsJSON, _ := json.Marshal(events)
	sub.Events = string(eventsJSON)
	sub.IsActive = isActive
	if err := repository.UpdateWebhookSubscription(sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func DeleteWebhook(id string) error {
	return repository.DeleteWebhookSubscription(id)
}

func TestWebhook(id string) error {
	sub, err := repository.GetWebhookSubscriptionByID(id)
	if err != nil {
		return err
	}
	payload := map[string]any{
		"event":     "TEST",
		"timestamp": time.Now().UTC(),
		"data":      map[string]string{"message": "test webhook delivery"},
	}
	return sendToSubscription(sub, "TEST", payload)
}

func GetWebhookLogs(subscriptionID string) ([]models.WebhookDeliveryLog, error) {
	return repository.GetWebhookDeliveryLogs(subscriptionID, 50)
}

// FireWebhookEvent sends an event to all active subscriptions asynchronously.
func FireWebhookEvent(event string, item *models.InventoryItem) {
	go func() {
		subs, err := repository.GetActiveWebhookSubscriptions()
		if err != nil {
			log.Printf("webhook: failed to get subscriptions: %v", err)
			return
		}

		var affectedProducts []string
		if item != nil {
			products, _ := repository.GetProductsByInventoryItemID(item.ID)
			for _, p := range products {
				affectedProducts = append(affectedProducts, p.POSProductID)
			}
		}

		payload := map[string]any{
			"event":     event,
			"timestamp": time.Now().UTC(),
		}
		if item != nil {
			payload["data"] = map[string]any{
				"inventory_item_id":     item.ID,
				"sku":                   item.SKU,
				"name":                  item.Name,
				"quantity_in_stock":     item.QuantityInStock,
				"unit":                  item.Unit,
				"affected_pos_products": affectedProducts,
			}
		}

		for i := range subs {
			sub := &subs[i]
			var allowedEvents []string
			json.Unmarshal([]byte(sub.Events), &allowedEvents)
			if !containsEvent(allowedEvents, event) {
				continue
			}
			if err := sendToSubscription(sub, event, payload); err != nil {
				log.Printf("webhook: failed to send to %s: %v", sub.URL, err)
			}
		}
		log.Printf("webhook: event %s sent to %d subscriptions", event, len(subs))
	}()
}

func sendToSubscription(sub *models.WebhookSubscription, event string, payload any) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	signature := computeHMAC(sub.Secret, payloadBytes)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", sub.URL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Inventory-Event", event)
	req.Header.Set("X-Inventory-Signature", "sha256="+signature)

	deliveryLog := &models.WebhookDeliveryLog{
		WebhookSubscriptionID: sub.ID,
		EventType:             event,
		Payload:               string(payloadBytes),
	}

	resp, err := client.Do(req)
	if err != nil {
		deliveryLog.IsSuccess = false
		deliveryLog.ResponseBody = err.Error()
		repository.CreateWebhookDeliveryLog(deliveryLog)
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	deliveryLog.ResponseStatus = resp.StatusCode
	deliveryLog.ResponseBody = string(body)
	deliveryLog.IsSuccess = resp.StatusCode >= 200 && resp.StatusCode < 300
	repository.CreateWebhookDeliveryLog(deliveryLog)
	return nil
}

func computeHMAC(secret string, payload []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

func containsEvent(events []string, target string) bool {
	for _, e := range events {
		if strings.EqualFold(e, target) {
			return true
		}
	}
	return false
}
