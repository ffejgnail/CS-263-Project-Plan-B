package main

type Direction uint8

const (
	Up Direction = iota
	Left
	Down
	Right
)

type Face uint8

type Agent struct {
	Health, Age int
	Face        Face
	IsAttacked  bool
	IsAttacking bool
	Brain       Brain
	Direction   Direction
}

func isPowerOfTwo(x uint8) bool {
	return (x & (x - 1)) == 0
}

func isSimilar(a, b Face) bool {
	return isPowerOfTwo(uint8(a) ^ uint8(b))
}

func relLoc(x, y, rx, ry int) (int, int) {
	return (x + rx + EnvSize) % EnvSize, (y + ry + EnvSize) % EnvSize
}

func nextLoc(x, y int, dir Direction) (int, int) {
	switch dir {
	case Up:
		return relLoc(x, y, 0, -1)
	case Left:
		return relLoc(x, y, -1, 0)
	case Down:
		return relLoc(x, y, 0, 1)
	case Right:
		return relLoc(x, y, 1, 0)
	default:
		return x, y
	}
}

func (ag *Agent) Act(x, y int, env *Environment) {
	x2, y2 := nextLoc(x, y, ag.Direction)
	ag2 := env.Cell[x2][y2].Agent
	input := uint32(0)
	if ag.IsAttacked {
		input += 16
	}
	if env.Cell[x2][y2].Food > env.Cell[x][y].Food {
		input += 8
	}
	if ag2 == nil {
		input += 4
	}
	if ag2 != nil && ag2.Health > ag.Health {
		input += 2
	}
	if ag2 != nil && isSimilar(ag2.Face, ag.Face) {
		input++
	}
	decision := (ag.Brain.react(input)) & 1
	ag.IsAttacked = false
	ag.IsAttacking = false
	if decision == 0 {
		ag.Direction = (ag.Direction + 1) & 3
	} else {
		if ag2 == nil {
			env.Cell[x2][y2].Agent = ag
			env.Cell[x][y].Agent = nil
		} else {
			ag2.Health -= int(1+env.Cell[x][y].Food)
			ag2.IsAttacked = true
			ag.IsAttacking = true
			ag2.Direction = ag.Direction ^ 2
			env.Aggressiveness[0]++
		}
	}
}