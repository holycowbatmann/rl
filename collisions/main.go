package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	BLOCK_SIZE = 1
)

type Block struct {
	position rl.Vector3
	color    rl.Color
	model    rl.Model
}

func newBlock(position rl.Vector3, color rl.Color) *Block {
	b := &Block{}
	b.position = position
	b.color = color
	b.model = rl.LoadModelFromMesh(rl.GenMeshCube(BLOCK_SIZE, BLOCK_SIZE, BLOCK_SIZE))
	b.updateTransform()
	return b
}

func (b *Block) updateTransform() {
	b.model.Transform = rl.MatrixTranslate(b.position.X, b.position.Y, b.position.Z)
}

func (b *Block) draw() {
	b.updateTransform()
	rl.DrawModel(b.model, rl.NewVector3(0, 0, 0), 1, b.color)
}

type BlockManager struct {
	blocks []*Block
}

func newBlockManager() *BlockManager {
	bm := &BlockManager{}

	for i := 0; i <= 50; i++ {
		bm.blocks = append(bm.blocks, newBlock(rl.NewVector3(float32(rl.GetRandomValue(0, 25)), 0, float32(rl.GetRandomValue(0, 25))), rl.Green))
	}

	return bm
}

func (bm *BlockManager) draw() {
	for _, b := range bm.blocks {
		b.draw()
	}
}

type Player struct {
	position rl.Vector3
	model    rl.Model
	color    rl.Color
}

func newPlayer(position rl.Vector3) *Player {
	p := &Player{}
	p.position = position
	p.color = rl.Blue
	p.model = rl.LoadModelFromMesh(rl.GenMeshCube(1, 1, 1))
	return p
}

func (p *Player) updateTransform() {
	p.model.Transform = rl.MatrixTranslate(p.position.X, p.position.Y, p.position.Z)
}

func (p *Player) update(g *Game) {

	g.camera.Target = rl.NewVector3(p.position.X, p.position.Y, p.position.Z)

	speed := float32(0.05)
	if rl.IsKeyDown(rl.KeyW) {
		p.position.Z += speed
	}
	if rl.IsKeyDown(rl.KeyS) {
		p.position.Z -= speed
	}
	if rl.IsKeyDown(rl.KeyA) {
		p.position.X += speed
	}
	if rl.IsKeyDown(rl.KeyD) {
		p.position.X -= speed
	}

	p.updateTransform()

	bbPlayer := rl.GetModelBoundingBox(p.model)
	for _, b := range g.blockManager.blocks {
		if rl.CheckCollisionBoxes(bbPlayer, rl.GetModelBoundingBox(b.model)) {
			b.color = rl.Red
			p.resolveCollision(b)
		} else {
			b.color = rl.Green
		}
	}
}

type CubeSide int

const (
	Z CubeSide = iota
	ZMinus
	X
	XMinus
	Y
	YMinus
)

type CubeSideValue struct {
	side  CubeSide
	value float32
}

func newCubeSideValue(side CubeSide, value float32) CubeSideValue {
	v := CubeSideValue{}
	v.side = side
	v.value = value
	return v
}

func (p *Player) resolveCollision(b *Block) {
	bbBlock := rl.GetModelBoundingBox(b.model)

	switch p.getCollisionWhichSide(bbBlock) {
	case Z:
		depth := p.getCollisionDepth(bbBlock, Z)
		p.position.Z += depth
	case ZMinus:
		depth := p.getCollisionDepth(bbBlock, ZMinus)
		p.position.Z -= depth
	case X:
		depth := p.getCollisionDepth(bbBlock, X)
		p.position.X += depth
	case XMinus:
		depth := p.getCollisionDepth(bbBlock, XMinus)
		p.position.X -= depth
	case Y:
		depth := p.getCollisionDepth(bbBlock, Y)
		p.position.Y += depth
	case YMinus:
		depth := p.getCollisionDepth(bbBlock, YMinus)
		p.position.Y -= depth
	}

	p.updateTransform()
}

func (p *Player) getCollisionWhichSide(bbBlock rl.BoundingBox) CubeSide {
	var depths []CubeSideValue
	depths = append(depths, newCubeSideValue(Z, p.getCollisionDepth(bbBlock, Z)))
	depths = append(depths, newCubeSideValue(ZMinus, p.getCollisionDepth(bbBlock, ZMinus)))
	depths = append(depths, newCubeSideValue(X, p.getCollisionDepth(bbBlock, X)))
	depths = append(depths, newCubeSideValue(XMinus, p.getCollisionDepth(bbBlock, XMinus)))
	depths = append(depths, newCubeSideValue(Y, p.getCollisionDepth(bbBlock, Y)))
	depths = append(depths, newCubeSideValue(YMinus, p.getCollisionDepth(bbBlock, YMinus)))

	min := depths[0]
	for _, d := range depths {
		if d.value < min.value {
			min = d
		}
	}

	return min.side
}

func (p *Player) getCollisionDepth(bbBlock rl.BoundingBox, side CubeSide) float32 {
	oldPosition := p.position

	var dir rl.Vector3
	switch side {
	case Z:
		dir = rl.NewVector3(0, 0, 1)
	case ZMinus:
		dir = rl.NewVector3(0, 0, -1)
	case X:
		dir = rl.NewVector3(1, 0, 0)
	case XMinus:
		dir = rl.NewVector3(-1, 0, 0)
	case Y:
		dir = rl.NewVector3(0, 1, 0)
	case YMinus:
		dir = rl.NewVector3(0, 1, 0)

	}

	step := float32(0.025)
	stepVec := rl.Vector3Scale(dir, step)

	totalStep := step
	p.position = rl.Vector3Add(p.position, stepVec)
	p.updateTransform()
	for rl.CheckCollisionBoxes(rl.GetModelBoundingBox(p.model), bbBlock) {
		totalStep += step
		p.position = rl.Vector3Add(p.position, stepVec)
		p.updateTransform()
	}

	p.position = oldPosition
	p.updateTransform()

	return totalStep
}

func (p *Player) draw() {
	p.updateTransform()
	rl.DrawModel(p.model, rl.Vector3Zero(), 1, p.color)
}

type Game struct {
	camera       rl.Camera3D
	blockManager *BlockManager
	player       *Player
}

func initGame() *Game {
	g := &Game{}

	camera := rl.Camera3D{}
	camera.Position = rl.NewVector3(12, 3, 0)
	camera.Target = rl.NewVector3(12, 0.0, 12)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Fovy = 60.0
	camera.Projection = rl.CameraPerspective
	g.camera = camera

	g.blockManager = newBlockManager()

	g.player = newPlayer(rl.NewVector3(12, 0, 5))
	return g
}

func (g *Game) update() {
	g.player.update(g)
}

func (g *Game) draw() {
	rl.BeginMode3D(g.camera)
	g.blockManager.draw()
	g.player.draw()
	rl.EndMode3D()
}

func main() {
	rl.InitWindow(1024, 768, "Collisions")
	rl.SetTargetFPS(60)

	g := initGame()

	for !rl.WindowShouldClose() {

		g.update()

		rl.BeginDrawing()

		rl.ClearBackground(rl.NewColor(0, 0, 0, 255))

		g.draw()

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
