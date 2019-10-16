# Examples

The folder that contains this README.md file has example YML files that can be used to create records in WNS.

```
$ wnscli tx nameservice set protocol.yml --from root
$ wnscli tx nameservice set pad.yml --from root
$ wnscli tx nameservice set bot.yml --from root

```

```
$ wnscli query nameservice list
[
  {
    "id": "QmNgCCwB2AGQADe1X4P1kVjd1asdTandaRHbRp1fKrDH9i",
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
    "id": "QmStVv79TJRoG9emnZn3ZXKaoYd2tfJuFvQ3NBQkBnAB62",
    "owners": [
      "6ee3328f65c8566cd5451e49e97a767d10a8adf7"
    ],
    "attributes": {
      "displayName": "ChessPad",
      "name": "wireline.io/chess-pad",
      "protocol": {
        "id": "Qmb8beMVtdjQpRiN9bFLUsisByun8EHPKPo2g1jSUUnHqe",
        "type": "wrn:reference"
      },
      "type": "wrn:pad",
      "version": "5.1.0"
    }
  },
  {
    "id": "Qmb8beMVtdjQpRiN9bFLUsisByun8EHPKPo2g1jSUUnHqe",
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

If the content in any of the YML files is changed, the file needs to be updated with a new signature. To generate a new signature, pass the `--sign-only` flag to the `set` command. For example:

```
$ wnscli tx nameservice set pad.yml --from root --sign-only
Password to sign with 'root':
CID       : QmStVv79TJRoG9emnZn3ZXKaoYd2tfJuFvQ3NBQkBnAB62
Address   : 6ee3328f65c8566cd5451e49e97a767d10a8adf7
PubKey    : 61rphyEC6tEq0pxTI2Sy97VlWCSZhA/PRaUfFlQjhQcpYfTfYtg=
Signature : LbuUEZhp88Ukvj3XWWUG0B9sk2tOrDB9jUvbbEkJBjAihsWzwtTN3W9IAJ9fsVQI2URlrJn0YUvTvVM28cPxag==
SigData   : {"displayName":"ChessPad","name":"wireline.io/chess-pad","protocol":{"id":"Qmb8beMVtdjQpRiN9bFLUsisByun8EHPKPo2g1jSUUnHqe","type":"wrn:reference"},"type":"wrn:pad","version":"5.1.0"}
```
