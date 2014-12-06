#pragma once

class Brain {
public:
	virtual bool React(bool isAttacked, bool moreGrass, bool isEmpty, bool isStronger, bool isSimilar) const = 0;
	virtual Brain * Mutate(void) const = 0;
};