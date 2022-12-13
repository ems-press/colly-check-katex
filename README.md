# colly-check-katex

Linkcheker that crawls ems.press and checks for katex rendering errors. Usage:

```
# go run main.go -h
Usage: check-katex [-u|--start-url ...] [-h|--help]

  -u, --start-url Starting point of the crawler (default: https://ems.press/journals)
  -h, --help      prints help information

Will only check URLs deeper than the given start URL. All errors are printed
stderr, verbose request information is printed to stdout.

Examples:

   go run main.go > /dev/null # only print errors
   go run main.go > /dev/null  2>&1 | tee errors.log # save all errors to a file
```
