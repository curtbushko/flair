Feature: Color blending and manipulation
  As a theme developer
  I need color operations like blending and lightening
  So that I can derive consistent color variations

  Scenario: Blend 50% black and white
    Given two colors "#000000" and "#ffffff"
    When I blend them with ratio 0.5
    Then the result should be approximately "#808080"

  Scenario: Blend 0.0 returns source
    Given two colors "#ff0000" and "#00ff00"
    When I blend them with ratio 0.0
    Then the result should be approximately "#ff0000"

  Scenario: Blend 1.0 returns target
    Given two colors "#ff0000" and "#00ff00"
    When I blend them with ratio 1.0
    Then the result should be approximately "#00ff00"

  Scenario: Blend 25% creates subtle mix
    Given two colors "#1a1b26" and "#7aa2f7"
    When I blend them with ratio 0.25
    Then the result should be approximately "#323d5a"

  Scenario: Lighten a dark color
    Given the color "#1a1b26"
    When I lighten it by 0.1
    Then the luminance should be approximately 0.032

  Scenario: Darken a light color
    Given the color "#c0caf5"
    When I darken it by 0.1
    Then the luminance should be approximately 0.394

  Scenario: Desaturate removes color
    Given the color "#ff0000"
    When I desaturate it by 1.0
    Then I convert it to HSL
    And the HSL values should be H=0 S=0.0 L=0.5

  Scenario: Shift hue by 120 degrees
    Given the color "#ff0000"
    When I shift hue by 120 degrees
    Then the result should be approximately "#00ff00"

  Scenario: Shift hue wraps around
    Given the color "#ff0000"
    When I shift hue by 360 degrees
    Then the result should be approximately "#ff0000"
