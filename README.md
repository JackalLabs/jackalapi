# Jackal API Server

The Jackal API Server is an HTTP server designed to act as a centralized upload point for the Jackal Protocol. This can
be run as-is locally or integrated into existing tech-stacks through a series of HTTP requests.

You will need to supply the API with a seed phrase that corresponds to a Jackal account with $JKL funding and a storage
plan active. You can activate that storage plan by heading to the dashboard with a wallet connected sharing the
seed-phrase of the API.

## Installation

```shell
git clone https://github.com/JackalLabs/jackalapi.git
cd jackalapi
git checkout v0.2.0
go make install
```

## Usage

All variables are set by environment variables, this includes the Jackal RPC, the wallet seed-phrase, chain-id and the
port to run on. You will be required to enter a seed phrase.

### Env Variables

Jackal Network (Defaults to Testnet)

* JAPI_SEED - (none)
* JAPI_RPC - `https://jackal-testnet-rpc.polkachu.com:443` (possibly rate limited)
* JAPI_CHAIN - `lupulella-2`

Other Networks

* JAPI_IPFS_GATEWAY - `https://ipfs.io/ipfs/`

Root Directories

* JAPI_OP_ROOT - `JAPI`
* JAPI_IPFS_ROOT - `s/JAPI/IPFS`
* JAPI_BULK_ROOT - `s/JAPI/Bulk`

Misc Variables

* JAPI_PORT - `3535`
* JAPI_MAX_FILE - (none)

### Uploading File

```shell
curl -F "file=@FILENAME" http://localhost:3535/upload
```

### Checking IPFS File

For every time the API is hit with a CID request, it will first check the Jackal network for the file, if the file does
not exist on the Jackal network, it will download the file, upload it to the Jackal network and then forward the file.
If it does have the file, it will download the file from the Jackal network and forward it to you.

In any browser or CLI, you can visit http://localhost:3535/ipfs/{CID/PATH}.
