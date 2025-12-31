package utils

import (
	"testing"
)

type TestStruct struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
	Country   string `json:"country"`
	City      string `json:"city"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	Status    bool   `json:"status"`
	CreatedAt string `json:"created_at"`
}

type ComplexStruct struct {
	ID          int64   `json:"id"`
	Username    string  `json:"username"`
	Email       string  `json:"email"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Age         int     `json:"age"`
	Height      float64 `json:"height"`
	Weight      float64 `json:"weight"`
	Country     string  `json:"country"`
	City        string  `json:"city"`
	Address     string  `json:"address"`
	PostalCode  string  `json:"postal_code"`
	Phone       string  `json:"phone"`
	MobilePhone string  `json:"mobile_phone"`
	Status      int     `json:"status"`
	IsActive    bool    `json:"is_active"`
	IsVerified  bool    `json:"is_verified"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	LastLogin   string  `json:"last_login"`
}

func BenchmarkConvertSimpleStruct(b *testing.B) {
	testData := TestStruct{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		Age:       30,
		Country:   "USA",
		City:      "New York",
		Address:   "123 Main St",
		Phone:     "+1234567890",
		Status:    true,
		CreatedAt: "2024-01-01",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = Convert(testData)
	}
}

func BenchmarkConvertComplexStruct(b *testing.B) {
	testData := ComplexStruct{
		ID:          1,
		Username:    "johndoe",
		Email:       "john@example.com",
		FirstName:   "John",
		LastName:    "Doe",
		Age:         30,
		Height:      180.5,
		Weight:      75.5,
		Country:     "USA",
		City:        "New York",
		Address:     "123 Main St",
		PostalCode:  "10001",
		Phone:       "+1234567890",
		MobilePhone: "+0987654321",
		Status:      1,
		IsActive:    true,
		IsVerified:  true,
		CreatedAt:   "2024-01-01",
		UpdatedAt:   "2024-01-02",
		LastLogin:   "2024-01-03",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = Convert(testData)
	}
}

func BenchmarkConvertPointer(b *testing.B) {
	testData := &TestStruct{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		Age:       30,
		Country:   "USA",
		City:      "New York",
		Address:   "123 Main St",
		Phone:     "+1234567890",
		Status:    true,
		CreatedAt: "2024-01-01",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = Convert(testData)
	}
}

func BenchmarkConvertParallel(b *testing.B) {
	testData := ComplexStruct{
		ID:          1,
		Username:    "johndoe",
		Email:       "john@example.com",
		FirstName:   "John",
		LastName:    "Doe",
		Age:         30,
		Height:      180.5,
		Weight:      75.5,
		Country:     "USA",
		City:        "New York",
		Address:     "123 Main St",
		PostalCode:  "10001",
		Phone:       "+1234567890",
		MobilePhone: "+0987654321",
		Status:      1,
		IsActive:    true,
		IsVerified:  true,
		CreatedAt:   "2024-01-01",
		UpdatedAt:   "2024-01-02",
		LastLogin:   "2024-01-03",
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _, _ = Convert(testData)
		}
	})
}

func TestConvertCorrectness(t *testing.T) {
	testData := TestStruct{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		Age:       30,
		Country:   "USA",
		City:      "New York",
		Address:   "123 Main St",
		Phone:     "+1234567890",
		Status:    true,
		CreatedAt: "2024-01-01",
	}

	jsonNames, values, err := Convert(testData)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	expectedFields := []string{"id", "name", "email", "age", "country", "city", "address", "phone", "status", "created_at"}
	if len(jsonNames) != len(expectedFields) {
		t.Errorf("Expected %d fields, got %d", len(expectedFields), len(jsonNames))
	}

	for i, name := range expectedFields {
		if jsonNames[i] != name {
			t.Errorf("Expected field name %s at index %d, got %s", name, i, jsonNames[i])
		}
	}

	if len(values) != len(expectedFields) {
		t.Errorf("Expected %d values, got %d", len(expectedFields), len(values))
	}
}

func TestConvertPointer(t *testing.T) {
	testData := &TestStruct{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	jsonNames, values, err := Convert(testData)
	if err != nil {
		t.Fatalf("Convert failed with pointer: %v", err)
	}

	if len(jsonNames) == 0 {
		t.Error("Expected non-empty field names")
	}

	if len(values) == 0 {
		t.Error("Expected non-empty values")
	}
}

func TestConvertError(t *testing.T) {
	testData := "not a struct"

	_, _, err := Convert(testData)
	if err == nil {
		t.Error("Expected error for non-struct input")
	}
}

func TestCacheConsistency(t *testing.T) {
	testData := TestStruct{
		ID:   1,
		Name: "John",
	}

	jsonNames1, values1, _ := Convert(testData)
	jsonNames2, values2, _ := Convert(testData)

	if len(jsonNames1) != len(jsonNames2) {
		t.Error("Inconsistent field names between calls")
	}

	for i := range jsonNames1 {
		if jsonNames1[i] != jsonNames2[i] {
			t.Errorf("Field name mismatch at index %d: %s != %s", i, jsonNames1[i], jsonNames2[i])
		}
	}

	if len(values1) != len(values2) {
		t.Error("Inconsistent value count between calls")
	}
}
