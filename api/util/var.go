package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/env"
	"git.edenfarm.id/cuxs/orm"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	URLFileToUpload             = env.GetString("API_URL_FILE_TO_UPLOAD", "http://storage.edenfarm.tech:8080/v1/create")
	ResponseURLFromPostUpload   = env.GetString("RESPONSE_URL_FROM_POST_UPLOAD", "http://storage.edenfarm.tech/")
	ErrorLogs                   = env.GetString("ERROR_LOG", "http://127.0.0.1:8082/v1/error/log")
	ExportDirectory             = env.GetString("EXPORT_DIRECTORY", "")
	ServerKeyFireBase           = env.GetString("FIREBASE_TOKEN_NOTIFICATION", "AAAASVXBPss:APA91bFJyw6Ixqiq-4muQAwnZ230PdHrwMx8qSXyvrjhFi0i1yOmAco42yFNaAmmG2tyzXMXV7oB_thlg0sMoikcY08wV2lc6reV9ett1iavJXi7oQjFNbLxjjRZMBnMGlvEaPD4nTjB")
	PickingServerKeyFireBase    = env.GetString("PICKING_FIREBASE_TOKEN_NOTIFICATION", "AAAAQshVU9k:APA91bF98A9beYvZIQkbyBH1qZHUB_bFLkJ9tV8ENZ65U6CjHXCAQiXQWGXZCOBn0vBX_PjSqUZvOhWdcAS3hDk8c9sZtuyTiZbq7Kb5wuuygkF8I_pJkwgWk24vsoEplZZvyIOOY-2n")
	PurchaserServerKeyFireBase  = env.GetString("PURCHASER_FIREBASE_TOKEN_NOTIFICATION", "AAAA8NkG9Qo:APA91bFsq0vCEIPuzBw3GJqJAxgvQBwsWjusFS_DOmGNCSl2HU2vU4bGc-U6NIRnwwKvFuVD_PrvO-E3Nc63dEcVLSv7e_MThDDJYgW7lWY-iI5MK8c2DaNCm1pHfuD2ZQhAcTOgustv")
	FieldSalesServerKeyFireBase = env.GetString("FIELD_SALES_FIREBASE_TOKEN_NOTIFICATION", "AAAA8DTSIuA:APA91bF3Q0Dj7_NlkMAZvJbaN0kJ4VNZbupNLviE1DyjTumWxbE4bAy4y3TRLl8nfIHw4kRlZaQ5rkKvdDZ_Pe4f7eOhmI_xXvngmr4KqHihQevrIQ7cImVx-cSlQLKoCWpwEocVAU8v")
	CampaignServerKeyFireBase   = env.GetString("CAMPAIGN_FIREBASE_TOKEN_NOTIFICATION", "AAAA8ZYGg4M:APA91bEfZBIRG5sO6Eeic0NUTmVoep69bkOftQLtfI2891yAo94Xfdm7xeWrztAJcQTM894wmYp74wlnJjbuOMx0Q4g80mjsU_gfLUadEFpK_NoRP4npqZhwqvJsuP20s-InlCvf8wbC")
	PostNotifURL                = env.GetString("POST_NOTIFICATION_URL", "http://apinotif.edenfarm.tech/v1/notification/message")
	PostPickingNotifURL         = env.GetString("POST_PICKING_NOTIFICATION_URL", "http://apinotif.edenfarm.tech/v1/notification/message_picking")
	PostPurchaserNotifURL       = env.GetString("POST_PURCHASER_NOTIFICATION_URL", "http://apinotif.edenfarm.tech/v1/notification/message_purchaser")
	PostFieldSalesNotifURL      = env.GetString("POST_FS_NOTIFICATION_URL", "http://apinotif.edenfarm.tech/v1/notification/message_fs")
	PostCampaignNotifURL        = env.GetString("POST_NOTIFICATION_CAMPAIGN_URL", "http://apinotif.edenfarm.tech/v1/notification/message_campaign")
	XenditKey                   = env.GetString("XENDIT_KEY", "xnd_development_1DKBEmYRrzZkQvKgBAycLOqZ2nJljryY6mdXkGrgAkLpmUZCPaGqUlyWQaGFcJ")
	UrlPrint                    = env.GetString("URL_PRINT", "http://127.0.0.1:8000/api/")
	UrlOapi                     = env.GetString("OAPI_HOST", "127.0.0.1:8082")

	S3endpoint        = env.GetString("S3_ENDPOINT", "sgp1.digitaloceanspaces.com")
	S3accessKeyID     = env.GetString("S3_ACCESS_KEY_ID", "NP6LWL6WRHY5DHGX7H72")
	S3secretAccessKey = env.GetString("S3_SECRET_ACCESS_KEY", "/GhdbcC8bNPmOuPbxiOQHt9Cnbsri1Iymf/URmhUppg")
	S3bucketName      = env.GetString("S3_BUCKET_NAME", "file-temp-dev-eden")
	S3bucketNameImage = env.GetString("S3_BUCKET_NAME_IMAGE", "image-dev-eden")
	VroomUrl          = env.GetString("VROOM_URL", "")
	OsrmCar           = env.GetString("OSRM_CAR", "")
	OsrmBike          = env.GetString("OSRM_BIKE", "")

	//talon.one
	TalonHost          = env.GetString("TALON_HOST", "https://edenfarm.asia-southeast1.talon.one")
	TalonApiKey        = env.GetString("TALON_API_KEY", "3c4e95aec441b175c44e3fa5185d92da240b0690c05900e8a479cec25fe21d18")
	TalonApplicationID = env.GetString("TALON_APPLICATION_ID", "1")
	TalonCampaignID    = env.GetString("TALON_CAMPAIGN_ID", "1")
	TalonLoyaltyID     = env.GetString("TALON_LOYALTY_ID", "1")
	TalonEmail         = env.GetString("TALON_EMAIL", "")
	TalonPass          = env.GetString("TALON_PASS", "")
	TalonToken         = env.GetString("TALON_TOKEN", "")
	TalonFile          = env.GetString("TALON_FILE", "")
)

