# gofig

## What's the Deal With gofig?
**gofig** is a **configuration library** for **config-heavy apps**. This was created because:
- As projects grow, configuration options grow and change. **gofig** makes it really hard to **not** document configuration options accurately and comprehensively. 
- **gofig** tries to prevent configuration related bugs by adding safeguards to fail early on config initialization.
- Configuration options can be scattered around the codebase. **gofig** centralizes all configuration options in one place. This makes it easy to see all the configuration options at a glance.

In **gofig**, the code *is* the documentation.

## Design Philosophies
1. It should be easy to use. It basically only has two API functions: `Init` and the `Get`-family of functions.
1. It should help with self-documenting app configuration.
1. It should be basically immutable. Once you've initialized your configuration, you can't change it. 
1. It should fail early. Misconfigurations should result in the program crashing before starting it's real work.   

## Demonstration
Below is a quick demo of how to use the **gofig**.
```go
func whatever() error {
    // Declare
    gf := gofig.Gofig{}

    // Initialization
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
			Required:    false,
			Default:     "I am a default value",
			IdPtr:       &bazId,
		},
	)
    if err != nil {
        log.Fatalf("Failed to initialize configuration: %v", err)
    }

    // Getting
    foo, err := gf.Get(fooId)
    if err != nil {
        return err
    }
    bar, err := gf.Get(barId)
    if err != nil {
        return err
    }
    baz, err := gf.Get(bazId)
    if err != nil {
        return err
    }
} 
```

In practice, you might choose to make your config global. In which case it might look like this:
```go
var gf gofig.Gofig

func whatever() {
    // Initialization
    err := gf.Init(
        gofig.InitOpt{
            Name:        "FOO",
            Description: "This is a foo. It is used for blah blah blah",
            Type:        gofig.TypeBool,
            Required:    true,
            Default:     false,
            IdPtr:       &fooId,
        },
        ...
    )
    if err != nil {
        log.Fatalf("Failed to initialize configuration: %v", err)
    }
} 

func idk() {
    // Getting
    foo, err := gf.Get(fooId).(bool)
    if err != nil {
        log.Fatalf("Failed to get foo: %v", err)
    }
    ...
}
```

For additional examples, see the [test](test) directory. This has a comprehensive set of test cases that can serve as an example. With that being said, there isn't much to demonstrate as the library is quite simple right now.

Use gofig and you'll find adding configuration options to config-heavy apps was never this easy. Go figure.

## Future Features
1. Support for complex validation rules of config values
  - For example:
    - if you have a config option that is supposed to be an email address, you can add a validation rule that checks if the value is in an email address format
    - if one config option is dependent on another, it will cause an error
    - if two config options are mutually exclusive, it will cause an error
1. Support for configuration from flags
1. Support for configuration from YAML files
