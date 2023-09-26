package global

import (
	"context"
	"errors"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/env"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"github.com/sirupsen/logrus"
)

func LoginToMicrosoftDynamicGP() (resp dto.LoginResponse, err error) {
	env, e := env.Env("env")
	if e != nil {
		err = e
	}
	req := new(dto.LoginRequest)

	req.UserName = env.GetString("gp.username")
	req.Password = env.GetString("gp.password")

	err = HttpRestApiToMicrosoftGP("POST", "Tokens/access", req, &resp)
	if err != nil {
		logrus.Error(err.Error())
		err = errors.New("Connection to the server could not be established")
		return
	}

	if resp.Code != 200 {
		logrus.Error("Error Login: " + resp.Message)
		err = errors.New("Connection to the server could not be established")
		return
	}

	// PING and add Token data to Redis
	if _, e := Setup.Common.Redisx.Ping(context.TODO()); e == nil {
		Setup.Common.Redisx.SetCache(context.TODO(), "gp", resp.Data, time.Hour*23)
	}

	return resp, err
}
