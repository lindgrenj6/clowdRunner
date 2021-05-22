# clowdRunner

A tool for running commands in a clowder pod, basically the program reads the cdappconfig.json and sets the ENV vars for a sub-command. 

### Building
type `make`

then you'll have a `clowdRun` binary. Copy that into your container for one-off debugging or host it somewhere and pull it during your docker build.
