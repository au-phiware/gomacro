/*
 * gomacro - A Go interpreter with Lisp-like macros
 *
 * Copyright (C) 2017-2018 Massimiliano Ghilardi
 *
 *     This program is free software: you can redistribute it and/or modify
 *     it under the terms of the GNU Lesser General Public License as published
 *     by the Free Software Foundation, either version 3 of the License, or
 *     (at your option) any later version.
 *
 *     This program is distributed in the hope that it will be useful,
 *     but WITHOUT ANY WARRANTY; without even the implied warranty of
 *     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *     GNU Lesser General Public License for more details.
 *
 *     You should have received a copy of the GNU Lesser General Public License
 *     along with this program.  If not, see <https://www.gnu.org/licenses/lgpl>.
 *
 *
 * untyped.go
 *
 *  Created on May 27, 2017
 *      Author Massimiliano Ghilardi
 */

package base

import (
	"fmt"
	"go/constant"
	"go/token"
	"go/types"
	"math/big"
	"reflect"
	"strings"
)

type UntypedVal struct {
	Kind reflect.Kind
	Val  constant.Value
}

func UntypedKindToReflectKind(gkind types.BasicKind) reflect.Kind {
	var kind reflect.Kind
	switch gkind {
	case types.UntypedBool:
		kind = reflect.Bool
	case types.UntypedInt:
		kind = reflect.Int
	case types.UntypedRune:
		kind = reflect.Int32
	case types.UntypedFloat:
		kind = reflect.Float64
	case types.UntypedComplex:
		kind = reflect.Complex128
	case types.UntypedString:
		kind = reflect.String
	case types.UntypedNil:
		kind = reflect.Invalid
	default:
		Errorf("unsupported types.BasicKind: %v", gkind)
	}
	return kind
}

func MarshalUntyped(kind types.BasicKind, val constant.Value) string {
	rkind := UntypedKindToReflectKind(kind)
	lit := UntypedVal{rkind, val}
	return lit.Marshal()
}

func UnmarshalUntyped(marshalled string) (reflect.Kind, constant.Value) {
	lit := UnmarshalUntypedLit(marshalled)
	return lit.Kind, lit.Val
}

func (lit *UntypedVal) Marshal() string {
	// untyped constants have arbitrary precision... they may overflow integers
	val := lit.Val
	var s string
	switch lit.Kind {
	case reflect.Invalid:
		s = "nil"
	case reflect.Bool:
		if constant.BoolVal(val) {
			s = "bool:true"
		} else {
			s = "bool:false"
		}
	case reflect.Int:
		s = fmt.Sprintf("int:%s", val.ExactString())
	case reflect.Int32:
		s = fmt.Sprintf("rune:%s", val.ExactString())
	case reflect.Float64:
		s = fmt.Sprintf("float:%s", val.ExactString())
	case reflect.Complex128:
		s = fmt.Sprintf("complex:%s:%s", constant.Real(val).ExactString(), constant.Imag(val).ExactString())
	case reflect.String:
		s = fmt.Sprintf("string:%s", constant.StringVal(val))
	}
	return s
}

func UnmarshalUntypedLit(marshalled string) *UntypedVal {
	var skind, str string
	if sep := strings.IndexByte(marshalled, ':'); sep >= 0 {
		skind = marshalled[:sep]
		str = marshalled[sep+1:]
	} else {
		skind = marshalled
	}

	var kind reflect.Kind
	var val constant.Value
	switch skind {
	case "bool":
		kind = reflect.Bool
		if str == "true" {
			val = constant.MakeBool(true)
		} else {
			val = constant.MakeBool(false)
		}
	case "int":
		kind = reflect.Int
		val = constant.MakeFromLiteral(str, token.INT, 0)
	case "rune":
		kind = reflect.Int32
		val = constant.MakeFromLiteral(str, token.INT, 0)
	case "float":
		kind = reflect.Float64
		val = unmarshalUntypedFloat(str)
	case "complex":
		kind = reflect.Complex128
		if sep := strings.IndexByte(str, ':'); sep >= 0 {
			re := unmarshalUntypedFloat(str[:sep])
			im := unmarshalUntypedFloat(str[sep+1:])
			val = constant.BinaryOp(constant.ToComplex(re), token.ADD, constant.MakeImag(im))
		} else {
			val = constant.ToComplex(unmarshalUntypedFloat(str))
		}
	case "string":
		kind = reflect.String
		val = constant.MakeString(str)
	case "nil":
		kind = reflect.Invalid
	default:
		kind = reflect.Invalid
	}
	return &UntypedVal{kind, val}
}

// generalization of constant.MakeFromLiteral, accepts the fractions generated by
// constant.Value.ExactString() for floating-point values
func unmarshalUntypedFloat(str string) constant.Value {
	if sep := strings.IndexByte(str, '/'); sep >= 0 {
		x := constant.MakeFromLiteral(str[:sep], token.FLOAT, 0)
		y := constant.MakeFromLiteral(str[sep+1:], token.FLOAT, 0)
		return constant.BinaryOp(x, token.QUO, y)
	}
	return constant.MakeFromLiteral(str, token.FLOAT, 0)
}

func (lit *UntypedVal) BigInt() (*big.Int, error) {
	val := lit.Val
	switch lit.Kind {
	case reflect.Int, reflect.Int32:
		if i, ok := constant.Int64Val(val); ok {
			return big.NewInt(i), nil
		}
		if bi, ok := new(big.Int).SetString(val.ExactString(), 10); ok {
			return bi, nil
		}
	}
	return nil, makeRuntimeError("cannot convert untyped %s to math/big.Int: %v", lit.Kind, lit.Val)
}

func (lit *UntypedVal) BigRat() (*big.Rat, error) {
	val := lit.Val
	switch lit.Kind {
	case reflect.Int, reflect.Int32:
		if i, ok := constant.Int64Val(val); ok {
			return big.NewRat(i, 1), nil
		}
		fallthrough
	case reflect.Float64:
		if br, ok := new(big.Rat).SetString(val.ExactString()); ok {
			return br, nil
		}
	}
	return nil, makeRuntimeError("cannot convert untyped %s to math/big.Rat: %v", lit.Kind, lit.Val)
}