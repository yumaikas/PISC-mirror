package pisc

import (
	"fmt"
)

var ModVectorCore = Module{
	Author:    "Andrew Owen",
	Name:      "VectorCore",
	License:   "MIT",
	DocString: "Words for creating and manipulating vectors",
	Load:      loadVectorCore,
}

func loadVectorCore(m *Machine) error {
	m.AddGoWordWithStack("vec-set-at",
		" ( vec val idx -- vec ) ",
		"Set the element of vec at idx to val",
		vecSetAt)
	m.AddGoWordWithStack("vec-at",
		" ( vec idx -- elem ) ",
		"Get the element at idx in vec",
		vecAt)
	m.AddGoWordWithStack("<vector>",
		" ( -- vector )",
		"Construct an empty vector",
		makeVec)
	m.AddGoWordWithStack("vec-each",
		" ( vec quot -- .. ) ",
		"Execute quot for each element in vec",
		vecEach)
	m.AddGoWordWithStack("vec-append",
		" ( vec elem -- vec ) ",
		"Push elem to the end of vec, leaving vec on the stack",
		vecAppend)
	m.AddGoWordWithStack("vec-push",
		" ( vec elem -- )",
		"Push elem to the end of vec, taking vec off the stack",
		vecPush)
	// TODO: Consider removing vec-pushback
	m.AddGoWordWithStack("vec-pushback",
		" ( vec elem -- ) ",
		"Same as vec-push",
		vecPush)
	m.AddGoWordWithStack("vec-prepend",
		" ( vec elem -- vec ) ",
		"Push elem to the front of vec, leaving vec on the stack",
		vecPrepend)
	m.AddGoWordWithStack("vec-popback",
		" ( vec -- vec elem ) ",
		"Pop elem of the end of vec, leaving both on the stack",
		vecPopback)
	m.AddGoWordWithStack("vec-popfront",
		" ( vec -- vec elem ) ",
		"Pop elem off the front of vec, leaving both on the stack",
		vecPopfront)
	// TODO: vec-slice ( vec begin end -- sliced-vec )
	// Return success if we can load out PISC file as well.
	return m.ImportPISCAsset("stdlib/vectors.pisc")
}

func vecSetAt(m *Machine) error {
	idx := int(m.PopValue().(Integer))
	val := m.PopValue()
	elems := m.Values[len(m.Values)-1].(*Vector).Elements
	if idx > len(elems)-1 || idx < 0 {
		return fmt.Errorf("index out of bounds: %v", idx)
	}
	elems[idx] = val
	return nil
}

func vecAt(m *Machine) error {
	idx := int(m.PopValue().(Integer))
	elems := m.PopValue().(*Vector).Elements
	if idx > len(elems)-1 || idx < 0 {
		return fmt.Errorf("index out of bounds: %v", idx)
	}
	m.PushValue(elems[idx])
	return nil
}

func makeVec(m *Machine) error {
	vec := &Vector{Elements: make([]StackEntry, 0)}
	m.PushValue(vec)
	return nil
}

func vecEach(m *Machine) error {
	quot := m.PopValue().(*Quotation)
	vect := m.PopValue().(*Vector)
	for _, elem := range vect.Elements {
		m.PushValue(elem)
		m.PushValue(quot)
		err := m.ExecuteQuotation()
		if IsLoopError(err) && LoopShouldEnd(err) {
			return nil
		}
		if IsLoopError(err) && !LoopShouldEnd(err) {
			continue
		}
		if err != nil {
			fmt.Println("ERROR HAPPENED")
			return err
		}
	}
	return nil
}

func vecPush(m *Machine) error {
	toAppend := m.PopValue()
	vec := m.PopValue().(*Vector)
	newElems := append(vec.Elements, toAppend)
	vec.Elements = newElems
	return nil
}

func vecAppend(m *Machine) error {
	toAppend := m.PopValue()
	vec := m.PopValue().(*Vector)
	newElems := append(vec.Elements, toAppend)
	vec.Elements = newElems
	m.PushValue(vec)
	return nil
}

func vecPrepend(m *Machine) error {
	toPrepend := m.PopValue()
	arr := m.PopValue().(*Vector)
	newElems := append([]StackEntry{toPrepend}, (arr.Elements)...)
	arr.Elements = newElems
	m.PushValue(arr)
	return nil
}

func vecPopback(m *Machine) error {
	vec := m.PopValue().(*Vector)
	if len(vec.Elements) < 1 {
		return fmt.Errorf("no elements in array")
	}
	val := (vec.Elements)[len(vec.Elements)-1]
	vec.Elements = vec.Elements[:len(vec.Elements)-1]
	m.PushValue(vec)
	m.PushValue(val)
	return nil
}

func vecPopfront(m *Machine) error {
	vec := m.PopValue().(*Vector)
	if len(vec.Elements) < 1 {
		return fmt.Errorf("no elements in array")
	}
	val := vec.Elements[0]
	vec.Elements = vec.Elements[1:]
	m.PushValue(vec)
	m.PushValue(val)
	return nil
}
