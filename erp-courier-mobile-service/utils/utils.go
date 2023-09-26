package util

import (
	"strconv"
	"strings"
)

var (
	ErrInvalidCredential    = "Please recheck Email or Password."
	ErrInvalidCredentialInd = "Mohon mengecek kembali email dan password anda."
)

// ini untuk membuat character pertama menjadi huruf besar
func changeFirstCharToUpper(fieldName string) string {
	return strings.Title(strings.ToLower(fieldName))
}

func ErrorInputRequired(fieldName string) string {
	return "Please enter " + fieldName + "."
}

func ErrorSelectRequired(fieldName string) string {
	return "Please select " + fieldName + "."
}

func ErrorInputRequiredIndo(fieldName string) string {
	return "Silahkan isi " + fieldName + "."
}

func ErrorSelectRequiredIndo(fieldName string) string {
	return "Silahkan pilih " + fieldName + "."
}

func ErrorDuplicate(fieldName string) string {
	return changeFirstCharToUpper(fieldName) + " is duplicate. Please enter another " + fieldName + "."
}

func ErrorDuplicateID(fieldName string) string {
	return "There are " + changeFirstCharToUpper(fieldName) + " duplicate id exist. Can not have duplicate id."
}

func ErrorUnique(fieldName string) string {
	return changeFirstCharToUpper(fieldName) + " has been registered. Please enter another " + fieldName + "."
}

func ErrorInvalidData(fieldName string) string {
	return "Invalid " + fieldName + "."
}

func ErrorInvalidDataInd(fieldName string) string {
	return fieldName + " tidak valid."
}

func ErrorAlreadyScanned(fieldName string) string {
	return changeFirstCharToUpper(fieldName) + " already scanned"
}

func ErrorAlreadyScannedInd(fieldName string) string {
	return changeFirstCharToUpper(fieldName) + " sudah dipindai"
}

func ErrorEqualLess(fieldName1, fieldName2 string) string {
	return changeFirstCharToUpper(fieldName1) + " must be equal or less than " + fieldName2
}

func ErrorEqualLessNotZeroInd(fieldName1, fieldName2 string) string {
	return changeFirstCharToUpper(fieldName1) + " harus sama atau kurang dari " + fieldName2 + " dan juga tidak kurang dari 0."
}

func ErrorLessNotZeroInd(fieldName1, fieldName2 string) string {
	return changeFirstCharToUpper(fieldName1) + " harus kurang dari " + fieldName2 + " dan juga tidak kurang dari 0."
}

func ErrorMustExistInActive(fieldName1, fieldName2 string) string {
	return changeFirstCharToUpper(fieldName1) + " must exist in active " + fieldName2
}

func ErrorActive(fieldName string) string {
	return changeFirstCharToUpper(fieldName) + " must be active."
}

func ErrorArchived(fieldName string) string {
	return changeFirstCharToUpper(fieldName) + " must be archived."
}

func ErrorDeleted(fieldName string) string {
	return changeFirstCharToUpper(fieldName) + " must be deleted."
}

func ErrorDraft(fieldName string) string {
	return changeFirstCharToUpper(fieldName) + " must be draft."
}

func ErrorCharLength(fieldName string, length int) string {
	return changeFirstCharToUpper(fieldName) + " must contain " + strconv.Itoa(length) + " character(s)."
}

func ErrorMustTrue(fieldName string) string {
	return changeFirstCharToUpper(fieldName) + " must be true."
}

func ErrorMustContain(fieldName1, fieldName2 string) string {
	return changeFirstCharToUpper(fieldName1) + " must contain same values as " + changeFirstCharToUpper(fieldName2)
}

func ErrorSelectOne(fieldName string) string {
	return "Please select at least 1 " + changeFirstCharToUpper(fieldName) + "."
}

func ErrorMustExistWarehouse(fieldName1, fieldName2 string) string {
	return "Selected " + fieldName1 + " must exist in " + fieldName2
}

func ErrorNoAssignedWh(fieldName string) string {
	return changeFirstCharToUpper(fieldName) + " must not have assigned warehouse"
}

func ErrorMustZero(fieldName string) string {
	return changeFirstCharToUpper(fieldName) + " must be 0."
}

func ErrorRelated(stateName, fieldName1, fieldName2 string) string {
	return "There are still " + stateName + " " + fieldName1 + " related to this " + fieldName2 + "."
}

func ErrorMustBeSame(fieldName1, fieldName2 string) string {
	return strings.Title(fieldName1) + " must be same as " + fieldName2 + "."
}

func ErrorMustBeSameInd(fieldName1, fieldName2 string) string {
	return strings.Title(fieldName1) + " harus sama seperti " + fieldName2 + "."
}

