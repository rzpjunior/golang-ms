package dummy

import "github.com/bxcodec/faker/v3"

func NameList(length int) (list []string) {
	for i := 1; i <= length; i++ {
		list = append(list, faker.Name())
	}
	return
}

func FirstNameList(length int) (list []string) {
	for i := 1; i <= length; i++ {
		list = append(list, faker.FirstName())
	}
	return
}

func ParagraphList(length int) (list []string) {
	for i := 1; i <= length; i++ {
		list = append(list, faker.Paragraph())
	}
	return
}

func WordList(length int) (list []string) {
	for i := 1; i <= length; i++ {
		list = append(list, faker.Word())
	}
	return
}

func IntList(length int) (list []int) {
	list, _ = faker.RandomInt(length)
	return
}
