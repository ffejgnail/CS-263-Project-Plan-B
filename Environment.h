#pragma once
#include <vector>
using std::vector;

struct EnvCell {
	int agent;
	int grass;
};

struct AgentInfo {
	int x, y;
	class Agent * p;
};

class Environment {
	vector<AgentInfo> vAgents;
	vector<vector<EnvCell>> world;
	static const unsigned AgentNum = 160;
	void PickFreeLoc(int ag);
public:
	static const unsigned EnvSize = 16;
	void Move(int srcX, int srcY, int dstX, int dstY);
	int GetGrass(int x, int y);
	Agent * GetAgent(int x, int y);

	void Setup(void);
	void Run(int iter);
	bool Check(void);
	void Print(void);
	void Show(int n);
};