package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func Test_Raise2(t *testing.T) {
	c := &Computer{}
	assert.Equal(t, 8, 1<<3)
	assert.Equal(t, 1<<2, 4)
	assert.Equal(t, 8, c.pow2(3))
}

func Test_Xor3(t *testing.T) {
	c := &Computer{}
	assert.Equal(t, 2, c.xor(3, 1))
	assert.Equal(t, 4, c.xor(3, 7))
	assert.Equal(t, 1, c.xor(3, 2))
	assert.Equal(t, 3, c.xor(3, 0))
	assert.Equal(t, 7, c.xor(3, 4))
	assert.Equal(t, 5, c.xor(3, 6))
	assert.Equal(t, 0, c.xor(3, 3))
}

func Test_Output_Length(t *testing.T) {
	for i := 1; i < 16; i++ {
		c := &Computer{a: int(math.Pow(float64(7), float64(i))), b: 0, c: 0, program: []int{2, 4, 1, 3, 7, 5, 4, 1, 1, 3, 0, 3, 5, 5, 3, 0}}
		fmt.Println(c.Run())
		assert.Equal(t, i, len(c.output))
	}

}

func Test_Programs(t *testing.T) {

	tests := []struct {
		computer *Computer
		expected func(c *Computer) bool
	}{
		{
			computer: &Computer{a: 2024, b: 0, c: 0, output: []int{}, ptr: 0, program: []int{0, 3, 5, 4, 3, 0}},
			expected: func(c *Computer) bool {
				if c.String() == "5,7,3,0" {
					return true
				}
				return false
			},
		},
		{
			computer: &Computer{a: 117440, b: 0, c: 0, output: []int{}, ptr: 0, program: []int{0, 3, 5, 4, 3, 0}},
			expected: func(c *Computer) bool {
				if c.String() == "0,3,5,4,3,0" {
					return true
				}
				return false
			},
		}, {
			computer: &Computer{a: 108107566389757, b: 0, c: 0, output: []int{}, ptr: 0, program: []int{2, 4, 1, 3, 7, 5, 4, 1, 1, 3, 0, 3, 5, 5, 3, 0}},
			expected: func(c *Computer) bool {
				if c.String() == "2,4,1,3,7,5,4,1,1,3,0,3,5,5,3,0" {
					return true
				}
				return false
			},
		},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			tc.computer.Run()
			assert.True(t, tc.expected(tc.computer))
		})
	}
}
