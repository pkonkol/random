use std::collections::VecDeque;
use std::fs::File;
use std::io::{BufRead, BufReader};


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
    let mut map = Vec::<Vec<char>>::new();
    let (mut starty, mut startx): (usize, usize) = (0, 0);
    let (mut endy, mut endx): (usize, usize) = (0, 0);
    for (y, line) in reader.lines().enumerate() {
        map.push(Vec::<char>::new());
        let l = map.last_mut().unwrap();
        for (x, ch) in line.unwrap().chars().enumerate() {
            if ch == 'S' {
                (starty, startx) = (y, x);
            } else if ch == 'E' {
                (endy, endx) = (y, x);
                l.push('z');
                continue
            }
            l.push(ch);
        }
    }
    let len = bfs(starty, startx, endy, endx, &map);
    print!("path length is {}\n", len);
}

fn task2(){
    let reader = get_input();
    let mut map = Vec::<Vec<char>>::new();
    // let (mut starty, mut startx): (usize, usize) = (0, 0);
    let (mut endy, mut endx): (usize, usize) = (0, 0);
    for (y, line) in reader.lines().enumerate() {
        map.push(Vec::<char>::new());
        let l = map.last_mut().unwrap();
        for (x, ch) in line.unwrap().chars().enumerate() {
            if ch == 'S' {
                l.push('a');
                // (starty, startx) = (y, x);
                continue
            } else if ch == 'E' {
                (endy, endx) = (y, x);
                l.push('z');
                continue
            }
            l.push(ch);
        }
    }
    let mut lengths = Vec::<usize>::new();
    for (y, line) in map.iter().enumerate() {
        for (x, ch)in line.iter().enumerate() {
            if *ch == 'a' {
                lengths.push(bfs(y, x, endy, endx, &map));
                // print!("`");
            }
        }
    }
    // print!("lenghts: {:?}", lengths);
    print!("\nshortest path length is {}\n", lengths.iter().min().unwrap());
}


fn bfs(starty: usize, startx: usize, endy: usize, endx: usize, map: &Vec<Vec<char>>) -> usize {
    let mut out = 1usize << 63;
    let mut visited = vec![vec![false; map[0].len()]; map.len()];
    let mut queue = VecDeque::<(usize, usize, usize)>::new();
    queue.push_front((starty, startx, 0));
    let dirs: [(i32, i32); 4]= [(1, 0), (-1, 0), (0, 1), (0, -1)];
    // let mut prev_depth = 0;
    while !queue.is_empty() {
        // print_visited(&visited);
        let (ny, nx, depth) = queue.pop_back().unwrap();
        if visited[ny][nx] {
            continue
        }
        visited[ny][nx] = true;
        if ny == endy && nx == endx {
            out = depth;
            break
        }
        // if prev_depth < depth {
        //     print!("{},", depth);
        //     prev_depth = depth;
        // }
        for (dy, dx) in dirs {
            let (sy, sx) = (ny as i32 + dy, nx as i32+dx);
            if !(sy >= 0 && sy < map.len() as i32 && sx >= 0 && sx < map[0].len() as i32) {
                continue
            }
            let (sy, sx) = ((ny as i32+dy) as usize, (nx as i32+dx) as usize);
            if !visited[sy][sx] && get_val(map[ny][nx])+1 >= get_val(map[sy][sx]) {
                queue.push_front((sy, sx, depth+1));
            }
        }
    }
    // print_visited(&visited);
    out
}

fn print_visited(v: &Vec<Vec<bool>>) {
    for line in v.iter() {
        for b in line {
            print!("{}", if *b == true {"#"} else {"."} )
        }
        println!();
    }
}

// its a fucking DFS but might be useful later so let's leave it
fn dfs(y: usize, x: usize, endy: usize, endx: usize, dist: usize, map: &Vec<Vec<char>>, visited: &mut Vec<Vec<bool>>) -> usize {
    if y == endy && x == endx {
        return dist;
    }
    visited[y][x] = true;
    let dirs: [(i32, i32); 4]= [(1, 0), (-1, 0), (0, 1), (0, -1)];
    for (dy, dx) in dirs {
        let (newy, newx) = (y as i32 + dy, x as i32+dx);
        if !(newy >= 0 && newy < map.len() as i32 && newx >= 0 && newx < map[0].len() as i32) {
            continue
        }
        let (newy, newx) = ((y as i32+dy) as usize, (x as i32+dx) as usize);
        if get_val(map[y][x])+1 >= get_val(map[newy][newx]) && !visited[newy][newx] {
            let out = dfs(newy, newx, endy, endx, dist+1, map, visited);
            if out != 0 {
                return out
            }
        }
    }   
    0usize
}

fn get_val(c: char) -> i32 {
    if c.is_ascii_lowercase() {
        return (c as u8 - 'a' as u8) as i32
    }
    -1
}

#[cfg(test)]
mod tests {
    use crate::*;
    #[test]
    fn get_val_works() {
        assert_eq!(0, get_val('a'));
        assert_eq!(25, get_val('z'));
    }
}