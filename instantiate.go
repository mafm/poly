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

package poly

import (
	"github.com/wdamron/poly/types"
)

func (ctx *commonContext) instantiate(level int, t types.Type) types.Type {
	if !t.IsGeneric() {
		return t
	}
	ctx.clearInstantiationLookup()
	ctx.clearRecursiveInstantiationLookup()
	return ctx.visitInstantiate(level, t)
}

func (ctx *commonContext) visitInstantiate(level int, t types.Type) types.Type {
	if !t.IsGeneric() {
		return t
	}

	switch t := t.(type) {
	case *types.Var:
		id := t.Id()
		if tv, ok := ctx.instLookup[id]; ok {
			return tv
		}

		tv := ctx.varTracker.New(level)
		if t.IsWeakVar() {
			tv.SetWeak()
		}

		constraints := t.Constraints()
		constraints2 := make([]types.InstanceConstraint, len(constraints))
		copy(constraints2, constraints)
		tv.SetConstraints(constraints2)
		ctx.instLookup[id] = tv
		return tv

	case *types.RecursiveLink:
		rec := t.Recursive
		if next := ctx.recInstLookup[t.Recursive.Id]; next != nil {
			return &types.RecursiveLink{Recursive: next, Index: t.Index}
		}
		next := &types.Recursive{
			Id:           rec.Id,
			Types:        make([]*types.App, len(rec.Types)),
			Names:        rec.Names,
			Indexes:      rec.Indexes,
			Instantiated: true,
			Flags:        rec.Flags,
		}
		copy(next.Types, rec.Types)
		ctx.recInstLookup[t.Recursive.Id] = next
		for i, ti := range next.Types {
			next.Types[i] = ctx.visitInstantiate(level, ti).(*types.App)
		}
		next.Flags &^= types.ContainsGenericVars
		return &types.RecursiveLink{Recursive: next, Index: t.Index}

	case *types.App:
		args := make([]types.Type, len(t.Args))
		for i, arg := range t.Args {
			args[i] = ctx.visitInstantiate(level, arg)
		}
		var underlying types.Type
		if t.Underlying != nil {
			underlying = ctx.visitInstantiate(level, t.Underlying)
		}
		return &types.App{Const: ctx.visitInstantiate(level, t.Const), Args: args, Underlying: underlying, Flags: t.Flags &^ types.ContainsGenericVars}

	case *types.Arrow:
		args := make([]types.Type, len(t.Args))
		for i, arg := range t.Args {
			args[i] = ctx.visitInstantiate(level, arg)
		}
		return &types.Arrow{Args: args, Return: ctx.visitInstantiate(level, t.Return), Method: t.Method}

	case *types.Method:
		arrow := ctx.visitInstantiate(level, t.TypeClass.Methods[t.Name]).(*types.Arrow)
		arrow.Method = t
		return arrow

	case *types.Record:
		return &types.Record{Row: ctx.visitInstantiate(level, t.Row)}

	case *types.Variant:
		return &types.Variant{Row: ctx.visitInstantiate(level, t.Row)}

	case *types.RowExtend:
		m := t.Labels
		mb := m.Builder()
		m.Range(func(label string, ts types.TypeList) bool {
			lb := ts.Builder()
			ts.Range(func(i int, t types.Type) bool {
				lb.Set(i, ctx.visitInstantiate(level, t))
				return true
			})
			mb.Set(label, lb.Build())
			return true
		})
		row := t.Row
		if row == nil {
			row = types.RowEmptyPointer
		} else if _, ok := row.(*types.RowEmpty); !ok {
			row = ctx.visitInstantiate(level, t.Row)
		}
		return &types.RowExtend{Row: row, Labels: mb.Build()}
	}
	panic("unexpected generic type " + t.TypeName())
}
