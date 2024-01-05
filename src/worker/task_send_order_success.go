package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Chien179/NMCBookstoreBE/src/models"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskSendOrderSuccess = "task:send_order_success"

type PayloadSendOrderSuccess struct {
	Order models.OrderReponse `json:"order"`
}

func (distributior *RedisTaskDistributor) DistributeTaskSendOrderSuccess(
	ctx context.Context,
	payload *PayloadSendOrderSuccess,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendOrderSuccess, jsonPayload, opts...)
	info, err := distributior.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("enqueued task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendOrderSuccess(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendOrderSuccess
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err, asynq.SkipRetry)
	}

	user, err := processor.store.GetUser(ctx, payload.Order.Username)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	subject := "Order Confirmation"
	content := fmt.Sprintf(`<!DOCTYPE html>
	<html>
	<head>
	<title></title>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<meta http-equiv="X-UA-Compatible" content="IE=edge" />
	<style type="text/css">
	
	body, table, td, a { -webkit-text-size-adjust: 100%; -ms-text-size-adjust: 100%; }
	table, td { mso-table-lspace: 0pt; mso-table-rspace: 0pt; }
	img { -ms-interpolation-mode: bicubic; }
	
	img { border: 0; height: auto; line-height: 100%; outline: none; text-decoration: none; }
	table { border-collapse: collapse !important; }
	body { height: 100% !important; margin: 0 !important; padding: 0 !important; width: 100% !important; }
	
	
	a[x-apple-data-detectors] {
		color: inherit !important;
		text-decoration: none !important;
		font-size: inherit !important;
		font-family: inherit !important;
		font-weight: inherit !important;
		line-height: inherit !important;
	}
	
	@media screen and (max-width: 480px) {
		.mobile-hide {
			display: none !important;
		}
		.mobile-center {
			text-align: center !important;
		}
	}
	div[style*="margin: 16px 0;"] { margin: 0 !important; }
	</style>
	<body style="margin: 0 !important; padding: 0 !important; background-color: #eeeeee;" bgcolor="#eeeeee">
	
	
	<div style="display: none; font-size: 1px; color: #fefefe; line-height: 1px; font-family: Open Sans, Helvetica, Arial, sans-serif; max-height: 0px; max-width: 0px; opacity: 0; overflow: hidden;">
	For what reason would it be advisable for me to think about business content? That might be little bit risky to have crew member like them. 
	</div>
	
	<table border="0" cellpadding="0" cellspacing="0" width="100%">
		<tr>
			<td align="center" style="background-color: #eeeeee;" bgcolor="#eeeeee">
			
			<table align="center" border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width:600px;">
				<tr>
					<td align="center" valign="top" style="font-size:0; padding: 35px;" bgcolor="#F44336">
				   
					<div style="display:inline-block; max-width:50%; min-width:100px; vertical-align:top; width:100%;">
						<table align="left" border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width:300px;">
							<tr>
								<td align="left" valign="top" style="font-family: Open Sans, Helvetica, Arial, sans-serif; font-size: 36px; font-weight: 800; line-height: 48px;" class="mobile-center">
									<h1 style="font-size: 36px; font-weight: 800; margin: 0; color: #ffffff;">NMC Bookstore</h1>
								</td>
							</tr>
						</table>
					</div>
					
					<div style="display:inline-block; max-width:50%; min-width:100px; vertical-align:top; width:100%;" class="mobile-hide">
						<table align="left" border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width:300px;">
							<tr>
								<td align="right" valign="top" style="font-family: Open Sans, Helvetica, Arial, sans-serif; font-size: 48px; font-weight: 400; line-height: 48px;">
									<table cellspacing="0" cellpadding="0" border="0" align="right">
										<tr>
											<td style="font-family: Open Sans, Helvetica, Arial, sans-serif; font-size: 18px; font-weight: 400;">
												<p style="font-size: 18px; font-weight: 400; margin: 0; color: #ffffff;"><a href="#" target="_blank" style="color: #ffffff; text-decoration: none;"></a></p>
											</td>
											<td style="font-family: Open Sans, Helvetica, Arial, sans-serif; font-size: 18px; font-weight: 400; line-height: 24px;">
												<a href="#" target="_blank" style="color: #ffffff; text-decoration: none;"><img src="https://img.icons8.com/color/48/000000/small-business.png" width="27" height="23" style="display: block; border: 0px;"/></a>
											</td>
										</tr>
									</table>
								</td>
							</tr>
						</table>
					</div>
				  
					</td>
				</tr>
				<tr>
					<td align="center" style="padding: 35px 35px 20px 35px; background-color: #ffffff;" bgcolor="#ffffff">
					<table align="center" border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width:600px;">
						<tr>
							<td align="center" style="font-family: Open Sans, Helvetica, Arial, sans-serif; font-size: 16px; font-weight: 400; line-height: 24px; padding-top: 25px;">
								<img src="https://img.icons8.com/carbon-copy/100/000000/checked-checkbox.png" width="125" height="120" style="display: block; border: 0px;" /><br>
								<h2 style="font-size: 30px; font-weight: 800; line-height: 36px; color: #333333; margin: 0;">
									Thank You For Your Order!
								</h2>
							</td>
						</tr>
						<tr>
							<td align="left" style="padding-top: 20px;">
								<table cellspacing="0" cellpadding="0" border="0" width="100%">
									<tr>
										<td width="75%" align="left" bgcolor="#eeeeee" style="font-family: Open Sans, Helvetica, Arial, sans-serif; font-size: 16px; font-weight: 800; line-height: 24px; padding: 10px;">
											Order Confirmation #
										</td>
										<td width="25%" align="left" bgcolor="#eeeeee" style="font-family: Open Sans, Helvetica, Arial, sans-serif; font-size: 16px; font-weight: 800; line-height: 24px; padding: 10px;">
											%d
										</td>
									</tr>
								</table>
							</td>
						</tr>
						<tr>
							<td align="left" style="padding-top: 20px;">
								<table cellspacing="0" cellpadding="0" border="0" width="100%">
									<tr>
										<td width="75%" align="left" style="font-family: Open Sans, Helvetica, Arial, sans-serif; font-size: 16px; font-weight: 800; line-height: 24px; padding: 10px; border-top: 3px solid #eeeeee; border-bottom: 3px solid #eeeeee;">
											TOTAL
										</td>
										<td width="25%" align="left" style="font-family: Open Sans, Helvetica, Arial, sans-serif; font-size: 16px; font-weight: 800; line-height: 24px; padding: 10px; border-top: 3px solid #eeeeee; border-bottom: 3px solid #eeeeee;">
											%f
										</td>
									</tr>
								</table>
							</td>
						</tr>
					</table>
					
					</td>
				</tr>
				 <tr>
					<td align="center" height="100%" valign="top" width="100%" style="padding: 0 35px 35px 35px; background-color: #ffffff;" bgcolor="#ffffff">
					<table align="center" border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width:660px;">
						<tr>
							<td align="center" valign="top" style="font-size:0;">
								<div style="display:inline-block; max-width:50%; min-width:240px; vertical-align:top; width:100%;">
	
									<table align="left" border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width:300px;">
										<tr>
											<td align="left" valign="top" style="font-family: Open Sans, Helvetica, Arial, sans-serif; font-size: 16px; font-weight: 400; line-height: 24px;">
												<p style="font-weight: 800;">Delivery Address</p>
												<p>%s</p>
											</td>
										</tr>
									</table>
								</div>
							</td>
						</tr>
					</table>
					</td>
				</tr>
				<tr>
					<td align="center" style="padding: 35px; background-color: #ffffff;" bgcolor="#ffffff">
					<table align="center" border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width:600px;">
						<tr>
							<td align="center">
								<img src="logo-footer.png" width="37" height="37" style="display: block; border: 0px;"/>
							</td>
						</tr>
						<tr>
							<td align="center" style="font-family: Open Sans, Helvetica, Arial, sans-serif; font-size: 14px; font-weight: 400; line-height: 24px; padding: 5px 0 10px 0;">
								<p style="font-size: 14px; font-weight: 800; line-height: 18px; color: #333333;">
									So 1 Vo Van Ngan, Linh Chieu, Thu Duc
								</p>
							</td>
						</tr>
					</table>
					</td>
				</tr>
			</table>
			</td>
		</tr>
	</table>
		
	</body>
	</html>
	`, payload.Order.ID, payload.Order.SubTotal, payload.Order.ToAddress)
	to := []string{user.Email}

	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send report review email %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("email", user.Email).
		Msg("processed task")
	return nil
}
