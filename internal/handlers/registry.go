package handlers

import (
	"fmt"
	"time"
)

type HandleFunc func(map[string]any) error

var Registry = map[string]HandleFunc{
	"send_email":    SendEmail,
	"update_logs":   UpdateLogs,
	"refresh_token": RefreshToken,
}

func SendEmail(payload map[string]any) error {
	email, ok := payload["email"].(string)
	if !ok {
		return fmt.Errorf("invalid email payload")
	}

	fmt.Printf("📨 Sending email to %s...\n", email)
	time.Sleep(2 * time.Second)
	fmt.Println("✅ Email sent")
	return nil
}

func UpdateLogs(payload map[string]any) error {
	logs, ok := payload["logs"].(string)
	if !ok {
		return fmt.Errorf("❌ invalid log format")
	}

	fmt.Printf("Updating logs...%s/n", logs)
	time.Sleep(2 * time.Second)
	fmt.Println("✅ Logs updated")
	return nil
}

func RefreshToken(payload map[string]any) error {
	userID, ok := payload["user_id"].(string)
	if !ok {
		return fmt.Errorf("invalid user_id payload")
	}

	fmt.Printf("🔄 Refreshing token for user: %s...\n", userID)
	time.Sleep(2 * time.Second)
	newToken := "new_token_123456"
	fmt.Printf("✅ Token refreshed: %s\n", newToken)
	return nil
}