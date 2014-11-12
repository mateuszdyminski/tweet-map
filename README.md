# Install all js dependencies 
bower install

# Install all go dependencies

# Run app
export GOPATH=$PWD && go build -o bin/main main && bin/main --log_dir="/tmp" --logtostderr=1 --stderrthreshold=0 -p 9000 -h localhost -dir app