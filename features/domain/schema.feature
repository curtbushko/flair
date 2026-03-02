Feature: Schema version management
  As a theme developer
  I need schema versioning for file compatibility
  So that file format changes can be detected and handled

  Scenario: CurrentVersion for palette
    Given file kind "palette"
    Then CurrentVersion should return 1

  Scenario: CurrentVersion for tokens
    Given file kind "tokens"
    Then CurrentVersion should return 1

  Scenario: CurrentVersion for vim-mapping
    Given file kind "vim-mapping"
    Then CurrentVersion should return 1

  Scenario: CurrentVersion for css-mapping
    Given file kind "css-mapping"
    Then CurrentVersion should return 1

  Scenario: CurrentVersion for gtk-mapping
    Given file kind "gtk-mapping"
    Then CurrentVersion should return 1

  Scenario: CurrentVersion for qss-mapping
    Given file kind "qss-mapping"
    Then CurrentVersion should return 1

  Scenario: CurrentVersion for stylix-mapping
    Given file kind "stylix-mapping"
    Then CurrentVersion should return 1

  Scenario: All FileKind constants have version greater than 0
    Given all FileKind constants
    Then each should have a version greater than 0
