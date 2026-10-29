import { Block } from "./block";
import { Blockchain } from "./blockchain";

// const block = new Block(
//   0,
//   "Genesis Block",
//   "0",
// );

// console.log(block);
// console.log("Calculated hash:", block.calculateHash());


const blockchain = new Blockchain();

blockchain.addBlock("Alice pays Bob 5000");
blockchain.addBlock("Bob pays Carol 2000");

console.log("Blockchain:");
console.dir(blockchain.chain, { depth: null });

console.log("Is blockchain valid?", blockchain.isValid());

const tamperedBlock = blockchain.chain[1] as { data: string } | undefined;
if (tamperedBlock) {
  tamperedBlock.data = "Alice pays Bob 50000";
}

console.log(
  "Valid after tampering?",
  blockchain.isValid(),
);