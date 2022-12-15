use std::collections::{HashMap};
use std::fs::File;
use std::hash::Hash;
use std::io::{BufRead, BufReader};


fn main() {
    println!("{:-<15}task1{:-<15}", "", "");
    task1();
    println!("{:-<15}task2{:-<15}", "", "");
    task2();
}

enum Tile {
    Rock,
    Sand,
    Spawn,
}

#[derive(Eq, Hash, PartialEq)]
struct Point {
    x: i32,
    y: i32,
}

fn get_input(filename: &str) -> BufReader<File> {
    let file = File::open(filename).unwrap();
    BufReader::new(file)
}

fn new_sparse_map() -> (HashMap<Point, Tile>, i32,i32,i32,i32) {
    let reader = get_input("input");
    let (mut maxy, mut maxx): (i32, i32) = (0, 500);
    let (mut miny, mut minx): (i32, i32) = (0, 500);
    let mut sparse_map: HashMap<Point, Tile>= HashMap::new();
    sparse_map.insert(Point {x: 500, y: 0}, Tile::Spawn);

    for line in reader.lines() {
        let mut prev: Option<(i32, i32)> = None;
        for rock in line.unwrap().split(" -> ") {
            if let None = prev {
                let mut cords = rock.split(",");
                let (x, y) = (cords.next().unwrap().parse::<i32>().unwrap(), cords.next().unwrap().parse::<i32>().unwrap()) ;
                prev = Some((x, y));
                continue
            }
            let mut cords = rock.split(",");
            let (x, y) = (cords.next().unwrap().parse::<i32>().unwrap(), cords.next().unwrap().parse::<i32>().unwrap()) ;
            let (prevx, prevy) = prev.unwrap();
            minx = [minx, x, prevx].iter().min().unwrap().clone();
            miny = [miny, y, prevy].iter().min().unwrap().clone();
            maxx = [maxx, x, prevx].iter().max().unwrap().clone();
            maxy = [maxy, y, prevy].iter().max().unwrap().clone();
            for y in range(prevy, y) {
                for x in range(prevx, x) {
                    sparse_map.insert(Point {x, y}, Tile::Rock);
                }
            }
            prev = Some((x, y));
        }
    }
    (sparse_map, minx, miny ,maxx, maxy)
}

fn task1() {
    let (mut sparse_map, minx,miny,maxx,maxy) = new_sparse_map();

    // simulate sand
    'outer: loop {
        let (mut sy, mut sx) = (0, 500);
        loop {
            if sy > maxy+1 {
                break 'outer
            } else if let None = sparse_map.get(&Point {x: sx, y: sy+1}) {
                (sx, sy) = (sx, sy+1);
            } else if let None = sparse_map.get(&Point {x: sx-1, y: sy+1}) {
                (sx, sy) = (sx-1, sy+1);
            } else if let None = sparse_map.get(&Point {x: sx+1, y: sy+1}) {
                (sx, sy) = (sx+1, sy+1);
            } else {
                sparse_map.insert(Point {x: sx, y: sy}, Tile::Sand);
                break
            }
        }
    }
    print_map(minx, maxx, miny, maxy, &sparse_map);
    let mut cnt = 0;
    for t in sparse_map.values() {
        if let Tile::Sand = *t {
            cnt += 1;
        }
    }
    print!("total grains of sand is {}\n", cnt);
}

fn print_map(minx: i32, maxx: i32, miny: i32, maxy: i32, sparse_map: &HashMap<Point, Tile>) {
    print!("x:    _   ");
    for x in minx-1..=maxx+1 {
        if x %4 == 0 {
            print!("{: <4}", x);
        }
    }
    for y in miny..maxy+3 {
        println!();
        print!("y: {: <2} |", y);
        for x in minx-1..=maxx+1 {
            print!("{}", match sparse_map.get(&Point {x, y}) {
                Some(Tile::Rock) => '#',
                Some(Tile::Sand) => 'o',
                Some(Tile::Spawn) => '+',
                None => '.'
            });
        }
    }
    println!();
}

fn range(a: i32, b: i32) -> impl Iterator<Item = i32> {
    if a > b {
        return b..=a
    }
    a..=b
}

fn task2(){
    let (mut sparse_map, minx,miny,maxx,maxy) = new_sparse_map();

    // simulate sand
    let floory = maxy + 1 + 2;
    'outer: loop {
        let (mut sy, mut sx) = (0, 500);
        loop {
            if sy+2 == floory {
                sparse_map.insert(Point {x: sx, y: sy}, Tile::Sand);
                break
            } else if let None = sparse_map.get(&Point {x: sx, y: sy+1}) {
                (sx, sy) = (sx, sy+1);
            } else if let None = sparse_map.get(&Point {x: sx-1, y: sy+1}) {
                (sx, sy) = (sx-1, sy+1);
            } else if let None = sparse_map.get(&Point {x: sx+1, y: sy+1}) {
                (sx, sy) = (sx+1, sy+1);
            } else if (sy, sx) == (0, 500) {
                break 'outer
            } else {
                sparse_map.insert(Point {x: sx, y: sy}, Tile::Sand);
                break
            }
        }
        // print_map(minx, maxx, miny, maxy, &sparse_map);
    }
    print_map(minx, maxx, miny, maxy, &sparse_map);
    let mut cnt = 0;
    for t in sparse_map.values() {
        if let Tile::Sand = *t {
            cnt += 1;
        }
    }
    print!("total grains of sand is {}\n", cnt+1);
}
