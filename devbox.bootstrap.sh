#!/bin/bash
set -e

if [ ! -f /etc/sudoers.d/`whoami` ]; then
	echo "`whoami` ALL=(ALL:ALL) NOPASSWD:ALL" | sudo tee /etc/sudoers.d/`whoami`
fi
sudo apt-get update
sudo apt-get install -y curl software-properties-common python-software-properties \
  python-pycurl vim htop build-essential git golang mercurial nodejs bundler unzip ruby
# Load RVM into a shell session *as a function*
if [[ -s "$HOME/.rvm/scripts/rvm" ]] ; then
  # First try to load from a user install
  source "$HOME/.rvm/scripts/rvm"
  RVMB="$HOME/.rvm/bin/rvm"
elif [[ -s "/usr/local/rvm/scripts/rvm" ]] ; then
  # Then try to load from a root install
  source "/usr/local/rvm/scripts/rvm"
  RVMB="/usr/local/rvm/bin/rvm"
else
 curl -sSL https://get.rvm.io | bash
 source "$HOME/.rvm/scripts/rvm"
 RVMB="$HOME/.rvm/bin/rvm"
fi
rvm install ruby-2.1-head
rvm use 2.1-head
if [ ! -d /opt/decc ]; then
 sudo mkdir -p /opt/decc
 sudo chown `whoami`:`whoami` /opt/decc
fi
cd /opt/decc
if [ ! -d decc_2050_model ]; then
 git clone https://github.com/craigmj/decc_2050_model.git
fi
cd decc_2050_model
# Undo any local changes caused by build (in particular bundle exec rake)
git reset --hard HEAD
git pull origin master
gem install bundler
bundle
#bundle exec rake
rm -f *.gem
gem build model.gemspec
gem uninstall -a decc_2050_model
gem install `ls -Rt decc_2050_model*.gem | head -1`
cd ..

## Now install the deccgem utility
if [ ! -d decc_ansible ]; then
	git clone https://github.com/craigmj/decc_ansible.git
fi
pushd decc_ansible
git pull origin master
cd deccgem

popd

## Install the twenty-fifty system
if [ ! -d twenty-fifty ]; then
git clone https://github.com/craigmj/twenty-fifty.git
fi
pushd twenty-fifty
git pull origin master

if [ ! -f deccgem/deccgem ]; then
	cd deccgem
	GOPATH=`pwd` go build src/cmd/deccgem.go
	sudo rm -f /usr/bin/deccgem
	sudo rm -f /usr/local/bin/deccgem
	if [ ! -f /usr/bin/deccgem ]; then
	  sudo ln -s `pwd`/deccgem /usr/local/bin
	fi
	cd ..
fi

rm -f Gemfile.lock
deccgem/deccgem gemfile > Gemfile
bundle install

if [ ! -f /etc/init/decc2050.conf ]; then
  sudo tee /etc/init/decc2050.conf > /dev/null << EOUPSTART
# decc2050 server
#
description     "decc2050 server"

start on (runlevel [2345])
stop on runlevel [!2345]

respawn
expect fork

script
  sudo -u `whoami` -- /bin/bash -l -c "cd /opt/decc/twenty-fifty; rvm 2.1-head do bundle exec rackup" &
end script

emits decc2050_starting
EOUPSTART

  sudo start decc2050
else
  sudo restart decc2050
fi
popd
