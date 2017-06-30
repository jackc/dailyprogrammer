use std::cmp;

pub fn is_embedded(world: &str, test: &str) -> bool {
    let mut world_iter = world.chars();

    for c in test.chars() {
        loop {
            match world_iter.next() {
                Some(wc) => {
                    if c == wc {
                        break;
                    }
                }
                None => return false
            }
        }
    }

    true
}

pub fn embed(dest: &str, word: &str) -> String {
    embed_slide(dest, word);
    let l = embed_left(dest, word);
    let r = embed_right(dest, word);
    if l.len() < r.len() {
        return l;
    }
    r
}

// embed_append is the simplest possible working embed
fn embed_append(dest: &str, word: &str) -> String {
    let mut s = dest.to_owned();
    s.push_str(word);
    s
}

fn embed_left(dest: &str, word: &str) -> String {
    if dest.len() == 0 {
        return word.to_owned();
    }
    if word.len() == 0 {
        return dest.to_owned();
    }

    let mut s = String::new();

    let mut word_iter = word.chars().rev();
    let mut wc_option = word_iter.next();

    for dc in dest.chars().rev() {
        s.insert(0, dc);
        if let Some(wc) = wc_option {
            if wc == dc {
                wc_option = word_iter.next()
            }
        }
    }

    if let Some(wc) = wc_option {
        s.insert(0, wc);
        for wc in word_iter {
            s.insert(0, wc);
        }
    }

    s
}

fn embed_right(dest: &str, word: &str) -> String {
    if dest.len() == 0 {
        return word.to_owned();
    }
    if word.len() == 0 {
        return dest.to_owned();
    }

    let mut s = String::new();

    let mut word_iter = word.chars();
    let mut wc_option = word_iter.next();

    for dc in dest.chars() {
        s.push(dc);
        if let Some(wc) = wc_option {
            if wc == dc {
                wc_option = word_iter.next()
            }
        }
    }

    if let Some(wc) = wc_option {
        s.push(wc);
        for wc in word_iter {
            s.push(wc);
        }
    }

    s
}

#[derive(Copy, Clone, Debug)]
struct Insert {
    beforeIdx: usize,
    value: char,
}

fn embed_insert(dest: &[char], word: &[char], destOffset: usize) -> Vec<Insert> {


    vec![Insert{beforeIdx: 0, value: ' '}]
}

fn embed_slide(dest: &str, word: &str) -> String {
    let window_size = cmp::min(4, word.len());

    let chars: Vec<_> = word.chars().collect();
    println!("{:?}", chars);
    "foo".to_owned()
}


#[cfg(test)]
mod test {
    use super::*;

    #[test]
    fn test_is_embedded() {
        assert!(is_embedded("tahsifs", "this"));
        assert!(!is_embedded("tahsifs", "hat"));
    }

    #[test]
    fn test_embed() {
        let mut s = embed("", "this");
        println!("{}", s);
        assert!(is_embedded(&s, "this"));

        s = embed(&s, "is");
        println!("{}", s);
        assert!(is_embedded(&s, "this"));
        assert!(is_embedded(&s, "is"));

        s = embed(&s, "that");
        println!("{}", s);
        assert!(is_embedded(&s, "this"));
        assert!(is_embedded(&s, "is"));
        assert!(is_embedded(&s, "that"));
    }

}
