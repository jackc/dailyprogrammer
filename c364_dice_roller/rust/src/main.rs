extern crate rand;
extern crate regex;

use rand::prelude::*;
use std::io::BufRead;
use regex::Regex;
use std::io;

fn main() {
    let stdin = io::stdin();
    let mut rng = thread_rng();

    for line in stdin.lock().lines() {
        match line {
            Ok(line) => {
                let (dice_count, dice_sides) = match parse_dice_def(&line) {
                    Ok(t) => t,
                    Err(s) => {
                        println!("{}", s);
                        continue
                    }
                };

                let mut dice_results: Vec<u32> = Vec::with_capacity(dice_count as usize);
                let mut sum: u32 = 0;

                for _ in 0..dice_count {
                    let roll = rng.gen_range(1, dice_sides+1);
                    sum += roll;
                    dice_results.push(roll);
                }



                print!("{}:", sum);

                for roll in dice_results.iter() {
                    print!(" {}", roll);
                }

                println!();
            },
            Err(err) => panic!("IO error: {}", err),
        }
    }
}

fn parse_dice_def(s: &str) -> Result<(u32, u32), &str> {
    let re = Regex::new(r"\A(\d{1,3})d(\d{1,3})\z").unwrap();

    let caps = match re.captures(&s) {
        Some(caps) => caps,
        None => return Err("Invalid dice format"),
    };

    let count: u32 = caps.get(1).unwrap().as_str().parse().unwrap();
    let sides: u32 = caps.get(2).unwrap().as_str().parse().unwrap();

    Ok((count, sides))
}
