#include <algorithm>
#include <chrono>
#include <cstdint>
#include <cstdlib>
#include <iomanip>
#include <iostream>
#include <random>
#include <string>
#include <thread>
#include <utility>

#include <boost/program_options.hpp>

#include "world.h"

struct config {
  int32_t width;
  int32_t height;

  int16_t target;
  int16_t penalty;
  int16_t reward;
  int16_t max;

  int seed;

  int steptime;
  int64_t benchmark;
};

config parse_command_line(int argc, char** argv) {
  using namespace boost::program_options;

  options_description desc{"Options"};
  desc.add_options()
    ("benchmark", value<int64_t>()->default_value(0), "run benchmark of N steps and quit")
    ("height", value<int32_t>()->default_value(30), "height of the world")
    ("help,h", "Print this help message")
    ("max", value<int16_t>()->default_value(15), "max value")
    ("penalty", value<int16_t>()->default_value(1), "penalty value")
    ("reward", value<int16_t>()->default_value(3), "reward value")
    ("seed", value<int>()->default_value(-1), "seed")
    ("steptime", value<int>()->default_value(250), "time per step in milliseconds")
    ("target", value<int16_t>()->default_value(5), "target value")
    ("width", value<int32_t>()->default_value(30), "width of the world")
  ;

  variables_map vm;
  store(parse_command_line(argc, argv, desc), vm);
  notify(vm);

  if (vm.count("help")) {
    std::cout << desc << '\n';
    std::exit(0);
  }

  config c = {};
  c.benchmark = vm["benchmark"].as<int64_t>();
  c.height = vm["height"].as<int32_t>();
  c.max = vm["max"].as<int16_t>();
  c.penalty = vm["penalty"].as<int16_t>();
  c.reward = vm["reward"].as<int16_t>();
  c.seed = vm["seed"].as<int>();
  c.steptime = vm["steptime"].as<int>();
  c.target = vm["target"].as<int16_t>();
  c.width = vm["width"].as<int32_t>();

  return c;
}

template <class CellType>
void print(world<CellType> w) {
  std::cout << std::string(80, '\n');

  for (int32_t y = 0; y < w.get_height(); y++) {
    for (int32_t x = 0; x < w.get_width(); x++) {
      std::cout << std::right << std::setw(3) << w.get(x, y);
    }
    std::cout << "\n";
  }

  std::cout.flush();
}

template <class CellType>
bool detect_subset_sum(world<CellType>& w, int32_t x, int32_t y, CellType want) {
  std::array<CellType, 8> cells = {
    w.get(x-1, y-1),
    w.get(x, y-1),
    w.get(x+1, y-1),

    w.get(x-1, y),
    w.get(x+1, y),

    w.get(x-1, y+1),
    w.get(x, y+1),
    w.get(x+1, y+1),
  };

  for (uint32_t i = 1; i < 256; i++) {
    CellType sum = 0;
    for (uint32_t j = 0; j < 8; j++) {
      if ((i&(1<<j)) != 0) {
        sum += cells[j];
      }
    }

    if (sum == want) return true;
  }

  return false;
}

template <class CellType>
void step(world<CellType>& curr, world<CellType> &next, CellType target, CellType reward, CellType penalty, CellType max) {
  for (int32_t y = 0; y < curr.get_height(); y++) {
    for (int32_t x = 0; x < curr.get_width(); x++) {
      auto new_value = curr.get(x, y);
      if (detect_subset_sum(curr, x, y, target)) {
        new_value += reward;
      } else {
        new_value -= penalty;
      }

      new_value = std::max(CellType(0), new_value);
      new_value = std::min(max, new_value);

      next.set(x, y, new_value);
    }
  }
}

int main(int argc, char** argv) {
  try {
    auto config = parse_command_line(argc, argv);

    world<int16_t> w(config.width, config.height);
    world<int16_t> w_scratch(config.width, config.height);

    std::default_random_engine prng;
    if (config.seed != -1) {
      prng.seed(config.seed);
    } else {
      std::random_device r;
      prng.seed(r());
    }
    std::uniform_int_distribution<int16_t> distribution(0, config.target + config.reward - 1);

    for (int32_t y = 0; y < w.get_height(); y++) {
      for (int32_t x = 0; x < w.get_width(); x++) {
        w.set(x, y, distribution(prng));
      }
    }

    if (config.benchmark > 0) {
      for (int64_t i = 0; i < config.benchmark; i++) {
        step(w, w_scratch, config.target, config.reward, config.penalty, config.max);
        std::swap(w, w_scratch);
      }
      print(w); // Access results to ensure entire calculation cannot be removed by the optimizer.
      std::exit(0);
    }

    while (true) {
      print(w);
      step(w, w_scratch, config.target, config.reward, config.penalty, config.max);
      std::swap(w, w_scratch);
      std::this_thread::sleep_for(std::chrono::milliseconds(config.steptime));
    }
  }
  catch (const std::exception &ex) {
    std::cerr << ex.what() << '\n';
    return 1;
  }

  return 0;
}
