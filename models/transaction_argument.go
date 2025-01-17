package models

import (
	"encoding/hex"
	"math/big"
	"strconv"

	"github.com/the729/lcs"
)

type TransactionArgument interface {
	ToString() string
}

type TxArgU8 struct {
	U8 uint8
}

type TxArgU64 struct {
	U64 uint64
}

type TxArgU128 struct {
	Higher, Lower uint64
}

type TxArgAddress struct {
	Addr AccountAddress
}

type TxArgU8Vector struct {
	Bytes []byte
}

type TxArgBool struct {
	Bool bool
}

var _ = lcs.RegisterEnum(
	(*TransactionArgument)(nil),
	TxArgU8{},
	TxArgU64{},
	TxArgU128{},
	TxArgAddress{},
	TxArgU8Vector{},
	TxArgBool{},
)

func (t TxArgU8) ToString() string {
	return strconv.FormatUint(uint64(t.U8), 10)
}

func (t TxArgU64) ToString() string {
	return strconv.FormatUint(t.U64, 10)
}

func (t TxArgU128) ToString() string {
	higherBytes := new(big.Int).SetUint64(t.Higher).Bytes()
	lowerBytes := new(big.Int).SetUint64(t.Lower).Bytes()
	return new(big.Int).SetBytes(append(higherBytes, lowerBytes...)).String()
}

func (t TxArgAddress) ToString() string {
	return hex.EncodeToString(t.Addr[:])
}

func (t TxArgU8Vector) ToString() string {
	return hex.EncodeToString(t.Bytes)
}

func (t TxArgBool) ToString() string {
	if t.Bool {
		return "true"
	}
	return "false"
}
