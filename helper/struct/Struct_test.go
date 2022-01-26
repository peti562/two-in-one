package helper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type simpleStruct struct {
	ID int `gorm:"column:id"`
}

type currencyStruct struct {
	ID             int       `gorm:"column:cu_id;primary_key" json:"id"`
	Name           string    `gorm:"column:cu_name" json:"name"`
	ShortName      string    `gorm:"column:cu_short_name" json:"shortName"`
	Symbol         string    `gorm:"column:cu_symbol" json:"symbol"`
	IsoCode        int       `gorm:"column:cu_iso_code" json:"isoCode"`
	ExchangeRate   float64   `gorm:"column:cu_exchange_rate" json:"exchangeRate"`
	ReverseRate    float64   `gorm:"column:cu_reverse_rate" json:"reverseRate"`
	RoundingFactor float64   `gorm:"column:cu_rounding_factor" json:"roundingFactor"`
	UpdatedDate    time.Time `gorm:"column:cu_updated_date" json:"updatedDate"`
	Status         bool      `gorm:"column:cu_status" json:"status"`
}

// @todo Test float32 :(

func TestMapAsGorm(t *testing.T) {

	timestamp := time.Now()

	tests := []struct {
		name   string
		values interface{}
		want   map[string]interface{}
	}{
		{
			name: "Simple struct",
			values: &simpleStruct{
				ID: 1,
			},
			want: map[string]interface{}{
				"id": 1,
			},
		},
		{
			name: "Currency struct",
			values: &currencyStruct{
				ID:             1,
				Name:           "Test",
				ShortName:      "TST",
				Symbol:         "$",
				IsoCode:        100,
				ExchangeRate:   1.05,
				ReverseRate:    0.85,
				RoundingFactor: 0.05,
				UpdatedDate:    timestamp,
				Status:         true,
			},
			want: map[string]interface{}{
				"cu_id":              1,
				"cu_name":            "Test",
				"cu_short_name":      "TST",
				"cu_symbol":          "$",
				"cu_iso_code":        100,
				"cu_exchange_rate":   1.05,
				"cu_reverse_rate":    0.85,
				"cu_rounding_factor": 0.05,
				"cu_updated_date":    timestamp,
				"cu_status":          true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapAsGorm(tt.values)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestWithPrefix(t *testing.T) {

	prefix := "test_"

	// Create an option interface
	argOption := WithPrefix(prefix)

	assert.Equal(t, prefix, argOption.GetValue())
	assert.Equal(t, "prefix", argOption.GetType())
	assert.Implements(t, (*OptionInterface)(nil), argOption)
}

func TestCombine(t *testing.T) {

	type testStruct struct {
		ID  int    `gorm:"column:id"`
		Key string `gorm:"column:key"`
	}

	// Combine the structs together
	data := Combine(
		MapAsGorm(&testStruct{
			ID:  1,
			Key: "test",
		}),
		MapAsGorm(&testStruct{
			ID:  2,
			Key: "testing",
		}, WithPrefix("b_")),
	)

	// Our second struct should have a prefix
	wantData := map[string]interface{}{
		"id":    1,
		"key":   "test",
		"b_id":  2,
		"b_key": "testing",
	}

	assert.Equal(t, wantData, data)
}
