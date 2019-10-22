package cmd

// noTrimPathPrefix determines if the output paths should
// have the input path trimmed from the prefix
var noTrimPathPrefix bool

// indentString is used in folder_map to determine which string to
// use as an indent for the sub-sections
var indentString string

// format is the format of our string output. Defaults to json
var format string

var sourceAddress string
var sourceToken string
var sourceNamespace string

var targetAddress string
var targetToken string
var targetNamespace string

var useSourceTarget bool
