Light Tower
============

Tool to watch codeship builds and send dbus notifications

Note
----
LightTower will not be v1 until codeship's v2 api is finalized, expect rapid version changes up to that time, until then all feature request will be considered so don't be afraid to ask or contribute.

Installation
------------
```
go get github.com/jckimble/lighttower
cd lighttower
go build -o lighttower main.go
```

Usage
-----
```
lighttower create -u [username] -p [password]
lighttower watch
```

Config (~/.lighttower.yaml)
-------
```
auth: (basic auth string)
SuccessImage: (can be image path or dbus image)
SuccessMessage: (message you want echoed on successful build)
SuccessSound: (can be sound path or dbus sound)
ErrorImage: (can be image path or dbus image)
ErrorMessage: (message you want echoed on failed build)
ErrorSound: (can be sound path or dbus sound)
```

License
-------
[GPLv3][gpl3.0]

[gpl3.0]: https://www.gnu.org/licenses/gpl-3.0.txt

Contributing
------------
Please follow the [Open Code of Conduct][code-of-conduct].

[code-of-conduct]: http://todogroup.org/opencodeofconduct

To make sure your pull request will be accepted, please open an issue in the issue tracker before starting work where we can talk to make sure a feature or bug fix is done right.
