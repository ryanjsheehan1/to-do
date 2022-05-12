package main

import (
	"bufio"
	"flag"
	"fmt"

	"io"
	"os"

	"strings"

	"github.com/ryanjsheehan1/todo"
)

// Default file name
var todoFileName = ".todo.json"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "ToDo Command-Line App\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage information:\n")
		flag.PrintDefaults()
	}

	// Parsing command line flags
	add := flag.Bool("add", false, "Items to be included in the ToDo list. Items are added as command line arguments or standard input.")
	list := flag.Bool("list", false, "List all items")
	complete := flag.Int("complete", 0, "Items to be completed")
	delete := flag.Int("del", 0, "Items to be deleted")
	verbose := flag.Bool("v", false, "Verbose output")
	incomplete := flag.Bool("i", false, "Display incomplete items only")

	flag.Parse()

	// Check if the user defined the ENV VAR for a custom file name
	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	// Define an items list
	l := &todo.List{}

	// Use the Get method to read ToDo items from file
	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Decide what to do based on the flags provided
	switch {
	case *list && *verbose && *incomplete:
		fmt.Print(l.Incomplete(true))
	case *list && *verbose:
		// List all ToDo items with verbose
		fmt.Print(l.Verbose())
	case *list && *incomplete:
		// List incomplete ToDo items
		fmt.Print(l.Incomplete(false))
	case *list:
		// List all ToDo items
		fmt.Print(l)
	case *complete > 0:
		// Complete the given item
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *delete > 0:
		// Delete the given item
		if err := l.Delete(*delete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *add:
		// When any arguments (excluding flags) are provided,
		// they will be used as the new task
		t, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		l.Add(t)

		// Save the new lsit
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		// Invalid flag provided
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}

// getTask function decides where to get the description
// for a new task from: arguments or STDIN
func getTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	s := bufio.NewScanner(r)
	s.Scan()
	if err := s.Err(); err != nil {
		return "", err
	}

	if len(s.Text()) == 0 {
		return "", fmt.Errorf("task cannot be blank")
	}

	return s.Text(), nil
}
