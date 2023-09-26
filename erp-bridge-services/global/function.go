package global

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/env"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"github.com/sirupsen/logrus"
)

func init() {
	env, _ := env.Env("env")
	EnvDatabaseGP = env.GetString("gp.database")
}

// HttpRestApiToMicrosoftGP sends an HTTP request to Microsoft GP API and parses the response into the provided callback struct.
func HttpRestApiToMicrosoftGP(method, endpoint string, bodyRequest, callback interface{}, params ...map[string]string) (err error) {
	// Retrieve environment variables
	env, err := env.Env("env")
	if err != nil {
		return err
	}

	timeSecond := time.Duration(env.GetInt("retry_function.time_second")) * time.Second
	countRetry := env.GetInt("retry_function.count_retry")
	var (
		counter int
	)
	// Retry function if there's error
	return Retry(func() error {
		var (
			hostGP          = env.GetString("gp.host")
			client          = &http.Client{}
			request         = &http.Request{}
			secondRequest   = &http.Request{}
			response        = &http.Response{}
			dataToken       dto.LoginResponse
			redisConnected  bool = true
			paramStr        string
			bodyRequestByte []byte
			statusCallback  map[int]bool
		)
		statusCallback = make(map[int]bool)
		statusCallback[200] = true
		statusCallback[201] = true
		var retry bool
		if counter > 0 {
			retry = true
		}
		if params != nil {
			paramStr = buildParams(params[0], retry)
		}
		counter = counter + 1

		// If bodyRequest not nil, it means the method http request is not GET
		if bodyRequest != nil {
			if bodyRequestByte, err = json.Marshal(bodyRequest); err != nil {
				return err
			}

			request, err = http.NewRequest(method, hostGP+endpoint+paramStr, bytes.NewBuffer(bodyRequestByte))
			secondRequest, err = http.NewRequest(method, hostGP+endpoint+paramStr, bytes.NewBuffer(bodyRequestByte))
		} else {
			request, err = http.NewRequest(method, hostGP+endpoint+paramStr, nil)
			secondRequest, err = http.NewRequest(method, hostGP+endpoint+paramStr, nil)
			fmt.Println(hostGP + endpoint + paramStr)
		}
		if err != nil {
			return err
		}

		// Set request headers
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "*/*")
		request.Header.Set("Connection", "keep-alive")

		if endpoint == "Tokens/access" { // LOGIN
			if response, err = client.Do(request); err != nil {
				return err
			}
			fmt.Println("=========responseLogin===========", response)
			if !statusCallback[response.StatusCode] {
				err = errors.New("Connection to the server could not be established")
				return err
			}

			if err = unmarshalBody(response, &callback); err != nil {
				return err
			}

			defer response.Body.Close()

			return err
		} else { // NON LOGIN
			// Ping redis
			_, err := Setup.Common.Redisx.Ping(context.TODO())
			if err != nil {
				redisConnected = false
				logrus.Error(err)
			}
			fmt.Println(redisConnected, ">>>>>>>>>bodyRequest>>>>>>>>>>>", string(bodyRequestByte))

			if redisConnected {
				if err = Setup.Common.Redisx.GetCache(context.TODO(), "gp", &dataToken.Data); err != nil {
					err = errors.New("Token Not Found not retry Login")
					dataToken.Data.Token = ""
					fmt.Println("----------", err)
				}
			}

			// Set Authorization header with token
			request.Header.Set("Authorization", "Bearer "+dataToken.Data.Token)
			if response, err = client.Do(request); err != nil {
				return err
			}

			// if request catch error response, try login GP
			if response.StatusCode == 401 {
				if dataToken, err = LoginToMicrosoftDynamicGP(); err != nil {
					return err
				}

				if redisConnected {
					if err = Setup.Common.Redisx.GetCache(context.TODO(), "gp", &dataToken.Data); err != nil {
						err = errors.New("Token Not Found retry Login")
						dataToken.Data.Token = ""
						fmt.Println("----------", err)
					}
				}
				secondRequest.Header.Set("Content-Type", "application/json")
				secondRequest.Header.Set("Accept", "*/*")
				secondRequest.Header.Set("Connection", "keep-alive")
				secondRequest.Header.Set("Authorization", "Bearer "+dataToken.Data.Token)

				if response, err = client.Do(secondRequest); err != nil {
					return err
				}

				// if error still persist, return error
				if response.StatusCode == 401 {
					err = errors.New("Data Not Found")
					return err
				}
			}
		}

		if err = unmarshalBody(response, &callback); err != nil {
			return err
		}
		fmt.Println(">>>>>>>>>>callback>>>>>>>>>>", callback)
		defer response.Body.Close()

		return err

	}, countRetry, timeSecond)
}
