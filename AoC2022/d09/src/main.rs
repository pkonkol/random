use std::collections::{HashSet};
use std::fs::File;
use std::io::{BufRead, BufReader};

fn main() {
    println!("{:-<15}task1{:-<15}", "", "");
    task1();
    println!("{:-<15}task2{:-<15}", "", "");
    task2();
}

type Point = (i32, i32);

fn get_input() -> BufReader<File> {
    let filename = "input";
    let file = File::open(filename).unwrap();
    BufReader::new(file)
}

fn rope_move(dir: char, h: Point, t: Point) -> (Point, Point) {
    let new_h = match dir {
        'R' => (h.0, h.1+1),
        'L' => (h.0, h.1-1),
        'U' => (h.0+1, h.1),
        'D' => (h.0-1, h.1),
        _ => (0, 0)
    };
    let dx: i32 = new_h.1 - t.1;
    let dy: i32 = new_h.0 - t.0;
    let new_t: Point;
    new_t = match (dy, dx) {
        (y, x) if y <=1 && x <= 1 && y >= -1 && x >= -1 => (t.0, t.1),
        (y, x) => (t.0+y.clamp(-1, 1), t.1+x.clamp(-1, 1))
    };
    (new_h, new_t)
}

fn snake_move(h: Point, t: Point) -> Point {
    let dx: i32 = h.1 - t.1;
    let dy: i32 = h.0 - t.0;
    let new_t: Point;
    new_t = match (dy, dx) {
        (y, x) if y <=1 && x <= 1 && y >= -1 && x >= -1 => (t.0, t.1),
        (y, x) => (t.0+y.clamp(-1, 1), t.1+x.clamp(-1, 1))
    };
    new_t
}


fn task1() {
    let mut visited: HashSet<Point> = HashSet::new();
    let mut pos_h: Point = (0, 0);
    let mut pos_t: Point = (0, 0);

    let reader = get_input();
    for line in reader.lines() {
        let mut l = line.unwrap();
        let mov = l.chars().nth(0).unwrap();
        let times = l.split_off(2).parse::<u32>().unwrap();
        // print!("move is {} times {}\n", mov, times);

        for _ in 0..times {
            (pos_h, pos_t) = rope_move(mov, pos_h, pos_t);
            visited.insert(pos_t);
        }
    }
    print!("visited len: {}\n", visited.len());
}

fn task2(){
    let mut visited: HashSet<Point> = HashSet::new();
    //let mut pos_h: Point = (0, 0);
    let mut snake: Vec<Point> = vec![(0,0); 10]; // pos_h is at 0

    let reader = get_input();
    for line in reader.lines() {
        let mut l = line.unwrap();
        let mov = l.chars().nth(0).unwrap();
        let times = l.split_off(2).parse::<u32>().unwrap();
        // print!("move is {} times {}\n", mov, times);

        for _ in 0..times {

            let (phead, ptail) = (snake.get(0).unwrap().clone(), snake.get(1).unwrap().clone());
            let (new_h, new_t) = rope_move(mov, phead, ptail);
            *snake.get_mut(0).unwrap() = (new_h.0, new_h.1);
            *snake.get_mut(1).unwrap() = (new_t.0, new_t.1);

            for i in 2..snake.len() {
                let (pfront, pback) = (snake.get(i-1).unwrap().clone(), snake.get(i).unwrap().clone());
                // print!("i: {} pfront: {:?} pback:{:?}\n", i, pfront, pback);
                let new_t = snake_move(pfront, pback);
                *snake.get_mut(i).unwrap() = (new_t.0, new_t.1);
            }
            // print!("tail pos is:{:?}, {:-<15}\n", snake.last().unwrap(), "-");
            visited.insert(snake.last().unwrap().clone());
        }
    }
    print!("visited len: {}\n", visited.len());
}
