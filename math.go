/*
Copyright 2023 Milan Suk

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import "math"

type Vec struct {
	x, y, z float32
}

func (a Vec) Add(b Vec) Vec {
	return Vec{a.x + b.x, a.y + b.y, a.z + b.z}
}

func (a *Vec) EqAdd(b Vec) {
	a.x += b.x
	a.y += b.y
	a.z += b.z
}

func (a Vec) Sub(b Vec) Vec {
	return Vec{a.x - b.x, a.y - b.y, a.z - b.z}
}

func (v Vec) Mult(t float32) Vec {
	return Vec{v.x * t, v.y * t, v.z * t}
}

func (v Vec) Neg() Vec {
	return Vec{-v.x, -v.y, -v.z}
}

func (v Vec) Length() float32 {
	return float32(math.Sqrt(float64(v.x*v.x + v.y*v.y + v.z*v.z)))
}
