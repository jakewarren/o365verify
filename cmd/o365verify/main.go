package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/jakewarren/o365verify"
	"github.com/remeh/sizedwaitgroup"
	"github.com/spf13/pflag"
)

var (
	// build information set by ldflags
	appName    = "o365verify"
	appVersion = "(ﾉ☉ヮ⚆)ﾉ ⌒*:･ﾟ✧"
	commit     = "(ﾉ☉ヮ⚆)ﾉ ⌒*:･ﾟ✧"
	buildDate  = "(ﾉ☉ヮ⚆)ﾉ ⌒*:･ﾟ✧"
)

type config struct {
	Threads int
}

const usageMessage = `Usage: o365verify [flags] <email address...>`

func main() {
	var c config

	pflag.Usage = func() {
		_, _ = fmt.Fprintln(os.Stderr, usageMessage)
		_, _ = fmt.Fprintln(os.Stderr, "")
		_, _ = fmt.Fprintln(os.Stderr, "Flags:")
		pflag.PrintDefaults()
		_, _ = fmt.Fprintln(os.Stderr, "")
		_, _ = fmt.Fprintln(os.Stderr, "URL: https://github.com/jakewarren/o365verify")
	}

	displayHelp := pflag.BoolP("help", "h", false, "display help")
	displayVersion := pflag.BoolP("version", "V", false, "display version information")

	pflag.IntVarP(&c.Threads, "threads", "t", 10, "number of threads to run with")
	pflag.Parse()

	// override the default usage display
	if *displayHelp {
		pflag.Usage()
		os.Exit(0)
	}

	if *displayVersion {
		fmt.Printf(`%s:
    version     : %s
    git hash    : %s
    build date  : %s 
    go version  : %s
    go compiler : %s
    platform    : %s/%s
`, appName, appVersion, commit, buildDate, runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	if pflag.NArg() == 0 {
		log.Fatalln("format string not provided")
	}

	swg := sizedwaitgroup.New(c.Threads)
	var resultWg sync.WaitGroup

	answers := make(chan *o365verify.Result)
	results := make([]*o365verify.Result, pflag.NArg())

	for _, e := range pflag.Args() {
		swg.Add()
		go func(e string) {
			defer swg.Done()
			r, err := o365verify.VerifyAddress(e)
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
				return
			}
			answers <- r
			resultWg.Add(1)
		}(e)

	}

	go func() {
		i := 0
		for r := range answers {
			results[i] = r
			i++
			resultWg.Done()
		}
	}()

	swg.Wait()
	close(answers)
	resultWg.Wait()

	writeJSON(results)
}

func writeJSON(r []*o365verify.Result) {
	out, _ := json.MarshalIndent(r, "", "  ")
	fmt.Println(string(out))
}
