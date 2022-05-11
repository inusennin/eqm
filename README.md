# eqm
=====

Go In Memory Event Queue Manager

Package eqm is an implementation framework for handling in memory event queues.

The premise is to allow a unified control module for handling multiple threads
that each handle management of a shared resource without the use of mutex locks
or wait groups.

# Supported Versions
--------------------

As of writing this, the implementation the builds are tested on is version 1.12.
To the author's knowledge, all subsequent versions should work without issue.

# Disclaimer
------------

This module is currently undocumented and a work in progress.

# License
---------
> Copyright (c) 2022 The eqm Authors. All rights reserved.
> Use of this source code is governed by a BSD-style
> license that can be found in the LICENSE file.
