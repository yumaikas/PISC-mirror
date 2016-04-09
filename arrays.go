package main

func isVectorWord(w word) bool {
	return w == "at" ||
		w == "<vector>" ||
		w == "#>array" || // Take n elements from the stack into a new array
		w == "+>" || // Add to end of array
		w == "<+" || // Add to front of array
		w == "<-" || // Remove from from of array
		w == "->" || // Remove from end of array
		w == "#>" // || // Pop n elements from stack into end of array
	//w == "#->" || // Pop n elements from end array into stack
	//w == "each" // execute a quotation for each method on the array
	//
}
func (m *machine) executeVectorWord(w word) error {
	switch w {
	case "at":
		idx := m.popValue().(Integer)
		arr := m.popValue().(Array)
		m.pushValue(arr[idx])
	case "<vector>":
		m.pushValue(Array(make([]stackEntry, 0)))
	case "+>": // ( vec elem -- newVect )
		toAppend := m.popValue()
		arr := m.popValue().(Array)
		arr = append(arr, toAppend)
		m.pushValue(arr)
	case "<+": // ( elem vec -- newVect )
		toPrepend := m.popValue()
		arr := m.popValue().(Array)
		arr = append([]stackEntry{toPrepend}, arr...)
		m.pushValue(arr)
	case "->": // ( vec -- sliced elem )
		arr := m.popValue().(Array)
		val := arr[len(arr)-1]
		arr = arr[:len(arr)-1]
		m.pushValue(arr)
		m.pushValue(val)
	case "<-": // ( vec -- sliced elem )
		arr := m.popValue().(Array)
		val := arr[0]
		arr = arr[1:]
		m.pushValue(arr)
		m.pushValue(val)
	}
	return nil

}
