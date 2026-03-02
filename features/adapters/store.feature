Feature: ThemeStore filesystem operations
  As the flair CLI
  I need to manage theme directories and symlinks
  So that themes can be stored and selected

  Scenario: EnsureThemeDir creates directory
    Given theme "tokyonight" does not exist
    When I call EnsureThemeDir("tokyonight")
    Then the theme directory should exist

  Scenario: OpenWriter creates writable file
    Given theme "tokyonight" does not exist
    When I call EnsureThemeDir("tokyonight")
    And I call OpenWriter("tokyonight", "test.yaml") and write "hello world"
    Then FileExists("tokyonight", "test.yaml") should return true

  Scenario: OpenReader reads existing file
    Given theme "tokyonight" exists with file "tokens.yaml"
    When I call OpenWriter("tokyonight", "tokens.yaml") and write "test content"
    And I call OpenReader("tokyonight", "tokens.yaml")
    Then the content should be "test content"

  Scenario: Select creates symlinks
    Given theme "tokyonight" does not exist
    When I call EnsureThemeDir("tokyonight")
    And I call Select("tokyonight")
    Then symlink "style.lua" should point to "tokyonight/style.lua"
    And symlink "style.json" should point to "tokyonight/style.json"

  Scenario: SelectedTheme reads symlink target
    Given theme "tokyonight" does not exist
    When I call EnsureThemeDir("tokyonight")
    And I call Select("tokyonight")
    Then SelectedTheme should return "tokyonight"

  Scenario: FileExists returns false for missing file
    Given theme "tokyonight" does not exist
    When I call EnsureThemeDir("tokyonight")
    Then FileExists("tokyonight", "nonexistent.yaml") should return false
