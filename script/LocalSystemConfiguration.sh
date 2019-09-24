# ==================================================================================================================== #
# To build all processes locally...
# ==================================================================================================================== #

# WP -- WPG #0

for x in 0 1 2
do
    go build  -o $HOME/WP/worker-$x/wcworker-$x $HOME/go/src/SDCC-Project/aftmapreduce/main/wcserver.go
    echo '{"ZookeeperServersPrivateIPs": ["localhost"], "NodeID": '$x', "NodeGroupID": 0, "NodeClass": "Worker"}' > ./WP/worker-$x/conf.json
done

# WP -- WPG #1

for x in 3 4 5
do
    go build  -o $HOME/WP/worker-$x/wcworker-$x $HOME/go/src/SDCC-Project/aftmapreduce/main/wcserver.go
    echo '{"ZookeeperServersPrivateIPs": ["localhost"], "NodeID": '$x', "NodeGroupID": 1, "NodeClass": "Worker"}' > ./WP/worker-$x/conf.json
done

# PP

for x in 1 2 3
do
    go build  -o $HOME/PP/primary-$x/wcprimary-$x $HOME/go/src/SDCC-Project/aftmapreduce/main/wcserver.go
    echo '{"ZookeeperServersPrivateIPs": ["localhost"], "NodeID": '$x', "NodeGroupID": 0, "NodeClass": "Primary"}' > ./PP/primary-$x/conf.json
done

# Client
go build -o $HOME/wcclient $HOME/go/src/SDCC-Project/aftmapreduce/process/client/main/wcclient.go


# ==================================================================================================================== #
# To start all processes locally...
# ==================================================================================================================== #

for x in 1 2 3
do
    konsole --new-tab --noclose --workdir $HOME/PP/primary-$x -e $HOME/PP/primary-$x/wcprimary-$x local &
done

for x in 0 1 2 3 4 5
do
    konsole --new-tab --noclose --workdir $HOME/WP/worker-$x -e $HOME/WP/worker-$x/wcworker-$x local &
done
