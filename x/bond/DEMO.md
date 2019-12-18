# Bonds Demo

## Setup

### Machine

Setup the machine as documented in https://github.com/wirelineio/wns#setup-machine.

### Code


```bash
# Clone `wns` repo. 
$ git clone git@github.com:wirelineio/wns.git
$ cd wns

# Switch to the `feature-bonds` branch.
$ git checkout feature-bonds

# Build and install the binaries.
$ make install
```

### Blockchain

```bash
# Delete old folders.
$ rm -rf ~/.wnsd ~/.wnscli

# Init the chain.
$ wnsd init my-node --chain-id wireline
```

```bash
# Update genesis params.
# Note: On Linux, use just `-i` instead of `-i ''`.

# Change staking token to uwire.
$ sed -i '' 's/stake/uwire/g' ~/.wnsd/config/genesis.json

# Change gov proposal pass timeout to 5 mins.
$ sed -i '' 's/172800000000000/300000000000/g' ~/.wnsd/config/genesis.json

# Change max bond amount.
$ sed -i '' 's/10wire/10000wire/g' ~/.wnsd/config/genesis.json
```

```bash
# Create accounts/keys.
$ echo "temp12345\nsalad portion potato insect unknown exile lion soft layer evolve flavor hollow emerge celery ankle sponsor easy effort flush furnace life maximum rotate apple" | wnscli keys add root --recover
$ echo "temp12345\nquestion cause van artefact belt dish turkey badge twenty bronze breeze visa" | wnscli keys add alice --recover
$ echo "temp12345\nhospital speak toward arrange tide universe attend surround useless nerve true nasty" | wnscli keys add bob --recover
```

```bash
# Add genesis accounts to chain.
$ wnsd add-genesis-account $(wnscli keys show root -a) 100000000000000uwire
$ wnsd add-genesis-account $(wnscli keys show alice -a) 100000000000000uwire
$ wnsd add-genesis-account $(wnscli keys show bob -a) 100000000000000uwire
```

```bash
# CLI config.
$ wnscli config chain-id wireline
$ wnscli config output json
$ wnscli config indent true
$ wnscli config trust-node true
```

```bash
# Setup genesis transactions.
$ echo "temp12345" | wnsd gentx --name root --amount 10000000000000uwire
$ wnsd collect-gentxs
$ wnsd validate-genesis
```

```bash
# Start the chain.
$ wnsd start
```

## Run

Create bonds.

```bash
# Two bonds from the `root` account.
$ echo temp12345 | wnscli tx bond create 1000wire --from root --yes -b block
$ echo temp12345 | wnscli tx bond create 1000wire --from root --yes -b block

# One bond each from the `alice` and `bob` accounts.
$ echo temp12345 | wnscli tx bond create 1000wire --from alice --yes -b block
$ echo temp12345 | wnscli tx bond create 1000wire --from bob --yes -b block
```

List bonds.

```bash
# Note: Bond ID, owner and balance.
$ wnscli query bond list
[
  {
    "id": "319fdfc0f4198643f2eb8adf602fb2e7f08682cb4d123d44d1935c87b554b959",
    "owner": "cosmos1razr52gj62vvgqtneqmys8hklm02mynv0exky4",
    "balance": [
      {
        "denom": "uwire",
        "amount": "1000000000"
      }
    ]
  },
  {
    "id": "614b4affedc58705ba7eb8fac1d0fbcc9cffaeeaf787bd042b27a7447b37177e",
    "owner": "cosmos1zk8etz23phxgtse8re6tggsr3nrfk2vtsesegy",
    "balance": [
      {
        "denom": "uwire",
        "amount": "1000000000"
      }
    ]
  },
  {
    "id": "8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3",
    "owner": "cosmos1wh8vvd0ymc5nt37h29z8kk2g2ays45ct2qu094",
    "balance": [
      {
        "denom": "uwire",
        "amount": "1000000000"
      }
    ]
  },
  {
    "id": "e205a46f6ec6f662cbfad84f4f926973422bf6217d8d2c2eebff94d148fd486d",
    "owner": "cosmos1wh8vvd0ymc5nt37h29z8kk2g2ays45ct2qu094",
    "balance": [
      {
        "denom": "uwire",
        "amount": "1000000000"
      }
    ]
  }
]
```

