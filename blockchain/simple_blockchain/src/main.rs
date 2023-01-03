// use std::error::Error;
use data_encoding::HEXUPPER;
use serde::{Serialize, Deserialize};
use ring::digest::{Context, Digest, SHA256};
use std::{time::{SystemTime}, fmt::Display, str::FromStr};
use owo_colors::OwoColorize;
use std::fmt;
use k256::{
    ecdsa::{SigningKey, Signature, signature::Signer, VerifyingKey, signature::Verifier},
    EncodedPoint,
    SecretKey,
    elliptic_curve::sec1::ToEncodedPoint, //schnorr::signature::Signature,
    // ecdsa::{},
};
use generic_array::{ArrayLength, GenericArray};

const POW_BITS: usize = 2;
const MY_PUBKEY: u64 = 1;
const MASTER_SEED: [u8; 8] = [0; 8];
const PENDING_THRESHOLD: usize = 3;

#[derive(Clone,  Deserialize, Serialize, PartialEq, Debug)] // Copy,
struct Tx {
    from: Pubkey,
    to: Pubkey,
    amount: u128,
    signature: Sig,
}

#[derive(Debug)]
struct Cb {
    to: Pubkey,
    amount: u128,
}

// type Txs = [Tx; BLOCK_DATA_SIZE];
// #[derive(Serialize)]
type BlockSHA = [u8; 32];

// #[derive(Serialize)]
type Pubkey = Vec<u8>;
type Sig = String;

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
    blocks: Vec<Block>,
    pending: Vec<Tx>,
}

struct Balance {
    address: Pubkey,
    amount: u64,
    tx_count: u64,
}

type StateCache = Vec::<Balance>;

impl Tx {
    fn new(privkey: SigningKey, to: Pubkey, amount: u128) -> Self {
        let pubkey = VerifyingKey::from(&privkey); // Serialize with `::to_encoded_point()`

        let mut data = bincode::serialize(pubkey.to_encoded_point(false).as_bytes()).unwrap();
        data.append(&mut bincode::serialize(&Vec::from(to.clone())).unwrap());
        data.append(&mut bincode::serialize(&amount).unwrap());
        let signature: Signature = privkey.sign(&data);

        Tx {
            from: Vec::from(pubkey.to_encoded_point(false).as_bytes()),
            to: to,
            amount: amount,
            signature: signature.to_string(),
        }
    }
    
    fn verify(&self) -> bool {
        let pubkey = VerifyingKey::from_sec1_bytes(&self.from).unwrap(); // Serialize with `::to_encoded_point()`

        let mut data = bincode::serialize(&self.from).unwrap();
        data.append(&mut bincode::serialize(&self.to).unwrap());
        data.append(&mut bincode::serialize(&self.amount).unwrap());

        let t = Signature::from_str(&self.signature).unwrap();
        pubkey.verify(&data, &t).is_ok()
    }
}

impl Block {
    fn new(nonce: u64, ts: Vec<Tx>, prev: &Block, timestamp: u64, cbaddr: &Pubkey) -> Self {
        let ph  = <BlockSHA>::try_from(prev.sha256().as_ref()).unwrap();
        let b = Self {
            id: prev.id + 1,
            timestamp: timestamp,
            coinbase: vec![Cb{to: cbaddr.clone(), amount: 100}],
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

impl Display for Block {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{}: {}, timestamp: {}, tx len: {}, prev: {}",
            "Block ".green(), self.id, self.timestamp, self.tx.len(), HEXUPPER.encode(self.prev_hash.as_ref()))
    }
}

impl BlockChain {
    fn new(cbaddr: Pubkey) -> Self {
        BlockChain { 
            blocks: vec![BlockChain::init(cbaddr)],
            pending: Vec::new(),
        }
    }
    fn head(&self) -> &Block {
        self.blocks.first().unwrap()
    }
    fn tail(&self) -> &Block {
        self.blocks.last().unwrap()
    }
    fn init(cbaddr: Pubkey) -> Block {
        Block {
            id: 0,
            nonce: 0,
            timestamp: SystemTime::now().duration_since(SystemTime::UNIX_EPOCH).unwrap().as_secs(),
            coinbase: vec![Cb{to: cbaddr, amount: 100_000_000},],
            tx: Vec::new(),
            prev_hash: [0; 32],
        }
    }
    fn mine(&mut self, cbaddr: &Pubkey) {
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
            println!("{} hash: {}, nonce:{} ", "Mined a block".red(), HEXUPPER.encode(new_sha), b.nonce);
            self.blocks.push(b);
            return
        }
    }

    fn verify(&self) -> bool {
        let mut prev = <BlockSHA>::try_from(self.blocks.first().unwrap().sha256().clone().as_ref()).unwrap();
        let mut cur;
        for b in self.blocks.iter().skip(1) {
            cur = <BlockSHA>::try_from(b.sha256().as_ref()).unwrap();
            if b.prev_hash != prev {
                print!("found not matching hashes: {:?} prev: {:?}\n",
                    HEXUPPER.encode(cur.as_ref()).red(), HEXUPPER.encode(prev.as_ref()).yellow());
                return false
            }
            prev = cur;
        }
        print!("{}", "Blockchain is valid :)\n".green());
        true
    }

    fn add_tx(&mut self, tx: Tx) {
        if !tx.verify() {
            print!("{}: {:?}", "TX failed to verify\n".on_red().white(), tx);
            return
        }
        self.pending.push(tx);
    }
}

impl fmt::Display for BlockChain {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{}: len: {}\n", "BlockChain".yellow(), self.blocks.len());
        for b in self.blocks.iter() {
            write!(f, "\t{}\n", b);
        }
        write!(f, "")
    }
}

fn main() {
    println!("{}", "start".on_green());
    let master_privkey = SigningKey::from_bytes(&MASTER_SEED).unwrap();
    let master_pubkey = VerifyingKey::from(&master_privkey);
    let master_pubkey_vec = Vec::from(master_pubkey.to_encoded_point(false).as_bytes());
    let mut c = BlockChain::new(master_pubkey_vec.clone());
    c.mine(&master_pubkey_vec);
    c.mine(&master_pubkey_vec);
    c.mine(&master_pubkey_vec);
    print!("{}", c);
    c.verify();
    println!();


    let signing_key = SigningKey::from_bytes(&MASTER_SEED).unwrap();
    let message = b"ECDSA proves knowledge of a secret number in the context of a single message";
    let signature: Signature = signing_key.sign(message);
    let verifying_key = VerifyingKey::from(&signing_key); // Serialize with `::to_encoded_point()`
    print!("test sig is: {}", signature.to_string());

    let privkey_1 = SigningKey::from_bytes(&[1; 8]).unwrap();
    let pubkey_1 = VerifyingKey::from(&master_privkey);
    let new_tx = Tx::new(master_privkey, Vec::from(pubkey_1.to_encoded_point(false).as_bytes()), 10);
    c.add_tx(new_tx);


    // TODO signatures and pubkey based accounts & transactions
    // TODO later p2p nodes mesh
}