use std::collections::{VecDeque, HashMap};
use std::fs::File;
use std::io::{BufReader, Read};

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

    let mut deq = VecDeque::<char>::new();
    let mut cnt = HashMap::<char, i32>::new();
    let mut c: char;
    let mut old_c: char;
    for (i, b) in reader.bytes().enumerate() {
        c = b.unwrap() as char;
        *cnt.entry(c).or_insert(0) += 1;
        if i < 4 {
            deq.push_front(c);
            continue
        }
        deq.push_front(c);
        old_c = deq.pop_back().unwrap();
        *cnt.get_mut(&old_c).unwrap() -= 1;
        if *cnt.get(&old_c).unwrap() == 0 {
            cnt.remove(&old_c);
        }
        if cnt.len() == 4 {
            print!("len == 4 for {}, result is {}\n", i, i+1);
            break;
        }
        print!("for {}:{} --- deq is {:?} --- cnt is {:?}\n", i, c, deq, cnt);
    }
}

fn task2(){
    let filename = "input";
    let file = File::open(filename).unwrap();
    let reader = BufReader::new(file);

    let mut deq = VecDeque::<char>::new();
    let mut cnt = HashMap::<char, i32>::new();
    let mut c: char;
    let mut old_c: char;
    for (i, b) in reader.bytes().enumerate() {
        c = b.unwrap() as char;
        *cnt.entry(c).or_insert(0) += 1;
        if i < 14 {
            deq.push_front(c);
            continue
        }
        deq.push_front(c);
        old_c = deq.pop_back().unwrap();
        *cnt.get_mut(&old_c).unwrap() -= 1;
        if *cnt.get(&old_c).unwrap() == 0 {
            cnt.remove(&old_c);
        }
        if cnt.len() == 14 {
            print!("len == 4 for {}, result is {}\n", i, i+1);
            break;
        }
        print!("for {}:{} --- deq is {:?} --- cnt is {:?}\n", i, c, deq, cnt);
    }
}
