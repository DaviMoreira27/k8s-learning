[Enables](https://github.com/nicolaka/netshoot) networking trouble-shooting. It pushes a temporarly pod that enables the use of common network execs and bash commands. 
```sh
kubectl run netshoot --rm -it --image=nicolaka/netshoot -- bash
```
