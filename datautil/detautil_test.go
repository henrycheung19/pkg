package datautil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindStrArr1(t *testing.T) {
	arr := []string{"I", "am", "an", "apple", "!"}
	fmt.Println(FindStrArr(arr, "apple"))
	assert.Equal(t, 3, FindStrArr(arr, "apple"))
}

func TestFindStrArr2(t *testing.T) {
	arr := []string{"I", "am", "an", "apple", "!"}
	assert.Equal(t, -1, FindStrArr(arr, "orange"))
}

func ExampleFindStrArr() {
	arr := []string{"I", "am", "an", "apple", "!"}
	fmt.Println(FindStrArr(arr, "apple"))
	fmt.Println(FindStrArr(arr, "orange"))
	// Output:
	// 3
	// -1
}
