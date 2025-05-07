package algorithm

import (
	"backend/search"
	"encoding/json"
	"os"
	"strings"
)

/*
   Struktur JSON
   -------------
   TargetElement : [                       // array‑of‑paths
     [                                       // satu path
       { "Product1": [[IngA, IngB]] },
       { "Product2": [[IngC, IngD]] },
       ...
     ],
     ...
   ]
*/

// StepNested = { "Mud": [["Water","Earth"]] }
type StepNested map[string][][]string

// PathsToNestedJSON menulis JSON dengan format di atas.
func PathsToNestedJSON(target string, paths [][]*search.ElementNode) (string, error) {
	out := make(map[string][]([]StepNested))
	for _, p := range paths {
		var onePath []StepNested
		for i := 0; i+2 < len(p); i += 3 {
			prod, a, b := p[i], p[i+1], p[i+2]
			onePath = append(onePath, StepNested{
				prod.Name: {{a.Name, b.Name}},
			})
		}
		out[target] = append(out[target], onePath)
	}

	fileName := "paths_" + strings.ReplaceAll(strings.ToLower(target), " ", "_") + ".json"
	f, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		return "", err
	}
	return fileName, nil
}