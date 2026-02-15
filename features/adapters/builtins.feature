Feature: Built-in palette source
  As a theme developer
  I need access to built-in palettes
  So that I can generate themes without external files

  Scenario: List returns all embedded palette names
    When I call List() on the built-in source
    Then the result should contain "tokyo-night-dark"
    And the result should contain "gruvbox-dark"
    And the result should contain "catppuccin-mocha"
    And the result should be sorted alphabetically

  Scenario: Get returns YAML bytes for tokyo-night-dark
    When I call Get("tokyo-night-dark") on the built-in source
    Then I should receive valid YAML bytes

  Scenario: Get returns YAML bytes for gruvbox-dark
    When I call Get("gruvbox-dark") on the built-in source
    Then I should receive valid YAML bytes

  Scenario: Get returns YAML bytes for catppuccin-mocha
    When I call Get("catppuccin-mocha") on the built-in source
    Then I should receive valid YAML bytes

  Scenario: Get unknown name returns error
    When I call Get("nonexistent-theme") on the built-in source
    Then Get should return an error

  Scenario: Has returns true for tokyo-night-dark
    When I call Has("tokyo-night-dark") on the built-in source
    Then Has should return true

  Scenario: Has returns true for gruvbox-dark
    When I call Has("gruvbox-dark") on the built-in source
    Then Has should return true

  Scenario: Has returns false for unknown
    When I call Has("my-custom-theme") on the built-in source
    Then Has should return false