List balance across all bonds.

```bash
$ wnscli query bond balance
{
  "bond": [
    {
      "denom": "uwire",
      "amount": "4198000000"
    }
  ]
}
```

Get bond by ID.

```bash
$ wnscli query bond get 319fdfc0f4198643f2eb8adf602fb2e7f08682cb4d123d44d1935c87b554b959
{
  "id": "319fdfc0f4198643f2eb8adf602fb2e7f08682cb4d123d44d1935c87b554b959",
  "owner": "cosmos1razr52gj62vvgqtneqmys8hklm02mynv0exky4",
  "balance": [
    {
      "denom": "uwire",
      "amount": "1000000000"
    }
  ]
}
```

Query bonds by owner.

```bash
# Uses a secondary index: Owner -> Bond ID.
$ wnscli query bond query-by-owner $(wnscli keys show -a root)
[
  {
    "id": "8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3",
    "owner": "cosmos1wh8vvd0ymc5nt37h29z8kk2g2ays45ct2qu094",
    "balance": [
      {
        "denom": "uwire",
        "amount": "1000000000"
      }
    ]
  },
  {
    "id": "e205a46f6ec6f662cbfad84f4f926973422bf6217d8d2c2eebff94d148fd486d",
    "owner": "cosmos1wh8vvd0ymc5nt37h29z8kk2g2ays45ct2qu094",
    "balance": [
      {
        "denom": "uwire",
        "amount": "1000000000"
      }
    ]
  }
]
```

Refill bond.

```bash
$ echo temp12345 | wnscli tx bond refill 8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3 500wire --from root --yes -b block

$ wnscli query bond get 8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3
{
  "id": "8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3",
  "owner": "cosmos1wh8vvd0ymc5nt37h29z8kk2g2ays45ct2qu094",
  "balance": [
    {
      "denom": "uwire",
      "amount": "1500000000"
    }
  ]
}
```

Withdraw funds from bond.

```bash
# Transfers the funds back into the bond owner account.
$ echo temp12345 | wnscli tx bond withdraw 8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3 300wire --from root --yes -b block

$ wnscli query bond get 8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3
{
  "id": "8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3",
  "owner": "cosmos1wh8vvd0ymc5nt37h29z8kk2g2ays45ct2qu094",
  "balance": [
    {
      "denom": "uwire",
      "amount": "1200000000"
    }
  ]
}
```

Publish a record (w/ bond).

```bash
$ cd x/nameservice/examples

$ echo temp12345 | wnscli tx nameservice set protocol.yml 8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3 --from root --yes -b block

$ echo temp12345 | wnscli tx nameservice set bot.yml 8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3 --from root --yes -b block

# Note: bondID and expiryTime attributes on the records.
$ wnscli query nameservice list
[
  {
    "id": "QmNgCCwB2AGQADe1X4P1kVjd1asdTandaRHbRp1fKrDH9i",
    "bondId": "8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3",
    "expiryTime": "2020-12-17T11:00:54.288525Z",
    "owners": [
      "6ee3328f65c8566cd5451e49e97a767d10a8adf7"
    ],
    "attributes": {
      "accessKey": "7db6a0c2b8bc79e733612b5cfd45e9f69bec4d05d424076826a0d08a2a62641c",
      "displayName": "ChessBot",
      "name": "wireline.io/chess-bot",
      "protocol": {
        "id": "Qmb8beMVtdjQpRiN9bFLUsisByun8EHPKPo2g1jSUUnHqe",
        "type": "wrn:reference"
      },
      "type": "wrn:bot",
      "version": "2.0.0"
    }
  },
  {
    "id": "Qmb8beMVtdjQpRiN9bFLUsisByun8EHPKPo2g1jSUUnHqe",
    "bondId": "8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3",
    "expiryTime": "2020-12-17T11:00:44.173405Z",
    "owners": [
      "6ee3328f65c8566cd5451e49e97a767d10a8adf7"
    ],
    "attributes": {
      "displayName": "Chess",
      "name": "wireline.io/chess",
      "type": "wrn:protocol",
      "version": "1.0.0"
    }
  }
]

# Note: Rent has been deducted from the bond.
$ wnscli query bond get 8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3
{
  "id": "8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3",
  "owner": "cosmos1wh8vvd0ymc5nt37h29z8kk2g2ays45ct2qu094",
  "balance": [
    {
      "denom": "uwire",
      "amount": "1198000000"
    }
  ]
}

# Note: Check balance of bond and record rent module accounts.
$ wnscli query bond balance
{
  "bond": [
    {
      "denom": "uwire",
      "amount": "4198000000"
    }
  ],
  "record_rent": [
    {
      "denom": "uwire",
      "amount": "2000000"
    }
  ]
}
```

