package static

import (
	"fmt"
	"math/big"
	"sort"
)

// define a type for slice of *big.Int
type BigIntSlice []*big.Int

func (p BigIntSlice) Len() int           { return len(p) }
func (p BigIntSlice) Less(i, j int) bool { return p[i].Cmp(p[j]) < 0 }
func (p BigIntSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Funktion zum Sortieren von Hex-Strings nach ihrer numerischen Größe
func SortHexStrings(hexStrings []string) ([]string, error) {
	var nums BigIntSlice
	for _, hexStr := range hexStrings {
		num := new(big.Int)
		_, success := num.SetString(hexStr, 16)
		if !success {
			return nil, fmt.Errorf("invalid hex string: %s", hexStr)
		}
		nums = append(nums, num)
	}

	sort.Sort(nums)

	var sortedHexStrings []string
	for _, num := range nums {
		sortedHexStrings = append(sortedHexStrings, fmt.Sprintf("%x", num))
	}

	return sortedHexStrings, nil
}
