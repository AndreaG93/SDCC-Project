sudo zkSrv.sh start
konsole --new-tab --noclose --workdir $HOME/PP/primary-1 -e $HOME/PP/primary-1/wcprimary-1 local &

sleep 1

for x in 2 3
do
    konsole --new-tab --noclose --workdir $HOME/PP/primary-$x -e $HOME/PP/primary-$x/wcprimary-$x local &
done

for x in 0 1 2 3 4 5
do
    konsole --new-tab --noclose --workdir $HOME/WP/worker-$x -e $HOME/WP/worker-$x/wcworker-$x local &
done

sleep 2

$HOME/wcclient $HOME/input localhost:2181
