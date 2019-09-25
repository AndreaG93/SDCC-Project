# ==================================================================================================================== #
# To build all processes locally...
# ==================================================================================================================== #

# WP -- WPG #0

for x in 0 1 2
do
    go build  -o $HOME/WP/worker-$x/wcworker-$x $HOME/go/src/SDCC-Project/aftmapreduce/main/wcserver.go
    echo '{"ZookeeperServersPrivateIPs": ["localhost"], "NodeID": '$x', "NodeGroupID": 0, "NodeClass": "Worker"}' > $HOME/WP/worker-$x/conf.json
done

# WP -- WPG #1

for x in 3 4 5
do
    go build  -o $HOME/WP/worker-$x/wcworker-$x $HOME/go/src/SDCC-Project/aftmapreduce/main/wcserver.go
    echo '{"ZookeeperServersPrivateIPs": ["localhost"], "NodeID": '$x', "NodeGroupID": 1, "NodeClass": "Worker"}' > $HOME/WP/worker-$x/conf.json
done

# PP

for x in 1 2 3
do
    go build  -o $HOME/PP/primary-$x/wcprimary-$x $HOME/go/src/SDCC-Project/aftmapreduce/main/wcserver.go
    echo '{"ZookeeperServersPrivateIPs": ["localhost"], "NodeID": '$x', "NodeGroupID": 0, "NodeClass": "Primary"}' > $HOME/PP/primary-$x/conf.json
done

go build -o $HOME/wcclient $HOME/go/src/SDCC-Project/aftmapreduce/process/client/main/wcclient.go
cp $HOME/go/src/SDCC-Project/test-input-data/input1.txt $HOME/input
