Feature: End-to-end pipeline
  As a theme developer
  I need the full pipeline to work correctly
  So that I can generate complete themes from palettes

  Scenario: Tokyo Night Dark full pipeline produces all 12 files
    When I run the full pipeline for "tokyo-night-dark"
    Then all 12 files should be created

  Scenario: Gruvbox Dark full pipeline produces all 12 files
    When I run the full pipeline for "gruvbox-dark"
    Then all 12 files should be created

  Scenario: Catppuccin Mocha full pipeline produces all 12 files
    When I run the full pipeline for "catppuccin-mocha"
    Then all 12 files should be created

  Scenario: Pipeline produces deterministic output
    When I run the full pipeline for "tokyo-night-dark"
    Then running the pipeline again should produce identical output
