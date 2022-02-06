# Deploy


### Kubernetes deployment

The _k8s_manifest.yaml_ file creates a configMap, service (LB) and a deployment with the latest version of _Gotcha_. The only necessary adjustment is in the **configMap** which has the settings to run the app, so you need to put your information for the application to work.  

Finally, just run a _kubectl apply -f **file.yaml**_