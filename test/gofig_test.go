package gofig

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/ippontech/gofig"
)

/***************
* +-------------------+
* | Test Errors       |
* +-------------------+
****************/

var ErrExpectedNoError = func(actual error) error {
	return fmt.Errorf("expected no error, got: `%v`", actual)
}
var ErrExpectedError = errors.New("expected an error")
var ErrErrorsDoNotMatch = func(expected, actual error) error {
	return fmt.Errorf("expected error: `%v`, got: `%v`", expected, actual)
}
var ErrNotImplemented = errors.New("not implemented")

/***************
* +-------------------+
* | Test Functions    |
* +-------------------+
****************/

/*
Test_Init_ is just a general example test to show usage of the library.
Does a few checks in one
*/
func Test_Init_General(t *testing.T) {
	t.Setenv("FOO", "true")
	t.Setenv("BAR", "10")
	t.Setenv("BAZ", "hello")

	var fooId gofig.Id
	var barId gofig.Id
	var bazId gofig.Id

	gf, err := gofig.Init([]gofig.InitOpt{
		{
			Name:        "FOO",
			Description: "This is a foo. It is used for blah blah blah",
			Type:        gofig.TypeBool,
			Required:    true,
			IdPtr:       &fooId,
		},
		{
			Name:        "BAR",
			Description: "This is a bar. It is used for blah blah blah",
			Type:        gofig.TypeInt,
			Required:    true,
			IdPtr:       &barId,
		},
		{
			Name:        "BAZ",
			Description: "This is a baz. It is used for blah blah blah",
			Type:        gofig.TypeString,
			Required:    true,
			IdPtr:       &bazId,
		},
	})

	if err != nil {
		t.Error(ErrExpectedNoError(err))
	}

	// Now you can access the values like this:
	foo, errActual := gf.Get(fooId)
	if errActual != nil {
		fmt.Println(errActual)
	}
	fmt.Println(foo)

	bar, errActual := gf.Get(barId)
	if errActual != nil {
		fmt.Println(errActual)
	}
	fmt.Println(bar)

	baz, errActual := gf.Get(bazId)
	if errActual != nil {
		fmt.Println(errActual)
	}
	fmt.Println(baz)

}

func Test_Init_ErrNil_WhenIntDefatultTypeCorrect(t *testing.T) {
	t.Setenv("FOO", "10")

	var fooId gofig.Id

	_, errActual := gofig.Init([]gofig.InitOpt{
		{
			Name:        "FOO",
			Description: "This is a foo. It is used for blah blah blah",
			Type:        gofig.TypeInt,
			Required:    false,
			Default:     10,
			IdPtr:       &fooId,
		},
	})

	if errActual != nil {
		t.Error(ErrExpectedNoError(errActual))
	}

}

func Test_Init_Err_When_IntDefaultTypeIncorrect(t *testing.T) {
	var fooId gofig.Id

	badInitOpt := gofig.InitOpt{
		Name:        "FOO",
		Description: "This is a foo. It is used for blah blah blah",
		Type:        gofig.TypeInt,
		Required:    false,
		Default:     "incorrect type",
		IdPtr:       &fooId,
	}

	_, errActual := gofig.Init([]gofig.InitOpt{
		badInitOpt,
	})

	errExpected := gofig.ErrDefaultValueIsWrongTypeWhenNotRequired(badInitOpt)
	if errActual.Error() != errExpected.Error() {
		t.Error(ErrErrorsDoNotMatch(errExpected, errActual))
	}
}

func Test_Init_ErrNil_When_BoolDefaultTypeCorrect(t *testing.T) {
	var fooId gofig.Id

	_, errActual := gofig.Init([]gofig.InitOpt{
		{
			Name:        "FOO",
			Description: "This is a foo. It is used for blah blah blah",
			Type:        gofig.TypeBool,
			Required:    false,
			Default:     true,
			IdPtr:       &fooId,
		},
	})

	if errActual != nil {
		t.Error(ErrExpectedNoError(errActual))
	}
}

