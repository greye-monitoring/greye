




mkdir greye
cd greye
echo '
notification:
  telegram:
    destination: "27664341"
    token: "6524810996:AAEu_XhUBmhDe4HPRQ90qGNmguvmZVi3zAk"
  email:
    destination: "27664341"
    token: "6524810996:AAEu_XhUBmhDe4HPRQ90qGNmguvmZVi3zAk"
monitoring:
  enabled: false
cluster:
  intervalSeconds: 60
  timeoutSeconds: 5
  maxFailedRequests: 3
  myIp: "192.168.1.24"
  ip:
  - 192.168.1.26:32473
'> values.yaml
helm repo update
helm upgrade greye greye/greye -f values.yaml -n greye