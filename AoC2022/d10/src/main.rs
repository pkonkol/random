use std::fs::File;
use std::io::{BufRead, BufReader};

const HEIGHT: i32 = 7;

fn main() {
    println!("{:-<15}task1{:-<15}", "", "");
    task1();
    println!("{:-<15}task2{:-<15}", "", "");
    task2();
}

fn get_input() -> BufReader<File> {
    let filename = "input";
    let file = File::open(filename).unwrap();
    BufReader::new(file)
}

fn task1() {
    let reader = get_input();
    let mut x: i32 = 1;
    let mut cycle: i32 = 1;
    let mut result: i32 = 0;
    for line in reader.lines() {
        let l = line.unwrap();
        if l.starts_with("addx") {
            cycle += 1;
            result += check_cycle(cycle, x);
            cycle += 1;
            x += l.split_ascii_whitespace().skip(1).next().unwrap().parse::<i32>().unwrap();
            result += check_cycle(cycle, x);
            // print!("at cycle {} x: {} res: {}\n", cycle, x, result);
        } else if l.starts_with("noop") {
            cycle += 1;
            result += check_cycle(cycle, x);
        }
    }
    print!("result score is: {}\n", result);
}

fn check_cycle(c: i32, x: i32) -> i32 {
    //  20th, 60th, 100th, 140th, 180th, and 220th cycles
    return match c {
        20|60|100|140|180|220 => {
            print!("adding {} from {}*{} to the result\n", x*c, x, c);x*c
        } ,
        _ => 0
    }
}

fn draw_pixel(c: i32, x: i32, s: &mut Vec<String>) {
    // screen is 6*40 pixels (h*w)
    // print!("index is {}\n", (c-1)/40);
    let l = s.get_mut(((c-1) / 40) as usize).unwrap();
    let val = (c-1) % 40;
    let ch: char;
    if val == x-1 || val == x || val == x+1 {
        ch = '#';
    } else {
        ch = '.';
    }
    l.push(ch);
}

fn task2(){
    let reader = get_input();
    let mut screen: Vec<String> = Vec::new();
    for _ in 0..HEIGHT {
        screen.push("".to_string());
    }
    print!("screen is {:?}, len {}, last element {}\n", screen, screen.len(), screen.last().unwrap());
    let mut x: i32 = 1;
    let mut cycle: i32 = 1;
    draw_pixel(cycle, x, &mut screen);
    for line in reader.lines() {
        let l = line.unwrap();
        if l.starts_with("addx") {
            cycle += 1;
            draw_pixel(cycle, x, &mut screen);
            cycle += 1;
            x += l.split_ascii_whitespace().skip(1).next().unwrap().parse::<i32>().unwrap();
            draw_pixel(cycle, x, &mut screen);
            // print!("at cycle {} x: {} res: {}\n", cycle, x, result);
        } else if l.starts_with("noop") {
            cycle += 1;
            draw_pixel(cycle, x, &mut screen);
        }
    }
    print!("screen is\n");
    for s in screen.iter() {
        print!("{}\n", s)
    }
}
