package twilio

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	types "github.com/twilio/twilio-go"
)

func resourceChatService() *schema.Resource {
	return &schema.Resource{
		Create: resourceChatServiceCreate,
		Read:   resourceChatServiceRead,
		Update: resourceChatServiceUpdate,
		Delete: resourceChatServiceDelete,

		Schema: map[string]*schema.Schema{
			"friendly_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default_service_role_sid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"default_channel_role_sid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"default_channel_creator_role_sid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"read_status_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"reachability_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"typing_indicator_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"consumption_report_interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"pre_webhook_url": {
				Type:     schema.TypeString,
				Optional: true,
				// ValidateFunc: validation.IsURLWithHTTPorHTTPS(), v1.6.0
			},
			"post_webhook_url": {
				Type:     schema.TypeString,
				Optional: true,
				// ValidateFunc: validation.IsURLWithHTTPorHTTPS(), v1.6.0
			},
			"webhook_method": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"webhook_filters": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"pre_webhook_retry_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"post_webhook_retry_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"account_sid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_update": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"limits": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"media": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size_limit_mb": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"compatibility_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"notifications": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"removed_from_channel": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"sound": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"template": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"added_to_channel": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"sound": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"template": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"invited_to_channel": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"sound": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"template": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"new_message": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"sound": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"template": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"log_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
					},
				},
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func expandBaseNotification(base interface{}) (*types.BaseNotification, error) {
	if base == nil {
		return nil, nil
	}

	baseL := base.([]interface{})

	if len(baseL) > 1 {
		return nil, errors.New("cannot specify more than one notification of each type")
	}

	notification := new(types.BaseNotification)

	for _, n := range baseL {
		setting := n.(map[string]interface{})
		notification.Enabled = setting["enabled"].(bool)
		notification.Sound = setting["sound"].(string)
		notification.Template = setting["template"].(string)
	}

	return notification, nil
}

func expandNewMessage(base interface{}) (*types.NewMessage, error) {
	if base == nil {
		return nil, nil
	}

	messageL := base.([]interface{})

	if len(messageL) > 1 {
		return nil, errors.New("cannot specify more than new message notification")
	}

	notification := new(types.NewMessage)

	for _, n := range messageL {
		message := n.(map[string]interface{})
		notification.Enabled = message["enabled"].(bool)
		notification.Sound = message["sound"].(string)
		notification.Template = message["template"].(string)

		if message["badge_count_enabled"] != nil {
			notification.BadgeCountEnabled = message["badge_count_enabled"].(bool)
		}
	}

	return notification, nil
}

func expandNotifications(d *schema.ResourceData) (*types.Notifications, error) {
	if n, ok := d.GetOk("notifications"); ok {
		notifications := new(types.Notifications)
		nL := n.([]interface{})

		if len(nL) > 1 {
			return nil, errors.New("cannot specify notifications more than one time")
		}

		for _, notification := range nL {
			m := notification.(map[string]interface{})
			removedFromChannel, err := expandBaseNotification(m["removed_from_channel"])
			if err != nil {
				return nil, err
			}

			addedToChannel, err := expandBaseNotification(m["added_to_channel"])
			if err != nil {
				return nil, err
			}

			invitedToChannel, err := expandBaseNotification(m["invited_to_channel"])
			if err != nil {
				return nil, err
			}

			newMessage, err := expandNewMessage(m["new_message"])
			if err != nil {
				return nil, err
			}

			notifications.RemovedFromChannel = removedFromChannel
			notifications.AddedToChannel = addedToChannel
			notifications.InvitedToChannel = invitedToChannel
			notifications.NewMessage = newMessage

			if m["log_enabled"] != nil {
				notifications.LogEnabled = m["log_enabled"].(bool)
			}

			return notifications, nil
		}
	}
	return nil, nil
}

func resourceChatServiceParams(d *schema.ResourceData) *types.ChatServiceParams {
	notifications, err := expandNotifications(d)

	if err != nil {
		log.Printf("[DEBUG] Notification Error: %v", err)
	}

	return &types.ChatServiceParams{
		FriendlyName:                 d.Get("friendly_name").(string),
		DefaultServiceRoleSid:        d.Get("default_service_role_sid").(string),
		DefaultChannelRoleSid:        d.Get("default_channel_role_sid").(string),
		DefaultChannelCreatorRoleSid: d.Get("default_channel_creator_role_sid").(string),
		ReadStatusEnabled:            d.Get("read_status_enabled").(bool),
		ReachabilityEnabled:          d.Get("reachability_enabled").(bool),
		TypingIndicatorTimeout:       d.Get("typing_indicator_timeout").(int),
		ConsumptionReportInterval:    d.Get("consumption_report_interval").(int),
		PreWebhookURL:                d.Get("pre_webhook_url").(string),
		PostWebhookURL:               d.Get("post_webhook_url").(string),
		WebhookMethod:                d.Get("webhook_method").(string),
		WebhookFilters:               d.Get("webhook_filters").(*schema.Set).List(),
		PreWebhookRetryCount:         d.Get("pre_webhook_retry_count").(int),
		PostWebhookRetryCount:        d.Get("post_webhook_retry_count").(int),
		Notifications:                notifications,
	}
}

