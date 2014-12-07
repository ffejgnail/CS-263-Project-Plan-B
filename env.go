package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math/rand"
)

type EnvCell struct {
	Food   uint8 // 0, 1, 2, 3, 4
	Agent *Agent
}

type Environment struct {
	Cell    [EnvSize][EnvSize]EnvCell
	record  gif.GIF
}

func (env *Environment) WriteTo(w io.Writer) error { // generating GIF
	return gif.EncodeAll(w, &env.record)
}

func (env *Environment) relCell(x, y, rx, ry int) *EnvCell {
	x2, y2 := relLoc(x, y, rx, ry)
	return &env.Cell[x2][y2]
}

func (env *Environment) pickFreeLoc () (int, int) {
	x := rand.Intn(EnvSize)
	y := rand.Intn(EnvSize)
	for env.Cell[x][y].Agent != nil {
		x = rand.Intn(EnvSize)
		y = rand.Intn(EnvSize)
	}
	return x, y
}

func NewEnvironment() *Environment {
	env := new(Environment)

	cx := EnvSize / 2
	cy := EnvSize / 2
	for i := -4; i <= 4; i++ {
		for j := -4; j <= 4; j++ {
			for k := 0; k <= 4; k++ {
				if i*i+j*j <= k*k {
					env.relCell(cx, cy, i, j).Food = uint8(4 - k)
					break
				}
			}
		}
	}

	for i := uint8(0); i < InitAnimatNum; i++ {
		brain := new(SimpleBrain)
		brain.lut = ^uint32(0)
		x, y := env.pickFreeLoc()
		ag := &Agent{
			Health:      InitHealth,
			Age:         Iteration,
			Face:        Face(i & ((1 << FaceLength) - 1)),
			IsAttacked:  false,
			IsAttacking: false,
			Brain:       brain,
			Direction:   Direction(rand.Intn(4)),
		}
		env.Cell[x][y].Agent = ag
	}

	env.record.Image = make([]*image.Paletted, RecordIteration)
	env.record.Delay = make([]int, RecordIteration)
	for i := 0; i < RecordIteration; i++ {
		env.record.Delay[i] = RecordDelay
	}
	return env
}

var (
	backgroundColor = color.RGBA{255, 255, 255, 255}
	gridColor       = color.RGBA{0, 0, 0, 255}
	grassColor1     = color.RGBA{30, 60, 30, 255}
	grassColor2     = color.RGBA{60, 120, 60, 255}
	grassColor3     = color.RGBA{90, 180, 90, 255}
	grassColor4     = color.RGBA{120, 240, 120, 255}
	grassColor5     = color.RGBA{240, 240, 120, 255}
	attackColor     = color.RGBA{255, 0, 0, 255}
	normalColor     = color.RGBA{0, 255, 0, 255}
)

func grassColor(grass uint8) color.Color {
	if grass > 3 {
		return grassColor1
	}
	if grass > 2 {
		return grassColor2
	}
	if grass > 1 {
		return grassColor3
	}
	if grass > 0 {
		return grassColor4
	}
	return grassColor5
}

func (env *Environment) Run(iter int) {
	type Location struct {
		Agent *Agent
		X     int
		Y     int
	}
	var list []*Location
	for i := 0; i < EnvSize; i++ {
		for j := 0; j < EnvSize; j++ {
			cell := &env.Cell[i][j]
			if cell.Agent == nil {
				continue
			}
			list = append(list, &Location{
				Agent: cell.Agent,
				X:     i,
				Y:     j,
			})
		}
	}
	for i := range list {
		if list[i].Agent.Health > 0 {
			continue
		}
		longestLife := -1
		bestBrain := ^uint32(0)
		var bestFace Face
		for j := range list {
			if list[j].Agent.Health > 0 && list[j].Agent.Age > longestLife{
				longestLife = list[j].Agent.Age
				bestBrain = list[j].Agent.Brain.getGene()
				bestFace = list[j].Agent.Face
			}
		}
		env.Cell[list[i].X][list[i].Y].Agent = nil
		list[i].X, list[i].Y = env.pickFreeLoc()
		env.Cell[list[i].X][list[i].Y].Agent = list[i].Agent
		list[i].Agent.Health = InitHealth
		list[i].Agent.Age = Iteration - iter
		list[i].Agent.Face = Face(uint8(bestFace) ^ uint8(1 << uint8(rand.Intn(FaceLength))))
		list[i].Agent.IsAttacked = false
		list[i].Agent.IsAttacking = false
		list[i].Agent.Brain.resetWithGene(bestBrain ^ (1 << uint32(rand.Intn(32))))
		list[i].Agent.Direction = Direction(rand.Intn(4))
	}
	for i := range list { // _, r := range rand.Perm(len(list))
		list[i].Agent.Act(list[i].X, list[i].Y, env)
	}
	if Iteration-RecordIteration <= iter {
		env.drawFrame(iter - Iteration + RecordIteration)
	}
}

