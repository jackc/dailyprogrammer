#[derive(Debug)]
pub enum Op {
    Add,
    Sub,
    Mul,
    Div,
}

#[derive(Debug)]
pub struct Term {
    pub op: Op,
    pub val: i32,
}



impl Term {
    pub fn apply(&self, prev: i32) -> i32 {
        match self.op {
           Op::Add => prev + self.val,
           Op::Sub => prev - self.val,
           Op::Mul => prev * self.val,
           Op::Div => prev / self.val,
        }
    }
}

#[derive(Debug)]
pub struct Formula {
    pub terms: Vec<Term>,
}

impl Formula {
    pub fn eval(&self) -> i32 {
        let mut term_iter = self.terms.iter();
        let first_val = term_iter.next().unwrap().val;
        term_iter.fold(first_val, |acc, term| term.apply(acc))
    }
}

// 1 2
// + 1 + 2
// + 1 - 2
// + 1 / 2
// + 1 * 2
// + 2 + 1
// + 2 - 1
// + 2 / 1
// + 2 * 1

pub fn solve(numbers: &[i32], goal: i32) -> Formula {
    // operations
    // cartesian of a

    // let mut all_possibilies = someiterator(numbers, ops)
    let mut terms = numbers.iter().map(|number| Term{op: Op::Div, val: *number}).collect();

    Formula{terms: terms}
}

pub struct RepeatedPermutationIterator<T> {
    indices: Vec<usize>,
    alphabet: Vec<T>,
}

impl<T> RepeatedPermutationIterator<T> {
    pub fn new(alphabet: Vec<T>, size: usize) -> RepeatedPermutationIterator<T> {
        let mut iter = RepeatedPermutationIterator{
            indices: vec![0, size],
            alphabet: alphabet,
        };
        iter.indices[0] -= 1;
        iter
    }
}

impl<T> Iterator for RepeatedPermutationIterator<T> {
    type Item = Vec<T>;

    fn next(&mut self) -> Option<Vec<T>> {
        self.indices[0] += 1;

        let mut carry_idx = 0;
        while carry_idx < self.indices.len() && self.indices[carry_idx] == self.alphabet.len() {
            self.indices[carry_idx] = 0;
            carry_idx += 1;
            if carry_idx < self.indices.len() {
                self.indices[carry_idx] += 1;
            }
        }

        if carry_idx < self.indices.len() {
            Some(self.indices.iter().map(|i| &self.alphabet[*i].clone()).collect())
        } else {
            None
        }
    }
}