func Test_Init_Err_When_BoolDefaultTypeIncorrect(t *testing.T) {
	var fooId gofig.Id

	badInitOpt := gofig.InitOpt{
		Name:        "FOO",
		Description: "This is a foo. It is used for blah blah blah",
		Type:        gofig.TypeBool,
		Required:    false,
		Default:     "incorrect type",
		IdPtr:       &fooId,
	}

	_, errActual := gofig.Init([]gofig.InitOpt{
		badInitOpt,
	})

	errExpected := gofig.ErrDefaultValueIsWrongTypeWhenNotRequired(badInitOpt)
	if errActual.Error() != errExpected.Error() {
		t.Error(ErrErrorsDoNotMatch(errExpected, errActual))
	}
}

func Test_Init_ErrNil_WhenFloatDefaultTypeCorrect(t *testing.T) {
	t.Setenv("FOO", "10.0")

	var fooId gofig.Id

	_, errActual := gofig.Init([]gofig.InitOpt{
		{
			Name:        "FOO",
			Description: "This is a foo. It is used for blah blah blah",
			Type:        gofig.TypeFloat,
			Required:    false,
			Default:     10.0,
			IdPtr:       &fooId,
		},
	})

	if errActual != nil {
		t.Error(ErrExpectedNoError(errActual))
	}

}

func Test_Init_Err_When_FloatDefaultTypeIncorrect(t *testing.T) {
	var fooId gofig.Id

	badInitOpt := gofig.InitOpt{
		Name:        "FOO",
		Description: "This is a foo. It is used for blah blah blah",
		Type:        gofig.TypeFloat,
		Required:    false,
		Default:     "incorrect type",
		IdPtr:       &fooId,
	}

	_, errActual := gofig.Init([]gofig.InitOpt{
		badInitOpt,
	})

	errExpected := gofig.ErrDefaultValueIsWrongTypeWhenNotRequired(badInitOpt)
	if errActual.Error() != errExpected.Error() {
		t.Error(ErrErrorsDoNotMatch(errExpected, errActual))
	}
}

func Test_Init_ErrNil_WhenStringDefaultTypeCorrect(t *testing.T) {
	var fooId gofig.Id

	_, errActual := gofig.Init([]gofig.InitOpt{
		{
			Name:        "FOO",
			Description: "This is a foo. It is used for blah blah blah",
			Type:        gofig.TypeString,
			Required:    false,
			Default:     "hello",
			IdPtr:       &fooId,
		},
	})

	if errActual != nil {
		t.Error(ErrExpectedNoError(errActual))
	}
}

func Test_Init_Err_When_StringDefaultTypeIncorrect(t *testing.T) {
	var fooId gofig.Id

	badInitOpt := gofig.InitOpt{
		Name:        "FOO",
		Description: "This is a foo. It is used for blah blah blah",
		Type:        gofig.TypeString,
		Required:    false,
		Default:     10,
		IdPtr:       &fooId,
	}

	_, errActual := gofig.Init([]gofig.InitOpt{
		badInitOpt,
	})

	errExpected := gofig.ErrDefaultValueIsWrongTypeWhenNotRequired(badInitOpt)
	if errActual.Error() != errExpected.Error() {
		t.Error(ErrErrorsDoNotMatch(errExpected, errActual))
	}
}

func Test_Init_Err_When_EnvVarCannotBeConvertedToInt(t *testing.T) {
	val := "not an int"
	t.Setenv("FOO", val)

	var fooId gofig.Id

	_, errActual := gofig.Init([]gofig.InitOpt{
		{
			Name:        "FOO",
			Description: "This is a foo. It is used for blah blah blah",
			Type:        gofig.TypeInt,
			Required:    true,
			IdPtr:       &fooId,
		},
	})

	if errActual == nil {
		t.Error(ErrExpectedError)
	}

	errExpected := gofig.ErrWrongTypeSetInEnvironment(gofig.InitOpt{
		Name:     "FOO",
		Type:     gofig.TypeInt,
		Required: true,
	}, val)

	if errActual.Error() != errExpected.Error() {
		t.Error(ErrErrorsDoNotMatch(errExpected, errActual))
	}
}

