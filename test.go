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

func WorldTest() *World {
	var world World
	world.gravitation.y = 9.81 //negative ...
	world.airFriction = 0.02

	const SPRING_LEN = 0.05

	/*//adds objects
	const N = 200
	p := Vec{x: 3, y: 1}
	for i := 0; i < N; i++ {
		world.objs = append(world.objs, &Obj{mass: 0.05, pos: p})
		p.x += SPRING_LEN
	}

	//adds springs
	for i := 1; i < N; i++ {
		world.springs = append(world.springs, &Spring{a: world.objs[i-1], b: world.objs[i], constant: 10000, length: SPRING_LEN, friction: 0.2})
	}*/

	//adds objects
	const N = 50
	p := Vec{x: 3, y: 1}
	for y := 0; y < N; y++ {
		for x := 0; x < N; x++ {
			world.objs = append(world.objs, &Obj{mass: 0.05, pos: Vec{x: 3 + float32(x)*SPRING_LEN, y: 1 + float32(y)*SPRING_LEN}})
			p.x += SPRING_LEN
		}
	}

	//adds springs
	for y := 0; y < N; y++ {
		for x := 1; x < N; x++ {
			iA := y*N + (x - 1)
			iB := y*N + (x + 0)
			world.springs = append(world.springs, &Spring{a: world.objs[iA], b: world.objs[iB], constant: 10000, length: SPRING_LEN, friction: 0.2})
		}
	}
	for y := 1; y < N; y++ {
		for x := 0; x < N; x++ {
			iA := (y-1)*N + x
			iB := (y-0)*N + x
			world.springs = append(world.springs, &Spring{a: world.objs[iA], b: world.objs[iB], constant: 10000, length: SPRING_LEN, friction: 0.2})
		}
	}

	return &world
}
