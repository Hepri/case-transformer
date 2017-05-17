package case_transformer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAlphaUpper(t *testing.T) {
	for i := 'A'; i < 'Z'; i++ {
		assert.True(t, IsAlphaUpper(i))
	}
	for i := 'a'; i < 'z'; i++ {
		assert.False(t, IsAlphaUpper(i))
	}
	// some special
	assert.False(t, IsAlphaUpper('0'))
	assert.False(t, IsAlphaUpper('?'))
	assert.False(t, IsAlphaUpper('!'))
	assert.False(t, IsAlphaUpper('#'))
}

func TestIsAlphaLower(t *testing.T) {
	for i := 'A'; i < 'Z'; i++ {
		assert.False(t, IsAlphaLower(i))
	}
	for i := 'a'; i < 'z'; i++ {
		assert.True(t, IsAlphaLower(i))
	}
	// some special
	assert.False(t, IsAlphaLower('0'))
	assert.False(t, IsAlphaLower('?'))
	assert.False(t, IsAlphaLower('!'))
	assert.False(t, IsAlphaLower('#'))
}

func TestIsDigit(t *testing.T) {
	assert.True(t, IsDigit('0'))
	assert.True(t, IsDigit('1'))
	assert.True(t, IsDigit('2'))
	assert.True(t, IsDigit('3'))
	assert.True(t, IsDigit('4'))
	assert.True(t, IsDigit('5'))
	assert.True(t, IsDigit('6'))
	assert.True(t, IsDigit('7'))
	assert.True(t, IsDigit('8'))
	assert.True(t, IsDigit('9'))

	assert.False(t, IsDigit('z'))
	assert.False(t, IsDigit('k'))
	assert.False(t, IsDigit('l'))
	assert.False(t, IsDigit('A'))
	assert.False(t, IsDigit('B'))
}

func TestIsDelimiter(t *testing.T) {
	assert.True(t, IsDelimiter('-'))
	assert.True(t, IsDelimiter('_'))
	assert.True(t, IsDelimiter(' '))
	assert.True(t, IsDelimiter('?'))
	assert.True(t, IsDelimiter('!'))
	assert.True(t, IsDelimiter('*'))
	assert.True(t, IsDelimiter('('))
	assert.True(t, IsDelimiter(')'))

	assert.False(t, IsDelimiter('A'))
	assert.False(t, IsDelimiter('a'))
	assert.False(t, IsDelimiter('1'))
}

func TestIsAbbreviation(t *testing.T) {
	assert.True(t, IsAbbreviation("ABC"))
	assert.True(t, IsAbbreviation("JSON"))

	assert.False(t, IsAbbreviation("json"))
	assert.False(t, IsAbbreviation("String"))
	assert.False(t, IsAbbreviation("someString"))
	assert.False(t, IsAbbreviation("ABC-ABC"))
	assert.False(t, IsAbbreviation("ABC_abc"))
}

func TestInitCap(t *testing.T) {
	assert.Equal(t, "Abc", InitCap("abc"))
	assert.Equal(t, "Abc", InitCap("Abc"))
	assert.Equal(t, "Abc", InitCap("aBC"))
	// abbreviation
	assert.Equal(t, "ABC", InitCap("ABC"))
}

func testStringParts(t *testing.T, str string, parts []string) {
	s := SplitToParts(str)
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