// RandomInt Returns an int >= min, < max
func RandomInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}

func ReplaceSpace(str string) string {
	return strings.ReplaceAll(str, " ", "")
}

func ReplaceUnderscore(str string) string {
	return strings.ReplaceAll(str, " ", "_")
}

func ReplaceNotificationSalesOrder(message, types, code string) string {
	return strings.ReplaceAll(message, types, code)
}

func ReplaceCodeString(message string, replacer map[string]interface{}) string {
	for k, v := range replacer {
		replacerString := strings.NewReplacer(k, v.(string))
		message = replacerString.Replace(message)
	}

	return message
}

func ArrayToString(a []int64, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

func GetTypePromo(statusPromo string) (status string, e error) {
	if statusPromo == "1" {
		status = "voucher"
	} else if statusPromo == "2" {
		status = "delivery"
	}

	return
}

func ConvertStatusMaster(status int8) (st string) {
	if status == 1 {
		st = "active"
	} else if status == 2 {
		st = "archived"
	} else if status == 3 {
		st = "deleted"
	}
	return st
}

func ConvertStatusPicking(status int8) (st string) {
	if status == 1 {
		st = "New"
	} else if status == 2 {
		st = "Finished"
	} else if status == 3 {
		st = "On Progress"
	} else if status == 4 {
		st = "Need Approval"
	} else if status == 5 {
		st = "Picked"
	} else if status == 6 {
		st = "Checking"
	} else if status == 7 {
		st = "Cancelled"
	}
	return st
}

func ConvertPurchaseOrderItemTaxStatus(status int8) (st string) {
	switch status {
	case 1:
		st = "Yes"
	default:
		st = "No"
	}

	return st
}

func ConvertPurchaseInvoiceItemTaxStatus(status int8) (st string) {
	switch status {
	case 1:
		st = "Yes"
	default:
		st = "No"
	}

	return st
}

func ConvertStatusDoc(status int8) (st string) {
	if status == 1 {
		st = "active"
	} else if status == 2 {
		st = "finished"
	} else if status == 3 {
		st = "cancelled"
	} else if status == 4 {
		st = "deleted"
	} else if status == 5 {
		st = "draft"
	} else if status == 6 {
		st = "partial"
	} else if status == 7 {
		st = "on_delivery"
	} else if status == 8 {
		st = "delivered"
	} else if status == 9 {
		st = "invoiced_not_delivered"
	} else if status == 10 {
		st = "invoiced_on_delivery"
	} else if status == 11 {
		st = "invoiced_delivered"
	} else if status == 12 {
		st = "paid_not_delivered"
	} else if status == 13 {
		st = "paid_on_delivery"
	} else if status == 14 {
		st = "new"
	} else if status == 15 {
		st = "registered"
	} else if status == 16 {
		st = "declined"
	}
	return st
}

func CheckPassword(ps string) string {
	if len(ps) < 8 {
		return "Password at least 8 characters"
	}
	//num := `[0-9]{1}`
	//A_Z := `[A-Z]{1}`
	//symbol := `[!@#~$%^&*()+|_]{1}`
	//if b, err := regexp.MatchString(num, ps); !b || err != nil {
	//	return "Must have one number"
	//}
	//if b, err := regexp.MatchString(A_Z, ps); !b || err != nil {
	//	return "Must have one uppercase character"
	//}
	//if b, err := regexp.MatchString(symbol, ps); !b || err != nil {
	//	return "Must have one special character [!@#$%]"
	//}
	return ""
}

func RandomStr(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

// EncIdInStr : convert a string of real id separated by comma into a string of encrypted id separated by comma
func EncIdInStr(str string) (returnStr string) {
	var encId string
	if str == "" {
		return str
	}

	arrStr := strings.Split(str, ",")
	returnStr = common.Encrypt(arrStr[0]) + ","

	for _, v := range arrStr[1:] {
		encId = common.Encrypt(v)

		returnStr = returnStr + encId + ","
	}

	return strings.TrimSuffix(returnStr, ",")
}

// DecryptIdInStr : convert a string of encrypted id separated by comma into a string of real id separated by comma
func DecryptIdInStr(str string) (returnStr string) {
	var realId int64
	if str == "" {
		return str
	}

	arrStr := strings.Split(str, ",")
	returnStr = common.Encrypt(arrStr[0]) + ","

	for _, v := range arrStr[1:] {
		realId, _ = common.Decrypt(v)

		returnStr = returnStr + strconv.Itoa(int(realId)) + ","
	}

	return strings.TrimSuffix(returnStr, ",")
}

//function to get same/duplicate value from array
func GetSameValue(poi []int64) []int64 {
	// Use map to record duplicates as we find them.
	duplicate := map[int64]bool{}
	var result []int64
	var storageTemp []int64
	for _, v := range poi {
		if duplicate[v] == true {
			// data duplicate will append to storageTemp
			storageTemp = append(storageTemp, v)

		} else {
			// Record this element as an encountered element.
			duplicate[v] = true
			result = append(result, v)
			// Append to result slice.
		}
	}
	// return same value data.
	return storageTemp
}

//function to get unique value from array
func GetUniqueValue(poi []int64) []int64 {

	// given sample array[1,1,2,2,3]
	keys := make(map[int64]int64)
	list := []int64{}

	// form mapping of array, to count each index
	// sample mapping [1:2, 2:2, 3:1]
	for _, entry := range poi {
		keys[entry]++
	}

	// get unique-value array from index with count result = 1 or [n:1]
	for k, v := range keys {
		if v == 1 {
			list = append(list, k)
		}
	}
	// return unique-value array
	return list
}

//function to distinct value from array
func RemoveDuplicateValuesString(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
func OpenFile(f string) *os.File {
	r, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE, 0775)
	if err != nil {
		fmt.Println(err)
	}
	return r
}

func UploadToS3(objectName string, filePath string, contentType string) (string, error) {
	ctx := context.Background()

	// Initialize minio client object.
	minioClient, err := minio.New(S3endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(S3accessKeyID, S3secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		return "", err
	}

	_, err = minioClient.FPutObject(ctx, S3bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", err
	}

	// Retrieve URL valid for 60 minutes
	reqParams := make(url.Values)
	preSignedURL, err := minioClient.PresignedGetObject(context.Background(), S3bucketName, objectName, time.Second*3600, reqParams)
	if err != nil {
		return "", err
	}
	return preSignedURL.String(), nil
}

func PostFormValue(client *http.Client, values map[string]io.Reader) (err error) {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client.Transport = tr
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err
		}

	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", URLFileToUpload, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close() // MUST CLOSED THIS
	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}

	_, err = io.Copy(ioutil.Discard, res.Body) // WE READ THE BODY
	if err != nil {
		return err
	}
	return err
}

type getDataPrint struct {
	Data string `json:"data"`
}

// func sendPrint(req *model.SalesOrder) io.ReadCloser {
func SendPrint(req map[string]interface{}, url string) string {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	var client = &http.Client{Transport: tr}

	jsonReq, _ := json.Marshal(req)

	request, _ := http.NewRequest("POST", UrlPrint+url, bytes.NewBuffer(jsonReq))

	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	//req := UnmarshalBody(response)
	defer response.Body.Close() // MUST CLOSED THIS

	var bodyBytes []byte
	if response.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(response.Body)
	}

	var get getDataPrint
	json.Unmarshal(bodyBytes, &get)
	response.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	_, err = io.Copy(ioutil.Discard, response.Body) // WE READ THE BODY
	if err != nil {
		return "read the body"
	}
	return get.Data
}
func GenerateRandomDoc(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	t := time.Now().Unix()
	s := strconv.FormatInt(t, 10)
	return s + string(b)
}

