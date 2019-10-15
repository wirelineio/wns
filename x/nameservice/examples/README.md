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
    "id": "QmWmL8D7nT1VDbGFrsroib87EYhTL3zhh8M5z79PGU2aRz",
    "owners": [
      "6ee3328f65c8566cd5451e49e97a767d10a8adf7"
    ],
    "attributes": {
      "name": "wireline.io/chess-bot",
      "protocol": {
        "id": "QmdeazkS38aCrqG6qKwaio2fQnShE6RGpmNdqStLkkZcQN",
        "type": "wrn:reference"
      },
      "type": "wrn:bot",
      "version": "2.0.0"
    },
    "extension": {
      "accessKey": "7db6a0c2b8bc79e733612b5cfd45e9f69bec4d05d424076826a0d08a2a62641c",
      "name": "ChessBot"
    }
  },
  {
    "id": "QmdeazkS38aCrqG6qKwaio2fQnShE6RGpmNdqStLkkZcQN",
    "owners": [
      "6ee3328f65c8566cd5451e49e97a767d10a8adf7"
    ],
    "attributes": {
      "name": "wireline.io/chess",
      "type": "wrn:protocol",
      "version": "1.0.0"
    },
    "extension": {
      "name": "Chess"
    }
  },
  {
    "id": "QmeEeW7WVDcmhDfgjNiRpfJR7tH29WCYRwpcPpVDm7qSYA",
    "owners": [
      "6ee3328f65c8566cd5451e49e97a767d10a8adf7"
    ],
    "attributes": {
      "name": "wireline.io/chess-pad",
      "protocol": {
        "id": "QmdeazkS38aCrqG6qKwaio2fQnShE6RGpmNdqStLkkZcQN",
        "type": "wrn:reference"
      },
      "type": "wrn:pad",
      "version": "5.1.0"
    },
    "extension": {
      "name": "ChessPad"
    }
  }
]
```

If the content in any of the YML files is changed, the file needs to be updated with a new signature. To generate a new signature, pass the `--sign-only` flag to the `set` command. For example:

```
$ wnscli tx nameservice set pad.yml --from root --sign-only
Password to sign with 'root':
Address   : 6ee3328f65c8566cd5451e49e97a767d10a8adf7
PubKey    : 61rphyEC6tEq0pxTI2Sy97VlWCSZhA/PRaUfFlQjhQcpYfTfYtg=
Signature : TRJzxB7J3bwlsYXwvMtiUs2DyW3GBO9BvDQiHKRjkIcXCNblemcrTDgYJ5cFE7fir/ZKp7y8wiKWJ1oPwraIWg==
```
