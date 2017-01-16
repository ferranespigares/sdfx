package sdf

import (
	"github.com/deadsy/pt/pt"
)

func Render(s *SDF3, render_floor bool) {

	scene := pt.Scene{}

	light := pt.LightMaterial(pt.White, 180)

	d := 4.0
	scene.Add(pt.NewSphere(pt.V(-1, -1, 0.5).Normalize().MulScalar(d), 0.25, light))
	scene.Add(pt.NewSphere(pt.V(0, -1, 0.25).Normalize().MulScalar(d), 0.25, light))
	scene.Add(pt.NewSphere(pt.V(-1, 1, 0).Normalize().MulScalar(d), 0.25, light))

	material := pt.GlossyMaterial(pt.HexColor(0x468966), 1.2, pt.Radians(20))

	s0 := NewPtSDF(s)
	//s0 = pt.NewTransformSDF(s0, pt.Translate(pt.V(0, 0, 0.2)))
	//s0 = pt.NewTransformSDF(s0, pt.Rotate(pt.V(0, 0, 1), pt.Radians(30)))

	scene.Add(pt.NewSDFShape(s0, material))

	if render_floor {
		bb := s0.BoundingBox()
		z_min := bb.Min.Z
		z_height := bb.Max.Z - bb.Min.Z
		z_gap := z_height * 0.1 // 10% of height

		floor := pt.GlossyMaterial(pt.HexColor(0xFFF0A5), 1.2, pt.Radians(20))
		floor_plane := pt.V(0, 0, z_min-z_gap)
		floor_normal := pt.V(0, 0, 1)

		scene.Add(pt.NewPlane(floor_plane, floor_normal, floor))
	}

	camera := pt.LookAt(pt.V(-3, 0, 1), pt.V(0, 0, 0), pt.V(0, 0, 1), 35)
	sampler := pt.NewSampler(4, 4)
	sampler.LightMode = pt.LightModeAll
	sampler.SpecularMode = pt.SpecularModeAll
	renderer := pt.NewRenderer(&scene, &camera, sampler, 800, 600)
	renderer.IterativeRender("out%03d.png", 10)
}
