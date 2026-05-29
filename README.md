# Testing Dockerized EVM

This app test a docker Geth installation. It uses `dev` as the blockchain, and persists state between runs in a `/data` folder.

## Steps to run

### Run the blockchain

Run this in your terminal:
```bash
docker-compose up -d
```
Verify that a `/data` folder is created

Run:
```bash
docker logs geth-local | tail -n 50
```
That should give you the last 50 lines from Geth log. Look for a line that reads something like
```bash
Using developer account   address=0x71562b71999873DB5b286dF957af199Ec94617F7
```
This is your developer account
### Create 2 test accounts

Create 2 new accounts. Leave the passwords empty (if you want to specify passwords, note them somewhere):
```bash
docker exec -it geth-local geth account new --datadir /root/.ethereum
docker exec -it geth-local geth account new --datadir /root/.ethereum
```
### Recover a private key from a local Geth keystore

After creating local accounts with `geth account new`, you can recover the private key for an account you own by decrypting its local keystore file.

Run:

```bash
go run ./cmd/decrypt_keystore/main.go -address 0xYOUR_ADDRESS -password 'your-password'
```

By default, the utility reads from `./data/geth/keystore`.

If the keystore password is an empty string, pass an empty shell argument:

```bash
go run decrypt_keystore.go -address 0xYOUR_ADDRESS -password ''
```

The utility prints the private key in hex without the `0x` prefix, so it can be copied directly into `.env` file. So create am `.env` file at the root of the project with the following contents:

```env
GETH_RPC_URL=http://localhost:8545
PRIVATE_KEY_1=abcdef...
PRIVATE_KEY_2=123456...
```
Feel free to create more addresses, and add them to the `.env` file - just make sure you load them in `config.go`.

### Fund the 2 accounts we created

Get into the geth CLI:
```bash
docker exec -it geth-local geth attach http://127.0.0.1:8545
```
Check to see which accounts are available (one should be the dev account, and 2 are the new ones you just created):
```bash
> eth.accounts
["0x<DEV account?0>", "0x<Account 1>", "0x<Account2>"]
```
We'll verify the dev account has funds:
```
> eth.getBalance("0x<dev account>")
1.15792089237316195423570985008687907853269984665640564039257583973451624232927e+77
```
That's a lot of ETH (see the e+77 part?). Now let's send 100 ETH to each one of our accounts:
```bash
> eth.sendTransaction({from:"0x<Dev account>", to: "0x<Account1>", value: web3.toWei(100, "ether")})
"0x06f353a7af6a5a9e523e149aafcdacab4208e5ec779f2b04203da545f3c8e6ea"
> eth.sendTransaction({from:"0x<Dev account>", to: "0x<Account2>", value: web3.toWei(100, "ether")})
"0x06f353a7af6a5a9e523e149aafcdacab4208e5ec779f2b04203da545f3c8e6ea"
> eth.getBalance("0x<Account1>")
100000000000000000000
> eth.getBalance("0x<Account2>")
100000000000000000000
```

Exit the CLI by clickcing CTRL+D, and now we're ready to run our app!

### Running the application

1. `go run main.go` - run main app
1. `go test -v` - run tests

### Shutting down the blockchain

To shut down the blockchain, run `docker-compose down`. This will still preserve the `/data` directory, so we can reuse the same accounts next time.

To fully remove the blockchain, and reset the accounts:
```bash
docker-compose down -v
rm -rf /data
```

And next time, you'll have to start from scratch

## Notes

You can always verify the chain is up by using the shell scripts provided:
```bash
cd ./test_scripts
chmod + x *.sh
./chain_id.sh   #returns the chain ID
./gas_price.sh  #returns gas price
cd ..
```
