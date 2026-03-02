Feature: Reader/Writer wrappers
  As the flair pipeline
  I need versioned writers and validating readers
  So that schema versions are handled consistently

  Scenario: VersionedWriter prepends schema header
    Given a VersionedWriter for kind "tokens" and theme "tokyonight"
    When I write "tokens:\n  syntax.keyword: '#bb9af7'"
    Then the output should start with "schema_version: 1"
    And the output should contain "kind: tokens"
    And the output should contain "theme_name: tokyonight"

  Scenario: VersionedWriter includes correct version for palette
    Given a VersionedWriter for kind "palette" and theme "gruvbox"
    When I write "palette data"
    Then the output should start with "schema_version: 1"
    And the output should contain "kind: palette"

  Scenario: VersionedWriter includes correct version for vim-mapping
    Given a VersionedWriter for kind "vim-mapping" and theme "catppuccin"
    When I write "highlights:"
    Then the output should start with "schema_version: 1"
    And the output should contain "kind: vim-mapping"

  Scenario: ValidatingReader passes valid schema version
    Given YAML with schema_version 1 for kind "tokens"
    When I wrap it in ValidatingReader and read
    Then reading should succeed

  Scenario: ValidatingReader rejects outdated schema version
    Given YAML with schema_version 0 for kind "tokens"
    When I wrap it in ValidatingReader and read
    Then reading should fail with SchemaVersionError
    And NeedsUpgrade should be false

  Scenario: ValidatingReader rejects future schema version
    Given YAML with schema_version 99 for kind "tokens"
    When I wrap it in ValidatingReader and read
    Then reading should fail with SchemaVersionError
    And NeedsUpgrade should be true

  Scenario: ValidatingReader works with vim-mapping
    Given YAML with schema_version 1 for kind "vim-mapping"
    When I wrap it in ValidatingReader and read
    Then reading should succeed

  Scenario: ValidatingReader works with stylix-mapping
    Given YAML with schema_version 1 for kind "stylix-mapping"
    When I wrap it in ValidatingReader and read
    Then reading should succeed