func ErrorEqualGreater(fieldName1, fieldName2 string) string {
	return strings.Title(fieldName1) + " must be equal or greater than " + fieldName2 + "."
}

func ErrorGreater(fieldName1, fieldName2 string) string {
	return strings.Title(fieldName1) + " must be greater than " + fieldName2 + "."
}

func ErrorGreaterInd(fieldName1, fieldName2 string) string {
	return strings.Title(fieldName1) + " harus lebih besar dari " + fieldName2 + "."
}

func ErrorMustExistInDirectory(fieldName string) string {
	return strings.Title(fieldName) + " must exist in directory"
}

func ErrorSelectMax(maxCount, fieldName string) string {
	return "Please select at most " + maxCount + " " + fieldName + "(s)."
}

func ErrorLater(fieldName1, fieldName2 string) string {
	return strings.Title(fieldName1) + " must be later than " + fieldName2 + "."
}

func ErrorPaymentCombination() string {
	return "Combination Payment Term & Invoice Term is not true"
}

func ErrorCreateDoc(fieldName string, targetDoc string) string {
	return fieldName + " can not be created. This " + targetDoc + " already has valid " + fieldName
}
func ErrorCreateDocStatus(fieldName string, targetDoc string, status string) string {
	return fieldName + " can not be created." + targetDoc + " is " + status
}

func ErrorNotFound(fieldName string) string {
	return strings.Title(fieldName) + " is not found"
}

func ErrorNotInPeriod(fieldName string) string {
	return strings.Title(fieldName) + " is not in active period."
}

func ErrorOutOfPeriod(fieldName string) string {
	return strings.Title(fieldName) + " is out of active period."
}

func ErrorFullyUsed(fieldName string) string {
	return strings.Title(fieldName) + " is fully used."
}

func ErrorNotValidFor(fieldName1, fieldName2 string) string {
	return strings.Title(fieldName1) + " is not valid for this " + fieldName2 + "."
}

func ErrorProductMustAvailable() string {
	return "Product must be available in the selected warehouse"
}

func ErrorAddDocument(fieldName1, fieldName2, fieldName3 string) string {
	return strings.Title(fieldName1) + " can not be created. " + strings.Title(fieldName2) + " is " + fieldName3
}

func ErrorStatusDoc(fieldName1, fieldName2, fieldName3 string) string {
	return strings.Title(fieldName1) + " can not be " + fieldName2 + ". Please check " + fieldName3 + " status"
}

func ErrorDocStatus(fieldName1, fieldName2 string) string {
	return strings.Title(fieldName1) + " status must be " + fieldName2
}

func ErrorPORelatedDoc() string {
	return "Purchase Order can not be cancelled. This Purchase Order already has valid Goods Receipt or Purchase Invoice"
}

func ErrorOneActiveSameCategory() string {
	return "There is active stock opname with the same category"
}

func ErrorOneActiveInWarehouse() string {
	return "There is active stock opname in this warehouse"
}

func ErrorReturnStockCannot0Qty() string {
	return "Return Good Stock and Return Waste Stock can not have 0 quantity at the same time"
}

func ErrorEqualLater(fieldName1, fieldName2 string) string {
	return strings.Title(fieldName1) + " must be equal or later than " + fieldName2
}

func ErrorInputCannotBeSame(fieldName1, fieldName2 string) string {
	return strings.Title(fieldName1) + " cannot be same as " + fieldName2 + ". Please enter another value."
}

func ErrorRangeChar(fieldName1, fieldName2, fieldName3 string) string {
	return strings.Title(fieldName1) + " must contains " + fieldName2 + " - " + fieldName3 + " characters."
}

func ErrorActiveIsPackable(state string) string {
	return "Product must be active and " + state + "."
}

func ErrorExistActivePackingOrder() string {
	return "Product must not exist in active packing order"
}

func ErrorUniqueProduct() string {
	return "Product name with same UOM has been registered."
}

func ErrorHelperAssign() string {
	return "Some helper already have assigned item"
}

func ErrorPhoneNumber() string {
	return "Merchant phone number has been registered."
}

func ErrorSalesOrderOnPicking() string {
	return "Gagal disimpan, Sales Order sudah tidak valid"
}

func ErrorSalesOrderCannotBeEmpty() string {
	return "Selected sales order cannot be empty"
}

func ErrorEqualGreaterInd(fieldName1, fieldName2 string) string {
	return strings.Title(fieldName1) + " harus besar sama dengan " + fieldName2
}

func ErrorNotValidTermConditions() string {
	return "Voucher is not valid, adjust with voucher terms and conditions."
}

func ErrorPickingStatus(fieldName1, fieldName2 string) string {
	return "Status picking order harus " + fieldName1 + " atau " + fieldName2
}

func ErrorPickingSingleStatus(fieldName1 string) string {
	return "Status picking order harus " + fieldName1
}

