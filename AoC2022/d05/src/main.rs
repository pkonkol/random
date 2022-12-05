use std::fs::File;
use std::io::{BufRead, BufReader, Read};

fn main() {
    println!("{:-<15}task1{:-<15}", "", "");
    task1();
    println!("{:-<15}task2{:-<15}", "", "");
    task2();
}

fn task1() {
    let filename = "input";
    let file = File::open(filename).unwrap();
    let mut reader = BufReader::new(file);
    // read until empty line
    let mut init: Vec<String> = Vec::new();
    for line in reader.by_ref().lines().take_while(|c| {c.as_ref().unwrap() != ""}) {
        let l = line.unwrap();
        if l == "" {
            break
        }
        init.push(l);
    }
    print!("len(l): {} n stacks: {}\n", init[0].len(), (init[0].len()+1)/4);
    print!("stacks are: {:?}\n", init);
    let mut stacks: Vec<Vec<char>> = Vec::new();
    for _ in 0..(init[0].len()+1)/4 {
        stacks.push(Vec::<char>::new());
    }
    for l in init.iter().rev().skip(1) {
        for (i, x) in l.chars().skip(1).step_by(4).enumerate() {
            //println!("taking stack row with {} at col {} ", x, i);
            if x.is_ascii_uppercase() {
                stacks[i].push(x);
            }
        }
        //println!("next row ");
    }
    print!("stacks are: {:?}\n", stacks);

    // read the rest, take columns 2,4,6
    for line in reader.lines() {
        let l = line.unwrap();
        let mut iter = l.split_whitespace().skip(1).step_by(2);
        let n = iter.next().unwrap().parse().unwrap();
        let src: usize = iter.next().unwrap().parse().unwrap();
        let dst: usize = iter.next().unwrap().parse().unwrap();
        for _ in 0..n {
            let x = stacks[src-1].pop();
            stacks[dst-1].push(x.unwrap());
            println!("repeating action for: move {} src {} dst {}", n, src, dst);
        }
        //println!("action is: move {} src {} dst {}", n, src, dst);
    }
    print!("stacks after operations are: {:?}\n", stacks);
    for s in stacks.iter() {
        print!("{}", s.last().unwrap());
    }
    println!("");
}

fn task2(){
    let filename = "input";
    let file = File::open(filename).unwrap();
    let mut reader = BufReader::new(file);
    // read until empty line
    let mut init: Vec<String> = Vec::new();
    for line in reader.by_ref().lines().take_while(|c| {c.as_ref().unwrap() != ""}) {
        let l = line.unwrap();
        if l == "" {
            break
        }
        init.push(l);
    }
    print!("len(l): {} n stacks: {}\n", init[0].len(), (init[0].len()+1)/4);
    print!("stacks are: {:?}\n", init);
    let mut stacks: Vec<Vec<char>> = Vec::new();
    for _ in 0..(init[0].len()+1)/4 {
        stacks.push(Vec::<char>::new());
    }
    for l in init.iter().rev().skip(1) {
        for (i, x) in l.chars().skip(1).step_by(4).enumerate() {
            //println!("taking stack row with {} at col {} ", x, i);
            if x.is_ascii_uppercase() {
                stacks[i].push(x);
            }
        }
        //println!("next row ");
    }
    print!("stacks are: {:?}\n", stacks);

    // read the rest, take columns 2,4,6
    for line in reader.lines() {
        let l = line.unwrap();
        let mut iter = l.split_whitespace().skip(1).step_by(2);
        let n: usize = iter.next().unwrap().parse().unwrap();
        let src: usize = iter.next().unwrap().parse().unwrap();
        let dst: usize = iter.next().unwrap().parse().unwrap();
        //for _ in 0..n {
        let len = stacks[src-1].len();
        let mut x = stacks[src-1].split_off(len-n);
        println!("split off {:?} for: move {} src {} dst {}", x, n, src, dst);
        stacks[dst-1].append(&mut x);
        //}
        //println!("action is: move {} src {} dst {}", n, src, dst);
    }
    print!("stacks after operations are: {:?}\n", stacks);
    for s in stacks.iter() {
        print!("{}", s.last().unwrap());
    }
    println!("");
}
