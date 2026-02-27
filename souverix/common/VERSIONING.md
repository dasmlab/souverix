# Versioning Contract

## N+1 and N-1 Compatibility

The Souverix Common Library maintains strict versioning contracts to ensure components can evolve independently.

### Compatibility Rules

1. **N+1 Support**: Components can use a newer version of common library
   - New features are additive only
   - Existing APIs remain unchanged
   - Backward compatible extensions allowed

2. **N-1 Support**: Components can use an older version of common library
   - Older components continue to work
   - No breaking changes in patch/minor versions
   - Major version bumps require explicit migration

3. **Extension Support**: Components can extend common functionality locally
   - Component-specific extensions don't break others
   - Data payloads remain compatible
   - Extensions are additive, not replacements

### Versioning Strategy

- **Major (v1.0.0)**: Breaking changes, requires migration
- **Minor (v1.1.0)**: New features, backward compatible
- **Patch (v1.0.1)**: Bug fixes, backward compatible

### Migration Path

When upgrading:
1. Test with N+1 version locally
2. Update component go.mod
3. Remove local replace directive
4. Verify all tests pass
5. Deploy

### Extension Pattern

```go
// Component can extend common functionality
type ExtendedDiagnostics struct {
    *diagnostics.Diagnostics
    // Component-specific extensions
}
```

This pattern ensures:
- Common contract remains intact
- Extensions don't break other components
- Data compatibility maintained
