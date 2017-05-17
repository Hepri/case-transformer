# case-transformer: golang library for code case transformation [![Build Status](https://travis-ci.org/Hepri/case-transformer.png?branch=master)](https://travis-ci.org/Hepri/case-transformer)

## Installation

Install:

    go get github.com/Hepri/case-transformer

## Usage

Possible transformations:
- camelCase
- PascalCase
- snake_case (underscore_case)
- kebab-case

```
package main

import (
    "fmt"
    "github.com/Hepri/case-transformer"
)

func main() {
    // myTestString
    fmt.Println(case_transformer.StringToCamelCase("my test string"))
    // MyTestString
    fmt.Println(case_transformer.StringToPascalCase("my test string"))
    // my_test_string
    fmt.Println(case_transformer.StringToSnakeCase("my test string"))
    // my_test_string
    fmt.Println(case_transformer.StringToUnderscore("my test string"))
    // my-test-string
    fmt.Println(case_transformer.StringToKebabCase("my test string"))
}

```
