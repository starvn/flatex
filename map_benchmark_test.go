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
	"testing"
)

var result interface{}

func BenchmarkIntegration(b *testing.B) {
	var res interface{}

	for _, size := range []int{1, 5, 50, 500} {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {

			in := getInputData(size)

			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				flatten, _ := Flatten(in, DefaultTokenizer)
				flatten.Move("a.*.b.*.c", "a.*.b.*.x")
				flatten.Move("b.b", "b.c")
				flatten.Del("b")
				res = flatten.Expand()
			}
		})
	}
	result = res
}

func BenchmarkExpand(b *testing.B) {
	var res interface{}

	for _, size := range []int{1, 5, 50, 500} {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {

			flatten, _ := Flatten(getInputData(size), DefaultTokenizer)
			flatten.Move("a.*.b.*.c", "a.*.b.*.x")
			flatten.Move("b.b", "b.c")
			flatten.Del("b")

			b.ResetTimer()

			for n := 0; n < b.N; n++ {
				res = flatten.Expand()
			}
		})
	}
	result = res
}

func BenchmarkNew(b *testing.B) {
	var res interface{}

	for _, size := range []int{1, 5, 50, 500} {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {

			in := getInputData(size)

			b.ResetTimer()

			for n := 0; n < b.N; n++ {
				res, _ = Flatten(in, DefaultTokenizer)
			}
		})
	}
	result = res
}

func BenchmarkMove(b *testing.B) {
	var res *Map

	for _, size := range []int{1, 5, 50, 500} {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {

			res, _ = Flatten(getInputData(size), DefaultTokenizer)

			b.ResetTimer()

			for n := 0; n < b.N; n++ {
				if n%2 == 0 {
					res.Move("a.*.b.*.c", "a.*.b.*.x")
					res.Move("b.b", "b.c")
				} else {
					res.Move("a.*.b.*.x", "a.*.b.*.c")
					res.Move("b.c", "b.b")
				}

			}
		})
	}
	result = res
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