Query records by bond.

```bash
# Note: Uses a secondary index: Bond ID -> Record ID. 
$ wnscli query nameservice query-by-bond 8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3
[
  {
    "id": "QmNgCCwB2AGQADe1X4P1kVjd1asdTandaRHbRp1fKrDH9i",
    "bondId": "8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3",
    "expiryTime": "2020-12-17T11:00:54.288525Z",
    "owners": [
      "6ee3328f65c8566cd5451e49e97a767d10a8adf7"
    ],
    "attributes": {
      "accessKey": "7db6a0c2b8bc79e733612b5cfd45e9f69bec4d05d424076826a0d08a2a62641c",
      "displayName": "ChessBot",
      "name": "wireline.io/chess-bot",
      "protocol": {
        "id": "Qmb8beMVtdjQpRiN9bFLUsisByun8EHPKPo2g1jSUUnHqe",
        "type": "wrn:reference"
      },
      "type": "wrn:bot",
      "version": "2.0.0"
    }
  },
  {
    "id": "Qmb8beMVtdjQpRiN9bFLUsisByun8EHPKPo2g1jSUUnHqe",
    "bondId": "8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3",
    "expiryTime": "2020-12-17T11:00:44.173405Z",
    "owners": [
      "6ee3328f65c8566cd5451e49e97a767d10a8adf7"
    ],
    "attributes": {
      "displayName": "Chess",
      "name": "wireline.io/chess",
      "type": "wrn:protocol",
      "version": "1.0.0"
    }
  }
]
```

Dissociate bond from record.

```bash
$ echo temp12345 | wnscli tx nameservice dissociate-bond QmNgCCwB2AGQADe1X4P1kVjd1asdTandaRHbRp1fKrDH9i --from root --yes -b block

$ wnscli query nameservice query-by-bond 8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3
[
  {
    "id": "Qmb8beMVtdjQpRiN9bFLUsisByun8EHPKPo2g1jSUUnHqe",
    "bondId": "8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3",
    "expiryTime": "2020-12-17T11:00:44.173405Z",
    "owners": [
      "6ee3328f65c8566cd5451e49e97a767d10a8adf7"
    ],
    "attributes": {
      "displayName": "Chess",
      "name": "wireline.io/chess",
      "type": "wrn:protocol",
      "version": "1.0.0"
    }
  }
]
```

Associate bond with record.

