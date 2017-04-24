/*
 * gomacro - A Go interpreter with Lisp-like macros
 *
 * Copyright (C) 2017 Massimiliano Ghilardi
 *
 *     This program is free software you can redistribute it and/or modify
 *     it under the terms of the GNU General Public License as published by
 *     the Free Software Foundation, either version 3 of the License, or
 *     (at your option) any later version.
 *
 *     This program is distributed in the hope that it will be useful,
 *     but WITHOUT ANY WARRANTY; without even the implied warranty of
 *     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *     GNU General Public License for more details.
 *
 *     You should have received a copy of the GNU General Public License
 *     along with this program.  If not, see <http//www.gnu.org/licenses/>.
 *
 * unary_plus.go
 *
 *  Created on Apr 07, 2017
 *      Author Massimiliano Ghilardi
 */

package fast

import (
	"go/ast"
	r "reflect"

	. "github.com/cosmos72/gomacro/base"
)

func (c *Comp) UnaryPlus(node *ast.UnaryExpr, xe *Expr) *Expr {
	if !IsCategory(xe.Type.Kind(), r.Int, r.Uint, r.Float64, r.Complex128) {
		return c.invalidUnaryExpr(node, xe)
	}
	return xe
}

func (c *Comp) UnaryMinus(node *ast.UnaryExpr, xe *Expr) *Expr {
	// if xe is constant, UnaryExpr will invoke EvalConst()
	// on our return value. no need to optimize that.
	x := xe.Fun
	var fun I
	switch x := x.(type) {
	case func(env *Env) int:
		fun = func(env *Env) int {
			return -x(env)
		}
	case func(env *Env) int8:
		fun = func(env *Env) int8 {
			return -x(env)
		}
	case func(env *Env) int16:
		fun = func(env *Env) int16 {
			return -x(env)
		}
	case func(env *Env) int32:
		fun = func(env *Env) int32 {
			return -x(env)
		}
	case func(env *Env) int64:
		fun = func(env *Env) int64 {
			return -x(env)
		}
	case func(env *Env) uint:
		fun = func(env *Env) uint {
			return -x(env)
		}
	case func(env *Env) uint8:
		fun = func(env *Env) uint8 {
			return -x(env)
		}
	case func(env *Env) uint16:
		fun = func(env *Env) uint16 {
			return -x(env)
		}
	case func(env *Env) uint32:
		fun = func(env *Env) uint32 {
			return -x(env)
		}
	case func(env *Env) uint64:
		fun = func(env *Env) uint64 {
			return -x(env)
		}
	case func(env *Env) uintptr:
		fun = func(env *Env) uintptr {
			return -x(env)
		}
	case func(env *Env) float32:
		fun = func(env *Env) float32 {
			return -x(env)
		}
	case func(env *Env) float64:
		fun = func(env *Env) float64 {
			return -x(env)
		}
	case func(env *Env) complex64:
		fun = func(env *Env) complex64 {
			return -x(env)
		}
	case func(env *Env) complex128:
		fun = func(env *Env) complex128 {
			return -x(env)
		}
	default:
		return c.invalidUnaryExpr(node, xe)
	}
	return exprFun(xe.Type, fun)
}

func (c *Comp) UnaryXor(node *ast.UnaryExpr, xe *Expr) *Expr {
	// if xe is constant, UnaryExpr will invoke EvalConst()
	// on our return value. no need to optimize that.
	x := xe.Fun
	var fun I
	switch x := x.(type) {
	case func(env *Env) int:
		fun = func(env *Env) int {
			return ^x(env)
		}
	case func(env *Env) int8:
		fun = func(env *Env) int8 {
			return ^x(env)
		}
	case func(env *Env) int16:
		fun = func(env *Env) int16 {
			return ^x(env)
		}
	case func(env *Env) int32:
		fun = func(env *Env) int32 {
			return ^x(env)
		}
	case func(env *Env) int64:
		fun = func(env *Env) int64 {
			return ^x(env)
		}
	case func(env *Env) uint:
		fun = func(env *Env) uint {
			return ^x(env)
		}
	case func(env *Env) uint8:
		fun = func(env *Env) uint8 {
			return ^x(env)
		}
	case func(env *Env) uint16:
		fun = func(env *Env) uint16 {
			return ^x(env)
		}
	case func(env *Env) uint32:
		fun = func(env *Env) uint32 {
			return ^x(env)
		}
	case func(env *Env) uint64:
		fun = func(env *Env) uint64 {
			return ^x(env)
		}
	case func(env *Env) uintptr:
		fun = func(env *Env) uintptr {
			return ^x(env)
		}
	default:
		return c.invalidUnaryExpr(node, xe)
	}
	return exprFun(xe.Type, fun)
}

func (c *Comp) UnaryNot(node *ast.UnaryExpr, xe *Expr) *Expr {
	// if xe is constant, UnaryExpr will invoke EvalConst()
	// on our return value. no need to optimize that.
	x := xe.Fun
	var fun I
	switch x := x.(type) {
	case func(env *Env) bool:
		fun = func(env *Env) bool {
			return !x(env)
		}
	default:
		return c.invalidUnaryExpr(node, xe)
	}
	return exprFun(xe.Type, fun)
}