func Test_Init_Err_When_EnvVarCannotBeConvertedToFloat64(t *testing.T) {
	val := "not a float"
	t.Setenv("FOO", val)

	var fooId gofig.Id

	_, errActual := gofig.Init([]gofig.InitOpt{
		{
			Name:        "FOO",
			Description: "This is a foo. It is used for blah blah blah",
			Type:        gofig.TypeFloat,
			Required:    true,
			IdPtr:       &fooId,
		},
	})

	if errActual == nil {
		t.Error(ErrExpectedError)
	}

	errExpected := gofig.ErrWrongTypeSetInEnvironment(gofig.InitOpt{
		Name:     "FOO",
		Type:     gofig.TypeFloat,
		Required: true,
	}, val)

	if errActual.Error() != errExpected.Error() {
		t.Error(ErrErrorsDoNotMatch(errExpected, errActual))
	}
}

func Test_Init_Err_When_NoInitOptsPassed(t *testing.T) {
	_, errActual := gofig.Init([]gofig.InitOpt{})

	errExpected := gofig.ErrNoInputOpts
	if errActual.Error() != errExpected.Error() {
		t.Errorf("expected error: `%v`, got: `%v`", errExpected, errActual)
	}
}

/*
Create a bunch of configuration Ids by calling Init with a bunch of correctly defined InitOpts with mixed types. Then test to see if all the values equal expected results.
*/
func Test_Get_ABunchOfCallsSuccess(t *testing.T) {
	const n = 30
	var ids [n]gofig.Id
	var opts [n]gofig.InitOpt
	var vals [n]any

	// generate a list of random init optsand corresponding environ variable values that were set
	for i := 0; i < n; i++ {
		opt, val := generateRandomInitOptWithEnvSet(&ids[i], t)
		opts[i] = opt
		vals[i] = val
	}

	gf, err := gofig.Init([]gofig.InitOpt{
		opts[0],
		opts[1],
		opts[2],
		opts[3],
		opts[4],
		opts[5],
		opts[6],
		opts[7],
		opts[8],
		opts[9],
		opts[10],
		opts[11],
		opts[12],
		opts[13],
		opts[14],
		opts[15],
		opts[16],
		opts[17],
		opts[18],
		opts[19],
		opts[20],
		opts[21],
		opts[22],
		opts[23],
		opts[24],
		opts[25],
		opts[26],
		opts[27],
		opts[28],
		opts[29],
	})

	if err != nil {
		t.Error(ErrExpectedNoError(err))
	}

	for i := 0; i < n; i++ {
		val, err := gf.Get(ids[i])
		if err != nil {
			t.Error(err)
		}
		if val != vals[i] {
			t.Errorf("expected: `%v`, got: `%v`", vals[i], val)
		}
	}
}

func Test_Get_Err_When_GofigNotInitialized(t *testing.T) {
	var fooId gofig.Id

	gf := gofig.Gofig{}

	_, errActual := gf.Get(fooId)
	errExpected := gofig.ErrNotInitialized
	if errActual.Error() != errExpected.Error() {
		t.Error(ErrErrorsDoNotMatch(errExpected, errActual))
	}
}

func Test_GetBool_Err_When_GofigNotInitialized(t *testing.T) {
	var fooId gofig.Id

	gf := gofig.Gofig{}

	_, errActual := gf.GetBool(fooId)
	errExpected := gofig.ErrNotInitialized
	if errActual.Error() != errExpected.Error() {
		t.Error(ErrErrorsDoNotMatch(errExpected, errActual))
	}
}

func Test_GetInt_Err_When_GofigNotInitialized(t *testing.T) {
	var fooId gofig.Id

	gf := gofig.Gofig{}

	_, errActual := gf.GetInt(fooId)
	errExpected := gofig.ErrNotInitialized
	if errActual.Error() != errExpected.Error() {
		t.Error(ErrErrorsDoNotMatch(errExpected, errActual))
	}
}

func Test_GetFloat_Err_When_GofigNotInitialized(t *testing.T) {
	var fooId gofig.Id

	gf := gofig.Gofig{}

	_, errActual := gf.GetFloat(fooId)
	errExpected := gofig.ErrNotInitialized
	if errActual.Error() != errExpected.Error() {
		t.Error(ErrErrorsDoNotMatch(errExpected, errActual))
	}
}

func Test_GetString_Err_When_GofigNotInitialized(t *testing.T) {
	var fooId gofig.Id

	gf := gofig.Gofig{}

	_, errActual := gf.GetString(fooId)
	errExpected := gofig.ErrNotInitialized
	if errActual.Error() != errExpected.Error() {
		t.Error(ErrErrorsDoNotMatch(errExpected, errActual))
	}
}

