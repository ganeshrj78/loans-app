# Small businesses Loans

This app is meant to demo multi-tenant deployment scenarios using Apache Knox

The application is a fairly contrived one where multiple lending institutions use the application as a
portal for originators to use to submit loan applications for small businesses.

The application itself is written in golang and uses Apache Knox in the following ways:

1. To proxy the UI of the application
2. To provide the SSO capabilities
3. To be the reverse proxy to the Hadoop cluster


# Setup

The demo setup expects to run on a laptop all locally (with one or more VMs and or docker containers).

To be able to simulate multiple domains the following is needed in the /etc/hosts file

127.0.0.1 www.unwise.com

127.0.0.1 www.goodloans.com

These two domains will be the two tenants of the application.

The topology files under /knox-files in this repo should be deployed to a knox instance.

The URLs to get to the app are therefore :

https://www.goodloans.com:8443/gateway/goodloans/loans/

and 

https://www.unwise.com:8443/gateway/unwise/loans/

The topology files assume a Hadoop cluster running on an Ambari Vagrant machine c6401.ambari.apache.org

The following setup is needed on that VM:

1. Add users and groups for the tenants

 groupadd loanscore

 useradd loanscore -g loanscore

 useradd admin -g loanscore
 
 groupadd unwise

 groupadd goodloans


2. Add the originator users to the tenant groups

 useradd bob_unwise -g unwise

 useradd bob_goodloans -g goodloans

3. Create the folders for the tenants (as user hdfs)

hdfs dfs -mkdir -p /unwise/applications

hdfs dfs -mkdir -p /goodloans/applications

hdfs dfs -mkdir -p /loanscore/scores

hdfs dfs -chmod 770 /unwise/applications

hdfs dfs -chmod 770 /loanscore/scores

hdfs dfs -chmod 770 /unwise/applications

hdfs dfs -chown -R hdfs:unwise /unwise

hdfs dfs -chown -R hdfs:goodloans /goodloans

hdfs dfs -chown -R loanscore:loanscore /loanscore

4. To seed the folders 

(in this local project)

cd hadoop-files

cd goodloans

$KNOXSHELL_HOME/bin/knoxshell.sh goodloans-seed.groovy

cd ../unwise

$KNOXSHELL_HOME/bin/knoxshell.sh unwise-seed.groovy 
