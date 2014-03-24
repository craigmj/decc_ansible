from fabric.api import *
from fabric.context_managers import shell_env

env.hosts = ["craig@2050.lateral.co.za"]

basePath = "/opt/decc/"
modelPath = basePath + "decc_2050_model"
webPath = basePath + "twenty_fifty"

def deploy():
	with cd(modelPath):
		run("git reset --hard")
		run("git pull")
		run("gem build model.gemspec")
		sudo("gem install decc_2050_model-3.5.1pre.gem")
	with cd(webPath):
		run("git reset --hard")
		run("git pull")
		run("cp Gemfile.v351pre Gemfile")
		run("gem clean")
		run("gem install")
		sudo("restart decc2050")
