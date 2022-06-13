# How to create a cdk environment

# Install node

```sh
cat <<EOF > install_node.yaml
- name: Installer
  hosts: all
  tasks:
  - args:
      creates: "{{ ansible_env.HOME }}/.nvm/versions/node/"
      executable: /usr/bin/bash
    name: Install node
    shell: |
      curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v{{ nvm_version }}/install.sh | bash
      source {{ ansible_env.HOME }}/.nvm/nvm.sh
      nvm install node
EOF
ansible-playbook -ilocalhost, -clocal -envm_version=0.39.1 install_node.yaml
source ~/.bashrc
```

# Install cdk
```sh
npm i -g aws-cdk
```

# Install go
https://go.dev/dl/
```sh
cat <<EOF > install_go.yaml
- name: Installer
  hosts: all
  tasks:
  - become: "yes"
    name: Install go
    unarchive:
      creates: /usr/local/go/
      dest: /usr/local/
      remote_src: "yes"
      src: |-
        https://go.dev/dl/go{{ version }}.linux-amd64.tar.gz
EOF
ansible-playbook -ilocalhost, -clocal -eversion=1.18.3 install_go.yaml -K
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```