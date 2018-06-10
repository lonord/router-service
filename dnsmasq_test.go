package main

import (
	"testing"
)

func TestCollectInternalArgs(t *testing.T) {
	args := collectInternalArgs(func(path string) (string, error) {
		return dsContent, nil
	})
	if len(args) != 4 {
		t.Error("arg length dismatch")
	}
	if args[3] != "--trust-anchor=.,19036,8,2,49AAC11D7B6F6446702E54A1607371607A1A41855200FD2CE1CDDE32F24E8FB5" {
		t.Errorf("ds content dismatch [%s]", args[3])
	}
}

const dsContent = `. IN DS 19036 8 2 49AAC11D7B6F6446702E54A1607371607A1A41855200FD2CE1CDDE32F24E8FB5
`
