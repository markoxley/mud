// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
// Package dtorm provides interfaces for database operations in the ORM.
package dtorm

// Updater defines the interface for objects that can be updated in the database.
type Updater interface {
	// Update generates an update query for the object using the provided manager.
	// Returns the generated query string and any error encountered.
	Update(mgr Manager) (string, error)
}

// Restorer defines the interface for objects that can be restored from the database.
type Restorer interface {
	// Restore generates a restore query for the object using the provided manager.
	// Returns the generated query string and any error encountered.
	Restore(mgr Manager) (string, error)
}

// Remover defines the interface for objects that can be removed from the database.
type Remover interface{}