func resourceChatServiceCreate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)

	chatService, err := m.(*Config).Client.Chat.Create(&types.ChatServiceParams{
		FriendlyName: d.Get("friendly_name").(string),
	})

	if err != nil {
		return fmt.Errorf("Error creating Chat Service: %s", err)
	}

	d.SetId(chatService.Sid)
	d.SetPartial("friendly_name")

	if _, err := m.(*Config).Client.Chat.Update(chatService.Sid, resourceChatServiceParams(d)); err != nil {
		return fmt.Errorf("Error creating Chat Service with optional parameters: %s", err)
	}

	d.SetPartial("default_service_role_sid")
	d.SetPartial("default_channel_role_sid")
	d.SetPartial("default_channel_creator_role_sid")
	d.SetPartial("read_status_enabled")
	d.SetPartial("reachability_enabled")
	d.SetPartial("typing_indicator_timeout")
	d.SetPartial("consumption_report_interval")
	d.SetPartial("pre_webhook_url")
	d.SetPartial("post_webhook_url")
	d.SetPartial("webhook_method")
	d.SetPartial("webhook_filters")

	d.Partial(false)

	return resourceChatServiceRead(d, m)
}

func resourceChatServiceRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	chatService, err := m.(*Config).Client.Chat.Read(id)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG]: ChatService Notifications %+v\n", chatService.Notifications.LogEnabled)

	d.Set("sid", chatService.Sid)
	d.Set("account_sid", chatService.AccountSid)
	d.Set("friendly_name", chatService.FriendlyName)
	d.Set("date_created", chatService.DateCreated)
	d.Set("date_updated", chatService.DateUpdated)
	d.Set("default_service_role_sid", chatService.DefaultServiceRoleSid)
	d.Set("default_channel_role_sid", chatService.DefaultChannelRoleSid)
	d.Set("default_channel_creator_role_sid", chatService.DefaultChannelCreatorRoleSid)
	d.Set("read_status_enabled", chatService.ReadStatusEnabled)
	d.Set("reachability_enabled", chatService.ReachabilityEnabled)
	d.Set("typing_indicator_timeout", chatService.TypingIndicatorTimeout)
	d.Set("consumption_report_interval", chatService.ConsumptionReportInterval)
	d.Set("limits", chatService.Limits)
	d.Set("pre_webhook_url", chatService.PreWebhookURL)
	d.Set("post_webhook_url", chatService.PostWebhookURL)
	d.Set("webhook_method", chatService.WebhookMethod)
	d.Set("webhook_filters", chatService.WebhookFilters)
	d.Set("pre_webhook_retry_count", chatService.PreWebhookRetryCount)
	d.Set("post_webhook_retry_count", chatService.PostWebhookRetryCount)
	d.Set("notifications", []interface{}{
		map[string]interface{}{
			"removed_from_channel": []interface{}{
				map[string]interface{}{
					"enabled":  chatService.Notifications.RemovedFromChannel.Enabled,
					"sound":    chatService.Notifications.RemovedFromChannel.Sound,
					"template": chatService.Notifications.RemovedFromChannel.Template,
				},
			},
			"log_enabled": chatService.Notifications.LogEnabled,
			"added_to_channel": []interface{}{
				map[string]interface{}{
					"enabled":  chatService.Notifications.AddedToChannel.Enabled,
					"sound":    chatService.Notifications.AddedToChannel.Sound,
					"template": chatService.Notifications.AddedToChannel.Template,
				},
			},
			"new_message": []interface{}{
				map[string]interface{}{
					"enabled":             chatService.Notifications.NewMessage.Enabled,
					"sound":               chatService.Notifications.NewMessage.Sound,
					"template":            chatService.Notifications.NewMessage.Template,
					"badge_count_enabled": chatService.Notifications.NewMessage.BadgeCountEnabled,
				},
			},
			"invited_to_channel": []interface{}{
				map[string]interface{}{
					"enabled":  chatService.Notifications.InvitedToChannel.Enabled,
					"sound":    chatService.Notifications.InvitedToChannel.Sound,
					"template": chatService.Notifications.InvitedToChannel.Template,
				},
			},
		},
	})
	d.Set("media", []interface{}{
		map[string]interface{}{
			"size_limit_mb":         chatService.Media.SizeLimitMB,
			"compatibility_message": chatService.Media.CompatibilityMessage,
		},
	})
	d.Set("url", chatService.URL)
	d.Set("links", chatService.Links)

	return nil
}

func resourceChatServiceUpdate(d *schema.ResourceData, m interface{}) error {
	if _, err := m.(*Config).Client.Chat.Update(d.Id(), resourceChatServiceParams(d)); err != nil {
		return fmt.Errorf("Error updating Chat Service: %s", err)
	}

	return resourceChatServiceRead(d, m)
}

func resourceChatServiceDelete(d *schema.ResourceData, m interface{}) error {
	if err := m.(*Config).Client.Chat.Delete(d.Id()); err != nil {
		return fmt.Errorf("Error deleting Chat Service: %s", err)
	}

	return nil
}
