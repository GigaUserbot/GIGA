# Adding Modules

## Quick Summary
This document will guide you on, How to add a module in the GIGA userbot. GIGA is a Go-based telegram userbot and Go is a compiled language, so it is unlikely that we add a module to GIGA during runtime. Instead, we can add up the `{module_name}.go` file to the source code and build it again. 

Please scroll down to know more about the modules of GIGA.

## How does the module system work in GIGA?

In order to make sure the automatic behaviour of loading modules, we use `reflect` package from the go standard packages to fetch all the methods of a struct named `module` located in the the [modules](./modules) directed inside `modules.go` file.


## How to load a module in GIGA?
A new module can be loaded to GIGA by defining a new method to `module` struct corresponding to the module being loaded.

### Example:
Let "Foo" be the name of the module, we want to import in the userbot. Hence, we will create a file named `foo.go` in the [modules](./modules) directory and define a method named `LoadFoo(dispatcher *dispatcher.CustomDispatcher)`

All the handlers will be allocated to the main dispatcher stream through this `LoadFoo` method.

```go
package modules

import (
    "github.com/anonyindian/gotgproto/dispatcher"
)

func (*module) LoadFoo(dp *dispatcher.CustomDispatcher) {
    // Handlers will be allocated to dispatcher through this function.
}
```

## What are the capabilities of `module` struct?
The `module` struct contains the basic needs of a module like `dispatcher`, `logger`, etc.

Here is the representation of `module` struct as Go code:
```go
type module struct {
	Logger *logger.Logger
}
```
Dispatcher is not included in the fields because it is passed through the method calls as we have seen above in the code snippet of loading a module.

### Creating logger for a module:
You can create logger for a module using the `Logger` field of `module` struct.

```go
var log = m.Logger.Create("FOO")
```

--- 
TO BE COMPLETED LATER

---