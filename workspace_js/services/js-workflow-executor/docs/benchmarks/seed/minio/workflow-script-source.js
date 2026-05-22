// Benchmark workflow code node script
// Simulates a realistic workflow code node that reads event data,
// applies transformations, and produces output + statePatch.

const eventType = event.eventType || 'unknown';
const timestamp = event.timestamp || new Date().toISOString();

// Transform event data
const transformed = {
    type: eventType,
    processedAt: new Date().toISOString(),
    originalTimestamp: timestamp,
    data: {},
};

if (event.data) {
    const keys = Object.keys(event.data);
    for (let i = 0; i < keys.length; i++) {
        transformed.data[keys[i]] = event.data[keys[i]];
    }
}

// Compute some derived values (simulate real work)
let checksum = 0;
const payload = JSON.stringify(event);
for (let i = 0; i < payload.length; i++) {
    checksum = ((checksum << 5) - checksum + payload.charCodeAt(i)) | 0;
}
transformed.data._checksum = checksum;

// Produce output and state patch
const result = {
    output: transformed,
    statePatch: {
        lastProcessedEvent: eventType,
        lastProcessedAt: transformed.processedAt,
        processedCount: (state.processedCount || 0) + 1,
    },
};
