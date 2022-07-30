package modules

import (
	"reflect"

	"github.com/anonyindian/gotgproto/dispatcher"
	"github.com/anonyindian/logger"
)

type module struct {
	Logger *logger.Logger
}

func Load(l *logger.Logger, dispatcher *dispatcher.CustomDispatcher) {
	l = l.Create("MODULES")
	defer l.ChangeLevel(logger.LevelMain).Println("LOADED")
	Type := reflect.TypeOf(&module{l})
	Value := reflect.ValueOf(&module{l})
	for i := 0; i < Type.NumMethod(); i++ {
		Type.Method(i).Func.Call([]reflect.Value{Value, reflect.ValueOf(dispatcher)})
	}
}
