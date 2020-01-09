# Ops

## Setup

* Install [Ansible](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html).
* Updating code from GitHub requires SSH agent [forwarding](https://developer.github.com/v3/guides/using-ssh-agent-forwarding/).

## Operations

Hard reset development testnet (Warning: Wipes all data on the testnet):

```bash
$ cd ~/wirelineio/wns/ops
$ ansible-playbook -i env/development reset-full.yml
```
