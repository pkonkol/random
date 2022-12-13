use std::fs::File;
use std::io::{BufRead, BufReader};
use serde::Deserialize;


fn main() {
    println!("{:-<15}task1{:-<15}", "", "");
    task1();
    println!("{:-<15}task2{:-<15}", "", "");
    // task2();
}

fn get_input() -> BufReader<File> {
    let filename = "example";
    let file = File::open(filename).unwrap();
    BufReader::new(file)
}

fn task1() {
    let reader = get_input();

    // let mut left = String::new();
    // let mut right = String::new();
    let mut lines = reader.lines();
    while let (Some(left), Some(right), _) = (lines.next(), lines.next(), lines.next()) {
        let left = left.unwrap();
        let right = right.unwrap();
        print!("left is  {}\nright is {}\n", left, right);
        // TODO deserialize it in serde so I can skip writing my own parser
        // let deserialized: Vec<_> = serde_json::from_str::<Vec<_>>(&left).unwrap();
    }
}

fn task2(){
}
