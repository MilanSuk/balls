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

import (
	"runtime"
	"sync"
)

const NUM_MAX_THREADS = 8

type World struct {
	objs    []*Obj
	springs []*Spring

	gravitation Vec
	airFriction float32
}

func NewWorld(num_threads int) *World {
	var world World
	return &world
}

func (world *World) solveSprings(thread_i int, num_threads int, wg *sync.WaitGroup) {

	defer wg.Done()

	jmpStrings := (len(world.springs) / num_threads) + 1
	st := thread_i * jmpStrings
	en := (thread_i + 1) * jmpStrings
	if en > len(world.springs) {
		en = len(world.springs)
	}
	for i := st; i < en; i++ {
		world.springs[i].Solve(thread_i)
	}
}

func (world *World) solveObjects(dt float32, thread_i int, num_threads int, wg *sync.WaitGroup) {

	defer wg.Done()

	jmpObjs := (len(world.objs) / num_threads) + 1
	st := thread_i * jmpObjs
	en := (thread_i + 1) * jmpObjs
	if en > len(world.objs) {
		en = len(world.objs)
	}

	for i := st; i < en; i++ {
		obj := world.objs[i]
		//applies enviroment forces
		obj.ApplyForce(world.gravitation.Mult(obj.mass), thread_i)
		obj.ApplyForce(obj.vel.Mult(-world.airFriction), thread_i)

		if !obj.static {
			obj.Solve(dt)
		}
	}
}

func (world *World) Solve(dt float32, num_threads int) {

	if num_threads < 0 {
		num_threads = runtime.NumCPU()
	}
	if num_threads > NUM_MAX_THREADS {
		num_threads = NUM_MAX_THREADS
	}

	var wg sync.WaitGroup

	//solves springs
	for i := 0; i < num_threads; i++ {
		wg.Add(1)
		go world.solveSprings(i, num_threads, &wg)
	}
	wg.Wait()

	//solve objects
	for i := 0; i < num_threads; i++ {
		wg.Add(1)
		go world.solveObjects(dt, i, num_threads, &wg)
	}
	wg.Wait()
}

type Obj struct {
	mass  float32
	pos   Vec
	vel   Vec
	force [NUM_MAX_THREADS]Vec

	static bool
}

func (obj *Obj) ApplyForce(f Vec, thread_i int) {
	obj.force[thread_i].EqAdd(f)
}
func (obj *Obj) Solve(dt float32) {

	//sum & reset
	var force Vec
	for i := 0; i < NUM_MAX_THREADS; i++ {
		force.EqAdd(obj.force[i])
		obj.force[i] = Vec{}
	}

	obj.vel.EqAdd(force.Mult(1 / obj.mass).Mult(dt))
	obj.pos.EqAdd(obj.vel.Mult(dt))
}

type Spring struct {
	a *Obj
	b *Obj

	constant float32
	length   float32
	friction float32
}

func (spring *Spring) Solve(thread_i int) {

	springVector := spring.a.pos.Sub(spring.b.pos)

	r := springVector.Length()

	var force Vec
	if r != 0 {
		force.EqAdd(springVector.Mult(-1 / r).Mult((r - spring.length) * spring.constant)) //-1 = neg
	}
	force.EqAdd(spring.a.vel.Sub(spring.b.vel).Mult(spring.friction * -1)) //-1 = neg

	spring.a.ApplyForce(force, thread_i)
	spring.b.ApplyForce(force.Neg(), thread_i)
}
