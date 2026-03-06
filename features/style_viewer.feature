Feature: Style Viewer TUI
  As a user
  I want to browse themes in an interactive viewer
  So that I can preview and select themes visually

  Background:
    Given a flair configuration directory exists
    And multiple themes are available

  # Viewer Launch (6.3, 6.4)

  Scenario: Launch style viewer with no arguments
    When I run "flair select" without arguments
    Then the style viewer TUI should launch
    And I should see a list of available themes

  Scenario: Launch viewer with --viewer flag
    When I run "flair select --viewer"
    Then the style viewer TUI should launch

  Scenario: Launch viewer with theme pre-selected
    When I run "flair select --viewer tokyo-night"
    Then the style viewer TUI should launch
    And "tokyo-night" should be highlighted in the theme list

  Scenario: Select theme directly without viewer
    When I run "flair select tokyo-night"
    Then the theme should be switched to "tokyo-night"
    And symlinks should be updated
    And the viewer should not launch

  # Theme Selection (6.3c, 6.3g)

  Scenario: Browse available themes
    Given the style viewer is running
    When I view the theme list
    Then I should see all available themes
    And the currently selected theme should be marked

  Scenario: Navigate theme list with keyboard
    Given the style viewer is running
    When I press "j" to move down
    Then the next theme should be highlighted
    When I press "k" to move up
    Then the previous theme should be highlighted

  Scenario: Select theme with Enter key
    Given the style viewer is running
    And I have highlighted "gruvbox" in the theme list
    When I press Enter
    Then "gruvbox" should become the selected theme
    And symlinks should be updated to point to "gruvbox"

  Scenario: Live preview on theme highlight
    Given the style viewer is running
    When I highlight a different theme in the list
    Then the style showcase should update to show that theme's colors
    And the preview should be immediate without confirmation

  # Display Pages (6.3d, 6.3e, 6.3f)

  Scenario: View palette display page
    Given the style viewer is running
    When I navigate to the palette page
    Then I should see base00 through base17 colors
    And each color should show a color swatch
    And each color should show its hex value

  Scenario: View token display page
    Given the style viewer is running
    When I navigate to the token page
    Then I should see semantic tokens grouped by category
    And I should see surface tokens
    And I should see text tokens
    And I should see status tokens
    And I should see syntax tokens

  Scenario: View lipgloss component showcase
    Given the style viewer is running
    When I navigate to the components page
    Then I should see styled button examples
    And I should see styled input examples
    And I should see styled list examples
    And I should see styled table examples

  Scenario: Component labels use token names
    Given the style viewer is running
    When I navigate to the components page
    Then component labels should use flair token names
    And examples should demonstrate the token's purpose

  # Keyboard Navigation (6.3j)

  Scenario: Navigate between pages with Tab
    Given the style viewer is running
    And I am on the palette page
    When I press Tab
    Then I should move to the next page
    When I press Shift+Tab
    Then I should move back to the previous page

  Scenario: Scroll content with j/k
    Given the style viewer is running
    And the current page has scrollable content
    When I press "j"
    Then the content should scroll down
    When I press "k"
    Then the content should scroll up

  Scenario: Quit viewer with q
    Given the style viewer is running
    When I press "q"
    Then the viewer should exit
    And the terminal should return to normal mode

  Scenario: Quit viewer with Escape
    Given the style viewer is running
    When I press Escape
    Then the viewer should exit

  Scenario: Show help with ?
    Given the style viewer is running
    When I press "?"
    Then I should see a help overlay
    And the help should list all keyboard shortcuts

  # Dynamic Theme Switching (6.3g)

  Scenario: Switch theme updates all displays
    Given the style viewer is running
    And I am viewing the components page
    When I select a different theme
    Then the palette page should update with new colors
    And the token page should update with new values
    And the components page should re-render with new styles

  # Error Handling

  Scenario: Viewer handles missing themes gracefully
    Given the style viewer is running
    And a theme directory is deleted externally
    When I try to select the deleted theme
    Then I should see an error message
    And the viewer should remain functional

  Scenario: Viewer handles corrupted theme files
    Given the style viewer is running
    And a theme has invalid YAML files
    When I highlight that theme
    Then I should see an error indicator
    And other themes should still be selectable
