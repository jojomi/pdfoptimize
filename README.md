# pdfoptimize

Optimize PDF files

Usage:
  pdfoptimize [input] [output] [flags]

Flags:
      --dpi int        image DPI (0 = auto)
  -e, --ebook          optimize for ebook
  -h, --help           help for pdfoptimize
  -i, --inplace        optimize file in-place
      --prepress       optimize for prepress
  -p, --print          optimize for print
  -s, --screen         optimize for screen
  -q, --silent         no output
      --style string   optimization style (screen, print, ebook)


## Technical details

This tool uses locally installed (https://www.ghostscript.com)[Ghostscript] to optimize PDF files.

In many cases this results in a much smaller file size, especially for documents containing non-optimized images (scanned documents and so on).

This CLI tool uses the (https://github.com/jojomi/pdfopt)[pdfopt] library that can be useful in your own Go projects too.

## Install

If you have a local Go environment, you can get the latest binary like this:

``` shell
go install github.com/jojomi/pdfoptimize@latest
```

Otherwise see the (https://github.com/jojomi/pdfoptimize/releases)[Releases section on Github].
>>>>>>> 9852d30 (Initial code)
