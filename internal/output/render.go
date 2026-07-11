package output

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"

	"go.yaml.in/yaml/v4"
)

// Render writes v in the requested format (json|yaml|table|csv). Nested structs render natively
// as json/yaml; table/csv flatten the top-level fields to key/value rows so a single record is
// still readable without a per-type column map.
func Render(w io.Writer, format string, v any) error {
	switch strings.ToLower(format) {
	case "json":
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(v)
	case "yaml":
		m, err := toMap(v)
		if err != nil {
			return err
		}
		b, err := yaml.Marshal(m)
		if err != nil {
			return err
		}
		_, err = w.Write(b)
		return err
	case "csv":
		m, err := toMap(v)
		if err != nil {
			return err
		}
		for _, k := range sortedKeys(m) {
			fmt.Fprintf(w, "%s,%v\n", k, scalar(m[k]))
		}
		return nil
	default: // table
		m, err := toMap(v)
		if err != nil {
			return err
		}
		tw := tabwriter.NewWriter(w, 0, 2, 2, ' ', 0)
		for _, k := range sortedKeys(m) {
			fmt.Fprintf(tw, "%s\t%v\n", k, scalar(m[k]))
		}
		return tw.Flush()
	}
}

// toMap marshals v to a top-level string-keyed map (via its JSON form). Non-objects (arrays,
// scalars) are wrapped under "value" so they still render.
func toMap(v any) (map[string]any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return map[string]any{"value": v}, nil
	}
	return m, nil
}

func sortedKeys(m map[string]any) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}

// scalar keeps primitives as-is and collapses nested objects/arrays to compact JSON for a cell.
func scalar(v any) any {
	switch v.(type) {
	case map[string]any, []any:
		b, _ := json.Marshal(v)
		return string(b)
	}
	return v
}
