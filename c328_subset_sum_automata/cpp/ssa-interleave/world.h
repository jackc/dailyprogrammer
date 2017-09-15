#pragma once

#include <cstdint>
#include <vector>

template <class CellType>
class world {
  std::vector<CellType> cells;
  int32_t width;
  int32_t height;

  std::size_t idx_from_coord(int32_t x, int32_t y);
public:
  world(int32_t width, int32_t height);
  int32_t get_width() { return width; };
  int32_t get_height() { return height; };

  void set(int32_t x, int32_t y, CellType val);
  CellType get(int32_t x, int32_t y);
};

template <class CellType>
world<CellType>::world(int32_t width, int32_t height) {
  this->width = width;
  this->height = height;
  cells.resize(width*height);
}

// idx_from_coord takes x and y coordinates and returns the index in cells.
// Coordinates wrap the boundaries of the world. e.g. Given world with a
// width of 10, then an x coordinate of -1 should be equal to 9.
template <class CellType>
std::size_t world<CellType>::idx_from_coord(int32_t x, int32_t y) {
  x = x % width;
  if (x < 0) {
    x += width;
  }
  y = y % height;
  if (y < 0) {
    y += height;
  }

  return y*width + x;
}

template <class CellType>
void world<CellType>::set(int32_t x, int32_t y, CellType val) {
  cells[idx_from_coord(x, y)] = val;
}

template <class CellType>
CellType world<CellType>::get(int32_t x, int32_t y) {
  return cells[idx_from_coord(x, y)];
}

