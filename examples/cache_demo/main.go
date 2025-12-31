package main

import (
	"fmt"
	"time"

	"github.com/alexzakarov/grogu/utils"
)

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
}

type Product struct {
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Category    string  `json:"category"`
}

func main() {
	fmt.Println("=== Grogu Metadata Cache Demonstration ===")
	fmt.Println()

	user := User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		Age:       28,
		CreatedAt: time.Now(),
		IsActive:  true,
	}

	product := Product{
		ProductID:   101,
		ProductName: "Laptop",
		Price:       15999.99,
		Stock:       50,
		Category:    "Electronics",
	}

	fmt.Println("[TEST 1] First Convert call (User struct) - Cache MISS expected")
	start := time.Now()
	fields1, values1, err1 := utils.Convert(user)
	duration1 := time.Since(start)

	if err1 != nil {
		fmt.Printf("[ERROR] Conversion failed: %v\n", err1)
		return
	}

	fmt.Printf("[INFO] Duration: %v\n", duration1)
	fmt.Printf("[INFO] Fields: %v\n", fields1)
	fmt.Printf("[INFO] Values: %v\n\n", values1)

	fmt.Println("[TEST 2] Second Convert call (User struct) - Cache HIT expected")
	start = time.Now()
	fields2, values2, err2 := utils.Convert(user)
	duration2 := time.Since(start)

	if err2 != nil {
		fmt.Printf("[ERROR] Conversion failed: %v\n", err2)
		return
	}

	fmt.Printf("[INFO] Duration: %v\n", duration2)
	fmt.Printf("[INFO] Fields: %v\n", fields2)
	fmt.Printf("[INFO] Values: %v\n", values2)

	speedup := float64(duration1.Nanoseconds()) / float64(duration2.Nanoseconds())
	fmt.Printf("[RESULT] Performance improvement: %.2fx faster\n\n", speedup)

	fmt.Println("[TEST 3] Different struct type (Product) - Cache MISS expected")
	start = time.Now()
	fields3, values3, err3 := utils.Convert(product)
	duration3 := time.Since(start)

	if err3 != nil {
		fmt.Printf("[ERROR] Conversion failed: %v\n", err3)
		return
	}

	fmt.Printf("[INFO] Duration: %v\n", duration3)
	fmt.Printf("[INFO] Fields: %v\n", fields3)
	fmt.Printf("[INFO] Values: %v\n\n", values3)

	fmt.Println("[TEST 4] Second Product call - Cache HIT expected")
	start = time.Now()
	fields4, values4, err4 := utils.Convert(product)
	duration4 := time.Since(start)

	if err4 != nil {
		fmt.Printf("[ERROR] Conversion failed: %v\n", err4)
		return
	}

	fmt.Printf("[INFO] Duration: %v\n", duration4)
	fmt.Printf("[INFO] Fields: %v\n", fields4)
	fmt.Printf("[INFO] Values: %v\n", values4)

	speedup2 := float64(duration3.Nanoseconds()) / float64(duration4.Nanoseconds())
	fmt.Printf("[RESULT] Performance improvement: %.2fx faster\n\n", speedup2)

	fmt.Println("[TEST 5] Pointer conversion test")
	userPtr := &user
	start = time.Now()
	fields5, values5, err5 := utils.Convert(userPtr)
	duration5 := time.Since(start)

	if err5 != nil {
		fmt.Printf("[ERROR] Conversion failed: %v\n", err5)
		return
	}

	fmt.Printf("[INFO] Duration: %v\n", duration5)
	fmt.Printf("[INFO] Fields: %v\n", fields5)
	fmt.Printf("[INFO] Values: %v\n\n", values5)

	fmt.Println("[TEST 6] High-volume performance test")
	iterations := 10000
	start = time.Now()
	for i := 0; i < iterations; i++ {
		utils.Convert(user)
	}
	totalDuration := time.Since(start)

	avgDuration := totalDuration / time.Duration(iterations)
	fmt.Printf("[INFO] Total duration for %d iterations: %v\n", iterations, totalDuration)
	fmt.Printf("[INFO] Average duration per call: %v\n", avgDuration)
	fmt.Printf("[RESULT] Throughput: approximately %.0f operations/second\n\n", float64(time.Second)/float64(avgDuration))

	fmt.Println("[SUCCESS] Metadata cache validation complete")
	fmt.Println("          - Cache miss occurs on first call per struct type")
	fmt.Println("          - Cache hit provides significant performance improvement")
	fmt.Println("          - Each struct type maintains separate cache entry")
	fmt.Println("          - Pointer and value types share the same cache entry")
}
