Feature: Palette validation rules
  As a theme developer
  I need palette validation
  So that color schemes meet quality requirements

  Scenario: Valid dark palette passes validation
    Given the Tokyo Night Dark palette from testdata
    When I validate the palette
    Then validation should pass with no errors

  Scenario: Validate luminance ordering for dark theme
    Given a dark palette where base00 is lighter than base05
    When I validate the palette
    Then validation should fail with luminance ordering error

  Scenario: Validate luminance ordering for light theme
    Given a light palette where base00 is darker than base05
    When I validate the palette
    Then validation should fail with luminance ordering error

  Scenario: Validate neutral ramp monotonicity warning
    Given a palette where base01 luminance is less than base00
    When I validate the palette
    Then validation should warn about monotonicity

  Scenario: Validate bright variants are brighter
    Given a palette where base12 is darker than base08
    When I validate the palette
    Then validation should warn about bright variant luminance
