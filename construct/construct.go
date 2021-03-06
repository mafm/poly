// The MIT License (MIT)
//
// Copyright (c) 2019 West Damron
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package construct

import (
	"github.com/wdamron/poly/ast"
	"github.com/wdamron/poly/types"
)

// Types

// Unit type: `()`
func TUnit() *types.Unit { return types.NewUnit() }

// Create a new type-variable with the given id and binding-level.
func TVar(id, level uint) *types.Var {
	return types.NewVar(id, level)
}

// Type constant: `int`, `bool`, etc
func TConst(name string) *types.Const {
	return &types.Const{Name: name}
}

// Size constant: `array[int, 8]`
func TSize(size int) types.Size {
	return types.Size(size)
}

// Recursive link to a type.
func TRecursiveLink(rec *types.Recursive, name string) *types.RecursiveLink {
	return &types.RecursiveLink{Recursive: rec, Index: rec.Indexes[name]}
}

// Type application: `list[int]`
func TApp(constructor types.Type, params ...types.Type) *types.App {
	return &types.App{Const: constructor, Params: params}
}

func TAlias(app *types.App, underlying types.Type) *types.App {
	return &types.App{Const: app.Const, Params: app.Params, Underlying: underlying}
}

func TRef(deref types.Type) *types.App {
	return types.NewRef(deref)
}

// Function type: `(int, int) -> int`
func TArrow(args []types.Type, ret types.Type) *types.Arrow {
	return &types.Arrow{Args: args, Return: ret}
}

// Function type: `int -> int`
func TArrow1(arg types.Type, ret types.Type) *types.Arrow {
	return &types.Arrow{Args: []types.Type{arg}, Return: ret}
}

// Function type: `(int, int) -> int`
func TArrow2(arg1, arg2 types.Type, ret types.Type) *types.Arrow {
	return &types.Arrow{Args: []types.Type{arg1, arg2}, Return: ret}
}

// Function type: `(int, int, int) -> int`
func TArrow3(arg1, arg2, arg3 types.Type, ret types.Type) *types.Arrow {
	return &types.Arrow{Args: []types.Type{arg1, arg2, arg3}, Return: ret}
}

// Type-class method type: `('a, int) -> 'a`
func TMethod(typeClass *types.TypeClass, name string) *types.Method {
	return &types.Method{TypeClass: typeClass, Name: name}
}

// Record type: `{...}`
func TRecord(row types.Type) *types.Record {
	return &types.Record{Row: row}
}

// Record type with fixed labels: `{...}`
func TRecordFlat(labels map[string]types.Type) *types.Record {
	return TRecord(TRowExtend(nil, TypeMap(labels)))
}

// Tagged (ad-hoc) variant-type: `[...]`
func TVariant(row types.Type) *types.Variant {
	return &types.Variant{Row: row}
}

// Row extension: `<a : _ , b : _ | ...>`
func TRowExtend(row types.Type, labels types.TypeMap) *types.RowExtend {
	if row == nil {
		row = types.RowEmptyPointer
	}
	return &types.RowExtend{Row: row, Labels: labels}
}

// Create a TypeMap with unscoped labels.
func TypeMap(m map[string]types.Type) types.TypeMap {
	return types.NewFlatTypeMap(m)
}

// Empty row: `<>`
func TRowEmpty() *types.RowEmpty {
	return types.RowEmptyPointer
}

// Expressions:

func Literal(syntax string, usingVars []string, constructType func(env types.TypeEnv, level uint, using []types.Type) (types.Type, error)) *ast.Literal {
	return &ast.Literal{Syntax: syntax, Construct: constructType}
}

// Variable
func Var(name string) *ast.Var {
	return &ast.Var{Name: name}
}

// Dereference: `*x`
func Deref(ref ast.Expr) *ast.Deref {
	return &ast.Deref{Ref: ref}
}

