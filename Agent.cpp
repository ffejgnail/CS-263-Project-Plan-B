#include "Agent.h"
#include "Environment.h"
#include "SimpleBrain.h"
#include "MyError.h"
#include <stdlib.h>

void Agent::CalcFrontPos(int & x, int & y) {
	const unsigned siz = Environment::EnvSize;
	switch(direction) {
		case 0: x--; break; // u
		case 1: y++; break; // r
		case 2: x++; break; // d
		case 3: y--; break; // l
		default: MyError("Unknown direction.");
	}
	x = (x + siz) % siz;
	y = (y + siz) % siz;
}

Agent::Agent(int face, Brain * brn, int a) {
	appearance = face;
	health = MaxHealth;
	age = a;
	direction = rand() & 3;
	isAttacked = false;
	if(brn)
		brain = brn;
	else
		brain = new SimpleBrain(~0); // magic number: 0xFFFFF300
}

Agent::~Agent() {
	delete brain;
}

unsigned Agent::SeeBrain(void) const {
	return ((const SimpleBrain*)brain)->Show();
}

int Agent::MutateFace(void) const {
	return appearance ^ (1 << (rand() % FaceLength));
}

Brain * Agent::MutateBrain(void) const {
	return brain->Mutate();
}

inline bool IsPowerOfTwo(unsigned x) {
	return (x & (x - 1)) == 0;
}

inline bool IsSimilar(unsigned a, unsigned b) {
	return IsPowerOfTwo(a ^ b);
}

void Agent::Act(int x, int y, Environment * env) {
	const int grass = env->GetGrass(x, y);
	int x2 = x, y2 = y;
	CalcFrontPos(x2, y2);
	Agent * ag2 = env->GetAgent(x2, y2);
	bool b = brain->React(isAttacked, env->GetGrass(x2, y2) > grass, ag2 == NULL, ag2 && ag2->health > health,
		ag2 && IsSimilar(ag2->appearance, appearance));
	isAttacked = false;
	if(b){
		if(ag2){
			ag2->health -= grass;
			ag2->isAttacked = true;
			ag2->direction = direction ^ 2;
		}
		else{
			env->Move(x, y, x2, y2);
		}
	}
	else{
		direction = (direction + 1) & 3;
	}
}