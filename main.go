package main

import (
	"fmt"
	"github.com/jojomi/pdfopt"
	"github.com/spf13/cobra"
	"math"
	"os"
	"strings"
)

type Style string

const (
	StyleScreen   Style = "screen"
	StyleEbook    Style = "ebook"
	StylePrint    Style = "print"
	StylePrepress Style = "prepress"
)

var (
	inplace      bool
	silent       bool
	style        Style = ""
	screen       bool
	ebook        bool
	printFlag    bool
	prepress     bool
	dpi          int
	styleWrapper = &styleValueWrapper{target: &style}
	rootCmd      = &cobra.Command{
		Use:   "pdfoptimize [input] [output]",
		Short: "Optimize PDF files",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 || len(args) > 2 {
				return fmt.Errorf("requires 1 or 2 arguments")
			}
			if inplace && len(args) > 1 {
				return fmt.Errorf("output argument not allowed with --inplace flag")
			}
			if style != "" && style != StyleScreen && style != StyleEbook && style != StylePrint && style != StylePrepress {
				return fmt.Errorf("style must be one of: screen, ebook, print, prepress")
			}
			return nil
		},
		Run: mainHandler,
	}
)

func init() {
	rootCmd.Flags().BoolVarP(&inplace, "inplace", "i", false, "optimize file in-place")
	rootCmd.Flags().BoolVarP(&silent, "silent", "q", false, "no output")
	rootCmd.Flags().VarP(styleWrapper, "style", "", "optimization style (screen, print, ebook)")
	rootCmd.Flags().BoolVarP(&printFlag, "print", "p", false, "optimize for print")
	rootCmd.Flags().BoolVarP(&ebook, "ebook", "e", false, "optimize for ebook")
	rootCmd.Flags().BoolVarP(&screen, "screen", "s", false, "optimize for screen")
	rootCmd.Flags().BoolVar(&prepress, "prepress", false, "optimize for prepress")
	rootCmd.Flags().IntVar(&dpi, "dpi", 0, "image DPI (0 = auto)")

	rootCmd.MarkFlagsMutuallyExclusive("style", "print", "ebook", "screen", "prepress")
}

type styleValueWrapper struct {
	target *Style
}

func (s *styleValueWrapper) String() string {
	if s.target != nil {
		return string(*s.target)
	}
	return ""
}

func (s *styleValueWrapper) Set(val string) error {
	switch Style(val) {
	case StyleScreen, StylePrint, StyleEbook, StylePrepress:
		if s.target != nil {
			*s.target = Style(val)
		}
		return nil
	default:
		return fmt.Errorf("style must be one of: screen, print, ebook")
	}
}

func (s *styleValueWrapper) Type() string {
	return "string"
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func mainHandler(cmd *cobra.Command, args []string) {
	input := args[0]
	output := input

	if style == "" {
		style = StyleScreen
	}
	if ebook {
		style = StyleEbook
	}
	if printFlag {
		style = StylePrint
	}
	if prepress {
		style = StylePrepress
	}

	if len(args) > 1 {
		output = args[1]
	} else if inplace {
		output = input
	} else {
		output = fmt.Sprintf("%s.%s.pdf", strings.TrimSuffix(input, ".pdf"), style)
	}

	p := pdfopt.NewPDFOpt(input)

	var (
	    err error
	    oldInfo os.FileInfo
    )
	if !silent {
		fmt.Printf("Optimizing %s to %s for %s...\n", input, output, style)

		oldInfo, err = os.Stat(input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading original file size:", err)
			os.Exit(1)
		}
	}

	var optimizeErr error

	switch {
	case style == StylePrint || printFlag:
		p.ForPrint()
	case style == StyleEbook || ebook:
		p.ForEbook()
	case style == StyleScreen || screen:
		p.ForScreen()
	case style == StylePrepress || prepress:
		p.ForPrepress()
	}

	if dpi > 0 {
		p.ImageDPI(dpi)
	}

	if inplace {
		optimizeErr = p.OptimizeInplace()
	} else {
		optimizeErr = p.Optimize(output)
	}

	if optimizeErr != nil {
		fmt.Fprintln(os.Stderr, "Error optimizing PDF:", optimizeErr.Error())
		os.Exit(1)
	}

	if !silent {
		newInfo, err := os.Stat(output)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading new file size:", err)
			os.Exit(1)
		}

		oldSize := oldInfo.Size()
		newSize := newInfo.Size()
		reduction := float64(oldSize-newSize) / float64(oldSize) * 100

		fmt.Printf("Original size: %s\n", formatBytes(oldSize))
		fmt.Printf("New size: %s\n", formatBytes(newSize))

		if reduction > 0 {
			fmt.Printf("Change: %.0f%% smaller\n", math.Abs(reduction))
		} else {
			fmt.Printf("Change: %.0f%% LARGER\n", math.Abs(reduction))
		}
	}
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
