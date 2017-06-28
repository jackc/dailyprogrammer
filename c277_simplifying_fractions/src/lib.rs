#[derive(Copy, Clone, Debug, PartialEq)]
pub struct Fraction {
    pub n: i64,
    pub d: i64,
}

impl Fraction {
    pub fn simplify(self) -> Fraction {
        let mut n = 42;
        let mut s = n;
        n = 50;

        let mut f = self;

        let mut pos_div = 2;

        while pos_div <= f.n && pos_div <= f.d {
            if f.n % pos_div == 0 && f.d % pos_div == 0 {
                f = Fraction{n: f.n / pos_div, d: f.d / pos_div};
            } else {
                pos_div += 1;
            }
        }
        f
    }
}

#[cfg(test)]
mod test {
    use super::*;

    #[test]
    fn fraction_simplify() {
        assert_eq!(Fraction {n: 4, d: 8}.simplify(), Fraction {n: 1, d: 2});
        assert_eq!(Fraction {n: 1, d: 1}.simplify(), Fraction {n: 1, d: 1});
        assert_eq!(Fraction {n: 4, d: 2}.simplify(), Fraction {n: 2, d: 1});
        assert_eq!(Fraction {n: 9, d: 54}.simplify(), Fraction {n: 1, d: 6});
        assert_eq!(Fraction {n: 1, d: 101}.simplify(), Fraction {n: 1, d: 101});
        assert_eq!(Fraction {n: 20, d: 100}.simplify(), Fraction {n: 1, d: 5});
    }
}
