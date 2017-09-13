let srand = require('srand');
const yargs = require('yargs')
  .option('benchmark', {describe: "run benchmark of N steps and quit", default: 0})
  .option('height', {describe: "height of the world", default: 30})
  .option('max', {describe: "max value", default: 15})
  .option('penalty', {describe: "penalty value", default: 1})
  .option('reward', {describe: "reward value", default: 3})
  .option('seed', {describe: "seed"})
  .option('steptime', {describe: "steptime", default: 250})
  .option('target', {describe: "target value", default: 5})
  .option('width', {describe: "width of the world", default: 30})
  .help()
  .argv

class World {
  constructor(height, width) {
    this.cells = new Array(height * width);
    this.height = height;
    this.width = width;
  }

  idx_from_coord(x, y) {
    x = x % this.width;
    if (x < 0) x += this.width;

    y = y % this.height;
    if (y < 0) y += this.height;

    return y*this.width + x;
  }

  get(x, y) {
    return this.cells[this.idx_from_coord(x, y)];
  }

  set(x, y, val) {
    this.cells[this.idx_from_coord(x, y)] = val;
  }
}

function detect_subset_sum(w, x, y, want) {
  cells = [
    w.get(x-1, y-1),
    w.get(x, y-1),
    w.get(x+1, y-1),

    w.get(x-1, y),
    w.get(x+1, y),

    w.get(x-1, y+1),
    w.get(x, y+1),
    w.get(x+1, y+1)
  ];

  for (let i = 1; i < 256; i++) {
    let sum = 0;
    for (let j = 0; j < 8; j++) {
      if (i&(1<<j) != 0) {
        sum += cells[j];
      }
    }

    if (sum == want) {
      return true;
    }
  }

  return false;
}

function step(curr, nxt, target, reward, penalty, max) {
  for (let y = 0; y < curr.height; y++) {
    for (let x = 0; x < curr.width; x++) {
      let new_value = curr.get(x, y);
      if (detect_subset_sum(curr, x, y, target)) {
        new_value += reward;
      } else {
        new_value -= penalty;
      }

      if (new_value < 0) new_value = 0;
      else if (new_value > max) new_value = max;

      if (new_value < 0) new_value = 0;
      else if (new_value > max) new_value = max;

      nxt.set(x, y, new_value);
    }
  }
}

function print_world(w) {
  for (let i = 0; i < 80; i++) {
    process.stdout.write("\n");
  }

  for (let y = 0; y < w.height; y++) {
    for (let x = 0; x < w.width; x++) {
      process.stdout.write(` ${w.get(x, y).toString(36)}`);
    }
    process.stdout.write("\n");
  }
}

if (yargs.seed) srand.seed(yargs.seed);
else srand.seed(Date.now());

let w = new World(yargs.width, yargs.height);
let w_scratch = new World(yargs.width, yargs.height);

for (let y = 0; y < w.height; y++) {
  for (let x = 0; x < w.width; x++) {
    w.set(x, y, srand.rand() % (yargs.target+yargs.reward));
  }
}

if (yargs.benchmark > 0) {
  for (let i = 0; i < yargs.benchmark; i++) {
    step(w, w_scratch, yargs.target, yargs.reward, yargs.penalty, yargs.max);
    w, w_scratch = w_scratch, w;
  }

  print_world(w);
  process.exit();
}

setInterval(function() {
  print_world(w);
  step(w, w_scratch, yargs.target, yargs.reward, yargs.penalty, yargs.max);
  let t = w;
  w = w_scratch;
  w_scratch = t;
}, yargs.steptime);
