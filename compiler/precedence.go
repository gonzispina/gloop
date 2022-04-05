package compiler

type Precedence int

const (
	precedenceNone       = iota
	precedenceAssigment  // =
	precedenceOr         // or
	precedenceAnd        // and
	precedenceEquality   // =
	precedenceComparison // < > <= >=
	precedenceTerm       // +
	precedenceFactor     // *
	precedenceUnary      // not
	precedenceCall       // ()
	precedencePrimary
)

type parseFunc func() error

type parseRule struct {
	prefix     parseFunc
	infix      parseFunc
	precedence Precedence
}

func getRule(c *Compiler, tt tokenType) parseRule {
	var rules = map[tokenType]parseRule{
		Plus:         {nil, c.binary, precedenceTerm},
		Star:         {nil, c.binary, precedenceFactor},
		Equal:        {nil, c.binary, precedenceEquality},
		Lesser:       {nil, c.binary, precedenceComparison},
		LesserEqual:  {nil, c.binary, precedenceComparison},
		Greater:      {nil, c.binary, precedenceComparison},
		GreaterEqual: {nil, c.binary, precedenceComparison},
		Not:          {c.unary, nil, precedenceNone},

		LeftArrow:          {nil, nil, precedenceNone},
		LeftParen:          {c.grouping, nil, precedenceNone},
		RightParen:         {nil, nil, precedenceNone},
		LeftSquareBracket:  {nil, nil, precedenceNone},
		RightSquareBracket: {nil, nil, precedenceNone},
		Comma:              {nil, nil, precedenceNone},

		// FALSE         : {literal,  nil,   precedenceNone},
		// TRUE          : {literal,  nil,   precedenceNone},
		DefineProcedure: {nil, nil, precedenceNone},
		If:              {nil, nil, precedenceNone},
		Then:            {nil, nil, precedenceNone},
		Else:            {nil, nil, precedenceNone},
		// And:           {nil, c.and, precedenceAnd},
		// Or:            {nil, c.or, precedenceOr},
		Loop:          {nil, nil, precedenceNone},
		AbortLoop:     {nil, nil, precedenceNone},
		Times:         {nil, nil, precedenceNone},
		EndProcedure:  {nil, nil, precedenceNone},
		QuitProcedure: {nil, nil, precedenceNone},

		Output: {nil, nil, precedenceNone},

		// Identifier: {c.variable, nil, precedenceNone},
		// Constant:   {c.constant, nil, precedenceNone},

		Eof: {nil, nil, precedenceNone},
	}

	return rules[tt]
}
