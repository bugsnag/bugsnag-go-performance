Feature: Dummy mazerunner tests

Scenario: Run first scenario
  When I start the service "app"
  Then I run "HandledScenario"
  And I should receive no errors