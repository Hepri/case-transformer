# case-transformer: golang library for code case transformation [![Build Status](https://travis-ci.org/Hepri/case-transformer.png?branch=master)](https://travis-ci.org/Hepri/case-transformer)

StringUp: Golang Extension Library for Strings!

## Installation

Install:

    go get github.com/Hepri/case-transformer

## Usage

Possible transformations:
- camelCase
- PascalCase
- snake_case (underscore)
- kebab-case

```
// -> 'someString'
StringToCamelCase('SomeString')

// -> 'SomeString'
StringToPascalCase('someString')

// -> 'some_string'
StringToSnakeCase('someString')

// -> 'some_string'
StringToUnderscore('someString')

// -> 'some-string'
StringToKebabCase('someString')
```
