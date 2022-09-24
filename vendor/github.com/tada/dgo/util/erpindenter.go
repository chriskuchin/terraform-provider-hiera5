package util

import (
	"github.com/tada/dgo/dgo"
)

type erpIndenter struct {
	dgo.Indenter
	seen []dgo.Indentable
}

// ToStringERP will produce an unindented string from an Indentable using an indenter returned
// by NewERPIndenter
func ToStringERP(ia dgo.Indentable) string {
	i := &erpIndenter{Indenter: NewIndenter(``)}
	i.seen = append(i.seen, ia)
	ia.AppendTo(i)
	return i.String()
}

// ToIndentedStringERP will produce a string from an Indentable using an indenter returned
// by NewERPIndenter that has been initialized with a two space indentation.
func ToIndentedStringERP(ia dgo.Indentable) string {
	i := NewERPIndenter(`  `).(*erpIndenter)
	i.seen = append(i.seen, ia)
	ia.AppendTo(i)
	return i.String()
}

// NewERPIndenter creates an endless recursion protected indenter capable of indenting self referencing
// values. When an endless recursion is encountered, the string <recursive self reference> is emitted
// rather than the value itself.
func NewERPIndenter(indent string) dgo.Indenter {
	return &erpIndenter{Indenter: NewIndenter(indent)}
}

func (i *erpIndenter) Indent() dgo.Indenter {
	return &erpIndenter{Indenter: i.Indenter.Indent(), seen: i.seen}
}

func (i *erpIndenter) AppendValue(v interface{}) {
	if vi, ok := v.(dgo.Indentable); ok {
		s := i.seen
		for n := range s {
			if s[n] == v {
				i.Append(`<recursive self reference`)
				var dv dgo.Value
				if dv, ok = v.(dgo.Value); ok {
					i.Append(` to `)
					i.Append(dv.Type().TypeIdentifier().String())
				}
				i.AppendRune('>')
				return
			}
		}
		i.seen = append(i.seen, vi)
		vi.AppendTo(i)
		i.seen = s
		return
	}
	i.Indenter.AppendValue(v)
}