// UploadImageToS3: function for upload image to S3 storage
func UploadImageToS3(fileName, filePath, types string) (string, error) {
	ctx := context.Background()

	// Initialize minio client object.
	minioClient, err := minio.New(S3endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(S3accessKeyID, S3secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		return "", err
	}

	userMetaData := map[string]string{"x-amz-acl": "public-read"}
	_, err = minioClient.FPutObject(ctx, S3bucketNameImage, types+"/"+fileName, filePath, minio.PutObjectOptions{UserMetadata: userMetaData})
	if err != nil {
		return "", err
	}

	u := &url.URL{
		Scheme: "https",
		Host:   S3endpoint,
		Path:   S3bucketNameImage,
	}

	if err != nil {
		return "", err
	}
	return u.String() + "/" + types + "/" + fileName, nil
}

// UploadImageToS3FieldPurchaser: function for upload image to S3 storage for Field Purchaser App
func UploadImageToS3FieldPurchaser(fileName, filePath, types string, public ...string) (data map[string]interface{}, err error) {
	ctx := context.Background()
	userMetaData := map[string]string{"x-amz-acl": "public-read"}

	// Initialize minio client object.
	minioClient, err := minio.New(S3endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(S3accessKeyID, S3secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		return data, err
	}
	if len(public) > 0 {
		if public[0] == "no" {
			userMetaData = map[string]string{}
		}
	}
	_, err = minioClient.FPutObject(ctx, S3bucketNameImage, types+"/"+fileName, filePath, minio.PutObjectOptions{UserMetadata: userMetaData})
	if err != nil {
		return data, err
	}

	u := &url.URL{
		Scheme: "https",
		Host:   S3endpoint,
		Path:   S3bucketNameImage,
	}

	if err != nil {
		return data, err
	}

	reqParams := make(url.Values)
	preSignedURLImage, _ := minioClient.PresignedGetObject(context.Background(), S3bucketNameImage, types+"/"+fileName, time.Second*60, reqParams)

	data = map[string]interface{}{"url": u.String() + "/" + types + "/" + fileName, "presigned_url": preSignedURLImage.String()}

	return data, nil
}

// GetWeekStart : function to get start date of a yearweek
func GetWeekStart(year, week int) time.Time {
	// start from the middle of the year:
	t := time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC)

	// roll back to Monday:
	if wd := t.Weekday(); wd == time.Sunday {
		t = t.AddDate(0, 0, -6)
	} else {
		t = t.AddDate(0, 0, -int(wd)+1)
	}

	// difference in weeks:
	_, w := t.ISOWeek()
	t = t.AddDate(0, 0, (week-w)*7)

	return t
}

// SendIDToOapi : function to send sales order id for webhook
func SendIDToOapi(soID string, token string) error {
	var (
		err      error
		request  *http.Request
		response *http.Response
	)

	client := new(http.Client)
	url := UrlOapi + "/v1/webhook/send"
	requestData := map[string]interface{}{"sales_order_id": soID}
	m, b := map[string]interface{}{"data": requestData}, new(bytes.Buffer)
	json.NewEncoder(b).Encode(m)
	if request, err = http.NewRequest("POST", url, b); err == nil {
		request.Header.Set("Authorization", "Bearer "+token)
		request.Header.Set("Content-Type", "application/json")
		response, err = client.Do(request)

		defer response.Body.Close()
	}

	return err
}

// GetOrderChannel : function to get order channel name from glossary
func GetOrderChannel(value ...string) (name string, e error) {
	orm := orm.NewOrm()
	orm.Using("read_only")
	var qMark string
	for _, _ = range value {
		qMark = qMark + "?,"
	}
	qMark = strings.TrimSuffix(qMark, ",")
	e = orm.Raw("select group_concat(g.value_name order by g.id) from glossary g where g.attribute = 'order_channel' and value_int in ("+qMark+")", value).QueryRow(&name)

	return
}

func IsOrderChannel(valueInt int) bool {
	orm := orm.NewOrm()
	orm.Using("read_only")
	id := 0

	err := orm.Raw("select id from glossary g where g.attribute = 'order_channel' and value_int = ?", valueInt).QueryRow(&id)
	if err != nil && id < 1 {
		return false
	}
	return true
}

func ItemExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)

	if arr.Kind() != reflect.Array {
		panic("Invalid data-type")
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

// ParsePhoneNumberPrefix : parse phone number if has prefix 62 or 0
func ParsePhoneNumberPrefix(value string) string {
	if strings.HasPrefix(value, "62") {
		return strings.Split(value, "62")[1]
	}

	if strings.HasPrefix(value, "+62") {
		return strings.Split(value, "+62")[1]
	}

	if strings.HasPrefix(value, "0") {
		return strings.Split(value, "0")[1]
	}

	return value
}

// PostErrorToSentry: for post error to sentry
func PostErrorToSentry(err error, tag, value string) {
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetTag(tag, value)
		scope.SetLevel(sentry.LevelFatal)
		sentry.CaptureException(err)
	})
}

//StringInSlice: to search string in an array of list string
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
