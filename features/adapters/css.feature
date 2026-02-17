Feature: CSS mapper and generator
  As a theme developer
  I need CSS output generation
  So that themes can be used in web applications

  Scenario: CSS mapper produces custom properties
    Given the Tokyo Night Dark palette from testdata
    And I derive tokens from the Tokyo Night Dark palette
    And I create a ResolvedTheme from Tokyo Night Dark
    When I map it with the CSS mapper
    Then the CssTheme should have custom properties

  Scenario: CSS mapper produces element rules
    Given the Tokyo Night Dark palette from testdata
    And I derive tokens from the Tokyo Night Dark palette
    And I create a ResolvedTheme from Tokyo Night Dark
    When I map it with the CSS mapper
    Then the CssTheme should have element rules

  Scenario: CSS generator produces valid CSS output
    Given the Tokyo Night Dark palette from testdata
    And I derive tokens from the Tokyo Night Dark palette
    And I create a ResolvedTheme from Tokyo Night Dark
    And I map it with the CSS mapper
    When I generate CSS output
    Then the generated output should contain ":root"
    And the generated output should contain "--"
