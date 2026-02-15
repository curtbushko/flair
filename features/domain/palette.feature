Feature: Base24 palette entity
  As a theme developer
  I need a palette entity that holds 24 color slots
  So that I can access colors by name or index

  Scenario: Full palette has 24 colors
    Given the Tokyo Night Dark palette from testdata
    Then the palette should have 24 colors

  Scenario: Access base00 slot
    Given the Tokyo Night Dark palette from testdata
    Then slot "base00" should be "#1a1b26"

  Scenario: Access base0D slot (blue)
    Given the Tokyo Night Dark palette from testdata
    Then slot "base0D" should be "#7aa2f7"

  Scenario: Access base0E slot (magenta)
    Given the Tokyo Night Dark palette from testdata
    Then slot "base0E" should be "#bb9af7"

  Scenario: Base index matches slot name
    Given the Tokyo Night Dark palette from testdata
    Then Base(0) should return the same as Slot("base00")
    And Base(13) should return the same as Slot("base0D")

  Scenario: Base16 fallback for base10 slot
    Given a base16 palette with only 16 colors
    Then base16 should be a fallback from base0

  Scenario: Base16 fallback for base12 slot
    Given a base16 palette with only 16 colors
    Then base18 should be a fallback from base8
