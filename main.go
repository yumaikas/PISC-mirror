// Posisition
// Independent
// Source
// Code
package main

import (
	"io"
	// "flag" TODO: Implement flags for file and burst modes
	"fmt"
	"os"
	"bytes"
	"strconv"
	"syscall"

	"log"

	"runtime/pprof"

	cli "gopkg.in/urfave/cli.v1"
	"github.com/pkg/term/termios"  // just for portable Tc[gs]etattr
)

const (
	esc = "\x1b"
)

var (
	col byte  // column last typed at on screen
)

// This function starts an interpertor
func main() {
	app := cli.NewApp()
	app.Author = "Andrew Owen, @yumaikas"
	app.Name = "PISC, aka Posisition Independent Source Code"
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
			Name:  "database, d",
			Usage: "Tells PISC to attach to the boltDB file at .piscdb",
		},
		cli.StringFlag{
			Name:  "file, f",
			Usage: "Execute a file as a bit of pisc, runs before -i or -c",
		},
		cli.BoolFlag{
			Name:   "benchmark",
			Hidden: true,
			Usage:  "Run various benchmarks, using pprof, and print out pertinent information",
		},
	}
	app.Run(os.Args)
}

func initMachine() *machine {
	m := &machine{
		values:               make([]stackEntry, 0),
		definedWords:         make(map[string]*codeQuotation),
		definedStackComments: make(map[string]string),
		predefinedWords:      make(map[string]GoWord),
		prefixWords:          make(map[string]*codeQuotation),
		helpDocs:             make(map[string]string),
	}
	return m
}

func handleFlags(ctx *cli.Context) {
	m := initMachine()
	// Execute this before benchmarking since we aren't yet benchmarking file loads
	if ctx.IsSet("benchmark") {
		err := m.loadForCLI()
		if err != nil {
			log.Fatalf("Unable to start benchmark due to error %v", err.Error())
			return
		}
		err = m.executeString(`"factorial.pisc" import`, codePosition{source: "pre-benchmark import"})
		if err != nil {
			log.Fatalf("Unable to start benchmark due to error %v", err.Error())
			return
		}
		f, err := os.Create("bench-cpu-recursion.prof")
		if err != nil {
			log.Fatal("Unable to create profiling file")
			return
		}
		pos := codePosition{source: "Benchmark recursive"}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("Unable to start CPU profile")
		}
		err = m.executeString("100000 [ 12 factorial drop ] times", pos)
		if err != nil {
			log.Fatal("Recursive benchmark failed:", err)
		}
		pprof.StopCPUProfile()
		f, err = os.Create("bench-cpu-iteration.prof")
		if err != nil {
			log.Fatal("Unable to create profiling file")
			return
		}
		pos = codePosition{source: "Benchmark loop"}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("Unable to start CPU profile")
			return
		}
		err = m.executeString("100000 [ 12 factorial-loop drop ] times", pos)
		if err != nil {
			log.Fatal("Recursive benchmark failed:", err)
			pprof.StopCPUProfile()
			return
		}
		pprof.StopCPUProfile()
		return
	}
	// Load PISC with libraries, according to the context
	if ctx.IsSet("file") || ctx.IsSet("command") || ctx.IsSet("interactive") {
		if ctx.IsSet("database") {
			err := m.loadForDB()
			if err != nil {
				fmt.Println(err.Error())
				log.Fatal("Error while loading modules")
			}
		} else {
			err := m.loadForCLI()
			if err != nil {
				fmt.Println(err.Error())
				log.Fatal("Error while loading modules")
			}
		}
		m.logAndResetDispatchCount(os.Stderr)
	}
	if ctx.IsSet("file") {
		m.pushValue(String(ctx.String("file")))
		err := m.executeString("import", codePosition{
			source: "argument line",
		})
		if err != nil {
			log.Println(err)
			log.Fatal("Error running file")
		}
		m.logAndResetDispatchCount(os.Stderr)
	}
	if ctx.IsSet("command") {
		line := ctx.String("command")
		p, err := stringToQuotation(line, codePosition{source: "args"})
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

func loadInteractive(m *machine) {

	// given_files := flag.Bool("f", false, "Sets the rest of the arguments to list of files")
	// Run command stuff here.

	in, err := syscall.Open("/dev/tty", syscall.O_RDONLY, 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error opening terminal")
		return
	}
	defer syscall.Close(in)

	fmt.Fprintln(
		os.Stderr,
		`Postion
Independent
Source
Code`)

	var orig_tios syscall.Termios
	err = termios.Tcgetattr(uintptr(in), &orig_tios)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error saving terminal settings")
		return
	}
	defer termios.Tcsetattr(uintptr(in), termios.TCSAFLUSH, &orig_tios)

	// minimal switch to raw mode
	tios := orig_tios
	tios.Lflag &^= syscall.ICANON
	termios.Tcsetattr(uintptr(in), termios.TCSAFLUSH, &tios)

	// interactive loop
	numEntries := 0
	for {
		fmt.Print(">> ");  col = 3
		line, err := readCommand(in, m)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		numEntries++
		// fmt.Println(words)

		err = m.executeString(line, codePosition{source: fmt.Sprint("stdin:", numEntries)})
		if err == ExitingProgram {
			fmt.Fprintln(os.Stderr, "Exiting program")
			return
		}
		if err != nil {
			fmt.Println("Error:")
			fmt.Println(err.Error())
			return
		}
		m.logAndResetDispatchCount(os.Stderr)
		fmt.Fprintln(os.Stderr, "Data Stack:")
		for _, val := range m.values {
			fmt.Println(val.String(), fmt.Sprint("<", val.Type(), ">"))
		}
	}
}

