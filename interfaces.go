// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
// Package mud provides interfaces for database operations in the ORM.
package mud

// Updater defines the interface for objects that can be updated in the database.
type Updater interface {
	// Update generates an update query for the object using the provided manager.
	// Returns the generated query string and any error encountered.
	Update(mgr Manager) error
}

// Restorer defines the interface for objects that can be restored from the database.
type Restorer interface {
	// Restore generates a restore query for the object using the provided manager.
	// Returns the generated query string and any error encountered.
	Restore(mgr Manager) error
}

// Tabler defines the interface for models that provide a custom table name.
// When implemented and non-empty, the returned name is used for all table
// operations instead of the derived struct type name.
type Tabler interface {
	TableName() string
}
