import yaml

NUM_SERVICES = 300
NAMESPACE = "te"
APP_NAME = "printall"
OUTPUT_FILE = "services.yaml"

def generate_service_yaml(index):
    return {
        "apiVersion": "v1",
        "kind": "Service",
        "metadata": {
            "name": f"{APP_NAME}-{index}",
            "namespace": NAMESPACE,
            "annotations": {
                "ge-enabled": "true",
                "ge-intervalSeconds": "30",
                "ge-paths": "|-\n/pippo\n/prova\n/cane"
               },
            "labels": {
                "app": APP_NAME`1
            }
        },
        "spec": {
            "ipFamilies": ["IPv4"],
            "ipFamilyPolicy": "SingleStack",
            "ports": [
                {
                    "name": "http",
                    "port": 80,
                    "protocol": "TCP",
                    "targetPort": 5001
                }
            ],
            "selector": {
                "app": APP_NAME
            },
            "sessionAffinity": "None",
            "type": "ClusterIP"
        },
        "status": {
            "loadBalancer": {}
        }
    }

if __name__ == "__main__":
    services = [generate_service_yaml(i) for i in range(1, NUM_SERVICES + 1)]
    with open(OUTPUT_FILE, "w") as f:
        yaml.dump_all(services, f, default_flow_style=False)
    print(f"Generated {NUM_SERVICES} services in {OUTPUT_FILE}")