func readCommand(in int, m *machine) (string, error) {
	buf := new(bytes.Buffer)
	for {
		b, err := readByte(in)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error reading")
			return "", err
		}
		if b[0] == /*ctrl-d*/4 {
			return "", io.EOF  // exit
		}
		if b[0] == /*ctrl-u*/21 {
			clear_line()
			start_of_line()
			fmt.Print(">> ");  col = 3
			buf.Truncate(0)
			continue
		}
		col++  // messes up on non-printing characters
		if b[0] == /*backspace*/127 {
			if buf.Len() > 0 {
				left(2)  // undo print of backspace code
				left(1)  // perform backspace
				clear_right()
				col -= 2  // undo backspace and previous character
				buf.Truncate(buf.Len()-1)
			}
		}
		if strconv.IsPrint(rune(b[0])) {
			buf.Write(b)
		}
		if b[0] == '\r' || b[0] == '\n' {
			break
		}
		execute_print_and_restore_cursor(*m, buf.String())
	}
	return buf.String(), nil
}

func execute_print_and_restore_cursor(m_copy machine, partial string) {
	m_copy.executeString(partial+"\n", codePosition{})
	var temporary_lines byte = 0
	start_of_next_line();  temporary_lines++
	clear_screen_below()
	fmt.Println("Data Stack:");  temporary_lines++
	for _, val := range m_copy.values {
		fmt.Println(val.String(), fmt.Sprint("<", val.Type(), ">"));  temporary_lines++
	}
	up(temporary_lines);  right(col)
}

func readByte(in int) ([]byte, error) {
	b := make([]byte, 1)
	_, err := syscall.Read(in, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// some vt100 escape sequences (http://www.uqac.ca/flemieux/PRO100/VT100_Escape_Codes.html)

func move_cursor(x, y int) {
	os.Stdout.Write([]byte(fmt.Sprint(esc,"[",x,";",y,"H")))
}

// 256-color mode
func color(c byte) {
	os.Stdout.Write([]byte(fmt.Sprint(esc,"[38;5;",c,"m")))
}
func bg(c byte) {
	os.Stdout.Write([]byte(fmt.Sprint(esc,"[48;5;",c,"m")))
}

func up(n byte) {
	os.Stdout.Write([]byte(fmt.Sprint(esc,"[",n,"A")))
}

func down(n byte) {  // doesn't scroll
	os.Stdout.Write([]byte(fmt.Sprint(esc,"[",n,"B")))
}

func left(n byte) {
	os.Stdout.Write([]byte(fmt.Sprint(esc,"[",n,"D")))
}

func right(n byte) {
	os.Stdout.Write([]byte(fmt.Sprint(esc,"[",n,"C")))
}

func clear_right() {
	os.Stdout.Write([]byte(fmt.Sprint(esc,"[K")))
}

func clear_line() {
	os.Stdout.Write([]byte(fmt.Sprint(esc,"[2K")))
}

func start_of_next_line() {
	os.Stdout.Write([]byte("\n"))
}

func start_of_line() {
	start_of_next_line()
	up(1)
}

func clear_screen() {
	os.Stdout.Write([]byte(fmt.Sprint(esc,"[2J")))
}

func clear_screen_below() {
	os.Stdout.Write([]byte(fmt.Sprint(esc,"[0J")))
}

// no good after scroll
func save_cursor() {
	os.Stdout.Write([]byte(fmt.Sprint(esc,"7")))
}
func restore_cursor() {
	os.Stdout.Write([]byte(fmt.Sprint(esc,"8")))
}
