import dotenv from 'dotenv';
dotenv.config();


import {ethers} from 'ethers';

console.log()

console.log("Node version:", process.version);
console.log("Ethers version:", ethers.version ? ethers.version : "unknown");
console.log("RPC URL loaded:", !!process.env.RPC_URL);