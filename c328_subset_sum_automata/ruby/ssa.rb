class World
  attr_reader :cells, :width, :height

  def initialize(width, height)
    @cells = Array.new(width*height) { 0 }
    @width = width
    @height = height
  end

  def idx_from_coord(x, y)
    x = x % width
    if x < 0
      x += width
    end
    y = y % height
    if y < 0
      y += height
    end

    y*width + x
  end

  def get(x, y)
    cells[idx_from_coord(x, y)]
  end

  def set(x, y, val)
    cells[idx_from_coord(x, y)] = val
  end
end

def detect_subset_sum(w, x, y, want)
  cells = [
    w.get(x-1, y-1),
    w.get(x, y-1),
    w.get(x+1, y-1),

    w.get(x-1, y),
    w.get(x+1, y),

    w.get(x-1, y+1),
    w.get(x, y+1),
    w.get(x+1, y+1),
  ]

  (1...256).each do |i|
    sum = 0
    8.times do |j|
      if i&(1<<j) != 0
        sum += cells[j]
      end
    end

    if sum == want
      return true
    end
  end

  false
end

def step(curr, nxt, target, reward, penalty, max)
  curr.height.times do |y|
    curr.width.times do |x|
      new_value = curr.get(x, y)
      if detect_subset_sum(curr, x, y, target)
        new_value += reward
      else
        new_value -= penalty
      end

      nxt.set(x, y, new_value.clamp(0, max))
    end
  end
end

def print_world(w)
  80.times { puts }

  w.height.times do |y|
    w.width.times do |x|
      print w.get(x, y).to_s(36).ljust(2)
    end
    puts
  end
end

require 'slop'

opts = Slop.parse do |o|
  o.integer '--benchmark', 'run benchmark of N steps and quit', default: 0
  o.integer '--height', 'height of the world', default: 30
  o.integer '--max', 'max value', default: 15
  o.integer '--penalty', 'penalty value', default: 1
  o.integer '--reward', 'reward value', default: 3
  o.integer '--seed', 'seed'
  o.integer '--steptime', 'time per step in milliseconds', default: 250
  o.integer '--target', 'target value', default: 5
  o.integer '--width', 'width of the world', default: 30
end

if opts[:seed]
  srand(opts[:seed])
else
  srand
end

w = World.new(opts[:width], opts[:height])
w_scratch = World.new(opts[:width], opts[:height])

w.height.times do |y|
  w.width.times do |x|
    w.set(x, y, rand(opts[:target]+opts[:reward]))
  end
end

if opts[:benchmark] > 0
  opts[:benchmark].times do
    step(w, w_scratch, opts[:target], opts[:reward], opts[:penalty], opts[:max])
    w, w_scratch = w_scratch, w
  end

  print_world(w)
  exit 0
end

while true
  print_world(w)
  step(w, w_scratch, opts[:target], opts[:reward], opts[:penalty], opts[:max])
  w, w_scratch = w_scratch, w
  sleep(opts[:steptime] * 0.001)
end
