package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

const (
	errOutputFormat = "error: unsupported output format"
	errOutputType   = "error: unsupported output type"
	errJSONMarshal  = "error: json marshal"
)

func (c *cli) output(out interface{}) {
	switch c.flagFormat {
	case "json":
		c.outputJSON(out)
	case "text":
		c.outputText(out)
	default:
		errW := c.cmd.ErrOrStderr()
		fmt.Fprintf(errW, "%s\n", errOutputFormat)
	}
}

func (c *cli) outputJSON(out interface{}) {
	outW := c.cmd.OutOrStdout()
	errW := c.cmd.ErrOrStderr()

	json, err := json.MarshalIndent(out, "", c.flagIndent)
	if err != nil {
		fmt.Fprintf(errW, "%s\n", errJSONMarshal)
		return
	}

	fmt.Fprintf(outW, "%s\n", json)
}

func (c *cli) outputText(out interface{}) {
	switch out := out.(type) {
	case string:
		c.outputTextString(out)
	case []string:
		c.outputTextList(out)
	case map[string]interface{}:
		c.outputTextMap(0, out)
	case map[string]map[string]interface{}:
		c.outputTextNestedMap(out)
	default:
		errW := c.cmd.ErrOrStderr()
		fmt.Fprintf(errW, "%s\n", errOutputType)
	}
}

func (c *cli) outputTextString(s string) {
	outW := c.cmd.OutOrStdout()
	fmt.Fprintf(outW, "%s\n", s)
}

func (c *cli) outputTextList(l []string) {
	if c.flagSort {
		sort.Strings(l)
	}

	for _, s := range l {
		c.outputTextString(s)
	}
}

func (c *cli) outputTextMap(indentTimes int, m map[string]interface{}) {
	outW := c.cmd.OutOrStdout()
	indent := strings.Repeat(c.flagIndent, indentTimes)

	keys := c.mapKeys(m)
	for _, k := range keys {
		fmt.Fprintf(outW, "%s%s => %+v\n", indent, k, m[k])
	}
}

func (c *cli) outputTextNestedMap(m map[string]map[string]interface{}) {
	outW := c.cmd.OutOrStdout()

	keys := c.nestedMapKeys(m)
	for _, k := range keys {
		fmt.Fprintf(outW, "%+v\n", k)
		c.outputTextMap(1, m[k])
	}
}

func (c *cli) mapKeys(m map[string]interface{}) []string {
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

func (c *cli) nestedMapKeys(m map[string]map[string]interface{}) []string {
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
