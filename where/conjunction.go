// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
// Package where provides functionality for building SQL WHERE clauses
package where

// conjunction represents a logical operator used to combine conditions in SQL WHERE clauses
type conjunction string

// Conjunction constants define the supported logical operators for combining conditions
const (
	conAnd conjunction = " AND " // Logical AND operator for combining conditions
	conOr  conjunction = " OR "  // Logical OR operator for combining conditions
)

