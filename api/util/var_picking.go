package util

// GetCustomerTag: Get customer tag name form merchant
func GetCustomerTag(tagSepComma []string) (s string, err error) {

	var res string
	if len(tagSepComma) == 0 {
		res = "-"
	}
	for _, v := range tagSepComma {
		if v == "1" {
			res = "NC"
			break
		} else if v == "8" {
			res = "PC"
			break
		} else {
			res = "-"
		}
	}
	return res, nil
}
