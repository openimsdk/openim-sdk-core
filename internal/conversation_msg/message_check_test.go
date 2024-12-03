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
	}

	for _, tt := range tests {
		result := mergeSortedArrays(tt.arr1, tt.arr2, tt.n, tt.isDescending)
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf(
				"mergeSortedArrays(%v, %v, %d) = %v; want %v",
				extractSendTimes(tt.arr1),
				extractSendTimes(tt.arr2),
				tt.n,
				extractSendTimes(result),
				extractSendTimes(tt.expected),
			)
		} else {
			fmt.Printf(
				"PASS: mergeSortedArrays(%v, %v, %d) = %v\n",
				extractSendTimes(tt.arr1),
				extractSendTimes(tt.arr2),
				tt.n,
				extractSendTimes(result),
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
