package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"

	"go.yaml.in/yaml/v4"
)

// Render writes v in the requested format (json|yaml|table|csv). A single record flattens its
// top-level fields to key/value rows; a JSON array of objects renders as a real table (one row
// per record, union of keys as columns) so lists and `history` export cleanly to CSV.
func Render(w io.Writer, format string, v any) error {
	switch strings.ToLower(format) {
	case "json":
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(v)
	case "yaml":
		n, err := normalize(v) // render arrays and objects natively
		if err != nil {
			return err
		}
		b, err := yaml.Marshal(n)
		if err != nil {
			return err
		}
		_, err = w.Write(b)
		return err
	case "csv":
		return renderCSV(w, v)
	default:
		return renderTable(w, v)
	}
}

func renderCSV(w io.Writer, v any) error {
	cw := csv.NewWriter(w)
	if headers, rows, ok := tabular(v); ok { // array of objects → header + a row per record
		_ = cw.Write(headers)
		for _, r := range rows {
			_ = cw.Write(stringsOf(r))
		}
	} else { // single record → key,value rows
		m, err := toMap(v)
		if err != nil {
			return err
		}
		for _, k := range sortedKeys(m) {
			_ = cw.Write([]string{k, fmt.Sprintf("%v", scalar(m[k]))})
		}
	}
	cw.Flush()
	return cw.Error()
}

func renderTable(w io.Writer, v any) error {
	tw := tabwriter.NewWriter(w, 0, 2, 2, ' ', 0)
	if headers, rows, ok := tabular(v); ok {
		fmt.Fprintln(tw, strings.Join(headers, "\t"))
		for _, r := range rows {
			fmt.Fprintln(tw, strings.Join(stringsOf(r), "\t"))
		}
		return tw.Flush()
	}
	m, err := toMap(v)
	if err != nil {
		return err
	}
	for _, k := range sortedKeys(m) {
		fmt.Fprintf(tw, "%s\t%v\n", k, scalar(m[k]))
	}
	return tw.Flush()
}

// tabular reports whether v is a non-empty JSON array of objects and, if so, returns the sorted
// union of keys as headers plus one cell-row per element (nested values collapsed to one line).
func tabular(v any) (headers []string, rows [][]any, ok bool) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, nil, false
	}
	var arr []map[string]any
	if json.Unmarshal(b, &arr) != nil || len(arr) == 0 {
		return nil, nil, false
	}
	seen := map[string]bool{}
	for _, m := range arr {
		for k := range m {
			seen[k] = true
		}
	}
	for k := range seen {
		headers = append(headers, k)
	}
	sort.Strings(headers)
	for _, m := range arr {
		row := make([]any, len(headers))
		for i, h := range headers {
			row[i] = scalar(m[h])
		}
		rows = append(rows, row)
	}
	return headers, rows, true
}

func stringsOf(row []any) []string {
	out := make([]string, len(row))
	for i, c := range row {
		out[i] = fmt.Sprintf("%v", c)
	}
	return out
}

// normalize returns v as generic maps/slices via a JSON round-trip, so YAML renders nested
// structs and arrays natively.
func normalize(v any) (any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var n any
	if err := json.Unmarshal(b, &n); err != nil {
		return nil, err
	}
	return n, nil
}

// toMap marshals v to a top-level string-keyed map (via its JSON form). Non-objects (arrays,
// scalars) are wrapped under "value" so they still render.
func toMap(v any) (map[string]any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if json.Unmarshal(b, &m) == nil {
		return m, nil
	}
	// Not a JSON object (array/scalar/null) — wrap so it still renders.
	return map[string]any{"value": v}, nil
}

func sortedKeys(m map[string]any) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}

// scalar keeps primitives as-is and collapses nested objects/arrays to compact one-line JSON so a
// single cell never breaks the table layout.
func scalar(v any) any {
	switch t := v.(type) {
	case map[string]any, []any:
		b, _ := json.Marshal(v)
		return cellOneLine(string(b))
	case string:
		return cellOneLine(t)
	}
	return v
}
