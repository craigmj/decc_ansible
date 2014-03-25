from fabric.api import *
from fabric.context_managers import shell_env

env.hosts = ["craig@2050.lateral.co.za"]

basePath = "/opt/decc/"
modelPath = basePath + "decc_2050_model"
webPath = basePath + "twenty_fifty"

def deploy():
	with cd(modelPath):
		run("gem clean")
		run("git reset --hard")
		run("git pull origin master")
		sudo("gem install `deccgem whichgem`")
	with cd(webPath):
		run("git reset --hard")
		run("git pull origin master")
		run("deccgem gemfile > Gemfile")
		run("bundle install")
		sudo("restart decc2050")

# startup.sh
# #!/bin/bash

# gem clean
# cd ../decc_2050_model
# gem build model.gemspec
# gem install decc_2050_model-0.71.20140319pre.gem
# cd ../twenty-fifty
# bundle install
# rackup
