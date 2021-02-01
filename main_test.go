package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetMapKeyASC(t *testing.T) {
	require := require.New(t)
	ramdonMap := map[string][]MarkDownRepo{
		"a": []MarkDownRepo{},
		"#": []MarkDownRepo{},
		"1": []MarkDownRepo{},
		"z": []MarkDownRepo{},
	}
	keys := GetMapKeyASC(ramdonMap)
	require.Equal([]string{"#", "1", "a", "z"}, keys)
}
