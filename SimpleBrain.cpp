#include "SimpleBrain.h"
#include <stdlib.h>

bool SimpleBrain::React(bool isAttacked, bool moreGrass, bool isEmpty, bool isStronger, bool isSimilar) const {
	int index = 0;
	if(isAttacked)
		index += 16;
	if(moreGrass)
		index += 8;
	if(isEmpty)
		index += 4;
	if(isStronger)
		index += 2;
	if(isSimilar)
		index ++;
	return ((1 << index) & lut) != 0;
}

Brain * SimpleBrain::Mutate(void) const {
	unsigned n = lut;
	for(int i = 0; i < MutationRate; i++) {
		n ^= (1 << (rand() & 31));
	}
	return new SimpleBrain(n);
}