// Dereference and assign: `*x = y`
func DerefAssign(ref ast.Expr, value ast.Expr) *ast.DerefAssign {
	return &ast.DerefAssign{Ref: ref, Value: value}
}

// Application: `f(x)`
func Call(f ast.Expr, args ...ast.Expr) *ast.Call {
	return &ast.Call{Func: f, Args: args}
}

// Abstraction: `fn (x, y) -> x`
func Func(args []string, body ast.Expr) *ast.Func {
	return &ast.Func{ArgNames: args, Body: body}
}

// Abstraction: `fn (x) -> x`
func Func1(arg string, body ast.Expr) *ast.Func {
	return &ast.Func{ArgNames: []string{arg}, Body: body}
}

// Abstraction: `fn (x, y) -> x`
func Func2(arg1, arg2 string, body ast.Expr) *ast.Func {
	return &ast.Func{ArgNames: []string{arg1, arg2}, Body: body}
}

// Abstraction: `fn (x, y, z) -> x`
func Func3(arg1, arg2, arg3 string, body ast.Expr) *ast.Func {
	return &ast.Func{ArgNames: []string{arg1, arg2, arg3}, Body: body}
}

// Control flow graph
func ControlFlow(name string, locals ...string) *ast.ControlFlow {
	return ast.NewControlFlow(name, locals...)
}

// Pipeline: `pipe $ = xs |> fmap($, fn (x) -> to_y(x)) |> fmap($, fn (y) -> to_z(y))`
func Pipe(as string, sequence ...ast.Expr) *ast.Pipe {
	return &ast.Pipe{Source: sequence[0], As: as, Sequence: sequence[1:]}
}

// Let-binding: `let a = 1 in e`
func Let(varName string, value ast.Expr, body ast.Expr) *ast.Let {
	return &ast.Let{Var: varName, Value: value, Body: body}
}

// Grouped let-bindings: `let a = 1 and b = 2 in e`
func LetGroup(vars []ast.LetBinding, body ast.Expr) *ast.LetGroup {
	return &ast.LetGroup{Vars: vars, Body: body}
}

// Paired identifier and value
func LetBinding(varName string, value ast.Expr) ast.LetBinding {
	return ast.LetBinding{Var: varName, Value: value}
}

// Selecting value of label: `r.a`
func RecordSelect(record ast.Expr, label string) *ast.RecordSelect {
	return &ast.RecordSelect{Record: record, Label: label}
}

// Deleting label: `{r - a}`
func RecordRestrict(record ast.Expr, label string) *ast.RecordRestrict {
	return &ast.RecordRestrict{Record: record, Label: label}
}

// Extending record: `{a = 1, b = 2 | r}`
func RecordExtend(record ast.Expr, labels ...ast.LabelValue) *ast.RecordExtend {
	if record == nil {
		record = RecordEmpty()
	}
	return &ast.RecordExtend{Record: record, Labels: labels}
}

// Paired label and value
func LabelValue(label string, value ast.Expr) ast.LabelValue {
	return ast.LabelValue{Label: label, Value: value}
}

// Empty record: `{}`
func RecordEmpty() *ast.RecordEmpty {
	return &ast.RecordEmpty{}
}

// Tagged (ad-hoc) variant: `:X a`
func Variant(label string, value ast.Expr) *ast.Variant {
	return &ast.Variant{Label: label, Value: value}
}

// Pattern-matching case expression over tagged (ad-hoc) variant-types:
//
//  match e {
//      :X a -> expr1
//    | :Y b -> expr2
//    |  ...
//    | z -> default_expr (optional)
//  }
func Match(value ast.Expr, cases []ast.MatchCase, defaultCase *ast.MatchCase) *ast.Match {
	return &ast.Match{Value: value, Cases: cases, Default: defaultCase}
}

// Case expression within Match: `:X a -> expr1`
func MatchCase(label string, varName string, value ast.Expr) ast.MatchCase {
	return ast.MatchCase{Label: label, Var: varName, Value: value}
}
