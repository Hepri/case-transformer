package case_transformer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAlphaUpper(t *testing.T) {
	for i := 'A'; i < 'Z'; i++ {
		assert.True(t, isAlphaUpper(i))
	}
	for i := 'a'; i < 'z'; i++ {
		assert.False(t, isAlphaUpper(i))
	}
	// some special
	assert.False(t, isAlphaUpper('0'))
	assert.False(t, isAlphaUpper('?'))
	assert.False(t, isAlphaUpper('!'))
	assert.False(t, isAlphaUpper('#'))
}

func TestIsAlphaLower(t *testing.T) {
	for i := 'A'; i < 'Z'; i++ {
		assert.False(t, isAlphaLower(i))
	}
	for i := 'a'; i < 'z'; i++ {
		assert.True(t, isAlphaLower(i))
	}
	// some special
	assert.False(t, isAlphaLower('0'))
	assert.False(t, isAlphaLower('?'))
	assert.False(t, isAlphaLower('!'))
	assert.False(t, isAlphaLower('#'))
}

func TestIsDigit(t *testing.T) {
	assert.True(t, isDigit('0'))
	assert.True(t, isDigit('1'))
	assert.True(t, isDigit('2'))
	assert.True(t, isDigit('3'))
	assert.True(t, isDigit('4'))
	assert.True(t, isDigit('5'))
	assert.True(t, isDigit('6'))
	assert.True(t, isDigit('7'))
	assert.True(t, isDigit('8'))
	assert.True(t, isDigit('9'))

	assert.False(t, isDigit('z'))
	assert.False(t, isDigit('k'))
	assert.False(t, isDigit('l'))
	assert.False(t, isDigit('A'))
	assert.False(t, isDigit('B'))
}

func TestIsDelimiter(t *testing.T) {
	assert.True(t, isDelimiter('-'))
	assert.True(t, isDelimiter('_'))
	assert.True(t, isDelimiter(' '))
	assert.True(t, isDelimiter('?'))
	assert.True(t, isDelimiter('!'))
	assert.True(t, isDelimiter('*'))
	assert.True(t, isDelimiter('('))
	assert.True(t, isDelimiter(')'))

	assert.False(t, isDelimiter('A'))
	assert.False(t, isDelimiter('a'))
	assert.False(t, isDelimiter('1'))
}

func TestIsAbbreviation(t *testing.T) {
	assert.True(t, isAbbreviation("ABC"))
	assert.True(t, isAbbreviation("JSON"))

	assert.False(t, isAbbreviation("json"))
	assert.False(t, isAbbreviation("String"))
	assert.False(t, isAbbreviation("someString"))
	assert.False(t, isAbbreviation("ABC-ABC"))
	assert.False(t, isAbbreviation("ABC_abc"))
}

func TestInitCap(t *testing.T) {
	assert.Equal(t, "Abc", initCap("abc"))
	assert.Equal(t, "Abc", initCap("Abc"))
	assert.Equal(t, "Abc", initCap("aBC"))
	// abbreviation
	assert.Equal(t, "ABC", initCap("ABC"))
}

func testStringParts(t *testing.T, str string, parts []string) {
	s := splitToParts(str)
	assert.EqualValues(t, parts, s)
}

func TestSplitToParts_SnakeCase(t *testing.T) {
	testStringParts(t, "some_string", []string{"some", "string"})
}

func TestSplitToParts_CamelCase(t *testing.T) {
	testStringParts(t, "someString", []string{"some", "String"})
}

func TestSplitToParts_PascalCase(t *testing.T) {
	testStringParts(t, "SomeString", []string{"Some", "String"})
}

func TestSplitToParts_KebabCase(t *testing.T) {
	testStringParts(t, "some-string", []string{"some", "string"})
	testStringParts(t, "Some-String", []string{"Some", "String"})
}

func TestSplitToParts_Delimiters(t *testing.T) {
	testStringParts(t, "a-b", []string{"a", "b"})
	testStringParts(t, "A-B", []string{"A", "B"})
	testStringParts(t, "A-b", []string{"A", "b"})
	testStringParts(t, "a_b", []string{"a", "b"})
	testStringParts(t, "A_B", []string{"A", "B"})
	testStringParts(t, "A_b", []string{"A", "b"})
}

func TestSplitToParts_Digits(t *testing.T) {
	testStringParts(t, "some_string12", []string{"some", "string12"})
	testStringParts(t, "12some22_string12", []string{"12some22", "string12"})

	testStringParts(t, "SomeString12", []string{"Some", "String12"})
	testStringParts(t, "someString12", []string{"some", "String12"})
}

func TestSplitToParts_Abbreviation(t *testing.T) {
	testStringParts(t, "JSONString", []string{"JSON", "String"})
	testStringParts(t, "JSON_String", []string{"JSON", "String"})
	testStringParts(t, "JSON-String", []string{"JSON", "String"})
	testStringParts(t, "StringABCString", []string{"String", "ABC", "String"})
	testStringParts(t, "ABCStringABCString", []string{"ABC", "String", "ABC", "String"})
	testStringParts(t, "ABC_String-ABCString", []string{"ABC", "String", "ABC", "String"})
}

var transformArr = [][]string{
	// [0] = initial
	// [1] = camel
	// [2] = pascal
	// [3] = snake (underscore)
	// [4] = kebab
	[]string{"some_string", "someString", "SomeString", "some_string", "some-string"},
	[]string{"some-string", "someString", "SomeString", "some_string", "some-string"},
	[]string{"Some_String", "someString", "SomeString", "some_string", "some-string"},
	[]string{"Some_string", "someString", "SomeString", "some_string", "some-string"},
	[]string{"someString", "someString", "SomeString", "some_string", "some-string"},
	[]string{"SomeString", "someString", "SomeString", "some_string", "some-string"},
	[]string{"JSONString", "jsonString", "JSONString", "json_string", "json-string"},
	[]string{"StringJSON", "stringJSON", "StringJSON", "string_json", "string-json"},
}

func testTransform(t *testing.T, transformIdx int, fn func(str string) string) {
	for _, p := range transformArr {
		assert.Equal(t, p[transformIdx], fn(p[0]))
	}
}

func TestStringToCamelCase(t *testing.T) {
	testTransform(t, 1, StringToCamelCase)
}

func TestStringToPascalCase(t *testing.T) {
	testTransform(t, 2, StringToPascalCase)
}

func TestStringToSnakeCase(t *testing.T) {
	testTransform(t, 3, StringToSnakeCase)
}

func TestStringToUnderscore(t *testing.T) {
	testTransform(t, 3, StringToUnderscore)
}

func TestStringToKebabCase(t *testing.T) {
	testTransform(t, 4, StringToKebabCase)
}
