#!/bin/bash
set -e

# Setup sudo permissions for current user, and load some general utilities we'll require
if [ ! -f /etc/sudoers.d/`whoami` ]; then
	echo "`whoami` ALL=(ALL:ALL) NOPASSWD:ALL" | sudo tee /etc/sudoers.d/`whoami`
fi
sudo apt-get update
sudo apt-get install -y curl software-properties-common python-software-properties \
  python-pycurl vim htop build-essential git golang mercurial nodejs bundler unzip ruby

# Install rvm so that we can somehow contain ruby madness
gpg --keyserver hkp://keys.gnupg.net --recv-keys D39DC0E3
curl -sSL https://get.rvm.io | bash -s stable --ruby

# Load RVM into our shell session *as a function*
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
rvm install ruby-2.1.5
rvm use 2.1.5


# Create directory for all decc files
if [ ! -d /opt/decc ]; then
 sudo mkdir -p /opt/decc
 sudo chown `whoami`:`whoami` /opt/decc
fi
cd /opt/decc

## Install the deccgem utility
# It will add itself to /usr/bin and hence hopefully be available to all scripts
if [ ! -d deccgem ]; then
  git clone https://github.com/craigmj/deccgem.git
fi
pushd deccgem
git pull origin master
./build.sh
popd

# Install the decc 2050 model
if [ ! -d decc_2050_model ]; then
 git clone https://github.com/craigmj/decc_2050_model.git
fi
pushd decc_2050_model
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
popd

## Install the twenty-fifty system
if [ ! -d twenty-fifty ]; then
git clone https://github.com/craigmj/twenty-fifty.git
fi
pushd twenty-fifty
git pull origin master

rm -f Gemfile.lock
deccgem gemfile > Gemfile
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

setuid `whoami`
script
  chdir /opt/decc/twenty-fifty
  rvm 2.1.5 do bundle exec rackup -o 0.0.0.0 &
end script

emits decc2050_starting
EOUPSTART

  sudo service decc2050 start
else
  sudo service decc2050 restart
fi
popd
