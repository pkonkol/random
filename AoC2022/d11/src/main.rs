use std::cell::RefCell;

fn main() {
    println!("{:-<15}task1{:-<15}", "", "");
    task1();
    println!("{:-<15}task2{:-<15}", "", "");
    task2();
}

const ROUNDS: i128 = 20;
const ROUNDS2: i128 = 10000;

struct Monkey {
    items: Vec<i128>, //stack
    operation: fn(i128) -> i128, //func
    test: fn(i128) -> i128,  //func
    cnt: i128,
}

impl Monkey {
    // fn investigate(&mut self, monkeys: &mut Vec<Monkey>) {
    //     for mut i in self.items.into_iter() {
    //         i = i/3;
    //         i = (self.operation)(i);
    //         self.cnt += 1;
    //         monkeys.get_mut((self.test)(i) as usize).unwrap().items.push(i);
    //     }
    // }

    fn print(&self) {
        print!("items are: {:?}\n", self.items);
    }
}

fn task1() {
    let monkeys = vec![
        RefCell::new(Monkey {items: vec![57,], operation: |x| x*13, test: |x| if x%11==0 {3} else {2}, cnt: 0}),
        RefCell::new(Monkey {items: vec![58, 93, 88, 81, 72, 73, 65], operation: |x| x+2, test: |x| if x%7==0 {6} else {7}, cnt: 0}),
        RefCell::new(Monkey {items: vec![65, 95], operation: |x| x+6, test: |x| if x%13==0 {3} else {5}, cnt: 0}),
        RefCell::new(Monkey {items: vec![58, 80, 81, 83], operation: |x| x*x, test: |x| if x%5==0 {4} else {5}, cnt: 0}),
        RefCell::new(Monkey {items: vec![58, 89, 90, 96, 55], operation: |x| x+3, test: |x| if x%3==0 {1} else {7}, cnt: 0}),
        RefCell::new(Monkey {items: vec![66, 73, 87, 58, 62, 67], operation: |x| x*7, test: |x| if x%17==0 {4} else {1}, cnt: 0}),
        RefCell::new(Monkey {items: vec![85, 55, 89], operation: |x| x+4, test: |x| if x%2==0 {2} else {0}, cnt: 0}),
        RefCell::new(Monkey {items: vec![73, 80, 54, 94, 90, 52, 69, 58], operation: |x| x+7, test: |x| if x%19==0 {6} else {0}, cnt: 0}),
    ];
    for i in 0..monkeys.len() {
        print!("monkey {} ", i);
        monkeys.get(i).unwrap().borrow().print();
    }
    for round in 1..=ROUNDS {
        print!(" {} {:-<30}\n", round ,'-');
        for i in 0..monkeys.len() {
            let m = monkeys.get(i).unwrap();
            let items = m.borrow().items.iter().cloned().collect::<Vec<_>>();

            print!("monkey {} ---- ", i);
            monkeys.get(i).unwrap().borrow().print();

            m.borrow_mut().items.clear();
            for mut i in items {
                print!("{} ", i);
                i = (m.borrow().operation)(i);
                print!("{} ", i);
                i = i/3;
                print!("{}>", i);
                m.borrow_mut().cnt += 1;
                monkeys.get((m.borrow().test)(i) as usize).unwrap().borrow_mut().items.push(i);
                print!("{}", (m.borrow().test)(i));
                print!("| ")
            }
            println!();
        }
        for i in 0..monkeys.len() {
            print!("monkey {} ", i);
            monkeys.get(i).unwrap().borrow().print();
        }
    }
    for (i, m) in monkeys.iter().enumerate() {
        print!("{:-<30}\nmonkey {} investigated {} times\n", "-", i, m.borrow().cnt);
    }
    // let max1 = monkeys.iter().enumerate().max_by(|(i, x)| x.cnt).unwrap();
    // let max2 = monkeys.iter().max_by(|x| x.cnt);
    //print!("monkey business is is: {}\n", max1 * max2 );
}

fn task2(){
    let monkeys = vec![
        RefCell::new(Monkey {items: vec![57,], operation: |x| x*13, test: |x| if x%11==0 {3} else {2}, cnt: 0}),
        RefCell::new(Monkey {items: vec![58, 93, 88, 81, 72, 73, 65], operation: |x| x+2, test: |x| if x%7==0 {6} else {7}, cnt: 0}),
        RefCell::new(Monkey {items: vec![65, 95], operation: |x| x+6, test: |x| if x%13==0 {3} else {5}, cnt: 0}),
        RefCell::new(Monkey {items: vec![58, 80, 81, 83], operation: |x| x*x, test: |x| if x%5==0 {4} else {5}, cnt: 0}),
        RefCell::new(Monkey {items: vec![58, 89, 90, 96, 55], operation: |x| x+3, test: |x| if x%3==0 {1} else {7}, cnt: 0}),
        RefCell::new(Monkey {items: vec![66, 73, 87, 58, 62, 67], operation: |x| x*7, test: |x| if x%17==0 {4} else {1}, cnt: 0}),
        RefCell::new(Monkey {items: vec![85, 55, 89], operation: |x| x+4, test: |x| if x%2==0 {2} else {0}, cnt: 0}),
        RefCell::new(Monkey {items: vec![73, 80, 54, 94, 90, 52, 69, 58], operation: |x| x+7, test: |x| if x%19==0 {6} else {0}, cnt: 0}),
    ];
    // LCM of test divisors is 9,699,690
    let lcm = 9_699_690;
    for _ in 1..=ROUNDS2 {
        for i in 0..monkeys.len() {
            let m = monkeys.get(i).unwrap();
            let items = m.borrow().items.iter().cloned().collect::<Vec<_>>();
            m.borrow_mut().items.clear();
            for mut i in items {
                i = (m.borrow().operation)(i) % lcm;
                m.borrow_mut().cnt += 1;
                monkeys.get((m.borrow().test)(i) as usize).unwrap().borrow_mut().items.push(i);
            }
        }
    }
    for (i, m) in monkeys.iter().enumerate() {
        print!("{:-<30}\nmonkey {} investigated {} times\n", "-", i, m.borrow().cnt);
    }
}
