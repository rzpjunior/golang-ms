package reportx

import (
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
)

type User struct {
	ID        int64
	Name      string
	IsEnable  bool
	Status    int8
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GenerateData() (users []User) {
	for i := 1; i < 100; i++ {
		users = append(users, User{
			ID:        int64(i),
			Name:      faker.FirstName(),
			IsEnable:  true,
			Status:    2,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}
	return
}

func Test_GenerateXlsx(t *testing.T) {

	var ex Excelx
	header := []string{
		"ID",
		"Name",
		"IsEnable",
		"Status",
		"Created At",
		"Updated At",
	}

	data := GenerateData()

	var cells []interface{}
	for _, item := range data {
		cells = append(cells, item)
	}

	ex.Sheets = append(ex.Sheets, Sheet{
		WithNumbering: true,
		Name:          "Sheet1",
		Headers:       header,
		Bodys:         cells,
	})

	dir, err := GenerateXlsx("sample.xlsx", ex)

	assert.Nil(t, err)
	assert.NotNil(t, dir)
}
