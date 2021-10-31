/*
 * Copyright (c) 2021 Huy Duc Dao
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package tree

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

const tabSize = 4

var sb = new(strings.Builder)

func (n *node) writeTo(bd *strings.Builder) {
	for i, e := range n.edges {
		e.writeTo(bd, []bool{i == len(n.edges)-1})
	}
}

func (e *edge) writeTo(bd *strings.Builder, tabList []bool) {
	length := len(tabList)
	isLast, tlist := tabList[length-1], tabList[:length-1]
	for _, hasTab := range tlist {
		if hasTab {
			bd.Write(bytes.Repeat([]byte(" "), tabSize))
			continue
		}
		bd.WriteRune('│')
		bd.Write(bytes.Repeat([]byte(" "), tabSize-1))
	}
	if !isLast {
		bd.WriteRune('├')
	} else {
		bd.WriteRune('└')
	}
	bd.WriteString("── ")
	bd.WriteString(e.label)
	if e.n.IsLeaf() {
		fmt.Fprintf(bd, "\t%+v", e.n.Value)
	} else if e.n.isCollection {
		bd.WriteString(" []")
	}
	bd.WriteByte('\n')
	for i, next := range e.n.edges {
		if len(tabList) < next.n.depth {
			tabList = append(tabList, i == len(e.n.edges)-1)
		} else {
			tabList[next.n.depth-1] = i == len(e.n.edges)-1
		}
		next.writeTo(bd, tabList)
	}
}

func (t *Tree) String() string {
	sb.Reset()
	sb.WriteByte('\n')
	t.root.writeTo(sb)
	return sb.String()
}

func TestTree_Del(t *testing.T) {
	for _, tc := range []struct {
		name    string
		pattern string
		in      map[string]interface{}
		out     string
	}{
		{
			name:    "unknown",
			pattern: "abc",
			in: map[string]interface{}{
				"sonic": 42,
				"turbo": false,
			},
			out: `
├── sonic	42
└── turbo	false
`,
		},
		{
			name:    "empty slice",
			pattern: "data.*.password",
			in: map[string]interface{}{
				"data": []interface{}{},
			},
			// but not
			// ── data	<nil>
			out: `
└── data	[]
`,
		},
		{
			name:    "empty map",
			pattern: "data.*.password",
			in: map[string]interface{}{
				"data": map[string]interface{}{},
			},
			// but not
			// ── data    <nil>
			out: `
└── data	map[]
`,
		},
		{
			name:    "plain",
			pattern: "sonic",
			in: map[string]interface{}{
				"sonic": 42,
				"turbo": false,
			},
			out: `
└── turbo	false
`,
		},
		{
			name:    "element_in_struct",
			pattern: "internal.sonic",
			in: map[string]interface{}{
				"internal": map[string]interface{}{
					"sonic": 42,
					"turbo": false,
				},
				"turbo": false,
			},
			out: `
├── internal
│   └── turbo	false
└── turbo	false
`,
		},
		{
			name:    "element_in_struct_with_wildcard",
			pattern: "a.*.sonic",
			in: map[string]interface{}{
				"a": map[string]interface{}{
					"first": map[string]interface{}{
						"sonic": 42,
						"turbo": false,
					},
					"last": map[string]interface{}{
						"sonic": 42,
						"turbo": false,
					},
				},
				"turbo": false,
			},
			out: `
├── a
│   ├── first
│   │   └── turbo	false
│   └── last
│       └── turbo	false
└── turbo	false
`,
		},
		{
			name:    "struct",
			pattern: "internal",
			in: map[string]interface{}{
				"internal": map[string]interface{}{
					"sonic": 42,
					"turbo": false,
				},
				"turbo": false,
			},
			out: `
└── turbo	false
`,
		},
		{
			name:    "element_in_substruct",
			pattern: "internal.internal.sonic",
			in: map[string]interface{}{
				"internal": map[string]interface{}{
					"sonic": 42,
					"turbo": false,
					"internal": map[string]interface{}{
						"sonic": 42,
						"turbo": false,
					},
				},
				"turbo": false,
			},
			out: `
├── internal
│   ├── internal
│   │   └── turbo	false
│   ├── sonic	42
│   └── turbo	false
└── turbo	false
`,
		},
		{
			name:    "similar_names",
			pattern: "a.a.a",
			in: map[string]interface{}{
				"a": map[string]interface{}{
					"a": map[string]interface{}{
						"a": map[string]interface{}{
							"a": 1,
						},
						"aa": 1,
					},
					"aa": 1,
				},
				"turbo": false,
			},
			out: `
├── a
│   ├── a
│   │   └── aa	1
│   └── aa	1
└── turbo	false
`,
		},
		{
			name:    "collection_element_attributes",
			pattern: "a.*.a",
			in: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{
						"a": map[string]interface{}{
							"a": map[string]interface{}{
								"a": 1,
							},
							"aa": 1,
						},
						"aa": 1,
					},
					map[string]interface{}{
						"a":  42,
						"aa": 1,
					},
				},
				"turbo": false,
			},
			out: `
├── a []
│   ├── 0
│   │   └── aa	1
│   └── 1
│       └── aa	1
└── turbo	false
`,
		},
		{
			name:    "nested_collection_element_attributes",
			pattern: "a.*.b.*.c",
			in: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{
						"b": []interface{}{
							map[string]interface{}{
								"c": map[string]interface{}{
									"a": 1,
								},
								"aa": 1,
							},
							map[string]interface{}{
								"c": map[string]interface{}{
									"a": 2,
								},
								"aa": 1,
							},
						},
						"aa": 1,
					},
					map[string]interface{}{
						"b": []interface{}{
							map[string]interface{}{
								"c": map[string]interface{}{
									"a": 1,
								},
								"aa": 1,
							},
						},
						"aa": 1,
					},
				},
				"turbo": false,
			},
			out: `
├── a []
│   ├── 0
│   │   ├── aa	1
│   │   └── b []
│   │       ├── 0
│   │       │   └── aa	1
│   │       └── 1
│   │           └── aa	1
│   └── 1
│       ├── aa	1
│       └── b []
│           └── 0
│               └── aa	1
└── turbo	false
`,
		},
		{
			name:    "large_collection_element_attributes",
			pattern: "a.*.a",
			in: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{
						"a":  1,
						"aa": 1,
					},
					map[string]interface{}{
						"a":  2,
						"aa": 1,
					},
					map[string]interface{}{
						"a":  1,
						"aa": 1,
					},
					map[string]interface{}{
						"a":  2,
						"aa": 1,
					},
					map[string]interface{}{
						"a":  1,
						"aa": 1,
					},
					map[string]interface{}{
						"a":  2,
						"aa": 1,
					},
					map[string]interface{}{
						"a":  1,
						"aa": 1,
					},
					map[string]interface{}{
						"a":  2,
						"aa": 1,
					},
					map[string]interface{}{
						"a":  1,
						"aa": 1,
					},
					map[string]interface{}{
						"a":  2,
						"aa": 1,
					},
					map[string]interface{}{
						"a":  1,
						"aa": 1,
					},
					map[string]interface{}{
						"a":  2,
						"aa": 1,
					},
				},
				"turbo": false,
			},
			out: `
├── a []
│   ├── 0
│   │   └── aa	1
│   ├── 1
│   │   └── aa	1
│   ├── 2
│   │   └── aa	1
│   ├── 3
│   │   └── aa	1
│   ├── 4
│   │   └── aa	1
│   ├── 5
│   │   └── aa	1
│   ├── 6
│   │   └── aa	1
│   ├── 7
│   │   └── aa	1
│   ├── 8
│   │   └── aa	1
│   ├── 9
│   │   └── aa	1
│   ├── 10
│   │   └── aa	1
│   └── 11
│       └── aa	1
└── turbo	false
`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			res, _ := New(tc.in)
			res.Sort()

			res.Del(strings.Split(tc.pattern, "."))
			tree := res.String()
			if tree != tc.out {
				t.Errorf("unexpected result (%s):'%s'\n'%s'", tc.pattern, tree, tc.out)
			}
		})
	}
}

func TestTree_Get(t *testing.T) {
	in := map[string]interface{}{
		"a": []interface{}{
			map[string]interface{}{
				"a":  0,
				"aa": 1,
			},
			map[string]interface{}{
				"a":  1,
				"aa": 1,
			},
			map[string]interface{}{
				"a":  2,
				"aa": 1,
			},
			map[string]interface{}{
				"a":  3,
				"aa": 1,
			},
			map[string]interface{}{
				"a":  4,
				"aa": 1,
			},
			map[string]interface{}{
				"a":  5,
				"aa": 1,
			},
			map[string]interface{}{
				"a":  6,
				"aa": 1,
			},
			map[string]interface{}{
				"a":  7,
				"aa": 1,
			},
			map[string]interface{}{
				"a":  8,
				"aa": 1,
			},
			map[string]interface{}{
				"a":  9,
				"aa": 1,
			},
			map[string]interface{}{
				"a":  10,
				"aa": 1,
			},
			map[string]interface{}{
				"a":  11,
				"aa": 1,
			},
		},
		"b": map[string]interface{}{
			"c": []interface{}{1, 2, 3, 4},
		},
		"turbo": false,
	}
	tree, _ := New(in)
	v := tree.Get([]string{"a", "1", "a"})
	i, ok := v.(int)
	if !ok {
		t.Errorf("unexpected result type: %v", v)
		return
	}
	if i != 1 {
		t.Errorf("unexpected result %d", i)
	}

	l := tree.Get([]string{"a", "*", "a"}).([]interface{})
	if len(l) != 12 {
		t.Errorf("unexpected number of returned values")
		return
	}
	for i, v := range l {
		if i != v.(int) {
			t.Errorf("unexpected value #%d: %v", i, v)
		}
	}
}

func TestTree_Move(t *testing.T) {
	for _, tc := range []struct {
		name string
		src  string
		dst  string
		in   map[string]interface{}
		out  string
	}{
		{
			name: "plain",
			src:  "a",
			dst:  "b",
			in:   map[string]interface{}{"a": 42},
			out: `
└── b	42
`,
		},
		{
			name: "in_struct",
			src:  "b.a",
			dst:  "b.c",
			in: map[string]interface{}{
				"a": 1,
				"b": map[string]interface{}{"a": 42},
			},
			out: `
├── a	1
└── b
    └── c	42
`,
		},
		{
			name: "in_struct_depth",
			src:  "b.a",
			dst:  "b.b.a.b.a.b.c",
			in: map[string]interface{}{
				"a": 1,
				"b": map[string]interface{}{"a": 42},
			},
			out: `
├── a	1
└── b
    └── b
        └── a
            └── b
                └── a
                    └── b
                        └── c	42
`,
		},
		{
			name: "from_struct",
			src:  "b.a",
			dst:  "c",
			in: map[string]interface{}{
				"a": 1,
				"b": map[string]interface{}{"a": 42},
			},
			out: `
├── a	1
├── b	<nil>
└── c	42
`,
		},
		{
			name: "from_struct_with_wildcard",
			src:  "b.*.c",
			dst:  "b.*.x",
			in: map[string]interface{}{
				"c": 42,
				"b": map[string]interface{}{
					"first": map[string]interface{}{"c": map[string]interface{}{"d": 42}},
					"last":  map[string]interface{}{"m": 42, "c": map[string]interface{}{"d": 42}},
				},
			},
			out: `
├── b
│   ├── first
│   │   └── x
│   │       └── d	42
│   └── last
│       ├── m	42
│       └── x
│           └── d	42
└── c	42
`,
		},
		{
			name: "from_struct_with_wildcard_deep",
			src:  "b.*.c",
			dst:  "b.*.c.b.x",
			in: map[string]interface{}{
				"c": 42,
				"b": map[string]interface{}{
					"first": map[string]interface{}{"c": map[string]interface{}{"d": 42}},
					"last":  map[string]interface{}{"m": 42, "c": map[string]interface{}{"d": 42}},
				},
			},
			out: `
├── b
│   ├── first
│   │   └── c
│   │       └── b
│   │           └── x
│   │               └── d	42
│   └── last
│       ├── c
│       │   └── b
│       │       └── x
│       │           └── d	42
│       └── m	42
└── c	42
`,
		},
		{
			name: "from_collection",
			src:  "b.*.c",
			dst:  "b.*.x",
			in: map[string]interface{}{
				"a": 42,
				"b": []interface{}{
					map[string]interface{}{"c": 42},
					map[string]interface{}{"c": map[string]interface{}{"d": 42}},
				},
			},
			out: `
├── a	42
└── b []
    ├── 0
    │   └── x	42
    └── 1
        └── x
            └── d	42
`,
		},
		{
			name: "from_struct_nested",
			src:  "b.b",
			dst:  "c",
			in: map[string]interface{}{
				"a": 42,
				"b": map[string]interface{}{
					"a":  42,
					"bb": true,
					"b":  map[string]interface{}{"a": 42},
				},
			},
			out: `
├── a	42
├── b
│   ├── a	42
│   └── bb	true
└── c
    └── a	42
`,
		},
		{
			name: "collection",
			src:  "a.*.b",
			dst:  "a.*.c",
			in: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{
						"b": []interface{}{
							map[string]interface{}{
								"c": map[string]interface{}{
									"a": 1,
								},
								"aa": 1,
							},
							map[string]interface{}{
								"c": map[string]interface{}{
									"a": 2,
								},
								"aa": 1,
							},
						},
						"aa": 1,
					},
					map[string]interface{}{
						"b": []interface{}{
							map[string]interface{}{
								"c": map[string]interface{}{
									"a": 1,
								},
								"aa": 1,
							},
						},
						"aa": 1,
					},
				},
				"turbo": false,
			},
			out: `
├── a []
│   ├── 0
│   │   ├── aa	1
│   │   └── c []
│   │       ├── 0
│   │       │   ├── aa	1
│   │       │   └── c
│   │       │       └── a	1
│   │       └── 1
│   │           ├── aa	1
│   │           └── c
│   │               └── a	2
│   └── 1
│       ├── aa	1
│       └── c []
│           └── 0
│               ├── aa	1
│               └── c
│                   └── a	1
└── turbo	false
`,
		},
		{
			name: "recursive_collection",
			src:  "a.*.b.*.c",
			dst:  "a.*.b.*.x",
			in: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{
						"b": []interface{}{
							map[string]interface{}{
								"c": map[string]interface{}{
									"a": 1,
								},
								"aa": 1,
							},
							map[string]interface{}{
								"c": map[string]interface{}{
									"a": 2,
								},
								"aa": 1,
							},
						},
						"aa": 1,
					},
					map[string]interface{}{
						"b": []interface{}{
							map[string]interface{}{
								"c": map[string]interface{}{
									"a": 1,
								},
								"aa": 1,
							},
						},
						"aa": 1,
					},
				},
				"turbo": false,
			},
			out: `
├── a []
│   ├── 0
│   │   ├── aa	1
│   │   └── b []
│   │       ├── 0
│   │       │   ├── aa	1
│   │       │   └── x
│   │       │       └── a	1
│   │       └── 1
│   │           ├── aa	1
│   │           └── x
│   │               └── a	2
│   └── 1
│       ├── aa	1
│       └── b []
│           └── 0
│               ├── aa	1
│               └── x
│                   └── a	1
└── turbo	false
`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			res, _ := New(tc.in)
			res.Sort()
			original := res.String()

			res.Move(strings.Split(tc.src, "."), strings.Split(tc.dst, "."))

			res.Sort()

			if tree := res.String(); tree != tc.out {
				t.Errorf("unexpected result (%s -> %s) from:%s\nhave:%s\nwant:%s", tc.src, tc.dst, original, tree, tc.out)
			}
		})
	}
}

func TestTree_Append(t *testing.T) {
	for _, tc := range []struct {
		name string
		src  string
		dst  string
		in   map[string]interface{}
		out  string
	}{
		{
			name: "plain",
			src:  "a",
			dst:  "b",
			in:   map[string]interface{}{"a": []interface{}{42}, "b": []interface{}{1}},
			out: `
└── b []
    ├── 0	1
    └── 1	42
`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			res, _ := New(tc.in)
			res.Sort()
			original := res.String()

			res.Append(strings.Split(tc.src, "."), strings.Split(tc.dst, "."))

			res.Sort()

			if tree := res.String(); tree != tc.out {
				t.Errorf("unexpected result (%s -> %s) from:%s\nhave:%s\nwant:%s", tc.src, tc.dst, original, tree, tc.out)
			}
		})
	}
}