func (env *Environment) drawFrame(iter int) {
	img := image.NewPaletted(image.Rect(0, 0, (CellPixel+1)*EnvSize-1, (CellPixel+1)*EnvSize-1), []color.Color{
		backgroundColor,
		gridColor,
		grassColor1,
		grassColor2,
		grassColor3,
		grassColor4,
		grassColor5,
		attackColor,
		normalColor,
	})

	for i := 1; i < EnvSize; i++ {
		for j := 0; j < (CellPixel+1)*EnvSize-1; j++ {
			img.Set((CellPixel+1)*i, j, gridColor)
			img.Set(j, (CellPixel+1)*i, gridColor)
		}
	}
	for i := 0; i < EnvSize; i++ {
		for j := 0; j < EnvSize; j++ {
			cell := &env.Cell[i][j]

			for ii := i*(CellPixel+1) + 1; ii < (i+1)*(CellPixel+1); ii++ {
				for jj := j*(CellPixel+1) + 1; jj < (j+1)*(CellPixel+1); jj++ {
					img.Set(ii, jj, grassColor(cell.Food))
				}
			}

			if cell.Agent == nil {
				continue
			}
			a1 := i*(CellPixel+1) + MarginPixel + 1
			a2 := (i+1)*(CellPixel+1) - MarginPixel
			b1 := j*(CellPixel+1) + MarginPixel + 1
			b2 := (j+1)*(CellPixel+1) - MarginPixel
			for ii := a1; ii < a2; ii++ {
				for jj := b1; jj < b2; jj++ {
					if ii == a1 || ii == a2-1 ||
						jj == b1 || jj == b2-1 {
						img.Set(ii, jj, gridColor)
					} else {
						img.Set(ii, jj, backgroundColor)
					}
				}
			}
			ag := cell.Agent
			switch ag.Direction {
			case Up:
				for ii := a1 + 1; ii < a2-1; ii++ {
					for jj := b1 + 1; jj < b1+1+HeadPixel; jj++ {
						if ag.IsAttacking {
							img.Set(ii, jj, attackColor)
						} else {
							img.Set(ii, jj, normalColor)
						}
					}
				}
			case Left:
				for ii := a1 + 1; ii < a1+1+HeadPixel; ii++ {
					for jj := b1 + 1; jj < b2-1; jj++ {
						if ag.IsAttacking {
							img.Set(ii, jj, attackColor)
						} else {
							img.Set(ii, jj, normalColor)
						}
					}
				}
			case Down:
				for ii := a1 + 1; ii < a2-1; ii++ {
					for jj := b2 - 2; jj > b2-2-HeadPixel; jj-- {
						if ag.IsAttacking {
							img.Set(ii, jj, attackColor)
						} else {
							img.Set(ii, jj, normalColor)
						}
					}
				}
			case Right:
				for ii := a2 - 2; ii > a2-2-HeadPixel; ii-- {
					for jj := b1 + 1; jj < b2-1; jj++ {
						if ag.IsAttacking {
							img.Set(ii, jj, attackColor)
						} else {
							img.Set(ii, jj, normalColor)
						}
					}
				}
			}
		}
	}
	env.record.Image[iter] = img
}