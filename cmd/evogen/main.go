package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/jncornett/beans-engine/evo/cli"
	"github.com/jncornett/beans-engine/evo/genome"
)

const (
	defaultLength         = 10
	defaultLengthVariance = 0.2
	defaultEncoding       = cli.EncodingEvo
)

func main() {
	var args struct {
		Length         uint
		LengthVariance float64
		Format         cli.Encoding `help:"output format (one of {json,evo,evox})"`
		Filename       string       `arg:"positional"`
	}
	args.Length = defaultLength
	args.LengthVariance = defaultLengthVariance
	args.Filename = cli.StdioFilename
	arg.MustParse(&args)
	log.Printf("%+v\n", args)
	if args.Filename == cli.StdioFilename && args.Format == cli.EncodingNone {
		args.Format = defaultEncoding
	}
	rand.Seed(time.Now().Unix())
	n := randomLength(args.Length, args.LengthVariance)
	code := genome.SampleN(genome.Default, n)
	if err := cli.Save(args.Filename, args.Format, code); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}

func randomLength(hint uint, variance float64) int {
	half := float64(hint) * variance
	min := float64(hint) - half
	if min < 0 {
		min = 0
	}
	return int(min) + rand.Intn(int(2*half))
}
