package main

import (
	"testing"
	// . "github.com/onsi/ginkgo"
	// . "github.com/onsi/ginkgo/extensions/table"
	// . "github.com/onsi/gomega"
)

// Approach-1 with go testing library
// func TestToReadableSize(t *testing.T) {

// 	got := toReadableSize(4000)
// 	want := "4 KB"

// 	if got != want {
// 		t.Errorf("got %q, wanted %q", got, want)
// 	}
// }

// // Approach-2 with go testing library with multiple multiple tests scenerios(Table Driven Test)
// type toReadableSizeTest struct {
// 	nbytes   int
// 	expected string
// }

// var toReadableSizeTests = []toReadableSizeTest{
// 	{1000, "1000 B"},
// 	{1000 * 1000, "1000 KB"},
// 	{1000 * 1000 * 1000, "1000 MB"},
// 	{1000 * 1000 * 1000 * 1000, "1000 GB"},
// }

// func TestToReadableSizeMultiple(t *testing.T) {

// 	for _, test := range toReadableSizeTests {
// 		if output := toReadableSize(int64(test.nbytes)); output != test.expected {
// 			t.Errorf("Output %q not equal to expected %q", output, test.expected)
// 		}
// 	}
// }

// Approach-3 with go testing library with multiple multiple tests scenerios and name for each scenerio(Table Driven Test)
func TestToReadableSizeMultipleApproachTwo(t *testing.T) {
	tt := []struct {
		name     string
		input    int64
		expected string
	}{
		{"byte_return", 125, "125 B"},
		{"kilobyte_return", 1010, "1 KB"},
		{"megabyte_return", 1988909, "1 MB"},
		{"gigabyte_return", 29121988909, "29 GB"},
		{"gigabyte_return", 890929121988909, "890 TB"},
		{"kilobyte_return2", 1000, "1000 B"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output := toReadableSize(tc.input)
			if output != tc.expected {
				t.Errorf("input %d, unexpected output: %s", tc.input, output)
			}
		})
	}
}

// func Test_traverseDir() {

// }

func Test_calculateHash(t *testing.T) {
	tt := []struct {
		name     string
		input    []byte
		expected string
	}{
		{"hashOfNumber", []byte{124}, "3eb416223e9e69e6bb8ee19793911ad1ad2027d8"},
		{"hashOfString", []byte("Hello"), "f7ff9e8b7bb2e09b70935a5d785e0cc5d9d0abf0"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			output := CalculateHash(tc.input)
			if output != tc.expected {
				t.Errorf("input %d, unexpected output: %s", tc.input, output)
			}
		})
	}
}

// Approach-4 with ginko ang gomega
// func TestToReadableSizeWithGinkoAndGomega(t *testing.T) {
// 	RegisterFailHandler(Fail)
// 	RunSpecs(t, "To readable size")
// }

// var _ = Describe("Main", func() {
// 	Context("To readable size", func() {
// 		It("should return in bytes if nbytes is less than 1000", func() {
// 			result := toReadableSize(125)
// 			Expect(result).To(Equal("125 B"))
// 		})

// 		It("should return in kilo bytes if nbytes is less than MB", func() {
// 			result := toReadableSize(10000)
// 			Expect(result).To(Equal("10 KB"))
// 		})
// 	})
// })

// // Approach-5 with ginko ang gomega(Table Driven Test)
// var _ = Describe("Main", func() {
// 	Context("To readable size", func() {

// 		DescribeTable("To readable size", func(nbytes int64, expectedResult string) {
// 			actualResult := toReadableSize(nbytes)
// 			Expect(actualResult).To(Equal(expectedResult))
// 		},
// 			Entry("should return in bytes if nbytes is less than 1000", int64(125), "125 B"),
// 			Entry("should return in kilo bytes if nbytes is less than MB", int64(10000), "10 KB"),
// 		)
// 	})

// })
