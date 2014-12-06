#pragma once

class Agent {
	int appearance, health, age;
	char direction;
	bool isAttacked;
	class Brain * brain;
	static const int MaxHealth = 8;
	void CalcFrontPos(int & x, int & y);
public:
	static const int FaceLength = 3;
	explicit Agent(int face, Brain * brn = 0, int a = 0);
	~Agent();
	int GetAppearance(void) const { return appearance ; }
	int GetHealth(void) const { return health ; }
	int GetAge(void) const { return age ; }
	unsigned SeeBrain(void) const ;
	int MutateFace(void) const ;
	Brain * MutateBrain(void) const ;
	void Act(int x, int y, class Environment * env);
};