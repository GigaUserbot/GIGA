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

## How to write new handler responses?

Handlers are mapped functions which are called when a new telegram update is processed. Though we may not need to call all of them on every update so there exist some precoditional ones and filters for the purpose.

For example, to create a handler for handling a command named "foo", having usage like `.foo`, we will use command handler as follows:
```go
dp.AddHandler(handlers.NewCommand("foo", responseFunction))
```
where responseFunction is function which would look like:
```go
func responseFcuntion(ctx *ext.Context, u *ext.Update) error {
    return nil
}
```

There are many such conditional handlers provided by the GoTGProto package.

## How to add help section for a module?
In order to add help section for your module, you have to call `SetModuleHelp` function of `helpmaker` package in the main module loading function.

### Example:
```go
package modules

import (
    "github.com/anonyindian/gotgproto/dispatcher"
)

func (*module) LoadFoo(dp *dispatcher.CustomDispatcher) {
    var l = m.Logger.Create("FOO")
	defer l.ChangeLevel(logger.LevelInfo).Println("LOADED")
    // foo is the module name and second argument is the help string corresponding to it.
    helpmaker.SetModuleHelp("foo", "help here")
}
```

## Important Points
- **ALWAYS CREATE** a logger and "LOADED" print line statement everytime you create a new module.
- **ALWAYS ADD** message handlers to unique handler groups.
- **NEVER USE** `dispatcher.EndGroups` in return statement of message handlers.
- **NEVER ADD** any other argument to Load method of the module than the `*dispatcher.CustomDispatcher` 

## Example Module
Here is an example module to understand the concept better, named as "Foo":

```go
package modules

import (
	"fmt"

	"github.com/anonyindian/gotgproto/dispatcher"
    "github.com/anonyindian/gotgproto/dispatcher/handlers"
	"github.com/anonyindian/gotgproto/dispatcher/handlers/filters"
	"github.com/anonyindian/gotgproto/ext"
	"github.com/anonyindian/logger"
	"github.com/gigauserbot/giga/bot/helpmaker"
)

func (m *module) LoadFoo(dispatcher *dispatcher.CustomDispatcher) {
	var l = m.Logger.Create("FOO")
	defer l.ChangeLevel(logger.LevelInfo).Println("LOADED")
	helpmaker.SetModuleHelp("foo", `
	Foo is an example module.
	
	<b>Commands</b>:
	 • <code>.foo</code>: Example command for foo module   
`)
	// authorised() decorator ensures that the command is not usabe by anyone
	// other than the authorised users i.e. you and sudos.
	dispatcher.AddHandler(handlers.NewCommand("foo", authorised(foo)))
	dispatcher.AddHandlerToGroup(handlers.NewMessage(filters.Message.All, fooMessage), 12)
}

func foo(ctx *ext.Context, u *ext.Update) error {
	fmt.Println("You used .foo command")
	return dispatcher.EndGroups
}

func fooMessage(ctx *ext.Context, u *ext.Update) error {
	fmt.Println("You got a new message update")
	// NEVER use EndGroups in a message handler response,
	// use nil or ContinueGroups instead.
	return nil
}

```