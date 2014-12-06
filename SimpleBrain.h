#pragma once
#include "Brain.h"

class SimpleBrain : public Brain {
	unsigned lut;
	static const int MutationRate = 1;
public:
	SimpleBrain(unsigned n) : lut(n) { ; }
	bool React(bool isAttacked, bool moreGrass, bool isEmpty, bool isStronger, bool isSimilar) const ;
	Brain * Mutate(void) const ;
	unsigned Show(void) const { return lut ; }
};