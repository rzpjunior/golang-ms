package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-notification-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-notification-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-notification-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-notification-service/internal/app/repository"
	"github.com/NaySoftware/go-fcm"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type INotificationCampaignService interface {
	Get(ctx context.Context, req *dto.GetNotificationCampaignRequest) (res []*dto.NotificationCampaignResponse, total int64, err error)
	SendNotificationCampaign(ctx context.Context, req *dto.SendNotificationCampaignRequest) (res *dto.NotificationStatus, e error)
	UpdateRead(ctx context.Context, req *dto.UpdateReadNotificationCampaignRequest) (err error)
	CountUnread(ctx context.Context, req *dto.CountUnreadNotificationCampaignRequest) (count int64, err error)
}

type NotificationCampaignService struct {
	opt                            opt.Options
	RepositoryNotificationCampaign repository.INotificationCampaignRepository
}

func NewNotificationCampaignService() INotificationCampaignService {
	return &NotificationCampaignService{
		opt:                            global.Setup.Common,
		RepositoryNotificationCampaign: repository.NewNotificationCampaignRepository(),
	}
}

func (s *NotificationCampaignService) Get(ctx context.Context, req *dto.GetNotificationCampaignRequest) (res []*dto.NotificationCampaignResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "NotificationCampaignService.Get")
	defer span.End()

	var notificationTransactions []*model.NotificationCampaign
	notificationTransactions, total, err = s.RepositoryNotificationCampaign.Get(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, notificationTransaction := range notificationTransactions {
		mongoId := notificationTransaction.ID.Hex()
		res = append(res, &dto.NotificationCampaignResponse{
			ID:                     mongoId,
			NotificationCampaignID: notificationTransaction.NotificationCampaignID,
			CustomerID:             notificationTransaction.CustomerID,
			UserCustomerID:         notificationTransaction.UserCustomerID,
			FirebaseToken:          notificationTransaction.FirebaseToken,
			RedirectTo:             notificationTransaction.RedirectTo,
			RedirectToName:         notificationTransaction.RedirectToName,
			RedirectValue:          notificationTransaction.RedirectValue,
			RedirectValueName:      notificationTransaction.RedirectValueName,
			Sent:                   notificationTransaction.Sent,
			Opened:                 notificationTransaction.Opened,
			Conversion:             notificationTransaction.Conversion,
			CreatedAt:              notificationTransaction.CreatedAt,
			UpdatedAt:              notificationTransaction.UpdatedAt,
			RetryCount:             notificationTransaction.RetryCount,
			FcmResultStatus:        notificationTransaction.FcmResultStatus,
		})
	}

	return
}

func (s *NotificationCampaignService) UpdateRead(ctx context.Context, req *dto.UpdateReadNotificationCampaignRequest) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "NotificationCampaignService.UpdateRead")
	defer span.End()

	filter := &model.NotificationCampaign{
		CustomerID:             req.CustomerID,
		NotificationCampaignID: req.NotificationCampaignID,
	}

	span.AddEvent("Update read notification transaction")
	err = s.RepositoryNotificationCampaign.UpdateRead(ctx, filter)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *NotificationCampaignService) CountUnread(ctx context.Context, req *dto.CountUnreadNotificationCampaignRequest) (count int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "NotificationCampaignService.CountUnread")
	defer span.End()

	filter := &model.NotificationCampaign{
		CustomerID: req.CustomerID,
		Opened:     2,
	}

	span.AddEvent("Count unread notification transaction")
	count, err = s.RepositoryNotificationCampaign.CountUnread(ctx, filter)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

