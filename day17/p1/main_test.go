package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Raise2(t *testing.T) {
	c := &Computer{}
	assert.Equal(t, 8, 1<<3)
	assert.Equal(t, 1<<2, 4)
	assert.Equal(t, 8, c.pow2(3))
}

func Test_Programs(t *testing.T) {

	tests := []struct {
		computer *Computer
		expected func(c *Computer) bool
	}{
		{
			computer: &Computer{a: 8, b: 0, c: 0, output: []int{}, ptr: 0, program: []int{0, 3, 5, 4}},
			expected: func(c *Computer) bool {
				if c.a == 1 && c.String() == "1" {
					return true
				}
				return false
			},
		},
		{
			computer: &Computer{a: 8, b: 0, c: 0, output: []int{}, ptr: 0, program: []int{6, 3, 5, 5}},
			expected: func(c *Computer) bool {
				if c.b == 1 && c.String() == "1" {
					return true
				}
				return false
			},
		},
		{
			computer: &Computer{a: 8, b: 0, c: 0, output: []int{}, ptr: 0, program: []int{7, 3, 5, 6}},
			expected: func(c *Computer) bool {
				if c.c == 1 && c.String() == "1" {
					return true
				}
				return false
			},
		},
		{
			computer: &Computer{a: 0, b: 0, c: 9, output: []int{}, ptr: 0, program: []int{2, 6}},
			expected: func(c *Computer) bool {
				if c.b == 1 {
					return true
				}
				return false
			},
		},
		{
			computer: &Computer{a: 10, b: 0, c: 0, output: []int{}, ptr: 0, program: []int{5, 0, 5, 1, 5, 4}},
			expected: func(c *Computer) bool {
				if c.String() == "0,1,2" {
					return true
				}
				return false
			},
		},
		{
			computer: &Computer{a: 0, b: 29, c: 0, output: []int{}, ptr: 0, program: []int{1, 7}},
			expected: func(c *Computer) bool {
				if c.b == 26 {
					return true
				}
				return false
			},
		},
		{
			computer: &Computer{a: 0, b: 2024, c: 43690, output: []int{}, ptr: 0, program: []int{4, 0}},
			expected: func(c *Computer) bool {
				if c.b == 44354 {
					return true
				}
				return false
			},
		},
		{
			computer: &Computer{a: 2024, b: 0, c: 0, output: []int{}, ptr: 0, program: []int{0, 1, 5, 4, 3, 0}},
			expected: func(c *Computer) bool {
				if c.String() == "4,2,5,6,7,7,7,7,3,1,0" && c.a == 0 {
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
