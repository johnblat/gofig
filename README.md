# gofig

## What's the Deal With gofig?

**gofig** is a self-documenting **configuration library** for **config-heavy apps**. It was created because:
- As projects grow, configuration options grow and change. **gofig** makes it really hard to **not** document configuration options accurately and comprehensively. 
- **gofig** tries to prevent configuration related bugs by adding safeguards to fail early on config initialization.
- Configuration options can be scattered around the codebase. **gofig** centralizes all configuration options in one place. This makes it easy to see all the configuration options at a glance.

In **gofig**, the code *is* the documentation.

## Goals
1. It should help with self-documenting app configuration options.
1. Config values should be basically immutable. Once you've initialized your configuration, you can't change it. 
1. It should fail early. Misconfigurations should result in the program crashing before starting it's real work.   

## Quick Usage Summary
- `gofig.InitOpt` is a struct used to define a particular configuration option.
    ```go
    type InitOpt struct {
        Name        string // The name of the config option (e.g. "ENV_VAR_A")
        Description string // A description of the config option
        Type        GfType // The type of the config option (e.g. TypeBool, TypeInt, TypeFloat, TypeString)
        Required    bool   // Whether the config option is required
        Default     any    // The default value of the config option. Doesn't do anything if the config option is required.
        IdPtr       *Id    // Pointer to the Id of the config option. This is where you store the Id. The Id value be set after the call to Init.
    }
    ```
- `gofig.Init` is a function that initializes a `gofig` an array of `gofig.InitOpt`s passed in.
    - All of your `gofig.Id`s will be set to their computed values after initialization and will be ready to go.
- `gofig.Get` is a function that retrieves the value of a configuration option given a `gofig.Id`.
 

## Demonstration
Below is a quick demo of how to use the **gofig**.
```go
func whatever() error {
    // --------------------------------
    // Assume env vars are set as follows:
    // export FOO=true
    // export BAR=42
    // export BAZ="I am a baz"
    // --------------------------------

    // Declare Ids
    var fooId gofig.Id
    var barId gofig.Id
    var bazId gofig.Id


    // Initialize gofig
    gf, err := gofig.Init(
	  gofig.InitOpt{
	  	Name:        "FOO",
	  	Description: "This is a foo. It is used for blah",
	  	Type:        gofig.TypeBool,
	  	Required:    true,
	  	IdPtr:       &fooId,
	  },
	  gofig.InitOpt{
	  	Name:        "BAR",
	  	Description: "This is a bar. It is used for blah",
	  	Type:        gofig.TypeInt,
	  	Required:    true,
	  	IdPtr:       &barId,
	  },
	  gofig.InitOpt{
	  	Name:        "BAZ",
	  	Description: "This is a baz. It is used for blah",
	  	Type:        gofig.TypeString,
	  	Required:    false,
	  	Default:     "I am a default value",
	  	IdPtr:       &bazId,
	  },
	)
    if err != nil {
        log.Fatalf("Failed to initialize configuration: %v", err)
    }

    // Get values
    foo, err := gf.GetBool(fooId) // foo = true
    if err != nil {
        return err
    }
    bar, err := gf.GetInt(barId) // bar = 42
    if err != nil {
        return err
    }
    baz, err := gf.GetString(bazId) // baz = "I am a baz"
    if err != nil {
        return err
    }
} 
```

In practice, you might choose to make your config global.


For additional examples, see the [example](example) and [test](test) directory. 


## Future Features
1. Support for complex validation rules of config values
    - For example:
      - if you have a config option that is supposed to be an email address, you can add a validation rule that checks if the value is in an email address format
      - if one config option is dependent on another, it will cause an error
      - if two config options are mutually exclusive, it will cause an error
1. Support for configuration from flags
1. Support for configuration from YAML files
1. Support for comma-separated values and arrays
    - Comma-separated for environment variables and command-line args
    - Arrays for JSON/YAML files

## FAQ
1. Why are you using `gofig.Id` instead of just using the name of the configuration option?
    - This is to prevent typos and usage of raw strings. Lots of raw strings means lots of find-and-replacing. Using `gofig.Id` will prevent this.
    - This forces the user to declare a variable to be used as an id.
1. Why have `Get` and `GetType` family of functions instead of just having one `Get` function that returns an `interface{}`/`any`?
    - This is for convenience as you will likely have to convert the returned `any` value to the correct type anyway.