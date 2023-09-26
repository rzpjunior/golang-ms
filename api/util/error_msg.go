// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package util

import (
	"strconv"
	"strings"
)

var (
	ErrInvalidCredential = "Please recheck Email or Password."
)

// ini untuk membuat character pertama menjadi huruf besar
func changeFirstCharToUpper(fieldName string) string {
	return strings.Title(strings.ToLower(fieldName))
}

func ErrorInputRequired(fieldName string) string {
	return "Please enter " + fieldName + "."
}

func ErrorAlphaNum(fieldName string) string {
	return changeFirstCharToUpper(fieldName) + " must be alpha numeric."
}

func ErrorSelectRequired(fieldName string) string {
	return "Please select " + fieldName + "."
}

func ErrorInputRequiredIndo(fieldName string) string {
	return "Silakan isi " + fieldName + "."
}

func ErrorSelectRequiredIndo(fieldName string) string {
	return "Silakan pilih " + fieldName + "."
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

func ErrorAlreadyScanned(fieldName string) string {
	return changeFirstCharToUpper(fieldName) + " already scanned"
}

func ErrorEqualLess(fieldName1, fieldName2 string) string {
	return changeFirstCharToUpper(fieldName1) + " must be equal or less than " + fieldName2
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

func ErrorEqualGreater(fieldName1, fieldName2 string) string {
	return strings.Title(fieldName1) + " must be equal or greater than " + fieldName2 + "."
}

func ErrorGreater(fieldName1, fieldName2 string) string {
	return strings.Title(fieldName1) + " must be greater than " + fieldName2 + "."
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
	return field + " has been changed"
}

func ErrorRunOut(field string) string {
	return field + " has run out"
}

func ErrorLess(fieldName1, fieldName2 string) string {
	return changeFirstCharToUpper(fieldName1) + " must be less than " + fieldName2
}

func ErrorSuspended(fieldName string) string {
	return changeFirstCharToUpper(fieldName) + " is suspended."
}

func ErrorIsBeingUsed(fieldName string) string {
	return strings.Title(fieldName) + " is being used."
}

func ErrorRole(field1, field2 string) string {
	return changeFirstCharToUpper(field1) + " role must be " + field2
}

func ErrorNoChanges(field1 string) string {
	return "There has to be at least one changes in the " + changeFirstCharToUpper(field1)
}

func ErrorRoleAlreadySigned(field1, field2 string) string {
	return "Job function " + changeFirstCharToUpper(field1) + " already signed by " + field2
}

func ErrorNumeric(fieldName string) string {
	return changeFirstCharToUpper(fieldName) + " must be numeric."
}

func ErrorCannotUpdateAfter(fieldName1, fieldName2 string) string {
	return strings.Title(fieldName1) + " cannot be update after " + fieldName2
}

func ErrorCreditLimitExceeded(merchantName string) string {
	return "The amount exceeds the Customer " + merchantName + "'s credit limit."
}

func ErrorGetRemainingInvoice() string {
	return "Failed to get remaining invoice."
}

func ErrorInvalidOrderChannelRestriction(fieldName string) string {
	return "Invalid " + fieldName + "."
}

func ErrorStaffBusy() string {
	return "One of the staffs are currently working at something else"
}

func ErrorType(field1, field2 string) string {
	return changeFirstCharToUpper(field1) + " has to be " + changeFirstCharToUpper(field2) + " type"
}

func ErrorPickingListStaff() string {
	return "The picking list is not this staff job"
}

func ErrorOrderChannelRestriction() string {
	return "this product not eligable to be sold in dashboard."
}

func ErrorExistActivePurchasePlan() string {
	return "Purchase order must not exist in active purchase plan"
}

func ErrorRangeValue(fieldName1, fieldName2, fieldName3 string) string {
	return strings.Title(fieldName1) + " value must between " + fieldName2 + " - " + fieldName3
}

func ErrorCannotCrossBusinessType() string {
	return "Archetype is not match with business type"
}

func ErrorSelectAnother(fieldName1, fieldName2, fieldName3 string) string {
	return "The selected " + fieldName1 + " already " + fieldName2 + " please select another " + fieldName3
}

func ErrorNotAllowedFor(fieldName1, fieldName2 string) string {
	return strings.Title(fieldName1) + " is not allowed for this " + fieldName2
}

func ErrorCannotCancelAfter(fieldName1, fieldName2 string) string {
	return strings.Title(fieldName1) + " cannot be cancel after " + fieldName2
}

func ErrorCreditCustomMinus(fieldName1 string) string {
	return strings.Title(fieldName1) + " value cannot be minus"
}

func InsidePolygon() string {
	return "Pinpoint has to be inside the polygon"
}

func HubMainWarehouse() string {
	return "Hub cannot be a main warehouse"
}

func HubNeedParent() string {
	return "A hub need its parent to serve the subdistrict first"
}

func HubOnlyOne() string {
	return "There can only be one hub for each parent warehouse in a subdistrict"
}

func MainWarehouseDelete() string {
	return "Cannot delete main warehouse"
}

func HubStillExist() string {
	return "There's a hub associated with this warehouse"
}

func IsMainWarehouse() string {
	return "This warehouse is already a main warehouse"
}

func ErrorExistWarehouseCoverage() string {
	return "This warehouse is already serving this subdistrict"
}

func ErrorMustDifferenctSalesPerson(branchCode, fsCode string) string {
	return "Branch " + branchCode + " already have salesperson " + fsCode
}

func ErrorWarehouseCoverage() string {
	return "This hub doesn't have warehouse coverage"
}

func ErrorCreationInProgress(field string) string {
	return "There's a " + field + " creation in progress"
}

func ErrorReferrerCode() string {
	return "Referrer code cannot be same with referral code of merchant"
}

func ErrorRoutingOnGoing() string {
	return "There's a routing on going for this Picking List"
}

func ErrorCannotCancelSalesPayment() string {
	return "Cancel sales payment paid off first"
}

func ErrorCannotCreatePaidOff() string {
	return "Confirm or cancel sales payment in progress first"
}

func ErrorCannotOverPay() string {
	return "Your payment exceeds the bill"
}

func ErrorCannotConfirmBulkPaidOff() string {
	return "Paid off must be on the last sales payment list"
}

func ErrorIntersect(activeStatus, startDate, endDate string) string {
	return "There are already " + activeStatus + " data between " + startDate + " and " + endDate
}

func ErrorSelfPickUp(fieldName string) string {
	return "Choose " + fieldName + " for self pickup order"
}

func ErrorAreaSelfPickUp(areaName string) string {
	return "Area " + areaName + " not available for self pickup"
}

func ErrorOnlyValidFor(fieldName1, fieldName2, fieldName3 string) string {
	return strings.Title(fieldName1) + " is only valid for " + fieldName2 + " " + fieldName3 + "."
}
