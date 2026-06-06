package mongo

// CollectionName is the Mongo collection that holds platform KEK records.
// One document per context; enforced by a unique index ensured at adapter
// construction time. Seeded by the mongodb-init container (kek-bootstrap.sh).
const CollectionName = "encryptionKeys"
