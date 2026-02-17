Feature: QSS mapper and generator
  As a theme developer
  I need QSS output generation
  So that themes can be used in Qt applications

  Scenario: QSS mapper produces widget rules
    Given the Tokyo Night Dark palette from testdata
    And I derive tokens from the Tokyo Night Dark palette
    And I create a ResolvedTheme from Tokyo Night Dark
    When I map it with the QSS mapper
    Then the QssTheme should have widget rules

  Scenario: QSS generator produces valid QSS output
    Given the Tokyo Night Dark palette from testdata
    And I derive tokens from the Tokyo Night Dark palette
    And I create a ResolvedTheme from Tokyo Night Dark
    And I map it with the QSS mapper
    When I generate QSS output
    Then the generated output should contain "QWidget"
    And the generated output should contain "#"
