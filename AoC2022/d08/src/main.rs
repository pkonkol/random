use std::collections::{HashSet};
use std::fs::File;
use std::io::{BufRead, BufReader};

fn main() {
    println!("{:-<15}task1{:-<15}", "", "");
    task1();
    println!("{:-<15}task2{:-<15}", "", "");
    task2();
}

fn get_forest() -> Vec<Vec<u8>> {
    let filename = "input";
    let file = File::open(filename).unwrap();
    let reader = BufReader::new(file);

    let mut forest: Vec<Vec<u8>> = Vec::new();

    for line in reader.lines() {
        let l = line.unwrap();
        let mut row: Vec<u8> = Vec::new();
        for c in l.chars() {
            row.push(c.to_digit(10).unwrap().try_into().unwrap());
        }
        forest.push(row);
    }
    forest
}

fn task1() {
    let forest = get_forest();
    let mut visible: HashSet<(usize, usize)> = HashSet::new();
    let width = forest.get(0).unwrap().len();
    let height = forest.len();

    // iterate from left
    for y in 0..height {
        let mut prev_top = -1;
        let vec_y = forest.get(y).unwrap();
        for x in 0..width {
            let cur = vec_y.get(x).unwrap().clone() as i32;
            if cur > prev_top {
                visible.insert((y, x));
                prev_top = cur;
            }
        }
    }
    // iterate from right
    for y in 0..height {
        let mut prev_top = -1;
        let vec_y = forest.get(y).unwrap();
        for x in (0..width).rev() {
            let cur = vec_y.get(x).unwrap().clone() as i32;
            if cur > prev_top {
                visible.insert((y, x));
                prev_top = cur;
            }
        }
    }
    // iterate from top
    for x in 0..width {
        let mut prev_top = -1;
        for y in 0..height {
            let vec_y = forest.get(y).unwrap();
            let cur = vec_y.get(x).unwrap().clone() as i32;
            if cur > prev_top {
                visible.insert((y, x));
                prev_top = cur;
            }
        }
    }
    // iterate from bottom
    for x in 0..width {
        let mut prev_top = -1;
        for y in (0..height).rev() {
            let vec_y = forest.get(y).unwrap();
            let cur = vec_y.get(x).unwrap().clone() as i32;
            if cur > prev_top {
                visible.insert((y, x));
                prev_top = cur;
            }
        }
    }
    print!("visible trees {}\n", visible.len());
}

fn visile_from(start_y: usize, start_x: usize, forest: &Vec<Vec<u8>>) -> i32 {
    let width = forest.get(0).unwrap().len();
    let height = forest.len();
    let start_top = forest.get(start_y).unwrap().get(start_x).unwrap().clone() as i32;
    
    if start_x == 0 || start_y == 0 || start_x == width-1 || start_y == height-1 {
        print!("skipping at border x:{}, y:{}\n", start_x, start_y);
        return 0
    }

    // print!("startx:{}; ", start_x);
    // for x in (0..start_x).rev() {
    //     print!("{};", x);
    // }
    // println!();

    let vec_y = forest.get(start_y).unwrap();
    // iterate to the left
    let mut score_l = 1;
    for x in (0..start_x).rev() {
        let cur = vec_y.get(x).unwrap().clone() as i32;
        if cur >= start_top {
            score_l = start_x - x;
            break;
        }
        if x == 0 && cur <= start_top {
            score_l = start_x;
        }
    }
    // iterate to the right
    let mut score_r = 1;
    for x in (start_x..width).skip(1) {
        let cur = vec_y.get(x).unwrap().clone() as i32;
        if cur >= start_top {
            score_r = x - start_x;
            break;
        }
        if x == width - 1 && cur < start_top {
            score_r = width - start_x - 1;
        }
    }
    
    // iterate up
    let mut score_u = 1;
    for y in (0..start_y).rev() {
        let vec_y = forest.get(y).unwrap();
        let cur = vec_y.get(start_x).unwrap().clone() as i32;
        if cur >= start_top {
            score_u = start_y - y;
            break;
        }
        if y == 0 && cur < start_top {
            score_u = start_y;
        }
    }
    
    // iterate down
    let mut score_d = 1;
    for y in (start_y..height).skip(1) {
        let vec_y = forest.get(y).unwrap();
        let cur = vec_y.get(start_x).unwrap().clone() as i32;
        if cur >= start_top {
            score_d = y - start_y;
            break;
        }
        if y == height - 1 && cur < start_top {
            score_d = height - start_y - 1;
        }
    }
    print!("x:{},y:{},h:{},w:{},l:{},r:{},u:{},d:{}\n", start_x, start_y, height, width, score_l, score_r, score_u, score_d);
    (score_l * score_r * score_u * score_d) as i32
}

fn task2(){
    let forest = get_forest();
    let width = forest.get(0).unwrap().len();
    let height = forest.len();

    let mut max = -1;
    let mut max_cords = (-1, -1);
    for y in 0..height {
        for x in 0..width {
            let res = visile_from(y, x, &forest);
            if res > max {
                max = res;
                max_cords = (y as i32, x as i32);
            }
        }
    }
    print!("best visibility is {} from y:{} x:{}\n", max, max_cords.0, max_cords.1);
}
