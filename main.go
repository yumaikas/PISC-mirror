// Position
// Independent
// Source
// Code
package pisc

// "flag" TODO: Implement flags for file and burst modes

/*
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


func initMachine() *Machine {
	m := &Machine{
		Values:               make([]StackEntry, 0),
		DefinedWords:         make(map[string]*CodeQuotation),
		DefinedStackComments: make(map[string]string),
		PredefinedWords:      make(map[string]GoWord),
		PrefixWords:          make(map[string]*CodeQuotation),
		HelpDocs:             make(map[string]string),
	}
	return m
}

func benchmark(m *Machine) {
	err := m.loadForCLI()
	if err != nil {
		log.Fatalf("Unable to start benchmark due to error %v", err.Error())
		return
	}
	err = m.ExecuteString(`"factorial.pisc" import`, CodePosition{Source: "pre-benchmark import"})
	if err != nil {
		log.Fatalf("Unable to start benchmark due to error %v", err.Error())
		return
	}
	f, err := os.Create("bench-cpu-recursion.prof")
	if err != nil {
		log.Fatal("Unable to create profiling file")
		return
	}
	pos := CodePosition{Source: "Benchmark recursive"}
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
	pos = CodePosition{Source: "Benchmark loop"}
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

func handleFlags(ctx *cli.Context) {
	m := initMachine()
	fmt.Println("???")
	// Execute this before benchmarking since we aren't yet benchmarking file loads
	if ctx.IsSet("benchmark") {
		benchmark(m)
	}
	// Load PISC with libraries, according to the context
	if ctx.IsSet("file") || ctx.IsSet("command") || ctx.IsSet("interactive") {
		if ctx.IsSet("chatbot") {
			err := m.loadForChatbot()
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				log.Fatal("Error while loading modules")
			}
		}
		if ctx.IsSet("boltdb") {
			err := m.loadForDB()
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				log.Fatal("Error while loading modules")
			}
		} else {
			err := m.loadForCLI()
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				log.Fatal("Error while loading modules")
			}
		}
		m.logAndResetDispatchCount(os.Stderr)
	}
	if ctx.IsSet("file") {
		m.PushValue(String(ctx.String("file")))
		err := m.ExecuteString("import", CodePosition{
			Source: "argument line",
		})
		if err != nil {
			log.Println(err)
			log.Fatal("Error running file")
		}
		m.logAndResetDispatchCount(os.Stderr)
	}
	if ctx.IsSet("command") {
		line := ctx.String("command")
		p, err := stringToQuotation(line, CodePosition{Source: "args"})
		if err != nil {
			log.Fatal("Error in command: ", err)
		}
		err = m.execute(p)
		if err != nil {
			log.Fatal("Error in command: ", err)
		}
		m.logAndResetDispatchCount(os.Stderr)
	}
	if ctx.IsSet("interactive") {
		loadInteractive(m)
	}
}

func loadInteractive(m *Machine) {

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

		err = m.ExecuteString(line, CodePosition{Source: fmt.Sprint("stdin:", numEntries)})
		if err == ExitingProgram {
			fmt.Println("Exiting program")
			return
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err.Error())
			return
		}
		m.logAndResetDispatchCount(os.Stderr)
		fmt.Println("Data Stack:")
		for _, val := range m.Values {
			fmt.Println(val.String(), fmt.Sprint("<", val.Type(), ">"))
		}
	}

}
*/
