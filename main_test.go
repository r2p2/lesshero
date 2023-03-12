// bench mark lessHero
package main

import (
	"testing"
	"time"

	"github.com/go-echarts/go-echarts/v2/opts"
)

func equalLineData(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func valuesFromLineData(line []opts.LineData) []interface{} {
	values := make([]interface{}, len(line))

	for i, l := range line {
		values[i] = l.Value
	}

	return values
}

func equalCommits(a, b []Commit) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func BenchmarkHero(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lessHero(".")
	}
}

func Test_lessHeroOrder(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCommits, _, err := lessHero(".")
			if (err != nil) != tt.wantErr {
				t.Errorf("lessHero() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// check wantCommits are in order
			for i := 0; i < len(gotCommits)-1; i++ {
				if gotCommits[i].date.After(gotCommits[i+1].date) {
					t.Errorf("lessHero() commits not in order")
				}
			}
		})
	}
}

func Test_calcRunningTotal(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name             string
		commits          []Commit
		expected_commits []Commit
	}{
		{
			name:             "no commits",
			commits:          []Commit{},
			expected_commits: []Commit{},
		},
		{
			name: "single commit",
			commits: []Commit{
				{hash: "A",
					author: "Bob",
					date:   now,
					total:  1},
			},
			expected_commits: []Commit{
				{hash: "A",
					author:       "Bob",
					date:         now,
					total:        1,
					runningTotal: 1},
			},
		},
		{
			name: "multiple commits",
			commits: []Commit{
				{hash: "A",
					author: "Bob",
					date:   now,
					total:  1},
				{hash: "A",
					author: "Alice",
					date:   now.Add(time.Second * 1),
					total:  5},
				{hash: "A",
					author: "Bob",
					date:   now.Add(time.Second * 2),
					total:  -3},
				{hash: "A",
					author: "Bob",
					date:   now.Add(time.Second * 3),
					total:  2},
			},
			expected_commits: []Commit{
				{hash: "A",
					author:       "Bob",
					date:         now,
					total:        1,
					runningTotal: 1},
				{hash: "A",
					author:       "Alice",
					date:         now.Add(time.Second * 1),
					total:        5,
					runningTotal: 6},
				{hash: "A",
					author:       "Bob",
					date:         now.Add(time.Second * 2),
					total:        -3,
					runningTotal: 3},
				{hash: "A",
					author:       "Bob",
					date:         now.Add(time.Second * 3),
					total:        2,
					runningTotal: 5},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calcRunningTotal(tt.commits)

			if !equalCommits(tt.commits, tt.expected_commits) {
				t.Errorf("calcRunningTotal() '%v' failed: %v != %v", tt.name, tt.commits, tt.expected_commits)
				return
			}
		})
	}
}

func Test_getSlocks(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name                     string
		commits                  []Commit
		expected_total_values    []interface{}
		expected_decrease_values []interface{}
	}{
		{
			name:                     "no commits",
			commits:                  []Commit{},
			expected_total_values:    []interface{}{},
			expected_decrease_values: []interface{}{},
		},
		{
			name: "multiple commits",
			commits: []Commit{
				{hash: "A",
					author: "Bob",
					date:   now,
					total:  0},
				{hash: "A",
					author: "Alice",
					date:   now,
					total:  0},
				{hash: "A",
					author: "Bob",
					date:   now,
					total:  0},
				{hash: "A",
					author: "Bob",
					date:   now,
					total:  1},
				{hash: "A",
					author: "Bob",
					date:   now,
					total:  4},
				{hash: "A",
					author: "Bob",
					date:   now,
					total:  1},
				{hash: "A",
					author: "Bob",
					date:   now,
					total:  0},
				{hash: "A",
					author: "Alice",
					date:   now,
					total:  0},
				{hash: "A",
					author: "Bob",
					date:   now,
					total:  0},
				{hash: "A",
					author: "Bob",
					date:   now,
					total:  -3},
				{hash: "A",
					author: "Bob",
					date:   now,
					total:  -2},
				{hash: "A",
					author: "Bob",
					date:   now,
					total:  -5},
				{hash: "A",
					author: "Alice",
					date:   now,
					total:  0},
				{hash: "A",
					author: "Bob",
					date:   now,
					total:  0},
				{hash: "A",
					author: "Bob",
					date:   now,
					total:  0},
			},
			expected_total_values:    []interface{}{0, 0, 0, 1, 5, 6, 6, 6, 6, 3, 1, -4, -4, -4, -4},
			expected_decrease_values: []interface{}{"-", "-", "-", "-", "-", "-", "-", "-", 6, 3, 1, -4, "-", "-", "-"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calcRunningTotal(tt.commits)
			line_total, line_decrease := getSlocs(tt.commits)
			values_total := valuesFromLineData(line_total)
			values_decrease := valuesFromLineData(line_decrease)

			if !equalLineData(values_total, tt.expected_total_values) {
				t.Errorf("getSlocs() line_total %v != %v", values_total, tt.expected_total_values)
				return
			}
			if !equalLineData(values_decrease, tt.expected_decrease_values) {
				t.Errorf("getSlocs() line_decrease %v != %v", values_decrease, tt.expected_decrease_values)
				return
			}
		})
	}
}