```bash
$ echo temp12345 | wnscli tx nameservice associate-bond QmNgCCwB2AGQADe1X4P1kVjd1asdTandaRHbRp1fKrDH9i 8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3 --from root --yes -b block

$ wnscli query nameservice query-by-bond 8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3
[
  {
    "id": "QmNgCCwB2AGQADe1X4P1kVjd1asdTandaRHbRp1fKrDH9i",
    "bondId": "8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3",
    "expiryTime": "2020-12-17T11:00:54.288525Z",
    "owners": [
      "6ee3328f65c8566cd5451e49e97a767d10a8adf7"
    ],
    "attributes": {
      "accessKey": "7db6a0c2b8bc79e733612b5cfd45e9f69bec4d05d424076826a0d08a2a62641c",
      "displayName": "ChessBot",
      "name": "wireline.io/chess-bot",
      "protocol": {
        "id": "Qmb8beMVtdjQpRiN9bFLUsisByun8EHPKPo2g1jSUUnHqe",
        "type": "wrn:reference"
      },
      "type": "wrn:bot",
      "version": "2.0.0"
    }
  },
  {
    "id": "Qmb8beMVtdjQpRiN9bFLUsisByun8EHPKPo2g1jSUUnHqe",
    "bondId": "8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3",
    "expiryTime": "2020-12-17T11:00:44.173405Z",
    "owners": [
      "6ee3328f65c8566cd5451e49e97a767d10a8adf7"
    ],
    "attributes": {
      "displayName": "Chess",
      "name": "wireline.io/chess",
      "type": "wrn:protocol",
      "version": "1.0.0"
    }
  }
]
```

Dissociate bond from all records.

```bash
$ echo temp12345 | wnscli tx nameservice dissociate-records 8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3 --from root --yes -b block

# Note: No records found.
$ wnscli query nameservice query-by-bond 8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3
```

Reassociate bond.

```bash
# First, associate records with a bond.
$ echo temp12345 | wnscli tx nameservice associate-bond QmNgCCwB2AGQADe1X4P1kVjd1asdTandaRHbRp1fKrDH9i 8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3 --from root --yes -b block

# First, associate records with a bond.
$ echo temp12345 | wnscli tx nameservice associate-bond Qmb8beMVtdjQpRiN9bFLUsisByun8EHPKPo2g1jSUUnHqe 8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3 --from root --yes -b block

# Check both records as associated with bond 8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3.
$ wnscli query nameservice list
[
  {
    "id": "QmNgCCwB2AGQADe1X4P1kVjd1asdTandaRHbRp1fKrDH9i",
    "bondId": "8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3",
    "expiryTime": "2020-12-17T11:00:54.288525Z",
    "owners": [
      "6ee3328f65c8566cd5451e49e97a767d10a8adf7"
    ],
    "attributes": {
      "accessKey": "7db6a0c2b8bc79e733612b5cfd45e9f69bec4d05d424076826a0d08a2a62641c",
      "displayName": "ChessBot",
      "name": "wireline.io/chess-bot",
      "protocol": {
        "id": "Qmb8beMVtdjQpRiN9bFLUsisByun8EHPKPo2g1jSUUnHqe",
        "type": "wrn:reference"
      },
      "type": "wrn:bot",
      "version": "2.0.0"
    }
  },
  {
    "id": "Qmb8beMVtdjQpRiN9bFLUsisByun8EHPKPo2g1jSUUnHqe",
    "bondId": "8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3",
    "expiryTime": "2020-12-17T11:00:44.173405Z",
    "owners": [
      "6ee3328f65c8566cd5451e49e97a767d10a8adf7"
    ],
    "attributes": {
      "displayName": "Chess",
      "name": "wireline.io/chess",
      "type": "wrn:protocol",
      "version": "1.0.0"
    }
  }
]

# List of bonds.
$ wnscli query bond query-by-owner $(wnscli keys show -a root)
[
  {
    "id": "8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3",
    "owner": "cosmos1wh8vvd0ymc5nt37h29z8kk2g2ays45ct2qu094",
    "balance": [
      {
        "denom": "uwire",
        "amount": "1198000000"
      }
    ]
  },
  {
    "id": "e205a46f6ec6f662cbfad84f4f926973422bf6217d8d2c2eebff94d148fd486d",
    "owner": "cosmos1wh8vvd0ymc5nt37h29z8kk2g2ays45ct2qu094",
    "balance": [
      {
        "denom": "uwire",
        "amount": "1000000000"
      }
    ]
  }
]

# Switch to bond e205a46f6ec6f662cbfad84f4f926973422bf6217d8d2c2eebff94d148fd486d.
$ echo temp12345 | wnscli tx nameservice reassociate-records 8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3 e205a46f6ec6f662cbfad84f4f926973422bf6217d8d2c2eebff94d148fd486d --from root --yes -b block

# Note: Records are now associated with bond e205a46f6ec6f662cbfad84f4f926973422bf6217d8d2c2eebff94d148fd486d.
$ wnscli query nameservice list
[
  {
    "id": "QmNgCCwB2AGQADe1X4P1kVjd1asdTandaRHbRp1fKrDH9i",
    "bondId": "e205a46f6ec6f662cbfad84f4f926973422bf6217d8d2c2eebff94d148fd486d",
    "expiryTime": "2020-12-17T11:00:54.288525Z",
    "owners": [
      "6ee3328f65c8566cd5451e49e97a767d10a8adf7"
    ],
    "attributes": {
      "accessKey": "7db6a0c2b8bc79e733612b5cfd45e9f69bec4d05d424076826a0d08a2a62641c",
      "displayName": "ChessBot",
      "name": "wireline.io/chess-bot",
      "protocol": {
        "id": "Qmb8beMVtdjQpRiN9bFLUsisByun8EHPKPo2g1jSUUnHqe",
        "type": "wrn:reference"
      },
      "type": "wrn:bot",
      "version": "2.0.0"
    }
  },
  {
    "id": "Qmb8beMVtdjQpRiN9bFLUsisByun8EHPKPo2g1jSUUnHqe",
    "bondId": "e205a46f6ec6f662cbfad84f4f926973422bf6217d8d2c2eebff94d148fd486d",
    "expiryTime": "2020-12-17T11:00:44.173405Z",
    "owners": [
      "6ee3328f65c8566cd5451e49e97a767d10a8adf7"
    ],
    "attributes": {
      "displayName": "Chess",
      "name": "wireline.io/chess",
      "type": "wrn:protocol",
      "version": "1.0.0"
    }
  }
]
```

