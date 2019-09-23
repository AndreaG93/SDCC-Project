
foreach ($x in @("0","1","2","3","4","5")) {
	go build  -o ./WP/worker-$x/wcworker-$x.exe $HOME/go/src/SDCC-Project/aftmapreduce/main/wcserver.go
}

foreach ($x in @("0","1","2")) {
	go build  -o ./PP/primary-$x/wcprimary-$x.exe $HOME/go/src/SDCC-Project/aftmapreduce/main/wcserver.go
}

foreach ($x in @("0","1","2")) {
	echo "{""ZookeeperServersPrivateIPs"": [""localhost""], ""NodeID"": $x, ""NodeGroupID"": 0, ""NodeClass"": ""Worker""}" | Out-File -Encoding ascii ./WP/worker-$x/conf.json
}

foreach ($x in @("3","4","5")) {
	echo "{""ZookeeperServersPrivateIPs"": [""localhost""], ""NodeID"": $x, ""NodeGroupID"": 1, ""NodeClass"": ""Worker""}" | Out-File -Encoding ascii ./WP/worker-$x/conf.json
}

# For starting...
foreach ($x in @("0","1","2","3","4","5")) {
	$cmd="-NoExit ""./wcworker-$x.exe local""" 
	start powershell $cmd -WorkingDirectory .\WP\worker-$x\
}

foreach ($x in @("0","1","2")) {
	$cmd="-NoExit ""./wcprimary-$x.exe local""" 
	start powershell $cmd -WorkingDirectory .\PP\primary-$x\
}
	

	
cat localfile.conf | ssh user@hostname 'cat -> /tmp/remotefile.conf'

echo "{""ZookeeperServersPrivateIPs"": [""localhost""], ""NodeID"": $x, ""NodeGroupID"": 1, ""NodeClass"": ""Worker""}" | Out-File -Encoding ascii ./conf.json
cat ./conf.json | ssh -i "graziani.pem" ubuntu@ec2-54-172-118-113.compute-1.amazonaws.com 'cat -> ./conf.conf'