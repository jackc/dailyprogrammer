#include <iostream>
#include <sstream>
#include <string>
#include <array>
#include <vector>
#include <algorithm>
#include <cstdint>

typedef std::array<int64_t, 3> answer;

int main() {
  for (std::string line; std::getline(std::cin, line);) {
    std::vector<int64_t> nums;
    std::stringstream ss(line);
    while (ss && !ss.eof()) {
      int64_t num;
      ss >> num;
      nums.push_back(num);
    }
    std::sort(nums.begin(), nums.end());

    std::vector<answer> answers;
    for (auto i = nums.begin(); i != nums.end(); ++i) {
      for (auto j = std::next(i, 1); j != nums.end(); ++j) {
        for (auto k = std::next(j, 1); k != nums.end(); ++k) {
          auto sum = *i + *j + *k;
          if (sum == 0) {
            answers.push_back(answer{{*i, *j, *k}});
          }
        }
      }
    }
    std::sort(answers.begin(), answers.end());
    auto new_end = std::unique(answers.begin(), answers.end());
    answers.resize(std::distance(answers.begin(), new_end));

    for (auto a : answers) {
      std::cout << a[0] << ' ' << a[1] << ' ' << a[2] << std::endl;
    }
    std::cout << std::endl;
  }
}
