package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	db "github.com/OCD-Labs/store-hub/db/sqlc"
	"github.com/OCD-Labs/store-hub/template/email_tmpl"
	"github.com/OCD-Labs/store-hub/util"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	// TaskSendAccessInvitationEmail represents the name of the task that sends the store access invitation email.
	TaskSendAccessInvitationEmail = "task:send_access_invitation_email"
)

type PayloadSendAccessInvitation struct {
	InviterID        int64  `json:"inviter_id"`
	InviteeAccountID string `json:"invitee_account_id"`
	InviteeProfileImgUrl string `json:"invitee_profile_img_url"`
	InviteeID        int64  `json:"invitee_id"`
	InviteeEmail     string `json:"invitee_email"`
	AccessLevel      int32  `json:"access_level"`
	StoreID          int64  `json:"store_id"`
	ClientIp         string `json:"client_ip"`
}

type TokenExtra struct {
	AccessLevel int32 `json:"access_level"`
	InviteeID   int64 `json:"invitee_id"`
	StoreID     int64 `json:"store_id"`
}

// DistributeTaskSendAccessInvitation enqueues the given task to be processed by a worker. It returns an error if the task could
// not be enqueued.
func (distributor *RedisTaskDistributor) DistributeTaskSendAccessInvitation(
	ctx context.Context,
	payload *PayloadSendAccessInvitation,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendAccessInvitationEmail, jsonPayload, opts...)

	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("enqueued task")

	return nil
}

// ProcessTaskSendVerifyEmail processes a 'TaskSendAccessInvitationEmail' task.
func (processor *RedisTaskProcessor) ProcessTaskSendAccessInvitation(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendAccessInvitation
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	inviter, err := processor.dbStore.GetUserByID(ctx, payload.InviterID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	extra := &TokenExtra{
		AccessLevel: payload.AccessLevel,
		InviteeID:   payload.InviteeID,
		StoreID:     payload.StoreID,
	}

	token, tokenPayload, err := processor.tokenMaker.CreateToken(inviter.ID, inviter.AccountID, 25*time.Minute, extra)
	if err != nil {
		return err
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	accessInvitationSession, err := processor.dbStore.CreateSession(ctx, db.CreateSessionParams{
		ID:        uuid,
		UserID:    inviter.ID,
		Token:     util.Extract(token),
		Scope:     "access_invitation_email",
		ClientIp:  payload.ClientIp,
		IsBlocked: false,
		ExpiresAt: tokenPayload.ExpiredAt,
	})
	if err != nil {
		werr := fmt.Errorf("failed to create access invitation email session: %s", err.Error())
		return werr
	}

	store, err := processor.dbStore.GetStoreByID(ctx, payload.StoreID)
	if err != nil {
		return fmt.Errorf("failed to get store: %w", err)
	}

	inviteeAcctId := util.SanitizeAccountID(payload.InviteeAccountID, processor.configs.NEARNetwork)
	inviterAcctId := util.SanitizeAccountID(inviter.AccountID, processor.configs.NEARNetwork)
	accessLevelDescription, detailedExplanation := generateAccessLevelInfo(payload.AccessLevel)

	acceptInvitationURL := fmt.Sprintf(
		"http://store-hub-frontend.vercel.app/access-confirmation?sth_code=%s&payload=%s",
		accessInvitationSession.Token,
		util.EncodeToBase64([]byte(fmt.Sprintf(
			`{
				"store_name": "%s",
				"store_id": %d,
				"store_access_level": "%s",
				"store_profile_img_url": "%s",
				"inviter_account_id": "%s",
				"invitee_account_id": "%s",
				"invitee_profile_img_url": "%s",
				"inviter_profile_img_url": "%s",
				"message": "%s"
			}`, 
			store.Name, store.ID, accessLevelDescription, 
			store.ProfileImageUrl, inviterAcctId, 
			inviteeAcctId, inviter.ProfileImageUrl.String, 
			payload.InviteeProfileImgUrl, detailedExplanation,
		))),
	)

	subject := fmt.Sprintf("Invitation to Manage %s", store.Name)
	content := email_tmpl.CoOwnershipAccessTmpl(
		inviteeAcctId, 
		inviterAcctId, 
		store.Name, 
		accessLevelDescription, 
		detailedExplanation, 
		acceptInvitationURL,
	)
	to := []string{payload.InviteeEmail}
	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send access invitation email: %w", err)
	}

	log.Info().Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("receiver_email", payload.InviteeEmail).
		Str("sender_email", inviter.Email).
		Msg("processed task")

	return nil
}

func generateAccessLevelInfo(accessLevel int32) (accessLevelDescription, detailedExplanation string) {
	switch accessLevel {
	case util.FULLACCESS:
		accessLevelDescription = "Full Access"
		detailedExplanation = "This access level grants you complete control over the store. You can manage products, view and manage sales, orders, and have access to financial data. Essentially, you have the same privileges as the primary store owner."
	case util.PRODUCTINVENTORYACCESS:
		accessLevelDescription = "Product Inventory Access"
		detailedExplanation = "This access level allows you to manage the store's product inventory. You can add, update, or remove products, and view stock levels. However, you won't have access to sales, orders, or financial data."
	case util.SALESACCESS:
		accessLevelDescription = "Sales Access"
		detailedExplanation = "With this access level, you can view and manage the store's sales data. This includes viewing sales reports, trends, and customer data related to sales. However, you won't have access to product inventory, orders, or financial data."
	case util.ORDERSACCESS:
		accessLevelDescription = "Orders Access"
		detailedExplanation = "This access level grants you the ability to view and manage customer orders. You can process orders, handle shipping, and manage customer inquiries related to their orders. You won't have access to product inventory, sales, or financial data."
	case util.FINANCIALACCESS:
		accessLevelDescription = "Financial Access"
		detailedExplanation = "With Financial Access, you can view and manage the store's financial data. This includes revenue, expenses, and profit reports. However, you won't have access to product inventory, sales, or orders."
	}

	return
}
