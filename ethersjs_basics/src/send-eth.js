import dotenv from 'dotenv';
import { ethers } from 'ethers';

dotenv.config();

async function main() {
    const rpcURL = process.env.RPC_URL;
    const privateKey = process.env.PRIVATE_KEY;
    const balanceWallet = '0x0f15cFEFbB8fa32aEeB368D73ae5AA42b57D5853';

    if (!rpcURL || !privateKey) {
        throw new Error("Set RPC_URL and PRIVATE_KEY in .env");
    }

    try {
        // Add connection options with timeout and retry
        const provider = new ethers.JsonRpcProvider(rpcURL, null, {
            staticNetwork: true,
            batchMaxCount: 1
        });

        console.log("Provider created, testing connection...");
        
        // Test the connection first
        const network = await provider.getNetwork();
        console.log(`Connected to network: ${network.name} (chainId: ${network.chainId})`);

        // Get balance
        const balance = await provider.getBalance(balanceWallet);
        console.log(`Balance for wallet: ${ethers.formatEther(balance)} ETH`);

    } catch (error) {
        console.error("Error details:", {
            message: error.message,
            code: error.code,
            operation: error.operation
        });
        throw error;
    }
}

main().catch(console.error);