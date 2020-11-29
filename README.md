### go-backup

A simple daemon that will back up your project every two minutes using Git

#### Installation:

```
go get
go install -o backup
mv backup /usr/local/bin
```

#### Usage:

To start backing up

```
cd {your project}
backup start
```

To stop

```
backup stop
```

To restore

```
backup restore {project_name}
```


#### Optional parameters

You can override the default values

```
cd {project}
backup start --path {backup root dir} --name {project name} --ignore_path  {folder to ignore}
```

Default values

```
--path ~/.backups
--name {current folder}
--ignore_path node_modules
--ignore_path dist
--ignore_path build
```


Pedro Enrique  
MIT