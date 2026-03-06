Feature: Token Overrides
  As a user
  I want to override specific semantic tokens in my palette
  So that I can customize my theme without editing intermediate files

  Background:
    Given a flair configuration directory exists
    And a theme "tokyo-night" exists with a valid palette.yaml

  # Parsing Overrides (8.2)

  Scenario: Parse palette with no overrides
    Given a palette.yaml without an overrides section
    When I parse the palette
    Then the palette should have an empty overrides map
    And tokenization should proceed normally

  Scenario: Parse palette with color override
    Given a palette.yaml with the following overrides:
      """
      overrides:
        syntax.keyword:
          color: "#ff00ff"
      """
    When I parse the palette
    Then the override for "syntax.keyword" should have color "#ff00ff"

  Scenario: Parse palette with style overrides
    Given a palette.yaml with the following overrides:
      """
      overrides:
        syntax.comment:
          color: "#666666"
          italic: true
          bold: false
      """
    When I parse the palette
    Then the override for "syntax.comment" should have color "#666666"
    And the override for "syntax.comment" should have italic true
    And the override for "syntax.comment" should have bold false

  Scenario: Parse palette with multiple overrides
    Given a palette.yaml with overrides for:
      | token              | color     | bold  | italic |
      | syntax.keyword     | #ff00ff   | false | false  |
      | syntax.comment     | #666666   | false | true   |
      | status.error       | #ff0000   | true  | false  |
      | surface.background | #000000   | false | false  |
    When I parse the palette
    Then all 4 overrides should be parsed correctly

  Scenario: Invalid override color format
    Given a palette.yaml with override color "not-a-color"
    When I parse the palette
    Then I should receive a validation error
    And the error should mention invalid color format

  Scenario: Unknown token path in override
    Given a palette.yaml with override for "invalid.token.path"
    When I parse the palette
    Then I should receive a validation warning
    And the warning should mention unknown token path

  # Applying Overrides (8.3)

  Scenario: Apply color override during tokenization
    Given a palette with override:
      """
      syntax.keyword:
        color: "#ff00ff"
      """
    When I tokenize the palette
    Then the "syntax.keyword" token should have color "#ff00ff"
    And other tokens should have their derived values

  Scenario: Apply style override during tokenization
    Given a palette with override:
      """
      syntax.function:
        italic: true
        underline: true
      """
    When I tokenize the palette
    Then the "syntax.function" token should have italic true
    And the "syntax.function" token should have underline true
    And the "syntax.function" token should retain its derived color

  Scenario: Override takes precedence over derivation
    Given a palette with base24 colors
    And the default derivation would set "syntax.keyword" to "#7aa2f7"
    And an override sets "syntax.keyword" to "#ff00ff"
    When I tokenize the palette
    Then the "syntax.keyword" token should have color "#ff00ff"

  Scenario: Partial override merges with derivation
    Given a palette with base24 colors
    And the default derivation would set "syntax.string" with color "#9ece6a"
    And an override sets only italic for "syntax.string"
    When I tokenize the palette
    Then the "syntax.string" token should have color "#9ece6a"
    And the "syntax.string" token should have italic true

  Scenario: Override surface token
    Given a palette with override:
      """
      surface.background:
        color: "#000000"
      """
    When I tokenize the palette
    Then the "surface.background" token should have color "#000000"

  Scenario: Override terminal ANSI color
    Given a palette with override:
      """
      terminal.red:
        color: "#ff0000"
      """
    When I tokenize the palette
    Then the "terminal.red" token should have color "#ff0000"

  # CLI Commands (8.6)

  Scenario: Add new override via CLI
    Given a theme "tokyo-night" with no overrides
    When I run "flair override tokyo-night syntax.keyword #ff00ff"
    Then the palette.yaml should contain an override for "syntax.keyword"
    And the override color should be "#ff00ff"

  Scenario: Update existing override via CLI
    Given a theme "tokyo-night" with override "syntax.keyword: #ff00ff"
    When I run "flair override tokyo-night syntax.keyword #00ff00"
    Then the override for "syntax.keyword" should be updated to "#00ff00"

  Scenario: Add override with style options
    When I run "flair override tokyo-night syntax.comment #666666 --italic"
    Then the palette.yaml should contain an override for "syntax.comment"
    And the override should have color "#666666"
    And the override should have italic true

  Scenario: List current overrides
    Given a theme "tokyo-night" with overrides:
      | token          | color   |
      | syntax.keyword | #ff00ff |
      | syntax.comment | #666666 |
    When I run "flair override tokyo-night --list"
    Then I should see "syntax.keyword" with color "#ff00ff"
    And I should see "syntax.comment" with color "#666666"

  Scenario: List overrides when none exist
    Given a theme "tokyo-night" with no overrides
    When I run "flair override tokyo-night --list"
    Then I should see a message indicating no overrides are set

  Scenario: Remove override via CLI
    Given a theme "tokyo-night" with override "syntax.keyword: #ff00ff"
    When I run "flair override tokyo-night --remove syntax.keyword"
    Then the palette.yaml should not contain an override for "syntax.keyword"

  Scenario: Remove non-existent override
    Given a theme "tokyo-night" with no overrides
    When I run "flair override tokyo-night --remove syntax.keyword"
    Then I should see a message indicating the override does not exist

  Scenario: Override invalid token path via CLI
    When I run "flair override tokyo-night invalid.path #ff00ff"
    Then I should see an error about invalid token path
    And the palette.yaml should not be modified

  Scenario: Override with invalid color via CLI
    When I run "flair override tokyo-night syntax.keyword not-a-color"
    Then I should see an error about invalid color format
    And the palette.yaml should not be modified

  # Validation (8.7)

  Scenario: Validate command checks overrides
    Given a theme with valid overrides
    When I run "flair validate tokyo-night"
    Then the validation should pass
    And I should see override information in the output

  Scenario: Validate warns about shadowed values
    Given a theme with override for "syntax.keyword"
    When I run "flair validate tokyo-night"
    Then I should see a warning about the override shadowing derived value

  Scenario: Validate fails for invalid override color
    Given a palette.yaml with override color "invalid"
    When I run "flair validate tokyo-night"
    Then the validation should fail
    And the error should mention invalid override color

  # Regeneration (8.5b)

  Scenario: Overrides preserved during regeneration
    Given a theme "tokyo-night" with overrides
    When I run "flair regenerate tokyo-night"
    Then the overrides in palette.yaml should be preserved
    And the generated files should reflect the overrides

  Scenario: Preview shows overridden tokens
    Given a theme "tokyo-night" with override "syntax.keyword: #ff00ff"
    When I run "flair preview tokyo-night"
    Then the preview should show "syntax.keyword" with color "#ff00ff"
    And overridden tokens should be visually indicated

  # Round-trip (8.4c)

  Scenario: Parse and write preserves overrides
    Given a palette.yaml with overrides
    When I parse the palette
    And I write the palette back to YAML
    And I parse the written file
    Then the overrides should be identical to the original

  Scenario: Generate from palette with overrides
    Given a palette.yaml with override "syntax.keyword: #ff00ff"
    When I run "flair generate tokyo-night"
    Then the generated style.lua should use "#ff00ff" for keyword highlights
    And the generated universal.yaml should reflect the override