Cancel bond.

```bash
$ wnscli query bond query-by-owner $(wnscli keys show -a root)
[
  {
    "id": "8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3",
    "owner": "cosmos1wh8vvd0ymc5nt37h29z8kk2g2ays45ct2qu094",
    "balance": [
      {
        "denom": "uwire",
        "amount": "1198000000"
      }
    ]
  },
  {
    "id": "e205a46f6ec6f662cbfad84f4f926973422bf6217d8d2c2eebff94d148fd486d",
    "owner": "cosmos1wh8vvd0ymc5nt37h29z8kk2g2ays45ct2qu094",
    "balance": [
      {
        "denom": "uwire",
        "amount": "1000000000"
      }
    ]
  }
]

# Note: Cancel fails if there are associated records.
$ echo temp12345 | wnscli tx bond cancel e205a46f6ec6f662cbfad84f4f926973422bf6217d8d2c2eebff94d148fd486d --from root --yes -b block

# Note: Cancel works if bond doesn't have associated records.
$ echo temp12345 | wnscli tx bond cancel 8e340dd7cf6fc91c27eeefce9cca1406c262e93fd6f3a4f3b1e99b01161fcef3 --from root --yes -b block

# Note: Cancelled bond is deleted.
$ wnscli query bond query-by-owner $(wnscli keys show -a root)
[
  {
    "id": "e205a46f6ec6f662cbfad84f4f926973422bf6217d8d2c2eebff94d148fd486d",
    "owner": "cosmos1wh8vvd0ymc5nt37h29z8kk2g2ays45ct2qu094",
    "balance": [
      {
        "denom": "uwire",
        "amount": "1000000000"
      }
    ]
  }
]
```

Consensus params.

```bash
$ wnscli query bond params
{
  "max_bond_amount": "10000wire"
}

$ wnscli query nameservice params
{
  "record_rent": "1wire",
  "record_expiry_time": "31536000000000000"
}
```

Create proposal to change param.

```bash
$ echo temp12345 | wnscli tx gov submit-proposal param-change params/update_max_bond_amount.json --from root --yes -b block
```

Query proposals.

