use std::collections::HashSet;
use std::fs::File;
use std::io::{BufRead, BufReader};
use itertools::Itertools;

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

    let mut sum: i32 = 0;
    for line in reader.lines() {
        let line = line.unwrap();
        let (a, b) = line.split_at(line.len()/2);
        println!("{} : {} ", a, b);
        let mut set: HashSet<char> = HashSet::new();
        for c in a.chars() {
            set.insert(c);
        };
        for c in b.chars() {
            if set.contains(&c) {
                sum += get_score(&c);
                println!("found mismactch for {} w/score {} in {} : {}", c, get_score(&c), a, b);
                break
            };
        }
    }
    println!("the sum is {}", sum)
}

fn task2(){
    let filename = "input";
    let file = File::open(filename).unwrap();
    let reader = BufReader::new(file);

    let mut sum: i32 = 0;
    for (l1, l2, l3) in reader.lines().tuples() {
        let (l1, l2, l3) = (l1.unwrap(), l2.unwrap(), l3.unwrap());
        let mut set1: HashSet<char> = HashSet::new();
        let mut set2: HashSet<char> = HashSet::new();
        for c in l1.chars() {
            set1.insert(c);
        };
        for c in l2.chars() {
            set2.insert(c);
        };
        for c in l3.chars() {
            if set1.contains(&c) && set2.contains(&c) {
                sum += get_score(&c);
                print!("found badge for {} w/score {} in {} : {} : {}\n", c, get_score(&c), l1, l2, l3);
                break
            };
        };
    }
    println!("the sum is {}", sum)

}

fn get_score(c: &char) -> i32 {
    // let b = *c as u8;
    if c.is_ascii_lowercase() {
        return *c as i32 - 'a' as i32 + 1
    } else if c.is_ascii_uppercase() {
        return *c as i32 - 'A' as i32 + 27
    }
    -1
}