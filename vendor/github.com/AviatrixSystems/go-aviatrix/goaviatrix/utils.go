package goaviatrix

import (
	"fmt"
)

var ErrNotFound = fmt.Errorf("ErrNotFound")

func ExpandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, v.(string))
		}
	}
	return vs
}

// difference returns the elements in a that aren't in b
func Difference(a, b []string) []string {
	mb := map[string]bool{}
	for _, x := range b {
		mb[x] = true
	}
	ab := make([]string, 0)
	for _, x := range a {
		if _, ok := mb[x]; !ok {
			ab = append(ab, x)
		}
	}
	return ab
}

// DifferenceSlice returns the one-dimension elements in two-dimension slice a that aren't in two-dimension b
func DifferenceSlice(a, b [][]string) [][]string {
	if a == nil || len(a) == 0 || b == nil || len(b) == 0 {
		return a
	}

	aa := make([]string, 0)
	for i := range a {
		temp := ""
		for j := range a[i] {
			temp += a[i][j]
		}
		aa = append(aa, temp)
	}

	bb := make([]string, 0)
	for t := range b {
		temp := ""
		for m := range b[t] {
			temp += b[t][m]
		}
		bb = append(bb, temp)
	}

	mb := map[string]bool{}
	for x := range bb {
		mb[bb[x]] = true
	}
	ab := make([][]string, 0)
	for x := range aa {
		if _, ok := mb[aa[x]]; !ok {
			ab = append(ab, a[x])
		}
	}
	return ab
}
