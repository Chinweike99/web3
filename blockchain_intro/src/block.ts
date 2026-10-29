import {createHash} from "node:crypto";

export class Block {
 public readonly index: number;
 public readonly timestamp: number;
 public readonly data: string;
 public readonly previousHash: string;
 public hash: string;
 
 
 constructor (
    index: number,
    data: string,
    previousHash: string,
 ){
    this.index = index;
    this.timestamp = Date.now();
    this.data = data;
    this.previousHash = previousHash;
    this.hash = this.calculateHash();
 }

 public calculateHash(): string {
    const blockContents = [
        this.index,
        this.timestamp,
        this.data,
        this.previousHash
    ].join("|");
    return createHash("sha256").update(blockContents).digest("hex");
 }


}