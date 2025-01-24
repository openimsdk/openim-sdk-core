package conversation_msg

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
)

func TestMergeSortedArrays(t *testing.T) {
	array := []*model_struct.LocalChatLog{
		{SendTime: 2, Content: "Message 2"},
		{SendTime: 4, Content: "Message 4"},
		{SendTime: 6, Content: "Message 6"},
	}
	reverse(array)

	tests := []struct {
		arr1, arr2   []*model_struct.LocalChatLog
		n            int
		isDescending bool
		expected     []*model_struct.LocalChatLog
	}{
		{
			// Test merging two descending arrays
			arr1: []*model_struct.LocalChatLog{
				{SendTime: 9, Content: "Message 9"},
				{SendTime: 5, Content: "Message 5"},
				{SendTime: 3, Content: "Message 3"},
			},
			arr2: []*model_struct.LocalChatLog{
				{SendTime: 8, Content: "Message 8"},
				{SendTime: 6, Content: "Message 6"},
				{SendTime: 2, Content: "Message 2"},
			},
			n:            4, // Limit result to first 4 elements
			isDescending: true,
			expected: []*model_struct.LocalChatLog{
				{SendTime: 9, Content: "Message 9"},
				{SendTime: 8, Content: "Message 8"},
				{SendTime: 6, Content: "Message 6"},
				{SendTime: 5, Content: "Message 5"},
			},
		},
		{
			// Test merging an empty array and a descending array
			arr1: []*model_struct.LocalChatLog{},
			arr2: []*model_struct.LocalChatLog{
				{SendTime: 5, Content: "Message 5"},
				{SendTime: 3, Content: "Message 3"},
				{SendTime: 1, Content: "Message 1"},
			},
			n:            3,
			isDescending: true,
			expected: []*model_struct.LocalChatLog{
				{SendTime: 5, Content: "Message 5"},
				{SendTime: 3, Content: "Message 3"},
				{SendTime: 1, Content: "Message 1"},
			},
		},
		{
			// Test merging two empty arrays
			arr1:         []*model_struct.LocalChatLog{},
			arr2:         []*model_struct.LocalChatLog{},
			n:            0,
			isDescending: true,
			expected:     []*model_struct.LocalChatLog{},
		},
		{
			// Test merging a descending array and an ascending array
			arr1: []*model_struct.LocalChatLog{
				{SendTime: 7, Content: "Message 7"},
				{SendTime: 5, Content: "Message 5"},
				{SendTime: 3, Content: "Message 3"},
			},
			arr2:         array,
			n:            5, // Limit result to first 5 elements
			isDescending: true,
			// Expected result: merged in descending order
			expected: []*model_struct.LocalChatLog{
				{SendTime: 7, Content: "Message 7"},
				{SendTime: 6, Content: "Message 6"},
				{SendTime: 5, Content: "Message 5"},
				{SendTime: 4, Content: "Message 4"},
				{SendTime: 3, Content: "Message 3"},
			},
		},

		{
			// Test merging a descending array and an ascending array
			arr1: []*model_struct.LocalChatLog{
				{SendTime: 1, Content: "Message 1"},
				{SendTime: 5, Content: "Message 5"},
				{SendTime: 7, Content: "Message 7"},
			},
			arr2: []*model_struct.LocalChatLog{
				{SendTime: 2, Content: "Message 2"},
				{SendTime: 6, Content: "Message 6"},
				{SendTime: 9, Content: "Message 9"},
			},
			n:            5, // Limit result to first 5 elements
			isDescending: false,
			// Expected result: merged in descending order
			expected: []*model_struct.LocalChatLog{
				{SendTime: 1, Content: "Message 1"},
				{SendTime: 2, Content: "Message 2"},
				{SendTime: 5, Content: "Message 5"},
				{SendTime: 6, Content: "Message 6"},
				{SendTime: 7, Content: "Message 7"},
			},
		},
		{
			// Test merging a descending array and an ascending array
			arr1: []*model_struct.LocalChatLog{
				{SendTime: 0, Content: "Message 54", Seq: 54},
				{SendTime: 0, Content: "Message 53", Seq: 53},
				{SendTime: 0, Content: "Message 52", Seq: 52},
				{SendTime: 0, Content: "Message 4", Seq: 4},
			},
			arr2: []*model_struct.LocalChatLog{
				{SendTime: 0, Content: "Message 51", Seq: 51},
				{SendTime: 0, Content: "Message 50", Seq: 50},
				{SendTime: 0, Content: "Message 49", Seq: 49},
				{SendTime: 0, Content: "Message 48", Seq: 48},
				{SendTime: 0, Content: "Message 47", Seq: 47},
				{SendTime: 0, Content: "Message 46", Seq: 46},
				{SendTime: 0, Content: "Message 45", Seq: 45},
				{SendTime: 0, Content: "Message 44", Seq: 44},
				{SendTime: 0, Content: "Message 43", Seq: 43},
				{SendTime: 0, Content: "Message 42", Seq: 42},
				{SendTime: 0, Content: "Message 41", Seq: 41},
				{SendTime: 0, Content: "Message 40", Seq: 40},
				{SendTime: 0, Content: "Message 39", Seq: 39},
				{SendTime: 0, Content: "Message 38", Seq: 38},
				{SendTime: 0, Content: "Message 37", Seq: 37},
				{SendTime: 0, Content: "Message 36", Seq: 36},
				{SendTime: 0, Content: "Message 35", Seq: 35},
				{SendTime: 0, Content: "Message 34", Seq: 34},
				{SendTime: 0, Content: "Message 33", Seq: 33},
				{SendTime: 0, Content: "Message 32", Seq: 32},
				{SendTime: 0, Content: "Message 31", Seq: 31},
				{SendTime: 0, Content: "Message 30", Seq: 30},
				{SendTime: 0, Content: "Message 29", Seq: 29},
				{SendTime: 0, Content: "Message 28", Seq: 28},
				{SendTime: 0, Content: "Message 27", Seq: 27},
				{SendTime: 0, Content: "Message 26", Seq: 26},
				{SendTime: 0, Content: "Message 25", Seq: 25},
				{SendTime: 0, Content: "Message 24", Seq: 24},
				{SendTime: 0, Content: "Message 23", Seq: 23},
				{SendTime: 0, Content: "Message 22", Seq: 22},
				{SendTime: 0, Content: "Message 21", Seq: 21},
				{SendTime: 0, Content: "Message 20", Seq: 20},
				{SendTime: 0, Content: "Message 19", Seq: 19},
				{SendTime: 0, Content: "Message 18", Seq: 18},
				{SendTime: 0, Content: "Message 17", Seq: 17},
				{SendTime: 0, Content: "Message 16", Seq: 16},
				{SendTime: 0, Content: "Message 15", Seq: 15},
				{SendTime: 0, Content: "Message 14", Seq: 14},
				{SendTime: 0, Content: "Message 13", Seq: 13},
				{SendTime: 0, Content: "Message 12", Seq: 12},
				{SendTime: 0, Content: "Message 11", Seq: 11},
				{SendTime: 0, Content: "Message 10", Seq: 10},
				{SendTime: 0, Content: "Message 9", Seq: 9},
				{SendTime: 0, Content: "Message 8", Seq: 8},
				{SendTime: 0, Content: "Message 7", Seq: 7},
				{SendTime: 0, Content: "Message 6", Seq: 6},
				{SendTime: 0, Content: "Message 5", Seq: 5},
			},
			n:            20, // Limit result to first 5 elements
			isDescending: true,
			// Expected result: merged in descending order
			expected: []*model_struct.LocalChatLog{
				{SendTime: 0, Content: "Message 54", Seq: 54},
				{SendTime: 0, Content: "Message 53", Seq: 53},
				{SendTime: 0, Content: "Message 52", Seq: 52},
				{SendTime: 0, Content: "Message 51", Seq: 51},
				{SendTime: 0, Content: "Message 50", Seq: 50},
				{SendTime: 0, Content: "Message 49", Seq: 49},
				{SendTime: 0, Content: "Message 48", Seq: 48},
				{SendTime: 0, Content: "Message 47", Seq: 47},
				{SendTime: 0, Content: "Message 46", Seq: 46},
				{SendTime: 0, Content: "Message 45", Seq: 45},
				{SendTime: 0, Content: "Message 44", Seq: 44},
				{SendTime: 0, Content: "Message 43", Seq: 43},
				{SendTime: 0, Content: "Message 42", Seq: 42},
				{SendTime: 0, Content: "Message 41", Seq: 41},
				{SendTime: 0, Content: "Message 40", Seq: 40},
				{SendTime: 0, Content: "Message 39", Seq: 39},
				{SendTime: 0, Content: "Message 38", Seq: 38},
				{SendTime: 0, Content: "Message 37", Seq: 37},
				{SendTime: 0, Content: "Message 36", Seq: 36},
				{SendTime: 0, Content: "Message 35", Seq: 35},
			},
		},
	}

	for _, tt := range tests {
		result := mergeSortedArrays(tt.arr1, tt.arr2, tt.n, tt.isDescending)
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf(
				"mergeSortedArrays(%v, %v, %d) = %v; want %v",
				extractSendSeqs(tt.arr1),
				extractSendSeqs(tt.arr2),
				tt.n,
				extractSendSeqs(result),
				extractSendSeqs(tt.expected),
			)
		} else {
			fmt.Printf(
				"PASS: mergeSortedArrays(%v, %v, %d) = %v\n",
				extractSendSeqs(tt.arr1),
				extractSendSeqs(tt.arr2),
				tt.n,
				extractSendSeqs(result),
			)
		}
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		input    []int
		expected []int
	}{
		{input: []int{1, 2, 3, 4, 5}, expected: []int{5, 4, 3, 2, 1}},
		{input: []int{10, 20, 30}, expected: []int{30, 20, 10}},
		{input: []int{1, 2}, expected: []int{2, 1}},
		{input: []int{100}, expected: []int{100}},
		{input: []int{}, expected: []int{}},
	}

	for _, tt := range tests {

		inputCopy := make([]int, len(tt.input))
		copy(inputCopy, tt.input)

		reverse(inputCopy)
		if !reflect.DeepEqual(inputCopy, tt.expected) {
			t.Errorf("reverse(%v) = %v; want %v", tt.input, inputCopy, tt.expected)
		}
	}
}

func extractSendTimes(arr []*model_struct.LocalChatLog) []int64 {
	sendTimes := make([]int64, len(arr))
	for i, log := range arr {
		sendTimes[i] = log.SendTime
	}
	return sendTimes
}
func extractSendSeqs(arr []*model_struct.LocalChatLog) []int64 {
	sendTimes := make([]int64, len(arr))
	for i, log := range arr {
		sendTimes[i] = log.Seq
	}
	return sendTimes
}
