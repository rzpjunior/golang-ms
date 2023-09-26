package edenlabs

import (
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/edenlabs/edenlabs/validation"
)

func ErrorValidation(field string, message string) error {
	v := bindValidator
	v.lazyinit()
	o := validation.SetError(field, message)
	return o
}

func ErrorRequired(field string) error {
	return ErrorValidation(field, "The "+utils.RemoveUnderscore(field)+" is required.")
}

func ErrorDuplicate(field string) error {
	return ErrorValidation(field, "The "+utils.RemoveUnderscore(field)+" is duplicated.")
}

func ErrorExists(field string) error {
	return ErrorValidation(field, "The "+utils.RemoveUnderscore(field)+" is existed.")
}

func ErrorScanned(field string) error {
	return ErrorValidation(field, "The "+utils.RemoveUnderscore(field)+" is already scanned.")
}

func ErrorInvalid(field string) error {
	return ErrorValidation(field, "The "+utils.RemoveUnderscore(field)+" is invalid value.")
}

func ErrorNotFound(field string) error {
	return ErrorValidation(field, "The "+utils.RemoveUnderscore(field)+" is not found")
}

func ErrorMustEqualOrLess(field, otherField string) error {
	return ErrorValidation(field, "The "+utils.RemoveUnderscore(field)+" is must be equal or less than "+utils.RemoveUnderscore(otherField)+".")
}

func ErrorMustExistInActive(field, otherField string) error {
	return ErrorValidation(field, "The "+utils.RemoveUnderscore(field)+" is must exist in active "+utils.RemoveUnderscore(otherField)+".")
}

func ErrorMustStatus(field string, status string) error {
	return ErrorValidation(field, "The "+utils.RemoveUnderscore(field)+" is must be "+status+".")
}

func ErrorMustActive(field string) error {
	return ErrorValidation(field, "The "+utils.RemoveUnderscore(field)+" is must be active.")
}

func ErrorMustArchived(field string) error {
	return ErrorValidation(field, "The "+utils.RemoveUnderscore(field)+" is must be archived.")
}

func ErrorMustDeleted(field string) error {
	return ErrorValidation(field, "The "+utils.RemoveUnderscore(field)+" is must be deleted.")
}

func ErrorMustDraft(field string) error {
	return ErrorValidation(field, "The "+utils.RemoveUnderscore(field)+" is must be draft.")
}

func ErrorMustTrue(field string) error {
	return ErrorValidation(field, "The "+utils.RemoveUnderscore(field)+" is must be true.")
}

func ErrorMustContain(field, otherField string) error {
	return ErrorValidation(field, "The "+utils.RemoveUnderscore(field)+" is must contain same values as "+utils.RemoveUnderscore(otherField)+".")
}

func ErrorMustZero(field string) error {
	return ErrorValidation(field, "The "+utils.RemoveUnderscore(field)+" is must be 0.")
}

func ErrorMustSame(field, otherField string) error {
	return ErrorValidation(field, "The "+utils.RemoveUnderscore(field)+" is must be same as "+utils.RemoveUnderscore(otherField)+".")
}

func ErrorMustEqualOrGreater(field, otherField string) error {
	return ErrorValidation(field, "The "+utils.RemoveUnderscore(field)+" is must be equal or greater than "+utils.RemoveUnderscore(otherField)+".")
}

func ErrorMustGreater(field, otherField string) error {
	return ErrorValidation(field, "The "+utils.RemoveUnderscore(field)+" is must be greater than "+utils.RemoveUnderscore(otherField)+".")
}

func ErrorMustLater(field, otherField string) error {
	return ErrorValidation(field, "The "+utils.RemoveUnderscore(field)+" is must be later than "+utils.RemoveUnderscore(otherField)+".")
}

func ErrorMustExistInDirectory(field string) error {
	return ErrorValidation(field, "The "+utils.RemoveUnderscore(field)+" is must exist in directory")
}

func ErrorRpcNotFound(service string, domain string) error {
	return ErrorValidation("id", "Failed to get "+utils.RemoveUnderscore(domain)+" or not found in "+utils.RemoveUnderscore(service)+" service")
}

func ErrorRpc(message string) error {
	return ErrorValidation("id", message)
}

func ErrorRowRequired(parent string, index int, field string) error {
	rowField := parent + "." + strconv.Itoa(index) + "." + field + ".invalid"
	return ErrorValidation(rowField, "The "+utils.RemoveUnderscore(field)+" is required.")
}

func ErrorRowDuplicate(parent string, index int, field string) error {
	rowField := parent + "." + strconv.Itoa(index) + "." + field + ".invalid"
	return ErrorValidation(rowField, "The "+utils.RemoveUnderscore(field)+" is duplicated.")
}

func ErrorRowExists(parent string, index int, field string) error {
	rowField := parent + "." + strconv.Itoa(index) + "." + field + ".invalid"
	return ErrorValidation(rowField, "The "+utils.RemoveUnderscore(field)+" is existed.")
}

