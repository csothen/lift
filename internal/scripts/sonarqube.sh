#!/bin/bash

# Script to setup an instance of sonarqube in a machine

# Variables
VERSION=8.9
READING_URLS=false
PLUGIN_URLS=()

# prints the help text
function usage() {
    cat <<EOF
    Usage: $0 [-v version] [-p list of plugin download urls]

    Options:
        -v --version:   sonarqube version to install
        -p --plugins:   sonarqube plugin download urls that are to be installed
        -h --help:      show this help

    Description:
        $0 will install sonarqube on an Ubuntu instance with all necessary dependencies and the desired plugins.
        The script will take a list of download URLs of the plugins which will be used to download them and move
        them to their correct location for sonarqube to use them.
EOF
    exit 1
}

# reads the options from the arguments
function read_options() {
    while [ "$1" != "" ]; do
        case $1 in
        -v | --version)
            shift
            READING_URLS=false
            VERSION=$1
            ;;
        -p | --plugins)
            shift
            READING_URLS=true
            PLUGIN_URLS+=("$1")
            ;;
        -h | --help)
            usage
            exit 1
            ;;
        *)
            if [ $READING_URLS = true ]; then
                PLUGIN_URLS+=("$1")
            else
                usage
                exit 1
            fi
            ;;
        esac
        shift
    done
}

# Read all the options from the arguments passed
read_options $@

# Set up a non-root user for Sonarqube
# create the group "sonar" for the user
sudo groupadd sonar
# create the user and add it to the group "sonar"
sudo useradd -c "Sonarqube User" -d /opt/sonarqube -g sonar -s /bin/bash sonar
# set its password to "sonarqubepassword"
usermod --password $(echo sonarqubepassword | openssl passwd -1 -stdin) sonar
# add the user to sudoers
sudo usermod -a -G sonar ec2-user

# Install Java 11

# download it
curl -O https://download.java.net/java/GA/jdk11/13/GPL/openjdk-11.0.1_linux-x64_bin.tar.gz
# unzip and move files to correct place
tar zxvf openjdk-11.0.1_linux-x64_bin.tar.gz
sudo mv jdk-11.0.1 /usr/local/
# change access from the jdk folder
sudo chmod -R 755 /usr/local/jdk-11.0.1
# set up JAVA_HOME
export JAVA_HOME=/usr/local/jdk-11.0.1
export PATH=$JAVA_HOME/bin:$PATH
# load changes
source /etc/profile

# Install Sonarqube

# download it
# TODO: Replace the hardcoded version with something dynamic
wget https://binaries.sonarsource.com/Distribution/sonarqube/sonarqube-8.9.1.44547.zip
#unzip it
unzip sonarqube-*.zip
# move sources to the approriate folder
sudo mv -v sonarqube-*/* /opt/sonarqube
# change ownership of all the sonarqube files to user sonar
sudo chown -R sonar:sonar /opt/sonarqube
# change file access privileges
sudo chmod -R 775 /opt/sonarqube

# Install the Sonarqube plugins

# download every plugin
for plugin_url in "$PLUGIN_URLS"
do
    wget $plugin_url
done
# move them all to the plugins folder of sonarqube
sudo mv *.jar /opt/sonarqube/extensions/plugins

# Configure the Sonarqube instance
# TODO: Add the configuration logic

# Start Sonarqube
echo sonarqubepassword | /opt/sonarqube/bin/linux-x86-64/sonar.sh start