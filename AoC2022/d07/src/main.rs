use std::cell::RefCell;
use std::collections::{VecDeque};
use std::fs::File;
use std::io::{BufReader, BufRead};
use std::rc::{Rc};

const ELVEN_SIZE_LIMIT: i128 = 100_000; //i128 not necessary but it's cool
const ELVEN_MEM_TO_FREE: i128 = 30_000_000; // haven't used it before
const ELVEN_TOTAL_MEM: i128 = 70_000_000;

fn main() {
    let root = gen_dir_tree();
    println!("{:-<15}task1{:-<15}", "", "");
    task1(root.clone());
    println!("{:-<15}task2{:-<15}", "", "");
    task2(root);
}

#[derive(Debug)]
struct DirTree {
    pub name: String,
    pub parent: Option<Rc<RefCell<DirTree>>>,
    pub children: Vec<Rc<RefCell<DirTree>>>,
    pub files: Vec<(String, i32)>,
    pub size: Option<i32>,
}

impl DirTree {
    pub fn new(p: Option<Rc<RefCell<DirTree>>>, name: String) -> Self {
        DirTree {
            name,
            parent: p,
            children: Vec::new(),
            files: Vec::new(),
            size: None,
        }
    }

}

fn get_dirs_dfs(root: Rc<RefCell<DirTree>>) -> VecDeque<Rc<RefCell<DirTree>>> {
    let mut deq = VecDeque::new();
    let mut out = VecDeque::new();

    deq.push_back(root.clone());
    out.push_back(root.clone());
    while !deq.is_empty() { // DFS to generate the size for dirs
        let n = deq.pop_front().unwrap();
        print!("at {} with files {} and kids {}\n", 
            n.borrow().name,
            n.borrow().files.iter().map(|(a,b)| { format!("{}-{}", a, b) }).collect::<Vec<_>>().join(";"),
            n.borrow().children.iter().map(|n| { n.borrow().name.clone() }).collect::<Vec<_>>().join(";"),
        );
        for c in n.borrow().children.iter() {
            deq.push_front(c.clone());
            out.push_front(c.clone());
        }
    }
    out
}

fn gen_dir_tree() -> Rc<RefCell<DirTree>> {
    let filename = "input";
    let file = File::open(filename).unwrap();
    let reader = BufReader::new(file);

    let root = DirTree::new(None, "/".to_string());
    let mut cur = Rc::new(RefCell::new(root));
    let root = cur.clone();

    for line in reader.lines().skip(1) {
        let mut line = line.unwrap();
        if line.starts_with('$') {
            line.drain(..2);
            if line.starts_with("ls") {
                print!("found ls command \n");
            } else if line.starts_with("cd") {
                line.drain(..3);
                print!("found cd command to {}\n", line);
                if line == ".." {
                    let next = cur.borrow().parent.clone().unwrap();
                    cur = next;
                } else {
                    let not_known = cur
                        .borrow()
                        .children
                        .iter()
                        .find(|x| x.borrow().name == line)
                        .is_none();


                    if not_known {
                        let new = Rc::new(RefCell::new(DirTree::new(Some(cur.clone()), line)));
                        cur.borrow_mut().children.push(new);
                        let next = cur.borrow().children.last().unwrap().clone();
                        cur = next;
                    } else {
                        let next = cur
                            .borrow()
                            .children
                            .iter()
                            .map(|x| x.clone())
                            .find(|x| x.borrow().name == line)
                            .unwrap();
                        cur = next;
                    }
                }
            }
        } else if line.chars().nth(0).unwrap().is_numeric() {
            let offset = line.find(' ').unwrap_or(line.len());
            let size: i32 = line.drain(..offset).collect::<String>().parse().unwrap();
            let not_known = cur
                .borrow()
                .files
                .iter()
                .find(|&x| {x.0 == line})
                .is_none();
            if not_known {
                cur.borrow_mut().files.push((line.clone(), size));
            }
            print!("found a file {} size {}\n", line.trim(), size);
        } else if line.starts_with("dir") {
            line.drain(..4);
            print!("found a dir {}\n", line);
            let not_known = cur
                .borrow()
                .children
                .iter()
                .find(|&x| x.borrow().name == line)
                .is_none();
                // .unwrap_or(&mut DirTree::new(Some(cur), line));
            if not_known {
                cur.borrow_mut().children.push(Rc::new(RefCell::new(DirTree::new(Some(cur.clone()), line))));
            }
        }
    }
    root
}

fn task1(root: Rc<RefCell<DirTree>>) {
    let mut size_deq = get_dirs_dfs(root);
    
    // size_deq has the leafs at the front, going backwards one level in depth
    print!("summing sizes {:-<15}\n", "-");
    let mut total: i128 = 0;
    while !size_deq.is_empty() {
        let n = size_deq.pop_front().unwrap();
        let mut size: i32 = 0;
        size += n.borrow().files.iter().map(|(_,b)| b).sum::<i32>();
        print!("dir={} size of files: {}  ", n.borrow().name,  size);
        size += n.borrow().children.iter().map(|x| x.clone()).map(|x| x.borrow().size.unwrap()).sum::<i32>();
        print!("size after adding kids: {}  ", size);
        n.borrow_mut().size = Some(size);
        if (size as i128) < ELVEN_SIZE_LIMIT {
            print!("matched limid with {} for {}", total, n.borrow().name);
            total += size as i128;
        }
        println!("");
    }
    println!("total of dirs smaller than {} is {}", ELVEN_SIZE_LIMIT, total);
}

fn task2(root: Rc<RefCell<DirTree>>){
    let size_deq = get_dirs_dfs(root.clone());
    
    let needed_space = ELVEN_MEM_TO_FREE - (ELVEN_TOTAL_MEM - root.borrow().size.unwrap() as i128);
    print!("seeking smallest dir to remove {:-<15}\n", "-");
    let min = size_deq
        .iter().cloned()
        .filter(|x| (x.borrow().size.unwrap() as i128) > needed_space)
        .map(|x| x.borrow().size.unwrap())
        .min()
        .unwrap();

    println!("size of root is {}", root.borrow().size.unwrap());
    println!("smallest dir larger than {} is {}", needed_space, min);
}
