import { Block } from "./block";

export class Blockchain {
    public readonly chain: Block[];

    constructor() {
        this.chain = [this.createGenesisBlock()];
    }

    private createGenesisBlock(): Block {
        return new Block(0, "Genesis Block", "0");
    }

    public getLatestBlock(): Block {
        const latestBlock = this.chain[this.chain.length - 1];
        if (!latestBlock) {
            throw new Error("Blockchain is empty");
        }
        return latestBlock;
    }

    public addBlock(data: string): void {
        const previousBlock = this.getLatestBlock();
        const newBlock = new Block(
            this.chain.length,
            data,
            previousBlock.hash,
        );
        this.chain.push(newBlock)
    }

    public isValid(): boolean {
        for (let i = 0; i < this.chain.length; i++) {
            const currentBlock = this.chain[i];

            if (!currentBlock) {
                return false;
            }

            if (currentBlock.hash !== currentBlock.calculateHash()) {
                return false;
            }

            if (i === 0) {
                continue;
            }

            const previousBlock = this.chain[i - 1];
            if (!previousBlock || currentBlock.previousHash !== previousBlock.hash) {
                return false;
            }
        }
        return true;
    }
}