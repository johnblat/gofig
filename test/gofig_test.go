package gofig

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/ippontech/gofig"
)

var ErrExpectedNoError = errors.New("expected no error")
var ErrExpectedError = errors.New("expected an error")
var ErrErrorsDoNotMatch = func(expected, actual error) error {
	return fmt.Errorf("expected error: `%v`, got: `%v`", expected, actual)
}

/*
TestInit is just a general example test to show usage of the library.
Does a few checks in one
*/
func TestInit(t *testing.T) {
	os.Setenv("FOO", "true")
	os.Setenv("BAR", "10")
	os.Setenv("BAZ", "hello")

	gf := gofig.Gofig{}

	var fooId gofig.Id
	var barId gofig.Id
	var bazId gofig.Id

	err := gf.Init(
		gofig.InitOpt{
			Name:        "FOO",
			Description: "This is a foo. It is used for blah blah blah",
			Type:        gofig.TypeBool,
			Required:    true,
			IdPtr:       &fooId,
		},
		gofig.InitOpt{
			Name:        "BAR",
			Description: "This is a bar. It is used for blah blah blah",
			Type:        gofig.TypeInt,
			Required:    true,
			IdPtr:       &barId,
		},
		gofig.InitOpt{
			Name:        "BAZ",
			Description: "This is a baz. It is used for blah blah blah",
			Type:        gofig.TypeString,
			Required:    true,
			IdPtr:       &bazId,
		},
	)

	if err != nil {
		t.Error(ErrExpectedNoError)
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

func TestGetFailsWhenNotInitialized(t *testing.T) {
	gf := gofig.Gofig{}

	var fooId gofig.Id

	_, errActual := gf.Get(fooId)
	if errActual == nil {
		t.Error(ErrExpectedError)
	}
	if errActual != gofig.ErrNotInitialized {
		t.Error(ErrErrorsDoNotMatch(gofig.ErrNotInitialized, errActual))
	}
}

func TestInitFailsWhenCalledMoreThanOnce(t *testing.T) {
	os.Setenv("FOO", "true")
	os.Setenv("BAR", "false")

	gf := gofig.Gofig{}

	var fooId gofig.Id

	_ = gf.Init( // not testing the error here
		gofig.InitOpt{
			Name:        "FOO",
			Description: "This is a foo. It is used for blah blah blah",
			Type:        gofig.TypeBool,
			Required:    true,
			IdPtr:       &fooId,
		},
	)

	errActual := gf.Init(
		gofig.InitOpt{
			Name:        "FOO",
			Description: "This is a foo. It is used for blah blah blah",
			Type:        gofig.TypeBool,
			Required:    true,
			IdPtr:       &fooId,
		},
	)

	if errActual == nil {
		t.Error(ErrExpectedError)
	}

	errExpected := gofig.ErrAlreadyInitialized
	if errActual != errExpected {
		t.Error(ErrErrorsDoNotMatch(gofig.ErrAlreadyInitialized, errActual))
	}
}

func TestInitErrIsNilWhenIntDefatultTypeCorrect(t *testing.T) {
	os.Setenv("FOO", "10")

	gf := gofig.Gofig{}

	var fooId gofig.Id

	errActual := gf.Init(
		gofig.InitOpt{
			Name:        "FOO",
			Description: "This is a foo. It is used for blah blah blah",
			Type:        gofig.TypeInt,
			Required:    false,
			Default:     10,
			IdPtr:       &fooId,
		},
	)

	if errActual != nil {
		t.Error(ErrExpectedNoError)
	}

}

func TestInitFailsWhenIntDefaultTypeIncorrect(t *testing.T) {
	gf := gofig.Gofig{}

	var fooId gofig.Id

	badInitOpt := gofig.InitOpt{
		Name:        "FOO",
		Description: "This is a foo. It is used for blah blah blah",
		Type:        gofig.TypeInt,
		Required:    false,
		Default:     "incorrect type",
		IdPtr:       &fooId,
	}

	errActual := gf.Init(
		badInitOpt,
	)

	errExpected := gofig.ErrDefaultValueIsWrongTypeWhenNotRequired(badInitOpt)
	if errActual.Error() != errExpected.Error() {
		t.Error(ErrErrorsDoNotMatch(errExpected, errActual))
	}
}

func TestInitFailsWhenBoolDefaultTypeIncorrect(t *testing.T) {
	gf := gofig.Gofig{}

	var fooId gofig.Id

	badInitOpt := gofig.InitOpt{
		Name:        "FOO",
		Description: "This is a foo. It is used for blah blah blah",
		Type:        gofig.TypeBool,
		Required:    false,
		Default:     "incorrect type",
		IdPtr:       &fooId,
	}

	errActual := gf.Init(
		badInitOpt,
	)

	errExpected := gofig.ErrDefaultValueIsWrongTypeWhenNotRequired(badInitOpt)
	if errActual.Error() != errExpected.Error() {
		t.Error(ErrErrorsDoNotMatch(errExpected, errActual))
	}
}

func TestInitFailsWhenFloatDefaultTypeIncorrect(t *testing.T) {
	gf := gofig.Gofig{}

	var fooId gofig.Id

	badInitOpt := gofig.InitOpt{
		Name:        "FOO",
		Description: "This is a foo. It is used for blah blah blah",
		Type:        gofig.TypeFloat,
		Required:    false,
		Default:     "incorrect type",
		IdPtr:       &fooId,
	}

	errActual := gf.Init(
		badInitOpt,
	)

	errExpected := gofig.ErrDefaultValueIsWrongTypeWhenNotRequired(badInitOpt)
	if errActual.Error() != errExpected.Error() {
		t.Error(ErrErrorsDoNotMatch(errExpected, errActual))
	}
}

func TestInitFailsWhenStringDefaultTypeIncorrect(t *testing.T) {
	gf := gofig.Gofig{}

	var fooId gofig.Id

	badInitOpt := gofig.InitOpt{
		Name:        "FOO",
		Description: "This is a foo. It is used for blah blah blah",
		Type:        gofig.TypeString,
		Required:    false,
		Default:     10,
		IdPtr:       &fooId,
	}

	errActual := gf.Init(
		badInitOpt,
	)

	errExpected := gofig.ErrDefaultValueIsWrongTypeWhenNotRequired(badInitOpt)
	if errActual.Error() != errExpected.Error() {
		t.Error(ErrErrorsDoNotMatch(errExpected, errActual))
	}
}

func TestInitFailsWhenEnvVarCannotBeConvertedToInt(t *testing.T) {
	val := "not an int"
	os.Setenv("FOO", val)

	gf := gofig.Gofig{}

	var fooId gofig.Id

	errActual := gf.Init(
		gofig.InitOpt{
			Name:        "FOO",
			Description: "This is a foo. It is used for blah blah blah",
			Type:        gofig.TypeInt,
			Required:    true,
			IdPtr:       &fooId,
		},
	)

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

func TestInitFailsWhenEnvVarCannotBeConvertedToFloat64(t *testing.T) {
	val := "not a float"
	os.Setenv("FOO", val)

	gf := gofig.Gofig{}

	var fooId gofig.Id

	errActual := gf.Init(
		gofig.InitOpt{
			Name:        "FOO",
			Description: "This is a foo. It is used for blah blah blah",
			Type:        gofig.TypeFloat,
			Required:    true,
			IdPtr:       &fooId,
		},
	)

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

func TestInitFailsWhenNoInitOptsPassed(t *testing.T) {
	gf := gofig.Gofig{}

	errActual := gf.Init()

	errExpected := gofig.ErrNoInputOpts
	if errActual.Error() != errExpected.Error() {
		t.Errorf("expected error: `%v`, got: `%v`", errExpected, errActual)
	}
}

/*
Create a bunch of configuration Ids by calling Init with a bunch of correctly defined InitOpts with mixed types. Then test to see if all the values equal expected results.
*/
func TestGetABunchOfCallsSuccess(t *testing.T) {
	gf := gofig.Gofig{}

	const n = 30
	var ids [n]gofig.Id
	var opts [n]gofig.InitOpt
	var vals [n]any

	// generate a list of random init optsand corresponding environ variable values that were set
	for i := 0; i < n; i++ {
		opt, val := generateRandomInitOptWithEnvSet(&ids[i])
		opts[i] = opt
		vals[i] = val
	}

	err := gf.Init(
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
	)

	if err != nil {
		t.Error(ErrExpectedNoError)
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
func generateRandomInitOptWithEnvSet(id *gofig.Id) (gofig.InitOpt, any) {
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
	os.Setenv(randomOptName, envValStr)

	return opt, randomVal

}