func ErrorSOLocked() string {
	return "Sales Order is locked. Please contact Sales Support."
}

func ErrorSOLockedInd() string {
	return "Sales Order telah terkunci. Mohon hubungi Sales Support terkait"
}

func ErrorExceedCutOffTime() string {
	return "Exceed cut off time."
}

func ErrorOrderTypeDraft() string {
	return "Sales order type must be draft."
}

func ErrorOrderTypeCantUpdate(fieldName1 string) string {
	return "You can't update this Sales Order because it is locked by " + fieldName1
}

func ErrorActiveInd(fieldName string) string {
	return changeFirstCharToUpper(fieldName) + " harus aktif."
}
func SalesOrderCannotBeUpdated() string {
	return "Order can not be updated. Order already in Checking Process"
}

func ErrorMustBeSameInOneDocument(field1, field2 string) string {
	return changeFirstCharToUpper(field1) + " must be same in one " + field2
}

func ErrorCannotBeSameInOneDocument(field1, field2 string) string {
	return changeFirstCharToUpper(field1) + " can not be same in one " + field2
}

func ErrorHasChange(field string) string {
	return field + " has been change."
}

func ErrorRunsOut(field string) string {
	return field + " has runs out."
}

func ErrorLess(fieldName1, fieldName2 string) string {
	return changeFirstCharToUpper(fieldName1) + " must be less than " + fieldName2
}

func ErrorSavingData(field1 string) string {
	return "Some error occured while saving " + field1 + " to database"
}

func ErrorDecode(field1 string) string {
	return "Some error occured while decoding " + field1
}

func ErrorRangeGreaterLess(field1, field2, field3 string) string {
	return changeFirstCharToUpper(field1) + "must be greater than " + field2 + " and less than " + field3
}

func ErrorWarehouseCoverage() string {
	return "This hub doesn't have warehouse coverage"
}

func ErrorCourierWarehouse() string {
	return "Courier need to be assigned to a warehouse"
}

func ErrorCourierWarehouseInd() string {
	return "Kurir harus memiliki warehouse"
}

func ErrorCourierVehicle() string {
	return "Courier need a vehicle"
}

func ErrorCourierVehicleInd() string {
	return "Kurir membutuhkan kendaraan"
}

func ErrorInvalidWarehouseRouteItem() string {
	return "courier and sales order's warehouse is not the same"
}

func ErrorInvalidWarehouseRouteItemInd() string {
	return "kurir dan warehouse sales order tidak sama"
}

func ErrorStatusNotAcceptable(field1 string) string {
	return changeFirstCharToUpper(field1) + "'s status is not acceptable"
}

func ErrorStatusNotAcceptableInd(field1 string) string {
	return changeFirstCharToUpper(field1) + " mempunyai status yang tidak dapat diproses."
}

func ErrorMustBeDelivery() string {
	return "type must be delivery"
}

func ErrorMustBeDeliveryInd() string {
	return "tipe harus pengantaran"
}

func ErrorJobCourier() string {
	return "this sales order is not this courier job"
}

func ErrorJobCourierInd() string {
	return "sales order ini bukan pekerjaan anda."
}

func ErrorMultipleJob() string {
	return "you have a delivery on going"
}

func ErrorMultipleJobInd() string {
	return "anda mempunyai pengiriman yang sedang berjalan"
}

func ErrorOnProgress(fieldName string) string {
	return changeFirstCharToUpper(fieldName) + " must be on progress."
}

func ErrorOnProgressInd(fieldName string) string {
	return changeFirstCharToUpper(fieldName) + " harus dalam keadaan sedang berlangsung."
}

func ErrorDataRequirement(str1, str2 string) string {
	return changeFirstCharToUpper(str1) + " need " + changeFirstCharToUpper(str2) + " to operate."
}

func ErrorDataRequirementInd(str1, str2 string) string {
	return changeFirstCharToUpper(str1) + " butuh " + changeFirstCharToUpper(str2) + " untuk beroperasi."
}

func ErrorInactive(fieldName string) string {
	return changeFirstCharToUpper(fieldName) + " must be inactive."
}

func ErrorNotAllowedForInd(fieldName1, fieldName2 string) string {
	return strings.Title(fieldName1) + " tidak diperbolehkan untuk " + fieldName2
}

func ReturnExistInd(fieldName string) string {
	return "Pengembalian pengiriman telah dibuat, tidak dapat " + fieldName + " pengiriman"
}

func RequiredDataInd(fieldName string) string {
	return strings.Title(fieldName) + " dibutuhkan untuk melakukan operasi ini."
}

func ErrorNotFoundInd(fieldName string) string {
	return strings.Title(fieldName) + " tidak ditemukan."
}
