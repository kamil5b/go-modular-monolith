// MongoDB migration for auth collections
// Run this in MongoDB shell or use mongosh

// Create auth_credentials collection with indexes
db.createCollection("auth_credentials");
db.auth_credentials.createIndex({ "username": 1 }, { unique: true, partialFilterExpression: { deleted_at: null } });
db.auth_credentials.createIndex({ "email": 1 }, { unique: true, partialFilterExpression: { deleted_at: null } });
db.auth_credentials.createIndex({ "user_id": 1 }, { partialFilterExpression: { deleted_at: null } });

// Create auth_sessions collection with indexes
db.createCollection("auth_sessions");
db.auth_sessions.createIndex({ "token": 1 }, { unique: true, partialFilterExpression: { revoked_at: null } });
db.auth_sessions.createIndex({ "user_id": 1 }, { partialFilterExpression: { revoked_at: null } });
db.auth_sessions.createIndex({ "expires_at": 1 });

// TTL index to automatically delete expired sessions after 30 days
db.auth_sessions.createIndex({ "expires_at": 1 }, { expireAfterSeconds: 2592000 });

print("Auth collections and indexes created successfully");
