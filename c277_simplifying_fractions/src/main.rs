extern crate c277_simplifying_fractions;

use c277_simplifying_fractions::Fraction;

fn main() {
    println!("{:?}", Fraction {n: 4, d: 6}.simplify());
}


