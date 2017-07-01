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
    let dest_chars: Vec<_> = dest.chars().collect();
    let word_chars: Vec<_> = word.chars().collect();
    let insert_list = embed_insert(&dest_chars, &word_chars);

    let mut insert_iter = insert_list.iter();
    let mut insert_option = insert_iter.next();

    let mut s = String::with_capacity(dest_chars.len() + word_chars.len());

    for (i, dc) in dest_chars.iter().enumerate() {
        loop {
            if let Some(insert) = insert_option {
                if insert.beforeIdx == i {
                    s.push(insert.value);
                    insert_option = insert_iter.next();
                } else {
                    break;
                }
            } else {
                break;
            }
        }
        s.push(*dc);
    }

    if let Some(insert) = insert_option {
        s.push(insert.value);
        for insert in insert_iter {
            s.push(insert.value);
        }
    }

    s
}

#[derive(Copy, Clone, Debug)]
struct Insert {
    beforeIdx: usize,
    value: char,
}

fn embed_insert(dest: &[char], word: &[char]) -> Vec<Insert> {
    if dest.len() == 0 {
        return word.iter().map(|c| { Insert{beforeIdx: 0, value: *c} }).collect();
    }

    if word.len() == 0 {
        return Vec::new();
    }

    let mut insert_list = Vec::with_capacity(word.len());

    let offset = dest.len();
    for insert in embed_insert(&dest[offset..], word) {
        insert_list.push(Insert{beforeIdx: offset+insert.beforeIdx, value: insert.value});
    }

    return insert_list;




    let mut word_iter = word.iter();
    let mut wc_option = word_iter.next();

    for dc in dest.iter() {
        if let Some(wc) = wc_option {
            if wc == dc {
                wc_option = word_iter.next()
            }
        }
    }

    if let Some(wc) = wc_option {
        insert_list.push(Insert{beforeIdx: dest.len(), value: *wc});
        for wc in word_iter {
            insert_list.push(Insert{beforeIdx: dest.len(), value: *wc});
        }
    }

    insert_list
}

// fn next_insert(dest: &[char], value: char) -> Insert {
//     for dc in dest.iter() {
//         if let Some(wc) = wc_option {
//             if wc == dc {
//                 wc_option = word_iter.next()
//             }
//         }
//     }
// }

// fn embed_slide(dest: &str, word: &str) -> String {
//     let dest_chars: Vec<_> = dest.chars().collect();
//     let word_chars: Vec<_> = word.chars().collect();
//     let insert_list = embed_insert(&dest_chars, &word_chars, 0);


//     let window_size = cmp::min(4, word.len());

//     let chars: Vec<_> = word.chars().collect();
//     println!("{} {} {:?}", dest, word, insert_list);
//     "foo".to_owned()
// }


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
