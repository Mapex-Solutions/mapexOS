// =============================================================
// MQTT Assets Seed — Auth Callout Benchmark
//
// Inserts MQTT-protocol assets for auth callout benchmarks.
// Count and startHex are injected by seed.sh via variable substitution.
//
// Each asset has:
//   - protocol.type = "mqtt"
//   - protocol.mqtt.username = "bench-mqtt-device-NNNNN"
//   - protocol.mqtt.password = "bench-mqtt-secret-NNNNN"
//   - assetUUID = "bench-mqtt-device-NNNNN"
//
// Username padding: 5 digits (00001..10000+) to support large counts.
//
// Usage (standalone):
//   mongosh "mongodb://localhost:27017/?replicaSet=rs0" --eval "
//     db = db.getSiblingDB('dev-assets');
//     var __COUNT__ = 10000;
//     var __START_HEX__ = 0x2001;
//     $(cat mqtt-assets.js)
//   "
//
// Or via seed.sh (recommended):
//   bash seed.sh setup
// =============================================================

var now = new Date();

// Variables are set by seed.sh before this script runs.
// Fallback defaults for standalone usage.
if (typeof __COUNT__ === 'undefined')      { var __COUNT__ = 100; }
if (typeof __START_HEX__ === 'undefined')  { var __START_HEX__ = 0x2001; }
if (typeof __ORG_ID__ === 'undefined')     { var __ORG_ID__ = '000000000000000000000099'; }
if (typeof __TEMPLATE_ID__ === 'undefined'){ var __TEMPLATE_ID__ = '000000000000000000000010'; }
if (typeof __CUSTOMER_ID__ === 'undefined'){ var __CUSTOMER_ID__ = '000000000000000000000098'; }
if (typeof __PATHKEY__ === 'undefined')    { var __PATHKEY__ = 'bench-mqtt'; }

var orgId      = ObjectId(__ORG_ID__);
var templateId = ObjectId(__TEMPLATE_ID__);
var customerId = ObjectId(__CUSTOMER_ID__);

// Insert in batches of 1000 to avoid BSON document size limits
var batchSize = 1000;
var totalInserted = 0;
var totalSkipped = 0;

for (var batch = 0; batch < __COUNT__; batch += batchSize) {
    var end = Math.min(batch + batchSize, __COUNT__);
    var docs = [];

    for (var i = batch; i < end; i++) {
        var n = i + 1;
        var id = ObjectId((__START_HEX__ + i).toString(16).padStart(24, '0'));
        var paddedN = n.toString().padStart(5, '0');
        var username = 'bench-mqtt-device-' + paddedN;
        var password = 'bench-mqtt-secret-' + paddedN;

        docs.push({
            _id: id,
            name: 'bench-mqtt-asset-' + n,
            enabled: true,
            debugEnabled: false,
            description: 'Benchmark MQTT asset ' + n + ' — auth callout benchmarks',
            assetUUID: username,
            assetTemplateId: templateId,
            orgId: orgId,
            pathKey: __PATHKEY__,
            customerId: customerId,
            routeGroupIds: [],
            protocol: {
                type: 'mqtt',
                mqtt: {
                    username: username,
                    password: password
                }
            },
            latitude: null,
            longitude: null,
            created: now,
            updated: now
        });
    }

    try {
        var result = db.assets.insertMany(docs, { ordered: false });
        totalInserted += Object.keys(result.insertedIds).length;
    } catch (e) {
        if (e.code === 11000 || (e.writeErrors && e.writeErrors.length > 0)) {
            totalSkipped += (end - batch);
        } else {
            throw e;
        }
    }
}

if (totalSkipped > 0) {
    print('Inserted ' + totalInserted + ' MQTT assets (' + totalSkipped + ' already existed, idempotent re-run)');
} else {
    print('Inserted ' + totalInserted + ' MQTT assets');
}
