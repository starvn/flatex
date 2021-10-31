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
	"fmt"
	"testing"
)

var _ interface{}

func BenchmarkIntegration(b *testing.B) {
	var res interface{}

	for _, size := range []int{1, 5, 50, 500} {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			in := getInputData(size)

			b.ResetTimer()

			for n := 0; n < b.N; n++ {
				flatten, _ := New(in)
				flatten.Move([]string{"a", "*", "b", "*", "c"}, []string{"a", "*", "b", "*", "x"})
				flatten.Move([]string{"b", "b"}, []string{"c"})
				flatten.Del([]string{"b"})
				res = flatten.Get([]string{})
			}
		})
	}
	_ = res
}

func BenchmarkExpand(b *testing.B) {
	var res interface{}

	for _, size := range []int{1, 5, 50, 500} {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			flatten, _ := New(getInputData(size))
			flatten.Move([]string{"a", "*", "b", "*", "c"}, []string{"a", "*", "b", "*", "x"})
			flatten.Move([]string{"b", "b"}, []string{"c"})
			flatten.Del([]string{"b"})

			b.ResetTimer()

			for n := 0; n < b.N; n++ {
				res = flatten.Get([]string{})
			}
		})
	}
	_ = res
}

func BenchmarkNew(b *testing.B) {
	var res interface{}

	for _, size := range []int{1, 5, 50, 500} {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			in := getInputData(size)

			b.ResetTimer()

			for n := 0; n < b.N; n++ {
				res, _ = New(in)
			}
		})
	}
	_ = res
}

func BenchmarkMove(b *testing.B) {
	var res *Tree

	for _, size := range []int{1, 5, 50, 500} {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {

			res, _ = New(getInputData(size))

			b.ResetTimer()

			for n := 0; n < b.N; n++ {
				if n%2 == 0 {
					res.Move([]string{"a", "*", "b", "*", "c"}, []string{"a", "*", "b", "*", "x"})
					res.Move([]string{"b", "b"}, []string{"b", "c"})
				} else {
					res.Move([]string{"a", "*", "b", "*", "x"}, []string{"a", "*", "b", "*", "c"})
					res.Move([]string{"b", "c"}, []string{"b", "b"})
				}
			}
		})
	}
	_ = res
}

func getInputData(size int) map[string]interface{} {
	first := map[string]interface{}{
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
	}
	second := map[string]interface{}{
		"b": []interface{}{
			map[string]interface{}{
				"c": map[string]interface{}{
					"a": 1,
				},
				"aa": 1,
			},
		},
		"aa": 1,
	}

	collection := make([]interface{}, 2*size)
	for i := 0; i < size; i++ {
		collection[2*i] = first
		collection[2*i+1] = second
	}

	return map[string]interface{}{
		"a": collection,
		"b": map[string]interface{}{
			"a":  42,
			"bb": true,
			"b":  map[string]interface{}{"a": 42},
		},
		"turbo": false,
	}
}
