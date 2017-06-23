// For this challenge I am purposely avoiding regular expresssions to make the
// challenge more interesting.

use std::string::String;
use std::vec::Vec;
use std::cmp;

pub fn condense(src: &str) -> String {
    let mut out : Vec<char> = Vec::new();

    for word in src.split(' ').map(|w| w.chars().collect::<Vec<char>>()) {
        let max_overlap = cmp::min(word.len(), out.len());
        let mut overlap_idx = 0;

        for i in 0..max_overlap {
            let possible_out_overlap = &out[(out.len()-max_overlap+i)..];
            let possible_word_overlap = &word[..(max_overlap-i)];

            println!("{} {:?} {:?}", i, possible_out_overlap.iter().collect::<String>(), possible_word_overlap.iter().collect::<String>());

            if possible_out_overlap == possible_word_overlap {
                overlap_idx = max_overlap - i;
                break;
            }
        }
        if overlap_idx == 0 && out.len() > 0 {
            out.push(' ');
        }

        for c in &word[overlap_idx..] {
            out.push(*c)
        }
    }

    return out.iter().collect::<String>();
}

#[cfg(test)]
mod test {
    use super::*;

    #[test]
    fn condense_sample() {
        let input = "Deep episodes of Deep Space Nine came on the television only after the news.
Digital alarm clocks scare area children.";

        let output = "Deepisodes of Deep Space Nine came on the televisionly after the news.
Digitalarm clockscarea children.";

        assert_eq!(output, condense(input));
    }
}
