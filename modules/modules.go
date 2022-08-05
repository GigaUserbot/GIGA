package modules

import (
	"reflect"

	"github.com/anonyindian/gotgproto/dispatcher"
	"github.com/anonyindian/gotgproto/dispatcher/handlers"
	"github.com/anonyindian/gotgproto/ext"
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

func authorised(cback handlers.CallbackResponse) handlers.CallbackResponse {
	return func(ctx *ext.Context, u *ext.Update) error {
		if u.EffectiveMessage.Out {
			return cback(ctx, u)
		}
		return dispatcher.EndGroups
	}
}

// func authorisedMessage(cback handlers.CallbackResponse) handlers.CallbackResponse {
// 	return func(ctx *ext.Context, u *ext.Update) error {
// 		if u.EffectiveUser() != nil && u.EffectiveUser().ID == gotgproto.Self.ID {
// 			return cback(ctx, u)
// 		}
// 		return nil
// 	}
// }
