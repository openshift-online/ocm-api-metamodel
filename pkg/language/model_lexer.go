// Code generated from ModelLexer.g4 by ANTLR 4.9.3. DO NOT EDIT.

package language

import (
	"fmt"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = unicode.IsLetter

var serializedLexerAtn = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 2, 31, 261,
	8, 1, 4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7,
	9, 7, 4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12,
	4, 13, 9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4,
	18, 9, 18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 4, 22, 9, 22, 4, 23,
	9, 23, 4, 24, 9, 24, 4, 25, 9, 25, 4, 26, 9, 26, 4, 27, 9, 27, 4, 28, 9,
	28, 4, 29, 9, 29, 4, 30, 9, 30, 4, 31, 9, 31, 4, 32, 9, 32, 3, 2, 3, 2,
	3, 2, 3, 2, 3, 2, 3, 2, 3, 2, 3, 2, 3, 2, 3, 2, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5,
	3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7,
	3, 8, 3, 8, 3, 8, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 3, 10, 3, 10, 3, 10, 3,
	10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11,
	3, 11, 3, 12, 3, 12, 3, 12, 3, 12, 3, 13, 3, 13, 3, 13, 3, 13, 3, 13, 3,
	13, 3, 13, 3, 13, 3, 13, 3, 13, 3, 14, 3, 14, 3, 14, 3, 14, 3, 14, 3, 14,
	3, 14, 3, 14, 3, 14, 3, 15, 3, 15, 3, 15, 3, 15, 3, 15, 3, 15, 3, 15, 3,
	16, 3, 16, 3, 16, 3, 16, 3, 16, 3, 16, 3, 16, 3, 17, 3, 17, 3, 17, 3, 17,
	3, 17, 3, 18, 3, 18, 3, 18, 3, 18, 3, 18, 3, 18, 3, 19, 3, 19, 3, 19, 3,
	19, 3, 19, 3, 19, 3, 19, 3, 19, 3, 19, 3, 20, 3, 20, 3, 21, 3, 21, 3, 22,
	3, 22, 3, 23, 3, 23, 3, 24, 3, 24, 3, 25, 6, 25, 195, 10, 25, 13, 25, 14,
	25, 196, 3, 26, 3, 26, 5, 26, 201, 10, 26, 3, 27, 3, 27, 7, 27, 205, 10,
	27, 12, 27, 14, 27, 208, 11, 27, 3, 27, 3, 27, 3, 28, 3, 28, 7, 28, 214,
	10, 28, 12, 28, 14, 28, 217, 11, 28, 3, 28, 3, 28, 3, 29, 3, 29, 7, 29,
	223, 10, 29, 12, 29, 14, 29, 226, 11, 29, 3, 30, 3, 30, 3, 30, 3, 30, 7,
	30, 232, 10, 30, 12, 30, 14, 30, 235, 11, 30, 3, 30, 3, 30, 3, 30, 3, 30,
	3, 31, 3, 31, 3, 31, 3, 31, 7, 31, 245, 10, 31, 12, 31, 14, 31, 248, 11,
	31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 32, 6, 32, 256, 10, 32, 13, 32,
	14, 32, 257, 3, 32, 3, 32, 5, 206, 215, 246, 2, 33, 3, 3, 5, 4, 7, 5, 9,
	6, 11, 7, 13, 8, 15, 9, 17, 10, 19, 11, 21, 12, 23, 13, 25, 14, 27, 15,
	29, 16, 31, 17, 33, 18, 35, 19, 37, 20, 39, 21, 41, 22, 43, 23, 45, 24,
	47, 25, 49, 26, 51, 27, 53, 2, 55, 2, 57, 28, 59, 29, 61, 30, 63, 31, 3,
	2, 7, 3, 2, 50, 59, 5, 2, 67, 92, 97, 97, 99, 124, 6, 2, 50, 59, 67, 92,
	97, 97, 99, 124, 4, 2, 12, 12, 15, 15, 5, 2, 11, 12, 15, 15, 34, 34, 2,
	266, 2, 3, 3, 2, 2, 2, 2, 5, 3, 2, 2, 2, 2, 7, 3, 2, 2, 2, 2, 9, 3, 2,
	2, 2, 2, 11, 3, 2, 2, 2, 2, 13, 3, 2, 2, 2, 2, 15, 3, 2, 2, 2, 2, 17, 3,
	2, 2, 2, 2, 19, 3, 2, 2, 2, 2, 21, 3, 2, 2, 2, 2, 23, 3, 2, 2, 2, 2, 25,
	3, 2, 2, 2, 2, 27, 3, 2, 2, 2, 2, 29, 3, 2, 2, 2, 2, 31, 3, 2, 2, 2, 2,
	33, 3, 2, 2, 2, 2, 35, 3, 2, 2, 2, 2, 37, 3, 2, 2, 2, 2, 39, 3, 2, 2, 2,
	2, 41, 3, 2, 2, 2, 2, 43, 3, 2, 2, 2, 2, 45, 3, 2, 2, 2, 2, 47, 3, 2, 2,
	2, 2, 49, 3, 2, 2, 2, 2, 51, 3, 2, 2, 2, 2, 57, 3, 2, 2, 2, 2, 59, 3, 2,
	2, 2, 2, 61, 3, 2, 2, 2, 2, 63, 3, 2, 2, 2, 3, 65, 3, 2, 2, 2, 5, 75, 3,
	2, 2, 2, 7, 81, 3, 2, 2, 2, 9, 86, 3, 2, 2, 2, 11, 91, 3, 2, 2, 2, 13,
	97, 3, 2, 2, 2, 15, 103, 3, 2, 2, 2, 17, 106, 3, 2, 2, 2, 19, 111, 3, 2,
	2, 2, 21, 119, 3, 2, 2, 2, 23, 126, 3, 2, 2, 2, 25, 130, 3, 2, 2, 2, 27,
	140, 3, 2, 2, 2, 29, 149, 3, 2, 2, 2, 31, 156, 3, 2, 2, 2, 33, 163, 3,
	2, 2, 2, 35, 168, 3, 2, 2, 2, 37, 174, 3, 2, 2, 2, 39, 183, 3, 2, 2, 2,
	41, 185, 3, 2, 2, 2, 43, 187, 3, 2, 2, 2, 45, 189, 3, 2, 2, 2, 47, 191,
	3, 2, 2, 2, 49, 194, 3, 2, 2, 2, 51, 200, 3, 2, 2, 2, 53, 202, 3, 2, 2,
	2, 55, 211, 3, 2, 2, 2, 57, 220, 3, 2, 2, 2, 59, 227, 3, 2, 2, 2, 61, 240,
	3, 2, 2, 2, 63, 255, 3, 2, 2, 2, 65, 66, 7, 99, 2, 2, 66, 67, 7, 118, 2,
	2, 67, 68, 7, 118, 2, 2, 68, 69, 7, 116, 2, 2, 69, 70, 7, 107, 2, 2, 70,
	71, 7, 100, 2, 2, 71, 72, 7, 119, 2, 2, 72, 73, 7, 118, 2, 2, 73, 74, 7,
	103, 2, 2, 74, 4, 3, 2, 2, 2, 75, 76, 7, 101, 2, 2, 76, 77, 7, 110, 2,
	2, 77, 78, 7, 99, 2, 2, 78, 79, 7, 117, 2, 2, 79, 80, 7, 117, 2, 2, 80,
	6, 3, 2, 2, 2, 81, 82, 7, 101, 2, 2, 82, 83, 7, 113, 2, 2, 83, 84, 7, 102,
	2, 2, 84, 85, 7, 103, 2, 2, 85, 8, 3, 2, 2, 2, 86, 87, 7, 103, 2, 2, 87,
	88, 7, 112, 2, 2, 88, 89, 7, 119, 2, 2, 89, 90, 7, 111, 2, 2, 90, 10, 3,
	2, 2, 2, 91, 92, 7, 103, 2, 2, 92, 93, 7, 116, 2, 2, 93, 94, 7, 116, 2,
	2, 94, 95, 7, 113, 2, 2, 95, 96, 7, 116, 2, 2, 96, 12, 3, 2, 2, 2, 97,
	98, 7, 104, 2, 2, 98, 99, 7, 99, 2, 2, 99, 100, 7, 110, 2, 2, 100, 101,
	7, 117, 2, 2, 101, 102, 7, 103, 2, 2, 102, 14, 3, 2, 2, 2, 103, 104, 7,
	107, 2, 2, 104, 105, 7, 112, 2, 2, 105, 16, 3, 2, 2, 2, 106, 107, 7, 110,
	2, 2, 107, 108, 7, 107, 2, 2, 108, 109, 7, 112, 2, 2, 109, 110, 7, 109,
	2, 2, 110, 18, 3, 2, 2, 2, 111, 112, 7, 110, 2, 2, 112, 113, 7, 113, 2,
	2, 113, 114, 7, 101, 2, 2, 114, 115, 7, 99, 2, 2, 115, 116, 7, 118, 2,
	2, 116, 117, 7, 113, 2, 2, 117, 118, 7, 116, 2, 2, 118, 20, 3, 2, 2, 2,
	119, 120, 7, 111, 2, 2, 120, 121, 7, 103, 2, 2, 121, 122, 7, 118, 2, 2,
	122, 123, 7, 106, 2, 2, 123, 124, 7, 113, 2, 2, 124, 125, 7, 102, 2, 2,
	125, 22, 3, 2, 2, 2, 126, 127, 7, 113, 2, 2, 127, 128, 7, 119, 2, 2, 128,
	129, 7, 118, 2, 2, 129, 24, 3, 2, 2, 2, 130, 131, 7, 114, 2, 2, 131, 132,
	7, 99, 2, 2, 132, 133, 7, 116, 2, 2, 133, 134, 7, 99, 2, 2, 134, 135, 7,
	111, 2, 2, 135, 136, 7, 103, 2, 2, 136, 137, 7, 118, 2, 2, 137, 138, 7,
	103, 2, 2, 138, 139, 7, 116, 2, 2, 139, 26, 3, 2, 2, 2, 140, 141, 7, 116,
	2, 2, 141, 142, 7, 103, 2, 2, 142, 143, 7, 117, 2, 2, 143, 144, 7, 113,
	2, 2, 144, 145, 7, 119, 2, 2, 145, 146, 7, 116, 2, 2, 146, 147, 7, 101,
	2, 2, 147, 148, 7, 103, 2, 2, 148, 28, 3, 2, 2, 2, 149, 150, 7, 117, 2,
	2, 150, 151, 7, 118, 2, 2, 151, 152, 7, 116, 2, 2, 152, 153, 7, 119, 2,
	2, 153, 154, 7, 101, 2, 2, 154, 155, 7, 118, 2, 2, 155, 30, 3, 2, 2, 2,
	156, 157, 7, 118, 2, 2, 157, 158, 7, 99, 2, 2, 158, 159, 7, 116, 2, 2,
	159, 160, 7, 105, 2, 2, 160, 161, 7, 103, 2, 2, 161, 162, 7, 118, 2, 2,
	162, 32, 3, 2, 2, 2, 163, 164, 7, 118, 2, 2, 164, 165, 7, 116, 2, 2, 165,
	166, 7, 119, 2, 2, 166, 167, 7, 103, 2, 2, 167, 34, 3, 2, 2, 2, 168, 169,
	7, 120, 2, 2, 169, 170, 7, 99, 2, 2, 170, 171, 7, 110, 2, 2, 171, 172,
	7, 119, 2, 2, 172, 173, 7, 103, 2, 2, 173, 36, 3, 2, 2, 2, 174, 175, 7,
	120, 2, 2, 175, 176, 7, 99, 2, 2, 176, 177, 7, 116, 2, 2, 177, 178, 7,
	107, 2, 2, 178, 179, 7, 99, 2, 2, 179, 180, 7, 100, 2, 2, 180, 181, 7,
	110, 2, 2, 181, 182, 7, 103, 2, 2, 182, 38, 3, 2, 2, 2, 183, 184, 7, 125,
	2, 2, 184, 40, 3, 2, 2, 2, 185, 186, 7, 127, 2, 2, 186, 42, 3, 2, 2, 2,
	187, 188, 7, 93, 2, 2, 188, 44, 3, 2, 2, 2, 189, 190, 7, 95, 2, 2, 190,
	46, 3, 2, 2, 2, 191, 192, 7, 63, 2, 2, 192, 48, 3, 2, 2, 2, 193, 195, 9,
	2, 2, 2, 194, 193, 3, 2, 2, 2, 195, 196, 3, 2, 2, 2, 196, 194, 3, 2, 2,
	2, 196, 197, 3, 2, 2, 2, 197, 50, 3, 2, 2, 2, 198, 201, 5, 53, 27, 2, 199,
	201, 5, 55, 28, 2, 200, 198, 3, 2, 2, 2, 200, 199, 3, 2, 2, 2, 201, 52,
	3, 2, 2, 2, 202, 206, 7, 36, 2, 2, 203, 205, 11, 2, 2, 2, 204, 203, 3,
	2, 2, 2, 205, 208, 3, 2, 2, 2, 206, 207, 3, 2, 2, 2, 206, 204, 3, 2, 2,
	2, 207, 209, 3, 2, 2, 2, 208, 206, 3, 2, 2, 2, 209, 210, 7, 36, 2, 2, 210,
	54, 3, 2, 2, 2, 211, 215, 7, 98, 2, 2, 212, 214, 11, 2, 2, 2, 213, 212,
	3, 2, 2, 2, 214, 217, 3, 2, 2, 2, 215, 216, 3, 2, 2, 2, 215, 213, 3, 2,
	2, 2, 216, 218, 3, 2, 2, 2, 217, 215, 3, 2, 2, 2, 218, 219, 7, 98, 2, 2,
	219, 56, 3, 2, 2, 2, 220, 224, 9, 3, 2, 2, 221, 223, 9, 4, 2, 2, 222, 221,
	3, 2, 2, 2, 223, 226, 3, 2, 2, 2, 224, 222, 3, 2, 2, 2, 224, 225, 3, 2,
	2, 2, 225, 58, 3, 2, 2, 2, 226, 224, 3, 2, 2, 2, 227, 228, 7, 49, 2, 2,
	228, 229, 7, 49, 2, 2, 229, 233, 3, 2, 2, 2, 230, 232, 10, 5, 2, 2, 231,
	230, 3, 2, 2, 2, 232, 235, 3, 2, 2, 2, 233, 231, 3, 2, 2, 2, 233, 234,
	3, 2, 2, 2, 234, 236, 3, 2, 2, 2, 235, 233, 3, 2, 2, 2, 236, 237, 8, 30,
	2, 2, 237, 238, 3, 2, 2, 2, 238, 239, 8, 30, 3, 2, 239, 60, 3, 2, 2, 2,
	240, 241, 7, 49, 2, 2, 241, 242, 7, 44, 2, 2, 242, 246, 3, 2, 2, 2, 243,
	245, 11, 2, 2, 2, 244, 243, 3, 2, 2, 2, 245, 248, 3, 2, 2, 2, 246, 247,
	3, 2, 2, 2, 246, 244, 3, 2, 2, 2, 247, 249, 3, 2, 2, 2, 248, 246, 3, 2,
	2, 2, 249, 250, 7, 44, 2, 2, 250, 251, 7, 49, 2, 2, 251, 252, 3, 2, 2,
	2, 252, 253, 8, 31, 3, 2, 253, 62, 3, 2, 2, 2, 254, 256, 9, 6, 2, 2, 255,
	254, 3, 2, 2, 2, 256, 257, 3, 2, 2, 2, 257, 255, 3, 2, 2, 2, 257, 258,
	3, 2, 2, 2, 258, 259, 3, 2, 2, 2, 259, 260, 8, 32, 3, 2, 260, 64, 3, 2,
	2, 2, 11, 2, 196, 200, 206, 215, 224, 233, 246, 257, 4, 3, 30, 2, 8, 2,
	2,
}

