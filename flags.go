package main

import (
	"encoding/hex"
	"errors"
	"strconv"
)

type rootFlag []byte

func (r *rootFlag) Set(s string) (err error) {
	root, err := hex.DecodeString(s)
	if err != nil {
		return
	}
	if len(root) != 32 {
		return errors.New("root must be a 32-byte hex string")
	}
	*r = root
	return
}

func (r rootFlag) String() string {
	return hex.EncodeToString(r)
}

type difficultyFlag struct{ n uint64 }

func (d *difficultyFlag) Set(s string) (err error) {
	target, err := strconv.ParseUint(s, 16, 0)
	if err == nil {
		d.n = target
	}
	return
}

func (d difficultyFlag) String() string {
	return strconv.FormatUint(d.n, 16)
}
