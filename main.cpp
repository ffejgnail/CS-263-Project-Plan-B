#include "Environment.h"
#include <iostream>
#include <time.h>
#include <stdlib.h>
using namespace std;

int main() {
	srand(unsigned(time(NULL)));
	//srand(29);
	Environment *p = new Environment;
	p->Setup();
	int n;
	cin >> n;
	for(int i = 0; i < n; i++) {
		p->Run(i);
	}
	if(!p->Check()) cout << "World matrix corrupted." << endl;
	system("pause");
	p->Print();
	p->Show(n);
	system("pause");
	return 0;
}