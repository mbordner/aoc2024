package main

import (
	"github.com/mbordner/aoc2024/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_LineMerge(t *testing.T) {
	tests := []struct {
		name     string
		lines    []Line
		expected *Line
	}{
		{
			name: "Test merging right dir adjacent lines",
			lines: []Line{
				{p1: common.Pos{X: 0, Y: 0}, p2: common.Pos{X: 1, Y: 0}},
				{p1: common.Pos{X: 1, Y: 0}, p2: common.Pos{X: 2, Y: 0}},
			},
			expected: &Line{p1: common.Pos{X: 0, Y: 0}, p2: common.Pos{X: 2, Y: 0}},
		},
		{
			name: "Test merging up dir adjacent lines",
			lines: []Line{
				{p1: common.Pos{X: 1, Y: 2}, p2: common.Pos{X: 1, Y: 1}},
				{p1: common.Pos{X: 1, Y: 1}, p2: common.Pos{X: 1, Y: 0}},
			},
			expected: &Line{p1: common.Pos{X: 1, Y: 2}, p2: common.Pos{X: 1, Y: 0}},
		},
		{
			name: "Test merging down dir adjacent lines",
			lines: []Line{
				{p1: common.Pos{X: 0, Y: 0}, p2: common.Pos{X: 0, Y: 1}},
				{p1: common.Pos{X: 0, Y: 1}, p2: common.Pos{X: 0, Y: 2}},
			},
			expected: &Line{p1: common.Pos{X: 0, Y: 0}, p2: common.Pos{X: 0, Y: 2}},
		},
		{
			name: "Test not merging right dir adjacent lines",
			lines: []Line{
				{p1: common.Pos{X: 0, Y: 0}, p2: common.Pos{X: 1, Y: 0}},
				{p1: common.Pos{X: 4, Y: 0}, p2: common.Pos{X: 5, Y: 0}},
			},
			expected: nil,
		},
		{
			name: "Test merging left dir adjacent lines",
			lines: []Line{
				{p1: common.Pos{X: 2, Y: 0}, p2: common.Pos{X: 1, Y: 0}},
				{p1: common.Pos{X: 1, Y: 0}, p2: common.Pos{X: 0, Y: 0}},
			},
			expected: &Line{p1: common.Pos{X: 2, Y: 0}, p2: common.Pos{X: 0, Y: 0}},
		},
		{
			name: "Test not merging lines in different dir",
			lines: []Line{
				{p1: common.Pos{X: 0, Y: 0}, p2: common.Pos{X: 1, Y: 0}},
				{p1: common.Pos{X: 1, Y: 0}, p2: common.Pos{X: 0, Y: 0}},
			},
			expected: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l := &tc.lines[0]
			o := &tc.lines[1]
			r := l.merge(o)
			if tc.expected == nil {
				assert.Nil(t, r)
			} else {
				assert.NotNil(t, r)
				assert.Equal(t, *tc.expected, *r)
			}
		})
	}
}
