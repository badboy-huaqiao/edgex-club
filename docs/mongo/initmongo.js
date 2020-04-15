// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0 

(function() {
    db = db.getSiblingDB('admin');
    try {
        db.createUser({ user: 'root', pwd: 'root', roles: [{ role: 'userAdminAnyDatabase', db: 'admin' }, { role: 'readWrite', db: 'admin' }] });
    } catch (e) {
        return;
    }
    db.auth('root', 'root');

    db = db.getSiblingDB('edgex-club');
    db.createUser({ user: 'edgex-club-user', pwd: '1234', roles: [{ role: 'readWrite', db: 'edgex-club' }] });
    db.auth('edgex-club-user', '1234');
}());