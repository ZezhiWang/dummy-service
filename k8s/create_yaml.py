import json
import sys


def generate_yaml(json_file, yaml_file, tear_down):
    with open(json_file, 'r') as f:
        obj = json.load(f)

    tear_down_list = ["#!/bin/bash\n"]
    yaml_list = []
    for svc in obj["svc"]:
        tear_down_list.append(
            "kubectl delete deployment " + svc + "\n" +
            "kubectl delete service " + svc + "\n"
        )

        yaml_list.append(
            "apiVersion: v1\n" +
            "kind: Service\n" +
            "metadata:\n" +
            "  name: " + svc + "\n" +
            "  labels:\n" +
            "    app: " + svc + "\n" +
            "spec:\n" +
            "  selector:\n" +
            "    app: " + svc + "\n" +
            "  ports:\n" +
            "    - port: 8888\n" +
            "      targetPort: 8888\n" +
            "---\n" +
            "apiVersion: apps/v1\n" +
            "kind: Deployment\n" +
            "metadata:\n" +
            "  name: " + svc + "\n" +
            "  labels:\n" +
            "    app: " + svc + "\n" +
            "spec:\n" +
            "  replicas: 1\n" +
            "  selector:\n" +
            "    matchLabels:\n" +
            "      app: " + svc + "\n" +
            "  template:\n" +
            "    metadata:\n" +
            "      labels:\n" +
            "        app: " + svc + "\n" +
            "    spec:\n" +
            "      containers:\n" +
            "        - name: " + svc + "\n" +
            "          image: docker.io/herryzzwang/dummy:latest\n" +
            "          imagePullPolicy: Always\n" +
            "          ports:\n" +
            "            - containerPort: 8888\n" +
            "          env:\n" +
            "            - name: CHILD_SERVICE\n" +
            "              value: \"" + obj["svc"][svc] + "\""
        )

    with open(yaml_file, 'w') as f:
        f.write("\n---\n".join(yaml_list))
    with open(tear_down, 'w') as f:
        f.write("\n".join(tear_down_list))


if __name__ == '__main__':
    generate_yaml("dummy.json", "dummy.yaml", "tear_down.sh")
