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
	"errors"
)

const wildcard = "*"

var errNoNilValuesAllowed = errors.New("no nil values allowed")

type Tree struct {
	root *node
}

func New(v interface{}) (*Tree, error) {
	if v == nil {
		return nil, errNoNilValuesAllowed
	}

	tr := &Tree{
		root: &node{},
	}

	tr.Add([]string{}, v)

	return tr, nil
}

func (t *Tree) Add(ks []string, v interface{}) {
	if v == nil {
		return
	}
	t.root.Add(ks, v)
}

func (t *Tree) Del(ks []string) {
	t.root.Del(ks...)
}

func (t *Tree) Append(src, dst []string) {
	elements1, ok := t.root.Get(src...).([]interface{})
	if !ok {
		return
	}
	elements2, ok := t.root.Get(dst...).([]interface{})
	if !ok {
		return
	}

	t.root.Add(dst, append(elements2, elements1...))

	t.root.Del(src...)
}

func (t *Tree) Get(ks []string) interface{} {
	return t.root.Get(ks...)
}

func (t *Tree) Move(src, dst []string) {
	next := []nodeAndPath{{n: t.root, p: []string{}}}

	prefixLen := len(src)

	if prefixLen > 1 {
		next = t.collectMoveCandidates(src[:prefixLen-1], next)
	}

	var edgesToMove []edgeToMove
	lenDst := len(dst)

	isEdgeRelabel := prefixLen == lenDst

	for _, nap := range next {
		for i, e := range nap.n.edges {
			if e.label != src[prefixLen-1] {
				continue
			}

			if isEdgeRelabel {
				e.label = dst[prefixLen-1]
				break
			}

			edgesToMove = append(edgesToMove, edgeToMove{nodeAndPath: nap, e: e})

			copy(nap.n.edges[i:], nap.n.edges[i+1:])
			nap.n.edges[len(nap.n.edges)-1] = nil
			nap.n.edges = nap.n.edges[:len(nap.n.edges)-1]

			break
		}
	}

	if isEdgeRelabel {
		return
	}

	if prefixLen > lenDst {
		t.promoteEdges(edgesToMove, dst)
		return
	}

	t.embeddingEdges(edgesToMove, dst[prefixLen-1:])
}

func (t *Tree) Sort() {
	t.root.sort()
}

func (t *Tree) collectMoveCandidates(src []string, next []nodeAndPath) []nodeAndPath {
	var acc []nodeAndPath
	for _, step := range src {
		if step == wildcard {
			for _, nap := range next {
				for _, e := range nap.n.edges {
					acc = append(acc, nodeAndPath{n: e.n, p: append(nap.p, e.label)})
				}
			}
		} else {
			for _, nap := range next {
				for _, e := range nap.n.edges {
					if step == e.label {
						acc = append(acc, nodeAndPath{n: e.n, p: append(nap.p, e.label)})
						break
					}
				}
			}
		}
		next, acc = acc, next[:0]
	}
	return next
}

func (t *Tree) promoteEdges(edgesToMove []edgeToMove, dst []string) {
	var l string
	lenDst := len(dst)
	for _, n := range edgesToMove {
		parent := t.root
		for i, path := range dst[:lenDst-1] {
			if path == wildcard {
				l = n.p[i]
			} else {
				l = path
			}

			found := false
			for _, e := range parent.edges {
				if e.label != l {
					continue
				}
				found = true
				parent = e.n
				break
			}

			if !found {
				break
			}
		}

		n.e.n.SetDepth(parent.depth + 1)
		n.e.label = dst[lenDst-1]
		parent.edges = append(parent.edges, n.e)
	}
}

func (t *Tree) embeddingEdges(edgesToMove []edgeToMove, dst []string) {
	lenDst := len(dst)
	for _, em := range edgesToMove {
		root := em.n
		for _, k := range dst[:lenDst-1] {
			found := false
			for _, e := range root.edges {
				if e.label != k {
					continue
				}
				found = true
				root = e.n
				break
			}
			if found {
				continue
			}
			child := newNode(root.depth + 1)
			root.edges = append(root.edges, &edge{label: k, n: child})
			root = child
		}
		em.e.label = dst[lenDst-1]
		em.e.n.SetDepth(root.depth + 1)
		root.edges = append(root.edges, em.e)
	}
}

type nodeAndPath struct {
	n *node
	p []string
}

type edgeToMove struct {
	nodeAndPath
	e *edge
}
