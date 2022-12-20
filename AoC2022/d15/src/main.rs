use std::collections::{HashMap, HashSet};
use std::fs::File;
use std::hash::Hash;
use std::io::{BufRead, BufReader};

const TY: i32 = 2_000_000;
const MAX: i32 = 4_000_000;

fn main() {
    println!("{:-<15}task1{:-<15}", "", "");
    // task1();
    println!("{:-<15}task2{:-<15}", "", "");
    task2();
}

#[derive(Eq, Hash, PartialEq, Debug, Clone)]
struct Point {
    x: i32,
    y: i32,
}

#[derive(Debug)]
enum Tile {
    Sensor(Point),
    Beacon,
    NoBeacon,
}

fn get_input(filename: &str) -> BufReader<File> {
    let file = File::open(filename).unwrap();
    BufReader::new(file)
}

fn new_sparse_map() -> HashMap<Point, Tile> {
    let reader = get_input("input");
    let mut sparse_map: HashMap<Point, Tile>= HashMap::new();

    for line in reader.lines() {
        let l = line.unwrap(); //
        let mut split = l.split_ascii_whitespace();
        let sx = split.nth(2).unwrap().trim_end_matches(",")[2..].parse::<i32>().unwrap();
        let sy = split.nth(0).unwrap().trim_end_matches(":")[2..].parse::<i32>().unwrap();
        let bx = split.nth(4).unwrap().trim_end_matches(",")[2..].parse::<i32>().unwrap();
        let by = split.nth(0).unwrap().trim_end_matches(":")[2..].parse::<i32>().unwrap();
        sparse_map.insert(Point {x: sx, y: sy}, Tile::Sensor(
            Point { x: bx, y: by }
        ));
        sparse_map.insert(Point {x: bx,y: by}, Tile::Beacon);
    }
    sparse_map
}

fn task1() {
    let mut sparse_map = new_sparse_map();
    let mut nonpresent: HashSet<i32> = HashSet::new();
    let ty = TY;
    for (k, v) in sparse_map.iter() {
        let d: i32;
        if let Tile::Sensor(p) = &v {
            d = distance(&k, p);
        } else {
            continue
        }
        if !(d > (k.y - ty).abs()) {
            continue
        } 

        let mut x = k.x;
        while distance(&k, &Point{x, y: ty}) <= d {
            // sparse_map.insert(Point{x: x, y: ty}, Tile::NoBeacon);
            // sparse_map.insert(Point{x: 2*k.x-x, y: ty}, Tile::NoBeacon);
            nonpresent.insert(x);
            nonpresent.insert(2*k.x - x);
            x += 1;
        };
    }
    for x in nonpresent.iter() {
        sparse_map.entry(Point{x: *x, y: ty}).or_insert(Tile::NoBeacon);
    }
    // print_map(-5, 25, -5, 22, &sparse_map);
    // print!("y: {ty} |");
    // for x in -6..=26 {
    //     print!("{}", match nonpresent.get(&x) {
    //         Some(_) => '#',
    //         _ => '.',
    //     });
    // }
    // println!();
    // print!("y: {ty} |");
    // for x in -6..=26 {
    //     print!("{}", match sparse_map.get(&Point {x, y: ty}) {
    //         Some(Tile::Sensor(_)) => 'S',
    //         Some(Tile::Beacon) => 'B',
    //         Some(Tile::NoBeacon) => '#',
    //         None => '.'
    //     });
    // }
    // println!();
    print!("total is {}\n", sparse_map.iter().filter(|(k, v)| {k.y == ty && matches!(v, Tile::NoBeacon) }).map(|_| 1).sum::<i32>());
}

