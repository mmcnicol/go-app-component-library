// pkg/storybook/notifications.go
package storybook

import (
	"time"
	
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// Notification represents a temporary toast/notification message
type Notification struct {
	ID        string
	Message   string
	Type      string // "info", "success", "warning", "error"
	Duration  time.Duration
	Timestamp time.Time
}

// NotificationManager handles displaying and managing notifications
type NotificationManager struct {
	app.Compo
	
	notifications []Notification
	maxDisplay    int
	autoDismiss   bool
}

var (
	// Global notification manager instance
	notificationManager *NotificationManager
	// Counter for generating unique IDs
	notificationCounter int
)

// ShowNotification displays a temporary notification message
func ShowNotification(ctx app.Context, message string, notificationType string) {
	ShowNotificationWithDuration(ctx, message, notificationType, 3*time.Second)
}

// ShowNotificationWithDuration displays a notification with custom duration
func ShowNotificationWithDuration(ctx app.Context, message string, notificationType string, duration time.Duration) {
	// Initialize manager if not exists
	if notificationManager == nil {
		notificationManager = &NotificationManager{
			maxDisplay:  5,
			autoDismiss: true,
		}
	}
	
	// Generate unique ID
	notificationCounter++
	id := "notification-" + time.Now().Format("150405") + "-" + string(rune(notificationCounter))
	
	// Create notification
	notification := Notification{
		ID:        id,
		Message:   message,
		Type:      notificationType,
		Duration:  duration,
		Timestamp: time.Now(),
	}
	
	// Dispatch to update UI
	ctx.Dispatch(func(ctx app.Context) {
		// Add to front of slice
		notificationManager.notifications = append([]Notification{notification}, notificationManager.notifications...)
		
		// Limit number of displayed notifications
		if len(notificationManager.notifications) > notificationManager.maxDisplay {
			notificationManager.notifications = notificationManager.notifications[:notificationManager.maxDisplay]
		}
		
		ctx.Update()
	})
	
	// Auto-dismiss after duration
	if duration > 0 {
		go func() {
			time.Sleep(duration)
			ctx.Dispatch(func(ctx app.Context) {
				RemoveNotification(ctx, id)
			})
		}()
	}
}

// RemoveNotification removes a specific notification by ID
func RemoveNotification(ctx app.Context, id string) {
	if notificationManager == nil {
		return
	}
	
	for i, n := range notificationManager.notifications {
		if n.ID == id {
			notificationManager.notifications = append(
				notificationManager.notifications[:i],
				notificationManager.notifications[i+1:]...,
			)
			ctx.Update()
			break
		}
	}
}

// ClearAllNotifications removes all active notifications
func ClearAllNotifications(ctx app.Context) {
	if notificationManager != nil {
		notificationManager.notifications = []Notification{}
		ctx.Update()
	}
}

// NotificationComponent renders the notification container and active notifications
type NotificationComponent struct {
	app.Compo
}

func (n *NotificationComponent) Render() app.UI {
	if notificationManager == nil || len(notificationManager.notifications) == 0 {
		return app.Div()
	}
	
	// Build notification elements
	notificationElems := make([]app.UI, len(notificationManager.notifications))
	for i, notification := range notificationManager.notifications {
		// Create notification div with key as an attribute
		notificationDiv := app.Div().
			Class("storybook-notification", "notification-"+notification.Type).
			Attr("data-key", notification.ID) // Use data attribute for key instead of Key() method
		
		// Add content
		notificationDiv.Body(
			app.Div().Class("notification-content").Body(
				app.Span().Class("notification-message").Text(notification.Message),
				app.Button().
					Class("notification-close").
					Text("×").
					OnClick(func(ctx app.Context, e app.Event) {
						RemoveNotification(ctx, notification.ID)
					}),
			),
		)
		
		// Add progress bar if duration > 0
		if notification.Duration > 0 {
			notificationDiv.Body(
				app.Div().Class("notification-progress").
					Style("animation-duration", notification.Duration.String()),
			)
		}
		
		notificationElems[i] = notificationDiv
	}
	
	return app.Div().Class("storybook-notifications").Body(
		notificationElems...,
	)
}