func Test_Get_Err_When_InvalidId(t *testing.T) {
	t.Setenv("FOO", "true")

	var fooId gofig.Id
	var invalidId gofig.Id // we never pass this into a call to Init

	initOpt := goodBoolInitOpt
	initOpt.IdPtr = &fooId

	gf, err := gofig.Init([]gofig.InitOpt{initOpt})

	if err != nil {
		t.Error(ErrExpectedNoError(err))
	}

	_, errActual := gf.Get(invalidId)
	errExpected := gofig.ErrInvalidId
	if errActual.Error() != errExpected.Error() {
		t.Error(ErrErrorsDoNotMatch(errExpected, errActual))
	}
}

func Test_GetBool_Err_When_InvalidId(t *testing.T) {
	t.Setenv("FOO", "true")

	var fooId gofig.Id
	var invalidId gofig.Id // we never pass this into a call to Init

	initOpt := goodBoolInitOpt
	initOpt.IdPtr = &fooId

	gf, err := gofig.Init([]gofig.InitOpt{initOpt})

	if err != nil {
		t.Error(ErrExpectedNoError(err))
	}

	_, errActual := gf.GetBool(invalidId)
	errExpected := gofig.ErrInvalidId
	if errActual.Error() != errExpected.Error() {
		t.Error(ErrErrorsDoNotMatch(errExpected, errActual))
	}
}

func Test_GetInt_Err_When_InvalidId(t *testing.T) {
	t.Setenv("FOO", "10")

	var fooId gofig.Id

	var invalidId gofig.Id // we never pass this into a call to Init

	initOpt := goodIntInitOpt
	initOpt.IdPtr = &fooId

	gf, err := gofig.Init([]gofig.InitOpt{initOpt})

	if err != nil {
		t.Error(ErrExpectedNoError(err))
	}

	_, errActual := gf.GetInt(invalidId)
	errExpected := gofig.ErrInvalidId
	if errActual.Error() != errExpected.Error() {
		t.Error(ErrErrorsDoNotMatch(errExpected, errActual))
	}
}

func Test_GetFloat_Err_When_InvalidId(t *testing.T) {
	t.Setenv("FOO", "10.0")

	var fooId gofig.Id

	var invalidId gofig.Id // we never pass this into a call to Init

	initOpt := goodFloatInitOpt
	initOpt.IdPtr = &fooId

	gf, err := gofig.Init([]gofig.InitOpt{initOpt})

	if err != nil {
		t.Error(ErrExpectedNoError(err))
	}

	_, errActual := gf.GetFloat(invalidId)
	errExpected := gofig.ErrInvalidId
	if errActual.Error() != errExpected.Error() {
		t.Error(ErrErrorsDoNotMatch(errExpected, errActual))
	}
}

func Test_GetString_Err_When_InvalidId(t *testing.T) {
	t.Setenv("FOO", "hello")

	var fooId gofig.Id

	var invalidId gofig.Id // we never pass this into a call to Init

	initOpt := goodStringInitOpt
	initOpt.IdPtr = &fooId

	gf, err := gofig.Init([]gofig.InitOpt{initOpt})

	if err != nil {
		t.Error(ErrExpectedNoError(err))
	}

	_, errActual := gf.GetString(invalidId)
	errExpected := gofig.ErrInvalidId

	if errActual.Error() != errExpected.Error() {
		t.Error(ErrErrorsDoNotMatch(errExpected, errActual))
	}
}

func Test_DocString_Matches_Expected(t *testing.T) {
	expectedDocStr := "FOO\n\tDescription: This is a foo. It is used for blah blah blah\n\tType: bool\n\tRequired: true\nFOO\n\tDescription: This is a foo. It is used for blah blah blah\n\tType: int\n\tRequired: true\nFOO\n\tDescription: This is a foo. It is used for blah blah blah\n\tType: float\n\tRequired: true\nFOO\n\tDescription: This is a foo. It is used for blah blah blah\n\tType: string\n\tRequired: true\n"

	initOpts := []gofig.InitOpt{
		goodBoolInitOpt,
		goodIntInitOpt,
		goodFloatInitOpt,
		goodStringInitOpt,
	}
	actualDocStr, err := gofig.DocString(initOpts)
	if err != nil {
		t.Error(ErrExpectedNoError(err))
	}
	if actualDocStr != expectedDocStr {
		t.Errorf("expected: `%v`, got: `%v`", expectedDocStr, actualDocStr)
	}
}