```bash
# Note: Proposal status = DepositPeriod.
$ wnscli query gov proposals
[
  {
    "content": {
      "type": "cosmos-sdk/ParameterChangeProposal",
      "value": {
        "title": "Bond Param Change",
        "description": "Update max bond amount.",
        "changes": [
          {
            "subspace": "bond",
            "key": "MaxBondAmount",
            "value": "\"15000wire\""
          }
        ]
      }
    },
    "id": "1",
    "proposal_status": "DepositPeriod",
    "final_tally_result": {
      "yes": "0",
      "abstain": "0",
      "no": "0",
      "no_with_veto": "0"
    },
    "submit_time": "2019-12-18T11:41:20.691337Z",
    "deposit_end_time": "2019-12-18T11:46:20.691337Z",
    "total_deposit": [
      {
        "denom": "uwire",
        "amount": "1000000"
      }
    ],
    "voting_start_time": "0001-01-01T00:00:00Z",
    "voting_end_time": "0001-01-01T00:00:00Z"
  }
]
```

Deposit sufficient funds to move the proposal into the voting stage.

```bash
# Note: Proposal ID = 1.
$ echo temp12345 | wnscli tx gov deposit 1 10000000uwire --from root -yes -b block

# Note: Proposal status = VotingPeriod.
$ wnscli query gov proposals
[
  {
    "content": {
      "type": "cosmos-sdk/ParameterChangeProposal",
      "value": {
        "title": "Bond Param Change",
        "description": "Update max bond amount.",
        "changes": [
          {
            "subspace": "bond",
            "key": "MaxBondAmount",
            "value": "\"15000wire\""
          }
        ]
      }
    },
    "id": "1",
    "proposal_status": "VotingPeriod",
    "final_tally_result": {
      "yes": "0",
      "abstain": "0",
      "no": "0",
      "no_with_veto": "0"
    },
    "submit_time": "2019-12-18T11:41:20.691337Z",
    "deposit_end_time": "2019-12-18T11:46:20.691337Z",
    "total_deposit": [
      {
        "denom": "uwire",
        "amount": "11000000"
      }
    ],
    "voting_start_time": "2019-12-18T11:42:46.695024Z",
    "voting_end_time": "2019-12-18T11:47:46.695024Z"
  }
]
```

Vote (yes) on proposal.

```bash
$ echo temp12345 | wnscli tx gov vote 1 yes --from root --yes -b block
```

Check votes.

```bash
$ wnscli query gov votes 1
[
  {
    "proposal_id": "1",
    "voter": "cosmos1wh8vvd0ymc5nt37h29z8kk2g2ays45ct2qu094",
    "option": "Yes"
  }
]

$ wnscli query gov tally 1
{
  "yes": "10000000000000",
  "abstain": "0",
  "no": "0",
  "no_with_veto": "0"
}
```

Wait for 5 mins. Proposal enters `Passed` status.

```bash
# Note: Proposal status = Passed.
$ wnscli query gov proposals
[
  {
    "content": {
      "type": "cosmos-sdk/ParameterChangeProposal",
      "value": {
        "title": "Bond Param Change",
        "description": "Update max bond amount.",
        "changes": [
          {
            "subspace": "bond",
            "key": "MaxBondAmount",
            "value": "\"15000wire\""
          }
        ]
      }
    },
    "id": "1",
    "proposal_status": "Passed",
    "final_tally_result": {
      "yes": "10000000000000",
      "abstain": "0",
      "no": "0",
      "no_with_veto": "0"
    },
    "submit_time": "2019-12-18T11:41:20.691337Z",
    "deposit_end_time": "2019-12-18T11:46:20.691337Z",
    "total_deposit": [
      {
        "denom": "uwire",
        "amount": "11000000"
      }
    ],
    "voting_start_time": "2019-12-18T11:42:46.695024Z",
    "voting_end_time": "2019-12-18T11:47:46.695024Z"
  }
]
```

Check updated value of param.

```bash
$ wnscli query bond params
{
  "max_bond_amount": "15000wire"
}
```
