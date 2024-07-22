package gofig

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	TypeBool   GfType = 0
	TypeInt    GfType = 1
	TypeFloat  GfType = 2
	TypeString GfType = 3
	numTypes   GfType = 4
)

var typeNames = []string{
	"bool",
	"int",
	"float",
	"string",
}

type GfType int

type Id struct {
	t      GfType // the GfType num will correspond to the index in the Gofig.valsByType slice
	valIdx int    // the index of the value in the slice of the corresponding GfType
}

/*
valsByType has this structure:
[

	[]bool,
	[]int,
	[]float64,
	[]string,

]
*/
type Gofig struct {
	valsByType []any // slice of slices corresponding to the different types the config options could be.
}

type InitOpt struct {
	Name        string // The name of the config option (e.g. "ENV_VAR_A")
	Description string // A description of the config option
	Type        GfType // The type of the config option (e.g. TypeBool, TypeInt, TypeFloat, TypeString)
	Required    bool   // Whether the config option is required
	Default     any    // The default value of the config option. Doesn't do anything if the config option is required.
	IdPtr       *Id    // Pointer to the Id of the config option. This is where you store the Id. The Id value be set after the call to Init.
}

/*
**********************
	+-----------------+
	|Error Definitions|
	+-----------------+
**********************
*/

var ErrNoInputOpts = errors.New("no initOpts provided. must provice initOpts to initialize a Gofig object")
var ErrInvalidId = errors.New("invalid id")
var ErrNotInitialized = errors.New("Gofig not initialized. Call Init() first")
var ErrDefaultValueIsWrongTypeWhenNotRequired = func(initOpt InitOpt) error {
	return fmt.Errorf(
		"config: `%v`. type: `%v`. default value of `%v` is not of type `%v`",
		initOpt.Name,
		typeNames[initOpt.Type],
		initOpt.Default,
		typeNames[initOpt.Type],
	)
}
var ErrRequiredConfigNotSet = func(name string) error {
	return fmt.Errorf("required config option %s not set", name)
}
var ErrDefaultNotNilWhenRequired = func(initOpt InitOpt) error {
	return fmt.Errorf("config: `%v`. required: true. default value: `%v`. default value must be nil when config is required", initOpt.Name, initOpt.Default)
}
var ErrDefaultIsNilWhenNotRequired = func(initOpt InitOpt) error {
	return fmt.Errorf("config: `%v`. required: false. default value: `nil`. default value must not be nil when config is not required", initOpt.Name)
}
var ErrWrongTypeSetInEnvironment = func(initOpt InitOpt, valFromEnviron string) error {
	return fmt.Errorf("config `%s` of type `%s` was not set as `%s` in environment. environment value: `%s`", initOpt.Name, typeNames[initOpt.Type], typeNames[initOpt.Type], valFromEnviron)
}

/**********************
    +-----------------+
    |Private functions|
	+-----------------+
***********************/

func isDefaultTypeCorrect(initOpt InitOpt) bool {
	var GfTypeMapReflectionKind = map[GfType]reflect.Kind{
		TypeBool:   reflect.Bool,
		TypeInt:    reflect.Int,
		TypeFloat:  reflect.Float64,
		TypeString: reflect.String,
	}

	if !initOpt.Required {
		defaultValKind := reflect.TypeOf(initOpt.Default).Kind()
		if defaultValKind != GfTypeMapReflectionKind[initOpt.Type] {
			return false
		}
	}
	return true
}

/***********************
	+---------------+
	|   Public API  |
	+---------------+
***********************/

/*
DocString returns a string that contains the documentation for the config options passed in.
*/
func DocString(initOpts []InitOpt) (string, error) {
	if len(initOpts) == 0 {
		return "", ErrNoInputOpts
	}

	var docs string

	for _, initOpt := range initOpts {
		docs += fmt.Sprintf(
			"%s\n\tDescription: %s\n\tType: %s\n\tRequired: %v\n",
			initOpt.Name,
			initOpt.Description,
			typeNames[initOpt.Type],
			initOpt.Required,
		)

		if !initOpt.Required {
			docs += fmt.Sprintf("\tDefault: %v\n", initOpt.Default)
		}
	}

	return docs, nil

}

