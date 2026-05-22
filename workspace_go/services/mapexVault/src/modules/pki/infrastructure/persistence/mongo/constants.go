package mongo

// CollectionName is the Mongo collection that holds platform CA records.
// Single document per Kind (root + intermediate); enforced by a unique
// index ensured at adapter construction time.
const CollectionName = "pkiCertificateAuthorities"