// send push notification campaign and save log into mongoDB.
func (s *NotificationCampaignService) SendNotificationCampaign(ctx context.Context, req *dto.SendNotificationCampaignRequest) (res *dto.NotificationStatus, e error) {
	ctx, span := s.opt.Trace.Start(ctx, "NotificationCampaignService.CountUnread")
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

	defer span.End()

	deeplink := s.opt.Env.GetString("campaign.deeplink_notification")

	serverKey := s.opt.Env.GetString("firebase.cma_server_key")

	limit := 500
	runtime.GOMAXPROCS(4)

	c1 := make(chan int)
	c2 := make(chan int)

	totalCustomers := len(req.UserCustomer)
	totalSuccessSent := 0
	totalFailedSent := 0

	// execute pararell
	go func(userCustomers []*dto.UserCustomer) {

		// send to Lark bot process is started
		e = s.PostLarkBotNotificationCampaign("start", req.NotificationCampaignCode, req.NotificationCampaignName, req.Title, req.Message, totalCustomers, 0, 0)
		if e != nil {
			return
		}

		var wg sync.WaitGroup

		for i := 0; i < totalCustomers; i += limit {
			j := i + limit
			if j > totalCustomers {
				j = totalCustomers
			}

			wg.Add(1)

			go func(i int, j int) {
				for _, userCustomer := range userCustomers[i:j] {
					notificationCampaign := &model.NotificationCampaign{
						ID:                     primitive.NewObjectID(),
						NotificationCampaignID: req.NotificationCampaignID,
						CustomerID:             utils.ToString(userCustomer.CustomerID),
						RedirectTo:             req.RedirectTo,
						RedirectToName:         strings.ToLower(req.RedirectToName),
						FirebaseToken:          userCustomer.FirebaseToken,
						Sent:                   2,
						Opened:                 2,
						Conversion:             2,
						CreatedAt:              time.Now(),
						RetryCount:             0,
						FcmResultStatus:        "",
					}

					// set value name by redirect_to
					switch req.RedirectTo {
					case 1:
						notificationCampaign.RedirectValue = deeplink + "product/" + req.RedirectValue
						notificationCampaign.RedirectValueName = "product"
					case 2:
						notificationCampaign.RedirectValue = deeplink + "category/" + req.RedirectValue
						notificationCampaign.RedirectValueName = "category"
					case 3:
						notificationCampaign.RedirectValue = deeplink + "cart"
						notificationCampaign.RedirectValueName = "cart"
					case 4:
						notificationCampaign.RedirectValue = deeplink + "webview"
						notificationCampaign.RedirectValueName = req.RedirectValue
					case 5:
						notificationCampaign.RedirectValue = deeplink + "category/0"
						notificationCampaign.RedirectValueName = "promo"
					default:
						notificationCampaign.RedirectValue = deeplink + "home"
						notificationCampaign.RedirectValueName = "home"
					}

					var sent bool
					var status string
					maxRetry := int8(3)

					// retry until 3x
					for {
						// sent to fcm
						sent, status, e = s.PostFCMCampaignNotif(serverKey, userCustomer.FirebaseToken, req.NotificationCampaignName, req.Title, req.Message, notificationCampaign)
						notificationCampaign.FcmResultStatus = status
						if e != nil {
							fmt.Println(e)
							break
						}

						if sent {
							notificationCampaign.Sent = 1
							totalSuccessSent += 1

							// save log to mongoDB
							e = s.RepositoryNotificationCampaign.Insert(ctx, notificationCampaign)
							if e != nil {
								break
							}
							break
						} else {
							notificationCampaign.Sent = 2
							notificationCampaign.RetryCount += 1
							if maxRetry == notificationCampaign.RetryCount {
								totalFailedSent += 1
								// save log to mongoDB
								e = s.RepositoryNotificationCampaign.Insert(ctx, notificationCampaign)
								if e != nil {
									break
								}
								break
							}
						}

					}
				}
				wg.Done()
			}(i, j)
		}

		wg.Wait()

		// send to Lark bot process is finished
		go func() {
			e = s.PostLarkBotNotificationCampaign("finish", req.NotificationCampaignCode, req.NotificationCampaignName, req.Title, req.Message, totalCustomers, totalSuccessSent, totalFailedSent)
			if e != nil {
				return
			}
		}()

		c1 <- totalSuccessSent
		c2 <- totalFailedSent

	}(req.UserCustomer)

	res = &dto.NotificationStatus{
		SuccessSent: int64(<-c1),
		FailedSent:  int64(<-c2),
	}
	return
}

// PostFCMCampaignNotif: function for send notification campaign to
func (s *NotificationCampaignService) PostFCMCampaignNotif(serverKey string, token string, campaignName string, title string, message string, notificationCampaign *model.NotificationCampaign) (sent bool, status string, err error) {
	c := fcm.NewFcmClient(serverKey)
	var NP fcm.NotificationPayload
	var SP fcm.FcmMsg

	NP.Title = title
	NP.Body = message
	NP.Sound = "default"
	SP.Priority = "high"

	data := map[string]string{
		"notification_campaign_id":   notificationCampaign.NotificationCampaignID,
		"notification_campaign_name": campaignName,
		"title":                      title,
		"message":                    message,
		"redirect_to":                strconv.Itoa(int(notificationCampaign.RedirectTo)),
		"redirect_to_name":           notificationCampaign.RedirectToName,
		"redirect_value":             notificationCampaign.RedirectValue,
		"redirect_value_name":        notificationCampaign.RedirectValueName,
		"sent":                       strconv.Itoa(int(notificationCampaign.Sent)),
		"opened":                     strconv.Itoa(int(notificationCampaign.Opened)),
		"conversion":                 strconv.Itoa(int(notificationCampaign.Conversion)),
		"created_at":                 notificationCampaign.CreatedAt.Format(time.RFC3339),
	}

	c.NewFcmMsgTo(token, data)
	c.SetNotificationPayload(&NP)
	res, err := c.Send()
	if err != nil {
		fmt.Println(err)
		return
	}

	results := res.Results
	resultsJson, err := json.Marshal(results)
	if err != nil {
		fmt.Println(err)
		return
	}

	if res.Success == 1 {
		sent = true
		status = string(resultsJson)
	} else {
		sent = false
		status = string(resultsJson)
	}

	return
}

