# How to create a cdk environment

# Create a VM
```sh
(
more_memory=4096
name=cdk

curl -o preseed.cfg https://www.debian.org/releases/bullseye/example-preseed.txt
cat <<EOF >> preseed.cfg
d-i grub-installer/bootdev string default
d-i netcfg/get_hostname string debian11-$name
d-i passwd/root-login boolean false
d-i passwd/user-fullname string whs
d-i passwd/user-password password magic
d-i passwd/user-password-again password magic
d-i passwd/username string whs
d-i preseed/late_command string apt-install openssh-server
popularity-contest popularity-contest/participate boolean false
tasksel tasksel/first multiselect standard
EOF
virt-install \
  --disk size=20 \
  --initrd-inject preseed.cfg \
  --location https://deb.debian.org/debian/dists/bullseye/main/installer-amd64/ \
  --memory ${more_memory:-2048} \
  --name debian11-$name \
  --network bridge:virbr0 \
  --os-variant debian11 \
  --vcpus 2 \
  ;
) &
```

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

# Install go (development only)
You do **not** need go to run `cdk deploy`

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

# Install vim-go (optional)
```sh
git clone https://github.com/fatih/vim-go.git ~/.vim/pack/plugins/start/vim-go
vim +GoInstallBinaries +q
```

# Install build essential
```sh
sudo apt-get install -y build-essential
```

# Install bazel
https://github.com/bazelbuild/bazelisk/releases
```sh
cat <<EOF > install_bazelisk.yaml
- name: Installer
  hosts: all
  tasks:
  - become: "yes"
    get_url:
      dest: /usr/local/bin/bazelisk
      mode: "0755"
      url: |-
        https://github.com/bazelbuild/bazelisk/releases/download/v{{ version }}/bazelisk-linux-amd64
    name: Install bazelisk
EOF
ansible-playbook -ilocalhost, -clocal -eversion=1.12.0 install_bazelisk.yaml -K
```

# Install buildifier
```sh
go install github.com/bazelbuild/buildtools/buildifier@latest
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
source ~/.bashrc
```

# Install awscli
```sh
# Copied from https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install
```
