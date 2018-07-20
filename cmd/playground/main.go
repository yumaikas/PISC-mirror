package main

import (
	"pisc"

	"github.com/gopherjs/gopherjs/js"
)

func log(message string) {
	js.Global.Get("console").Call("log", message)
}

func log_error(message string) {
	js.Global.Get("console").Call("error", message)
}

func main() {

	m := &pisc.Machine{
		Values:               make([]pisc.StackEntry, 0),
		DefinedWords:         make(map[string]*pisc.CodeQuotation),
		DefinedStackComments: make(map[string]string),
		PredefinedWords:      make(map[string]pisc.GoWord),
		PrefixWords:          make(map[string]*pisc.CodeQuotation),
		HelpDocs:             make(map[string]string),
	}
	m.LoadModules(append(
		pisc.StandardModules,
		pisc.ModDebugCore,
		pisc.ModIOCore,
		pisc.ModHelpCode,
		ModPlayground)...)

	printStack := func(this *js.Object, arguments []*js.Object) interface{} {
		for _, val := range m.Values {
			printVal := val.String() + " <" + val.Type() + ">"
			log(printVal)
		}
		return js.Undefined
	}

	getStack := func(this *js.Object, arguments []*js.Object) interface{} {
		var jsVals = make([]*js.Object, 0)
		for _, val := range m.Values {
			jsVals = append(jsVals, js.MakeWrapper(val))
		}
		return jsVals
	}

	eval := func(this *js.Object, arguments []*js.Object) interface{} {
		if len(arguments) < 1 {
			log("Need a string to eval")
		}
		code := arguments[0].String()
		err := m.ExecuteString(code, pisc.CodePosition{Source: "User Input"})
		if err != nil {
			log_error(err.Error())
			return js.MakeWrapper(err)
		}
		/*
			for idx, val := range m.Values {
				printVal := val.String() + "<" + val.Type() + "> : " + strconv.Itoa(idx)
				log(printVal)
			}
		*/
		return js.Undefined
	}

	js.Global.Set("pisc_eval", js.MakeFunc(eval))
	js.Global.Set("pisc_show_stack", js.MakeFunc(printStack))
	js.Global.Set("pisc_get_stack", js.MakeFunc(getStack))
}