func (s *NotificationCampaignService) PostLarkBotNotificationCampaign(typePost string, campaignCode string, campaignName string, title string, message string, totalCustomer int, totalSuccess int, totalFailed int) (err error) {
	larkBotURL := s.opt.Env.GetString("lark.host")
	var client = &http.Client{}

	var elements []dto.LarkBotMessageElements
	var cardTitle string
	var template string

	switch typePost {
	case "start":
		cardTitle = "✨ START PUSH NOTIF"
		template = "yellow"

		elements = append(elements, dto.LarkBotMessageElements{
			Tag: "div",
			Text: dto.LarkBotMessageText{
				Content: "🕐 " + time.Now().Format("02 Jan 2006 15:04:05"),
				Tag:     "lark_md",
			},
		})

		elements = append(elements, dto.LarkBotMessageElements{
			Tag: "div",
			Text: dto.LarkBotMessageText{
				Content: "📢 Campaign : \n**Code** : \n" + campaignCode + "\n**Name** : \n" + campaignName,
				Tag:     "lark_md",
			},
		})

		elements = append(elements, dto.LarkBotMessageElements{
			Tag: "div",
			Text: dto.LarkBotMessageText{
				Content: "🔔 Notification : \n**Title** : \n" + title + "\n**Message** : \n" + message,
				Tag:     "lark_md",
			},
		})

		totalMerchantStr := strconv.Itoa(totalCustomer)

		elements = append(elements, dto.LarkBotMessageElements{
			Tag: "div",
			Text: dto.LarkBotMessageText{
				Content: "📄 Report : \n**Total Merchant** : \n" + totalMerchantStr,
				Tag:     "lark_md",
			},
		})

	case "finish":
		cardTitle = "☑️ FINISH PUSH NOTIF"
		template = "primary"

		elements = append(elements, dto.LarkBotMessageElements{
			Tag: "div",
			Text: dto.LarkBotMessageText{
				Content: "🕐 " + time.Now().Format("02 Jan 2006 15:04:05"),
				Tag:     "lark_md",
			},
		})

		elements = append(elements, dto.LarkBotMessageElements{
			Tag: "div",
			Text: dto.LarkBotMessageText{
				Content: "📢 Campaign : \n**Code** : \n" + campaignCode + "\n**Name** : \n" + campaignName,
				Tag:     "lark_md",
			},
		})

		elements = append(elements, dto.LarkBotMessageElements{
			Tag: "div",
			Text: dto.LarkBotMessageText{
				Content: "🔔 Notification : \n**Title** : \n" + title + "\n**Message** : \n" + message,
				Tag:     "lark_md",
			},
		})

		totalMerchantStr := strconv.Itoa(totalCustomer)
		totalSuccessStr := strconv.Itoa(totalSuccess)
		totalFailedStr := strconv.Itoa(totalFailed)

		elements = append(elements, dto.LarkBotMessageElements{
			Tag: "div",
			Text: dto.LarkBotMessageText{
				Content: "📄 Report : \n**Total Merchant** : \n" + totalMerchantStr + "\n**Success Sent** : \n" + totalSuccessStr + "\n**Failed Sent** : \n" + totalFailedStr,
				Tag:     "lark_md",
			},
		})
	}

	msg := &dto.LarkBotMessage{
		MsgType: "interactive",
		Card: dto.LarkBotMessageCard{
			Config: dto.LarkBotMessageConfig{
				WideScreenMode: true,
				EnableForward:  true,
			},
			Header: dto.LarkBotMessageHeader{
				Template: template,
				Title: dto.LarkBotMessageHeaderTitle{
					Tag:     "plain_text",
					Content: cardTitle,
				},
			},
			Elements: elements,
		},
	}

	jsonReq, _ := json.Marshal(msg)
	request, err := http.NewRequest("POST", larkBotURL, bytes.NewBuffer(jsonReq))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))

	defer response.Body.Close()
	_, err = io.Copy(ioutil.Discard, response.Body) // WE READ THE BODY
	if err != nil {
		return err
	}
	return err
}
