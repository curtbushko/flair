Feature: Token derivation and theme generation
  As a theme developer
  I need to derive semantic tokens and generate output files
  So that I can create consistent themes across targets

  Scenario: Derive tokens from Tokyo Night Dark
    Given the Tokyo Night Dark palette from testdata
    When I derive tokens from the Tokyo Night Dark palette
    Then the TokenSet should have at least 87 tokens

  Scenario: Surface tokens have correct values
    Given the Tokyo Night Dark palette from testdata
    When I derive tokens from the Tokyo Night Dark palette
    Then token "surface.background" should have color "#1a1b26"

  Scenario: Syntax keyword token
    Given the Tokyo Night Dark palette from testdata
    When I derive tokens from the Tokyo Night Dark palette
    Then token "syntax.keyword" should have color "#bb9af7"

  Scenario: Syntax comment is italic
    Given the Tokyo Night Dark palette from testdata
    When I derive tokens from the Tokyo Night Dark palette
    Then token "syntax.comment" should be italic

  Scenario: Markup heading is bold
    Given the Tokyo Night Dark palette from testdata
    When I derive tokens from the Tokyo Night Dark palette
    Then token "markup.heading" should be bold

  Scenario: Stylix mapper produces at least 60 values
    Given the Tokyo Night Dark palette from testdata
    And I derive tokens from the Tokyo Night Dark palette
    And I create a ResolvedTheme from Tokyo Night Dark
    When I map it with the Stylix mapper
    Then the StylixTheme should have at least 60 values

  Scenario: Vim mapper produces at least 200 highlight groups
    Given the Tokyo Night Dark palette from testdata
    And I derive tokens from the Tokyo Night Dark palette
    And I create a ResolvedTheme from Tokyo Night Dark
    When I map it with the Vim mapper
    Then the VimTheme should have at least 200 highlight groups
    And the VimTheme should have 16 terminal colors

  Scenario: Stylix generator produces valid JSON
    Given the Tokyo Night Dark palette from testdata
    And I derive tokens from the Tokyo Night Dark palette
    And I create a ResolvedTheme from Tokyo Night Dark
    And I map it with the Stylix mapper
    When I generate Stylix output
    Then the output should be valid JSON
    And the JSON should contain key "base00"

  Scenario: Vim generator produces Lua output
    Given the Tokyo Night Dark palette from testdata
    And I derive tokens from the Tokyo Night Dark palette
    And I create a ResolvedTheme from Tokyo Night Dark
    And I map it with the Vim mapper
    When I generate Vim output
    Then the generated output should contain "nvim_set_hl"
    And the generated output should contain "vim.cmd"