/*
Init initializes the Gofig object with the config options passed in.
If Gofig has already been initialized, Init will return an error.
*/
func Init(initOpts []InitOpt) (Gofig, error) {
	gf := Gofig{}

	var valsBool []bool
	var valsInt []int
	var valsFloat []float64
	var valsString []string

	if len(initOpts) == 0 {
		return gf, ErrNoInputOpts
	}

	for _, initOpt := range initOpts {
		if initOpt.Required && initOpt.Default != nil {
			return gf, ErrDefaultNotNilWhenRequired(initOpt)
		}
		if !initOpt.Required && initOpt.Default == nil {
			return gf, ErrDefaultIsNilWhenNotRequired(initOpt)
		}

		initOpt.IdPtr.t = initOpt.Type

		switch initOpt.Type {
		case TypeBool:

			if ok := isDefaultTypeCorrect(initOpt); !ok {
				return gf, ErrDefaultValueIsWrongTypeWhenNotRequired(initOpt)
			}

			var val bool
			if !initOpt.Required {
				val = initOpt.Default.(bool)
			}

			valStr, exists := os.LookupEnv(initOpt.Name)
			if !exists && initOpt.Required {
				return gf, ErrRequiredConfigNotSet(initOpt.Name)
			}

			valConv := strings.ToUpper(valStr) == "TRUE"

			if valConv != val {
				val = valConv
			}

			initOpt.IdPtr.valIdx = len(valsBool)
			valsBool = append(valsBool, val)

		case TypeInt:

			if ok := isDefaultTypeCorrect(initOpt); !ok {
				return gf, ErrDefaultValueIsWrongTypeWhenNotRequired(initOpt)
			}

			var val int
			if !initOpt.Required {
				val = initOpt.Default.(int)
			}

			valStr, exists := os.LookupEnv(initOpt.Name)
			if !exists && initOpt.Required {
				return gf, ErrRequiredConfigNotSet(initOpt.Name)
			}

			valConv, err := strconv.Atoi(valStr)
			if err != nil {
				return gf, ErrWrongTypeSetInEnvironment(initOpt, valStr)
			}

			if valConv != val {
				val = valConv
			}

			initOpt.IdPtr.valIdx = len(valsInt)
			valsInt = append(valsInt, val)

		case TypeFloat:

			if ok := isDefaultTypeCorrect(initOpt); !ok {
				return gf, ErrDefaultValueIsWrongTypeWhenNotRequired(initOpt)
			}

			var val float64
			if !initOpt.Required {
				val = initOpt.Default.(float64)
			}

			valStr, exists := os.LookupEnv(initOpt.Name)
			if !exists && initOpt.Required {
				return gf, ErrRequiredConfigNotSet(initOpt.Name)
			}

			valConv, err := strconv.ParseFloat(valStr, 64)
			if err != nil {
				return gf, ErrWrongTypeSetInEnvironment(initOpt, valStr)
			}

			if valConv != val {
				val = valConv
			}

			initOpt.IdPtr.valIdx = len(valsFloat)
			valsFloat = append(valsFloat, val)

		case TypeString:

			if ok := isDefaultTypeCorrect(initOpt); !ok {
				return gf, ErrDefaultValueIsWrongTypeWhenNotRequired(initOpt)
			}

			var val string
			if !initOpt.Required {
				val = initOpt.Default.(string)
			}

			valStr, exists := os.LookupEnv(initOpt.Name)
			if !exists && initOpt.Required {
				return gf, ErrRequiredConfigNotSet(initOpt.Name)
			}

			if valStr != val {
				val = valStr
			}

			initOpt.IdPtr.valIdx = len(valsString)
			valsString = append(valsString, val)
		}
	}

	gf.valsByType = []any{
		valsBool,
		valsInt,
		valsFloat,
		valsString,
	}

	return gf, nil
}

/*
Get returns the value of the config option corresponding to the Id passed in.
If the Id is invalid, Get will return an error.
If Gofig has not been initialized, Get will return an error.
*/
func (gf *Gofig) Get(id Id) (any, error) {
	if id.t < 0 || id.t >= numTypes {
		return nil, ErrInvalidId
	}

	if id.valIdx < 0 {
		return nil, ErrInvalidId
	}

	switch id.t {
	case TypeBool:
		if id.valIdx >= len(gf.valsByType[id.t].([]bool)) {
			return nil, ErrInvalidId
		}
		return gf.valsByType[id.t].([]bool)[id.valIdx], nil

	case TypeInt:
		if id.valIdx >= len(gf.valsByType[id.t].([]int)) {
			return nil, ErrInvalidId
		}
		return gf.valsByType[id.t].([]int)[id.valIdx], nil

	case TypeFloat:
		if id.valIdx >= len(gf.valsByType[id.t].([]float64)) {
			return nil, ErrInvalidId
		}
		return gf.valsByType[id.t].([]float64)[id.valIdx], nil

	case TypeString:
		if id.valIdx >= len(gf.valsByType[id.t].([]string)) {
			return nil, ErrInvalidId
		}
		return gf.valsByType[id.t].([]string)[id.valIdx], nil

	}

	// if somehow we get here, just return not valid id
	return nil, ErrInvalidId
}

// More Get-family functions for bool, int, float64, and string

/*
GetBool
*/
func (gf *Gofig) GetBool(id Id) (bool, error) {
	if id.t != TypeBool {
		return false, ErrInvalidId
	}
	if id.valIdx < 0 {
		return false, ErrInvalidId
	}
	if id.valIdx >= len(gf.valsByType[id.t].([]bool)) {
		return false, ErrInvalidId
	}
	return gf.valsByType[id.t].([]bool)[id.valIdx], nil
}

func (gf *Gofig) GetInt(id Id) (int, error) {
	if id.t != TypeInt {
		return 0, ErrInvalidId
	}
	if id.valIdx < 0 {
		return 0, ErrInvalidId
	}
	if id.valIdx >= len(gf.valsByType[id.t].([]int)) {
		return 0, ErrInvalidId
	}
	return gf.valsByType[id.t].([]int)[id.valIdx], nil
}

func (gf *Gofig) GetFloat(id Id) (float64, error) {
	if id.t != TypeFloat {
		return 0.0, ErrInvalidId
	}
	if id.valIdx < 0 {
		return 0.0, ErrInvalidId
	}
	if id.valIdx >= len(gf.valsByType[id.t].([]float64)) {
		return 0.0, ErrInvalidId
	}
	return gf.valsByType[id.t].([]float64)[id.valIdx], nil
}

func (gf *Gofig) GetString(id Id) (string, error) {
	if id.t != TypeString {
		return "", ErrInvalidId
	}
	if id.valIdx < 0 {
		return "", ErrInvalidId
	}
	if id.valIdx >= len(gf.valsByType[id.t].([]string)) {
		return "", ErrInvalidId
	}
	return gf.valsByType[id.t].([]string)[id.valIdx], nil
}
