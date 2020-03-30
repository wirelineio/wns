# Ops

## Setup

* Install [Ansible](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html).
* Updating code from GitHub requires SSH agent [forwarding](https://developer.github.com/v3/guides/using-ssh-agent-forwarding/).
* Add the local SSH public key to the `~/.ssh/authorized_keys` file on the remote machines (`ubuntu` account).

## Operations

Hard reset development testnet (Warning: Wipes all data on the testnet):

```bash
$ cd ~/wirelineio/wns/ops
$ ansible-playbook -i env/development reset-full.yml
```

Post full reset, create a bond.

```bash
$ ansible-playbook -i env/development ./create-bond.yml
```
