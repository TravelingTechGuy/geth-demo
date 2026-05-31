#!/bin/bash
set -e

# Navigate to the project root relative to the script location
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

echo "🔍 Fetching developer account address..."
DEV_ADDR=$(docker exec geth-local geth --exec "eth.coinbase" attach http://127.0.0.1:8545 | tr -d '"\r')

if [ -z "$DEV_ADDR" ]; then
    echo "❌ Error: Could not find developer account. Is Geth running?"
    exit 1
fi
echo "✅ Developer account: $DEV_ADDR"

create_account() {
    # Pipes two newlines to handle the password and password confirmation prompts
    printf "\n\n" | docker exec -i geth-local geth account new --datadir /root/.ethereum | grep -oE '0x[a-fA-F0-9]{40}'
}

echo "🆕 Creating test accounts..."
ADDR1=$(create_account)
echo "   Account 1: $ADDR1"
ADDR2=$(create_account)
echo "   Account 2: $ADDR2"

echo "🔓 Decrypting private keys..."
# We run the go utility from the root so it can find the ./data directory
PK1=$(go run cmd/decrypt_keystore/main.go -address "$ADDR1" -password "" | awk '/Private key/ {print $NF}')
PK2=$(go run cmd/decrypt_keystore/main.go -address "$ADDR2" -password "" | awk '/Private key/ {print $NF}')

echo "💰 Funding accounts with 100 ETH each..."
docker exec geth-local geth --exec "eth.sendTransaction({from: '$DEV_ADDR', to: '$ADDR1', value: web3.toWei(100, 'ether')})" attach http://127.0.0.1:8545 > /dev/null
docker exec geth-local geth --exec "eth.sendTransaction({from: '$DEV_ADDR', to: '$ADDR2', value: web3.toWei(100, 'ether')})" attach http://127.0.0.1:8545 > /dev/null

echo "📝 Generating .env file..."
cat <<EOF > .env
GETH_RPC_URL=http://localhost:8545
PRIVATE_KEY_1=$PK1
PRIVATE_KEY_2=$PK2
EOF

echo "✨ Done!"
echo "Project configured with:"
echo "  - Private Key 1: $PK1"
echo "  - Private Key 2: $PK2"
echo "You can now run 'go run main.go' or 'go test -v'."
