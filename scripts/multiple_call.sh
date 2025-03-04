echo "localhost:8080"
curl -s http://localhost:8080/api/v1/cluster/status | jq '.["localhost:7080"] | {timestamp, found_by: .error.found_by, error_count: .error.count}'

echo "localhost:8070"
curl -s http://localhost:8070/api/v1/cluster/status | jq '.["localhost:7080"] | {timestamp, found_by: .error.found_by, error_count: .error.count}'

echo "localhost:8060"
curl -s http://localhost:8060/api/v1/cluster/status | jq '.["localhost:7080"] | {timestamp, found_by: .error.found_by, error_count: .error.count}'

echo "localhost:8050"
curl -s http://localhost:8050/api/v1/cluster/status | jq '.["localhost:7080"] | {timestamp, found_by: .error.found_by, error_count: .error.count}'

echo "localhost:8040"
curl -s http://localhost:8040/api/v1/cluster/status | jq '.["localhost:7080"] | {timestamp, found_by: .error.found_by, error_count: .error.count}'

echo "localhost:8030"
curl -s http://localhost:8030/api/v1/cluster/status | jq '.["localhost:7080"] | {timestamp, found_by: .error.found_by, error_count: .error.count}'
echo "localhost:8020"
curl -s http://localhost:8020/api/v1/cluster/status | jq '.["localhost:7080"] | {timestamp, found_by: .error.found_by, error_count: .error.count}'

echo "localhost:8010"
curl -s http://localhost:8010/api/v1/cluster/status | jq '.["localhost:7080"] | {timestamp, found_by: .error.found_by, error_count: .error.count}'

echo "localhost:8000"
curl -s http://localhost:8000/api/v1/cluster/status | jq '.["localhost:7080"] | {timestamp, found_by: .error.found_by, error_count: .error.count}'

echo "localhost:7090"
curl -s http://localhost:7090/api/v1/cluster/status | jq '.["localhost:7080"] | {timestamp, found_by: .error.found_by, error_count: .error.count}'
