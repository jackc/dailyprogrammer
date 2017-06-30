extern crate c313_embedded_word_list;

use std::io::{self, Read};
use c313_embedded_word_list::*;

fn main() {
    let mut src_buffer = String::new();
    io::stdin().read_to_string(&mut src_buffer).expect("failed to read stdin");

    let mut dest = String::new();
    for line in src_buffer.lines() {
        dest = embed(&dest, line);
        println!("{}", line);
    }

    println!("{} (length {})", dest, dest.len());
}
