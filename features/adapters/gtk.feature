Feature: GTK mapper and generator
  As a theme developer
  I need GTK CSS output generation
  So that themes can be used in GTK applications

  Scenario: GTK mapper produces color definitions
    Given the Tokyo Night Dark palette from testdata
    And I derive tokens from the Tokyo Night Dark palette
    And I create a ResolvedTheme from Tokyo Night Dark
    When I map it with the GTK mapper
    Then the GtkTheme should have color definitions

  Scenario: GTK mapper produces widget rules
    Given the Tokyo Night Dark palette from testdata
    And I derive tokens from the Tokyo Night Dark palette
    And I create a ResolvedTheme from Tokyo Night Dark
    When I map it with the GTK mapper
    Then the GtkTheme should have widget rules

  Scenario: GTK generator produces valid GTK CSS
    Given the Tokyo Night Dark palette from testdata
    And I derive tokens from the Tokyo Night Dark palette
    And I create a ResolvedTheme from Tokyo Night Dark
    And I map it with the GTK mapper
    When I generate GTK output
    Then the generated output should contain "@define-color"
