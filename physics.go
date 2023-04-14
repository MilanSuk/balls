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

type World struct {
	objs    []*Obj
	springs []*Spring

	gravitation Vec
	airFriction float32
}

func (world *World) Solve(dt float32) {

	//solves springs
	for _, spring := range world.springs {
		spring.Solve()
	}

	//applies enviroment forces
	for _, obj := range world.objs {
		obj.ApplyForce(world.gravitation.Mult(obj.mass))
		obj.ApplyForce(obj.vel.Mult(-world.airFriction))
	}

	world.objs[0].force = Vec{} //static(no force, no movement)

	//solves objects
	for _, obj := range world.objs {
		obj.Solve(dt)
	}
}

type Obj struct {
	mass  float32
	pos   Vec
	vel   Vec
	force Vec
}

func (obj *Obj) ApplyForce(f Vec) {
	obj.force = obj.force.Add(f)
}
func (obj *Obj) Solve(dt float32) {
	obj.vel = obj.vel.Add(obj.force.Mult(1 / obj.mass).Mult(dt))
	obj.pos = obj.pos.Add(obj.vel.Mult(dt))

	obj.force = Vec{}
}

type Spring struct {
	a *Obj
	b *Obj

	constant float32
	length   float32
	friction float32
}

func (spring *Spring) Solve() {

	springVector := spring.a.pos.Sub(spring.b.pos)

	r := springVector.Length()

	var force Vec
	if r != 0 {
		force = force.Add(springVector.Mult(-1 / r).Mult((r - spring.length) * spring.constant)) //-1 = neg
	}
	force = force.Add(spring.a.vel.Sub(spring.b.vel).Mult(spring.friction * -1)) //-1 = neg

	spring.a.ApplyForce(force)
	spring.b.ApplyForce(force.Neg())
}