var lexerChannelNames = []string{
	"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
}

var lexerModeNames = []string{
	"DEFAULT_MODE",
}

var lexerLiteralNames = []string{
	"", "'attribute'", "'class'", "'code'", "'enum'", "'error'", "'false'",
	"'in'", "'link'", "'locator'", "'method'", "'out'", "'parameter'", "'resource'",
	"'struct'", "'target'", "'true'", "'value'", "'variable'", "'{'", "'}'",
	"'['", "']'", "'='",
}

var lexerSymbolicNames = []string{
	"", "ATTRIBUTE", "CLASS", "CODE", "ENUM", "ERROR", "FALSE", "IN", "LINK",
	"LOCATOR", "METHOD", "OUT", "PARAMETER", "RESOURCE", "STRUCT", "TARGET",
	"TRUE", "VALUE", "VARIABLE", "LEFT_CURLY_BRACKET", "RIGHT_CURLY_BRACKET",
	"LEFT_SQUARE_BRACKET", "RIGHT_SQUARE_BRACKET", "EQUALS_SIGN", "INTEGER_LITERAL",
	"STRING_LITERAL", "IDENTIFIER", "LINE_COMMENT", "BLOCK_COMMENT", "WS",
}

var lexerRuleNames = []string{
	"ATTRIBUTE", "CLASS", "CODE", "ENUM", "ERROR", "FALSE", "IN", "LINK", "LOCATOR",
	"METHOD", "OUT", "PARAMETER", "RESOURCE", "STRUCT", "TARGET", "TRUE", "VALUE",
	"VARIABLE", "LEFT_CURLY_BRACKET", "RIGHT_CURLY_BRACKET", "LEFT_SQUARE_BRACKET",
	"RIGHT_SQUARE_BRACKET", "EQUALS_SIGN", "INTEGER_LITERAL", "STRING_LITERAL",
	"SHORT_STRING", "LONG_STRING", "IDENTIFIER", "LINE_COMMENT", "BLOCK_COMMENT",
	"WS",
}

type ModelLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

// NewModelLexer produces a new lexer instance for the optional input antlr.CharStream.
//
// The *ModelLexer instance produced may be reused by calling the SetInputStream method.
// The initial lexer configuration is expensive to construct, and the object is not thread-safe;
// however, if used within a Golang sync.Pool, the construction cost amortizes well and the
// objects can be used in a thread-safe manner.
func NewModelLexer(input antlr.CharStream) *ModelLexer {
	l := new(ModelLexer)
	lexerDeserializer := antlr.NewATNDeserializer(nil)
	lexerAtn := lexerDeserializer.DeserializeFromUInt16(serializedLexerAtn)
	lexerDecisionToDFA := make([]*antlr.DFA, len(lexerAtn.DecisionToState))
	for index, ds := range lexerAtn.DecisionToState {
		lexerDecisionToDFA[index] = antlr.NewDFA(ds, index)
	}
	l.BaseLexer = antlr.NewBaseLexer(input)
	l.Interpreter = antlr.NewLexerATNSimulator(l, lexerAtn, lexerDecisionToDFA, antlr.NewPredictionContextCache())

	l.channelNames = lexerChannelNames
	l.modeNames = lexerModeNames
	l.RuleNames = lexerRuleNames
	l.LiteralNames = lexerLiteralNames
	l.SymbolicNames = lexerSymbolicNames
	l.GrammarFileName = "ModelLexer.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// ModelLexer tokens.
const (
	ModelLexerATTRIBUTE            = 1
	ModelLexerCLASS                = 2
	ModelLexerCODE                 = 3
	ModelLexerENUM                 = 4
	ModelLexerERROR                = 5
	ModelLexerFALSE                = 6
	ModelLexerIN                   = 7
	ModelLexerLINK                 = 8
	ModelLexerLOCATOR              = 9
	ModelLexerMETHOD               = 10
	ModelLexerOUT                  = 11
	ModelLexerPARAMETER            = 12
	ModelLexerRESOURCE             = 13
	ModelLexerSTRUCT               = 14
	ModelLexerTARGET               = 15
	ModelLexerTRUE                 = 16
	ModelLexerVALUE                = 17
	ModelLexerVARIABLE             = 18
	ModelLexerLEFT_CURLY_BRACKET   = 19
	ModelLexerRIGHT_CURLY_BRACKET  = 20
	ModelLexerLEFT_SQUARE_BRACKET  = 21
	ModelLexerRIGHT_SQUARE_BRACKET = 22
	ModelLexerEQUALS_SIGN          = 23
	ModelLexerINTEGER_LITERAL      = 24
	ModelLexerSTRING_LITERAL       = 25
	ModelLexerIDENTIFIER           = 26
	ModelLexerLINE_COMMENT         = 27
	ModelLexerBLOCK_COMMENT        = 28
	ModelLexerWS                   = 29
)

func (l *ModelLexer) comment() {
	comments[l.GetLine()] = l.GetText()
}

func (l *ModelLexer) Action(localctx antlr.RuleContext, ruleIndex, actionIndex int) {
	switch ruleIndex {
	case 28:
		l.LINE_COMMENT_Action(localctx, actionIndex)

	default:
		panic("No registered action for: " + fmt.Sprint(ruleIndex))
	}
}

func (l *ModelLexer) LINE_COMMENT_Action(localctx antlr.RuleContext, actionIndex int) {
	this := l
	_ = this

	switch actionIndex {
	case 0:
		l.comment()

	default:
		panic("No registered action for: " + fmt.Sprint(actionIndex))
	}
}
