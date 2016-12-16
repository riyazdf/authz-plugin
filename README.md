## Build

`docker build -t . rootfs`

`docker plugin create riyaz/authz-no-volume-plugin .`


## Ship

`docker plugin push riyaz/authz-no-volume-plugin`


## Run

`docker plugin install riyaz/authz-no-volume-plugin`

`dockerd --authorization-plugin=riyaz/authz-no-volume-plugin:latest`
