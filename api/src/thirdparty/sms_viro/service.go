package sms_viro

import "time"

func updateOTPSmsViro(r updateRequest) (e error) {
	r.OTPOutGoing.DeliveryStatus = r.Results[0].Status.GroupID
	r.OTPOutGoing.UpdatedAt = time.Now()
	r.OTPOutGoing.Save("DeliveryStatus", "UpdatedAt")
	return
}
