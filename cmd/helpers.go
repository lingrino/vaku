package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
)

var (
	errOutputFormat  = errors.New("unsupported output format")
	errOutputType    = errors.New("unsupported output type")
	errJSONMarshal   = errors.New("json marshal")
	errJSONUnmarshal = errors.New("json unmarshal")
)

// combineErr combines two errors to be output later.
func (c *cli) combineErr(e1, e2 error) error {
	if e1 == nil && e2 == nil {
		return nil
	}
	if e2 == nil {
		return e1
	}
	if e1 == nil {
		return e2
	}
	return fmt.Errorf("%s\n%s%s", e1, c.flagIndent, e2) //nolint:errorlint
}

// output handles outputting all of our messages (regular output or errors).
func (c *cli) output(out any) {
	outW := c.cmd.OutOrStdout()
	errW := c.cmd.ErrOrStderr()

	if out == nil {
		return
	}

	switch c.flagFormat {
	case "json":
		c.outputJSON(outW, out)
	case "text":
		c.outputText(outW, out)
	default:
		c.outputText(errW, errOutputFormat)
	}
}

// outputJSON handles output when flagFormat == json.
func (c *cli) outputJSON(w io.Writer, out any) {
	var jsonOut any

	switch out := out.(type) {
	case error:
		w = c.cmd.ErrOrStderr()
		jsonOut = map[string]string{
			"error": out.Error(),
		}
	default:
		jsonOut = out
	}

	json, err := json.MarshalIndent(jsonOut, "", c.flagIndent)
	if err != nil {
		c.outputText(c.cmd.ErrOrStderr(), errJSONMarshal)
		return
	}

	fmt.Fprintf(w, "%s\n", json)
}

// outputJSON handles output when flagFormat == text.
func (c *cli) outputText(w io.Writer, out any) {
	switch out := out.(type) {
	case error:
		w = c.cmd.ErrOrStderr()
		c.outputTextError(w, out)
	case string:
		c.outputTextString(w, out)
	case []string:
		c.outputTextList(w, out)
	case map[string]any:
		c.outputTextMap(w, 0, out)
	case map[string]map[string]any:
		c.outputTextNestedMap(w, out)
	default:
		c.outputTextError(c.cmd.ErrOrStderr(), errOutputType)
	}
}

// outputTextError outputs errors.
func (c *cli) outputTextError(w io.Writer, e error) {
	fmt.Fprintf(w, "ERROR: %s\n", e)
}

// outputTextString outputs strings.
func (c *cli) outputTextString(w io.Writer, s string) {
	fmt.Fprintf(w, "%s\n", s)
}

// outputTextList outputs lists of strings.
func (c *cli) outputTextList(w io.Writer, l []string) {
	if c.flagSort {
		sort.Strings(l)
	}

	for _, s := range l {
		c.outputTextString(w, s)
	}
}

// outputTextMap outputs maps of strings to interfaces.
func (c *cli) outputTextMap(w io.Writer, indentTimes int, m map[string]any) {
	indent := strings.Repeat(c.flagIndent, indentTimes)

	keys := c.mapKeys(m)
	for _, k := range keys {
		fmt.Fprintf(w, "%s%s: %+v\n", indent, k, m[k])
	}
}

// outputTextNestedMap outputs nested maps of maps of strings to interfaces.
func (c *cli) outputTextNestedMap(w io.Writer, m map[string]map[string]any) {
	keys := c.nestedMapKeys(m)
	for _, k := range keys {
		fmt.Fprintf(w, "%+v\n", k)
		c.outputTextMap(w, 1, m[k])
	}
}

// mapKeys gets a list of (optionally sorted) keys from a map.
func (c *cli) mapKeys(m map[string]any) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	if c.flagSort {
		sort.Strings(keys)
	}
	return keys
}

// nestedMapKeys gets a list of (optionally sorted) keys from a nested map.
func (c *cli) nestedMapKeys(m map[string]map[string]any) []string {
	keys := make([]string, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	if c.flagSort {
		sort.Strings(keys)
	}
	return keys
}
