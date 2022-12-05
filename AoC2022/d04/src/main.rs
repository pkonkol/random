use std::fs::File;
use std::io::{BufRead, BufReader};

fn main() {
    println!("{:-<15}task1{:-<15}", "", "");
    task1();
    println!("{:-<15}task2{:-<15}", "", "");
    task2();
}

fn task1() {
    let filename = "input";
    let file = File::open(filename).unwrap();
    let reader = BufReader::new(file);

    let mut cnt: i32 = 0;
    for line in reader.lines() {
        let line = line.unwrap();
        let mut split = line.split(",");
        let (a, b) = (split.next().unwrap(), split.next().unwrap());
        println!("{} : {} ", a, b);
        let mut split = a.split("-");
        let (al, ar): (i32, i32) = (split.next().unwrap().parse().unwrap(), split.next().unwrap().parse().unwrap());
        let mut split = b.split("-");
        let (bl, br): (i32, i32) = (split.next().unwrap().parse().unwrap(), split.next().unwrap().parse().unwrap());
        print!("{} : {}; al:ar:bl:br {}:{}:{}:{}", a, b, al, ar, bl, br);

        if al >= bl && ar <= br {
            cnt += 1;
            print!("; a contained in b\n")
        } else if bl >= al && br <= ar {
            cnt += 1;
            print!("; b contained in a\n")
        }
    }
    println!("the cnt is {}", cnt);
}

fn task2(){
    let filename = "input";
    let file = File::open(filename).unwrap();
    let reader = BufReader::new(file);

    let mut cnt: i32 = 0;
    for line in reader.lines() {
        let line = line.unwrap();
        let mut split = line.split(",");
        let (a, b) = (split.next().unwrap(), split.next().unwrap());
        println!("{} : {} ", a, b);
        let mut split = a.split("-");
        let (al, ar): (i32, i32) = (split.next().unwrap().parse().unwrap(), split.next().unwrap().parse().unwrap());
        let mut split = b.split("-");
        let (bl, br): (i32, i32) = (split.next().unwrap().parse().unwrap(), split.next().unwrap().parse().unwrap());
        print!("{} : {}; al:ar:bl:br {}:{}:{}:{}", a, b, al, ar, bl, br);

        if al < bl {
            if ar >= bl {
                cnt += 1;
            }
            print!("; a overlaping in b\n")
        } else if bl < al {
            if br >= al {
                cnt += 1;
            }
            print!("; b overlaping in a\n")
        } else {
            cnt += 1;
            print!("; al == bl\n")
        }
    }
    println!("the cnt is {}", cnt);
}
