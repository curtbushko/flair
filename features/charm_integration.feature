Feature: Charm Integration
  As a CLI developer
  I want to use flair themes with lipgloss
  So that my TUI applications have consistent, beautiful styling

  Background:
    Given a flair configuration directory exists
    And a theme "tokyo-night" is selected

  # Theme Loading (6.1)

  Scenario: Load selected theme from configuration
    When I load the current theme
    Then I should receive a Theme object
    And the theme name should be "tokyo-night"

  Scenario: Load theme when no theme is selected
    Given no theme is currently selected
    When I load the current theme
    Then I should receive an error indicating no theme is selected

  Scenario: Load theme from custom directory
    Given a custom flair directory at "/tmp/custom-flair"
    And a theme "gruvbox" exists in the custom directory
    When I load the theme from the custom directory
    Then I should receive a Theme object with name "gruvbox"

  # Color Accessors (6.1d)

  Scenario: Access surface colors
    When I load the current theme
    Then I can access surface.background color
    And I can access surface.raised color
    And I can access surface.sunken color
    And I can access surface.overlay color
    And I can access surface.popup color

  Scenario: Access text colors
    When I load the current theme
    Then I can access text.primary color
    And I can access text.secondary color
    And I can access text.muted color
    And I can access text.inverse color

  Scenario: Access status colors
    When I load the current theme
    Then I can access status.error color
    And I can access status.warning color
    And I can access status.success color
    And I can access status.info color

  Scenario: Access accent colors
    When I load the current theme
    Then I can access accent.primary color
    And I can access accent.secondary color

  # Lipgloss Style Builders (6.2)

  Scenario: Create lipgloss styles from theme
    When I create LipglossStyles from the theme
    Then I should receive a LipglossStyles object
    And the styles should be configured with theme colors

  Scenario: Surface styles are available
    When I create LipglossStyles from the theme
    Then I can access the Background style
    And I can access the Raised style
    And I can access the Sunken style
    And I can access the Overlay style
    And I can access the Popup style

  Scenario: Text styles are available
    When I create LipglossStyles from the theme
    Then I can access the Primary text style
    And I can access the Secondary text style
    And I can access the Muted text style
    And I can access the Inverse text style

  Scenario: Status styles are available
    When I create LipglossStyles from the theme
    Then I can access the Error style
    And I can access the Warning style
    And I can access the Success style
    And I can access the Info style

  Scenario: Border styles are available
    When I create LipglossStyles from the theme
    Then I can access the Default border style
    And I can access the Focus border style
    And I can access the Muted border style

  Scenario: Component styles are available
    When I create LipglossStyles from the theme
    Then I can access the Button style
    And I can access the Input style
    And I can access the List style
    And I can access the Table style
    And I can access the Dialog style

  Scenario: State styles are available
    When I create LipglossStyles from the theme
    Then I can access the Hover state style
    And I can access the Active state style
    And I can access the Disabled state style
    And I can access the Selected state style

  Scenario: Apply text style to content
    Given I have created LipglossStyles from the theme
    When I render "Hello World" with the Primary text style
    Then the output should contain ANSI color codes
    And the text color should match text.primary from the theme

  Scenario: Apply button style to content
    Given I have created LipglossStyles from the theme
    When I render "Submit" with the Button style
    Then the output should contain ANSI color codes
    And the button should have background color from the theme
    And the button should have foreground color from the theme
