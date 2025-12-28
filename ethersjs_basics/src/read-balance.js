import dotenv from 'dotenv';
dotenv.config();

import { ethers } from 'ethers';

async function main() {
  const rpcURL = process.env.RPC_URL;
  if (!rpcURL) throw new Error("RPC_URL not set in .env");

  const provider = new ethers.JsonRpcProvider(rpcURL);
  const address = process.env.ADDRESS;

  const balance = await provider.getBalance(address);
  console.log(`Balance of ${address}: ${ethers.formatEther(balance)} ETH`);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
