use std::io;
use std::io::Write;

fn main() {
    print!("Enter a number:");
    io::stdout().flush().unwrap();

    let mut num = String::new();

    io::stdin()
        .read_line(&mut num)
        .expect("Failed to read line");

    let num: i32 = num.trim().parse().expect("Must be a number");

    let grid = comp_spiral(num);
    print_spiral(&grid);
}

fn comp_spiral(num: i32) -> Grid {
    let mut run_len = num;
    let mut curr_run_len = 0;
    let mut decr_run_len = true;

    let dir = vec![Vec2 {x: 1, y: 0}, Vec2 {x: 0, y: 1}, Vec2 {x: -1, y: 0}, Vec2 {x: 0, y: -1}];
    let mut dir = dir.iter().cycle();
    let mut curr_dir = dir.next().unwrap();

    let mut pos = Vec2 {x: 0, y: 0};

    let mut grid = Grid::new(num as i32);

    for i in 1..(num*num+1) {
        grid.set(pos.x, pos.y, i);
        curr_run_len += 1;
        if curr_run_len == run_len {
            curr_run_len = 0;
            curr_dir = dir.next().unwrap();
            if decr_run_len {
                run_len -= 1;
            }
            decr_run_len = !decr_run_len;
        }
        pos.x += curr_dir.x;
        pos.y += curr_dir.y;
    }

    return grid;
}

fn print_spiral(grid: &Grid) {
    let cell_width = ((grid.size() * grid.size) as f32).log(10.0).ceil() as usize;

    for y in 0..grid.size() {
        for x in 0..grid.size() {
            print!("{:width$} ", grid.get(x, y), width = cell_width);
        }
        println!("");
    }
}

#[derive(Debug)]
pub struct Vec2 {
    x: i32,
    y: i32,
}

#[derive(Debug)]
pub struct Grid {
    size: i32,
    cells: Vec<i32>,
}

impl Grid {
    pub fn new(size: i32) -> Grid {
        Grid {
            size: size,
            cells: vec![0; (size as usize) * (size as usize)],
        }
    }

    pub fn size(&self) -> i32 { self.size }

    pub fn get(&self, x: i32, y: i32) -> i32 {
        self.cells[self.coord_to_idx(x, y)]
    }

    pub fn set(&mut self, x: i32, y: i32, val: i32) {
        let idx = self.coord_to_idx(x, y);
        self.cells[idx] = val;
    }

    fn coord_to_idx(&self, x: i32, y: i32) -> usize {
        (y as usize) * (self.size as usize) + (x as usize)
    }
}
