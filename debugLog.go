package pisc

import "fmt"

var xLOG = false

// This function is strictly for logging purposes.
func debugLog(args ...interface{}) {
	if xLOG {
		fmt.Println(args)
	}
}
