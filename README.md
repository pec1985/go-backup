### go-backup

A simple deamon that will back up your project every two minutes using Git

#### Usage:

To start backing up (after building go-backup and placing it in your PATH)

```
cd {your project}
go-backup start
```

To stop

```
go-backup stop
```

To restore

```
go-backup restore {project_name}
```


#### Optional parameters

You can override the default values

```
cd {project}
go-backup start --path {backup root dir} --name {project name} --ignore_path  {folder to ignore}
```

Default values

```
--path ~/backups
--name {current folder}
--ignore_path node_modules
--ignore_path dist
--ignore_path build
```


Pedro Enrique  
MIT