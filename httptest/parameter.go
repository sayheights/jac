package httptest

import (
	"sort"
	"strconv"
	"strings"
)

type parameter struct {
	key     string
	val     map[string]bool
	in      string
	matched bool
}

func (p parameter) diff(params []parameter) *diff {
	var (
		match parameter
		found bool
	)
	for i, param := range params {
		if p.key == param.key {
			newP := param
			newP.matched = true
			params[i] = newP
			match, found = param, true
			break
		}
	}

	if !found {
		return &diff{typ: missingKey, key: p.key, in: p.in, val: p.valString()}
	}

	return p.diffMatch(match)
}

func (p parameter) diffMatch(match parameter) *diff {
	typ := resolveDiffType(p.val, match.val)

	if val, missing := diffSet(p.val, match.val); missing {
		return &diff{typ: typ, val: val, key: p.key, in: p.in}
	}

	if val, missing := diffSet(match.val, p.val); missing {
		return &diff{typ: unknownValue, val: val, key: p.key, in: p.in}
	}

	return nil
}

func (p parameter) valString() string {
	var sb []string
	for k := range p.val {
		sb = append(sb, k)
	}
	sort.Strings(sb)
	val := strings.Join(sb, ",")

	return val
}

// byKey implement sort.Interface for []parameter based on the key field.
type byKey []parameter

func (b byKey) Len() int { return len(b) }

func (b byKey) Less(i int, j int) bool { return b[i].key < b[j].key }

func (b byKey) Swap(i int, j int) { b[i], b[j] = b[j], b[i] }

func fromMap(m map[string][]string, in string) []parameter {
	var params byKey
	for k, v := range m {
		param := parameter{key: k, val: toSet(v), in: in}
		params = append(params, param)
	}
	sort.Sort(params)

	return params
}

func fromPath(p string) []parameter {
	var params []parameter
	p = strings.Trim(p, "/")
	split := strings.Split(p, "/")
	for i, seg := range split {
		key := "SEGMENT_" + strconv.Itoa(i)
		segs := toSet([]string{seg})
		param := parameter{key: key, val: segs, in: "path"}
		params = append(params, param)

	}
	return params
}

func toSet(s []string) map[string]bool {
	set := map[string]bool{}
	for _, i := range s {
		set[i] = true
	}

	return set
}
