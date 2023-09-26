package xendit_transaction

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"time"

	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/datamodel/model"

	"git.edenfarm.id/project-version2/api/util"
)

type fixedVaRequest struct {
	AccountNumber     string                `json:"account_number" valid:"required"`
	TransactionDate   string                `json:"transaction_date"`
	TransactionTime   string                `json:"transaction_time"`
	PaidAmount        float64               `json:"paid_amount"`
	Token             string                `json:"token"`
	TransactionDateAt time.Time             `json:"-"`
	TransactionTimeAt time.Time             `json:"-"`
	MerchantAccNum    *model.MerchantAccNum `json:"-"`
	PaymentChannel    *model.PaymentChannel `json:"-"`
}

func (c *fixedVaRequest) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	var e error

	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	key := []byte("joyfuls joy j0y5")
	dec := decrypt(key, c.Token)
	token := fmt.Sprintf("%s", dec)
	if token != "hey please push on thursday, it will make me happy" {
		o.Failure("token", "invalid token")
	}

	if c.TransactionDate != "" {
		if c.TransactionDateAt, e = time.Parse("2006-01-02", c.TransactionDate); e != nil {
			o.Failure("transaction_date.invalid", "invalid date")
		}
	}
	if e = orSelect.Raw("SELECT * FROM merchant_acc_num WHERE account_number = ? ", c.AccountNumber).QueryRow(&c.MerchantAccNum); e != nil {
		o.Failure("account_number.invalid", util.ErrorInvalidData("merchant account number"))
	}

	return o
}

func (c *fixedVaRequest) Messages() map[string]string {
	return map[string]string{}
}

// decrypt from base64 to decrypted string
func decrypt(key []byte, cryptoText string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		fmt.Println("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}
