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
	"encoding/json"
	"fmt"
)

func ExampleFlatten() {
	sample := map[string]interface{}{
		"supu":  42,
		"turbo": false,
		"foo":   "bar",
		"a": map[string]interface{}{
			"b": true,
			"c": 42,
			"d": "turbo",
		},
		"collection": []interface{}{
			map[string]interface{}{
				"b": false,
				"d": "foobar",
			},
			map[string]interface{}{
				"b": true,
				"c": 42,
				"d": "turbo",
			},
		},
	}

	res, _ := Flatten(sample, DefaultTokenizer)
	b, _ := json.MarshalIndent(res.m, "", "\t")
	fmt.Println(string(b))

	// output:
	// {
	// 	"a.b": true,
	// 	"a.c": 42,
	// 	"a.d": "turbo",
	// 	"collection.#": 2,
	// 	"collection.0.b": false,
	// 	"collection.0.d": "foobar",
	// 	"collection.1.b": true,
	// 	"collection.1.c": 42,
	// 	"collection.1.d": "turbo",
	// 	"foo": "bar",
	// 	"supu": 42,
	// 	"turbo": false
	// }
}

func ExampleFlatten_collection() {
	sample := map[string]interface{}{
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
	}
	res, _ := Flatten(sample, DefaultTokenizer)

	b, _ := json.MarshalIndent(res.m, "", "\t")
	fmt.Println(string(b))

	// output:
	// {
	// 	"a.#": 2,
	// 	"a.0.aa": 1,
	// 	"a.0.b.#": 2,
	// 	"a.0.b.0.aa": 1,
	// 	"a.0.b.0.c.a": 1,
	// 	"a.0.b.1.aa": 1,
	// 	"a.0.b.1.c.a": 2,
	// 	"a.1.aa": 1,
	// 	"a.1.b.#": 1,
	// 	"a.1.b.0.aa": 1,
	// 	"a.1.b.0.c.a": 1,
	// 	"turbo": false
	// }
}