/***************
* +-------------------+
* | helper vars       |
* +-------------------+
****************/

var goodBoolInitOpt = gofig.InitOpt{
	Name:        "FOO",
	Description: "This is a foo. It is used for blah blah blah",
	Type:        gofig.TypeBool,
	Required:    true,
}

var goodIntInitOpt = gofig.InitOpt{
	Name:        "FOO",
	Description: "This is a foo. It is used for blah blah blah",
	Type:        gofig.TypeInt,
	Required:    true,
}

var goodFloatInitOpt = gofig.InitOpt{
	Name:        "FOO",
	Description: "This is a foo. It is used for blah blah blah",
	Type:        gofig.TypeFloat,
	Required:    true,
}

var goodStringInitOpt = gofig.InitOpt{
	Name:        "FOO",
	Description: "This is a foo. It is used for blah blah blah",
	Type:        gofig.TypeString,
	Required:    true,
}

/**************
* +-------------------+
* | Helper Functions  |
* +-------------------+
**************/

func generateRandomStrOfLenN(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		randLen, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		b[i] = letters[int(randLen.Int64())]
	}
	return string(b)
}

// Intn is a shortcut for generating a random integer between 0 and
// max using crypto/rand.
func intn(max int64) int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		panic(err)
	}
	return nBig.Int64()
}

// returns random init opt as well as the value of the random environment variable that was set
func generateRandomInitOptWithEnvSet(id *gofig.Id, t *testing.T) (gofig.InitOpt, any) {
	randOptNameLen, _ := rand.Int(rand.Reader, big.NewInt(63))
	randOptDescLen, _ := rand.Int(rand.Reader, big.NewInt(63))
	randIsRequiredBigInt, _ := rand.Int(rand.Reader, big.NewInt(2))
	randTypeBigInt, _ := rand.Int(rand.Reader, big.NewInt(4))
	randomOptName := generateRandomStrOfLenN(int(randOptNameLen.Int64()) + 1) // random string between 1 and 64 characters
	randomOptDesc := generateRandomStrOfLenN(int(randOptDescLen.Int64()) + 1) // random string between 1 and 64 characters
	randomRequired := randIsRequiredBigInt.Int64() == 1
	randomType := gofig.GfType(int(randTypeBigInt.Int64()))

	var randomVal interface{}
	switch randomType {
	case gofig.TypeBool:
		rval, _ := rand.Int(rand.Reader, big.NewInt(2))
		randomVal = rval.Int64() == 1
	case gofig.TypeInt:
		rval, _ := rand.Int(rand.Reader, big.NewInt(100))
		randomVal = int(rval.Int64())
	case gofig.TypeFloat:
		randomVal = float64(intn(1<<53)) / (1 << 53)
	case gofig.TypeString:
		rStrLen, _ := rand.Int(rand.Reader, big.NewInt(63))
		randomVal = generateRandomStrOfLenN(int(rStrLen.Int64())) // random string between 1 and 64 characters
	}

	var randomDefault interface{}
	if !randomRequired {
		switch randomType {
		case gofig.TypeBool:
			rval, _ := rand.Int(rand.Reader, big.NewInt(2))
			randomDefault = rval.Int64() == 1
		case gofig.TypeInt:
			rval, _ := rand.Int(rand.Reader, big.NewInt(100))
			randomDefault = int(rval.Int64())
		case gofig.TypeFloat:
			randomDefault = float64(intn(1<<53)) / (1 << 53)
		case gofig.TypeString:
			rStrLen, _ := rand.Int(rand.Reader, big.NewInt(63))
			randomDefault = generateRandomStrOfLenN(int(rStrLen.Int64())) // random string between 1 and 64 characters
		}
	}
	opt := gofig.InitOpt{
		Name:        randomOptName,
		Description: randomOptDesc,
		Type:        randomType,
		Required:    randomRequired,
		Default:     randomDefault,
		IdPtr:       id,
	}

	// convert randomVal to string depending on randomType and set in environ
	envValStr := fmt.Sprintf("%v", randomVal)
	t.Setenv(randomOptName, envValStr)

	return opt, randomVal

}
