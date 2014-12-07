package main

type Brain interface {
	react(uint32) uint8
	getGene() uint32
	resetWithGene(uint32)
}

type SimpleBrain struct{
	lut uint32
}

func (sb *SimpleBrain) react(input uint32) uint8 {
	return uint8((sb.lut >> input) & 1)
}

func (sb *SimpleBrain) getGene () uint32 {
	return sb.lut
}

func (sb *SimpleBrain) resetWithGene (newGene uint32) {
	sb.lut = newGene
}