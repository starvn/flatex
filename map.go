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

package flatex

import (
	"regexp"
	"strconv"
	"strings"
)

var defaultCollectionPattern = regexp.MustCompile(`\.\*\.`)

func newMap(t Tokenizer) (*Map, error) {
	sep := t.Separator()
	var hasWildcard *regexp.Regexp
	var err error
	if sep == "." {
		hasWildcard = defaultCollectionPattern
	} else {
		hasWildcard, err = regexp.Compile(sep + `\*` + sep)
	}
	if err != nil {
		return nil, err
	}
	return &Map{
		m:  make(map[string]interface{}),
		t:  t,
		re: hasWildcard,
	}, nil
}

type Map struct {
	m  map[string]interface{}
	t  Tokenizer
	re *regexp.Regexp
}

func (m *Map) Move(original, newKey string) {
	if v, ok := m.m[original]; ok {
		m.m[newKey] = v
		delete(m.m, original)
		return
	}

	if m.re.MatchString(original) {
		m.moveSliceAttribute(original, newKey)
		return
	}

	sep := m.t.Separator()

	for k := range m.m {
		if !strings.HasPrefix(k, original) {
			continue
		}

		if k[len(original):len(original)+1] != sep {
			continue
		}

		m.m[newKey+sep+k[len(original)+1:]] = m.m[k]
		delete(m.m, k)
	}
}

// Del deletes a key out of the map with the given prefix
func (m *Map) Del(prefix string) {
	if _, ok := m.m[prefix]; ok {
		delete(m.m, prefix)
		return
	}

	if m.re.MatchString(prefix) {
		m.delSliceAttribute(prefix)
		return
	}

	sep := m.t.Separator()

	for k := range m.m {
		if !strings.HasPrefix(k, prefix) {
			continue
		}

		if k[len(prefix):len(prefix)+1] != sep {
			continue
		}

		delete(m.m, k)
	}
}

func (m *Map) delSliceAttribute(prefix string) {
	i := strings.Index(prefix, "*")
	sep := m.t.Separator()
	prefixRemainder := prefix[i+1:]
	recursive := strings.Index(prefixRemainder, "*") > -1

	for k := range m.m {
		if len(k) < i+2 {
			continue
		}

		if !strings.HasPrefix(k, prefix[:i]) {
			continue
		}

		if recursive {
			newPref := k[:i+1+strings.Index(k[i+1:], sep)] + prefixRemainder
			m.Del(newPref)
			continue
		}

		keyRemainder := k[i+1+strings.Index(k[i+1:], sep):]
		if keyRemainder == prefixRemainder {
			delete(m.m, k)
			continue
		}

		if !strings.HasPrefix(keyRemainder, prefixRemainder+sep) {
			continue
		}

		delete(m.m, k)
	}
}

func (m *Map) moveSliceAttribute(original, newKey string) {
	i := strings.Index(original, "*")
	sep := m.t.Separator()
	originalRemainder := original[i+1:]
	recursive := strings.Index(originalRemainder, "*") > -1

	newKeyOffset := strings.Index(newKey, "*")
	newKeyRemainder := newKey[newKeyOffset+1:]
	newKeyPrefix := newKey[:newKeyOffset]

	for k := range m.m {
		if len(k) <= i+2 {
			continue
		}

		if !strings.HasPrefix(k, original[:i]) {
			continue
		}

		remainder := k[i:]
		idLen := strings.Index(remainder, sep)
		cleanRemainder := k[i+idLen:]
		keyPrefix := newKeyPrefix + k[i:i+idLen]

		if recursive {
			m.Move(k[:i+idLen]+originalRemainder, keyPrefix+newKeyRemainder)
			continue
		}

		if cleanRemainder == originalRemainder[1:] {
			m.m[keyPrefix+newKeyRemainder] = m.m[k]
			delete(m.m, k)
			continue
		}

		rPrefix := originalRemainder[1:] + sep

		if cleanRemainder != sep+originalRemainder[1:] && !strings.HasPrefix(cleanRemainder, sep+rPrefix) {
			continue
		}

		m.m[keyPrefix+newKeyRemainder+cleanRemainder[len(rPrefix):]] = m.m[k]
		delete(m.m, k)
	}
}

func (m *Map) Expand() map[string]interface{} {
	res := map[string]interface{}{}
	hasCollections := false
	for k, v := range m.m {
		ks := m.t.Keys(k)
		tr := res

		if ks[len(ks)-1] == "#" {
			hasCollections = true
		}
		for _, tk := range ks[:len(ks)-1] {
			trnew, ok := tr[tk]
			if !ok {
				trnew = make(map[string]interface{})
				tr[tk] = trnew
			}
			tr = trnew.(map[string]interface{})
		}
		tr[ks[len(ks)-1]] = v
	}

	if !hasCollections {
		return res
	}

	return m.expandNestedCollections(res).(map[string]interface{})
}

func (m *Map) expandNestedCollections(original map[string]interface{}) interface{} {
	for k, v := range original {
		if t, ok := v.(map[string]interface{}); ok {
			original[k] = m.expandNestedCollections(t)
		}
	}

	size, ok := original["#"]
	if !ok {
		return original
	}

	col := make([]interface{}, size.(int))
	for k := range col {
		col[k] = original[strconv.Itoa(k)]
	}
	return col
}
