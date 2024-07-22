package gofig

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

/************************************************************
	+--------------------------------------------------------+
	|            What's the deal with Gofig?				 |
	+--------------------------------------------------------+



gofig is a config management library that is designed to have:
1. Self-Documenting capabilities
2. A very simple API
3. Validation functionality for config values
4. Almost immutable config values
5. Ability to create multiple configuration sections for different parts of your application if needed

1. Self-Documenting Capabilities
gofig allows you to define all possible config options in one function-call. In theory, you could have one file that defines all possible config options for your application. This is useful for documentation purposes, as you can see all possible config options in one place. There's no need to update a readme file everytime you add a new config option - just tell the reader to check your go file where you initialize the config options.const

2. A very simple API
gofig contains only 2 functions:
- Init
- Get
There's not much to learn about gofig.

3. Validation functionality for config values
gofig allows you to optionally add validation functions to your config options. This provides more self-documenting code which will tell the reader the required values, the possible values, and even conditonal validation. For example, is `ENV_VAR_A` is set to `TRUE`, then `ENV_VAR_1`, `ENV_VAR_2`, and `ENV_VAR_3` must be set to something.

4. Almost immutable config values
gofig allows you to only call `Init` once per Gofig instance. Once `Init` has been called, the values Gofig contains can't be changed.

5. Ability to create multiple configuration sections for different parts of your application if needed
gofig allows you to create multiple Gofig instances. This is useful if you have different parts of your application that require different config options, or you want to separate sections of your application
*/

const (
	TypeBool   GfType = 0
	TypeInt    GfType = 1
	TypeFloat  GfType = 2
	TypeString GfType = 3
	numTypes   GfType = 4
)

type NoDefault struct{}

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

type Gofig struct {
	sync.Mutex        // prevents someone from doing something like having two threads and calling Init at the same time on the same Gofig object
	initialized bool  // prevents gofig from being initialized more than once. This attempts to make the config values read-only.
	valsByType  []any // slice of slices corresponding to the different types the config options could be.
	/*
		valsByType has this structure:
		[
			[]bool,
			[]int,
			[]float64,
			[]string,
		]
	*/
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

var ErrAlreadyInitialized = errors.New("Gofig object already initialized")
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

func lookupFromEnv(initOpt InitOpt) (string, error) {
	valStr, exists := os.LookupEnv(initOpt.Name)
	if !exists && initOpt.Required {
		return "", ErrRequiredConfigNotSet(initOpt.Name)
	}
	return valStr, nil
}

/***********************
	+---------------+
	|   Public API  |
	+---------------+
***********************/

func New(initOpts ...InitOpt) (*Gofig, error) {
	gf := &Gofig{}
	err := gf.Init(initOpts...)
	if err != nil {
		return nil, err
	}
	return gf, nil
}

/*
Init initializes the Gofig object with the config options passed in.
If Gofig has already been initialized, Init will return an error.
*/
func (gf *Gofig) Init(initOpts ...InitOpt) error {
	var valsBool []bool
	var valsInt []int
	var valsFloat []float64
	var valsString []string

	gf.Lock() // don't allow multiple threads to alter the Gofig object at the same time
	defer gf.Unlock()

	if gf.initialized {
		return ErrAlreadyInitialized
	}

	if len(initOpts) == 0 {
		return ErrNoInputOpts
	}

	for _, initOpt := range initOpts {
		if initOpt.Required && initOpt.Default != nil {
			return ErrDefaultNotNilWhenRequired(initOpt)
		}
		if !initOpt.Required && initOpt.Default == nil {
			return ErrDefaultIsNilWhenNotRequired(initOpt)
		}

		initOpt.IdPtr.t = initOpt.Type

		switch initOpt.Type {
		case TypeBool:

			var val bool
			if !initOpt.Required {
				if _, ok := initOpt.Default.(bool); !ok {
					return ErrDefaultValueIsWrongTypeWhenNotRequired(initOpt)
				}
				val = initOpt.Default.(bool)
			}

			valStr, err := lookupFromEnv(initOpt)
			if err != nil {
				return err
			}

			valConv := strings.ToUpper(valStr) == "TRUE"

			if valConv != val {
				val = valConv
			}

			initOpt.IdPtr.valIdx = len(valsBool)
			valsBool = append(valsBool, val)

		case TypeInt:

			var val int
			if !initOpt.Required {
				if _, ok := initOpt.Default.(int); !ok {
					return ErrDefaultValueIsWrongTypeWhenNotRequired(initOpt)
				}
				val = initOpt.Default.(int)
			}

			valStr, err := lookupFromEnv(initOpt)
			if err != nil {
				return err
			}

			valConv, err := strconv.Atoi(valStr)
			if err != nil {
				return ErrWrongTypeSetInEnvironment(initOpt, valStr)
			}

			if valConv != val {
				val = valConv
			}

			initOpt.IdPtr.valIdx = len(valsInt)
			valsInt = append(valsInt, val)

		case TypeFloat:

			var val float64
			if !initOpt.Required {
				if _, ok := initOpt.Default.(float64); !ok {
					return ErrDefaultValueIsWrongTypeWhenNotRequired(initOpt)
				}
				val = initOpt.Default.(float64)
			}

			valStr, err := lookupFromEnv(initOpt)
			if err != nil {
				return err
			}

			valConv, err := strconv.ParseFloat(valStr, 64)
			if err != nil {
				return ErrWrongTypeSetInEnvironment(initOpt, valStr)
			}

			if valConv != val {
				val = valConv
			}

			initOpt.IdPtr.valIdx = len(valsFloat)
			valsFloat = append(valsFloat, val)

		case TypeString:

			var val string
			if !initOpt.Required {
				if _, ok := initOpt.Default.(string); !ok {
					return ErrDefaultValueIsWrongTypeWhenNotRequired(initOpt)
				}
				val = initOpt.Default.(string)
			}

			valStr, err := lookupFromEnv(initOpt)
			if err != nil {
				return err
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

	gf.initialized = true

	return nil
}

/*
Get returns the value of the config option corresponding to the Id passed in.
If the Id is invalid, Get will return an error.
If Gofig has not been initialized, Get will return an error.
*/
func (gf *Gofig) Get(id Id) (any, error) {

	if !gf.initialized {
		return 0, ErrNotInitialized
	}

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
