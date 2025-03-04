
version=`cat ./VERSION`

version=$((version+1))

docker build -f deploy/Dockerfile -t 192.168.1.24:30515/cm:$version .
docker tag 192.168.1.24:30515/cm:$version ftrigari/greye:$version
docker tag 192.168.1.24:30515/cm:$version ftrigari/greye:latest

docker push 192.168.1.24:30515/cm:$version
docker push ftrigari/greye:$version
docker push ftrigari/greye:latest

echo $version > ./VERSION
helm repo update
helm upgrade greye greye/greye --set image.tag=$version

