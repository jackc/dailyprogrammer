P=ssa-interleave
CXXFLAGS = -g -Wall -Werror -O3 -std=c++14 -pthread
LDLIBS=-lboost_program_options
CXX=g++

$(P): main.o
	$(CXX) $(CXXFLAGS) main.o -o $@ $(LDFLAGS) $(LDLIBS)

main.o: main.cpp world.h
	$(CXX) $(CXXFLAGS) -c $< -o $@
