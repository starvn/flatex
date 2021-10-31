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
	"fmt"
	"strings"
)

type Tokenizer interface {
	Separator() string
	Token([]string) string
	Keys(string) []string
}

type StringTokenizer string

func (s StringTokenizer) Token(ks []string) string { return strings.Join(ks, string(s)) }

func (s StringTokenizer) Keys(ks string) []string { return strings.Split(ks, string(s)) }

func (s StringTokenizer) Separator() string { return string(s) }

var DefaultTokenizer = StringTokenizer(".")

func Flatten(m map[string]interface{}, tokenizer Tokenizer) (*Map, error) {
	result, err := newMap(tokenizer)
	if err != nil {
		return nil, err
	}
	flatten(m, []string{}, func(ks []string, v interface{}) {
		result.m[tokenizer.Token(ks)] = v
	})
	return result, nil
}

type updateFunc func([]string, interface{})

func flatten(i interface{}, ks []string, update updateFunc) {
	switch v := i.(type) {
	case map[string]interface{}:
		flattenMap(v, ks, update)
	case []interface{}:
		flattenSlice(v, ks, update)
	default:
		update(ks, v)
	}
}

func flattenMap(m map[string]interface{}, ks []string, update updateFunc) {
	for k, v := range m {
		flatten(v, append(ks, k), update)
	}
}

func flattenSlice(vs []interface{}, ks []string, update updateFunc) {
	update(append(ks, "#"), len(vs))
	for i, v := range vs {
		flatten(v, append(ks, fmt.Sprintf("%d", i)), update)
	}
}
