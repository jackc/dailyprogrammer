#include <algorithm>
#include <condition_variable>
#include <chrono>
#include <cstdint>
#include <cstdlib>
#include <iomanip>
#include <iostream>
#include <memory>
#include <mutex>
#include <queue>
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

  bool parallel;
};

config parse_command_line(int argc, char** argv) {
  using namespace boost::program_options;

  options_description desc{"Options"};
  desc.add_options()
    ("benchmark", value<int64_t>()->default_value(0), "run benchmark of N steps and quit")
    ("height", value<int32_t>()->default_value(30), "height of the world")
    ("help,h", "Print this help message")
    ("max", value<int16_t>()->default_value(15), "max value")
    ("parallel", "parallel execution")
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
  if (vm.count("parallel")) c.parallel = true;
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
class stepper {
protected:
  static void step_range(world<CellType>& curr, world<CellType>& next, CellType target, CellType reward, CellType penalty, CellType max, int32_t y_start, int32_t y_end);

public:
  virtual void step(world<CellType>& curr, world<CellType>& next, CellType target, CellType reward, CellType penalty, CellType max) = 0;
};

template <class CellType>
void stepper<CellType>::step_range(world<CellType>& curr, world<CellType>& next, CellType target, CellType reward, CellType penalty, CellType max, int32_t y_start, int32_t y_end) {
  for (int32_t y = y_start; y < y_end; y++) {
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

template <class CellType>
class stepper_serial : public stepper<CellType> {
public:
  virtual void step(world<CellType>& curr, world<CellType>& next, CellType target, CellType reward, CellType penalty, CellType max);
};

template <class CellType>
void stepper_serial<CellType>::step(world<CellType>& curr, world<CellType>& next, CellType target, CellType reward, CellType penalty, CellType max) {
  stepper<CellType>::step_range(curr, next, target, reward, penalty, max, 0, curr.get_height());
}

template <class CellType>
class stepper_parallel : public stepper<CellType> {
  // missing destructor to terminate threads as threads run entire time of program
  std::vector<std::thread> threads;

  void work();

  struct job {
    world<CellType>& curr;
    world<CellType>& next;
    int32_t y_start;
    int32_t y_end;
    CellType target;
    CellType reward;
    CellType penalty;
    CellType max;
  };

  std::mutex mutex;
  std::condition_variable cond_var;
  std::queue<job> job_queue;

  std::mutex pending_mutex;
  std::condition_variable pending_cond_var;
  int pending;

public:
  stepper_parallel();
  stepper_parallel(const stepper_parallel& src) = delete;
  stepper_parallel& operator=(const stepper_parallel& rhs) = delete;

  virtual void step(world<CellType>& curr, world<CellType>& next, CellType target, CellType reward, CellType penalty, CellType max);
};

template <class CellType>
stepper_parallel<CellType>::stepper_parallel() {
  auto thread_count = std::thread::hardware_concurrency();
  for (unsigned int i = 0; i < thread_count; i++) {
    threads.emplace_back( &stepper_parallel<CellType>::work, this );
  }
}

template <class CellType>
void stepper_parallel<CellType>::work() {
  std::unique_lock<std::mutex> lock(mutex);

  while (true) {
    cond_var.wait(lock);

    lock.unlock();

    while (true) {
      lock.lock();
      if (job_queue.empty()) {
        break;
      }

      auto j = job_queue.front();
      job_queue.pop();
      lock.unlock();

      stepper<CellType>::step_range(j.curr, j.next, j.target, j.reward, j.penalty, j.max, j.y_start, j.y_end);

      std::lock_guard<std::mutex> pending_lock(pending_mutex);
      pending--;
      pending_cond_var.notify_all();
    }
  }
}

template <class CellType>
void stepper_parallel<CellType>::step(world<CellType>& curr, world<CellType>& next, CellType target, CellType reward, CellType penalty, CellType max) {
  pending = 0;
  {
    std::lock_guard<std::mutex> lock(mutex);

    int32_t y_start = 0;
    int32_t rows_per_thread = (curr.get_height() / threads.size()) + 1;

    for (std::vector<std::thread>::size_type i = 0; i < threads.size(); i++) {
      int32_t y_end = y_start+rows_per_thread;
      if (y_end > curr.get_height()) y_end = curr.get_height();
      job_queue.push({curr, next, y_start, y_end, target, reward, penalty, max});
      y_start += rows_per_thread;
      pending++;
    }
  }
  cond_var.notify_all();

  std::unique_lock<std::mutex> lock(pending_mutex);
  while (pending) {
    pending_cond_var.wait(lock);
  }
}

int main(int argc, char** argv) {
  try {
    auto config = parse_command_line(argc, argv);

    world<int16_t> w(config.width, config.height);
    world<int16_t> w_scratch(config.width, config.height);
    std::unique_ptr<stepper<int16_t>> stepper;
    if (config.parallel) {
      stepper.reset(new stepper_parallel<int16_t>);
    } else {
      stepper.reset(new stepper_serial<int16_t>);
    }

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
        stepper->step(w, w_scratch, config.target, config.reward, config.penalty, config.max);
        std::swap(w, w_scratch);
      }
      print(w); // Access results to ensure entire calculation cannot be removed by the optimizer.
      std::exit(0);
    }

    while (true) {
      print(w);
      stepper->step(w, w_scratch, config.target, config.reward, config.penalty, config.max);
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
