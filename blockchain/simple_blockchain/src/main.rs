// use std::error::Error;
use data_encoding::HEXUPPER;
use serde::{Serialize, Deserialize};
use ring::digest::{Context, Digest, SHA256};
use std::time::{SystemTime};
use std::fmt;

// const BLOCK_DATA_SIZE: usize = 10;
const POW_BITS: usize = 2;

#[derive(Clone, Copy, Deserialize, Serialize, PartialEq, Debug)]
struct Tx {
    from: u64,
    to: u64,
    amount: u128,
}

#[derive(Debug)]
struct Cb {
    to: u64,
    amount: u128,
}

// #[derive(Serialize)]
// type Txs = [Tx; BLOCK_DATA_SIZE];
type BlockSHA = [u8; 32];

// #[derive(Debug)]
// struct BadSHAError {}
// impl Error for BadSHAError {}

#[derive(Debug)]
struct Block {
    id: u64,
    timestamp: u64,
    nonce: u64,
    coinbase: Vec<Cb>,
    tx: Vec<Tx>,
    prev_hash: BlockSHA,
}

struct BlockChain {
    blocks: Vec<Block>
}

impl Block {
    fn new(nonce: u64, ts: Vec<Tx>, prev: &Block, timestamp: u64, cbaddr: u64) -> Self {
        let ph  = <BlockSHA>::try_from(prev.sha256().as_ref()).unwrap();
        let b = Self {
            id: prev.id + 1,
            timestamp: timestamp,
            coinbase: vec![Cb{to: cbaddr, amount: 100}],
            nonce: nonce,
            tx: ts,
            prev_hash: ph,
        };
        b
    }
    fn sha256(&self) -> Digest {
        let mut context = Context::new(&SHA256);
        context.update(&bincode::serialize(&self.id).unwrap()[..]);
        context.update(&bincode::serialize(&self.nonce).unwrap()[..]);
        context.update(&bincode::serialize(&self.timestamp).unwrap()[..]);
        context.update(&bincode::serialize(&self.tx).unwrap()[..]);
        context.update(&bincode::serialize(&self.prev_hash).unwrap()[..]);
        context.finish()
    }
}

impl BlockChain {
    fn new() -> Self {
        BlockChain { blocks: vec![BlockChain::init()] }
    }
    fn head(&self) -> &Block {
        self.blocks.first().unwrap()
    }
    fn tail(&self) -> &Block {
        self.blocks.last().unwrap()
    }
    fn init() -> Block {
        Block {
            id: 0,
            nonce: 0,
            timestamp: SystemTime::now().duration_since(SystemTime::UNIX_EPOCH).unwrap().as_secs(),
            coinbase: vec![Cb{to: 0xf, amount: 100_000_000},],
            tx: vec![Tx{from: 0, to: 0, amount: 0},],
            prev_hash: [0; 32],
        }
    }
    fn mine(&mut self, cbaddr: u64) {
        let tsp = SystemTime::now().duration_since(SystemTime::UNIX_EPOCH).unwrap().as_secs() + 7200;
        let prev = self.tail();
        let mut new_sha: &[u8];
        let mut b = Block::new(0, Vec::<Tx>::new(), prev, tsp, cbaddr);
        let mut d: Digest;
        'out: for i in 0..std::u64::MAX {
            b.nonce = i;
            d = b.sha256();
            new_sha = d.as_ref();
            if new_sha[0..POW_BITS] != [0u8; POW_BITS as usize] {
                continue 'out
            }
            // for j in 0..POW_BITS {
            //     if new_sha[j as usize] != 0u8 {
            //         continue 'out
            //     }
            // }
            println!("mined a block, hash: {:?}, nonce:{} ",  new_sha, b.nonce);
            self.blocks.push(b);
            return
        }
    }

    fn verify(&self) -> bool {
        let mut prev = self.blocks.first().unwrap().sha256().clone();
        let mut cur: Digest;
        for b in self.blocks.iter() {
            let d = b.sha256();
            cur = d;//b.sha256().as_ref();
            if cur.as_ref() != prev.as_ref() {
                print!("fount not matching hashes: {:?} prev: {:?}\n", cur.as_ref(), prev.as_ref());
                return false
            }
            prev = cur;
        }
        true
    }
}

impl fmt::Display for BlockChain {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "Len of blockchain is: {}\n", self.blocks.len());
        for b in self.blocks.iter() {
            write!(f, "{:?}\n", b);
        }
        write!(f, "--------\n")
    }
}

fn main() {
    println!("Hello, world!");
    let my_address = 1u64;
    let b = BlockChain::init();
    let h = b.sha256();
    let ph  = <BlockSHA>::try_from(h.as_ref()).unwrap();
    print!(" h is {:?}\nph is {:?}\n", h.as_ref(), ph);
    print!("len of hash slice is {}\n", h.as_ref().len());
    print!("SHA256 for \n block {:?} \nis\n{}\n", b, HEXUPPER.encode(h.as_ref()));
    let mut c = BlockChain::new();
    c.mine(my_address);
    c.mine(my_address);
    c.mine(my_address);
    print!("c is: {}\n", c);
    print!("is the chain correct: {}", c.verify());
}