package dgo

type (
	// TypeOp describes the logical operation of an unary or ternary type
	TypeOp int

	// UnaryType represents a type that applies an operation on another type
	UnaryType interface {
		Type

		// Operator returns the unary operator OpNot
		Operator() TypeOp

		// Operand returns the operand type
		Operand() Type
	}

	// TernaryType represents a collection of types using a logical operator
	TernaryType interface {
		Type

		// Operator returns the ternary operator OpAnd or OpOr
		Operator() TypeOp

		// Operands returns the types that this ternary type operates on
		Operands() Array
	}
)

const (
	// OpNot is the unary negation operator
	OpNot = TypeOp(iota)

	// OpAnd is the ternary operator returned by the AllOf type
	OpAnd

	// OpOr is the ternary operator returned by the AnyOf type
	OpOr

	// OpOne is the ternary operator returned by the OneOf type
	OpOne

	// OpMeta is the unary operator returned by the Meta type
	OpMeta

	// OpSensitive is the unary operator returned by the Sensitive type
	OpSensitive
)
