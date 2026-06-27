package repository

import (
	"investory-management-backend/internal/database"
	"investory-management-backend/internal/models"
)

func GetWebhookSubscriptions() ([]models.WebhookSubscription, error) {
	var subs []models.WebhookSubscription
	err := database.DB.Order("created_at DESC").Find(&subs).Error
	return subs, err
}

func GetWebhookSubscriptionByID(id string) (*models.WebhookSubscription, error) {
	var sub models.WebhookSubscription
	err := database.DB.First(&sub, "id = ?", id).Error
	return &sub, err
}

func GetActiveWebhookSubscriptions() ([]models.WebhookSubscription, error) {
	var subs []models.WebhookSubscription
	err := database.DB.Where("is_active = ?", true).Find(&subs).Error
	return subs, err
}

func CreateWebhookSubscription(sub *models.WebhookSubscription) error {
	return database.DB.Create(sub).Error
}

func UpdateWebhookSubscription(sub *models.WebhookSubscription) error {
	return database.DB.Save(sub).Error
}

func DeleteWebhookSubscription(id string) error {
	return database.DB.Delete(&models.WebhookSubscription{}, "id = ?", id).Error
}

func CreateWebhookDeliveryLog(log *models.WebhookDeliveryLog) error {
	return database.DB.Create(log).Error
}

func GetWebhookDeliveryLogs(subscriptionID string, limit int) ([]models.WebhookDeliveryLog, error) {
	var logs []models.WebhookDeliveryLog
	err := database.DB.
		Where("webhook_subscription_id = ?", subscriptionID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}
