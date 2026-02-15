Feature: Color representation and parsing
  As a theme developer
  I need reliable color parsing and formatting
  So that color values are handled consistently throughout the pipeline

  Scenario: Parse a 6-digit hex color with hash
    Given the hex string "#7aa2f7"
    When I parse it as a Color
    Then the RGB values should be R=122 G=162 B=247

  Scenario: Parse a 6-digit hex color without hash
    Given the hex string "7aa2f7"
    When I parse it as a Color
    Then the RGB values should be R=122 G=162 B=247

  Scenario: Parse 3-digit shorthand
    Given the hex string "#f00"
    When I parse it as a Color
    Then the RGB values should be R=255 G=0 B=0

  Scenario: Parse white color
    Given the hex string "#ffffff"
    When I parse it as a Color
    Then the RGB values should be R=255 G=255 B=255

  Scenario: Parse black color
    Given the hex string "#000000"
    When I parse it as a Color
    Then the RGB values should be R=0 G=0 B=0

  Scenario: Reject invalid hex characters
    Given the hex string "#zzzzzz"
    When I parse it as a Color
    Then parsing should fail with a ParseError

  Scenario: Reject wrong length
    Given the hex string "#12345"
    When I parse it as a Color
    Then parsing should fail with a ParseError

  Scenario: Format color as hex
    Given the hex string "#7aa2f7"
    When I parse it as a Color
    Then the color formatted as hex should be "#7aa2f7"

  Scenario: RGB to HSL conversion for red
    Given the hex string "#ff0000"
    When I parse it as a Color
    And I convert it to HSL
    Then the HSL values should be H=0 S=1.0 L=0.5

  Scenario: HSL round-trip
    Given the hex string "#bb9af7"
    When I parse it as a Color
    And I convert it to HSL
    And I convert the HSL back to RGB
    Then the color formatted as hex should be "#bb9af7"

  Scenario: NONE color sentinel
    Given a NONE color
    Then IsNone should be true

  Scenario: Luminance of white
    Given the hex string "#ffffff"
    When I parse it as a Color
    Then the luminance should be approximately 1.0

  Scenario: Luminance of black
    Given the hex string "#000000"
    When I parse it as a Color
    Then the luminance should be approximately 0.0
