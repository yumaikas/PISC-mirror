package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"log"

	"runtime/pprof"

	"pisc"
	"pisc/libs/boltdb"

	"gopkg.in/readline.v1"
	cli "gopkg.in/urfave/cli.v1"
)

// This function starts an interpreter
func main() {
	app := cli.NewApp()
	app.Author = "Andrew Owen, @yumaikas"
	app.Name = "PISC, aka Position Independent Source Code"
	app.Usage = "A small stack based scripting langauge built for fun"
	app.Action = handleFlags
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "interactive, i",
			Usage: "Run the interactive version of PISC",
		},
		cli.StringFlag{
			Name:  "command, c",
			Usage: "Expressions to run from the command line, before -i, if it exists",
		},
		cli.BoolFlag{
			Name:  "boltdb, d",
			Usage: "Tells PISC to enable boltdb integration",
		},
		cli.StringFlag{
			Name:  "file, f",
			Usage: "Execute a file as a bit of pisc, runs before -i or -c",
		},
		cli.BoolFlag{
			Name:  "chatbot",
			Usage: "Load the chatbot modules before -c and -i",
		},
		cli.BoolFlag{
			Name:   "benchmark",
			Hidden: true,
			Usage:  "Run various benchmarks, using pprof, and print out pertinent information",
		},
	}
	app.Run(os.Args)
}

func initMachine() *pisc.Machine {
	m := &pisc.Machine{
		Values:               make([]pisc.StackEntry, 0),
		DefinedWords:         make(map[string]*pisc.CodeQuotation),
		DefinedStackComments: make(map[string]string),
		PredefinedWords:      make(map[string]pisc.GoWord),
		PrefixWords:          make(map[string]*pisc.CodeQuotation),
		HelpDocs:             make(map[string]string),
	}
	return m
}

func benchmark(m *pisc.Machine) {
	err := LoadForCLI(m)
	if err != nil {
		log.Fatalf("Unable to start benchmark due to error %v", err.Error())
		return
	}
	err = m.ExecuteString(`"factorial.pisc" import`, pisc.CodePosition{Source: "pre-benchmark import"})
	if err != nil {
		log.Fatalf("Unable to start benchmark due to error %v", err.Error())
		return
	}
	f, err := os.Create("bench-cpu-recursion.prof")
	if err != nil {
		log.Fatal("Unable to create profiling file")
		return
	}
	pos := pisc.CodePosition{Source: "Benchmark recursive"}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("Unable to start CPU profile")
	}
	err = m.ExecuteString("100000 [ 12 factorial drop ] times", pos)
	if err != nil {
		log.Fatal("Recursive benchmark failed:", err)
	}
	pprof.StopCPUProfile()
	f, err = os.Create("bench-cpu-iteration.prof")
	if err != nil {
		log.Fatal("Unable to create profiling file")
		return
	}
	pos = pisc.CodePosition{Source: "Benchmark loop"}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("Unable to start CPU profile")
		return
	}
	err = m.ExecuteString("100000 [ 12 factorial-loop drop ] times", pos)
	if err != nil {
		log.Fatal("Recursive benchmark failed:", err)
		pprof.StopCPUProfile()
		return
	}
	pprof.StopCPUProfile()
	return
}

func LoadForCLI(m *pisc.Machine) error {
	return m.LoadModules(append(pisc.StandardModules,
		pisc.ModIOCore, pisc.ModDebugCore, pisc.ModShellUtils)...)
}

func LoadForDB(m *pisc.Machine) error {
	return m.LoadModules(append(pisc.StandardModules,
		boltdb.ModBoltDB, pisc.ModIOCore, pisc.ModShellUtils)...)
}

func LoadForChatbot(m *pisc.Machine) error {
	return m.LoadModules(append(pisc.StandardModules,
		boltdb.ModBoltDB, pisc.ModIRCKit)...)
}

func handleFlags(ctx *cli.Context) {
	m := initMachine()
	// Execute this before benchmarking since we aren't yet benchmarking file loads
	if ctx.IsSet("benchmark") {
		benchmark(m)
	}
	// Load PISC with libraries, according to the context
	if ctx.IsSet("file") || ctx.IsSet("command") || ctx.IsSet("interactive") {
		if ctx.IsSet("chatbot") {
			err := LoadForChatbot(m)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				log.Fatal("Error while loading modules")
			}
		}
		if ctx.IsSet("boltdb") {
			err := LoadForDB(m)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				log.Fatal("Error while loading modules")
			}
		} else {
			err := LoadForCLI(m)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				log.Fatal("Error while loading modules")
			}
		}
		m.LogAndResetDispatchCount(os.Stderr)
	}
	if ctx.IsSet("file") {
		m.PushValue(pisc.String(ctx.String("file")))
		err := m.ExecuteString("import", pisc.CodePosition{
			Source: "argument line",
		})
		if err != nil {
			log.Println(err)
			log.Fatal("Error running file")
		}
		m.LogAndResetDispatchCount(os.Stderr)
	}
	if ctx.IsSet("command") {
		line := ctx.String("command")
		err := m.ExecuteString(line, pisc.CodePosition{Source: "args"})
		if err != nil {
			log.Fatal("Error in command: ", err)
		}
		m.LogAndResetDispatchCount(os.Stderr)
	}
	if ctx.IsSet("interactive") {
		loadInteractive(m)
	}
}

func loadInteractive(m *pisc.Machine) {

	// given_files := flag.Bool("f", false, "Sets the rest of the arguments to list of files")
	// Run command stuff here.

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          ">> ",
		HistoryFile:     "/tmp/readline.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Println(
		`Position
Independent
Source
Code`)
	numEntries := 0
	for {
		// fmt.Print(">> ")
		line, err := rl.Readline()
		if strings.TrimSpace(line) == "exit" {
			fmt.Println("Exiting")
			return
		}
		if err == io.EOF {
			fmt.Println("Exiting program")
			return
		}
		if err != nil {
			panic(err)
		}
		numEntries++
		// fmt.Println(words)

		err = m.ExecuteString(line, pisc.CodePosition{Source: fmt.Sprint("stdin:", numEntries)})
		if err == pisc.ExitingProgram {
			fmt.Println("Exiting program")
			return
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err.Error())
			return
		}
		m.LogAndResetDispatchCount(os.Stderr)
		fmt.Println("Data Stack:")
		for _, val := range m.Values {
			fmt.Println(val.String(), fmt.Sprint("<", val.Type(), ">"))
		}
	}

}