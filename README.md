JResponse: Encapsulates responses for various commands
executed on Juniper devices via NETCONF.
===============================================

- [Intro](#intro)
- [Installation](#installation)
- [User Guide](#user-guide)
- [License](#license)

Intro
-----

JResponse is a library for handling various responses to commands executed on Juniper devices
 via NETCONF. It was made to make transforming the response data into other formats easier.

RPCs currently supported are:

    - traceroute
    - show route protocol bgp

Responses can be converted from their native XML RPC reply into JSON, or they can be formatted
to look as they would look if run directly on the Junos CLI.

Installation
------------

1. Include the library path in your project.
2. From your projects directory, run `go get`
3. Use the library.


User Guide
----------


License
-------
This software is licensed under the BSD 2-clause “Simplified” License
 ©2015 JResponse contributors