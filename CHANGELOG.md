# Changelog

## Unreleased

This release converts all query generation to parameterized SQL and fixes a
number of correctness issues. Several changes are **breaking** — review the
items below before upgrading.

### Breaking changes

- **Raw SQL methods take argument slices.** `RawExecute`, `RawSelect`, and
  `RawScalar` now accept `args []interface{}` which are bound as query
  parameters. Pass `nil` (or omit variadic args where applicable) for
  parameter-free statements. For SQL Server use numbered `@p1`, `@p2`, ...
  placeholders; SQLite and MySQL use `?`.

- **`Manager` interface grew.** `Placeholder(pos *int) string` is a new
  required method, and `TableExistsQuery` now returns
  `(string, []interface{})`. Any custom `Manager` implementations must be
  updated.

- **WHERE builders return SQL + args.** `where.Builder.String` was replaced by
  `Build(operators []string, ph func(*int) string) (string, []interface{})`
  (and `BuildAt` for seeded placeholder positions). `Criteria.String` /
  `WhereString` were replaced by `Build` / `WhereClause`.

- **`TableName()` is now honored.** Models implementing `TableName() string`
  (or using `Model.SetTableName`) now read/write that table instead of the
  struct-derived name. **Warning:** if a model already had a `TableName()`
  method that was previously ignored, its table name changes on upgrade —
  existing data under the old struct-named table will not be found without
  migration.

- **`Remover` interface removed.** It was an empty interface with no
  references; any code mentioning `mud.Remover` will fail to compile.

- **`where` package internals removed.** The unexported `DBInt`, `DBFloat`,
  `DBNumeric`, `DBField` constraint types and `operatorType` bitmask table
  were deleted (the whole `datatypes.go` file). No API impact unless you
  vendored and referenced them via internals.

- **`Save` can now fail on standing-data errors.** Seed errors from
  `StandingData()` were previously swallowed; they are now propagated.

- **`updateCommand` errors when nothing is updatable.** An UPDATE with no
  assignable fields now returns `errors.New("no updatable fields")` instead
  of emitting malformed SQL.

### Behavior changes (fixes)

- **Unsigned integer fields now populate.** An inverted `err != nil` check
  meant `uint*` columns were never written back onto fetched models.

- **Nullable pointer fields work.** `*int`, `*float64`, `*bool`, `*time.Time`
  fields previously panicked or silently stayed nil on fetch; the populate
  path now allocates pointers correctly.

- **`char` fields no longer panic on empty strings.**

- **`size:N,M` tag now emits `DECIMAL(N,M)`.** The decimal component was
  previously discarded during tag parsing and never rendered into DDL.

- **New tables declare `PRIMARY KEY` on the ID column** (previously only an
  index). Existing tables are unaffected.

- **Custom `identity` fields are excluded from INSERT/UPDATE** column lists.

- **Struct tags are hardened.** `getDefs` no longer panics on malformed tag
  segments (`mud:"type"`), mistypes nil pointer fields, or recurses into nil
  pointer structs; unexported fields are skipped.

- **`Fetch`/`Range` resolve table names via `GetTableName`** like every
  other path, so overrides apply consistently.

- **Table-creation race fixed.** Concurrent first-touch of a table can no
  longer observe the cached definition before `CREATE TABLE` completes.

- **`Scan`/`Columns` errors are checked** instead of silently ignored during
  model population.

### New APIs

- **`Tabler` interface** — implement `TableName() string` to override the
  table name (documented in GUIDE.md).
- **`Model.TableName()` / `Model.SetTableName(name)`** — programmatic
  per-instance override; the private `tableName` field is now functional.
- **`utils.MakeArg`** — bind-safe value conversion used by the
  parameterized paths (bool → 1/0, `time.Time` → SQL text).
- **`ErrNoResults.Unwrap`** — enables `errors.Is`/`errors.As` chains.

### Schema/DDL note

New tables now emit `PRIMARY KEY` for the identity column and honor
`size:N,M` precision. Databases created before this release keep their
existing schema — Mud does not run migrations, so recreate or `ALTER`
tables manually if you need the new definitions.
