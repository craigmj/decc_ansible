#!/bin/bash
set -e
ME=`whoami`
if [ ! -f /etc/sudoers.d/`whoami` ]; then
	echo "$ME ALL=(ALL:ALL) NOPASSWD:ALL" | sudo tee /etc/sudoers.d/$ME
fi
sudo apt-get update
sudo apt-get install -y curl software-properties-common python-software-properties \
  python-pycurl vim htop build-essential git mercurial nodejs bundler unzip ruby
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
 sudo chown $ME:$ME /opt/decc
fi
cd /opt/decc
if [ ! -d decc_2050_model ]; then
 git clone https://github.com/craigmj/decc_2050_model.git
fi
cd decc_2050_model
# Undo any local changes caused by build (in particular bundle exec rake)
git revert --hard HEAD
git pull origin master
gem install bundler
bundle
##bundle exec rake
gem build model.gemspec
gem install decc_2050_model-0.80.20140325pre.gem
cd ..
if [ ! -d twenty-fifty ]; then
git clone https://github.com/craigmj/twenty-fifty.git
fi
cd twenty-fifty
git pull origin master
bundle install

cat >/opt/decc/decc2050.upstart.conf <<EOUPSTART
# decc2050 server
#
description     "decc2050 server"

start on (runlevel [2345])
stop on runlevel [!2345]

respawn
expect fork

script
  sudo -u $ME -- /bin/bash -l -c "cd /opt/decc/twenty-fifty; rvm 2.1-head do rackup" &
end script

emits decc2050_starting
EOUPSTART
sudo mv /opt/decc/decc2050.upstart.conf /etc/init/decc2050.conf
sudo start decc2050
