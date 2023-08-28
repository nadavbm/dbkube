# dbkube

kubernetes operator built with [kubebuilder](https://book.kubebuilder.io/introduction.html)

the operator will create cronjob by using crd inside a namespace

### create a builder with kubebuilder

this repository will work with several apis and controllers

to initialise with kubebuilder with multiple groups

to initialize kubebuilder project run:
```
kubebuilder init --domain etz.com --repo github.com/nadavbm/dbkube
kubebuilder edit --multigroup=true
```

create the relevant api for running a database service instance in kubernetes:
```
kubebuilder create api --group configmaps --version v1alpha1 --kind ConfigMap
kubebuilder create api --group secrets --version v1alpha1 --kind Secret
kubebuilder create api --group services --version v1alpha1 --kind Service
kubebuilder create api --group deployments --version v1alpha1 --kind Deployment
```

edit and generate crd (whenever a new api group was added):
```
make generate
make manifests
```

to push docker image (release a new version):
```
make docker-build docker-push IMG="nadavbm/dbkube:v0.0.1"
```

### testing operator with minikube

to test the operator, run minikube and the following commands:
```
sh test/cleanup.sh 
kubectl apply -f test/specs.yaml 
kubectl apply -f test/testcrd.yaml
```

test crd changes with (connect to the relevant namespace with `kubens`):
```
k edit jobs.cronjobs.example.com jobop
```
