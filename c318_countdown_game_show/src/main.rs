extern crate c318_countdown_game_show;

use std::io::{self};
use c318_countdown_game_show::*;

fn main() {
    let mut line = String::new();
    io::stdin().read_line(&mut line).unwrap();
    let numbers = line
        .split_whitespace()
        .map(|s| s.parse::<i32>().unwrap() )
        .collect::<Vec<i32>>();

    let goal = numbers[numbers.len()-1];
    let numbers = &numbers[..numbers.len()-1];

    let formula = solve(numbers, goal);

    println!("{:?} -> {}", formula, formula.eval());
}