func (c *Comp) UnaryRecv(node *ast.UnaryExpr, xe *Expr) *Expr {
	t := xe.Type
	if t.Kind() != r.Chan {
		return c.badUnaryExpr("expecting channel, found", node, xe)
	}
	var fun func(env *Env) (r.Value, []r.Value)
	switch x := xe.Fun.(type) {
	case func(env *Env) (r.Value, []r.Value):
		fun = func(env *Env) (r.Value, []r.Value) {
			channel, _ := x(env)
			ret, ok := channel.Recv()
			return ret, []r.Value{ret, r.ValueOf(ok)}
		}
	default:
		x1 := xe.AsX1()
		fun = func(env *Env) (r.Value, []r.Value) {
			ret, ok := x1(env).Recv()
			return ret, []r.Value{ret, r.ValueOf(ok)}
		}
	}
	types := []r.Type{t.Elem(), TypeOfBool}
	return exprXV(types, fun)
}

// StarExpr compiles unary operator * i.e. pointer dereference
func (c *Comp) StarExpr(node *ast.StarExpr) *Expr {
	addr := c.Expr1(node.X) // panics if addr returns zero values, warns if returns multiple values
	taddr := addr.Type
	if taddr.Kind() != r.Ptr {
		c.Errorf("unary operation * on non-pointer <%v>: %v", taddr, node)
	}
	return c.Deref(addr)
}

// Deref compiles unary operator * i.e. pointer dereference
func (c *Comp) Deref(addr *Expr) *Expr {
	taddr := addr.Type
	if taddr.Kind() != r.Ptr {
		c.Errorf("unary operation * on non-pointer <%v>", taddr)
	}
	x1 := addr.AsX1() // panics if addr returns zero values, warns if returns multiple values
	t := taddr.Elem()
	x := addr.Fun
	var fun I
	// fast interpreter expects that Exprs returning primitive types or string
	// do NOT wrap them into reflect.Value
	switch x := x.(type) {
	case func(env *Env) *bool:
		fun = func(env *Env) bool {
			return *x(env)
		}
	case func(env *Env) *int:
		fun = func(env *Env) int {
			return *x(env)
		}
	case func(env *Env) *int8:
		fun = func(env *Env) int8 {
			return *x(env)
		}
	case func(env *Env) *int16:
		fun = func(env *Env) int16 {
			return *x(env)
		}
	case func(env *Env) *int32:
		fun = func(env *Env) int32 {
			return *x(env)
		}
	case func(env *Env) *int64:
		fun = func(env *Env) int64 {
			return *x(env)
		}
	case func(env *Env) *uint:
		fun = func(env *Env) uint {
			return *x(env)
		}
	case func(env *Env) *uint8:
		fun = func(env *Env) uint8 {
			return *x(env)
		}
	case func(env *Env) *uint16:
		fun = func(env *Env) uint16 {
			return *x(env)
		}
	case func(env *Env) *uint32:
		fun = func(env *Env) uint32 {
			return *x(env)
		}
	case func(env *Env) *uint64:
		fun = func(env *Env) uint64 {
			return *x(env)
		}
	case func(env *Env) *uintptr:
		fun = func(env *Env) uintptr {
			return *x(env)
		}
	case func(env *Env) *float32:
		fun = func(env *Env) float32 {
			return *x(env)
		}
	case func(env *Env) *float64:
		fun = func(env *Env) float64 {
			return *x(env)
		}
	case func(env *Env) *complex64:
		fun = func(env *Env) complex64 {
			return *x(env)
		}
	default:
		fun = c.derefUnwrap(t, x1)
	}
	return exprFun(t, fun)
}

// deref0Unwrap compiles unary operator * on reflect.Value - unwraps reflect.Value.Elem() if possible
func (c *Comp) derefUnwrap(t r.Type, x1 func(*Env) r.Value) I {
	var fun I
	switch t.Kind() {
	case r.Bool:
		fun = func(env *Env) bool {
			return x1(env).Elem().Bool()
		}
	case r.Int:
		fun = func(env *Env) int {
			return int(x1(env).Elem().Int())
		}
	case r.Int8:
		fun = func(env *Env) int8 {
			return int8(x1(env).Elem().Int())
		}
	case r.Int16:
		fun = func(env *Env) int16 {
			return int16(x1(env).Elem().Int())
		}
	case r.Int32:
		fun = func(env *Env) int32 {
			return int32(x1(env).Elem().Int())
		}
	case r.Int64:
		fun = func(env *Env) int64 {
			return x1(env).Elem().Int()
		}
	case r.Uint:
		fun = func(env *Env) uint {
			return uint(x1(env).Elem().Uint())
		}
	case r.Uint8:
		fun = func(env *Env) uint8 {
			return uint8(x1(env).Elem().Uint())
		}
	case r.Uint16:
		fun = func(env *Env) uint16 {
			return uint16(x1(env).Elem().Uint())
		}
	case r.Uint32:
		fun = func(env *Env) uint32 {
			return uint32(x1(env).Elem().Uint())
		}
	case r.Uint64:
		fun = func(env *Env) uint64 {
			return x1(env).Elem().Uint()
		}
	case r.Uintptr:
		fun = func(env *Env) uintptr {
			return uintptr(x1(env).Elem().Uint())
		}
	case r.Float32:
		fun = func(env *Env) float32 {
			return float32(x1(env).Elem().Float())
		}
	case r.Float64:
		fun = func(env *Env) float64 {
			return x1(env).Elem().Float()
		}
	case r.Complex64:
		fun = func(env *Env) complex64 {
			return complex64(x1(env).Elem().Complex())
		}
	case r.Complex128:
		fun = func(env *Env) complex128 {
			return x1(env).Elem().Complex()
		}
	case r.String:
		fun = func(env *Env) string {
			return x1(env).Elem().String()
		}
	default:
		fun = func(env *Env) r.Value {
			return x1(env).Elem()
		}
	}
	return fun
}