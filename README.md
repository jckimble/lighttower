Light Tower
============

Tool to watch codeship builds and send dbus notifications

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

License
-------
[GPLv3][gpl3.0]

[gpl3.0]: https://www.gnu.org/licenses/gpl-3.0.txt

Contributing
------------
Please follow the [Open Code of Conduct][code-of-conduct].

[code-of-conduct]: http://todogroup.org/opencodeofconduct

To make sure your pull request will be accepted, please open an issue in the issue tracker before starting work where we can talk to make sure a feature or bug fix is done right.
