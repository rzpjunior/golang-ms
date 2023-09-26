package service

import (
	"context"
	rand2 "crypto/rand"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/repository"
)

type ICodeGeneratorService interface {
	GetGenerateCode(ctx context.Context, format string, domain string, length int) (newCode string, err error)
	GenerateCustomerCode(ctx context.Context, format string, domain string) (code string, e error)
	GenerateReferralCode(ctx context.Context, length int) (code string, e error)
}

type CodeGeneratorService struct {
	opt                     opt.Options
	RepositoryCodeGenerator repository.ICodeGeneratorRepository
}

func NewCodeGeneratorService() ICodeGeneratorService {
	return &CodeGeneratorService{
		opt:                     global.Setup.Common,
		RepositoryCodeGenerator: repository.NewCodeGeneratorRepository(),
	}
}

func (s *CodeGeneratorService) GetGenerateCode(ctx context.Context, format string, domain string, length int) (newCode string, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CodeGeneratorService.Get")
	defer span.End()

	var codeLength int
	if length > 0 {
		codeLength = length
	} else {
		codeLength = 4
	}

	var lastCode string
	var lastCodeLength int

	template := format + "#" + fmt.Sprintf("%0"+strconv.Itoa(codeLength)+"d", 1)
	formatLenght := len(format)

	if template != "" {
		var codeGenerator *model.CodeGenerator
		codeGenerator, err = s.RepositoryCodeGenerator.GetLastCode(ctx, domain, format)
		if err != nil {
			newCode = fmt.Sprintf("%s%s", utils.ToUpper(format), fmt.Sprintf("%0"+strconv.Itoa(codeLength)+"d", 1))
		} else {
			lastCode = codeGenerator.Code
			lastCodeLength = len(lastCode)
			tempIncrement := lastCode[formatLenght:lastCodeLength]
			increment, _ := strconv.Atoi(tempIncrement)
			increments := fmt.Sprintf("%0"+strconv.Itoa(codeLength)+"d", increment+1)
			newCode = fmt.Sprintf("%s%s", utils.ToUpper(format), increments)
		}
	}

	err = s.RepositoryCodeGenerator.Create(ctx, &model.CodeGenerator{
		Code:   newCode,
		Domain: domain,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)

		var isDuplicate = strings.Contains(err.Error(), "Duplicate entry")
		if isDuplicate {
			newCode, err = s.GetGenerateCode(ctx, format, domain, 0)
		}
	}

	return
}

// GenerateCustomerCode : function to generate new code for customer
func (s *CodeGeneratorService) GenerateCustomerCode(ctx context.Context, format string, domain string) (code string, e error) {
	randCode := GenerateRandomString("ABCDEFGHIJKLMNOPQRSTUVWXYZ", 4)
	randNumber := GenerateRandomString("0123456789", 4)

	code = format + randCode + randNumber

	if err := s.RepositoryCodeGenerator.Create(ctx, &model.CodeGenerator{
		Code:   code,
		Domain: domain,
	}); err != nil {
		var isDuplicate = strings.Contains(e.Error(), "Duplicate entry")
		if isDuplicate {
			code, e = s.GenerateCustomerCode(ctx, format, domain)
		}
	}

	return
}

func (s *CodeGeneratorService) GenerateReferralCode(ctx context.Context, length int) (string, error) {
	var e error
	var code string
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ret := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand2.Int(rand2.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}
	code = string(ret)
	if err := s.RepositoryCodeGenerator.CreateReferral(ctx, &model.CodeGeneratorReferral{
		Code:      code,
		CreatedAt: time.Now(),
	}); err != nil {
		var isDuplicate = strings.Contains(e.Error(), "Duplicate entry")
		if isDuplicate {
			code, e = s.GenerateReferralCode(ctx, 8)
		}
	}

	return code, nil
}

func GenerateRandomString(charset string, length int) string {
	byteString := make([]byte, length)
	for i := range byteString {
		byteString[i] = charset[rand.Intn(len(charset))]
	}

	return string(byteString)
}
