package main

type stackEntry interface{}

type Boolean bool
type Integer int
type Double float64
type Dict map[string]stackEntry
type Array []stackEntry
type String string

// This is a separate sematic from arrays.
type quotation []word
