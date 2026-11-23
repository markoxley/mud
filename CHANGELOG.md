# Changelog

## Unreleased

## v0.0.14 - 2026-11-23

## v0.0.13 - 2026-11-23

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

- **Column types are now dialect-specific.** Previously every database got
  the same declarations. SQL Server now emits `BIT`, `FLOAT`, `REAL`,
  `NVARCHAR`, and `DATETIME2`; MySQL emits `TINYINT(1)`, `FLOAT`, `CHAR(36)`;
  SQLite uses `INTEGER`/`REAL`/`VARCHAR` affinity types.

- **`UNSIGNED` is MySQL-only** and emitted only on numeric columns
  (`INT`/`BIGINT`/`DECIMAL`/`FLOAT`/`DOUBLE`). Previously it was appended
  unconditionally — a syntax error in T-SQL and meaningless in SQLite.
  On SQL Server, unsigned fields are **widened** instead: `uint`/`uint32` →
  `BIGINT`, `uint64` → `DECIMAL(20,0)`.

- **Size suffixes are dialect-aware.** `size:N` no longer produces
  `INT(N)` on SQL Server (invalid T-SQL); precision is only emitted for
  `DECIMAL` and character types there. MySQL keeps `INT(N)` display width.

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

- **`Refresh` could return the wrong model.** It passed `*string` to
  `where.Equal`, which is not a bindable argument, so the `ID = ?` clause was
  silently dropped and the query returned the first row instead.

- **Timestamps without fractional seconds failed to populate.**
  `utils.SQLToTime` returned `false` for `YYYY-MM-DD HH:MM:SS` because an empty
  nanoseconds string was passed to `strconv.Atoi`; every `time.Time` model
  field stayed zero on fetch. Missing nanoseconds now default to `0`.

- **UUID string criteria are detected case-insensitively.** `uuid.NewV4()`
  emits lowercase but the detection regex only matched uppercase, so a UUID
  string was treated as raw SQL and the query failed. The generated clause
  also targets `ID` rather than `id`.

- **`where.In` accepts pointers.** A `*[]int` now expands into per-element
  placeholders instead of producing an empty clause; a nil pointer yields an
  empty clause (which is dropped from the builder).

- **Fields of unexported embedded structs are modeled.** `getDefs` skipped
  anonymous fields of unexported struct types entirely, silently dropping
  their promoted exported fields. It now recurses into embedded structs
  regardless of the embedded type's visibility.

- **Ambiguous promoted fields no longer crash Save/Update.** Two embedded
  structs exposing the same column name made `FieldByName` return an invalid
  `reflect.Value`, and `insertCommand`/`updateCommand` panicked on it. They
  now skip such fields like `setFieldValue` already did.

- **Mid-iteration query failures are no longer silent.** `populateModel` and
  `RawSelect` now check `Rows.Err()` after iterating, so a context
  cancellation (or driver error) partway through a result set surfaces as an
  error instead of returning truncated data as success.

### New APIs

- **`Tabler` interface** — implement `TableName() string` to override the
  table name (documented in GUIDE.md).
- **`Model.TableName()` / `Model.SetTableName(name)`** — programmatic
  per-instance override; the private `tableName` field is now functional.
- **`utils.MakeArg`** — bind-safe value conversion used by the
  parameterized paths (bool → 1/0, `time.Time` → SQL text).
- **`ErrNoResults.Unwrap`** — enables `errors.Is`/`errors.As` chains.
- **`DB.WithContext(ctx)`** — returns a `*DB` clone whose queries honour the
  context for cancellation and deadlines (`QueryContext`/`ExecContext`/
  `BeginTx`). The clone shares the parent's pool, schema caches and locks, so
  it can be passed to any function taking `*DB`, including the generic
  helpers. A nil context behaves like `context.Background()`.

### Schema/DDL note

New tables now emit `PRIMARY KEY` for the identity column and honor
`size:N,M` precision. Databases created before this release keep their
existing schema — Mud does not run migrations, so recreate or `ALTER`
tables manually if you need the new definitions.
