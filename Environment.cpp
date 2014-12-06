#pragma once
#include "Environment.h"
#include "Agent.h"
#include "MyError.h"
#include <iostream>
#include <bitset>
#include <stdlib.h>
using namespace std;

void Environment::Move(int srcX, int srcY, int dstX, int dstY) {
	if(world[dstX][dstY].agent != -1) MyError("Agents cannot overlap.");
	AgentInfo & info = vAgents[world[srcX][srcY].agent];
	if(info.x != srcX || info.y != srcY) MyError("Already corrupted.");
	info.x = dstX;
	info.y = dstY;
	vAgents[world[srcX][srcY].agent].y = dstY;
	world[dstX][dstY].agent = world[srcX][srcY].agent;
	world[srcX][srcY].agent = -1;
}

int Environment::GetGrass(int x, int y) {
	return world[x][y].grass;
}

Agent * Environment::GetAgent(int x, int y) {
	const int index = world[x][y].agent;
	if(index == -1)
		return NULL;
	return vAgents[index].p;
}

void Environment::PickFreeLoc(int ag) {
	int x, y;
	do {
		x = rand() % EnvSize;
		y = rand() % EnvSize;
	} while(world[x][y].agent != -1);
	world[x][y].agent = ag;
	vAgents[ag].x = x;
	vAgents[ag].y = y;
}

void Environment::Setup(void) {
	vAgents.resize(AgentNum);
	world.resize(EnvSize);
	for(unsigned i = 0; i < EnvSize; i++) {
		world[i].resize(EnvSize);
		for(unsigned j = 0; j < EnvSize; j++) {
			world[i][j].grass = 10; //1+(rand() & 3);
			world[i][j].agent = -1;
		}
	}
	/*
	for(unsigned i = 0; i < 4; i++) {
		for(unsigned j = 0; j < EnvSize; j++) {
			world[4+i][j].grass = 4-i;
			world[3-i][j].grass = 4-i;
			world[j][4+i].grass = 4-i;
			world[j][3-i].grass = 4-i;
		}
	}*/
	for(unsigned i = 0; i < AgentNum; i++) {
		PickFreeLoc(i);
		vAgents[i].p = new Agent(i & ((1 << (Agent::FaceLength)) - 1));
	}
}

void Environment::Run(int iter) {
	for(unsigned i = 0; i < vAgents.size(); i++) {
		AgentInfo & info  = vAgents[i];
		if(info.p->GetHealth() < 0){
			delete info.p;
			world[info.x][info.y].agent = -1;
			PickFreeLoc(i);
			int bestIndex = -1, bestAge = -1, bestGrass = 0;
			for(unsigned j = 0; j < vAgents.size(); j++) {
				if(j == i) continue;
				int age = iter - vAgents[j].p->GetAge();
				int grass = world[vAgents[j].x][vAgents[j].y].grass;
				if(age > bestAge || (age == bestAge && grass < bestGrass)) {
					bestIndex = j;
					bestAge = age;
					bestGrass = grass;
				}
			}
			info.p = new Agent(vAgents[bestIndex].p->MutateFace(), vAgents[bestIndex].p->MutateBrain(), iter);
		}
		else{
			info.p->Act(info.x, info.y, this);
		}
	}
}

bool Environment::Check(void) {
	if(vAgents.size() != AgentNum) return false;
	int cnt = 0;
	for(unsigned i = 0; i < EnvSize; i++) {
		for(unsigned j = 0; j < EnvSize; j++) {
			if(world[i][j].agent == -1) continue;
			AgentInfo & info = vAgents[world[i][j].agent];
			if(info.x != i || info.y != j) {
				return false;
			}
			cnt++;
		}
	}
	if(cnt != AgentNum) return false;
	return true;
}

void Environment::Print(void) {
	for(unsigned i = 0; i < EnvSize; i++) {
		for(unsigned j = 0; j < EnvSize; j++) {
			if(world[i][j].agent == -1) {
				cout << "_|";
			}
			else{
				int hp = vAgents[world[i][j].agent].p->GetHealth();
				cout << (hp < 0 ? 0 : hp + 1) << "|";
			}
		}
		cout << endl;
	}
	cout << endl;

	for(unsigned i = 0; i < EnvSize; i++) {
		for(unsigned j = 0; j < EnvSize; j++) {
			if(world[i][j].agent == -1) {
				cout << "_|";
			}
			else{
				int face = vAgents[world[i][j].agent].p->GetAppearance();
				cout << face << "|";
			}
		}
		cout << endl;
	}
	cout << endl;
}

void Environment::Show(int n) {
	/*
	vector<vector<int>> v;
	v.resize(1 << (Agent::FaceLength));
	for(unsigned i = 0; i < vAgents.size(); i++) {
		unsigned brain = vAgents[i].p->SeeBrain();
		int favorSimilar = 0;
		for(int j = 0; j < 32; j++) {
			if((brain ^ j) & 1)
				favorSimilar++;
			else
				favorSimilar--;
			brain >>= 1;
		}
		v[vAgents[i].p->GetAppearance()].push_back(favorSimilar);
	}
	for(unsigned i = 0; i < v.size(); i++) {
		for(unsigned j = 0; j < v[i].size(); j++) {
			cout << v[i][j] << " ";
		}
		cout << endl;
	} */
	for(unsigned i = 0; i < vAgents.size(); i++) {
		cout << bitset<32>(vAgents[i].p->SeeBrain()) << " ";
	}
	cout << endl << endl;
	for(unsigned i = 0; i < vAgents.size(); i++) {
		cout << n - vAgents[i].p->GetAge() << " ";
	}
	cout << endl;
}