// it took few days running on the homelab in the backgroud to find the correct point
// but it did finally
fn task2() {
    let sparse_map = new_sparse_map();
    // let mut nonpresent: HashSet<i32> = HashSet::new();
    let max = MAX;
    let keys: Vec<(Point, i32)> = sparse_map.iter().filter_map(|(p, b)| {
        match b {
            Tile::Sensor(bp) => Some((p.clone(), distance(p, bp))),
            _ => None
        }
    }).collect(); 

    // let check = max/100;
    print!("starting to search through each point for matching one\n");
    // tested up uyntil 350k, not found xD
    'y: for y in (350_000..=max).rev() {
        print!(".");
        if y % 10000 == 0 {
            println!("next 10k done, at {y} \n");
        }
        'x: for x in 0..=max {
            for (k, d) in &keys {
                if (k.x - x).abs() + (k.y - y).abs() <= *d {
                    continue 'x
                } 
            }
            let freq: i128 = x as i128 *4_000_000+y as i128;
            print!("\nSpot at {x}:{y} with tuning freq {} didn't match any sensor range\n", freq);
            break 'y
        }
    }
}

// this even more naive approach would take like infinity instead of a few days
// both are slow
fn task2_slow(){
    let mut sparse_map = new_sparse_map();
    // let mut nonpresent: HashSet<i32> = HashSet::new();
    let max = 20;
    let keys: Vec<_> = sparse_map.keys().cloned().collect(); 
    print!("starting to insert NoBeacon nodes\n");
    for k in keys {
        print!("key {:?} started\n", k);
        let d: i32;
        if let Some(Tile::Sensor(p)) = sparse_map.get(&k) {
            d = distance(&k, p);
        } else {
            continue
        }
        let part = d/100;
        for x in (k.x-d)..=k.x {
            if part > 0 && x % part == 0 {
                print!("1");
            }
            for y in (k.y-d)..=k.y {
                if  (k.x - x).abs() + (k.y - y).abs() <= d {
                    if x >= 0 && x <= max && y >= 0 && y <= max {
                        sparse_map.entry(Point{x: x, y: y}).or_insert(Tile::NoBeacon);
                    }
                    if 2*k.x-x >= 0 && 2*k.x-x <= max && y >= 0 && y <= max {
                        sparse_map.entry(Point{x: 2*k.x-x, y: y}).or_insert(Tile::NoBeacon);
                    }
                    if x >= 0 && x <= max && 2*k.y-y >= 0 && 2*k.y-y <= max {
                        sparse_map.entry(Point{x: x, y: 2*k.y-y}).or_insert(Tile::NoBeacon);
                    }
                    if 2*k.x-x >= 0 && 2*k.x-x <= max && 2*k.y-y >= 0 && 2*k.y-y <= max {
                        sparse_map.entry(Point{x: 2*k.x-x, y: 2*k.y-y}).or_insert(Tile::NoBeacon);
                    }
                } 
            }
            println!();
        }
    } 
    
    // print_map(-5, 25, -5, 22, &sparse_map);
    // print!("y: {ty} |");
    // for x in -6..=26 {
    //     print!("{}", match nonpresent.get(&x) {
    //         Some(_) => '#',
    //         _ => '.',
    //     });
    // }
    // println!();
    // print!("y: {ty} |");
    // for x in -6..=26 {
    //     print!("{}", match sparse_map.get(&Point {x, y: ty}) {
    //         Some(Tile::Sensor(_)) => 'S',
    //         Some(Tile::Beacon) => 'B',
    //         Some(Tile::NoBeacon) => '#',
    //         None => '.'
    //     });
    // }
    // println!();
    print!("starting to seek the single empty node\n");
    let check = max/100;
    'y: for x in 0..=max {
        if x % check == 0 {
            print!("(1%)")
        }
        for y in 0..=max {
            if let None = sparse_map.get(&Point{x, y}) {
                let freq: i128 = x as i128 *4_000_000+y as i128;
                print!("\nFound the spot at {x}:{y} with tuning freq {}\n", freq);
                break 'y
            }
        }
    }
}

fn distance(pa: &Point, pb: &Point) -> i32 {
    (pa.x - pb.x).abs() + (pa.y - pb.y).abs() 
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
                Some(Tile::Sensor(_)) => 'S',
                Some(Tile::Beacon) => 'B',
                Some(Tile::NoBeacon) => '#',
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
