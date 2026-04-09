# ADR 0009: read_group for database-level aggregations

## Status

Accepted

## Context

The existing `search_read` command in gindoo only supports retrieving individual records with no grouping capabilities. This forces users to:

1. Load all records (inefficient for large datasets)
2. Manually group data in post-processing (error-prone, slow)

Odoo's ORM provides a `read_group` method that performs database-level grouping using SQL `GROUP BY`, which is:
- Fast and scalable (database handles aggregation)
- Supports built-in aggregations (count, avg, sum, min, max)
- Returns only grouped metadata, not all records

Common use cases include:
- Comparing product attributes across companies
- Counting records per category
- Analyzing aggregated data (e.g., average prices)

## Decision

Add a new `read_group` command to gindoo that exposes Odoo's `read_group` functionality:

```sh
gindoo read_group <model> <domain> <fields> <groupby> [--limit N] [--orderby "field"]
```

### Command Structure

- **Positional arguments**:
  - `model`: Technical model name (e.g., `product.template`)
  - `domain`: Odoo domain filter in list syntax (e.g., `[]`)
  - `fields`: Fields to aggregate with optional aggregation syntax (e.g., `['fine_weight:avg']`)
  - `groupby`: Fields to group by (e.g., `['default_code']`)

- **Optional flags**:
  - `--limit`: Maximum number of groups to return (default: 10)
  - `--orderby`: Field to sort groups (e.g., `default_code desc`)

### Implementation Details

1. **Odoo Integration**: Uses `godoorpc.ExecuteKW()` to call `read_group` method
2. **Force database query**: Sets `lazy=false` to ensure database-level grouping
3. **Argument parsing**: Follows existing patterns with flexible flag positioning
4. **Output format**: Consistent JSON structure with input parameters and results

### Example Usage

```sh
# Group products by default_code, show average fine_weight
gindoo read_group product.template "[]" "['fine_weight:avg']" "['default_code']"

# Count companies by country with domain filter
gindoo read_group res.partner "[('is_company', '=', True)]" "['id:count']" "['country_id']" --limit 5

# Multiple aggregations with sorting
gindoo read_group product.template "[]" "['fine_weight:avg', 'list_price:max']" "['categ_id']" --orderby "fine_weight desc"
```

## Consequences

### Positive

- **Performance**: Database-level grouping is orders of magnitude faster than client-side processing
- **Scalability**: Can handle large datasets without memory issues
- **Accuracy**: Built-in aggregations avoid manual calculation errors
- **Consistency**: Follows Odoo's native API patterns
- **Flexibility**: Supports complex aggregations (avg, sum, count, min, max)
- **Agent-friendly**: Clear JSON output format for programmatic use

### Negative

- **Complexity**: Additional command increases learning curve
- **Odoo dependency**: Relies on Odoo's `read_group` implementation specifics
- **Aggregation syntax**: Requires understanding of Odoo's aggregation format (e.g., `field:avg`)
- **Limited to grouping**: Not a replacement for `search_read` but a complement

## Alternatives Considered

1. **Extend search_read**: Add grouping flags to existing command
   - Rejected: Would complicate the simple `search_read` interface
   - Different use case: grouping vs. record retrieval

2. **Client-side grouping**: Implement grouping logic in gindoo
   - Rejected: Would be slow and memory-intensive for large datasets
   - Doesn't leverage database optimization

3. **Separate aggregation commands**: Create individual commands for count, avg, etc.
   - Rejected: Would fragment functionality and limit flexibility
   - Odoo's `read_group` already provides all aggregations in one method

## Notes

The `read_group` command complements rather than replaces `search_read`. Users should:
- Use `search_read` for retrieving individual records
- Use `read_group` for grouped data and aggregations
- Use `search_count` for simple record counting

The command respects the read-only principle of gindoo and follows all established patterns for argument parsing, error handling, and output formatting.