func ErrorRowScanned(parent string, index int, field string) error {
	rowField := parent + "." + strconv.Itoa(index) + "." + field + ".invalid"
	return ErrorValidation(rowField, "The "+utils.RemoveUnderscore(field)+" is already scanned.")
}

func ErrorRowInvalid(parent string, index int, field string) error {
	rowField := parent + "." + strconv.Itoa(index) + "." + field + ".invalid"
	return ErrorValidation(rowField, "The "+utils.RemoveUnderscore(field)+" is invalid value.")
}

func ErrorRowNotFound(parent string, index int, field string) error {
	rowField := parent + "." + strconv.Itoa(index) + "." + field + ".invalid"
	return ErrorValidation(rowField, "The "+utils.RemoveUnderscore(field)+" is not found.")
}

func ErrorRowMustEqualOrLess(parent string, index int, field string, otherField string) error {
	rowField := parent + "." + strconv.Itoa(index) + "." + field + ".invalid"
	return ErrorValidation(rowField, "The "+utils.RemoveUnderscore(field)+" is must be equal or less than "+utils.RemoveUnderscore(otherField)+".")
}

func ErrorRowMustExistInActive(parent string, index int, field string, otherField string) error {
	rowField := parent + "." + strconv.Itoa(index) + "." + field + ".invalid"
	return ErrorValidation(rowField, "The "+utils.RemoveUnderscore(field)+" is must exist in active "+utils.RemoveUnderscore(otherField)+".")
}

func ErrorRowMustActive(parent string, index int, field string) error {
	rowField := parent + "." + strconv.Itoa(index) + "." + field + ".invalid"
	return ErrorValidation(rowField, "The "+utils.RemoveUnderscore(field)+" is must be active.")
}

func ErrorRowMustArchived(parent string, index int, field string) error {
	rowField := parent + "." + strconv.Itoa(index) + "." + field + ".invalid"
	return ErrorValidation(rowField, "The "+utils.RemoveUnderscore(field)+" is must be archived.")
}

func ErrorRowMustDeleted(parent string, index int, field string) error {
	rowField := parent + "." + strconv.Itoa(index) + "." + field + ".invalid"
	return ErrorValidation(rowField, "The "+utils.RemoveUnderscore(field)+" is must be deleted.")
}

func ErrorRowMustDraft(parent string, index int, field string) error {
	rowField := parent + "." + strconv.Itoa(index) + "." + field + ".invalid"
	return ErrorValidation(rowField, "The "+utils.RemoveUnderscore(field)+" is must be draft.")
}

func ErrorRowMustTrue(parent string, index int, field string) error {
	rowField := parent + "." + strconv.Itoa(index) + "." + field + ".invalid"
	return ErrorValidation(rowField, "The "+utils.RemoveUnderscore(field)+" is must be true.")
}

func ErrorRowMustContain(parent string, index int, field string, otherField string) error {
	rowField := parent + "." + strconv.Itoa(index) + "." + field + ".invalid"
	return ErrorValidation(rowField, "The "+utils.RemoveUnderscore(field)+" is must contain same values as "+utils.RemoveUnderscore(otherField)+".")
}

func ErrorRowMustZero(parent string, index int, field string) error {
	rowField := parent + "." + strconv.Itoa(index) + "." + field + ".invalid"
	return ErrorValidation(rowField, "The "+utils.RemoveUnderscore(field)+" is must be 0.")
}

func ErrorRowMustSame(parent string, index int, field string, otherField string) error {
	rowField := parent + "." + strconv.Itoa(index) + "." + field + ".invalid"
	return ErrorValidation(rowField, "The "+utils.RemoveUnderscore(field)+" is must be same as "+utils.RemoveUnderscore(otherField)+".")
}

func ErrorRowMustEqualOrGreater(parent string, index int, field string, otherField string) error {
	rowField := parent + "." + strconv.Itoa(index) + "." + field + ".invalid"
	return ErrorValidation(rowField, "The "+utils.RemoveUnderscore(field)+" is must be equal or greater than "+utils.RemoveUnderscore(otherField)+".")
}

func ErrorRowMustGreater(parent string, index int, field string, otherField string) error {
	rowField := parent + "." + strconv.Itoa(index) + "." + field + ".invalid"
	return ErrorValidation(rowField, "The "+utils.RemoveUnderscore(field)+" is must be greater than "+utils.RemoveUnderscore(otherField)+".")
}

func ErrorRowMustLater(parent string, index int, field string, otherField string) error {
	rowField := parent + "." + strconv.Itoa(index) + "." + field + ".invalid"
	return ErrorValidation(rowField, "The "+utils.RemoveUnderscore(field)+" is must be later than "+utils.RemoveUnderscore(otherField)+".")
}

func ErrorRowMustExistInDirectory(parent string, index int, field string) error {
	rowField := parent + "." + strconv.Itoa(index) + "." + field + ".invalid"
	return ErrorValidation(rowField, "The "+utils.RemoveUnderscore(field)+" is must exist in directory.")
}
