# Ops

## Setup

* Install [Ansible](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html).

## Operations

```bash
cd ~/wirelineio/wns/ops
```

Restart `wnsd` on servers:

```bash
$ ansible-playbook -i env/development.yml restart.yml
```

## Update Code

```bash
$ ansible-playbook -i env/development.yml update.yml
```