use std::io::BufRead;
use std::io;

fn main() {
    let stdin = io::stdin();

    for line in stdin.lock().lines() {
        match line {
            Ok(line) => {
                let words: Vec<&str> = line.split_whitespace().collect();
                if words.len() != 2 {
                    eprintln!("Invalid input. Requires exactly two words.");
                    continue;
                }

                println!("{} {} => {}", words[0], words[1], funnel(words[0], words[1]));
            },
            Err(err) => panic!("IO error: {}", err),
        }
    }
}

fn funnel(source: &str, candidate: &str) -> bool {
    let expected_skips = 1;
    let mut skips = 0;

    let source: Vec<char> = source.chars().collect();
    let candidate: Vec<char> = candidate.chars().collect();

    if source.len() - expected_skips != candidate.len() {
        return false;
    }

    let mut candidate_iter = candidate.iter();
    let mut candidate_char = candidate_iter.next();

    for s in source {
        match candidate_char {
            Some(&c) if c == s => {
                candidate_char = candidate_iter.next();
            },
            _ => {
                skips += 1;
                if skips > expected_skips {
                    return false;
                }
            }

        }
    }

    true
}
