package algorithm

import (
	"backend/search"
	"fmt"
)

/* ---------- helpers ---------- */

// ---------- helpers (unchanged) ----------
var primordial = map[string]struct{}{
	"Air": {}, "Water": {}, "Fire": {}, "Earth": {},
}
func isBaseElement(n *search.ElementNode) bool {
	if n == nil {
		return false
	}
	if _, ok := primordial[n.Name]; ok {
		return true
	}
	return len(n.Recipes) == 0
}

// ---------- pretty‑printer (tiny tweak: fixed‑step loop) ----------
func PrintCraftingPath(path []*search.ElementNode) {
	if len(path) == 0 {
		fmt.Println("No crafting path found.")
		return
	}
	fmt.Printf("Crafting path for %s:\n", path[0].Name)

	for i := 0; i < len(path); {
		// sentinel (ID==0) marks a plain base line
		if path[i].ID == 0 { break }

		// last item in the slice can be "<product> (base element)"
		if i+1 < len(path) && path[i+1].ID == 0 {
			fmt.Printf("%s <= (base element)\n", path[i].Name)
			i += 2
			continue
		}

		// normal triple
		product, ing1, ing2 := path[i], path[i+1], path[i+2]
		sfx := func(n *search.ElementNode) string {
			if isBaseElement(n) { return " (base)" }
			return ""
		}
		fmt.Printf("%s <= %s%s + %s%s\n",
			product.Name, ing1.Name, sfx(ing1), ing2.Name, sfx(ing2))
		i += 3
	}
}

// ---------- DFS with correct path building & branch isolation ----------
func DFS(target *search.ElementNode, maxPaths int) [][]*search.ElementNode {
	visited := make(map[string]bool)
	var results [][]*search.ElementNode
	sentinel := &search.ElementNode{ID: 0, Name: "BASE"}

	var explore func(el *search.ElementNode, cur []*search.ElementNode)
	explore = func(el *search.ElementNode, cur []*search.ElementNode) {
		if len(results) >= maxPaths {
			return
		}

		// leaf ----------------------------------------------------------------
		if isBaseElement(el) {
			cp := append([]*search.ElementNode(nil), cur...)
			cp = append(cp, el, sentinel) // show "<leaf> <= (base element)"
			results = append(results, cp)
			return
		}

		// internal node -------------------------------------------------------
		for _, r := range el.Recipes {
			if len(r) != 2 { continue }
			ing1, ing2 := r[0], r[1]

			key := fmt.Sprintf("%d:%d:%d", el.ID, ing1.ID, ing2.ID)
			if visited[key] { continue }
			visited[key] = true

			// **add the whole triple**, not only the two ingredients
			base := append(append([]*search.ElementNode(nil), cur...),
				el, ing1, ing2)

			b1, b2 := isBaseElement(ing1), isBaseElement(ing2)

			switch {
			case b1 && b2:
				results = append(results, base)

			case b1 && !b2:
				explore(ing2, base)

			case !b1 && b2:
				explore(ing1, base)

			default: // both expandable – explore on two *independent* copies
				left  := append([]*search.ElementNode(nil), base...)
				right := append([]*search.ElementNode(nil), base...)
				explore(ing1, left)
				explore(ing2, right)
			}
		}
	}

	explore(target, nil)
	return results
}

// ---------- wrapper (unchanged except for name) ----------
func FindCraftingPaths(el *search.ElementNode, max int) [][]*search.ElementNode {
	paths := DFS(el, max)
	return paths
}